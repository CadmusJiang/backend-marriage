#!/bin/bash

echo "🧪 测试Teams API的日志系统"
echo "===================="

# 设置测试用的trace_id
TRACE_ID="test-teams-$(date +%s)"

echo "1. 使用Trace-ID: $TRACE_ID 调用Teams API..."
echo ""

# 调用teams API
echo "📡 调用 GET /api/v1/teams..."
RESPONSE=$(curl -s -H "X-Trace-ID: $TRACE_ID" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  "http://localhost:9999/api/v1/teams?current=1&pageSize=10")

echo "响应状态: $?"
echo "响应内容: ${RESPONSE:0:200}..."
echo ""

# 检查日志文件
echo "2. 检查日志文件中的trace_id记录..."
echo ""

# 查找包含trace_id的日志
echo "在logs/app.log中查找trace_id: $TRACE_ID"
if [ -f "logs/app.log" ]; then
    LOG_COUNT=$(grep -c "$TRACE_ID" logs/app.log 2>/dev/null || echo "0")
    echo "找到 $LOG_COUNT 条相关日志"
    
    if [ "$LOG_COUNT" -gt 0 ]; then
        echo ""
        echo "最近的几条日志:"
        grep "$TRACE_ID" logs/app.log | tail -3 | while read -r line; do
            echo "  ${line:0:100}..."
        done
    fi
else
    echo "日志文件 logs/app.log 不存在"
fi

echo ""

# 查找包含trace_id的访问日志
echo "在logs/marriage_system-access.log中查找trace_id: $TRACE_ID"
if [ -f "logs/marriage_system-access.log" ]; then
    ACCESS_LOG_COUNT=$(grep -c "$TRACE_ID" logs/marriage_system-access.log 2>/dev/null || echo "0")
    echo "找到 $ACCESS_LOG_COUNT 条访问日志"
    
    if [ "$ACCESS_LOG_COUNT" -gt 0 ]; then
        echo ""
        echo "最近的几条访问日志:"
        grep "$TRACE_ID" logs/marriage_system-access.log | tail -3 | while read -r line; do
            echo "  ${line:0:100}..."
        done
    fi
else
    echo "访问日志文件 logs/marriage_system-access.log 不存在"
fi

echo ""
echo "3. 测试其他Teams API端点..."
echo ""

# 测试创建团队
echo "📡 调用 POST /api/v1/teams..."
CREATE_RESPONSE=$(curl -s -H "X-Trace-ID: $TRACE_ID" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -H "Content-Type: application/json" \
  -d '{"belongGroupId": 1, "username": "test-team", "name": "测试团队"}' \
  "http://localhost:9999/api/v1/teams")

echo "创建团队响应: ${CREATE_RESPONSE:0:200}..."
echo ""

# 测试获取团队详情
echo "📡 调用 GET /api/v1/teams/1..."
DETAIL_RESPONSE=$(curl -s -H "X-Trace-ID: $TRACE_ID" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  "http://localhost:9999/api/v1/teams/1")

echo "获取团队详情响应: ${DETAIL_RESPONSE:0:200}..."
echo ""

echo "4. 检查日志文件中的新记录..."
echo ""

# 再次检查日志文件
if [ -f "logs/app.log" ]; then
    NEW_LOG_COUNT=$(grep -c "$TRACE_ID" logs/app.log 2>/dev/null || echo "0")
    echo "现在总共有 $NEW_LOG_COUNT 条相关日志"
    
    if [ "$NEW_LOG_COUNT" -gt 0 ]; then
        echo ""
        echo "包含operation字段的日志:"
        grep "$TRACE_ID" logs/app.log | grep "operation" | tail -3 | while read -r line; do
            echo "  ${line:0:120}..."
        done
    fi
fi

echo ""
echo "✅ 测试完成！"
echo ""
echo "💡 提示:"
echo "- 确保服务正在运行 (http://localhost:9999)"
echo "- 确保有有效的认证token"
echo "- 检查logs/app.log和logs/marriage_system-access.log文件"
echo "- 所有日志都应该包含trace_id: $TRACE_ID"
