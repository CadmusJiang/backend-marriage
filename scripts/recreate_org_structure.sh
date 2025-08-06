#!/bin/bash

# 重新创建正确的组织结构脚本
# 确保每个team只属于一个组

BASE_URL="http://localhost:9999"

echo "=========================================="
echo "开始重新创建正确的组织结构..."
echo "=========================================="

# 1. 创建组
echo "1. 创建组..."

echo "创建系统管理组..."
curl -s -X POST "$BASE_URL/api/v1/org-infos" \
  -H "Content-Type: application/json" \
  -d '{"orgName":"系统管理组","orgType":"1","orgPath":"/1","orgLevel":1,"maxCnt":10}' | jq '.'

echo "创建南京-天元大厦组..."
curl -s -X POST "$BASE_URL/api/v1/org-infos" \
  -H "Content-Type: application/json" \
  -d '{"orgName":"南京-天元大厦组","orgType":"1","orgPath":"/2","orgLevel":1,"maxCnt":50}' | jq '.'

echo "创建南京-南京南站组..."
curl -s -X POST "$BASE_URL/api/v1/org-infos" \
  -H "Content-Type: application/json" \
  -d '{"orgName":"南京-南京南站组","orgType":"1","orgPath":"/3","orgLevel":1,"maxCnt":50}' | jq '.'

# 2. 创建团队 - 每个团队只属于一个组
echo ""
echo "2. 创建团队..."

echo "创建系统管理团队 (属于系统管理组)..."
curl -s -X POST "$BASE_URL/api/v1/org-infos" \
  -H "Content-Type: application/json" \
  -d '{"orgName":"系统管理团队","orgType":"2","orgPath":"/1/1","orgLevel":2,"maxCnt":10}' | jq '.'

echo "创建营销团队A (属于南京-天元大厦组)..."
curl -s -X POST "$BASE_URL/api/v1/org-infos" \
  -H "Content-Type: application/json" \
  -d '{"orgName":"营销团队A","orgType":"2","orgPath":"/2/1","orgLevel":2,"maxCnt":20}' | jq '.'

echo "创建营销团队B (属于南京-南京南站组)..."
curl -s -X POST "$BASE_URL/api/v1/org-infos" \
  -H "Content-Type: application/json" \
  -d '{"orgName":"营销团队B","orgType":"2","orgPath":"/3/1","orgLevel":2,"maxCnt":20}' | jq '.'

echo "创建营销团队C (属于南京-天元大厦组)..."
curl -s -X POST "$BASE_URL/api/v1/org-infos" \
  -H "Content-Type: application/json" \
  -d '{"orgName":"营销团队C","orgType":"2","orgPath":"/2/2","orgLevel":2,"maxCnt":20}' | jq '.'

# 3. 验证结果
echo ""
echo "3. 验证结果..."

echo "查看所有teams:"
curl -s -X GET "$BASE_URL/api/v1/teams" \
  -H "Content-Type: application/json" | jq '.'

echo ""
echo "查看所有groups (包含teams):"
curl -s -X GET "$BASE_URL/api/v1/groups?includeTeams=true" \
  -H "Content-Type: application/json" | jq '.'

echo ""
echo "=========================================="
echo "重新创建完成！"
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
echo "每个team现在只属于一个组！" 