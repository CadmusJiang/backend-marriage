#!/bin/bash

# Marriage System 停止脚本

echo "=== Marriage System 停止脚本 ==="

# 检查PID文件是否存在
if [ ! -f "marriage_system.pid" ]; then
    echo "❌ PID文件不存在，应用可能未运行"
    exit 1
fi

# 读取PID
PID=$(cat marriage_system.pid)

# 检查进程是否存在
if ! ps -p $PID > /dev/null 2>&1; then
    echo "❌ 进程 $PID 不存在，应用可能已停止"
    rm -f marriage_system.pid
    exit 1
fi

# 停止进程
echo "正在停止进程 $PID..."
kill $PID

# 等待进程结束
sleep 2

# 检查进程是否已停止
if ps -p $PID > /dev/null 2>&1; then
    echo "强制停止进程 $PID..."
    kill -9 $PID
fi

# 删除PID文件
rm -f marriage_system.pid

echo "✅ Marriage System 已停止"
echo "进程 $PID 已终止" 