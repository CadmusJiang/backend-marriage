package logger

import (
	"go.uber.org/zap/zapcore"
)

// 日志级别常量
const (
	DebugLevel = zapcore.DebugLevel
	InfoLevel  = zapcore.InfoLevel
	WarnLevel  = zapcore.WarnLevel
	ErrorLevel = zapcore.ErrorLevel
)

// Config 日志配置
type Config struct {
	Level      zapcore.Level `json:"level" yaml:"level"`
	Format     string        `json:"format" yaml:"format"`           // json 或 console
	OutputPath string        `json:"output_path" yaml:"output_path"` // 输出路径
	ErrorPath  string        `json:"error_path" yaml:"error_path"`   // 错误输出路径
	MaxSize    int           `json:"max_size" yaml:"max_size"`       // 单个文件最大尺寸(MB)
	MaxAge     int           `json:"max_age" yaml:"max_age"`         // 最大保留天数
	MaxBackups int           `json:"max_backups" yaml:"max_backups"` // 最大备份数
	Compress   bool          `json:"compress" yaml:"compress"`       // 是否压缩
	Caller     bool          `json:"caller" yaml:"caller"`           // 是否显示调用者
	Stacktrace zapcore.Level `json:"stacktrace" yaml:"stacktrace"`   // 堆栈跟踪级别
}

// ConfigOption 配置选项函数（重命名避免与logger.go中的Option冲突）
type ConfigOption func(*Config)

// WithLevel 设置日志级别
func WithLevel(level zapcore.Level) ConfigOption {
	return func(c *Config) {
		c.Level = level
	}
}

// WithFormat 设置日志格式
func WithFormat(format string) ConfigOption {
	return func(c *Config) {
		c.Format = format
	}
}

// WithOutputPath 设置输出路径
func WithOutputPath(path string) ConfigOption {
	return func(c *Config) {
		c.OutputPath = path
	}
}

// WithErrorPath 设置错误输出路径
func WithErrorPath(path string) ConfigOption {
	return func(c *Config) {
		c.ErrorPath = path
	}
}

// WithMaxSize 设置单个文件最大尺寸
func WithMaxSize(size int) ConfigOption {
	return func(c *Config) {
		c.MaxSize = size
	}
}

// WithMaxAge 设置最大保留天数
func WithMaxAge(age int) ConfigOption {
	return func(c *Config) {
		c.MaxAge = age
	}
}

// WithMaxBackups 设置最大备份数
func WithMaxBackups(backups int) ConfigOption {
	return func(c *Config) {
		c.MaxBackups = backups
	}
}

// WithCompress 设置是否压缩
func WithCompress(compress bool) ConfigOption {
	return func(c *Config) {
		c.Compress = compress
	}
}

// WithCaller 设置是否显示调用者
func WithCaller(caller bool) ConfigOption {
	return func(c *Config) {
		c.Caller = caller
	}
}

// WithStacktrace 设置堆栈跟踪级别
func WithStacktrace(level zapcore.Level) ConfigOption {
	return func(c *Config) {
		c.Stacktrace = level
	}
}

// 便捷配置选项
var (
	// 开发环境配置
	Development = func() ConfigOption {
		return func(c *Config) {
			c.Level = DebugLevel
			c.Format = "console"
			c.Caller = true
			c.Stacktrace = WarnLevel
		}
	}

	// 生产环境配置
	Production = func() ConfigOption {
		return func(c *Config) {
			c.Level = InfoLevel
			c.Format = "json"
			c.Caller = false
			c.Stacktrace = ErrorLevel
		}
	}

	// 测试环境配置
	Testing = func() ConfigOption {
		return func(c *Config) {
			c.Level = DebugLevel
			c.Format = "console"
			c.Caller = true
			c.Stacktrace = ErrorLevel
		}
	}
)
