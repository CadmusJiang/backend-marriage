#!/bin/bash

# 项目状态检查脚本
# 检查服务运行状态

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
echo -e "${BLUE}    ${PROJECT_NAME} 状态检查${NC}"
echo -e "${BLUE}================================${NC}"

# 检查端口状态
check_port_status() {
    echo -e "${YELLOW}检查端口 ${PORT} 状态...${NC}"
    
    if lsof -Pi :$PORT -sTCP:LISTEN -t >/dev/null 2>&1; then
        echo -e "${GREEN}✓ 服务正在运行${NC}"
        echo -e "${BLUE}进程信息:${NC}"
        lsof -i :$PORT
        return 0
    else
        echo -e "${RED}✗ 服务未运行${NC}"
        return 1
    fi
}

# 检查服务健康状态
check_health() {
    echo -e "${YELLOW}检查服务健康状态...${NC}"
    
    # 尝试访问服务
    if command -v curl &> /dev/null; then
        HTTP_CODE=$(curl -s -o /dev/null -w "%{http_code}" http://localhost:$PORT/ 2>/dev/null)
        if [ "$HTTP_CODE" = "200" ] || [ "$HTTP_CODE" = "404" ]; then
            echo -e "${GREEN}✓ 服务响应正常 (HTTP $HTTP_CODE)${NC}"
        else
            echo -e "${YELLOW}⚠ 服务响应异常 (HTTP $HTTP_CODE)${NC}"
        fi
    else
        echo -e "${YELLOW}⚠ curl 未安装，无法检查HTTP响应${NC}"
    fi
}

# 检查日志文件
check_logs() {
    echo -e "${YELLOW}检查日志文件...${NC}"
    
    LOG_FILE="./logs/${PROJECT_NAME}-access.log"
    
    if [ -f "$LOG_FILE" ]; then
        echo -e "${GREEN}✓ 日志文件存在: $LOG_FILE${NC}"
        
        # 显示最近的日志行数
        LOG_LINES=$(wc -l < "$LOG_FILE" 2>/dev/null || echo "0")
        echo -e "${BLUE}日志行数: $LOG_LINES${NC}"
        
        # 显示最近的日志
        if [ "$LOG_LINES" -gt 0 ]; then
            echo -e "${BLUE}最近的日志:${NC}"
            tail -5 "$LOG_FILE" 2>/dev/null | while IFS= read -r line; do
                echo -e "${YELLOW}  $line${NC}"
            done
        fi
    else
        echo -e "${YELLOW}⚠ 日志文件不存在: $LOG_FILE${NC}"
    fi
}

# 检查系统资源
check_resources() {
    echo -e "${YELLOW}检查系统资源...${NC}"
    
    # 检查内存使用
    if command -v ps &> /dev/null; then
        PIDS=$(lsof -Pi :$PORT -sTCP:LISTEN -t 2>/dev/null)
        if [ -n "$PIDS" ]; then
            for PID in $PIDS; do
                MEMORY=$(ps -o rss= -p $PID 2>/dev/null | awk '{print $1/1024 " MB"}')
                CPU=$(ps -o %cpu= -p $PID 2>/dev/null)
                echo -e "${BLUE}进程 $PID:${NC}"
                echo -e "${BLUE}  内存使用: $MEMORY${NC}"
                echo -e "${BLUE}  CPU使用: ${CPU}%${NC}"
            done
        fi
    fi
}

# 显示服务信息
show_service_info() {
    echo -e "${BLUE}服务信息:${NC}"
    echo -e "${BLUE}  项目名称: $PROJECT_NAME${NC}"
    echo -e "${BLUE}  端口: $PORT${NC}"
    echo -e "${BLUE}  访问地址: http://localhost:$PORT${NC}"
    echo -e "${BLUE}  API文档: http://localhost:$PORT/docs/swagger.html${NC}"
    echo ""
}

# 主函数
main() {
    show_service_info
    
    if check_port_status; then
        check_health
        check_resources
    fi
    
    check_logs
    
    echo ""
    echo -e "${BLUE}================================${NC}"
    echo -e "${BLUE}状态检查完成${NC}"
    echo -e "${BLUE}================================${NC}"
}

# 执行主函数
main
