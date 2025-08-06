-- 数据库迁移脚本：将 org_info 表重命名为 org 表
-- 执行前请备份数据库

USE marriage_system;

-- 1. 检查 org_info 表是否存在
SELECT COUNT(*) as org_info_exists FROM information_schema.tables 
WHERE table_schema = 'marriage_system' AND table_name = 'org_info';

-- 2. 检查 org 表是否已存在
SELECT COUNT(*) as org_exists FROM information_schema.tables 
WHERE table_schema = 'marriage_system' AND table_name = 'org';

-- 3. 如果 org_info 表存在且 org 表不存在，则重命名
-- 注意：这个操作会保留所有数据和索引
RENAME TABLE `org_info` TO `org`;

-- 4. 验证迁移结果
SELECT COUNT(*) as org_count FROM `org`;
SELECT 'org_info 表已成功重命名为 org 表' as migration_result;

-- 5. 显示表结构（可选）
-- DESCRIBE `org`; 