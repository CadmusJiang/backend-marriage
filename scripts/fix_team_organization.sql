-- 修复team组织关系脚本
-- 确保每个team只属于一个组

USE marriage_system;

-- 1. 备份当前数据
CREATE TABLE IF NOT EXISTS `org_backup_$(date +%Y%m%d_%H%M%S)` AS SELECT * FROM `org`;

-- 2. 查看当前数据分布
SELECT '当前数据分布:' as info;
SELECT 
    id,
    org_name,
    org_type,
    org_path,
    org_level,
    CASE 
        WHEN org_type = 1 THEN '组'
        WHEN org_type = 2 THEN '团队'
        ELSE '未知'
    END as type_desc
FROM `org` 
ORDER BY org_type, org_level, id;

-- 3. 删除所有现有数据
DELETE FROM `org`;

-- 4. 重新插入正确的组织结构
-- 组结构：
-- 1. 系统管理组 (顶级组)
-- 2. 南京-天元大厦组 (顶级组)  
-- 3. 南京-南京南站组 (顶级组)

-- 团队分配：
-- 系统管理组 -> 系统管理团队
-- 南京-天元大厦组 -> 营销团队A, 营销团队C
-- 南京-南京南站组 -> 营销团队B

INSERT INTO `org` (
    `org_name`, 
    `org_type`, 
    `org_path`, 
    `org_level`, 
    `current_cnt`, 
    `max_cnt`, 
    `status`, 
    `ext_data`, 
    `created_at`, 
    `updated_at`, 
    `created_user`, 
    `updated_user`
) VALUES 
-- 组 (org_type = 1)
('系统管理组', 1, '/1', 1, 1, 10, 1, '{"prefix": "weilan_"}', 1705123200, 1705123200, '系统管理员', '系统管理员'),
('南京-天元大厦组', 1, '/2', 1, 15, 50, 1, '{"prefix": "weilan_"}', 1705123200, 1705123200, '系统管理员', '系统管理员'),
('南京-南京南站组', 1, '/3', 1, 12, 50, 1, '{"prefix": "weilan_"}', 1705123200, 1705123200, '系统管理员', '系统管理员'),

-- 团队 (org_type = 2) - 每个团队只属于一个组
('系统管理团队', 2, '/1/1', 2, 1, 10, 1, '{"team_leader": "系统管理员", "project": "系统维护"}', 1705123200, 1705123200, '系统管理员', '系统管理员'),
('营销团队A', 2, '/2/1', 2, 8, 20, 1, '{"team_leader": "张伟", "project": "产品推广"}', 1705123200, 1705123200, '系统管理员', '系统管理员'),
('营销团队B', 2, '/3/1', 2, 10, 20, 1, '{"team_leader": "刘强", "project": "市场拓展"}', 1705123200, 1705123200, '系统管理员', '系统管理员'),
('营销团队C', 2, '/2/2', 2, 6, 20, 1, '{"team_leader": "王芳", "project": "品牌建设"}', 1705123200, 1705123200, '系统管理员', '系统管理员');

-- 5. 验证修复结果
SELECT '修复后的数据分布:' as info;
SELECT 
    id,
    org_name,
    org_type,
    org_path,
    org_level,
    CASE 
        WHEN org_type = 1 THEN '组'
        WHEN org_type = 2 THEN '团队'
        ELSE '未知'
    END as type_desc
FROM `org` 
ORDER BY org_type, org_level, id;

-- 6. 验证每个团队只属于一个组
SELECT '团队归属验证:' as info;
SELECT 
    t.id as team_id,
    t.org_name as team_name,
    g.org_name as group_name,
    t.org_path as team_path
FROM `org` t
JOIN `org` g ON t.org_path LIKE CONCAT(g.org_path, '/%')
WHERE t.org_type = 2 AND g.org_type = 1
ORDER BY t.id;

-- 7. 显示统计信息
SELECT 
    '数据统计' as info,
    COUNT(*) as total_records,
    SUM(CASE WHEN org_type = 1 THEN 1 ELSE 0 END) as groups_count,
    SUM(CASE WHEN org_type = 2 THEN 1 ELSE 0 END) as teams_count
FROM `org`; 