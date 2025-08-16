package logs

import (
	"net/http"
	"strconv"
	"time"

	"github.com/xinliangnote/go-gin-api/configs"
	"github.com/xinliangnote/go-gin-api/internal/pkg/core"
	"github.com/xinliangnote/go-gin-api/internal/services/logs"
)

type Handler struct {
	logService *logs.LogService
}

func New() *Handler {
	return &Handler{
		logService: logs.New(),
	}
}

// GetLatestLogs 获取最近的日志
func (h *Handler) GetLatestLogs() core.HandlerFunc {
	return func(ctx core.Context) {
		const maxReadSize = int64(1 << 20) // 1MB

		lines, err := h.logService.GetLatestLogs(configs.ProjectAccessLogFile, maxReadSize, 10)
		if err != nil {
			ctx.AbortWithError(core.Error(http.StatusInternalServerError, 500, "获取日志失败").WithError(err))
			return
		}

		ctx.Payload(map[string]interface{}{
			"file":  configs.ProjectAccessLogFile,
			"count": len(lines),
			"lines": lines,
		})
	}
}

// GetUnifiedLogs 获取统一的日志数据
func (h *Handler) GetUnifiedLogs() core.HandlerFunc {
	return func(ctx core.Context) {
		const maxReadSize = int64(1 << 20) // 1MB

		logEntries, stats, err := h.logService.GetUnifiedLogs(configs.ProjectAccessLogFile, maxReadSize)
		if err != nil {
			ctx.AbortWithError(core.Error(http.StatusInternalServerError, 500, "获取统一日志失败").WithError(err))
			return
		}

		// 转换为map格式以兼容现有前端
		var logs []map[string]interface{}
		for _, entry := range logEntries {
			logs = append(logs, map[string]interface{}{
				"raw":       entry.Raw,
				"message":   entry.Message,
				"type":      entry.Type,
				"success":   entry.Success,
				"timestamp": entry.Timestamp,
				"details":   entry.Details,
			})
		}

		// 转换为map格式以兼容现有前端
		statsMap := map[string]interface{}{
			"total":      stats.Total,
			"http":       stats.HTTP,
			"mysql":      stats.MySQL,
			"redis":      stats.Redis,
			"system":     stats.System,
			"with_sql":   stats.WithSQL,
			"with_redis": stats.WithRedis,
		}

		ctx.Payload(map[string]interface{}{
			"success": true,
			"data": map[string]interface{}{
				"logs":      logs,
				"stats":     statsMap,
				"timestamp": time.Now().Format("2006-01-02 15:04:05"),
			},
		})
	}
}

// 新增：获取分页日志数据
func (h *Handler) GetPaginatedLogs() core.HandlerFunc {
	return func(ctx core.Context) {
		// 获取查询参数
		params := ctx.RequestInputParams()
		pageStr := params.Get("page")
		pageSizeStr := params.Get("page_size")
		logType := params.Get("type")

		// 设置默认值
		page := 1
		pageSize := 10
		if logType == "" {
			logType = "all"
		}

		// 解析分页参数
		if pageStr != "" {
			if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
				page = p
			}
		}
		if pageSizeStr != "" {
			if ps, err := strconv.Atoi(pageSizeStr); err == nil && ps > 0 && ps <= 100 {
				pageSize = ps
			}
		}

		const maxReadSize = int64(1 << 20) // 1MB

		response, err := h.logService.GetPaginatedLogs(configs.ProjectAccessLogFile, page, pageSize, logType)
		if err != nil {
			ctx.AbortWithError(core.Error(http.StatusInternalServerError, 500, "获取分页日志失败").WithError(err))
			return
		}

		ctx.Payload(response)
	}
}

// GetTraceLogs 获取指定Trace-ID的全链路日志
func (h *Handler) GetTraceLogs() core.HandlerFunc {
	return func(ctx core.Context) {
		// 获取查询参数
		params := ctx.RequestInputParams()
		traceID := params.Get("trace_id")

		if traceID == "" {
			ctx.AbortWithError(core.Error(http.StatusBadRequest, 400, "缺少trace_id参数"))
			return
		}

		const maxReadSize = int64(1 << 20) // 1MB

		logEntries, err := h.logService.GetTraceLogs(configs.ProjectAccessLogFile, traceID, maxReadSize)
		if err != nil {
			ctx.AbortWithError(core.Error(http.StatusInternalServerError, 500, "获取全链路日志失败").WithError(err))
			return
		}

		// 转换为map格式
		var logs []map[string]interface{}
		for _, entry := range logEntries {
			logs = append(logs, map[string]interface{}{
				"raw":         entry.Raw,
				"message":     entry.Message,
				"type":        entry.Type,
				"success":     entry.Success,
				"timestamp":   entry.Timestamp,
				"details":     entry.Details,
				"trace_id":    entry.TraceID,
				"method":      entry.Method,
				"path":        entry.Path,
				"status_code": entry.StatusCode,
				"duration":    entry.Duration,
			})
		}

		ctx.Payload(map[string]interface{}{
			"success": true,
			"data": map[string]interface{}{
				"trace_id":  traceID,
				"logs":      logs,
				"count":     len(logs),
				"timestamp": time.Now().Format("2006-01-02 15:04:05"),
			},
		})
	}
}

// GetTraceLogsByTimeRange 根据时间范围获取指定Trace-ID的日志
func (h *Handler) GetTraceLogsByTimeRange() core.HandlerFunc {
	return func(ctx core.Context) {
		// 获取查询参数
		params := ctx.RequestInputParams()
		traceID := params.Get("trace_id")
		startTimeStr := params.Get("start_time")
		endTimeStr := params.Get("end_time")

		if traceID == "" {
			ctx.AbortWithError(core.Error(http.StatusBadRequest, 400, "缺少trace_id参数"))
			return
		}

		// 解析时间参数
		startTime := time.Now().Add(-24 * time.Hour) // 默认查询最近24小时
		endTime := time.Now()

		if startTimeStr != "" {
			if t, err := time.Parse("2006-01-02 15:04:05", startTimeStr); err == nil {
				startTime = t
			}
		}
		if endTimeStr != "" {
			if t, err := time.Parse("2006-01-02 15:04:05", endTimeStr); err == nil {
				endTime = t
			}
		}

		const maxReadSize = int64(1 << 20) // 1MB

		logEntries, err := h.logService.GetTraceLogsByTimeRange(configs.ProjectAccessLogFile, traceID, startTime, endTime, maxReadSize)
		if err != nil {
			ctx.AbortWithError(core.Error(http.StatusInternalServerError, 500, "获取时间范围日志失败").WithError(err))
			return
		}

		// 转换为map格式
		var logs []map[string]interface{}
		for _, entry := range logEntries {
			logs = append(logs, map[string]interface{}{
				"raw":         entry.Raw,
				"message":     entry.Message,
				"type":        entry.Type,
				"success":     entry.Success,
				"timestamp":   entry.Timestamp,
				"details":     entry.Details,
				"trace_id":    entry.TraceID,
				"method":      entry.Method,
				"path":        entry.Path,
				"status_code": entry.StatusCode,
				"duration":    entry.Duration,
			})
		}

		ctx.Payload(map[string]interface{}{
			"success": true,
			"data": map[string]interface{}{
				"trace_id":   traceID,
				"start_time": startTime.Format("2006-01-02 15:04:05"),
				"end_time":   endTime.Format("2006-01-02 15:04:05"),
				"logs":       logs,
				"count":      len(logs),
				"timestamp":  time.Now().Format("2006-01-02 15:04:05"),
			},
		})
	}
}
