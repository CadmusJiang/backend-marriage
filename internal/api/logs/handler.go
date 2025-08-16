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

		// 解析页码
		if pageStr != "" {
			if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
				page = p
			}
		}

		// 解析页面大小
		if pageSizeStr != "" {
			if ps, err := strconv.Atoi(pageSizeStr); err == nil && ps > 0 && ps <= 100 {
				pageSize = ps
			}
		}

		// 调用服务获取分页日志
		response, err := h.logService.GetPaginatedLogs(configs.ProjectAccessLogFile, page, pageSize, logType)
		if err != nil {
			ctx.AbortWithError(core.Error(http.StatusInternalServerError, 500, "获取分页日志失败").WithError(err))
			return
		}

		// 转换为map格式以兼容现有前端
		var logs []map[string]interface{}
		for _, entry := range response.Logs {
			logs = append(logs, map[string]interface{}{
				"raw":        entry.Raw,
				"message":    entry.Message,
				"type":       entry.Type,
				"success":    entry.Success,
				"timestamp":  entry.Timestamp,
				"details":    entry.Details,
				"method":     entry.Method,
				"path":       entry.Path,
				"statusCode": entry.StatusCode,
				"duration":   entry.Duration,
				"traceId":    entry.TraceID,
			})
		}

		// 转换为map格式以兼容现有前端
		statsMap := map[string]interface{}{
			"total":      response.Stats.Total,
			"http":       response.Stats.HTTP,
			"mysql":      response.Stats.MySQL,
			"redis":      response.Stats.Redis,
			"system":     response.Stats.System,
			"with_sql":   response.Stats.WithSQL,
			"with_redis": response.Stats.WithRedis,
		}

		paginationMap := map[string]interface{}{
			"current_page":  response.Pagination.CurrentPage,
			"page_size":     response.Pagination.PageSize,
			"total_lines":   response.Pagination.TotalLines,
			"total_pages":   response.Pagination.TotalPages,
			"has_more":      response.Pagination.HasMore,
			"next_page":     response.Pagination.NextPage,
			"previous_page": response.Pagination.PreviousPage,
		}

		ctx.Payload(map[string]interface{}{
			"success": true,
			"data": map[string]interface{}{
				"logs":       logs,
				"stats":      statsMap,
				"pagination": paginationMap,
				"timestamp":  response.Timestamp,
			},
		})
	}
}
