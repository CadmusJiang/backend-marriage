#!/bin/bash

# 日志查看器测试脚本
BASE_URL="http://localhost:9999"

echo "=== 日志查看器测试 ==="

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

# 2. 测试获取日志列表
echo "2. 测试获取日志列表..."
curl -s -X GET "$BASE_URL/api/v1/logs" | jq '.'

echo ""

# 3. 测试获取实时日志
echo "3. 测试获取实时日志..."
curl -s -X GET "$BASE_URL/api/v1/logs/realtime" | jq '.'

echo ""

# 4. 测试访问日志查看页面
echo "4. 测试访问日志查看页面..."
echo "请在浏览器中访问: http://localhost:9999/logs"

echo ""

# 5. 生成一些测试请求来产生日志
echo "5. 生成测试请求..."
echo "正在发送一些测试请求来产生日志..."

# 发送一些测试请求
curl -s -X GET "$BASE_URL/api/v1/check-db" > /dev/null
curl -s -X GET "$BASE_URL/api/v1/auth/login" -H "Content-Type: application/json" -d '{"username":"test","password":"test"}' > /dev/null
curl -s -X GET "$BASE_URL/swagger/index.html" > /dev/null
curl -s -X GET "$BASE_URL/metrics" > /dev/null

echo "测试请求已发送"

echo ""

# 6. 再次获取日志查看是否有新日志
echo "6. 再次获取日志列表..."
curl -s -X GET "$BASE_URL/api/v1/logs" | jq '.'

echo ""

echo "=== 日志查看器测试完成 ==="
echo ""
echo "访问日志查看页面: http://localhost:9999/logs"
echo "API接口:"
echo "  - GET /api/v1/logs - 获取日志列表"
echo "  - GET /api/v1/logs/realtime - 获取实时日志" 