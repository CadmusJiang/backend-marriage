package helper

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/xinliangnote/go-gin-api/configs"
	"github.com/xinliangnote/go-gin-api/internal/code"
	"github.com/xinliangnote/go-gin-api/internal/pkg/core"
	"github.com/xinliangnote/go-gin-api/pkg/file"
	"go.uber.org/zap"
)

type logData struct {
	Level       string      `json:"level"`
	Time        string      `json:"time"`
	Path        string      `json:"path"`
	HTTPCode    int         `json:"http_code"`
	Method      string      `json:"method"`
	Msg         string      `json:"msg"`
	TraceID     string      `json:"trace_id"`
	Content     string      `json:"content"`
	CostSeconds float64     `json:"cost_seconds"`
	Success     bool        `json:"success"`
	Details     interface{} `json:"details"`
}

type logsResponse struct {
	Data     []logData `json:"data"`
	Total    int       `json:"total"`
	Success  bool      `json:"success"`
	LastTime string    `json:"last_time"`
}

type logParseData struct {
	Level        string  `json:"level"`
	Time         string  `json:"time"`
	Caller       string  `json:"caller"`
	Msg          string  `json:"msg"`
	Domain       string  `json:"domain"`
	Method       string  `json:"method"`
	Path         string  `json:"path"`
	HTTPCode     int     `json:"http_code"`
	BusinessCode int     `json:"business_code"`
	Success      bool    `json:"success"`
	CostSeconds  float64 `json:"cost_seconds"`
	TraceID      string  `json:"trace_id"`
}

// GetLogs 获取日志列表
// @Summary 获取日志列表
// @Description 获取最新的日志记录
// @Tags API.helper
// @Accept application/x-www-form-urlencoded
// @Produce json
// @Param limit query int false "获取条数" default(50)
// @Success 200 {object} logsResponse
// @Failure 400 {object} code.Failure
// @Router /api/v1/logs [get]
func (h *handler) GetLogs() core.HandlerFunc {
	return func(c core.Context) {
		readLineFromEnd, err := file.NewReadLineFromEnd(configs.ProjectAccessLogFile)
		if err != nil {
			h.logger.Error("NewReadLineFromEnd err", zap.Error(err))
			c.AbortWithError(core.Error(
				http.StatusInternalServerError,
				code.ServerError,
				code.Text(code.ServerError)).WithError(err),
			)
			return
		}

		var logs []logData
		var lastTime string

		// 读取指定数量的日志
		for i := 0; i < 100; i++ { // 最多读取100条
			content, _ := readLineFromEnd.ReadLine()
			if string(content) != "" {
				var logParse logParseData
				err = json.Unmarshal(content, &logParse)
				if err != nil {
					h.logger.Error("json Unmarshal err", zap.Error(err))
					continue
				}

				// 只处理trace-log类型的日志
				if logParse.Msg == "trace-log" {
					h.logger.Info("处理trace-log", zap.String("path", logParse.Path))
					// 过滤掉日志相关的请求
					isLog := isLogRelatedRequest(logParse.Path)
					h.logger.Info("检查是否为日志请求", zap.String("path", logParse.Path), zap.Bool("isLog", isLog))
					if isLog {
						h.logger.Info("过滤掉日志相关请求", zap.String("path", logParse.Path))
						continue
					}

					// 简化日志内容，只保留关键信息
					simplifiedContent := map[string]interface{}{
						"method":       logParse.Method,
						"path":         logParse.Path,
						"http_code":    logParse.HTTPCode,
						"success":      logParse.Success,
						"cost_seconds": logParse.CostSeconds,
						"trace_id":     logParse.TraceID,
					}

					contentBytes, _ := json.Marshal(simplifiedContent)

					data := logData{
						Content:     string(contentBytes),
						Level:       logParse.Level,
						Time:        logParse.Time,
						Path:        logParse.Path,
						Method:      logParse.Method,
						Msg:         logParse.Msg,
						HTTPCode:    logParse.HTTPCode,
						TraceID:     logParse.TraceID,
						CostSeconds: logParse.CostSeconds,
						Success:     logParse.Success,
					}
					logs = append(logs, data)

					if lastTime == "" {
						lastTime = logParse.Time
					}
				}
			}
		}

		// 限制返回数量
		if len(logs) > 50 {
			logs = logs[:50]
		}

		res := &logsResponse{
			Data:     logs,
			Total:    len(logs),
			Success:  true,
			LastTime: lastTime,
		}

		c.Payload(res)
	}
}

