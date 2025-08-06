-- 修复账户归属规则
-- 根据角色类型设置正确的归属组和团队

-- 1. company_manager不能有归属组和团队
UPDATE account 
SET belong_group_id = NULL, belong_team_id = NULL 
WHERE role_type = 'company_manager';

-- 2. group_manager不能有归属团队，但必须有归属组
UPDATE account 
SET belong_team_id = NULL 
WHERE role_type = 'group_manager';

-- 3. team_manager必须有归属组和团队（保持不变）
-- 4. employee必须有归属组，但可以没有归属团队（保持不变）

-- 验证修复结果
SELECT 
    username,
    nickname,
    role_type,
    belong_group_id,
    belong_team_id,
    CASE 
        WHEN role_type = 'company_manager' AND belong_group_id IS NULL AND belong_team_id IS NULL THEN '✅ 正确'
        WHEN role_type = 'group_manager' AND belong_group_id IS NOT NULL AND belong_team_id IS NULL THEN '✅ 正确'
        WHEN role_type = 'team_manager' AND belong_group_id IS NOT NULL AND belong_team_id IS NOT NULL THEN '✅ 正确'
        WHEN role_type = 'employee' AND belong_group_id IS NOT NULL THEN '✅ 正确'
        ELSE '❌ 错误'
    END as rule_check
FROM account 
ORDER BY role_type, username; 