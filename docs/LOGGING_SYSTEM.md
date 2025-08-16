# 日志系统使用说明

## 概述

本系统实现了一个符合行业标准的全链路日志追踪系统，支持结构化日志记录，包括meta信息和日志内容。所有日志都会自动关联到请求的Trace-ID，实现全链路追踪。

## 核心特性

- **结构化日志**: 支持JSON格式的结构化日志输出
- **全链路追踪**: 通过Trace-ID关联同一请求的所有日志
- **元数据支持**: 支持业务元数据，如用户ID、操作类型等
- **上下文感知**: 自动从请求头获取Trace-ID
- **性能优化**: 基于zap的高性能日志库

## Trace-ID获取策略

系统按以下优先级获取Trace-ID：

1. **X-Trace-ID** (推荐)
2. **Trace-ID** 
3. **X-Request-ID**

## 使用方法

### 1. 在Handler层使用

```go
func (h *handler) ListTeams() core.HandlerFunc {
    return func(c core.Context) {
        // 获取trace_id
        traceID := c.GetTraceID()
        
        // 记录日志
        h.logger.Info("获取团队列表",
            zap.String("trace_id", traceID),
            zap.String("operation", "ListTeams"),
        )
        
        // ... 业务逻辑
    }
}
```

### 2. 在Service层使用

```go
func (s *service) ListTeams(ctx Context, belongGroupId uint32, current, pageSize int, keyword string, scope *authz.AccessScope) ([]map[string]interface{}, int64, error) {
    // 获取trace_id用于日志关联
    traceID := ctx.GetTraceID()
    
    // 创建日志元数据
    meta := &logger.BaseMeta{
        TraceID: traceID,
        Service: "organization-service",
        Extra: map[string]interface{}{
            "operation":     "ListTeams",
            "belongGroupId": belongGroupId,
            "current":       current,
            "pageSize":      pageSize,
            "keyword":       keyword,
        },
    }

    // 记录调试日志
    logger.GetGlobalLogger().WithMeta(meta).Debug("开始获取团队列表",
        zap.Uint32("belongGroupId", belongGroupId),
        zap.Any("scope", scope),
    )
    
    // ... 业务逻辑
    
    // 记录成功日志
    logger.GetGlobalLogger().WithMeta(meta).Info("团队列表获取成功",
        zap.Int("total", len(out)),
        zap.Int64("totalCount", total),
    )
    
    return out, total, nil
}
```

### 3. 日志级别

- **Debug**: 调试信息，开发环境使用
- **Info**: 一般信息，记录重要操作
- **Warn**: 警告信息，需要注意但不影响功能
- **Error**: 错误信息，操作失败
- **Fatal**: 致命错误，程序无法继续运行

### 4. 元数据结构

```go
type BaseMeta struct {
    TraceID       string                 `json:"trace_id,omitempty"`
    RequestID     string                 `json:"request_id,omitempty"`
    UserID        string                 `json:"user_id,omitempty"`
    SessionID     string                 `json:"session_id,omitempty"`
    CorrelationID string                 `json:"correlation_id,omitempty"`
    Service       string                 `json:"service,omitempty"`
    Version       string                 `json:"version,omitempty"`
    Environment   string                 `json:"environment,omitempty"`
    Host          string                 `json:"host,omitempty"`
    IP            string                 `json:"ip,omitempty"`
    UserAgent     string                 `json:"user_agent,omitempty"`
    Extra         map[string]interface{} `json:"extra,omitempty"`
}
```

## 配置选项

### 初始化日志系统

```go
import "github.com/xinliangnote/go-gin-api/pkg/logger"

func main() {
    // 初始化日志系统
    logger.Init(
        logger.WithLevel(logger.InfoLevel),        // 日志级别
        logger.WithFormat("json"),                 // 输出格式
        logger.WithOutputPath("logs/app.log"),     // 输出文件
        logger.WithErrorPath("logs/error.log"),    // 错误日志文件
        logger.WithCaller(true),                   // 显示调用者
        logger.WithStacktrace(logger.ErrorLevel),  // 错误时显示堆栈
    )
}
```

