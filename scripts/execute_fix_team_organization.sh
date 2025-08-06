#!/bin/bash

# 修复team组织关系脚本
# 确保每个team只属于一个组

set -e

# 配置信息
DB_HOST="localhost"
DB_PORT="3306"
DB_USER="marriage_system"
DB_PASS="19970901Zyt"
DB_NAME="marriage_system"

echo "=========================================="
echo "开始修复team组织关系..."
echo "=========================================="

# 1. 备份当前数据
echo "1. 备份当前数据..."
BACKUP_FILE="org_backup_$(date +%Y%m%d_%H%M%S).sql"
mysqldump -h $DB_HOST -P $DB_PORT -u $DB_USER -p$DB_PASS $DB_NAME org > $BACKUP_FILE
echo "备份文件: $BACKUP_FILE"

# 2. 执行修复脚本
echo "2. 执行修复脚本..."
mysql -h $DB_HOST -P $DB_PORT -u $DB_USER -p$DB_PASS $DB_NAME < scripts/fix_team_organization.sql

echo "=========================================="
echo "修复完成！"
echo "=========================================="

# 3. 验证修复结果
echo "3. 验证修复结果..."
mysql -h $DB_HOST -P $DB_PORT -u $DB_USER -p$DB_PASS $DB_NAME -e "
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
FROM org 
ORDER BY org_type, org_level, id;
"

echo "=========================================="
echo "团队归属验证:"
echo "=========================================="
mysql -h $DB_HOST -P $DB_PORT -u $DB_USER -p$DB_PASS $DB_NAME -e "
SELECT 
    t.id as team_id,
    t.org_name as team_name,
    g.org_name as group_name,
    t.org_path as team_path
FROM org t
JOIN org g ON t.org_path LIKE CONCAT(g.org_path, '/%')
WHERE t.org_type = 2 AND g.org_type = 1
ORDER BY t.id;
"

echo "=========================================="
echo "修复完成！现在每个team只属于一个组。"
echo "==========================================" 