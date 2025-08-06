#!/bin/bash

# 执行组织数据初始化脚本
# 确保每个team只属于一个组

echo "=========================================="
echo "开始初始化正确的组织数据..."
echo "=========================================="

# 使用腾讯云数据库配置
DB_HOST="gz-cdb-pepynap7.sql.tencentcdb.com"
DB_PORT="63623"
DB_USER="marriage_system"
DB_PASS="19970901Zyt"
DB_NAME="marriage_system"

echo "1. 连接到腾讯云数据库..."
echo "主机: $DB_HOST"
echo "端口: $DB_PORT"
echo "数据库: $DB_NAME"
echo "用户: $DB_USER"

echo ""
echo "2. 执行数据初始化脚本..."

# 执行SQL脚本
mysql -h $DB_HOST -P $DB_PORT -u $DB_USER -p$DB_PASS $DB_NAME < scripts/init_correct_org_data.sql

if [ $? -eq 0 ]; then
    echo ""
    echo "=========================================="
    echo "数据初始化成功！"
    echo "=========================================="
    
    echo ""
    echo "3. 验证数据..."
    
    # 验证teams接口
    echo "验证 /api/v1/teams 接口:"
    curl -s -X GET "http://localhost:9999/api/v1/teams" \
      -H "Content-Type: application/json" | jq '.'
    
    echo ""
    echo "验证 /api/v1/groups 接口:"
    curl -s -X GET "http://localhost:9999/api/v1/groups?includeTeams=true" \
      -H "Content-Type: application/json" | jq '.'
    
    echo ""
    echo "=========================================="
    echo "初始化完成！"
    echo "=========================================="
    echo ""
    echo "新的组织结构："
    echo "1. 系统管理组"
    echo "   └── 系统管理团队"
    echo ""
    echo "2. 南京-天元大厦组"
    echo "   ├── 营销团队A"
    echo "   └── 营销团队C"
    echo ""
    echo "3. 南京-南京南站组"
    echo "   └── 营销团队B"
    echo ""
    echo "✅ 每个team现在只属于一个组！"
else
    echo ""
    echo "❌ 数据初始化失败！"
    echo "请检查数据库连接和权限。"
fi 