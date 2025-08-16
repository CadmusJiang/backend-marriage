#!/bin/bash

# 数据库强制清理脚本
# 提供更高级的数据库清理功能，包括强制断开连接、清理缓存、重置状态等

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
PURPLE='\033[0;35m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# 数据库配置（从配置文件读取）
DB_HOST="gz-cdb-pepynap7.sql.tencentcdb.com"
DB_PORT="63623"
DB_NAME="marriage_system"
DB_USER="marriage_system"
DB_PASS="19970901Zyt"

# 显示帮助信息
show_help() {
    echo -e "${BLUE}================================${NC}"
    echo -e "${BLUE}    数据库强制清理脚本${NC}"
    echo -e "${BLUE}================================${NC}"
    echo ""
    echo -e "${YELLOW}用法:${NC}"
    echo -e "  $0 [选项]"
    echo ""
    echo -e "${YELLOW}选项:${NC}"
    echo -e "  -h, --help              显示此帮助信息"
    echo -e "  -f, --force             强制执行所有清理操作"
    echo -e "  -c, --connections       强制断开所有数据库连接"
    echo -e "  -k, --cache             清理数据库缓存"
    echo -e "  -s, --status            重置数据库状态"
    echo -e "  -t, --temp              清理临时表"
    echo -e "  -o, --optimize          优化数据库表"
    echo -e "  -i, --increment         重置所有表的自增ID"
    echo -e "  -a, --all               执行所有清理操作（默认）"
    echo ""
    echo -e "${YELLOW}示例:${NC}"
    echo -e "  $0                    # 执行所有清理操作"
    echo -e "  $0 -c                 # 只断开连接"
    echo -e "  $0 -i                 # 只重置自增ID"
    echo -e "  $0 -f -a              # 强制执行所有操作"
    echo ""
}

# 检查MySQL客户端
check_mysql_client() {
    if ! command -v mysql &> /dev/null; then
        echo -e "${RED}错误: mysql客户端未安装${NC}"
        echo -e "${YELLOW}请安装mysql-client:${NC}"
        echo -e "${BLUE}  macOS: brew install mysql-client${NC}"
        echo -e "${BLUE}  Ubuntu/Debian: sudo apt-get install mysql-client${NC}"
        echo -e "${BLUE}  CentOS/RHEL: sudo yum install mysql-client${NC}"
        exit 1
    fi
    
    # 检查MySQL版本兼容性
    local mysql_version=$(mysql --version 2>/dev/null | grep -o 'mysql [0-9]\+\.[0-9]\+' | cut -d' ' -f2)
    if [[ "$mysql_version" == "9.4" ]]; then
        echo -e "${YELLOW}警告: 检测到MySQL 9.4版本，可能存在兼容性问题${NC}"
        echo -e "${YELLOW}建议使用MySQL 8.0或更早版本以获得更好的兼容性${NC}"
        echo -e "${BLUE}可以尝试: brew install mysql-client@8.0${NC}"
    fi
}

# 测试数据库连接
test_connection() {
    echo -e "${YELLOW}测试数据库连接...${NC}"
    
    # 尝试不同的连接方式
    local connection_methods=(
        "mysql -h\"$DB_HOST\" -P\"$DB_PORT\" -u\"$DB_USER\" -p\"$DB_PASS\" \"$DB_NAME\" -e \"SELECT 1;\""
        "mysql -h\"$DB_HOST\" -P\"$DB_PORT\" -u\"$DB_USER\" -p\"$DB_PASS\" \"$DB_NAME\" --ssl-mode=DISABLED -e \"SELECT 1;\""
        "mysql -h\"$DB_HOST\" -P\"$DB_PORT\" -u\"$DB_USER\" -p\"$DB_PASS\" \"$DB_NAME\" --protocol=TCP -e \"SELECT 1;\""
    )
    
    for method in "${connection_methods[@]}"; do
        echo -e "${BLUE}尝试连接方式: $method${NC}"
        if eval "$method" >/dev/null 2>&1; then
            echo -e "${GREEN}数据库连接成功${NC}"
            # 保存成功的连接方式
            MYSQL_CONNECTION_CMD="$method"
            return 0
        fi
    done
    
    echo -e "${RED}所有连接方式都失败了${NC}"
    echo -e "${YELLOW}可能的原因:${NC}"
    echo -e "${BLUE}  1. MySQL 9.4版本兼容性问题${NC}"
    echo -e "${BLUE}  2. 网络连接问题${NC}"
    echo -e "${BLUE}  3. 数据库配置问题${NC}"
    echo -e "${BLUE}  4. 防火墙或安全组设置${NC}"
    
    # 提供解决方案
    echo -e "${YELLOW}建议解决方案:${NC}"
    echo -e "${BLUE}  1. 安装MySQL 8.0客户端: brew install mysql-client@8.0${NC}"
    echo -e "${BLUE}  2. 检查网络连接: ping $DB_HOST${NC}"
    echo -e "${BLUE}  3. 检查数据库配置和权限${NC}"
    
    return 1
}