// GetLogsRealtime 获取实时日志
// @Summary 获取实时日志
// @Description 获取指定时间之后的日志记录
// @Tags API.helper
// @Accept application/x-www-form-urlencoded
// @Produce json
// @Param since query string false "起始时间" default("")
// @Success 200 {object} logsResponse
// @Failure 400 {object} code.Failure
// @Router /api/v1/logs/realtime [get]
func (h *handler) GetLogsRealtime() core.HandlerFunc {
	return func(c core.Context) {
		since := c.GetHeader("since")

		readLineFromEnd, err := file.NewReadLineFromEnd(configs.ProjectAccessLogFile)
		if err != nil {
			h.logger.Error("NewReadLineFromEnd err", zap.Error(err))
			c.AbortWithError(core.Error(
				http.StatusInternalServerError,
				code.ServerError,
				code.Text(code.ServerError)).WithError(err),
			)
			return
		}

		var logs []logData
		var lastTime string
		maxLogs := 10 // 只返回最新的10条记录

		// 读取日志文件，收集所有符合条件的日志
		var allLogs []logData
		for {
			content, err := readLineFromEnd.ReadLine()
			if err != nil {
				break
			}

			if len(content) == 0 {
				continue
			}

			var fullLog map[string]interface{}
			if err := json.Unmarshal(content, &fullLog); err != nil {
				continue
			}

			// 检查时间过滤
			if since != "" {
				if logTime, ok := fullLog["time"].(string); ok {
					if logTime <= since {
						break
					}
				}
			}

			// 提取基本信息
			log := logData{
				Time:        getString(fullLog, "time"),
				Method:      getString(fullLog, "method"),
				Path:        getString(fullLog, "path"),
				HTTPCode:    getInt(fullLog, "http_code"),
				Success:     getBool(fullLog, "success"),
				CostSeconds: getFloat(fullLog, "cost_seconds"),
				TraceID:     getString(fullLog, "trace_id"),
			}

			// 只添加有效的日志记录（有路径和方法的）
			if log.Path == "" || log.Method == "" {
				continue
			}

			// 过滤掉日志相关的请求
			isLog := isLogRelatedRequest(log.Path)
			h.logger.Info("实时日志检查是否为日志请求", zap.String("path", log.Path), zap.Bool("isLog", isLog))
			if isLog {
				h.logger.Info("实时日志过滤掉日志相关请求", zap.String("path", log.Path))
				continue
			}

			// 只处理trace-log类型的日志
			if getString(fullLog, "msg") != "trace-log" {
				continue
			}

			// 只提取关键的请求和响应信息，不包含完整的trace_info
			if traceInfo, ok := fullLog["trace_info"].(map[string]interface{}); ok {
				details := make(map[string]interface{})

				// 简化的请求信息
				if req, ok := traceInfo["request"].(map[string]interface{}); ok {
					requestDetails := make(map[string]interface{})
					requestDetails["method"] = getString(req, "method")
					requestDetails["url"] = getString(req, "decoded_url")

					// 只保留关键的请求头
					if headers, ok := req["header"].(map[string]interface{}); ok {
						keyHeaders := make(map[string]interface{})
						for key, value := range headers {
							if key == "Authorization" || key == "Content-Type" || key == "User-Agent" {
								keyHeaders[key] = value
							}
						}
						requestDetails["headers"] = keyHeaders
					}

					// 只保留请求体的前100个字符
					if body, ok := req["body"].(string); ok && len(body) > 0 {
						if len(body) > 100 {
							requestDetails["body"] = body[:100] + "..."
						} else {
							requestDetails["body"] = body
						}
					}

					details["request"] = requestDetails
				}

				// 简化的响应信息
				if resp, ok := traceInfo["response"].(map[string]interface{}); ok {
					responseDetails := make(map[string]interface{})
					responseDetails["http_code"] = getInt(resp, "http_code")
					responseDetails["http_code_msg"] = getString(resp, "http_code_msg")

					// 只保留关键的响应头
					if headers, ok := resp["header"].(map[string]interface{}); ok {
						keyHeaders := make(map[string]interface{})
						for key, value := range headers {
							if key == "Content-Type" || key == "Trace-Id" {
								keyHeaders[key] = value
							}
						}
						responseDetails["headers"] = keyHeaders
					}

					// 只保留响应体的前200个字符
					if body, ok := resp["body"]; ok {
						bodyStr := fmt.Sprintf("%v", body)
						if len(bodyStr) > 200 {
							responseDetails["body"] = bodyStr[:200] + "..."
						} else {
							responseDetails["body"] = bodyStr
						}
					}

					details["response"] = responseDetails
				}

				log.Details = details
			}

			// 收集所有符合条件的日志
			allLogs = append(allLogs, log)
		}

		// ReadLineFromEnd已经按时间倒序返回了，最新的在前面

		// 只取最新的10条记录
		if len(allLogs) > maxLogs {
			logs = allLogs[:maxLogs]
		} else {
			logs = allLogs
		}

		// 获取最新一条日志的时间作为lastTime（最新的在数组开头）
		if len(logs) > 0 {
			lastTime = logs[0].Time
		}

		res := &logsResponse{
			Data:     logs,
			Total:    len(logs),
			Success:  true,
			LastTime: lastTime,
		}

		c.Payload(res)
	}
}

