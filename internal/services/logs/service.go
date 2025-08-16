package logs

import (
	"bufio"
	"encoding/json"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type LogEntry struct {
	Raw       string                 `json:"raw"`
	Message   string                 `json:"message"`
	Type      string                 `json:"type"`
	Success   bool                   `json:"success"`
	Timestamp string                 `json:"timestamp"`
	Details   map[string]interface{} `json:"details,omitempty"`
	// 新增字段用于HTTP日志
	Method     string  `json:"method,omitempty"`
	Path       string  `json:"path,omitempty"`
	StatusCode int     `json:"statusCode,omitempty"`
	Duration   float64 `json:"duration,omitempty"`
	TraceID    string  `json:"traceId,omitempty"`
}

type LogStats struct {
	Total     int `json:"total"`
	HTTP      int `json:"http"`
	MySQL     int `json:"mysql"`
	Redis     int `json:"redis"`
	System    int `json:"system"`
	WithSQL   int `json:"with_sql"`
	WithRedis int `json:"with_redis"`
}

// 新增分页响应结构
type PaginatedLogsResponse struct {
	Logs       []LogEntry `json:"logs"`
	Stats      LogStats   `json:"stats"`
	Pagination struct {
		CurrentPage  int   `json:"current_page"`
		PageSize     int   `json:"page_size"`
		TotalLines   int64 `json:"total_lines"`
		TotalPages   int   `json:"total_pages"`
		HasMore      bool  `json:"has_more"`
		NextPage     int   `json:"next_page,omitempty"`
		PreviousPage int   `json:"previous_page,omitempty"`
	} `json:"pagination"`
	Timestamp string `json:"timestamp"`
}

type LogService struct{}

func New() *LogService {
	return &LogService{}
}

// GetLatestLogs 获取最近的日志（保持向后兼容）
func (s *LogService) GetLatestLogs(filePath string, maxBytes int64, n int) ([]string, error) {
	lines, err := s.readLastNLines(filePath, maxBytes, n)
	if err != nil {
		return nil, err
	}
	return lines, nil
}

// GetUnifiedLogs 获取统一的日志数据（保持向后兼容）
func (s *LogService) GetUnifiedLogs(filePath string, maxBytes int64) ([]LogEntry, LogStats, error) {
	lines, err := s.readLastNLines(filePath, maxBytes, 1000)
	if err != nil {
		return nil, LogStats{}, err
	}

	var logs []LogEntry
	var stats = LogStats{}

	for _, line := range lines {
		if line == "" {
			continue
		}

		// 解析日志行
		logEntry := s.parseLogLine(line)
		logs = append(logs, logEntry)

		// 更新统计信息
		s.updateStats(&stats, logEntry)
	}

	return logs, stats, nil
}

// 新增：获取分页日志数据
func (s *LogService) GetPaginatedLogs(filePath string, page, pageSize int, logType string) (*PaginatedLogsResponse, error) {
	// 获取文件总行数
	totalLines, err := s.getFileLineCount(filePath)
	if err != nil {
		return nil, err
	}

	// 计算分页信息
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100 // 限制最大页面大小
	}

	totalPages := int((totalLines + int64(pageSize) - 1) / int64(pageSize))
	if page > totalPages && totalPages > 0 {
		page = totalPages
	}

	// 计算起始和结束行号（从文件末尾开始，最新的日志在前面）
	startLine := totalLines - int64((page-1)*pageSize)
	endLine := totalLines - int64((page-1)*pageSize) + int64(pageSize)

	if startLine < 1 {
		startLine = 1
	}
	if endLine > totalLines {
		endLine = totalLines
	}

	// 使用现有的readLastNLines方法来读取日志，然后进行分页
	// 读取足够多的行来覆盖我们需要的范围，并且为过滤预留空间
	linesToRead := int(endLine-startLine+1) + 500                       // 增加缓冲行数，为过滤预留空间
	lines, err := s.readLastNLines(filePath, int64(1<<20), linesToRead) // 1MB
	if err != nil {
		return nil, err
	}

	// 反转lines数组，因为readLastNLines返回的是最新的行
	// 我们需要从最新的行开始计算
	reversedLines := make([]string, len(lines))
	for i, j := 0, len(lines)-1; i < len(lines); i, j = i+1, j-1 {
		reversedLines[i] = lines[j]
	}

	// 计算实际需要返回的行
	var resultLines []string
	if int64(len(reversedLines)) >= startLine {
		start := startLine - 1
		end := endLine
		if end > int64(len(reversedLines)) {
			end = int64(len(reversedLines))
		}
		if start < int64(len(reversedLines)) {
			resultLines = reversedLines[start:end]
		}
	}

	// 如果resultLines为空，尝试直接从文件末尾读取指定行数
	if len(resultLines) == 0 {
		// 直接从文件末尾读取指定行数
		resultLines, err = s.readLastNLines(filePath, int64(1<<20), pageSize*10) // 读取更多行用于过滤
		if err != nil {
			return nil, err
		}
	}

	// 解析日志
	var logs []LogEntry
	var stats = LogStats{}

	for _, line := range resultLines {
		if line == "" {
			continue
		}

		// 过滤掉系统内部的API调用日志
		if strings.Contains(line, "/system/logs/paginated") {
			continue
		}

		logEntry := s.parseLogLine(line)

		// 根据日志类型进行过滤
		if logType != "" && logType != "all" {
			if logType == "http" && logEntry.Type != "http" {
				continue
			}
			if (logType == "sql" || logType == "mysql") && logEntry.Type != "mysql" {
				continue
			}
			if logType == "redis" && logEntry.Type != "redis" {
				continue
			}
			if logType == "system" && logEntry.Type != "system" {
				continue
			}
		}

		logs = append(logs, logEntry)
		s.updateStats(&stats, logEntry)

		// 如果已经收集到足够的日志，就停止
		if len(logs) >= pageSize {
			break
		}
	}

	// 构建分页响应
	response := &PaginatedLogsResponse{
		Logs:      logs,
		Stats:     stats,
		Timestamp: time.Now().Format("2006-01-02 15:04:05"),
	}

	response.Pagination.CurrentPage = page
	response.Pagination.PageSize = pageSize
	response.Pagination.TotalLines = totalLines
	response.Pagination.TotalPages = totalPages
	response.Pagination.HasMore = page < totalPages
	if page < totalPages {
		response.Pagination.NextPage = page + 1
	}
	if page > 1 {
		response.Pagination.PreviousPage = page - 1
	}

	return response, nil
}

