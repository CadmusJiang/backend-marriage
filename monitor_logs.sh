#!/bin/bash

# 实时监控日志文件
echo "开始监控后端日志..."
echo "按 Ctrl+C 停止监控"
echo "----------------------------------------"

# 监控访问日志文件
tail -f ./logs/marriage_system-access.log | while read line; do
    echo "[$(date '+%Y-%m-%d %H:%M:%S')] $line"
done 