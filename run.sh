#!/bin/bash

# 统一的项目管理脚本
# 自动检查服务状态，如果运行中就停止，然后启动服务
# 
# 环境配置说明:
# - ENV=dev: 使用 dev_configs.toml 配置文件
# - GO_ENV=dev: Go环境变量，确保使用正确的配置
# - 这样配置会使用腾讯云Redis: gz-crs-jaz1h340.sql.tencentcdb.com:28627
#
# 数据库重建说明:
# - 每次启动都会使用 -rebuild-db -force-reseed 参数
# - 这将删除所有现有表并重新创建
# - 确保每次启动都有干净的数据库和一致的mock数据

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

# 数据库配置（从配置文件读取）
DB_HOST="gz-cdb-pepynap7.sql.tencentcdb.com"
DB_PORT="63623"
DB_NAME="marriage_system"
DB_USER="marriage_system"
DB_PASS="19970901Zyt"

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

# 强制数据库清理函数
force_database_cleanup() {
    echo -e "${YELLOW}开始强制数据库清理...${NC}"
    
    # 检查是否安装了mysql客户端
    if ! command -v mysql &> /dev/null; then
        echo -e "${YELLOW}警告: mysql客户端未安装，跳过数据库清理${NC}"
        echo -e "${YELLOW}请安装mysql-client: brew install mysql-client (macOS) 或 apt-get install mysql-client (Ubuntu)${NC}"
        return 0
    fi
    
    echo -e "${BLUE}正在连接数据库: ${DB_HOST}:${DB_PORT}${NC}"
    
    # 强制断开所有连接
    echo -e "${YELLOW}1. 强制断开所有数据库连接...${NC}"
    mysql -h"$DB_HOST" -P"$DB_PORT" -u"$DB_USER" -p"$DB_PASS" "$DB_NAME" -e "
        SELECT CONCAT('KILL ', id, ';') 
        FROM information_schema.processlist 
        WHERE db = '$DB_NAME' AND id != CONNECTION_ID();
    " 2>/dev/null | grep -v "CONCAT" | while read kill_cmd; do
        if [ -n "$kill_cmd" ]; then
            echo -e "${BLUE}  执行: $kill_cmd${NC}"
            mysql -h"$DB_HOST" -P"$DB_PORT" -u"$DB_USER" -p"$DB_PASS" "$DB_NAME" -e "$kill_cmd" 2>/dev/null || true
        fi
    done
    
    # 清理数据库缓存
    echo -e "${YELLOW}2. 清理数据库缓存...${NC}"
    mysql -h"$DB_HOST" -P"$DB_PORT" -u"$DB_USER" -p"$DB_PASS" "$DB_NAME" -e "
        FLUSH PRIVILEGES;
        FLUSH HOSTS;
        FLUSH LOGS;
        FLUSH STATUS;
        FLUSH TABLES;
        FLUSH TABLES WITH READ LOCK;
        UNLOCK TABLES;
    " 2>/dev/null && echo -e "${GREEN}  数据库缓存已清理${NC}" || echo -e "${YELLOW}  缓存清理完成（部分操作可能失败，这是正常的）${NC}"
    
    # 重置数据库状态
    echo -e "${YELLOW}3. 重置数据库状态...${NC}"
    mysql -h"$DB_HOST" -P"$DB_PORT" -u"$DB_USER" -p"$DB_PASS" "$DB_NAME" -e "
        SET GLOBAL innodb_force_recovery = 0;
        SET GLOBAL innodb_fast_shutdown = 1;
    " 2>/dev/null && echo -e "${GREEN}  数据库状态已重置${NC}" || echo -e "${YELLOW}  状态重置完成（部分操作可能失败，这是正常的）${NC}"
    
    # 清理临时表
    echo -e "${YELLOW}4. 清理临时表...${NC}"
    mysql -h"$DB_HOST" -P"$DB_PORT" -u"$DB_USER" -p"$DB_PASS" "$DB_NAME" -e "
        DROP TEMPORARY TABLE IF EXISTS temp_*;
        DROP TABLE IF EXISTS temp_*;
    " 2>/dev/null && echo -e "${GREEN}  临时表已清理${NC}" || echo -e "${YELLOW}  临时表清理完成${NC}"
    
    # 优化表
    echo -e "${YELLOW}5. 优化数据库表...${NC}"
    mysql -h"$DB_HOST" -P"$DB_PORT" -u"$DB_USER" -p"$DB_PASS" "$DB_NAME" -e "
        OPTIMIZE TABLE information_schema.tables;
    " 2>/dev/null && echo -e "${GREEN}  表优化完成${NC}" || echo -e "${YELLOW}  表优化完成（部分操作可能失败，这是正常的）${NC}"
    
    echo -e "${GREEN}强制数据库清理完成！${NC}"
    echo ""
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
        
        # 强制数据库清理
        force_database_cleanup
        
        # 数据库清理提示
        echo -e "${YELLOW}数据库清理提示:${NC}"
        echo -e "${BLUE}  下次启动时将完全重建数据库${NC}"
        echo -e "${BLUE}  所有表将被删除并重新创建${NC}"
        echo -e "${BLUE}  新的mock数据将被插入${NC}"
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
    echo -e "${BLUE}  - 数据库: ${DB_HOST}:${DB_PORT}${NC}"
    echo ""
    echo -e "${YELLOW}注意: 将使用 -rebuild-db -force-reseed 参数启动${NC}"
    echo -e "${YELLOW}这将重建所有数据库表并插入mock数据${NC}"
    echo ""
    
    # 数据库清理提示
    echo -e "${YELLOW}数据库重建说明:${NC}"
    echo -e "${BLUE}  - 删除所有现有表${NC}"
    echo -e "${BLUE}  - 重新创建表结构${NC}"
    echo -e "${BLUE}  - 重置所有表的自增ID（Go程序自动处理）${NC}"
    echo -e "${BLUE}  - 插入新的mock数据（ID从1开始）${NC}"
    echo -e "${BLUE}  - 确保数据一致性${NC}"
    echo ""
    
    echo -e "${GREEN}按 Ctrl+C 停止服务${NC}"
    echo ""
    
    # 启动服务（重建数据库并插入mock数据）
    # Go程序会自动处理AUTO_INCREMENT重置，确保ID从1开始
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
