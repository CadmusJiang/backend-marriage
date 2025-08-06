#!/bin/bash

# 数据库初始化脚本
# 用于删除旧表并创建新表结构

echo "开始初始化数据库..."

# 检查MySQL连接
mysql -u root -p -e "SELECT 1;" > /dev/null 2>&1
if [ $? -ne 0 ]; then
    echo "错误: 无法连接到MySQL数据库，请检查连接配置"
    exit 1
fi

# 执行SQL脚本
echo "执行数据库初始化脚本..."
mysql -u root -p < scripts/init_database.sql

if [ $? -eq 0 ]; then
    echo "数据库初始化成功！"
    echo "已删除旧表并创建新表结构"
    echo "已插入Mock数据"
else
    echo "错误: 数据库初始化失败"
    exit 1
fi

echo "数据库初始化完成！" 