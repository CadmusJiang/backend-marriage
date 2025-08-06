-- 修复账户和组织数据的关联关系

-- 1. 首先查看当前的组织数据
SELECT '=== 当前组织数据 ===' as info;
SELECT id, org_id, org_name, org_type, org_level FROM org ORDER BY id;

-- 2. 查看当前的账户数据
SELECT '=== 当前账户数据 ===' as info;
SELECT 
    id,
    account_id,
    username,
    nickname,
    belong_group_id,
    belong_group_nickname,
    belong_team_id,
    belong_team_nickname
FROM account ORDER BY id;

-- 3. 修复账户表中的组织关联
-- 根据组织名称匹配来设置正确的组织ID
UPDATE account SET 
    belong_group_id = (SELECT id FROM org WHERE org_name = '系统管理组' LIMIT 1),
    belong_group_username = (SELECT org_id FROM org WHERE org_name = '系统管理组' LIMIT 1),
    belong_group_nickname = '系统管理组',
    belong_group_created_timestamp = 1705123200,
    belong_group_modified_timestamp = 1705123200,
    belong_group_current_cnt = 1
WHERE username = 'admin';

UPDATE account SET 
    belong_group_id = (SELECT id FROM org WHERE org_name = '南京-天元大厦组' LIMIT 1),
    belong_group_username = (SELECT org_id FROM org WHERE org_name = '南京-天元大厦组' LIMIT 1),
    belong_group_nickname = '南京-天元大厦组',
    belong_group_created_timestamp = 1705123200,
    belong_group_modified_timestamp = 1705123200,
    belong_group_current_cnt = 15,
    belong_team_id = (SELECT id FROM org WHERE org_name = '营销团队A' LIMIT 1),
    belong_team_username = (SELECT org_id FROM org WHERE org_name = '营销团队A' LIMIT 1),
    belong_team_nickname = '营销团队A',
    belong_team_created_timestamp = 1705123200,
    belong_team_modified_timestamp = 1705123200,
    belong_team_current_cnt = 8
WHERE username = 'company_manager';

UPDATE account SET 
    belong_group_id = (SELECT id FROM org WHERE org_name = '北京办公室组' LIMIT 1),
    belong_group_username = (SELECT org_id FROM org WHERE org_name = '北京办公室组' LIMIT 1),
    belong_group_nickname = '北京办公室组',
    belong_group_created_timestamp = 1705123200,
    belong_group_modified_timestamp = 1705123200,
    belong_group_current_cnt = 12,
    belong_team_id = (SELECT id FROM org WHERE org_name = '销售团队B' LIMIT 1),
    belong_team_username = (SELECT org_id FROM org WHERE org_name = '销售团队B' LIMIT 1),
    belong_team_nickname = '销售团队B',
    belong_team_created_timestamp = 1705123200,
    belong_team_modified_timestamp = 1705123200,
    belong_team_current_cnt = 6
WHERE username = 'group_manager';

UPDATE account SET 
    belong_group_id = (SELECT id FROM org WHERE org_name = '上海分公司组' LIMIT 1),
    belong_group_username = (SELECT org_id FROM org WHERE org_name = '上海分公司组' LIMIT 1),
    belong_group_nickname = '上海分公司组',
    belong_group_created_timestamp = 1705123200,
    belong_group_modified_timestamp = 1705123200,
    belong_group_current_cnt = 20,
    belong_team_id = (SELECT id FROM org WHERE org_name = '技术团队C' LIMIT 1),
    belong_team_username = (SELECT org_id FROM org WHERE org_name = '技术团队C' LIMIT 1),
    belong_team_nickname = '技术团队C',
    belong_team_created_timestamp = 1705123200,
    belong_team_modified_timestamp = 1705123200,
    belong_team_current_cnt = 10
WHERE username = 'team_manager';

UPDATE account SET 
    belong_group_id = (SELECT id FROM org WHERE org_name = '广州办公室组' LIMIT 1),
    belong_group_username = (SELECT org_id FROM org WHERE org_name = '广州办公室组' LIMIT 1),
    belong_group_nickname = '广州办公室组',
    belong_group_created_timestamp = 1705123200,
    belong_group_modified_timestamp = 1705123200,
    belong_group_current_cnt = 8,
    belong_team_id = (SELECT id FROM org WHERE org_name = '客服团队D' LIMIT 1),
    belong_team_username = (SELECT org_id FROM org WHERE org_name = '客服团队D' LIMIT 1),
    belong_team_nickname = '客服团队D',
    belong_team_created_timestamp = 1705123200,
    belong_team_modified_timestamp = 1705123200,
    belong_team_current_cnt = 4
WHERE username = 'employee';

-- 4. 验证修复结果
SELECT '=== 修复后的账户-组织关联 ===' as info;
SELECT 
    a.account_id,
    a.username,
    a.nickname,
    a.belong_group_id,
    g.org_name as group_name,
    a.belong_team_id,
    t.org_name as team_name
FROM account a
LEFT JOIN org g ON a.belong_group_id = g.id
LEFT JOIN org t ON a.belong_team_id = t.id
ORDER BY a.id;

-- 5. 显示修复完成信息
SELECT '账户-组织关联修复完成！' as message; 