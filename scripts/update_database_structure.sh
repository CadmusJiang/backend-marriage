#!/bin/bash

# 更新数据库表结构脚本
# 将ID字段改为bigint unsigned，其他数字字段改为bigint

echo "开始更新数据库表结构..."

# 从配置文件读取数据库连接信息
DB_HOST=$(grep "addr" configs/dev_configs.toml | grep mysql | head -1 | cut -d'"' -f2)
DB_USER=$(grep "user" configs/dev_configs.toml | grep mysql | head -1 | cut -d'"' -f2)
DB_PASS=$(grep "pass" configs/dev_configs.toml | grep mysql | head -1 | cut -d'"' -f2)
DB_NAME=$(grep "name" configs/dev_configs.toml | grep mysql | head -1 | cut -d'"' -f2)

echo "数据库连接信息:"
echo "Host: $DB_HOST"
echo "User: $DB_USER"
echo "Database: $DB_NAME"

# 检查数据库连接
echo "检查数据库连接..."
if ! mysql -h "$DB_HOST" -u "$DB_USER" -p"$DB_PASS" -e "SELECT 1;" > /dev/null 2>&1; then
    echo "❌ 数据库连接失败"
    exit 1
fi
echo "✅ 数据库连接成功"

# 执行SQL更新脚本
echo "执行表结构更新..."
mysql -h "$DB_HOST" -u "$DB_USER" -p"$DB_PASS" "$DB_NAME" < scripts/update_table_structure.sql

if [ $? -eq 0 ]; then
    echo "✅ 表结构更新成功"
else
    echo "❌ 表结构更新失败"
    exit 1
fi

echo "验证表结构..."
curl -s http://localhost:9999/api/v1/check-db | jq '.account.columns[] | select(.Field == "id") | .Type'

echo "✅ 数据库表结构更新完成！" 