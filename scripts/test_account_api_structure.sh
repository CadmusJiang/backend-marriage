#!/bin/bash

# 测试账户API返回的数据结构是否与前端mock代码一致

BASE_URL="http://localhost:9999"

echo "=== 测试账户API数据结构 ==="

# 1. 测试获取账户列表（包含组织信息）
echo "1. 获取账户列表（包含组织信息）"
curl -s -X GET "${BASE_URL}/api/v1/accounts?current=1&pageSize=3" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <your_token>" \
  | jq '.'

echo -e "\n"

# 2. 测试获取账户列表（不包含组织信息）
echo "2. 获取账户列表（不包含组织信息）"
curl -s -X GET "${BASE_URL}/api/v1/accounts?current=1&pageSize=3&includeGroup=false&includeTeam=false" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <your_token>" \
  | jq '.'

echo -e "\n"

# 3. 测试获取特定账户详情（包含组织信息）
echo "3. 获取特定账户详情（包含组织信息）"
curl -s -X GET "${BASE_URL}/api/v1/accounts/0" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <your_token>" \
  | jq '.'

echo -e "\n"

# 4. 测试获取特定账户详情（不包含组织信息）
echo "4. 获取特定账户详情（不包含组织信息）"
curl -s -X GET "${BASE_URL}/api/v1/accounts/0?includeGroup=false&includeTeam=false" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <your_token>" \
  | jq '.'

echo -e "\n"

# 5. 测试搜索功能
echo "5. 测试搜索功能"
curl -s -X GET "${BASE_URL}/api/v1/accounts?current=1&pageSize=5&nickname=张伟" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <your_token>" \
  | jq '.'

echo -e "\n"

# 6. 测试角色筛选
echo "6. 测试角色筛选"
curl -s -X GET "${BASE_URL}/api/v1/accounts?current=1&pageSize=5&roleType=company_manager" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <your_token>" \
  | jq '.'

echo -e "\n"

echo "=== 数据结构验证 ==="
echo "验证返回的数据结构是否包含以下字段："
echo "- id: 账户ID"
echo "- username: 用户名"
echo "- nickname: 姓名"
echo "- phone: 手机号"
echo "- roleType: 角色类型"
echo "- belongGroup: 所属组对象（可选）"
echo "- belongTeam: 所属团队对象（可选）"
echo "- status: 状态"
echo "- createdTimestamp: 创建时间戳"
echo "- modifiedTimestamp: 修改时间戳"
echo "- lastLoginTimestamp: 最后登录时间戳"

echo -e "\n"

echo "=== 组织对象结构验证 ==="
echo "belongGroup 和 belongTeam 对象应包含："
echo "- id: 组织ID"
echo "- username: 组织用户名"
echo "- nickname: 组织名称"
echo "- createdTimestamp: 创建时间戳"
echo "- modifiedTimestamp: 修改时间戳"
echo "- currentCnt: 当前成员数量"

echo -e "\n"

echo "=== 测试完成 ==="
echo "如果返回的数据结构与前端mock代码一致，说明API实现正确。" 