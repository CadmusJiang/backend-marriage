#!/bin/bash

echo "🧪 测试日志解析逻辑"
echo "===================="

# 检查日志文件中的123456789记录
echo "1. 检查日志文件中的123456789记录数量..."
TOTAL_COUNT=$(grep -c "123456789" logs/marriage_system-access.log)
echo "总记录数: $TOTAL_COUNT"

# 检查API返回的记录数
echo ""
echo "2. 检查API返回的记录数..."
API_RESPONSE=$(curl -s "http://localhost:9999/api/v1/logs/trace?trace_id=123456789")
API_COUNT=$(echo "$API_RESPONSE" | jq -r '.data.data.logs | length' 2>/dev/null || echo "0")
echo "API返回记录数: $API_COUNT"

# 检查日志文件的行数
echo ""
echo "3. 检查日志文件信息..."
FILE_LINES=$(wc -l logs/marriage_system-access.log | awk '{print $1}')
FILE_SIZE=$(ls -lh logs/marriage_system-access.log | awk '{print $5}')
echo "文件行数: $FILE_LINES"
echo "文件大小: $FILE_SIZE"

# 检查最近的几条123456789记录
echo ""
echo "4. 检查最近的几条123456789记录..."
grep "123456789" logs/marriage_system-access.log | tail -3 | while read -r line; do
    echo "记录: ${line:0:100}..."
done

# 检查是否有损坏的日志行
echo ""
echo "5. 检查是否有损坏的日志行..."
grep "123456789" logs/marriage_system-access.log | grep -v "^{" | wc -l | awk '{print "非JSON格式记录数: " $1}'

echo ""
echo "测试完成！"
