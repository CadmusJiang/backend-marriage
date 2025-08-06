#!/bin/bash

# 重启脚本 - 管理9999端口的进程
# 使用方法: ./scripts/restart_better.sh [binary_name]

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 默认二进制文件名
DEFAULT_BINARY="main"
BINARY_NAME=${1:-$DEFAULT_BINARY}

# 端口号
PORT=9999

echo -e "${BLUE}=== 重启脚本 ===${NC}"
echo -e "${BLUE}目标端口: ${PORT}${NC}"
echo -e "${BLUE}二进制文件: ${BINARY_NAME}${NC}"
echo ""

# 检查二进制文件是否存在
if [ ! -f "../${BINARY_NAME}" ]; then
    echo -e "${RED}错误: 二进制文件 ${BINARY_NAME} 不存在${NC}"
    echo -e "${YELLOW}可用的二进制文件:${NC}"
    ls -la ../ | grep -E "(main|marriage|test_build)" | awk '{print $9}' || echo "没有找到可执行文件"
    exit 1
fi

# 函数：查找占用端口的进程
find_port_process() {
    echo -e "${YELLOW}正在查找占用端口 ${PORT} 的进程...${NC}"
    
    # 查找占用端口的进程
    PID=$(lsof -ti:${PORT} 2>/dev/null)
    
    if [ -z "$PID" ]; then
        echo -e "${GREEN}端口 ${PORT} 当前没有被占用${NC}"
        return 1
    else
        echo -e "${YELLOW}找到进程 PID: ${PID}${NC}"
        
        # 显示进程详细信息
        echo -e "${BLUE}进程详细信息:${NC}"
        ps -p $PID -o pid,ppid,user,command --no-headers 2>/dev/null || echo "无法获取进程信息"
        
        return 0
    fi
}

# 函数：杀掉进程
kill_process() {
    local pid=$1
    
    if [ -z "$pid" ]; then
        echo -e "${YELLOW}没有进程需要杀掉${NC}"
        return
    fi
    
    echo -e "${YELLOW}正在杀掉进程 PID: ${pid}${NC}"
    
    # 尝试优雅地杀掉进程
    kill $pid 2>/dev/null
    
    # 等待进程结束
    local count=0
    while kill -0 $pid 2>/dev/null && [ $count -lt 10 ]; do
        echo -e "${YELLOW}等待进程结束... (${count}/10)${NC}"
        sleep 1
        count=$((count + 1))
    done
    
    # 如果进程还在运行，强制杀掉
    if kill -0 $pid 2>/dev/null; then
        echo -e "${RED}进程仍在运行，强制杀掉...${NC}"
        kill -9 $pid 2>/dev/null
        sleep 2
    fi
    
    # 验证进程是否已经被杀掉
    if kill -0 $pid 2>/dev/null; then
        echo -e "${RED}错误: 无法杀掉进程 PID: ${pid}${NC}"
        return 1
    else
        echo -e "${GREEN}进程已成功杀掉${NC}"
        return 0
    fi
}

# 函数：启动新进程
start_new_process() {
    echo -e "${YELLOW}正在启动新的二进制文件: ${BINARY_NAME}${NC}"
    
    # 检查端口是否已经被释放
    if lsof -ti:${PORT} >/dev/null 2>&1; then
        echo -e "${RED}错误: 端口 ${PORT} 仍被占用${NC}"
        return 1
    fi
    
    # 启动新进程
    echo -e "${BLUE}启动命令: ../${BINARY_NAME}${NC}"
    
    # 在后台启动进程
    nohup ../${BINARY_NAME} > ../logs/app.log 2>&1 &
    NEW_PID=$!
    
    echo -e "${GREEN}新进程已启动，PID: ${NEW_PID}${NC}"
    
    # 等待几秒钟让进程完全启动
    sleep 3
    
    # 验证进程是否正在运行
    if kill -0 $NEW_PID 2>/dev/null; then
        echo -e "${GREEN}进程正在运行${NC}"
        
        # 检查端口是否被新进程占用
        if lsof -ti:${PORT} | grep -q $NEW_PID; then
            echo -e "${GREEN}端口 ${PORT} 已被新进程占用${NC}"
        else
            echo -e "${YELLOW}警告: 端口 ${PORT} 可能没有被新进程占用${NC}"
        fi
        
        return 0
    else
        echo -e "${RED}错误: 新进程启动失败${NC}"
        return 1
    fi
}

# 主执行流程
main() {
    echo -e "${BLUE}开始执行重启流程...${NC}"
    echo ""
    
    # 1. 查找占用端口的进程
    if find_port_process; then
        PID=$(lsof -ti:${PORT})
        
        # 2. 杀掉进程
        if kill_process $PID; then
            echo -e "${GREEN}进程已成功杀掉${NC}"
        else
            echo -e "${RED}杀掉进程失败${NC}"
            exit 1
        fi
    else
        echo -e "${GREEN}没有进程需要杀掉${NC}"
    fi
    
    echo ""
    
    # 3. 启动新进程
    if start_new_process; then
        echo -e "${GREEN}重启完成！${NC}"
        echo -e "${BLUE}新进程 PID: $(lsof -ti:${PORT} 2>/dev/null || echo 'N/A')${NC}"
        echo -e "${BLUE}日志文件: ../logs/app.log${NC}"
    else
        echo -e "${RED}启动新进程失败${NC}"
        exit 1
    fi
}

# 执行主函数
main "$@" 