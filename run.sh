#!/bin/bash

# 统一的项目管理脚本
# 自动检查服务状态，如果运行中就停止，然后启动服务
# 
# 环境配置说明:
# - ENV=dev: 使用 dev_configs.toml 配置文件
# - GO_ENV=dev: Go环境变量，确保使用正确的配置
# - 这样配置会使用腾讯云Redis: gz-crs-jaz1h340.sql.tencentcdb.com:28627

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

# 环境配置 - 使用腾讯云Redis
export ENV="dev"
export GO_ENV="dev"

echo -e "${BLUE}================================${NC}"
echo -e "${BLUE}    ${PROJECT_NAME} 统一管理脚本${NC}"
echo -e "${BLUE}================================${NC}"

# 检查端口是否被占用
check_port() {
    if lsof -Pi :$PORT -sTCP:LISTEN -t >/dev/null 2>&1; then
        return 0  # 端口被占用
    else
        return 1  # 端口可用
    fi
}

# 停止服务
stop_service() {
    echo -e "${YELLOW}检测到服务正在运行，正在停止...${NC}"
    
    # 获取占用端口的进程PID
    PIDS=$(lsof -Pi :$PORT -sTCP:LISTEN -t 2>/dev/null)
    
    if [ -n "$PIDS" ]; then
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
        
        # 等待端口释放
        echo -e "${YELLOW}等待端口释放...${NC}"
        for i in {1..10}; do
            if ! lsof -Pi :$PORT -sTCP:LISTEN -t >/dev/null 2>&1; then
                echo -e "${GREEN}端口已释放${NC}"
                break
            fi
            sleep 1
        done
        
        # 清理临时文件
        echo -e "${YELLOW}清理临时文件...${NC}"
        if [ -f "main" ]; then
            rm -f main
            echo -e "${GREEN}已删除可执行文件${NC}"
        fi
        
        # 清理Go缓存
        go clean -cache -modcache -testcache >/dev/null 2>&1
        echo -e "${GREEN}Go缓存已清理${NC}"
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
    echo -e "${YELLOW}当前环境配置:${NC}"
    echo -e "${BLUE}  - ENV: $ENV${NC}"
    echo -e "${BLUE}  - GO_ENV: $GO_ENV${NC}"
    echo -e "${BLUE}  - Redis: 腾讯云 (gz-crs-jaz1h340.sql.tencentcdb.com:28627)${NC}"
    echo -e "${BLUE}  - 配置文件: dev_configs.toml${NC}"
    echo ""
    echo -e "${YELLOW}注意: 将使用 -rebuild-db -force-reseed 参数启动${NC}"
    echo -e "${YELLOW}这将重建所有数据库表并插入mock数据${NC}"
    echo ""
    echo -e "${GREEN}按 Ctrl+C 停止服务${NC}"
    echo ""
    
    # 启动服务（重建数据库并插入mock数据）
    ./main -rebuild-db -force-reseed
}

# 主函数
main() {
    # 检查服务状态
    if check_port; then
        echo -e "${YELLOW}检测到服务正在运行${NC}"
        stop_service
        echo -e "${GREEN}服务已停止${NC}"
        echo ""
    else
        echo -e "${GREEN}服务未运行，直接启动${NC}"
        echo ""
    fi
    
    # 启动服务
    check_prerequisites
    download_deps
    build_project
    start_service
}

# 执行主函数
main