// extractRequestDetails 提取请求详情
func extractRequestDetails(fullLog map[string]interface{}) map[string]interface{} {
	request := make(map[string]interface{})

	// 从trace_info中提取请求详情
	if traceInfo, ok := fullLog["trace_info"].(map[string]interface{}); ok {
		if req, ok := traceInfo["request"].(map[string]interface{}); ok {
			// 提取请求头
			if headers, ok := req["header"].(map[string]interface{}); ok {
				request["headers"] = headers
			}

			// 提取请求体
			if body, ok := req["body"].(string); ok {
				request["body"] = body
			}

			// 提取解码后的URL
			if decodedURL, ok := req["decoded_url"].(string); ok {
				request["url"] = decodedURL
			}

			// 提取请求方法
			if method, ok := req["method"].(string); ok {
				request["method"] = method
			}
		}
	}

	// 从trace_info中提取响应详情
	if traceInfo, ok := fullLog["trace_info"].(map[string]interface{}); ok {
		if resp, ok := traceInfo["response"].(map[string]interface{}); ok {
			// 提取响应头
			if headers, ok := resp["header"].(map[string]interface{}); ok {
				request["response_headers"] = headers
			}

			// 提取响应体
			if body, ok := resp["body"].(map[string]interface{}); ok {
				request["response_body"] = body
			} else if body, ok := resp["body"].(string); ok {
				request["response_body"] = body
			}

			// 提取HTTP状态码
			if httpCode, ok := resp["http_code"].(float64); ok {
				request["response_http_code"] = int(httpCode)
			}
		}
	}

	return request
}

// extractResponseDetails 提取响应详情
func extractResponseDetails(fullLog map[string]interface{}) map[string]interface{} {
	response := make(map[string]interface{})

	// 从trace_info中提取响应详情
	if traceInfo, ok := fullLog["trace_info"].(map[string]interface{}); ok {
		if resp, ok := traceInfo["response"].(map[string]interface{}); ok {
			// 提取响应头
			if headers, ok := resp["header"].(map[string]interface{}); ok {
				response["headers"] = headers
			}

			// 提取响应体
			if body, ok := resp["body"].(map[string]interface{}); ok {
				response["body"] = body
			} else if body, ok := resp["body"].(string); ok {
				response["body"] = body
			}

			// 提取HTTP状态码
			if httpCode, ok := resp["http_code"].(float64); ok {
				response["http_code"] = int(httpCode)
			}

			// 提取HTTP状态码消息
			if httpCodeMsg, ok := resp["http_code_msg"].(string); ok {
				response["http_code_msg"] = httpCodeMsg
			}
		}
	}

	// 提取错误信息
	if error, ok := fullLog["error"].(string); ok {
		response["error"] = error
	}

	return response
}

// 辅助函数：安全地从map中获取字符串值
func getString(data map[string]interface{}, key string) string {
	if val, ok := data[key]; ok {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return ""
}

// 辅助函数：安全地从map中获取整数值
func getInt(data map[string]interface{}, key string) int {
	if val, ok := data[key]; ok {
		switch v := val.(type) {
		case int:
			return v
		case float64:
			return int(v)
		case string:
			if i, err := strconv.Atoi(v); err == nil {
				return i
			}
		}
	}
	return 0
}

// 辅助函数：安全地从map中获取布尔值
func getBool(data map[string]interface{}, key string) bool {
	if val, ok := data[key]; ok {
		if b, ok := val.(bool); ok {
			return b
		}
	}
	return false
}

// 辅助函数：安全地从map中获取浮点数值
func getFloat(data map[string]interface{}, key string) float64 {
	if val, ok := data[key]; ok {
		switch v := val.(type) {
		case float64:
			return v
		case int:
			return float64(v)
		case string:
			if f, err := strconv.ParseFloat(v, 64); err == nil {
				return f
			}
		}
	}
	return 0.0
}

// isLogRelatedRequest 检查是否为日志相关的请求
func isLogRelatedRequest(path string) bool {
	// 定义日志相关的路径（只过滤实际存在的URL）
	logPaths := []string{
		"/api/v1/logs",
		"/api/v1/logs/realtime",
		"/logs",
		"/tool/logs",
	}

	// 检查路径是否匹配日志相关的路径
	for _, logPath := range logPaths {
		if path == logPath {
			return true
		}
	}

	return false
}
