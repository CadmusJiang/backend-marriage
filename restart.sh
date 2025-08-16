#!/bin/bash

# 项目重启脚本
# 停止服务并重新启动

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
echo -e "${BLUE}    ${PROJECT_NAME} 重启脚本${NC}"
echo -e "${BLUE}================================${NC}"

# 停止服务
stop_service() {
    echo -e "${YELLOW}正在停止服务...${NC}"
    
    # 检查是否有服务在运行
    if lsof -Pi :$PORT -sTCP:LISTEN -t >/dev/null 2>&1; then
        # 获取占用端口的进程PID
        PIDS=$(lsof -Pi :$PORT -sTCP:LISTEN -t 2>/dev/null)
        
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
        
        # 等待端口释放
        echo -e "${YELLOW}等待端口释放...${NC}"
        for i in {1..10}; do
            if ! lsof -Pi :$PORT -sTCP:LISTEN -t >/dev/null 2>&1; then
                echo -e "${GREEN}端口已释放${NC}"
                break
            fi
            sleep 1
        done
    else
        echo -e "${YELLOW}没有发现运行中的服务${NC}"
    fi
}

# 启动服务
start_service() {
    echo -e "${YELLOW}正在启动服务...${NC}"
    
    # 检查端口是否可用
    if lsof -Pi :$PORT -sTCP:LISTEN -t >/dev/null 2>&1; then
        echo -e "${RED}错误: 端口 ${PORT} 仍被占用${NC}"
        return 1
    fi
    
    # 检查必要的文件
    if [ ! -f "go.mod" ]; then
        echo -e "${RED}错误: 未找到 go.mod 文件${NC}"
        return 1
    fi
    
    # 创建日志目录
    if [ ! -d "./logs" ]; then
        mkdir -p "./logs"
    fi
    
    # 下载依赖
    echo -e "${YELLOW}下载项目依赖...${NC}"
    go mod download
    if [ $? -ne 0 ]; then
        echo -e "${RED}依赖下载失败${NC}"
        return 1
    fi
    
    # 构建项目
    echo -e "${YELLOW}构建项目...${NC}"
    go build -o main .
    if [ $? -ne 0 ]; then
        echo -e "${RED}项目构建失败${NC}"
        return 1
    fi
    
    # 启动服务
    echo -e "${YELLOW}启动服务...${NC}"
    echo -e "${BLUE}服务将在端口 ${PORT} 上启动${NC}"
    echo -e "${BLUE}访问地址: http://localhost:${PORT}${NC}"
    echo -e "${BLUE}API文档: http://localhost:${PORT}/docs/swagger.html${NC}"
    echo ""
    echo -e "${GREEN}按 Ctrl+C 停止服务${NC}"
    echo ""
    
    # 启动服务
    ./main
}

# 主函数
main() {
    stop_service
    sleep 2
    start_service
}

# 执行主函数
main
