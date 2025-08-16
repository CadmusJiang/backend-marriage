package logger

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"
)

// ExampleUsage 使用示例
func ExampleUsage() {
	// 1. 创建日志记录器
	logger, err := NewJSONLogger(
		WithDebugLevel(),
		WithField("service", "example"),
		WithTimeLayout("2006-01-02 15:04:05"),
	)
	if err != nil {
		panic(err)
	}

	// 2. 基础用法
	logger.Info("应用启动成功")
	logger.Debug("调试信息", zap.String("module", "auth"))
	logger.Warn("警告信息", zap.Int("retry_count", 3))

	// 模拟错误
	err = fmt.Errorf("模拟错误")
	logger.Error("错误信息", zap.Error(err))

	// 3. 结构化日志
	logger.Info("用户登录",
		zap.String("action", "login"),
		zap.String("username", "john_doe"),
		zap.Time("login_time", time.Now()),
	)

	// 4. 带上下文的日志
	ctx := context.WithValue(context.Background(), "trace_id", "trace-789")
	logger.Info("数据库查询",
		zap.String("query", "SELECT * FROM users"),
		zap.Duration("duration", 150*time.Millisecond),
		zap.String("trace_id", ctx.Value("trace_id").(string)),
	)

	// 5. 链式调用
	logger.With(
		zap.String("component", "database"),
		zap.String("operation", "insert"),
	).Info("插入用户记录",
		zap.String("table", "users"),
		zap.Int("affected_rows", 1),
	)

	// 6. 错误日志示例
	logger.Error("数据库连接失败",
		zap.String("database", "mysql"),
		zap.String("host", "localhost:3306"),
		zap.Error(err),
		zap.String("solution", "检查网络连接和数据库状态"),
	)

	// 7. 性能日志示例
	start := time.Now()
	// ... 执行某些操作 ...
	duration := time.Since(start)

	logger.Info("操作完成",
		zap.String("operation", "data_processing"),
		zap.Duration("duration", duration),
		zap.Int("records_processed", 1000),
	)

	// 8. 业务日志示例
	logger.Info("订单创建成功",
		zap.String("order_id", "ORD-2024-001"),
		zap.Float64("amount", 99.99),
		zap.String("currency", "USD"),
		zap.String("payment_method", "credit_card"),
		zap.String("status", "pending"),
	)

	// 9. 同步日志
	logger.Sync()
}

// ExampleWithContext 上下文使用示例
func ExampleWithContext() {
	logger, err := NewJSONLogger(
		WithInfoLevel(),
		WithField("service", "context-example"),
	)
	if err != nil {
		panic(err)
	}

	// 创建带跟踪ID的上下文
	ctx := context.WithValue(context.Background(), "trace_id", "trace-123")
	ctx = context.WithValue(ctx, "user_id", "user-456")

	// 使用上下文记录日志
	logger.Info("用户操作",
		zap.String("action", "profile_update"),
		zap.String("field", "email"),
		zap.String("trace_id", ctx.Value("trace_id").(string)),
		zap.String("user_id", ctx.Value("user_id").(string)),
	)

	// 传递上下文到其他函数
	processUserData(ctx, logger)
}

// ExampleWithFields 字段使用示例
func ExampleWithFields() {
	logger, err := NewJSONLogger(
		WithWarnLevel(),
		WithField("service", "fields-example"),
	)
	if err != nil {
		panic(err)
	}

	// 添加固定字段
	logger = logger.With(
		zap.String("environment", "production"),
		zap.String("version", "1.0.0"),
	)

	// 记录带字段的日志
	logger.Warn("系统警告",
		zap.String("component", "cache"),
		zap.String("issue", "memory_usage_high"),
		zap.Int("memory_mb", 1024),
	)
}

// ExampleWithFile 文件输出示例
func ExampleWithFile() {
	logger, err := NewJSONLogger(
		WithErrorLevel(),
		WithFileP("logs/error.log"),
		WithDisableConsole(),
	)
	if err != nil {
		panic(err)
	}

	// 错误日志将写入文件
	logger.Error("严重错误",
		zap.String("component", "database"),
		zap.String("error", "connection_timeout"),
		zap.Int("timeout_seconds", 30),
	)
}

// ExampleWithRotation 日志轮转示例
func ExampleWithRotation() {
	logger, err := NewJSONLogger(
		WithInfoLevel(),
		WithFileRotationP("logs/app.log"),
	)
	if err != nil {
		panic(err)
	}

	// 记录大量日志，测试轮转功能
	for i := 0; i < 1000; i++ {
		logger.Info("测试日志",
			zap.Int("index", i),
			zap.String("message", "这是第"+fmt.Sprintf("%d", i)+"条日志"),
		)
	}
}

// 辅助函数
func processUserData(ctx context.Context, logger *zap.Logger) {
	// 从上下文获取值
	traceID := ctx.Value("trace_id").(string)
	userID := ctx.Value("user_id").(string)

	logger.Info("处理用户数据",
		zap.String("trace_id", traceID),
		zap.String("user_id", userID),
		zap.String("operation", "data_processing"),
	)
}
