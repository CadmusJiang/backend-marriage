#!/bin/bash

# API测试脚本
# 演示如何使用各个API接口

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 服务地址
BASE_URL="http://localhost:9999"
API_BASE="$BASE_URL/api/v1"

echo -e "${BLUE}=== Marriage System API 测试脚本 ===${NC}"
echo ""

# 检查服务是否运行
echo -e "${YELLOW}1. 检查服务状态...${NC}"
if curl -s "$BASE_URL/api/v1/check-db" > /dev/null; then
    echo -e "${GREEN}✓ 服务正在运行${NC}"
else
    echo -e "${RED}✗ 服务未运行，请先启动服务${NC}"
    echo "启动命令: ./scripts/start_marriage_system.sh"
    exit 1
fi

echo ""

# 测试登录
echo -e "${YELLOW}2. 测试用户登录...${NC}"
LOGIN_RESPONSE=$(curl -s -X POST "$API_BASE/auth/login" \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"123456"}')

TOKEN=$(echo $LOGIN_RESPONSE | grep -o '"token":"[^"]*"' | cut -d'"' -f4)

if [ -n "$TOKEN" ]; then
    echo -e "${GREEN}✓ 登录成功，获取到token${NC}"
    echo "Token: $TOKEN"
else
    echo -e "${RED}✗ 登录失败${NC}"
    echo "响应: $LOGIN_RESPONSE"
    exit 1
fi

echo ""

# 测试获取账户列表
echo -e "${YELLOW}3. 测试获取账户列表（批量获取）...${NC}"
ACCOUNTS_RESPONSE=$(curl -s -X GET "$API_BASE/accounts?current=1&pageSize=3" \
  -H "Authorization: Bearer $TOKEN")

if echo "$ACCOUNTS_RESPONSE" | grep -q '"success":true'; then
    echo -e "${GREEN}✓ 获取账户列表成功${NC}"
    echo "账户数量: $(echo $ACCOUNTS_RESPONSE | grep -o '"total":[0-9]*' | cut -d':' -f2)"
else
    echo -e "${RED}✗ 获取账户列表失败${NC}"
    echo "响应: $ACCOUNTS_RESPONSE"
fi

echo ""

# 测试获取单个账户详情
echo -e "${YELLOW}4. 测试获取单个账户详情...${NC}"
ACCOUNT_DETAIL_RESPONSE=$(curl -s -X GET "$API_BASE/accounts/acc_001" \
  -H "Authorization: Bearer $TOKEN")

if echo "$ACCOUNT_DETAIL_RESPONSE" | grep -q '"success":true'; then
    echo -e "${GREEN}✓ 获取账户详情成功${NC}"
    USERNAME=$(echo $ACCOUNT_DETAIL_RESPONSE | grep -o '"username":"[^"]*"' | cut -d'"' -f4)
    NICKNAME=$(echo $ACCOUNT_DETAIL_RESPONSE | grep -o '"nickname":"[^"]*"' | cut -d'"' -f4)
    echo "用户名: $USERNAME, 姓名: $NICKNAME"
else
    echo -e "${RED}✗ 获取账户详情失败${NC}"
    echo "响应: $ACCOUNT_DETAIL_RESPONSE"
fi

echo ""

# 测试获取组织列表
echo -e "${YELLOW}5. 测试获取组织列表（批量获取）...${NC}"
ORGS_RESPONSE=$(curl -s -X GET "$API_BASE/org-infos?current=1&pageSize=3" \
  -H "Authorization: Bearer $TOKEN")

if echo "$ORGS_RESPONSE" | grep -q '"success":true'; then
    echo -e "${GREEN}✓ 获取组织列表成功${NC}"
    echo "组织数量: $(echo $ORGS_RESPONSE | grep -o '"total":[0-9]*' | cut -d':' -f2)"
else
    echo -e "${RED}✗ 获取组织列表失败${NC}"
    echo "响应: $ORGS_RESPONSE"
fi

echo ""

# 测试获取单个组织详情
echo -e "${YELLOW}6. 测试获取单个组织详情...${NC}"
ORG_DETAIL_RESPONSE=$(curl -s -X GET "$API_BASE/org-infos/org_001" \
  -H "Authorization: Bearer $TOKEN")

if echo "$ORG_DETAIL_RESPONSE" | grep -q '"success":true'; then
    echo -e "${GREEN}✓ 获取组织详情成功${NC}"
    ORG_NAME=$(echo $ORG_DETAIL_RESPONSE | grep -o '"orgName":"[^"]*"' | cut -d'"' -f4)
    ORG_TYPE=$(echo $ORG_DETAIL_RESPONSE | grep -o '"orgType":"[^"]*"' | cut -d'"' -f4)
    echo "组织名称: $ORG_NAME, 类型: $ORG_TYPE"
