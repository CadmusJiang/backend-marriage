#!/bin/bash

# 美化日志显示脚本
echo "=== 后端请求日志监控 ==="
echo "按 Ctrl+C 停止监控"
echo "----------------------------------------"

# 实时监控并美化日志输出
tail -f ./logs/marriage_system-access.log 2>/dev/null | while read line; do
    # 检查是否是trace-log
    if echo "$line" | grep -q '"msg":"trace-log"'; then
        # 提取关键信息
        timestamp=$(echo "$line" | grep -o '"time":"[^"]*"' | cut -d'"' -f4)
        method=$(echo "$line" | grep -o '"method":"[^"]*"' | cut -d'"' -f4)
        path=$(echo "$line" | grep -o '"path":"[^"]*"' | cut -d'"' -f4)
        http_code=$(echo "$line" | grep -o '"http_code":[0-9]*' | cut -d':' -f2)
        cost_seconds=$(echo "$line" | grep -o '"cost_seconds":[0-9.]*' | cut -d':' -f2)
        success=$(echo "$line" | grep -o '"success":[^,]*' | cut -d':' -f2)
        trace_id=$(echo "$line" | grep -o '"trace_id":"[^"]*"' | cut -d'"' -f4)
        
        # 格式化输出
        status_color="\033[32m"  # 绿色
        if [ "$success" = "false" ]; then
            status_color="\033[31m"  # 红色
        fi
        
        echo -e "[${timestamp}] ${status_color}${method}${path}\033[0m -> \033[33m${http_code}\033[0m (${cost_seconds}s) [${status_color}${success}\033[0m] [${trace_id}]"
    else
        # 其他类型的日志
        timestamp=$(echo "$line" | grep -o '"time":"[^"]*"' | cut -d'"' -f4)
        level=$(echo "$line" | grep -o '"level":"[^"]*"' | cut -d'"' -f4)
        msg=$(echo "$line" | grep -o '"msg":"[^"]*"' | cut -d'"' -f4)
        
        level_color="\033[33m"  # 黄色
        if [ "$level" = "error" ] || [ "$level" = "fatal" ]; then
            level_color="\033[31m"  # 红色
        elif [ "$level" = "info" ]; then
            level_color="\033[32m"  # 绿色
        fi
        
        echo -e "[${timestamp}] ${level_color}[${level}]\033[0m ${msg}"
    fi
done 