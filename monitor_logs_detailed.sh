#!/bin/bash

# 详细日志监控脚本
echo "=== 后端日志监控工具 ==="
echo "1. 实时监控访问日志文件"
echo "2. 实时监控控制台输出"
echo "3. 显示最近的日志"
echo "4. 退出"
echo "----------------------------------------"

# 检查日志文件是否存在
if [ ! -f "./logs/marriage_system-access.log" ]; then
    echo "警告: 日志文件不存在，请确保服务器正在运行"
    echo "创建日志目录..."
    mkdir -p ./logs
fi

# 显示最近的日志
echo "最近的日志记录:"
echo "----------------------------------------"
if [ -f "./logs/marriage_system-access.log" ]; then
    tail -20 ./logs/marriage_system-access.log
else
    echo "日志文件不存在，请启动服务器后重试"
fi

echo ""
echo "实时监控已启动 (按 Ctrl+C 停止):"
echo "----------------------------------------"

# 实时监控日志文件
tail -f ./logs/marriage_system-access.log 2>/dev/null | while read line; do
    # 解析JSON日志并格式化输出
    if echo "$line" | grep -q '"level":"info"'; then
        # 提取关键信息
        method=$(echo "$line" | grep -o '"method":"[^"]*"' | cut -d'"' -f4)
        path=$(echo "$line" | grep -o '"path":"[^"]*"' | cut -d'"' -f4)
        http_code=$(echo "$line" | grep -o '"http_code":[0-9]*' | cut -d':' -f2)
        cost_seconds=$(echo "$line" | grep -o '"cost_seconds":[0-9.]*' | cut -d':' -f2)
        success=$(echo "$line" | grep -o '"success":[^,]*' | cut -d':' -f2)
        
        # 格式化输出
        echo "[$(date '+%H:%M:%S')] $method $path -> $http_code (${cost_seconds}s) [成功: $success]"
    else
        echo "[$(date '+%H:%M:%S')] $line"
    fi
done 