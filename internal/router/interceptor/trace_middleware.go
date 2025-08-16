package interceptor

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/xinliangnote/go-gin-api/pkg/trace"
	"go.uber.org/zap"
)

// TraceMiddleware 全链路追踪中间件
func TraceMiddleware(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求头获取Trace-ID，如果没有则自动生成
		traceID := c.GetHeader("X-Trace-ID")
		if traceID == "" {
			traceID = c.GetHeader("Trace-ID")
		}
		if traceID == "" {
			traceID = c.GetHeader("X-Request-ID")
		}

		// 创建新的Trace对象
		t := trace.New(traceID)

		// 设置请求信息
		t.WithRequest(&trace.Request{
			TTL:        c.GetHeader("X-Request-TTL"),
			Method:     c.Request.Method,
			DecodedURL: c.Request.URL.String(),
			Header:     c.Request.Header,
			Body:       c.Request.Body,
		})

		// 将Trace对象设置到gin.Context中，供后续中间件使用
		c.Set("_trace", t)

		// 设置Logger，确保所有日志都带有Trace-ID
		traceLogger := logger.With(
			zap.String("trace_id", t.ID()),
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.String("remote_addr", c.ClientIP()),
		)

		// 将Logger设置到gin.Context中
		c.Set("_trace_logger", traceLogger)

		// 在响应头中添加Trace-ID
		c.Header("X-Trace-ID", t.ID())

		// 记录请求开始日志
		traceLogger.Info("HTTP请求开始",
			zap.String("trace_id", t.ID()),
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.String("query", c.Request.URL.RawQuery),
			zap.String("user_agent", c.Request.UserAgent()),
		)

		// 继续处理请求
		c.Next()

		// 记录请求结束日志
		statusCode := c.Writer.Status()
		duration := float64(c.Writer.Size()) / 1024.0 // 响应大小KB

		traceLogger.Info("HTTP请求结束",
			zap.String("trace_id", t.ID()),
			zap.Int("status_code", statusCode),
			zap.Float64("response_size_kb", duration),
		)

		// 设置响应信息到Trace
		t.WithResponse(&trace.Response{
			Header:      c.Writer.Header(),
			HttpCode:    statusCode,
			HttpCodeMsg: http.StatusText(statusCode),
		})
	}
}

// GetTrace 从gin.Context获取Trace对象
func GetTrace(c *gin.Context) *trace.Trace {
	if t, exists := c.Get("_trace"); exists {
		return t.(*trace.Trace)
	}
	return nil
}

// GetTraceLogger 从gin.Context获取Trace Logger
func GetTraceLogger(c *gin.Context) *zap.Logger {
	if logger, exists := c.Get("_trace_logger"); exists {
		return logger.(*zap.Logger)
	}
	return nil
}
