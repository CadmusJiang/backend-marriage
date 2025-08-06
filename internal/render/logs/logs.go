package logs

import (
	"github.com/xinliangnote/go-gin-api/internal/pkg/core"
	"github.com/xinliangnote/go-gin-api/internal/repository/mysql"
	"github.com/xinliangnote/go-gin-api/internal/repository/redis"

	"go.uber.org/zap"
)

type handler struct {
	logger *zap.Logger
	db     mysql.Repo
	cache  redis.Repo
}

func New(logger *zap.Logger, db mysql.Repo, cache redis.Repo) *handler {
	return &handler{
		logger: logger,
		db:     db,
		cache:  cache,
	}
}

// Viewer 日志查看页面
func (h *handler) Viewer() core.HandlerFunc {
	return func(ctx core.Context) {
		// 直接返回HTML内容，不依赖模板系统
		html := `<!DOCTYPE html>
<html lang="zh">
<head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1, maximum-scale=1, user-scalable=no"/>
    <title>实时日志查看器</title>
    <link href="/assets/bootstrap/css/bootstrap.min.css" rel="stylesheet">
    <link href="/assets/bootstrap/css/materialdesignicons.min.css" rel="stylesheet">
    <link href="/assets/bootstrap/css/style.min.css" rel="stylesheet">
    <style>
        .log-entry {
            margin-bottom: 10px;
            padding: 10px;
            border-radius: 5px;
            border-left: 4px solid #007bff;
            background-color: #f8f9fa;
        }
        .log-entry.success {
            border-left-color: #28a745;
            background-color: #d4edda;
        }
        .log-entry.error {
            border-left-color: #dc3545;
            background-color: #f8d7da;
        }
        .log-method {
            font-weight: bold;
            color: #007bff;
        }
        .log-path {
            color: #6c757d;
            font-family: monospace;
        }
        .log-time {
            color: #6c757d;
            font-size: 0.9em;
        }
        .log-duration {
            color: #28a745;
            font-weight: bold;
        }
        .log-status {
            font-weight: bold;
        }
        .log-status.success {
            color: #28a745;
        }
        .log-status.error {
            color: #dc3545;
        }
        .log-details {
            margin-top: 5px;
            font-size: 0.9em;
            color: #6c757d;
        }
        .log-container {
            max-height: 80vh;
            overflow-y: auto;
        }
        .filter-bar {
            position: sticky;
            top: 0;
            background: white;
            z-index: 100;
            padding: 10px 0;
            border-bottom: 1px solid #dee2e6;
        }
        .detail-panel {
            margin-top: 10px;
            padding: 10px;
            background-color: #f8f9fa;
            border-radius: 5px;
            border: 1px solid #dee2e6;
        }
        .detail-section {
            margin-bottom: 15px;
        }
        .detail-section h6 {
            color: #495057;
            font-weight: bold;
            margin-bottom: 5px;
        }
        .detail-content {
            background-color: #ffffff;
            padding: 10px;
            border-radius: 3px;
            border: 1px solid #e9ecef;
            font-family: monospace;
            font-size: 0.85em;
            max-height: 200px;
            overflow-y: auto;
        }
        .json-content {
            white-space: pre-wrap;
            word-break: break-all;
        }
        .auto-refresh-indicator {
            position: fixed;
            top: 60px;
            right: 20px;
            z-index: 1000;
            background: rgba(40, 167, 69, 0.9);
            color: white;
            padding: 5px 10px;
            border-radius: 15px;
            font-size: 0.8em;
            display: none;
        }
        .refresh-interval {
            position: fixed;
            top: 90px;
            right: 20px;
            z-index: 1000;
            background: rgba(0, 123, 255, 0.9);
            color: white;
            padding: 5px 10px;
            border-radius: 15px;
            font-size: 0.8em;
        }
    </style>
</head>

<body>
<div class="container-fluid p-t-15">
    <div class="row">
        <div class="col-12">
            <div class="card">
                <header class="card-header">
                    <div class="card-title">
                        <i class="mdi mdi-console"></i> 实时日志查看器
                        <small class="text-muted">监控后端API请求日志</small>
                    </div>
                    <ul class="card-actions">
                        <li>
                            <button id="autoRefreshBtn" class="btn btn-sm btn-outline-primary">
                                <i class="mdi mdi-refresh"></i> 自动刷新
                            </button>
                        </li>
                        <li>
                            <button id="clearLogsBtn" class="btn btn-sm btn-outline-secondary">
                                <i class="mdi mdi-delete"></i> 清空
                            </button>
                        </li>
                        <li>
                            <select id="refreshInterval" class="form-control form-control-sm">
                                <option value="1000">1秒</option>
                                <option value="2000" selected>2秒</option>
                                <option value="5000">5秒</option>
                                <option value="10000">10秒</option>
                            </select>
                        </li>
                    </ul>
                </header>
                
                <div class="card-body">
                    <!-- 过滤器 -->
                    <div class="filter-bar">
                        <div class="row">
                            <div class="col-md-3">
                                <select id="methodFilter" class="form-control form-control-sm">
                                    <option value="">所有方法</option>
                                    <option value="GET">GET</option>
                                    <option value="POST">POST</option>
                                    <option value="PUT">PUT</option>
                                    <option value="DELETE">DELETE</option>
                                </select>
                            </div>
                            <div class="col-md-3">
                                <select id="statusFilter" class="form-control form-control-sm">
                                    <option value="">所有状态</option>
                                    <option value="success">成功</option>
                                    <option value="error">失败</option>
                                </select>
                            </div>
                            <div class="col-md-4">
                                <input type="text" id="pathFilter" class="form-control form-control-sm" placeholder="路径过滤...">
                            </div>
                            <div class="col-md-2">
                                <button id="applyFilter" class="btn btn-sm btn-primary">应用过滤</button>
                            </div>
                        </div>
                    </div>

                    <!-- 统计信息 -->
                    <div class="row mb-3">
                        <div class="col-md-3">
                            <div class="card bg-primary text-white">
                                <div class="card-body p-2">
                                    <small>总请求数</small>
                                    <div id="totalRequests">0</div>
                                </div>
                            </div>
                        </div>
                        <div class="col-md-3">
                            <div class="card bg-success text-white">
                                <div class="card-body p-2">
                                    <small>成功请求</small>
                                    <div id="successRequests">0</div>
                                </div>
                            </div>
                        </div>
                        <div class="col-md-3">
                            <div class="card bg-danger text-white">
                                <div class="card-body p-2">
                                    <small>失败请求</small>
                                    <div id="errorRequests">0</div>
                                </div>
                            </div>
                        </div>
                        <div class="col-md-3">
                            <div class="card bg-info text-white">
                                <div class="card-body p-2">
                                    <small>平均响应时间</small>
                                    <div id="avgResponseTime">0ms</div>
                                </div>
                            </div>
                        </div>
                    </div>

                    <!-- 日志容器 -->
                    <div id="logContainer" class="log-container">
                        <div class="text-center text-muted">
                            <i class="mdi mdi-loading mdi-spin"></i> 正在加载日志...
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </div>
</div>

<!-- 自动刷新指示器 -->
<div id="autoRefreshIndicator" class="auto-refresh-indicator">
    <i class="mdi mdi-refresh mdi-spin"></i> 自动刷新中
</div>

<div id="refreshIntervalDisplay" class="refresh-interval">
    刷新间隔: <span id="intervalText">2秒</span>
</div>

<script type="text/javascript" src="/assets/bootstrap/js/jquery.min.js"></script>
<script type="text/javascript" src="/assets/bootstrap/js/popper.min.js"></script>
<script type="text/javascript" src="/assets/bootstrap/js/bootstrap.min.js"></script>

<script>
let autoRefresh = false;
let refreshInterval;
let lastTime = '';
let allLogs = [];
let filteredLogs = [];
let currentRefreshInterval = 2000;

// 初始化
$(document).ready(function() {
    loadLogs();
    
    // 自动刷新按钮
    $('#autoRefreshBtn').click(function() {
        autoRefresh = !autoRefresh;
        if (autoRefresh) {
            $(this).html('<i class="mdi mdi-pause"></i> 停止刷新');
            $(this).removeClass('btn-outline-primary').addClass('btn-primary');
            $('#autoRefreshIndicator').show();
            startAutoRefresh();
        } else {
            $(this).html('<i class="mdi mdi-refresh"></i> 自动刷新');
            $(this).removeClass('btn-primary').addClass('btn-outline-primary');
            $('#autoRefreshIndicator').hide();
            stopAutoRefresh();
        }
    });

    // 清空日志
    $('#clearLogsBtn').click(function() {
        allLogs = [];
        filteredLogs = [];
        updateLogDisplay();
        updateStats();
    });

    // 应用过滤
    $('#applyFilter').click(function() {
        applyFilters();
    });

    // 回车键应用过滤
    $('#pathFilter').keypress(function(e) {
        if (e.which == 13) {
            applyFilters();
        }
    });

    // 刷新间隔选择
    $('#refreshInterval').change(function() {
        currentRefreshInterval = parseInt($(this).val());
        $('#intervalText').text((currentRefreshInterval / 1000) + '秒');
        if (autoRefresh) {
            stopAutoRefresh();
            startAutoRefresh();
        }
    });
});

// 加载日志
function loadLogs() {
    $.ajax({
        url: '/api/v1/logs/realtime',
        method: 'GET',
        success: function(response) {
            if (response.success) {
                allLogs = response.data;
                lastTime = response.last_time;
                applyFilters();
                updateStats();
            }
        },
        error: function(xhr, status, error) {
            $('#logContainer').html('<div class="alert alert-danger">加载日志失败: ' + error + '</div>');
        }
    });
}

// 加载实时日志
function loadRealtimeLogs() {
    $.ajax({
        url: '/api/v1/logs/realtime',
        method: 'GET',
        headers: {
            'since': lastTime
        },
        success: function(response) {
            if (response.success && response.data.length > 0) {
                // 添加新日志到开头
                allLogs = response.data.concat(allLogs);
                lastTime = response.last_time;
                
                // 限制日志数量
                if (allLogs.length > 1000) {
                    allLogs = allLogs.slice(0, 1000);
                }
                
                applyFilters();
                updateStats();
            }
        }
    });
}

// 应用过滤器
function applyFilters() {
    const methodFilter = $('#methodFilter').val();
    const statusFilter = $('#statusFilter').val();
    const pathFilter = $('#pathFilter').val().toLowerCase();

    filteredLogs = allLogs.filter(log => {
        let match = true;
        
        if (methodFilter && log.method !== methodFilter) {
            match = false;
        }
        
        if (statusFilter) {
            if (statusFilter === 'success' && !log.success) {
                match = false;
            } else if (statusFilter === 'error' && log.success) {
                match = false;
            }
        }
        
        if (pathFilter && !log.path.toLowerCase().includes(pathFilter)) {
            match = false;
        }
        
        return match;
    });

    updateLogDisplay();
}

// 更新日志显示
function updateLogDisplay() {
    const container = $('#logContainer');
    
    if (filteredLogs.length === 0) {
        container.html('<div class="text-center text-muted">暂无日志记录</div>');
        return;
    }

    let html = '';
    filteredLogs.forEach((log, index) => {
        const statusClass = log.success ? 'success' : 'error';
        const statusText = log.success ? '成功' : '失败';
        const methodClass = getMethodClass(log.method);
        
        html += '<div class="log-entry ' + statusClass + '">';
        html += '<div class="d-flex justify-content-between align-items-start">';
        html += '<div>';
        html += '<span class="log-method ' + methodClass + '">' + log.method + '</span>';
        html += '<span class="log-path">' + log.path + '</span>';
        html += '<span class="log-status ' + statusClass + '">' + statusText + '</span>';
        html += '<span class="log-duration">' + (log.cost_seconds * 1000).toFixed(0) + 'ms</span>';
        html += '</div>';
        html += '<div class="log-time">' + formatTime(log.time) + '</div>';
        html += '</div>';
        html += '<div class="log-details">';
        html += '<small>TraceID: ' + log.trace_id + ' | HTTP: ' + log.http_code + '</small>';
        html += '<button class="btn btn-sm btn-outline-info ml-2" onclick="toggleDetails(' + index + ')">';
        html += '<i class="mdi mdi-eye"></i> 查看详情';
        html += '</button>';
        html += '</div>';
        html += '<div id="details-' + index + '" class="detail-panel" style="display: none;">';
        html += generateDetailsHtml(log);
        html += '</div>';
        html += '</div>';
    });

    container.html(html);
}

// 生成详情HTML
function generateDetailsHtml(log) {
    let html = '';
    
    // 请求详情
    if (log.details && log.details.request) {
        html += '<div class="detail-section">';
        html += '<h6><i class="mdi mdi-arrow-up text-primary"></i> 请求详情</h6>';
        
        if (log.details.request.headers) {
            html += '<div class="detail-content">';
            html += '<strong>请求头:</strong><br>';
            html += '<div class="json-content">' + JSON.stringify(log.details.request.headers, null, 2) + '</div>';
            html += '</div>';
        }
        
        if (log.details.request.body) {
            html += '<div class="detail-content mt-2">';
            html += '<strong>请求体:</strong><br>';
            html += '<div class="json-content">' + log.details.request.body + '</div>';
            html += '</div>';
        }
        
        if (log.details.request.query) {
            html += '<div class="detail-content mt-2">';
            html += '<strong>查询参数:</strong><br>';
            html += '<div class="json-content">' + JSON.stringify(log.details.request.query, null, 2) + '</div>';
            html += '</div>';
        }
        
        html += '</div>';
    }
    
    // 响应详情
    if (log.details && log.details.response) {
        html += '<div class="detail-section">';
        html += '<h6><i class="mdi mdi-arrow-down text-success"></i> 响应详情</h6>';
        
        if (log.details.response.headers) {
            html += '<div class="detail-content">';
            html += '<strong>响应头:</strong><br>';
            html += '<div class="json-content">' + JSON.stringify(log.details.response.headers, null, 2) + '</div>';
            html += '</div>';
        }
        
        if (log.details.response.body) {
            html += '<div class="detail-content mt-2">';
            html += '<strong>响应体:</strong><br>';
            html += '<div class="json-content">' + log.details.response.body + '</div>';
            html += '</div>';
        }
        
        if (log.details.response.error) {
            html += '<div class="detail-content mt-2">';
            html += '<strong>错误信息:</strong><br>';
            html += '<div class="text-danger">' + log.details.response.error + '</div>';
            html += '</div>';
        }
        
        html += '</div>';
    }
    
    return html;
}

// 切换详情显示
function toggleDetails(index) {
    const detailsPanel = $('#details-' + index);
    if (detailsPanel.is(':visible')) {
        detailsPanel.hide();
    } else {
        detailsPanel.show();
    }
}

// 获取方法样式类
function getMethodClass(method) {
    const classes = {
        'GET': 'text-primary',
        'POST': 'text-success',
        'PUT': 'text-warning',
        'DELETE': 'text-danger',
        'PATCH': 'text-info'
    };
    return classes[method] || 'text-secondary';
}

// 格式化时间
function formatTime(timeStr) {
    const date = new Date(timeStr);
    return date.toLocaleTimeString();
}

// 更新统计信息
function updateStats() {
    const total = allLogs.length;
    const success = allLogs.filter(log => log.success).length;
    const error = total - success;
    const avgTime = allLogs.length > 0 
        ? allLogs.reduce((sum, log) => sum + log.cost_seconds, 0) / allLogs.length * 1000
        : 0;

    $('#totalRequests').text(total);
    $('#successRequests').text(success);
    $('#errorRequests').text(error);
    $('#avgResponseTime').text(avgTime.toFixed(0) + 'ms');
}

// 开始自动刷新
function startAutoRefresh() {
    refreshInterval = setInterval(function() {
        loadRealtimeLogs();
    }, currentRefreshInterval);
}

// 停止自动刷新
function stopAutoRefresh() {
    if (refreshInterval) {
        clearInterval(refreshInterval);
        refreshInterval = null;
    }
}
</script>
</body>
</html>`

		ctx.ResponseWriter().Write([]byte(html))
	}
}
