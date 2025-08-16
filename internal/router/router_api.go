package router

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/xinliangnote/go-gin-api/configs"
	"github.com/xinliangnote/go-gin-api/internal/api/account"
	"github.com/xinliangnote/go-gin-api/internal/api/analytics"
	"github.com/xinliangnote/go-gin-api/internal/api/cooperation_store"
	"github.com/xinliangnote/go-gin-api/internal/api/customer_authorization_record"
	"github.com/xinliangnote/go-gin-api/internal/api/customer_authorization_record_history"
	"github.com/xinliangnote/go-gin-api/internal/api/organization"
	"github.com/xinliangnote/go-gin-api/internal/code"
	"github.com/xinliangnote/go-gin-api/internal/pkg/core"
)

func setApiRouter(r *resource) {
	// account
	accountHandler := account.New(r.logger, r.db, r.cache)

	// cooperation_store
	cooperationStoreHandler := cooperation_store.New(r.logger, r.db, r.cache)

	// customer_authorization_record
	customerAuthRecordHandler := customer_authorization_record.New(r.logger, r.db, r.db)

	// customer_authorization_record_history
	customerAuthRecordHistoryHandler := customer_authorization_record_history.New(r.logger, r.db, r.cache)

	// helper 已移除

	// organization
	organizationHandler := organization.New(r.logger, r.db, r.cache)

	// analytics
	analyticsHandler := analytics.New(r.logger, r.db, r.db)

	// API v1 路由组 - 无需登录验证（CoreAuth）
	apiV1 := r.mux.Group("/api/v1")
	{
		// CoreAuth - 这些接口需要trace记录但不需要登录验证
		authAPI := apiV1.Group("/auth")
		{
			authAPI.POST("/login", accountHandler.Login())
			authAPI.POST("/logout", accountHandler.Logout())
		}

		// helper 接口已删除
	}

	// 需要签名验证、登录验证、RBAC 权限验证
	api := r.mux.Group("/api", core.WrapAuthHandler(r.interceptors.CheckLogin), r.interceptors.CheckRBAC())
	{
		// Organization
		api.GET("/v1/groups", organizationHandler.GetGroups())
		api.POST("/v1/groups", organizationHandler.CreateGroup())
		api.GET("/v1/groups/:orgId", organizationHandler.GetOrgInfoDetail())
		api.PUT("/v1/groups/:orgId", organizationHandler.UpdateGroup())
		api.GET("/v1/groups/:orgId/history", organizationHandler.GetGroupHistory())

		// teams endpoints
		api.GET("/v1/teams", organizationHandler.ListTeams())
		api.POST("/v1/teams", organizationHandler.CreateTeam())
		api.GET("/v1/teams/:teamId", organizationHandler.GetTeam())
		api.PUT("/v1/teams/:teamId", organizationHandler.UpdateTeam())
		api.GET("/v1/teams/:teamId/history", organizationHandler.GetTeamHistory())
		api.GET("/v1/teams/:teamId/members", organizationHandler.ListTeamMembers())
		api.POST("/v1/teams/:teamId/members", organizationHandler.AddTeamMember())
		api.DELETE("/v1/teams/:teamId/members/:memberId", organizationHandler.RemoveTeamMember())
		api.PUT("/v1/teams/:teamId/members/:memberId/role", organizationHandler.UpdateTeamMemberRole())
		api.GET("/v1/unassigned-account", organizationHandler.ListUnassignedAccounts())

		// Analytics
		api.GET("/v1/analytics/account/rankings", analyticsHandler.GetAccountRankings())
		api.GET("/v1/analytics/account/trends", analyticsHandler.GetAccountTrends())
		api.GET("/v1/analytics/teams/rankings", analyticsHandler.GetTeamsRankings())
		api.GET("/v1/analytics/teams/trends", analyticsHandler.GetTeamsTrends())
		api.GET("/v1/analytics/customer-authorization-record/composition", analyticsHandler.GetCustomerComposition())

		// Customer
		api.GET("/v1/customer-authorization-records", customerAuthRecordHandler.GetCustomerAuthorizationRecordList())
		api.POST("/v1/customer-authorization-records", customerAuthRecordHandler.CreateCustomerAuthorizationRecord())
		api.GET("/v1/customer-authorization-records/check-phone", customerAuthRecordHandler.CheckPhoneExistence())
		api.GET("/v1/customer-authorization-records/:id", customerAuthRecordHandler.GetCustomerAuthorizationRecordDetail())
		api.PUT("/v1/customer-authorization-records/:id", customerAuthRecordHandler.UpdateCustomerAuthorizationRecord())
		// Customer History
		api.GET("/v1/customer-authorization-record-histories", customerAuthRecordHistoryHandler.GetCustomerAuthorizationRecordHistoryList())
		api.GET("/v1/customer-authorization-record-histories/:historyId", customerAuthRecordHistoryHandler.GetCustomerAuthorizationRecordHistoryDetail())

		// CooperationStore
		api.GET("/v1/cooperation-stores", cooperationStoreHandler.GetCooperationStoreList())
		api.POST("/v1/cooperation-stores", cooperationStoreHandler.CreateCooperationStore())
		api.GET("/v1/cooperation-stores/:id", cooperationStoreHandler.GetCooperationStoreDetail())
		api.PUT("/v1/cooperation-stores/:id", cooperationStoreHandler.UpdateCooperationStore())
		api.GET("/v1/cooperation-stores/:id/history", cooperationStoreHandler.GetCooperationStoreHistory())

		// Account
		api.GET("/v1/accounts", accountHandler.GetAccountList())
		api.POST("/v1/accounts", accountHandler.CreateAccount())
		api.GET("/v1/accounts/:accountId", accountHandler.GetAccountDetail())
		api.PUT("/v1/accounts/:accountId", accountHandler.UpdateAccount())
		api.PUT("/v1/accounts/:accountId/password", accountHandler.UpdateAccountPassword())
		api.GET("/v1/account-histories", accountHandler.GetAccountHistories())

		// Me (Account)
		api.GET("/v1/me", accountHandler.Me())

		// org-infos 路由移除，统一使用 /v1/groups 系列
	}

	// system 工具接口：最近日志
	system := r.mux.Group("/system")
	{
		// 统一日志监控页面
		system.GET("/logs", func(ctx core.Context) {
			ctx.Payload(map[string]interface{}{
				"message":     "统一日志监控页面",
				"description": "请访问 /docs/logs.html 查看完整的日志监控面板",
				"features": []string{
					"HTTP请求日志",
					"MySQL访问日志",
					"Redis操作日志",
					"系统性能监控",
				},
				"urls": map[string]string{
					"monitor_page": "/docs/logs.html",
					"http_logs":    "/system/logs/latest",
					"mysql_logs":   "/system/mysql-logs",
					"redis_logs":   "/system/redis/status",
				},
			})
		})

		// 统一日志API接口
		system.GET("/logs/unified", func(ctx core.Context) {
			const maxReadSize = int64(1 << 20) // 1MB

			// 统一日志条目结构
			type UnifiedLogEntry struct {
				ID         string                   `json:"id"`         // 日志ID
				Timestamp  string                   `json:"timestamp"`  // 时间戳
				Type       string                   `json:"type"`       // 日志类型：http, mysql, redis, system
				Level      string                   `json:"level"`      // 日志级别
				Message    string                   `json:"message"`    // 日志消息
				TraceID    string                   `json:"traceId"`    // 追踪ID
				Method     string                   `json:"method"`     // HTTP方法
				Path       string                   `json:"path"`       // 请求路径
				StatusCode int                      `json:"statusCode"` // HTTP状态码
				Duration   float64                  `json:"duration"`   // 请求耗时（秒）
				Success    bool                     `json:"success"`    // 是否成功
				Error      string                   `json:"error"`      // 错误信息
				SQLs       []map[string]interface{} `json:"sqls"`       // SQL操作列表
				SQLCount   int                      `json:"sqlCount"`   // SQL操作数量
				RedisOps   []map[string]interface{} `json:"redisOps"`   // Redis操作列表
				RedisCount int                      `json:"redisCount"` // Redis操作数量
				Details    map[string]interface{}   `json:"details"`    // 详细信息
			}

			readLastNLines := func(filePath string, maxBytes int64, n int) ([]string, error) {
				f, err := os.Open(filePath)
				if err != nil {
					return nil, err
				}
				defer f.Close()

				st, err := f.Stat()
				if err != nil {
					return nil, err
				}

				start := int64(0)
				if st.Size() > maxBytes {
					start = st.Size() - maxBytes
				}

				if _, err := f.Seek(start, io.SeekStart); err != nil {
					return nil, err
				}

				data, err := io.ReadAll(f)
				if err != nil {
					return nil, err
				}

				s := string(data)
				if start > 0 {
					if idx := strings.IndexByte(s, '\n'); idx >= 0 {
						s = s[idx+1:]
					}
				}

				lines := strings.Split(s, "\n")
				if len(lines) > 0 && lines[len(lines)-1] == "" {
					lines = lines[:len(lines)-1]
				}
				if len(lines) > n {
					lines = lines[len(lines)-n:]
				}
				return lines, nil
			}

			// 解析统一日志行
			parseUnifiedLogLine := func(line string, index int) *UnifiedLogEntry {
				var logData map[string]interface{}
				if err := json.Unmarshal([]byte(line), &logData); err != nil {
					return nil
				}

				entry := &UnifiedLogEntry{
					ID: fmt.Sprintf("log_%d", index),
				}

				// 基础信息
				if timeStr, ok := logData["time"].(string); ok {
					entry.Timestamp = timeStr
				}
				if level, ok := logData["level"].(string); ok {
					entry.Level = level
				}
				if msg, ok := logData["msg"].(string); ok {
					entry.Message = msg
				}

				// 检查是否是trace-log类型的日志
				if entry.Message == "trace-log" {
					entry.Type = "http"

					// 解析trace_info
					if traceInfo, ok := logData["trace_info"].(map[string]interface{}); ok {
						if traceID, ok := traceInfo["trace_id"].(string); ok {
							entry.TraceID = traceID
						}

						// HTTP请求信息
						if request, ok := traceInfo["request"].(map[string]interface{}); ok {
							if method, ok := request["method"].(string); ok {
								entry.Method = method
							}
							if path, ok := request["decoded_url"].(string); ok {
								entry.Path = path
							}
						}

						if response, ok := traceInfo["response"].(map[string]interface{}); ok {
							if httpCode, ok := response["http_code"].(float64); ok {
								entry.StatusCode = int(httpCode)
							}
							if costSeconds, ok := response["cost_seconds"].(float64); ok {
								entry.Duration = costSeconds
							}
						}

						if success, ok := traceInfo["success"].(bool); ok {
							entry.Success = success
						}

						// SQL信息
						if sqls, ok := traceInfo["sqls"].([]interface{}); ok && len(sqls) > 0 {
							entry.SQLCount = len(sqls)
							sqlDetails := make([]map[string]interface{}, 0, len(sqls))
							for _, sqlInterface := range sqls {
								if sqlMap, ok := sqlInterface.(map[string]interface{}); ok {
									sqlDetails = append(sqlDetails, sqlMap)
								}
							}
							entry.SQLs = sqlDetails
							entry.Type = "http_with_sql"
						}

						// Redis信息
						if redis, ok := traceInfo["redis"].([]interface{}); ok && len(redis) > 0 {
							entry.RedisCount = len(redis)
							redisDetails := make([]map[string]interface{}, 0, len(redis))
							for _, redisInterface := range redis {
								if redisMap, ok := redisInterface.(map[string]interface{}); ok {
									redisDetails = append(redisDetails, redisMap)
								}
							}
							entry.RedisOps = redisDetails
							if entry.Type == "http_with_sql" {
								entry.Type = "http_with_sql_redis"
							} else {
								entry.Type = "http_with_redis"
							}
						}
					}
				} else {
					// 其他类型的日志
					if strings.Contains(entry.Message, "mysql") || strings.Contains(entry.Message, "sql") {
						entry.Type = "mysql"
					} else if strings.Contains(entry.Message, "redis") {
						entry.Type = "redis"
					} else {
						entry.Type = "system"
					}

					// 为系统日志设置success字段
					if entry.Type == "system" {
						// 根据日志级别和消息内容判断是否成功
						if entry.Level == "fatal" || entry.Level == "error" {
							entry.Success = false
						} else if strings.Contains(strings.ToLower(entry.Message), "成功") ||
							strings.Contains(strings.ToLower(entry.Message), "success") ||
							strings.Contains(strings.ToLower(entry.Message), "完成") {
							entry.Success = true
						} else if strings.Contains(strings.ToLower(entry.Message), "失败") ||
							strings.Contains(strings.ToLower(entry.Message), "error") ||
							strings.Contains(strings.ToLower(entry.Message), "fatal") {
							entry.Success = false
						} else {
							// 默认情况下，info和debug级别的日志认为是成功的
							entry.Success = entry.Level == "info" || entry.Level == "debug"
						}
					}
				}

				// 错误信息
				if errMsg, ok := logData["error"].(string); ok {
					entry.Error = errMsg
				}

				// 详细信息
				details := make(map[string]interface{})
				details["原始数据"] = logData
				entry.Details = details

				return entry
			}

			lines, err := readLastNLines(configs.ProjectAccessLogFile, maxReadSize, 50)
			if err != nil {
				ctx.AbortWithError(core.Error(http.StatusInternalServerError, code.ServerError, code.Text(code.ServerError)).WithError(err))
				return
			}

			// 解析日志条目
			var logEntries []UnifiedLogEntry
			for i, line := range lines {
				if strings.TrimSpace(line) != "" {
					entry := parseUnifiedLogLine(line, i+1)
					if entry != nil && entry.Path != "/system/logs/unified" {
						logEntries = append(logEntries, *entry)
					}
				}
			}

			// 按时间倒序排列（最新的在前）
			for i, j := 0, len(logEntries)-1; i < j; i, j = i+1, j-1 {
				logEntries[i], logEntries[j] = logEntries[j], logEntries[i]
			}

			// 统计信息
			stats := map[string]int{
				"total":      0,
				"http":       0,
				"mysql":      0,
				"redis":      0,
				"system":     0,
				"with_sql":   0,
				"with_redis": 0,
			}

			for _, entry := range logEntries {
				stats["total"]++
				switch entry.Type {
				case "http", "http_with_sql", "http_with_redis", "http_with_sql_redis":
					stats["http"]++
					if entry.SQLCount > 0 {
						stats["with_sql"]++
					}
					if entry.RedisCount > 0 {
						stats["with_redis"]++
					}
				case "mysql":
					stats["mysql"]++
				case "redis":
					stats["redis"]++
				case "system":
					stats["system"]++
				}
			}

			ctx.Payload(map[string]interface{}{
				"file":      configs.ProjectAccessLogFile,
				"total":     len(logEntries),
				"timestamp": time.Now().Format("2006-01-02 15:04:05"),
				"stats":     stats,
				"logs":      logEntries,
			})
		})
		system.GET("/logs/latest", func(ctx core.Context) {
			const maxReadSize = int64(1 << 20) // 1MB

			readLastNLines := func(filePath string, maxBytes int64, n int) ([]string, error) {
				f, err := os.Open(filePath)
				if err != nil {
					return nil, err
				}
				defer f.Close()

				st, err := f.Stat()
				if err != nil {
					return nil, err
				}

				start := int64(0)
				if st.Size() > maxBytes {
					start = st.Size() - maxBytes
				}

				if _, err := f.Seek(start, io.SeekStart); err != nil {
					return nil, err
				}

				data, err := io.ReadAll(f)
				if err != nil {
					return nil, err
				}

				s := string(data)
				if start > 0 {
					if idx := strings.IndexByte(s, '\n'); idx >= 0 {
						s = s[idx+1:]
					}
				}

				lines := strings.Split(s, "\n")
				if len(lines) > 0 && lines[len(lines)-1] == "" {
					lines = lines[:len(lines)-1]
				}
				if len(lines) > n {
					lines = lines[len(lines)-n:]
				}
				return lines, nil
			}

			lines, err := readLastNLines(configs.ProjectAccessLogFile, maxReadSize, 10)
			if err != nil {
				ctx.AbortWithError(core.Error(http.StatusInternalServerError, code.ServerError, code.Text(code.ServerError)).WithError(err))
				return
			}

			ctx.Payload(map[string]interface{}{
				"file":  configs.ProjectAccessLogFile,
				"count": len(lines),
				"lines": lines,
			})
		})
	}
}