### 便捷配置

```go
// 开发环境配置
logger.Init(logger.Development())

// 生产环境配置  
logger.Init(logger.Production())

// 测试环境配置
logger.Init(logger.Testing())
```

## 最佳实践

### 1. 日志记录原则

- **结构化**: 使用结构化字段而不是字符串拼接
- **有意义**: 记录有业务价值的操作和状态
- **适度**: 避免过度记录，影响性能
- **一致性**: 保持日志格式和字段的一致性

### 2. 错误日志

```go
logger.GetGlobalLogger().WithMeta(meta).Error("操作失败",
    zap.Error(err),                    // 错误对象
    zap.String("operation", "create"), // 操作类型
    zap.String("resource", "user"),    // 资源类型
    zap.String("solution", "检查参数"), // 解决建议
)
```

### 3. 性能日志

```go
start := time.Now()
// ... 执行操作
duration := time.Since(start)

logger.GetGlobalLogger().WithMeta(meta).Info("操作完成",
    zap.String("operation", "database_query"),
    zap.Duration("duration", duration),
    zap.Int("records_processed", count),
)
```

### 4. 业务日志

```go
logger.GetGlobalLogger().WithMeta(meta).Info("用户操作",
    zap.String("action", "login"),
    zap.String("username", username),
    zap.String("ip", clientIP),
    zap.String("result", "success"),
)
```

## 测试

使用提供的测试脚本验证日志系统：

```bash
# 测试Teams API的日志系统
./test_teams_logging.sh

# 测试日志解析
./test_log_parsing.sh
```

## 监控和查询

### 1. 通过Web界面查看

访问 `http://localhost:9999/docs/trace-logs.html` 查看全链路日志。

### 2. 通过API查询

```bash
# 查询指定Trace-ID的日志
curl "http://localhost:9999/api/v1/logs/trace?trace_id=YOUR_TRACE_ID"

# 按时间范围查询
curl "http://localhost:9999/api/v1/logs/trace/range?trace_id=YOUR_TRACE_ID&start_time=2024-01-01 00:00:00&end_time=2024-01-01 23:59:59"
```

### 3. 日志文件分析

```bash
# 查找特定Trace-ID的日志
grep "YOUR_TRACE_ID" logs/app.log

# 统计日志数量
grep -c "YOUR_TRACE_ID" logs/app.log

# 查看最近的日志
tail -f logs/app.log | grep "YOUR_TRACE_ID"
```

## 故障排除

### 1. 日志不显示

- 检查日志级别配置
- 确认日志文件路径正确
- 检查文件权限

### 2. Trace-ID不关联

- 确认请求头包含正确的Trace-ID
- 检查中间件是否正确设置
- 验证Context传递是否正确

### 3. 性能问题

- 调整日志级别
- 使用异步日志
- 检查磁盘I/O性能

## 扩展

### 1. 自定义元数据

```go
type CustomMeta struct {
    logger.BaseMeta
    BusinessID   string `json:"business_id"`
    OperationType string `json:"operation_type"`
}

func (m CustomMeta) ToFields() []Field {
    fields := m.BaseMeta.ToFields()
    fields = append(fields,
        zap.String("business_id", m.BusinessID),
        zap.String("operation_type", m.OperationType),
    )
    return fields
}
```

### 2. 自定义日志格式

```go
logger.Init(
    logger.WithFormat("console"),           // 控制台格式
    logger.WithTimeLayout("2006-01-02 15:04:05"), // 自定义时间格式
)
```

## 总结

新的日志系统提供了：

1. **全链路追踪**: 通过Trace-ID关联所有相关日志
2. **结构化输出**: JSON格式，便于解析和分析
3. **灵活配置**: 支持多种配置选项
4. **高性能**: 基于zap的高性能实现
5. **易于使用**: 简洁的API设计

通过合理使用这个日志系统，可以大大提高系统的可观测性和问题排查效率。
