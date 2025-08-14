package analytics

import (
	"net/http"
	"time"

	"github.com/xinliangnote/go-gin-api/internal/code"
	"github.com/xinliangnote/go-gin-api/internal/pkg/core"
)

type teamsTrendsRequest struct {
	Metric    string `form:"metric"`
	TimeRange string `form:"time"`
	Top       int    `form:"top"`
}

type teamTrendSeries struct {
	Name   string       `json:"name"`
	TeamId string       `json:"teamId"`
	Trends []trendPoint `json:"trends"`
}

type teamsTrendsResponse struct {
	Success bool `json:"success"`
	Data    struct {
		Trends []teamTrendSeries `json:"trends"`
		Meta   analyticsMeta     `json:"meta"`
	} `json:"data"`
}

// GetTeamsTrends 获取团队多序列趋势（多队折线）
// @Summary 获取团队多序列趋势（多队折线）
// @Description 返回多支团队在时间范围内的趋势
// @Tags Analytics
// @Accept application/x-www-form-urlencoded
// @Produce json
// @Param metric query string false "指标" Enums(paid_count, lead_count) default(paid_count)
// @Param time query string false "时间范围" Enums(today, week, last7days, month) default(month)
// @Param top query int false "返回前N条折线"
// @Success 200 {object} teamsTrendsResponse
// @Failure 400 {object} code.Failure
// @Router /api/v1/analytics/teams/trends [get]
func (h *handler) GetTeamsTrends() core.HandlerFunc {
	return func(c core.Context) {
		req := new(teamsTrendsRequest)
		if err := c.ShouldBindForm(req); err != nil {
			c.AbortWithError(core.Error(http.StatusBadRequest, code.ParamBindError, code.Text(code.ParamBindError)).WithError(err))
			return
		}
		if req.Metric == "" {
			req.Metric = "paid_count"
		}
		if req.TimeRange == "" {
			req.TimeRange = "month"
		}

		start, end := resolveTimeRange(req.TimeRange)
		var points []trendPoint
		for t := start; t.Before(end); t = t.Add(24 * time.Hour) {
			points = append(points, trendPoint{Date: t, Count: float64(10 + t.Day()%5)})
		}

		series := []teamTrendSeries{
			{Name: "团队A", TeamId: "t-1", Trends: points},
			{Name: "团队B", TeamId: "t-2", Trends: points},
			{Name: "团队C", TeamId: "t-3", Trends: points},
		}
		if req.Top > 0 && req.Top < len(series) {
			series = series[:req.Top]
		}

		res := new(teamsTrendsResponse)
		res.Success = true
		res.Data.Trends = series
		res.Data.Meta = analyticsMeta{Metric: req.Metric, Time: req.TimeRange, DateRange: &dateSpan{Start: start, End: end}}
		c.Payload(res)
	}
}
