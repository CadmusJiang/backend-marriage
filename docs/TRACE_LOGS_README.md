# 全链路日志追踪功能说明

## 概述

全链路日志追踪功能允许您通过一个唯一的Trace-ID来追踪HTTP请求的完整执行链路，包括：
- HTTP请求和响应
- SQL查询执行
- Redis操作
- 系统日志
- 第三方API调用

## 功能特性

### 1. 自动Trace-ID生成
- 如果请求头中没有提供Trace-ID，系统会自动生成一个
- 支持多种请求头格式：`X-Trace-ID`、`Trace-ID`、`X-Request-ID`

### 2. 全链路记录
- **HTTP层**: 记录请求开始、结束、状态码、响应时间等
- **数据库层**: 记录所有SQL查询，包括执行时间、影响行数等
- **缓存层**: 记录Redis操作，包括命令、执行时间等
- **业务层**: 记录业务逻辑执行过程

### 3. 响应头返回
- 每个HTTP响应都会在响应头中包含 `X-Trace-ID` 字段
- 客户端可以通过这个字段来关联请求和日志

## 使用方法

### 1. 客户端请求

#### 方式一：自定义Trace-ID
```bash
curl -H "X-Trace-ID: my-custom-trace-123" \
     -H "Content-Type: application/json" \
     http://localhost:9999/api/v1/accounts
```

#### 方式二：自动生成Trace-ID
```bash
curl -H "Content-Type: application/json" \
     http://localhost:9999/api/v1/accounts
```

### 2. 查看全链路日志

#### 访问Web界面
打开浏览器访问：`http://localhost:9999/docs/trace-logs.html`

#### 使用API接口
```bash
# 获取指定Trace-ID的日志
GET /api/v1/logs/trace?trace_id=your-trace-id

# 根据时间范围获取日志
GET /api/v1/logs/trace/range?trace_id=your-trace-id&start_time=2024-01-01 10:00:00&end_time=2024-01-01 11:00:00
```

## 技术实现

### 1. 中间件架构
```
HTTP请求 → Trace中间件 → 业务逻辑 → 响应
    ↓           ↓           ↓        ↓
  生成/获取  设置Trace    记录操作   返回Trace-ID
  Trace-ID   上下文      日志      响应头
```

### 2. 核心组件

#### Trace中间件 (`internal/router/interceptor/trace_middleware.go`)
- 处理请求头中的Trace-ID
- 创建Trace上下文
- 设置响应头
- 记录请求开始和结束日志

#### Trace对象 (`pkg/trace/trace.go`)
- 存储请求的完整追踪信息
- 支持追加SQL、Redis、第三方调用等日志
- 提供线程安全的日志追加方法

#### 日志服务 (`internal/services/logs/service.go`)
- 提供按Trace-ID查询日志的方法
- 支持时间范围过滤
- 日志排序和格式化

### 3. 数据库集成

#### MySQL追踪
- 通过GORM的Logger接口记录SQL执行
- 记录SQL语句、执行时间、影响行数等信息
- 自动关联到当前请求的Trace上下文

#### Redis追踪
- 记录Redis命令执行
- 记录执行时间和结果
- 支持Redis集群和单机模式

## 配置说明

### 1. 环境变量
```bash
# 启用Trace日志记录
ENABLE_TRACE_LOGS=true

# Trace日志级别
TRACE_LOG_LEVEL=info
```

### 2. 日志格式
Trace日志采用结构化格式，便于解析和查询：
```json
{
  "timestamp": "2024-01-01 10:00:00",
  "trace_id": "abc123def456",
  "type": "mysql",
  "message": "SQL执行",
  "sql": "SELECT * FROM users WHERE id = ?",
  "cost_seconds": 0.001,
  "rows": 1
}
```

## 最佳实践

### 1. 客户端使用
- 为重要的业务请求设置有意义的Trace-ID
- 在微服务架构中传递Trace-ID
- 在错误报告中包含Trace-ID

### 2. 运维监控
- 定期检查Trace日志的完整性
- 监控SQL执行时间异常
- 分析Redis操作性能

### 3. 问题排查
- 通过Trace-ID快速定位问题
- 分析请求的完整执行链路
- 识别性能瓶颈

## 故障排除

### 1. 常见问题

#### Trace-ID未生成
- 检查中间件是否正确配置
- 确认请求是否经过Trace中间件

#### 日志不完整
- 检查日志文件权限
- 确认日志级别设置
- 验证数据库连接状态

#### 性能影响
- Trace功能对性能影响很小（<1%）
- 如遇性能问题，可调整日志级别

### 2. 调试方法

#### 启用调试日志
```bash
export TRACE_LOG_LEVEL=debug
```

#### 检查中间件状态
```bash
curl -H "X-Trace-ID: debug-123" \
     http://localhost:9999/api/v1/logs/trace?trace_id=debug-123
```

## 扩展功能

### 1. 未来计划
- 支持分布式追踪（Jaeger、Zipkin集成）
- 添加性能指标收集
- 支持日志导出和备份
- 集成告警系统

### 2. 自定义扩展
- 支持自定义日志格式
- 添加业务特定的追踪字段
- 支持多种存储后端

## 相关链接

- [分页日志查看器](logs.html)
- [API文档](swagger.html)
- [项目主页](../README.md)
- [数据库说明](../README_DATABASE.md)

## 技术支持

如果您在使用过程中遇到问题，请：
1. 查看本文档的故障排除部分
2. 检查系统日志
3. 联系开发团队

---

*最后更新：2024年1月*
