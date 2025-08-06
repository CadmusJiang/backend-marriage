#!/bin/bash

# Groups API测试脚本
BASE_URL="http://localhost:9999"

echo "=== Groups API测试 ==="

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

# 2. 测试获取groups列表
echo "2. 测试获取groups列表..."
curl -s -X GET "$BASE_URL/api/v1/groups?current=1&pageSize=10" \
  -H "Authorization: Bearer $TOKEN" | jq '.'

echo ""

# 3. 测试按名称搜索groups
echo "3. 测试按名称搜索groups..."
curl -s -X GET "$BASE_URL/api/v1/groups?current=1&pageSize=5&orgName=系统" \
  -H "Authorization: Bearer $TOKEN" | jq '.'

echo ""

# 4. 测试按层级筛选groups
echo "4. 测试按层级筛选groups..."
curl -s -X GET "$BASE_URL/api/v1/groups?current=1&pageSize=5&orgLevel=1" \
  -H "Authorization: Bearer $TOKEN" | jq '.'

echo ""

# 5. 测试分页功能
echo "5. 测试分页功能..."
curl -s -X GET "$BASE_URL/api/v1/groups?current=1&pageSize=3" \
  -H "Authorization: Bearer $TOKEN" | jq '.'

echo ""

# 6. 测试获取所有groups（不分页）
echo "6. 测试获取所有groups（不分页）..."
curl -s -X GET "$BASE_URL/api/v1/groups?current=1&pageSize=100" \
  -H "Authorization: Bearer $TOKEN" | jq '.'

echo ""

echo "=== Groups API测试完成 ===" 