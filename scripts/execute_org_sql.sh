#!/bin/bash

# 执行org表创建脚本
echo "正在创建org表..."

# 使用腾讯云数据库连接信息
DB_HOST="gz-cdb-pepynap7.sql.tencentcdb.com"
DB_PORT="63623"
DB_USER="marriage_system"
DB_PASS="19970901Zyt"
DB_NAME="marriage_system"

# 执行SQL脚本
mysql -h $DB_HOST -P $DB_PORT -u $DB_USER -p$DB_PASS $DB_NAME < scripts/create_org_table.sql

if [ $? -eq 0 ]; then
    echo "org表创建成功！"
else
    echo "org表创建失败，请检查数据库连接和权限"
fi 