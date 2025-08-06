#!/bin/bash

echo "=== 数据库连接测试 ==="

# 测试MySQL连接
echo "正在测试MySQL连接..."
mysql -h gz-cdb-pepynap7.sql.tencentcdb.com -P 3306 -u root -p19970901Zyt -e "SELECT 1;" 2>/dev/null
if [ $? -eq 0 ]; then
    echo "✅ MySQL连接成功"
else
    echo "❌ MySQL连接失败"
fi

# 测试Redis连接
echo "正在测试Redis连接..."
redis-cli -h gz-cdb-pepynap7.sql.tencentcdb.com -p 6379 -a 19970901Zyt ping 2>/dev/null
if [ $? -eq 0 ]; then
    echo "✅ Redis连接成功"
else
    echo "❌ Redis连接失败"
fi

echo "=== 连接测试完成 ===" 