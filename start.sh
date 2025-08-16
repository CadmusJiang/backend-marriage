#!/bin/bash

# 项目启动脚本
# 检查端口占用情况并启动服务

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 项目配置
PROJECT_NAME="marriage_system"
PORT="9999"
LOG_DIR="./logs"

echo -e "${BLUE}================================${NC}"
echo -e "${BLUE}    ${PROJECT_NAME} 启动脚本${NC}"
echo -e "${BLUE}================================${NC}"

# 检查端口是否被占用
check_port() {
    echo -e "${YELLOW}检查端口 ${PORT} 是否被占用...${NC}"
    
    if lsof -Pi :$PORT -sTCP:LISTEN -t >/dev/null 2>&1; then
        echo -e "${RED}错误: 端口 ${PORT} 已被占用${NC}"
        echo -e "${YELLOW}占用端口的进程信息:${NC}"
        lsof -i :$PORT
        echo ""
        echo -e "${YELLOW}请先停止占用端口的进程，或者修改配置文件中的端口号${NC}"
        exit 1
    else
        echo -e "${GREEN}端口 ${PORT} 可用${NC}"
    fi
}

# 检查必要的目录和文件
check_prerequisites() {
    echo -e "${YELLOW}检查项目依赖...${NC}"
    
    # 检查 Go 是否安装
    if ! command -v go &> /dev/null; then
        echo -e "${RED}错误: Go 未安装${NC}"
        exit 1
    fi
    
    # 检查 go.mod 文件
    if [ ! -f "go.mod" ]; then
        echo -e "${RED}错误: 未找到 go.mod 文件${NC}"
        exit 1
    fi
    
    # 创建日志目录
    if [ ! -d "$LOG_DIR" ]; then
        echo -e "${YELLOW}创建日志目录: $LOG_DIR${NC}"
        mkdir -p "$LOG_DIR"
    fi
    
    echo -e "${GREEN}项目依赖检查完成${NC}"
}

# 下载依赖
download_deps() {
    echo -e "${YELLOW}下载项目依赖...${NC}"
    go mod download
    if [ $? -eq 0 ]; then
        echo -e "${GREEN}依赖下载完成${NC}"
    else
        echo -e "${RED}依赖下载失败${NC}"
        exit 1
    fi
}

# 构建项目
build_project() {
    echo -e "${YELLOW}构建项目...${NC}"
    go build -o main .
    if [ $? -eq 0 ]; then
        echo -e "${GREEN}项目构建完成${NC}"
    else
        echo -e "${RED}项目构建失败${NC}"
        exit 1
    fi
}

# 启动服务
start_service() {
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
    check_prerequisites
    check_port
    download_deps
    build_project
    start_service
}

# 执行主函数
main
