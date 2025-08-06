-- 检查账户和组织数据的关联关系

-- 1. 查看组织表数据
SELECT '=== 组织表数据 ===' as info;
SELECT id, org_id, org_name, org_type, org_level FROM org ORDER BY id;

-- 2. 查看账户表数据
SELECT '=== 账户表数据 ===' as info;
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

-- 3. 检查关联关系
SELECT '=== 账户-组织关联检查 ===' as info;
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

-- 4. 检查是否有未关联的组织
SELECT '=== 未关联的组织 ===' as info;
SELECT 
    o.id,
    o.org_id,
    o.org_name,
    o.org_type
FROM org o
WHERE o.id NOT IN (
    SELECT DISTINCT belong_group_id FROM account WHERE belong_group_id IS NOT NULL
    UNION
    SELECT DISTINCT belong_team_id FROM account WHERE belong_team_id IS NOT NULL
);

-- 5. 检查账户表中的组织ID是否有效
SELECT '=== 无效的组织ID ===' as info;
SELECT 
    a.account_id,
    a.username,
    a.belong_group_id,
    a.belong_team_id
FROM account a
WHERE (a.belong_group_id IS NOT NULL AND a.belong_group_id NOT IN (SELECT id FROM org))
   OR (a.belong_team_id IS NOT NULL AND a.belong_team_id NOT IN (SELECT id FROM org)); 