else
    echo -e "${RED}✗ 获取组织详情失败${NC}"
    echo "响应: $ORG_DETAIL_RESPONSE"
fi

echo ""

# 测试获取子组织
echo -e "${YELLOW}7. 测试获取子组织...${NC}"
CHILDREN_RESPONSE=$(curl -s -X GET "$API_BASE/org-infos/org_001/children" \
  -H "Authorization: Bearer $TOKEN")

if echo "$CHILDREN_RESPONSE" | grep -q '"success":true'; then
    echo -e "${GREEN}✓ 获取子组织成功${NC}"
    CHILDREN_COUNT=$(echo $CHILDREN_RESPONSE | grep -o '"data":\[[^]]*\]' | grep -o 'orgName' | wc -l)
    echo "子组织数量: $CHILDREN_COUNT"
else
    echo -e "${RED}✗ 获取子组织失败${NC}"
    echo "响应: $CHILDREN_RESPONSE"
fi

echo ""

# 测试获取父组织
echo -e "${YELLOW}8. 测试获取父组织...${NC}"
PARENT_RESPONSE=$(curl -s -X GET "$API_BASE/org-infos/org_006/parent" \
  -H "Authorization: Bearer $TOKEN")

if echo "$PARENT_RESPONSE" | grep -q '"success":true'; then
    echo -e "${GREEN}✓ 获取父组织成功${NC}"
    PARENT_NAME=$(echo $PARENT_RESPONSE | grep -o '"orgName":"[^"]*"' | cut -d'"' -f4)
    echo "父组织名称: $PARENT_NAME"
else
    echo -e "${RED}✗ 获取父组织失败${NC}"
    echo "响应: $PARENT_RESPONSE"
fi

echo ""

# 测试数据库检查
echo -e "${YELLOW}9. 测试数据库检查...${NC}"
DB_CHECK_RESPONSE=$(curl -s -X GET "$API_BASE/check-db")

if echo "$DB_CHECK_RESPONSE" | grep -q '"account"'; then
    echo -e "${GREEN}✓ 数据库检查成功${NC}"
    ACCOUNT_COUNT=$(echo $DB_CHECK_RESPONSE | grep -o '"account".*"count":[0-9]*' | grep -o '"count":[0-9]*' | cut -d':' -f2)
    ORG_COUNT=$(echo $DB_CHECK_RESPONSE | grep -o '"org".*"count":[0-9]*' | grep -o '"count":[0-9]*' | cut -d':' -f2)
    echo "账户表记录数: $ACCOUNT_COUNT"
    echo "组织表记录数: $ORG_COUNT"
else
    echo -e "${RED}✗ 数据库检查失败${NC}"
    echo "响应: $DB_CHECK_RESPONSE"
fi

echo ""
echo -e "${GREEN}=== API测试完成 ===${NC}"
echo ""
echo -e "${BLUE}可用的API接口：${NC}"
echo "1. 认证相关："
echo "   - POST $API_BASE/auth/login (用户登录)"
echo "   - POST $API_BASE/auth/logout (退出登录)"
echo ""
echo "2. 账户管理："
echo "   - GET $API_BASE/accounts (获取账户列表)"
echo "   - GET $API_BASE/accounts/{accountId} (获取单个账户)"
echo "   - POST $API_BASE/accounts (创建账户)"
echo "   - PUT $API_BASE/accounts/{accountId} (更新账户)"
echo "   - GET $API_BASE/account-histories (获取历史记录)"
echo ""
echo "3. 组织管理："
echo "   - GET $API_BASE/org-infos (获取组织列表)"
echo "   - GET $API_BASE/org-infos/{orgId} (获取单个组织)"
echo "   - POST $API_BASE/org-infos (创建组织)"
echo "   - PUT $API_BASE/org-infos/{orgId} (更新组织)"
echo "   - DELETE $API_BASE/org-infos/{orgId} (删除组织)"
echo "   - GET $API_BASE/org-infos/{orgId}/children (获取子组织)"
echo "   - GET $API_BASE/org-infos/{orgId}/parent (获取父组织)"
echo ""
echo "4. 系统检查："
echo "   - GET $API_BASE/check-db (检查数据库状态)"
echo ""
echo -e "${BLUE}Swagger文档：${NC}"
echo "http://localhost:8000/swagger/index.html" 