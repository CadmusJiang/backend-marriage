#!/bin/bash

# 更新组织表名脚本：将 org_info 表重命名为 org 表
# 执行前请确保已备份数据库

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${YELLOW}=== 组织表名更新脚本 ===${NC}"
echo "此脚本将把 org_info 表重命名为 org 表"
echo ""

# 检查参数
if [ $# -eq 0 ]; then
    echo -e "${RED}错误：请提供数据库连接参数${NC}"
    echo "用法: $0 -h HOST -P PORT -u USER -pPASSWORD"
    echo "示例: $0 -h localhost -P 3306 -u root -p123456"
    exit 1
fi

# 执行迁移
echo -e "${YELLOW}正在执行数据库迁移...${NC}"
mysql "$@" < scripts/migrate_org_info_to_org.sql

if [ $? -eq 0 ]; then
    echo -e "${GREEN}✓ 数据库迁移完成${NC}"
    echo ""
    echo -e "${GREEN}更新内容：${NC}"
    echo "- 表名从 org_info 改为 org"
    echo "- 保留了所有数据和索引"
    echo "- 内部字段 org_name 保持不变"
    echo ""
    echo -e "${YELLOW}注意事项：${NC}"
    echo "- 请重新生成数据库模型代码："
    echo "  go run cmd/gormgen/main.go -structs OrgInfo -input ./internal/repository/mysql/org_info/"
    echo "- 请重新编译应用："
    echo "  go build -o bin/marriage-system main.go"
    echo ""
    echo -e "${GREEN}迁移完成！${NC}"
else
    echo -e "${RED}✗ 数据库迁移失败${NC}"
    exit 1
fi 