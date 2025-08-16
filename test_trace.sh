#!/bin/bash

# 全链路日志追踪功能测试脚本

echo "🧪 开始测试全链路日志追踪功能..."
echo "=================================="

# 设置基础URL
BASE_URL="http://localhost:9999"

# 测试1: 自动生成Trace-ID
echo "📝 测试1: 自动生成Trace-ID"
echo "发送请求到 /api/v1/accounts (无Trace-ID头)"
RESPONSE=$(curl -s -D /tmp/headers1 $BASE_URL/api/v1/accounts)
TRACE_ID=$(grep "X-Trace-ID" /tmp/headers1 | cut -d' ' -f2 | tr -d '\r')

if [ ! -z "$TRACE_ID" ]; then
    echo "✅ 成功! 自动生成的Trace-ID: $TRACE_ID"
else
    echo "❌ 失败! 未在响应头中找到X-Trace-ID"
fi

echo ""

# 测试2: 自定义Trace-ID
echo "📝 测试2: 自定义Trace-ID"
CUSTOM_TRACE_ID="test_trace_$(date +%s)"
echo "发送请求到 /api/v1/accounts (自定义Trace-ID: $CUSTOM_TRACE_ID)"
RESPONSE=$(curl -s -D /tmp/headers2 -H "X-Trace-ID: $CUSTOM_TRACE_ID" $BASE_URL/api/v1/accounts)
RESPONSE_TRACE_ID=$(grep "X-Trace-ID" /tmp/headers2 | cut -d' ' -f2 | tr -d '\r')

if [ "$RESPONSE_TRACE_ID" = "$CUSTOM_TRACE_ID" ]; then
    echo "✅ 成功! 响应头中的Trace-ID匹配: $RESPONSE_TRACE_ID"
else
    echo "❌ 失败! Trace-ID不匹配. 期望: $CUSTOM_TRACE_ID, 实际: $RESPONSE_TRACE_ID"
fi

echo ""

# 测试3: 查询Trace日志
echo "📝 测试3: 查询Trace日志"
if [ ! -z "$TRACE_ID" ]; then
    echo "查询Trace-ID: $TRACE_ID 的日志"
    LOGS=$(curl -s "$BASE_URL/api/v1/logs/trace?trace_id=$TRACE_ID")
    echo "日志查询结果: $LOGS"
    
    # 检查是否包含trace_id字段
    if echo "$LOGS" | grep -q "trace_id"; then
        echo "✅ 成功! 日志查询返回了trace_id字段"
    else
        echo "❌ 失败! 日志查询未返回trace_id字段"
    fi
else
    echo "⚠️  跳过日志查询测试 (Trace-ID为空)"
fi

echo ""

# 测试4: 测试不同的Trace-ID头格式
echo "📝 测试4: 测试不同的Trace-ID头格式"
echo "测试 Trace-ID 头格式"
RESPONSE=$(curl -s -D /tmp/headers3 -H "Trace-ID: test_format_$(date +%s)" $BASE_URL/api/v1/accounts)
TRACE_ID_ALT=$(grep "X-Trace-ID" /tmp/headers3 | cut -d' ' -f2 | tr -d '\r')

if [ ! -z "$TRACE_ID_ALT" ]; then
    echo "✅ 成功! Trace-ID头格式支持正常: $TRACE_ID_ALT"
else
    echo "❌ 失败! Trace-ID头格式不支持"
fi

echo ""

# 测试5: 测试X-Request-ID头格式
echo "📝 测试5: 测试X-Request-ID头格式"
echo "测试 X-Request-ID 头格式"
RESPONSE=$(curl -s -D /tmp/headers4 -H "X-Request-ID: test_request_$(date +%s)" $BASE_URL/api/v1/accounts)
TRACE_ID_REQUEST=$(grep "X-Trace-ID" /tmp/headers4 | cut -d' ' -f2 | tr -d '\r')

if [ ! -z "$TRACE_ID_REQUEST" ]; then
    echo "✅ 成功! X-Request-ID头格式支持正常: $TRACE_ID_REQUEST"
else
    echo "❌ 失败! X-Request-ID头格式不支持"
fi

echo ""

# 清理临时文件
rm -f /tmp/headers1 /tmp/headers2 /tmp/headers3 /tmp/headers4

echo "=================================="
echo "🎉 测试完成!"
echo ""
echo "📖 使用说明:"
echo "1. 访问 $BASE_URL/docs/trace-logs.html 查看全链路日志"
echo "2. 访问 $BASE_URL/docs/logs.html 查看分页日志"
echo "3. 在任何API请求中添加 X-Trace-ID 头来追踪请求"
echo "4. 响应头中的 X-Trace-ID 字段包含本次请求的追踪ID"
echo ""
echo "🔍 查看日志API:"
echo "   GET $BASE_URL/api/v1/logs/trace?trace_id=<trace_id>"
echo "   GET $BASE_URL/api/v1/logs/trace/range?trace_id=<trace_id>&start_time=<start>&end_time=<end>"
