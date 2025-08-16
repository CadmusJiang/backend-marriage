#!/bin/bash

echo "🧪 测试全链路日志追踪功能"
echo "================================"

BASE_URL="http://localhost:9999"

# 检查服务是否运行
echo "1. 检查服务状态..."
if ! curl -s "$BASE_URL/api/v1/logs" > /dev/null; then
    echo "❌ 服务未运行，请先启动服务"
    exit 1
fi
echo "✅ 服务正在运行"

# 测试Trace日志查询
echo ""
echo "2. 测试Trace日志查询..."
echo "查询Trace-ID: 123456789"

RESPONSE=$(curl -s "$BASE_URL/api/v1/logs/trace?trace_id=123456789")
echo "API响应:"
echo "$RESPONSE" | jq '.' 2>/dev/null || echo "$RESPONSE"

# 测试时间范围查询
echo ""
echo "3. 测试时间范围查询..."
echo "查询最近1小时的日志"

START_TIME=$(date -d '1 hour ago' '+%Y-%m-%d %H:%M:%S')
END_TIME=$(date '+%Y-%m-%d %H:%M:%S')

echo "时间范围: $START_TIME ~ $END_TIME"

RESPONSE_RANGE=$(curl -s "$BASE_URL/api/v1/logs/trace/range?trace_id=123456789&start_time=$START_TIME&end_time=$END_TIME")
echo "API响应:"
echo "$RESPONSE_RANGE" | jq '.' 2>/dev/null || echo "$RESPONSE_RANGE"

echo ""
echo "✅ 测试完成！"
echo ""
echo "📝 使用说明："
echo "1. 访问 http://localhost:9999/docs/trace-logs.html"
echo "2. 输入Trace-ID: 123456789"
echo "3. 点击搜索按钮查看全链路日志"
