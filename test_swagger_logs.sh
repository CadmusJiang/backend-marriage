#!/bin/bash

echo "🧪 测试Swagger日志接口"
echo "===================="

# 检查服务是否运行
echo "1. 检查服务状态..."
if ! curl -s "http://localhost:9999/health" > /dev/null 2>&1; then
    echo "❌ 服务未运行，请先启动服务"
    exit 1
fi
echo "✅ 服务正在运行"
echo ""

# 测试日志接口
echo "2. 测试日志相关接口..."
echo ""

# 测试获取最近日志
echo "📡 测试 GET /api/v1/logs/latest..."
RESPONSE=$(curl -s "http://localhost:9999/api/v1/logs/latest")
if [ $? -eq 0 ]; then
    echo "✅ 接口响应正常"
    echo "响应内容: ${RESPONSE:0:200}..."
else
    echo "❌ 接口调用失败"
fi
echo ""

# 测试获取统一日志
echo "📡 测试 GET /api/v1/logs/unified..."
RESPONSE=$(curl -s "http://localhost:9999/api/v1/logs/unified")
if [ $? -eq 0 ]; then
    echo "✅ 接口响应正常"
    echo "响应内容: ${RESPONSE:0:200}..."
else
    echo "❌ 接口调用失败"
fi
echo ""

# 测试获取分页日志
echo "📡 测试 GET /api/v1/logs/paginated..."
RESPONSE=$(curl -s "http://localhost:9999/api/v1/logs/paginated?page=1&page_size=5")
if [ $? -eq 0 ]; then
    echo "✅ 接口响应正常"
    echo "响应内容: ${RESPONSE:0:200}..."
else
    echo "❌ 接口调用失败"
fi
echo ""

# 测试获取全链路日志
echo "📡 测试 GET /api/v1/logs/trace..."
RESPONSE=$(curl -s "http://localhost:9999/api/v1/logs/trace?trace_id=test-123")
if [ $? -eq 0 ]; then
    echo "✅ 接口响应正常"
    echo "响应内容: ${RESPONSE:0:200}..."
else
    echo "❌ 接口调用失败"
fi
echo ""

# 测试按时间范围获取日志
echo "📡 测试 GET /api/v1/logs/trace/range..."
RESPONSE=$(curl -s "http://localhost:9999/api/v1/logs/trace/range?trace_id=test-123&start_time=2024-01-01 00:00:00&end_time=2024-01-01 23:59:59")
if [ $? -eq 0 ]; then
    echo "✅ 接口响应正常"
    echo "响应内容: ${RESPONSE:0:200}..."
else
    echo "❌ 接口调用失败"
fi
echo ""

# 检查Swagger文档
echo "3. 检查Swagger文档..."
echo ""

# 检查swagger.yaml文件
if [ -f "docs/swagger.yaml" ]; then
    echo "✅ swagger.yaml 文件存在"
    
    # 检查是否包含日志相关定义
    if grep -q "logs.logEntry" docs/swagger.yaml; then
        echo "✅ 包含日志数据模型定义"
    else
        echo "❌ 缺少日志数据模型定义"
    fi
    
    if grep -q "/api/v1/logs/trace" docs/swagger.yaml; then
        echo "✅ 包含日志API接口定义"
    else
        echo "❌ 缺少日志API接口定义"
    fi
    
    if grep -q "name: Logs" docs/swagger.yaml; then
        echo "✅ 包含Logs标签定义"
    else
        echo "❌ 缺少Logs标签定义"
    fi
else
    echo "❌ swagger.yaml 文件不存在"
fi
echo ""

# 检查Swagger UI
echo "4. 检查Swagger UI..."
echo ""

echo "📖 访问Swagger文档:"
echo "   http://localhost:9999/docs/swagger.html"
echo ""
echo "🔍 在Swagger UI中查找 'Logs' 标签，应该包含以下接口："
echo "   - GET /api/v1/logs/latest - 获取最近日志"
echo "   - GET /api/v1/logs/unified - 获取统一日志数据"
echo "   - GET /api/v1/logs/paginated - 分页获取日志"
echo "   - GET /api/v1/logs/trace - 获取全链路日志"
echo "   - GET /api/v1/logs/trace/range - 按时间范围获取全链路日志"
echo ""

echo "✅ 测试完成！"
echo ""
echo "💡 提示:"
echo "- 所有日志接口都应该在Swagger文档中可见"
echo "- 接口应该包含完整的参数说明和响应模型"
echo "- 可以在Swagger UI中直接测试这些接口"
