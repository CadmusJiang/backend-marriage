#!/bin/bash

# 账户管理系统API测试脚本 - 真实数据库版本
BASE_URL="http://localhost:9999"

echo "=== 账户管理系统API测试 (真实数据库) ==="

# 1. 测试登录
echo "1. 测试登录..."
LOGIN_RESPONSE=$(curl -s -X POST "$BASE_URL/api/v1/auth/login" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "admin",
    "password": "123456"
  }')

echo "登录响应: $LOGIN_RESPONSE"
TOKEN=$(echo $LOGIN_RESPONSE | grep -o '"token":"[^"]*"' | cut -d'"' -f4)
echo "获取到的Token: $TOKEN"

echo ""

# 2. 测试获取账户列表
echo "2. 测试获取账户列表..."
curl -s -X GET "$BASE_URL/api/v1/accounts?current=1&pageSize=5" \
  -H "Authorization: Bearer $TOKEN" | jq '.'

echo ""

# 3. 测试获取账户详情
echo "3. 测试获取账户详情..."
curl -s -X GET "$BASE_URL/api/v1/accounts/acc_001" \
  -H "Authorization: Bearer $TOKEN" | jq '.'

echo ""

# 4. 测试创建账户
echo "4. 测试创建账户..."
curl -s -X POST "$BASE_URL/api/v1/accounts" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "username": "testuser",
    "nickname": "测试用户",
    "password": "123456",
    "phone": "13800138099",
    "roleType": "employee",
    "status": "enabled"
  }' | jq '.'

echo ""

# 5. 测试更新账户
echo "5. 测试更新账户..."
curl -s -X PUT "$BASE_URL/api/v1/accounts/acc_001" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "nickname": "更新后的管理员",
    "phone": "13800138088",
    "status": "enabled"
  }' | jq '.'

echo ""

# 6. 测试获取账户历史记录
echo "6. 测试获取账户历史记录..."
curl -s -X GET "$BASE_URL/api/v1/account-histories?current=1&pageSize=5" \
  -H "Authorization: Bearer $TOKEN" | jq '.'

echo ""

# 7. 测试搜索功能
echo "7. 测试搜索功能..."
echo "按用户名搜索:"
curl -s -X GET "$BASE_URL/api/v1/accounts?username=admin" \
  -H "Authorization: Bearer $TOKEN" | jq '.'

echo "按角色类型搜索:"
curl -s -X GET "$BASE_URL/api/v1/accounts?roleType=company_manager" \
  -H "Authorization: Bearer $TOKEN" | jq '.'

echo ""

# 8. 测试退出登录
echo "8. 测试退出登录..."
curl -s -X POST "$BASE_URL/api/v1/auth/logout" \
  -H "Authorization: Bearer $TOKEN" | jq '.'

echo ""
echo "=== 测试完成 ==="

echo ""
echo "=== 数据库验证 ==="
echo "请检查数据库中是否有以下数据："
echo "1. 账户表 (account) 中应该有5个测试账户"
echo "2. 历史记录表 (account_history) 中应该有5条创建记录"
echo "3. 新创建的账户应该有对应的历史记录"
echo ""
echo "可以使用以下SQL查询验证："
echo "SELECT COUNT(*) as account_count FROM go_gin_api.account;"
echo "SELECT COUNT(*) as history_count FROM go_gin_api.account_history;"
echo "SELECT * FROM go_gin_api.account WHERE username = 'testuser';" 