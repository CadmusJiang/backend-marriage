# 分页日志系统使用说明

## 概述

为了解决日志文件越来越大时，一次性加载所有日志导致的性能问题，我们实现了一个按需加载的分页日志系统。该系统支持：

- 分页加载日志数据
- 向下滚动自动加载更多
- 实时统计信息
- 美观的Web界面

## 功能特性

### 1. 分页加载
- 支持自定义每页显示数量（10、20、50、100条）
- 自动计算总页数和当前页位置
- 支持手动跳转到指定页面

### 2. 按需加载
- 初始只加载第一页数据
- 向下滚动到底部时自动加载下一页
- 避免一次性加载大量数据导致的性能问题

### 3. 标签页分类
- **📋 全部日志**: 显示所有类型的日志（HTTP + SQL + Redis + 系统）
- **🌐 HTTP请求**: 仅显示HTTP请求相关的日志，包含方法、路径、状态码、耗时、Trace ID
- **🗄️ SQL操作**: 仅显示MySQL数据库操作相关的日志
- **🔴 Redis操作**: 仅显示Redis缓存操作相关的日志
- **⚙️ 系统日志**: 仅显示系统级别的日志信息

### 4. 实时统计
- 总日志数量
- HTTP请求数量
- MySQL操作数量
- Redis操作数量
- 系统日志数量

### 5. 日志分类
- HTTP请求日志（包含方法、路径、状态码、耗时、Trace ID）
- MySQL操作日志
- Redis操作日志
- 系统日志

## API接口

### 分页日志接口

**URL:** `GET /system/logs/paginated`

**查询参数:**
- `page`: 页码（默认：1）
- `page_size`: 每页大小（默认：10，最大：100）
- `type`: 日志类型过滤（可选值：`all`、`http`、`sql`、`redis`、`system`，默认：`all`）

**使用示例:**
```bash
# 获取第一页，每页5条，全部类型日志
curl "http://localhost:9999/system/logs/paginated?page=1&page_size=5&type=all"

# 获取第一页，每页10条，仅HTTP请求日志
curl "http://localhost:9999/system/logs/paginated?page=1&page_size=10&type=http"

# 获取第二页，每页20条，仅系统日志
curl "http://localhost:9999/system/logs/paginated?page=2&page_size=20&type=system"
```

**响应示例:**
```json
{
  "success": true,
  "data": {
    "data": {
      "logs": [
        {
          "raw": "原始日志内容",
          "message": "解析后的消息",
          "type": "http",
          "success": true,
          "timestamp": "2025-08-16 09:52:00",
          "method": "GET",
          "path": "/api/v1/accounts",
          "statusCode": 200,
          "duration": 0.265,
          "traceId": "315dcf3347c312a4f233"
        }
      ],
      "pagination": {
        "current_page": 1,
        "has_more": true,
        "next_page": 2,
        "page_size": 5,
        "previous_page": 0,
        "total_lines": 11003,
        "total_pages": 2201
      },
      "stats": {
        "total": 5,
        "http": 3,
        "mysql": 0,
        "redis": 2,
        "system": 0,
        "with_redis": 5,
        "with_sql": 3
      },
      "timestamp": "2025-08-16 09:52:00"
    }
  }
}
```

### 其他接口

- `GET /system/logs/latest` - 获取最新日志（向后兼容）
- `GET /system/logs/unified` - 获取统一日志（向后兼容）

## Web界面

### 访问地址
`http://localhost:9999/docs/logs.html`

### 界面功能
1. **控制面板**
   - 每页显示数量选择
   - 当前页码输入
   - 加载日志按钮
   - 回到第一页按钮

2. **统计信息**
   - 实时显示各类日志数量
   - 总日志数量

3. **分页导航**
   - 当前页/总页数显示
   - 上一页/下一页按钮
   - 自动禁用不可用的按钮

4. **日志展示**
   - 按类型分类显示
   - 支持展开查看详细信息
   - 自动滚动加载更多

## 使用方法

### 1. 通过Web界面
1. 打开浏览器访问 `http://localhost:9999/docs/logs.html`
2. 选择每页显示数量
3. 点击"加载日志"按钮
4. 向下滚动查看更多日志
5. 使用分页控制跳转到指定页面

### 2. 通过API接口
```bash
# 获取第一页，每页5条
curl "http://localhost:9999/system/logs/paginated?page=1&page_size=5"

# 获取第二页，每页10条
curl "http://localhost:9999/system/logs/paginated?page=2&page_size=10"
```

### 3. 无限滚动
- 在Web界面中，向下滚动到底部时会自动加载下一页
- 无需手动点击分页按钮
- 提供流畅的用户体验

## 技术实现

### 后端实现
- 使用Go语言实现
- 支持大文件高效读取
- 智能分页算法
- 实时统计计算

### 前端实现
- 纯HTML/CSS/JavaScript
- 响应式设计
- 自动滚动加载
- 错误处理和用户反馈

### 性能优化
- 避免一次性读取整个日志文件
- 按需加载指定范围的数据
- 智能缓存和缓冲机制
- 支持大文件（GB级别）

## 注意事项

1. **文件大小限制**
   - 建议日志文件不超过1GB
   - 超大文件可能影响加载性能

2. **内存使用**
   - 每页数据会占用一定内存
   - 建议合理设置页面大小

3. **网络请求**
   - 每次翻页都会发起新的API请求
   - 建议在网络良好的环境下使用

4. **浏览器兼容性**
   - 支持现代浏览器（Chrome、Firefox、Safari、Edge）
   - 不支持IE浏览器

## 故障排除

### 常见问题

1. **页面显示空白**
   - 检查服务是否正常运行
   - 检查API接口是否可访问
   - 查看浏览器控制台错误信息

2. **分页信息显示错误**
   - 检查日志文件是否存在
   - 检查文件权限是否正确
   - 查看服务日志获取详细错误信息

3. **加载速度慢**
   - 检查日志文件大小
   - 适当减少每页显示数量
   - 检查系统资源使用情况

### 调试方法

1. **查看API响应**
   ```bash
   curl -v "http://localhost:9999/system/logs/paginated?page=1&page_size=5"
   ```

2. **检查服务状态**
   ```bash
   ./status.sh
   ```

3. **查看服务日志**
   ```bash
   tail -f logs/marriage_system-access.log
   ```

## 更新日志

- **v1.0.0** - 初始版本，支持基本分页功能
- **v1.1.0** - 添加无限滚动和自动加载
- **v1.2.0** - 优化性能和错误处理
- **v1.3.0** - 添加统计信息和分类显示

## 联系支持

如果遇到问题或需要技术支持，请联系开发团队或查看项目文档。
