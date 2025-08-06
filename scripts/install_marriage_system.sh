#!/bin/bash

# Marriage System 安装脚本
# 自动导入数据库和启动应用

echo "=== Marriage System 安装脚本 ==="

# 检查配置文件是否存在
if [ ! -f "configs/dev_configs.toml" ]; then
    echo "错误: 配置文件 configs/dev_configs.toml 不存在"
    exit 1
fi

# 从配置文件读取MySQL连接信息
echo "1. 读取数据库配置..."
MYSQL_HOST=$(grep -A 5 "\[mysql.read\]" configs/dev_configs.toml | grep "addr" | cut -d'"' -f2 | cut -d':' -f1)
MYSQL_PORT=$(grep -A 5 "\[mysql.read\]" configs/dev_configs.toml | grep "addr" | cut -d'"' -f2 | cut -d':' -f2)
MYSQL_USER=$(grep -A 5 "\[mysql.read\]" configs/dev_configs.toml | grep "user" | cut -d'"' -f2)
MYSQL_PASS=$(grep -A 5 "\[mysql.read\]" configs/dev_configs.toml | grep "pass" | cut -d'"' -f2)
MYSQL_DB=$(grep -A 5 "\[mysql.read\]" configs/dev_configs.toml | grep "name" | cut -d'"' -f2)

echo "数据库配置:"
echo "  主机: $MYSQL_HOST"
echo "  端口: $MYSQL_PORT"
echo "  用户: $MYSQL_USER"
echo "  数据库: $MYSQL_DB"
echo "  密码: ${MYSQL_PASS:0:3}***"

# 检查MySQL连接
echo "2. 测试MySQL连接..."
echo "正在连接远程数据库: $MYSQL_HOST:$MYSQL_PORT"

# 测试网络连接
echo "  检查网络连接..."
if ! nc -z $MYSQL_HOST $MYSQL_PORT 2>/dev/null; then
    echo "❌ 网络连接失败: 无法连接到 $MYSQL_HOST:$MYSQL_PORT"
    echo "请检查:"
    echo "1. 网络连接是否正常"
    echo "2. 防火墙是否阻止连接"
    echo "3. 数据库服务器是否运行"
    exit 1
fi
echo "✅ 网络连接正常"

# 测试数据库连接
echo "  测试数据库连接..."
if mysql -h $MYSQL_HOST -P $MYSQL_PORT -u $MYSQL_USER -p$MYSQL_PASS -e "SELECT 1;" > /dev/null 2>&1; then
    echo "✅ MySQL连接成功"
else
    echo "❌ MySQL连接失败"
    echo "请检查以下配置:"
    echo "  主机: $MYSQL_HOST"
    echo "  端口: $MYSQL_PORT"
    echo "  用户: $MYSQL_USER"
    echo "  密码: ${MYSQL_PASS:0:3}***"
    echo ""
    echo "可能的问题:"
    echo "1. 用户名或密码错误"
    echo "2. 数据库用户权限不足"
    echo "3. 数据库服务器配置问题"
    echo "4. 配置文件中的连接信息错误"
    echo ""
    echo "建议:"
    echo "1. 检查配置文件中的数据库连接信息"
    echo "2. 确认数据库用户有足够的权限"
    echo "3. 联系数据库管理员确认连接信息"
    exit 1
fi

# 导入数据库
echo "3. 导入数据库..."
if sed "s/__DB_NAME__/$MYSQL_DB/g" scripts/init_marriage_system_database.sql | mysql -h $MYSQL_HOST -P $MYSQL_PORT -u $MYSQL_USER -p$MYSQL_PASS; then
    echo "✅ 数据库导入成功"
else
    echo "❌ 数据库导入失败"
    exit 1
fi

# 检查Go环境
echo "4. 检查Go环境..."
if ! command -v go &> /dev/null; then
    echo "错误: Go未安装，请先安装Go"
    exit 1
fi

# 安装依赖
echo "5. 安装Go依赖..."
go mod tidy
if [ $? -eq 0 ]; then
    echo "✅ 依赖安装成功"
else
    echo "❌ 依赖安装失败"
    exit 1
fi

# 生成数据库模型代码
echo "6. 生成数据库模型代码..."
go run cmd/gormgen/main.go -structs Account -input ./internal/repository/mysql/account/
go run cmd/gormgen/main.go -structs AccountHistory -input ./internal/repository/mysql/account_history/
go run cmd/gormgen/main.go -structs OrgInfo -input ./internal/repository/mysql/org_info/

if [ $? -eq 0 ]; then
    echo "✅ 数据库模型代码生成成功"
else
    echo "❌ 数据库模型代码生成失败"
    exit 1
fi

# 编译应用
echo "7. 编译应用..."
go build -o marriage_system main.go
if [ $? -eq 0 ]; then
    echo "✅ 应用编译成功"
else
    echo "❌ 应用编译失败"
    exit 1
fi

# 启动应用
echo "8. 启动应用..."
echo "应用将在后台运行，端口: 9999"
echo "可以使用以下命令查看日志:"
echo "tail -f logs/marriage_system.log"
echo ""
echo "测试API:"
echo "./scripts/test_account_api_real.sh"
echo "./scripts/test_org_info_api.sh"
echo ""

# 后台启动应用
nohup ./marriage_system > logs/marriage_system.log 2>&1 &
PID=$!
echo $PID > marriage_system.pid

echo "✅ Marriage System 安装完成！"
echo "应用PID: $PID"
echo "访问地址: http://localhost:9999"
echo ""
echo "数据库信息:"
echo "- 数据库名: $MYSQL_DB"
echo "- 数据库主机: $MYSQL_HOST:$MYSQL_PORT"
echo "- 账户表: account"
echo "- 组织表: org"
echo "- 历史记录表: account_history"
echo ""
echo "测试账户:"
echo "- 用户名: admin, 密码: 123456"
echo "- 用户名: company_manager, 密码: 123456"
echo "- 用户名: group_manager, 密码: 123456"
echo "- 用户名: team_manager, 密码: 123456"
echo "- 用户名: employee, 密码: 123456" 