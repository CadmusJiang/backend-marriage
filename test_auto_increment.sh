#!/bin/bash

# 测试自增ID重置功能的脚本

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 数据库配置
DB_HOST="gz-cdb-pepynap7.sql.tencentcdb.com"
DB_PORT="63623"
DB_NAME="marriage_system"
DB_USER="marriage_system"
DB_PASS="19970901Zyt"

echo -e "${BLUE}================================${NC}"
echo -e "${BLUE}    自增ID重置功能测试${NC}"
echo -e "${BLUE}================================${NC}"
echo ""

# 检查MySQL客户端
if ! command -v mysql &> /dev/null; then
    echo -e "${RED}错误: mysql客户端未安装${NC}"
    echo -e "${YELLOW}请安装mysql-client: brew install mysql-client (macOS)${NC}"
    exit 1
fi

echo -e "${YELLOW}MySQL客户端版本:${NC}"
mysql --version
echo ""

# 测试数据库连接
echo -e "${YELLOW}测试数据库连接...${NC}"
if mysql -h"$DB_HOST" -P"$DB_PORT" -u"$DB_USER" -p"$DB_PASS" "$DB_NAME" -e "SELECT 1;" >/dev/null 2>&1; then
    echo -e "${GREEN}数据库连接成功${NC}"
elif mysql -h"$DB_HOST" -P"$DB_PORT" -u"$DB_USER" -p"$DB_PASS" "$DB_NAME" --ssl-mode=DISABLED -e "SELECT 1;" >/dev/null 2>&1; then
    echo -e "${GREEN}数据库连接成功（禁用SSL）${NC}"
    MYSQL_EXTRA_OPTS="--ssl-mode=DISABLED"
elif mysql -h"$DB_HOST" -P"$DB_PORT" -u"$DB_USER" -p"$DB_PASS" "$DB_NAME" --protocol=TCP -e "SELECT 1;" >/dev/null 2>&1; then
    echo -e "${GREEN}数据库连接成功（强制TCP）${NC}"
    MYSQL_EXTRA_OPTS="--protocol=TCP"
else
    echo -e "${RED}数据库连接失败${NC}"
    echo -e "${YELLOW}可能的原因: MySQL 9.4版本兼容性问题${NC}"
    echo -e "${BLUE}建议: 安装MySQL 8.0客户端: brew install mysql-client@8.0${NC}"
    exit 1
fi

echo ""

# 显示当前表信息
echo -e "${YELLOW}当前数据库表信息:${NC}"
mysql -h"$DB_HOST" -P"$DB_PORT" -u"$DB_USER" -p"$DB_PASS" "$DB_NAME" $MYSQL_EXTRA_OPTS -e "
    SELECT 
        TABLE_NAME as '表名',
        TABLE_ROWS as '行数',
        AUTO_INCREMENT as '当前自增值',
        CREATE_TIME as '创建时间'
    FROM information_schema.TABLES 
    WHERE TABLE_SCHEMA = '$DB_NAME' 
    ORDER BY TABLE_NAME;
" 2>/dev/null || echo -e "${YELLOW}无法获取表信息（可能是权限问题）${NC}"

echo ""

# 测试自增ID重置
echo -e "${YELLOW}测试自增ID重置功能...${NC}"
echo -e "${BLUE}使用脚本: ./cleanup_db.sh -i${NC}"
echo ""

# 执行自增ID重置
if [ -f "./cleanup_db.sh" ]; then
    echo -e "${GREEN}执行自增ID重置...${NC}"
    ./cleanup_db.sh -i
else
    echo -e "${RED}错误: 未找到 cleanup_db.sh 脚本${NC}"
    exit 1
fi

echo ""

# 显示重置后的表信息
echo -e "${YELLOW}重置后的数据库表信息:${NC}"
mysql -h"$DB_HOST" -P"$DB_PORT" -u"$DB_USER" -p"$DB_PASS" "$DB_NAME" $MYSQL_EXTRA_OPTS -e "
    SELECT 
        TABLE_NAME as '表名',
        TABLE_ROWS as '行数',
        AUTO_INCREMENT as '重置后自增值',
        CREATE_TIME as '创建时间'
    FROM information_schema.TABLES 
    WHERE TABLE_SCHEMA = '$DB_NAME' 
    ORDER BY TABLE_NAME;
" 2>/dev/null || echo -e "${YELLOW}无法获取表信息（可能是权限问题）${NC}"

echo ""
echo -e "${GREEN}================================${NC}"
echo -e "${GREEN}    测试完成！${NC}"
echo -e "${GREEN}================================${NC}"
