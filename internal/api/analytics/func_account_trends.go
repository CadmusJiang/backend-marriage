package analytics

import (
	"net/http"
	"time"

	"github.com/xinliangnote/go-gin-api/internal/authz"
	"github.com/xinliangnote/go-gin-api/internal/code"
	"github.com/xinliangnote/go-gin-api/internal/pkg/core"
)

type accountTrendsRequest struct {
	Metric        string `form:"metric"`
	TimeRange     string `form:"time"`
	Top           int    `form:"top"`
	BelongGroupId string `form:"belongGroupId"`
}

type trendPoint struct {
	Date  time.Time `json:"date"`
	Count float64   `json:"count"`
}

type accountTrendSeries struct {
	Name      string       `json:"name"`
	AccountId string       `json:"accountId"`
	Trends    []trendPoint `json:"trends"`
}

type accountTrendsResponse struct {
	Success bool `json:"success"`
	Data    struct {
		Trends []accountTrendSeries `json:"trends"`
		Meta   analyticsMeta        `json:"meta"`
	} `json:"data"`
}

// GetAccountTrends 获取个人趋势数据
// @Summary 获取个人趋势数据
// @Description 返回个人在开单或留资方面的趋势数据，支持按组筛选和时间范围
// @Tags Analytics
// @Accept application/x-www-form-urlencoded
// @Produce json
// @Param metric query string false "指标类型" Enums(paid_count, lead_count) default(paid_count) "paid_count: 开单数, lead_count: 留资数"
// @Param time query string false "时间范围" Enums(today, week, last7days, month) default(month) "today: 今天, week: 本周, last7days: 最近7天, month: 本月"
// @Param top query int false "返回前N名" default(10) "限制返回结果数量，默认10"
// @Param belongGroupId query string false "归属组ID" "按组筛选，不传则返回所有组数据"
// @Success 200 {object} accountTrendsResponse
// @Failure 400 {object} code.Failure
// @Router /api/v1/analytics/account/trends [get]
func (h *handler) GetAccountTrends() core.HandlerFunc {
	return func(c core.Context) {
		req := new(accountTrendsRequest)
		if err := c.ShouldBindForm(req); err != nil {
			c.AbortWithError(core.Error(http.StatusBadRequest, code.ParamBindError, code.Text(code.ParamBindError)).WithError(err))
			return
		}

		// 设置默认值和验证
		if req.Metric == "" {
			req.Metric = "paid_count" // 设置默认指标为开单数
		}

		// 验证指标类型
		validMetrics := map[string]bool{"paid_count": true, "lead_count": true}
		if !validMetrics[req.Metric] {
			c.AbortWithError(core.Error(http.StatusBadRequest, code.ParamBindError, "invalid metric, must be one of: paid_count, lead_count"))
			return
		}

		if req.TimeRange == "" {
			req.TimeRange = "month" // 默认时间范围为月
		}

		// 验证时间范围
		validTimeRanges := map[string]bool{"today": true, "week": true, "last7days": true, "month": true}
		if !validTimeRanges[req.TimeRange] {
			c.AbortWithError(core.Error(http.StatusBadRequest, code.ParamBindError, "invalid time range, must be one of: today, week, last7days, month"))
			return
		}

		// 设置默认Top值
		if req.Top <= 0 {
			req.Top = 10 // 默认返回前10名
		}

		// 范围上限：非公司管理员一律限定在“本组”
		scope, _ := authz.ComputeScope(c, h.db)
		if !scope.ScopeAll {
			if len(scope.AllowedGroupIDs) > 0 {
				req.BelongGroupId = formatFirst(scope.AllowedGroupIDs)
			} else {
				req.BelongGroupId = ""
			}
		}

		start, end := resolveTimeRange(req.TimeRange)
		// produce daily points between start and end for stub
		var points []trendPoint
		for t := start; t.Before(end); t = t.Add(24 * time.Hour) {
			points = append(points, trendPoint{Date: t, Count: float64(t.Day() % 7)})
		}

		series := []accountTrendSeries{
			{Name: "张三", AccountId: "acc-1", Trends: points},
			{Name: "李四", AccountId: "acc-2", Trends: points},
		}
		if req.Top > 0 && req.Top < len(series) {
			series = series[:req.Top]
		}

		res := new(accountTrendsResponse)
		res.Success = true
		res.Data.Trends = series
		res.Data.Meta = analyticsMeta{Metric: req.Metric, Time: req.TimeRange, BelongGroupId: req.BelongGroupId, DateRange: &dateSpan{Start: start, End: end}}
		c.Payload(res)
	}
}