# 执行MySQL命令的通用函数
execute_mysql_command() {
    local sql_command="$1"
    local description="$2"
    
    if [ -z "$MYSQL_CONNECTION_CMD" ]; then
        echo -e "${RED}错误: 没有可用的数据库连接${NC}"
        return 1
    fi
    
    # 构建完整的命令
    local full_cmd="${MYSQL_CONNECTION_CMD% -e *} -e \"$sql_command\""
    
    echo -e "${BLUE}  执行: $description${NC}"
    if eval "$full_cmd" >/dev/null 2>&1; then
        echo -e "${GREEN}  $description 成功${NC}"
        return 0
    else
        echo -e "${YELLOW}  $description 失败（可能是权限问题，这是正常的）${NC}"
        return 1
    fi
}

# 强制断开所有连接
force_kill_connections() {
    echo -e "${YELLOW}1. 强制断开所有数据库连接...${NC}"
    
    # 获取所有连接（除了当前连接）
    local connections=$(mysql -h"$DB_HOST" -P"$DB_PORT" -u"$DB_USER" -p"$DB_PASS" "$DB_NAME" -e "
        SELECT id, user, host, db, command, time, state, info 
        FROM information_schema.processlist 
        WHERE db = '$DB_NAME' AND id != CONNECTION_ID();
    " 2>/dev/null)
    
    if [ -n "$connections" ]; then
        echo -e "${BLUE}找到以下活动连接:${NC}"
        echo "$connections" | head -20
        
        # 强制断开连接
        local kill_commands=$(mysql -h"$DB_HOST" -P"$DB_PORT" -u"$DB_USER" -p"$DB_PASS" "$DB_NAME" -e "
            SELECT CONCAT('KILL ', id, ';') 
            FROM information_schema.processlist 
            WHERE db = '$DB_NAME' AND id != CONNECTION_ID();
        " 2>/dev/null | grep -v "CONCAT")
        
        if [ -n "$kill_commands" ]; then
            echo -e "${YELLOW}正在断开连接...${NC}"
            echo "$kill_commands" | while read kill_cmd; do
                if [ -n "$kill_cmd" ]; then
                    echo -e "${BLUE}  执行: $kill_cmd${NC}"
                    mysql -h"$DB_HOST" -P"$DB_PORT" -u"$DB_USER" -p"$DB_PASS" "$DB_NAME" -e "$kill_cmd" 2>/dev/null || true
                fi
            done
            echo -e "${GREEN}  所有连接已断开${NC}"
        else
            echo -e "${GREEN}  没有需要断开的连接${NC}"
        fi
    else
        echo -e "${GREEN}  没有活动连接${NC}"
    fi
}

# 清理数据库缓存
clean_cache() {
    echo -e "${YELLOW}2. 清理数据库缓存...${NC}"
    
    local flush_commands=(
        "FLUSH PRIVILEGES"
        "FLUSH HOSTS"
        "FLUSH LOGS"
        "FLUSH STATUS"
        "FLUSH TABLES"
        "FLUSH TABLES WITH READ LOCK"
        "UNLOCK TABLES"
    )
    
    for cmd in "${flush_commands[@]}"; do
        echo -e "${BLUE}  执行: $cmd${NC}"
        mysql -h"$DB_HOST" -P"$DB_PORT" -u"$DB_USER" -p"$DB_PASS" "$DB_NAME" -e "$cmd" 2>/dev/null || true
    done
    
    echo -e "${GREEN}  数据库缓存已清理${NC}"
}

# 重置数据库状态
reset_status() {
    echo -e "${YELLOW}3. 重置数据库状态...${NC}"
    
    local reset_commands=(
        "SET GLOBAL innodb_force_recovery = 0"
        "SET GLOBAL innodb_fast_shutdown = 1"
        "SET GLOBAL query_cache_size = 0"
        "SET GLOBAL query_cache_type = 0"
    )
    
    for cmd in "${reset_commands[@]}"; do
        echo -e "${BLUE}  执行: $cmd${NC}"
        mysql -h"$DB_HOST" -P"$DB_PORT" -u"$DB_USER" -p"$DB_PASS" "$DB_NAME" -e "$cmd" 2>/dev/null || true
    done
    
    echo -e "${GREEN}  数据库状态已重置${NC}"
}

# 清理临时表
clean_temp_tables() {
    echo -e "${YELLOW}4. 清理临时表...${NC}"
    
    # 查找临时表
    local temp_tables=$(mysql -h"$DB_HOST" -P"$DB_PORT" -u"$DB_USER" -p"$DB_PASS" "$DB_NAME" -e "
        SELECT table_name 
        FROM information_schema.tables 
        WHERE table_schema = '$DB_NAME' 
        AND (table_name LIKE 'temp_%' OR table_name LIKE '%_temp' OR table_name LIKE 'tmp_%');
    " 2>/dev/null | grep -v "table_name")
    
    if [ -n "$temp_tables" ]; then
        echo -e "${BLUE}找到以下临时表:${NC}"
        echo "$temp_tables"
        
        echo "$temp_tables" | while read table_name; do
            if [ -n "$table_name" ]; then
                echo -e "${BLUE}  删除表: $table_name${NC}"
                mysql -h"$DB_HOST" -P"$DB_PORT" -u"$DB_USER" -p"$DB_PASS" "$DB_NAME" -e "DROP TABLE IF EXISTS \`$table_name\`;" 2>/dev/null || true
            fi
        done
        
        echo -e "${GREEN}  临时表已清理${NC}"
    else
        echo -e "${GREEN}  没有找到临时表${NC}"
    fi
}

# 优化数据库表
optimize_tables() {
    echo -e "${YELLOW}5. 优化数据库表...${NC}"
    
    # 获取所有表
    local tables=$(mysql -h"$DB_HOST" -P"$DB_PORT" -u"$DB_USER" -p"$DB_PASS" "$DB_NAME" -e "
        SELECT table_name 
        FROM information_schema.tables 
        WHERE table_schema = '$DB_NAME' 
        AND table_type = 'BASE TABLE';
    " 2>/dev/null | grep -v "table_name")
    
    if [ -n "$tables" ]; then
        echo -e "${BLUE}找到 ${#tables[@]} 个表，开始优化...${NC}"
        
        local count=0
        echo "$tables" | while read table_name; do
            if [ -n "$table_name" ]; then
                count=$((count + 1))
                echo -e "${BLUE}  优化表 [$count]: $table_name${NC}"
                mysql -h"$DB_HOST" -P"$DB_PORT" -u"$DB_USER" -p"$DB_PASS" "$DB_NAME" -e "OPTIMIZE TABLE \`$table_name\`;" 2>/dev/null || true
            fi
        done
        
        echo -e "${GREEN}  表优化完成${NC}"
    else
        echo -e "${GREEN}  没有找到需要优化的表${NC}"
    fi
}

# 显示数据库状态
show_status() {
    echo -e "${YELLOW}6. 显示数据库状态...${NC}"
    
    echo -e "${BLUE}数据库连接信息:${NC}"
    mysql -h"$DB_HOST" -P"$DB_PORT" -u"$DB_USER" -p"$DB_PASS" "$DB_NAME" -e "
        SELECT 
            @@hostname as '主机名',
            @@port as '端口',
            @@version as '版本',
            @@datadir as '数据目录',
            @@innodb_version as 'InnoDB版本';
    " 2>/dev/null || echo -e "${YELLOW}  无法获取数据库状态${NC}"
    
    echo -e "${BLUE}当前连接数:${NC}"
    mysql -h"$DB_HOST" -P"$DB_PORT" -u"$DB_USER" -p"$DB_PASS" "$DB_NAME" -e "
        SELECT COUNT(*) as '连接数' FROM information_schema.processlist WHERE db = '$DB_NAME';
    " 2>/dev/null || echo -e "${YELLOW}  无法获取连接数${NC}"
}

# 重置自增ID
reset_auto_increment() {
    echo -e "${YELLOW}6. 重置所有表的自增ID...${NC}"
    
    if [ -z "$MYSQL_CONNECTION_CMD" ]; then
        echo -e "${RED}错误: 没有可用的数据库连接${NC}"
        return 1
    fi
    
    # 获取所有表名
    local tables=$(eval "${MYSQL_CONNECTION_CMD% -e *} -e \"SHOW TABLES;\" 2>/dev/null" | grep -v "Tables_in_$DB_NAME")
    
    if [ -n "$tables" ]; then
        echo -e "${BLUE}找到以下表:${NC}"
        echo "$tables"
        echo ""
        
        local count=0
        echo "$tables" | while read table_name; do
            if [ -n "$table_name" ]; then
                count=$((count + 1))
                echo -e "${BLUE}  重置表 [$count]: $table_name${NC}"
                
                # 重置AUTO_INCREMENT值
                eval "${MYSQL_CONNECTION_CMD% -e *} -e \"ALTER TABLE \`$table_name\` AUTO_INCREMENT = 1;\" 2>/dev/null" && \
                    echo -e "${GREEN}    ✓ $table_name 自增ID已重置${NC}" || \
                    echo -e "${YELLOW}    ⚠ $table_name 自增ID重置失败（可能是权限问题）${NC}"
            fi
        done
        
        echo -e "${GREEN}  自增ID重置完成${NC}"
        
        # 验证重置结果
        echo -e "${BLUE}  验证重置结果...${NC}"
        eval "${MYSQL_CONNECTION_CMD% -e *} -e \"
            SELECT 
                TABLE_NAME as '表名',
                AUTO_INCREMENT as '当前自增值'
            FROM information_schema.TABLES 
            WHERE TABLE_SCHEMA = '$DB_NAME' 
            AND AUTO_INCREMENT IS NOT NULL
            ORDER BY TABLE_NAME;
        \" 2>/dev/null" || echo -e "${YELLOW}    无法获取自增值信息（可能是权限问题）${NC}"
        
    else
        echo -e "${GREEN}  没有找到表，跳过自增ID重置${NC}"
    fi
}

# 主清理函数
main_cleanup() {
    local force_mode=false
    local operations=()
    
    # 解析命令行参数
    while [[ $# -gt 0 ]]; do
        case $1 in
            -h|--help)
                show_help
                exit 0
                ;;
            -f|--force)
                force_mode=true
                shift
                ;;
            -c|--connections)
                operations+=("connections")
                shift
                ;;
            -k|--cache)
                operations+=("cache")
                shift
                ;;
            -s|--status)
                operations+=("status")
                shift
                ;;
            -t|--temp)
                operations+=("temp")
                shift
                ;;
            -o|--optimize)
                operations+=("optimize")
                shift
                ;;
            -i|--increment)
                operations+=("increment")
                shift
                ;;
            -a|--all)
                operations=("connections" "cache" "status" "temp" "optimize" "increment")
                shift
                ;;
            *)
                echo -e "${RED}未知选项: $1${NC}"
                show_help
                exit 1
                ;;
        esac
    done
    
    # 如果没有指定操作，默认执行所有
    if [ ${#operations[@]} -eq 0 ]; then
        operations=("connections" "cache" "status" "temp" "optimize" "increment")
    fi
    
    echo -e "${BLUE}================================${NC}"
    echo -e "${BLUE}    开始数据库强制清理${NC}"
    echo -e "${BLUE}================================${NC}"
    echo -e "${BLUE}目标数据库: ${DB_HOST}:${DB_PORT}/${DB_NAME}${NC}"
    echo -e "${BLUE}执行操作: ${operations[*]}${NC}"
    echo -e "${BLUE}强制模式: $force_mode${NC}"
    echo ""
    
    # 检查前置条件
    check_mysql_client
    
    if ! test_connection; then
        echo -e "${RED}无法连接到数据库，请检查配置${NC}"
        exit 1
    fi
    
    # 确认操作
    if [ "$force_mode" = false ]; then
        echo -e "${YELLOW}警告: 即将执行数据库清理操作${NC}"
        echo -e "${YELLOW}这些操作可能会影响数据库性能${NC}"
        read -p "是否继续? (y/N): " -n 1 -r
        echo
        if [[ ! $REPLY =~ ^[Yy]$ ]]; then
            echo -e "${YELLOW}操作已取消${NC}"
            exit 0
        fi
    fi
    
    # 执行清理操作
    for op in "${operations[@]}"; do
        case $op in
            "connections")
                force_kill_connections
                ;;
            "cache")
                clean_cache
                ;;
            "status")
                reset_status
                ;;
            "temp")
                clean_temp_tables
                ;;
            "optimize")
                optimize_tables
                ;;
            "increment")
                reset_auto_increment
                ;;
        esac
        echo ""
    done
    
    # 显示最终状态
    show_status
    
    echo -e "${GREEN}================================${NC}"
    echo -e "${GREEN}    数据库清理完成！${NC}"
    echo -e "${GREEN}================================${NC}"
}

# 执行主函数
main_cleanup "$@"