// 新增：获取文件总行数（使用更简单的方法）
func (s *LogService) getFileLineCount(filePath string) (int64, error) {
	// 使用现有的readLastNLines方法，传入一个很大的行数限制来获取所有行
	lines, err := s.readLastNLines(filePath, int64(1<<30), 0) // 1GB, 0表示不限制行数
	if err != nil {
		return 0, err
	}
	return int64(len(lines)), nil
}

// 新增：按行号范围读取文件
func (s *LogService) readLinesByRange(filePath string, startLine, endLine int64) ([]string, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var lines []string
	scanner := bufio.NewScanner(f)
	// 设置更大的缓冲区来处理长行
	buf := make([]byte, 0, 64*1024) // 64KB缓冲区
	scanner.Buffer(buf, 1024*1024)  // 最大1MB

	var currentLine int64 = 0

	for scanner.Scan() {
		currentLine++
		if currentLine >= startLine && currentLine <= endLine {
			lines = append(lines, scanner.Text())
		}
		if currentLine > endLine {
			break
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return lines, nil
}

// readLastNLines 读取文件的最后N行
func (s *LogService) readLastNLines(filePath string, maxBytes int64, n int) ([]string, error) {
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

	content := string(data)
	if start > 0 {
		if idx := strings.IndexByte(content, '\n'); idx >= 0 {
			content = content[idx+1:]
		}
	}

	lines := strings.Split(content, "\n")
	if len(lines) > 0 && lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	}

	// 如果指定了行数限制，只返回最后N行
	if n > 0 && len(lines) > n {
		lines = lines[len(lines)-n:]
	}

	return lines, nil
}

// parseLogLine 解析单行日志
func (s *LogService) parseLogLine(line string) LogEntry {
	logEntry := LogEntry{
		Raw:       line,
		Message:   line,
		Type:      "system",
		Success:   true,
		Timestamp: time.Now().Format("2006-01-02 15:04:05"),
		Details:   make(map[string]interface{}),
	}

	// 尝试解析JSON格式的日志
	if s.parseJSONLog(line, &logEntry) {
		return logEntry
	}

	// 尝试解析trace-log格式的日志
	if s.parseTraceLog(line, &logEntry) {
		return logEntry
	}

	// 尝试解析HTTP请求日志
	if s.parseHTTPLog(line, &logEntry) {
		return logEntry
	}

	// 只有在没有通过其他方式识别类型时，才进行通用类型检测
	if logEntry.Type == "system" {
		// 优先检查trace-log
		if strings.Contains(line, "trace-log") {
			logEntry.Type = "http"
			logEntry.Message = "HTTP请求日志"
		} else if strings.Contains(line, "SQL") || strings.Contains(line, "mysql") || strings.Contains(line, "query") || strings.Contains(line, "sqls") {
			logEntry.Type = "mysql"
		} else if strings.Contains(line, "redis") || strings.Contains(line, "Redis") || strings.Contains(line, "with_redis") {
			logEntry.Type = "redis"
		} else if strings.Contains(line, "router") || strings.Contains(line, "core") || strings.Contains(line, "domain") || strings.Contains(line, "caller") {
			logEntry.Type = "system"
		} else if strings.Contains(line, "HTTP") || strings.Contains(line, "GET") || strings.Contains(line, "POST") || strings.Contains(line, "PUT") || strings.Contains(line, "DELETE") {
			logEntry.Type = "http"
		}
	}

	// 对于trace-log日志，进一步检查是否包含SQL或Redis信息
	if logEntry.Type == "http" && strings.Contains(line, "trace-log") {
		if strings.Contains(line, "sqls") && !strings.Contains(line, `"sqls":null`) {
			logEntry.Type = "mysql"
			logEntry.Message = "MySQL操作日志"
		} else if strings.Contains(line, "redis") && !strings.Contains(line, `"redis":null`) {
			logEntry.Type = "redis"
			logEntry.Message = "Redis操作日志"
		}
	}

	// 最后检查：如果日志包含trace-log但没有被识别为HTTP类型，强制设置为HTTP类型
	if strings.Contains(line, "trace-log") && logEntry.Type != "http" && logEntry.Type != "mysql" && logEntry.Type != "redis" {
		logEntry.Type = "http"
		logEntry.Message = "HTTP请求日志"
	}

	// 额外检查：如果日志包含sqls字段，优先识别为SQL类型
	if strings.Contains(line, "sqls") && !strings.Contains(line, `"sqls":null`) {
		logEntry.Type = "mysql"
		logEntry.Message = "MySQL操作日志"
	}

	// 检查trace-log日志中的SQL信息
	if strings.Contains(line, "trace-log") && strings.Contains(line, "sqls") {
		// 如果包含sqls字段且不是null，则识别为SQL类型
		if !strings.Contains(line, `"sqls":null`) {
			logEntry.Type = "mysql"
			logEntry.Message = "MySQL操作日志"
		}
	}

	return logEntry
}

// parseJSONLog 尝试解析JSON格式的日志
func (s *LogService) parseJSONLog(line string, logEntry *LogEntry) bool {
	var jsonData map[string]interface{}
	if err := json.Unmarshal([]byte(line), &jsonData); err != nil {
		// 如果JSON解析失败，尝试使用parseTraceLog
		return false
	}

	// 检查是否是trace-log格式
	if msg, ok := jsonData["msg"].(string); ok && msg == "trace-log" {
		logEntry.Type = "http"
		logEntry.Message = "HTTP请求日志"

		// 提取基本信息（从根级别）
		if method, ok := jsonData["method"].(string); ok {
			logEntry.Method = method
		}
		if path, ok := jsonData["path"].(string); ok {
			logEntry.Path = path
		}
		if httpCode, ok := jsonData["http_code"].(float64); ok {
			logEntry.StatusCode = int(httpCode)
			logEntry.Success = httpCode == 200
		}
		if costSeconds, ok := jsonData["cost_seconds"].(float64); ok {
			logEntry.Duration = costSeconds
		}
		if traceID, ok := jsonData["trace_id"].(string); ok {
			logEntry.TraceID = traceID
		}

		// 如果根级别没有找到，尝试从trace_info中提取
		if logEntry.Method == "" {
			if traceInfo, ok := jsonData["trace_info"].(map[string]interface{}); ok {
				if request, ok := traceInfo["request"].(map[string]interface{}); ok {
					if method, ok := request["method"].(string); ok {
						logEntry.Method = method
					}
					if path, ok := request["decoded_url"].(string); ok {
						logEntry.Path = path
					}
				}
				if response, ok := traceInfo["response"].(map[string]interface{}); ok {
					if httpCode, ok := response["http_code"].(float64); ok {
						logEntry.StatusCode = int(httpCode)
						logEntry.Success = httpCode == 200
					}
					if costSeconds, ok := response["cost_seconds"].(float64); ok {
						logEntry.Duration = costSeconds
					}
				}
				if traceID, ok := traceInfo["trace_id"].(string); ok {
					logEntry.TraceID = traceID
				}
			}
		}

		// 保存原始数据
		logEntry.Details["原始数据"] = jsonData

		return true
	}

	return false
}

// parseTraceLog 解析trace-log格式的日志
func (s *LogService) parseTraceLog(line string, logEntry *LogEntry) bool {
	// 检查是否包含trace-log关键词
	if !strings.Contains(line, "trace-log") {
		return false
	}

	logEntry.Type = "http"
	logEntry.Message = "HTTP请求日志"

	// 使用正则表达式提取信息
	// 提取方法
	if methodMatch := regexp.MustCompile(`"method":"([^"]+)"`).FindStringSubmatch(line); len(methodMatch) > 1 {
		logEntry.Method = methodMatch[1]
	}

	// 提取路径
	if pathMatch := regexp.MustCompile(`"path":"([^"]+)"`).FindStringSubmatch(line); len(pathMatch) > 1 {
		logEntry.Path = pathMatch[1]
	}

	// 提取HTTP状态码
	if statusMatch := regexp.MustCompile(`"http_code":(\d+)`).FindStringSubmatch(line); len(statusMatch) > 1 {
		if statusCode := s.parseInt(statusMatch[1]); statusCode > 0 {
			logEntry.StatusCode = statusCode
			logEntry.Success = statusCode == 200
		}
	}

	// 提取耗时
	if durationMatch := regexp.MustCompile(`"cost_seconds":([\d.]+)`).FindStringSubmatch(line); len(durationMatch) > 1 {
		if duration := s.parseFloat(durationMatch[1]); duration > 0 {
			logEntry.Duration = duration
		}
	}

	// 提取trace ID
	if traceMatch := regexp.MustCompile(`"trace_id":"([^"]+)"`).FindStringSubmatch(line); len(traceMatch) > 1 {
		logEntry.TraceID = traceMatch[1]
	}

	return true
}

// parseHTTPLog 解析HTTP请求日志
func (s *LogService) parseHTTPLog(line string, logEntry *LogEntry) bool {
	// 检查是否包含HTTP请求相关信息
	if !strings.Contains(line, "HTTP") && !strings.Contains(line, "GET") && !strings.Contains(line, "POST") &&
		!strings.Contains(line, "PUT") && !strings.Contains(line, "DELETE") && !strings.Contains(line, "PATCH") {
		return false
	}

	logEntry.Type = "http"
	logEntry.Message = "HTTP请求日志"

	// 提取HTTP方法
	methods := []string{"GET", "POST", "PUT", "DELETE", "PATCH", "HEAD", "OPTIONS"}
	for _, method := range methods {
		if strings.Contains(line, method) {
			logEntry.Method = method
			break
		}
	}

	// 提取路径（简单的URL模式匹配）
	if urlMatch := regexp.MustCompile(`/[a-zA-Z0-9/_-]+`).FindString(line); urlMatch != "" {
		logEntry.Path = urlMatch
	}

	// 提取状态码
	if statusMatch := regexp.MustCompile(`(\d{3})`).FindStringSubmatch(line); len(statusMatch) > 1 {
		if statusCode := s.parseInt(statusMatch[1]); statusCode >= 100 && statusCode < 600 {
			logEntry.StatusCode = statusCode
			logEntry.Success = statusCode >= 200 && statusCode < 400
		}
	}

	// 提取耗时
	if durationMatch := regexp.MustCompile(`(\d+\.\d+)s`).FindStringSubmatch(line); len(durationMatch) > 1 {
		if duration := s.parseFloat(durationMatch[1]); duration > 0 {
			logEntry.Duration = duration
		}
	}

	return true
}

// parseInt 安全地解析整数
func (s *LogService) parseInt(str string) int {
	if i, err := strconv.Atoi(str); err == nil {
		return i
	}
	return 0
}

// parseFloat 安全地解析浮点数
func (s *LogService) parseFloat(str string) float64 {
	// 先尝试直接解析为float64
	if f, err := strconv.ParseFloat(str, 64); err == nil {
		return f
	}
	// 如果失败，再尝试解析为时间格式
	if f, err := time.ParseDuration(str); err == nil {
		return f.Seconds()
	}
	return 0
}

// updateStats 更新统计信息
func (s *LogService) updateStats(stats *LogStats, logEntry LogEntry) {
	stats.Total++

	switch logEntry.Type {
	case "http":
		stats.HTTP++
	case "mysql":
		stats.MySQL++
	case "redis":
		stats.Redis++
	default:
		stats.System++
	}

	// 检测是否包含SQL或Redis操作
	if strings.Contains(logEntry.Raw, "SQL") || strings.Contains(logEntry.Raw, "mysql") || strings.Contains(logEntry.Raw, "query") {
		stats.WithSQL++
	}
	if strings.Contains(logEntry.Raw, "redis") || strings.Contains(logEntry.Raw, "Redis") {
		stats.WithRedis++
	}
}
