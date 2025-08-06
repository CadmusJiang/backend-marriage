#!/bin/bash

# 测试账户API的belongGroup和belongTeam字段

BASE_URL="http://localhost:9999"

echo "=== 测试账户API的belongGroup和belongTeam字段 ==="

# 1. 测试获取账户列表
echo "1. 获取账户列表（包含组织信息）"
curl -X GET "${BASE_URL}/api/v1/accounts?current=1&pageSize=5" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <your_token>" \
  | jq '.'

echo -e "\n"

# 2. 测试获取特定账户详情
echo "2. 获取特定账户详情"
curl -X GET "${BASE_URL}/api/v1/accounts/acc_001" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <your_token>" \
  | jq '.'

echo -e "\n"

# 3. 测试数据库查询
echo "3. 检查数据库中的账户组织信息"
echo "SELECT account_id, username, nickname, belong_group_id, belong_group_nickname, belong_team_id, belong_team_nickname FROM account LIMIT 5;"

echo -e "\n"

# 4. 测试组织API
echo "4. 获取组织列表"
curl -X GET "${BASE_URL}/api/v1/org-infos?current=1&pageSize=5" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <your_token>" \
  | jq '.'

echo -e "\n"

echo "=== 测试完成 ==="
echo "如果belongGroup和belongTeam字段仍然为空，请检查："
echo "1. 数据库中账户数据是否包含组织信息"
echo "2. 账户模型是否正确映射了数据库字段"
echo "3. API实现是否正确处理了组织信息" 