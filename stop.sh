#!/bin/bash

# 项目停止脚本
# 停止运行在指定端口的服务

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 项目配置
PROJECT_NAME="marriage_system"
PORT="9999"

echo -e "${BLUE}================================${NC}"
echo -e "${BLUE}    ${PROJECT_NAME} 停止脚本${NC}"
echo -e "${BLUE}================================${NC}"

# 检查端口是否被占用
check_port() {
    echo -e "${YELLOW}检查端口 ${PORT} 是否被占用...${NC}"
    
    if lsof -Pi :$PORT -sTCP:LISTEN -t >/dev/null 2>&1; then
        echo -e "${GREEN}发现运行在端口 ${PORT} 的进程${NC}"
        return 0
    else
        echo -e "${YELLOW}端口 ${PORT} 没有被占用${NC}"
        return 1
    fi
}

# 停止服务
stop_service() {
    echo -e "${YELLOW}正在停止服务...${NC}"
    
    # 获取占用端口的进程PID
    PIDS=$(lsof -Pi :$PORT -sTCP:LISTEN -t 2>/dev/null)
    
    if [ -z "$PIDS" ]; then
        echo -e "${YELLOW}没有找到运行在端口 ${PORT} 的进程${NC}"
        return
    fi
    
    # 显示进程信息
    echo -e "${YELLOW}找到以下进程:${NC}"
    lsof -i :$PORT
    
    # 停止进程
    for PID in $PIDS; do
        echo -e "${YELLOW}正在停止进程 PID: $PID${NC}"
        kill -TERM $PID 2>/dev/null
        
        # 等待进程结束
        for i in {1..10}; do
            if ! kill -0 $PID 2>/dev/null; then
                echo -e "${GREEN}进程 $PID 已停止${NC}"
                break
            fi
            sleep 1
        done
        
        # 如果进程还在运行，强制杀死
        if kill -0 $PID 2>/dev/null; then
            echo -e "${YELLOW}进程 $PID 仍在运行，强制停止...${NC}"
            kill -KILL $PID 2>/dev/null
            echo -e "${GREEN}进程 $PID 已强制停止${NC}"
        fi
    done
    
    # 再次检查端口
    sleep 2
    if lsof -Pi :$PORT -sTCP:LISTEN -t >/dev/null 2>&1; then
        echo -e "${RED}警告: 端口 ${PORT} 仍被占用${NC}"
        lsof -i :$PORT
    else
        echo -e "${GREEN}服务已成功停止${NC}"
    fi
}

# 清理临时文件
cleanup() {
    echo -e "${YELLOW}清理临时文件...${NC}"
    
    # 删除编译生成的可执行文件
    if [ -f "main" ]; then
        rm -f main
        echo -e "${GREEN}已删除可执行文件${NC}"
    fi
    
    # 自动清理Go缓存
    go clean -cache -modcache -testcache
    echo -e "${GREEN}Go缓存已清理${NC}"
}

# 主函数
main() {
    if check_port; then
        stop_service
    fi
    
    # 自动清理临时文件
    cleanup
    
    echo -e "${GREEN}操作完成${NC}"
}

# 执行主函数
main
