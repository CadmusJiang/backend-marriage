#!/bin/bash

# 组织信息API测试脚本
BASE_URL="http://localhost:9999"

echo "=== 组织信息API测试 ==="

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

# 2. 测试获取组织信息列表
echo "2. 测试获取组织信息列表..."
curl -s -X GET "$BASE_URL/api/v1/org-infos?current=1&pageSize=10" \
  -H "Authorization: Bearer $TOKEN" | jq '.'

echo ""

# 3. 测试获取组织信息详情
echo "3. 测试获取组织信息详情..."
curl -s -X GET "$BASE_URL/api/v1/org-infos/org_001" \
  -H "Authorization: Bearer $TOKEN" | jq '.'

echo ""

# 4. 测试创建组织信息
echo "4. 测试创建组织信息..."
curl -s -X POST "$BASE_URL/api/v1/org-infos" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "orgName": "测试组织",
    "orgType": "group",
    "orgLevel": 1,
    "orgDescription": "这是一个测试组织",
    "maxCnt": 20
  }' | jq '.'

echo ""

# 5. 测试更新组织信息
echo "5. 测试更新组织信息..."
curl -s -X PUT "$BASE_URL/api/v1/org-infos/org_001" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "orgName": "更新后的系统管理组",
    "orgDescription": "更新后的系统管理相关组织",
    "maxCnt": 15
  }' | jq '.'

echo ""

# 6. 测试获取子组织
echo "6. 测试获取子组织..."
curl -s -X GET "$BASE_URL/api/v1/org-infos/org_001/children" \
  -H "Authorization: Bearer $TOKEN" | jq '.'

echo ""

# 7. 测试获取父组织
echo "7. 测试获取父组织..."
curl -s -X GET "$BASE_URL/api/v1/org-infos/org_006/parent" \
  -H "Authorization: Bearer $TOKEN" | jq '.'

echo ""

# 8. 测试搜索功能
echo "8. 测试搜索功能..."
echo "按组织名称搜索:"
curl -s -X GET "$BASE_URL/api/v1/org-infos?orgName=系统" \
  -H "Authorization: Bearer $TOKEN" | jq '.'

echo "按组织类型搜索:"
curl -s -X GET "$BASE_URL/api/v1/org-infos?orgType=group" \
  -H "Authorization: Bearer $TOKEN" | jq '.'

echo "按组织层级搜索:"
curl -s -X GET "$BASE_URL/api/v1/org-infos?orgLevel=1" \
  -H "Authorization: Bearer $TOKEN" | jq '.'

echo ""

# 9. 测试删除组织信息
echo "9. 测试删除组织信息..."
curl -s -X DELETE "$BASE_URL/api/v1/org-infos/org_010" \
  -H "Authorization: Bearer $TOKEN" | jq '.'

echo ""

# 10. 测试退出登录
echo "10. 测试退出登录..."
curl -s -X POST "$BASE_URL/api/v1/auth/logout" \
  -H "Authorization: Bearer $TOKEN" | jq '.'

echo ""
echo "=== 测试完成 ==="

echo ""
echo "=== 数据库验证 ==="
echo "请检查数据库中是否有以下数据："
echo "1. 组织信息表 (org) 中应该有10个测试组织"
echo "2. 账户表 (account) 中应该有5个测试账户"
echo "3. 历史记录表 (account_history) 中应该有5条创建记录"
echo ""
echo "可以使用以下SQL查询验证："
echo "USE marriage_system;"
echo "SELECT COUNT(*) as org_count FROM org;"
echo "SELECT COUNT(*) as account_count FROM account;"
echo "SELECT COUNT(*) as history_count FROM account_history;"
echo "SELECT * FROM org WHERE org_name LIKE '%测试%';" 