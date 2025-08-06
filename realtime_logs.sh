#!/bin/bash

# 实时日志查看器
clear
echo "=========================================="
echo "    后端实时日志监控器"
echo "=========================================="
echo "正在监控前端请求..."
echo "按 Ctrl+C 停止监控"
echo "=========================================="

# 监控日志文件并实时显示
tail -f ./logs/marriage_system-access.log | while read line; do
    # 解析JSON日志
    if echo "$line" | grep -q '"msg":"trace-log"'; then
        # 提取信息
        time=$(echo "$line" | jq -r '.time' 2>/dev/null)
        method=$(echo "$line" | jq -r '.method' 2>/dev/null)
        path=$(echo "$line" | jq -r '.path' 2>/dev/null)
        http_code=$(echo "$line" | jq -r '.http_code' 2>/dev/null)
        cost=$(echo "$line" | jq -r '.cost_seconds' 2>/dev/null)
        success=$(echo "$line" | jq -r '.success' 2>/dev/null)
        trace_id=$(echo "$line" | jq -r '.trace_id' 2>/dev/null)
        
        # 格式化输出
        status="✅"
        if [ "$success" = "false" ]; then
            status="❌"
        fi
        
        echo "[$(date '+%H:%M:%S')] $status $method $path -> $http_code (${cost}s) [$trace_id]"
    fi
done 