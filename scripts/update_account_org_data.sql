-- 更新账户的组织信息

-- 1. 更新admin账户的组织信息
UPDATE account SET 
    belong_group_id = 1,
    belong_group_username = 'admin_group',
    belong_group_nickname = '系统管理组',
    belong_group_created_timestamp = 1705123200,
    belong_group_modified_timestamp = 1705123200,
    belong_group_current_cnt = 1,
    belong_team_id = 1,
    belong_team_username = 'admin_team',
    belong_team_nickname = '系统管理团队',
    belong_team_created_timestamp = 1705123200,
    belong_team_modified_timestamp = 1705123200,
    belong_team_current_cnt = 1
WHERE username = 'admin';

-- 2. 更新company_manager账户的组织信息
UPDATE account SET 
    belong_group_id = 2,
    belong_group_username = 'nanjing_tianyuan',
    belong_group_nickname = '南京-天元大厦组',
    belong_group_created_timestamp = 1705123200,
    belong_group_modified_timestamp = 1705123200,
    belong_group_current_cnt = 15,
    belong_team_id = 2,
    belong_team_username = 'marketing_team_a',
    belong_team_nickname = '营销团队A',
    belong_team_created_timestamp = 1705123200,
    belong_team_modified_timestamp = 1705123200,
    belong_team_current_cnt = 8
WHERE username = 'company_manager';

-- 3. 更新group_manager账户的组织信息
UPDATE account SET 
    belong_group_id = 3,
    belong_group_username = 'beijing_office',
    belong_group_nickname = '北京办公室组',
    belong_group_created_timestamp = 1705123200,
    belong_group_modified_timestamp = 1705123200,
    belong_group_current_cnt = 12,
    belong_team_id = 3,
    belong_team_username = 'sales_team_b',
    belong_team_nickname = '销售团队B',
    belong_team_created_timestamp = 1705123200,
    belong_team_modified_timestamp = 1705123200,
    belong_team_current_cnt = 6
WHERE username = 'group_manager';

-- 4. 更新team_manager账户的组织信息
UPDATE account SET 
    belong_group_id = 4,
    belong_group_username = 'shanghai_branch',
    belong_group_nickname = '上海分公司组',
    belong_group_created_timestamp = 1705123200,
    belong_group_modified_timestamp = 1705123200,
    belong_group_current_cnt = 20,
    belong_team_id = 4,
    belong_team_username = 'tech_team_c',
    belong_team_nickname = '技术团队C',
    belong_team_created_timestamp = 1705123200,
    belong_team_modified_timestamp = 1705123200,
    belong_team_current_cnt = 10
WHERE username = 'team_manager';

-- 5. 更新group_manager2账户的组织信息
UPDATE account SET 
    belong_group_id = 5,
    belong_group_username = 'guangzhou_office',
    belong_group_nickname = '广州办公室组',
    belong_group_created_timestamp = 1705123200,
    belong_group_modified_timestamp = 1705123200,
    belong_group_current_cnt = 8,
    belong_team_id = 5,
    belong_team_username = 'service_team_d',
    belong_team_nickname = '客服团队D',
    belong_team_created_timestamp = 1705123200,
    belong_team_modified_timestamp = 1705123200,
    belong_team_current_cnt = 4
WHERE username = 'group_manager2';

-- 6. 更新employee001账户的组织信息
UPDATE account SET 
    belong_group_id = 2,
    belong_group_username = 'nanjing_tianyuan',
    belong_group_nickname = '南京-天元大厦组',
    belong_group_created_timestamp = 1705123200,
    belong_group_modified_timestamp = 1705123200,
    belong_group_current_cnt = 15,
    belong_team_id = 2,
    belong_team_username = 'marketing_team_a',
    belong_team_nickname = '营销团队A',
    belong_team_created_timestamp = 1705123200,
    belong_team_modified_timestamp = 1705123200,
    belong_team_current_cnt = 8
WHERE username = 'employee001';

-- 7. 更新employee002账户的组织信息
UPDATE account SET 
    belong_group_id = 3,
    belong_group_username = 'beijing_office',
    belong_group_nickname = '北京办公室组',
    belong_group_created_timestamp = 1705123200,
    belong_group_modified_timestamp = 1705123200,
    belong_group_current_cnt = 12,
    belong_team_id = 3,
    belong_team_username = 'sales_team_b',
    belong_team_nickname = '销售团队B',
    belong_team_created_timestamp = 1705123200,
    belong_team_modified_timestamp = 1705123200,
    belong_team_current_cnt = 6
WHERE username = 'employee002';

-- 8. 更新employee003账户的组织信息
UPDATE account SET 
    belong_group_id = 4,
    belong_group_username = 'shanghai_branch',
    belong_group_nickname = '上海分公司组',
    belong_group_created_timestamp = 1705123200,
    belong_group_modified_timestamp = 1705123200,
    belong_group_current_cnt = 20,
    belong_team_id = 4,
    belong_team_username = 'tech_team_c',
    belong_team_nickname = '技术团队C',
    belong_team_created_timestamp = 1705123200,
    belong_team_modified_timestamp = 1705123200,
    belong_team_current_cnt = 10
WHERE username = 'employee003';

-- 验证更新结果
SELECT 
    id,
    username,
    nickname,
    belong_group_id,
    belong_group_nickname,
    belong_team_id,
    belong_team_nickname
FROM account 
ORDER BY id; 