package analytics

import (
	"net/http"

	"github.com/xinliangnote/go-gin-api/internal/code"
	"github.com/xinliangnote/go-gin-api/internal/pkg/core"
)

type teamRankingsRequest struct {
	Metric    string `form:"metric"`
	TimeRange string `form:"time"`
	Top       int    `form:"top"`
}

type teamRankingsResponse struct {
	Success bool `json:"success"`
	Data    struct {
		Rankings []rankingItem `json:"rankings"`
		Meta     analyticsMeta `json:"meta"`
	} `json:"data"`
}

// GetTeamsRankings 获取团队排名（全体团队维度）
// @Summary 获取团队排名（全体团队维度）
// @Description 按开单或留资对团队进行排名
// @Tags Analytics
// @Accept application/x-www-form-urlencoded
// @Produce json
// @Param metric query string false "指标" Enums(paid_count, lead_count) default(paid_count)
// @Param time query string false "时间范围" Enums(today, week, last7days, month) default(month)
// @Param top query int false "返回前N名"
// @Success 200 {object} teamRankingsResponse
// @Failure 400 {object} code.Failure
// @Router /api/v1/analytics/teams/rankings [get]
func (h *handler) GetTeamsRankings() core.HandlerFunc {
	return func(c core.Context) {
		req := new(teamRankingsRequest)
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

		rankings := []rankingItem{
			{Name: "团队A", Value: 350, RankChange: -1},
			{Name: "团队B", Value: 300, RankChange: 1},
			{Name: "团队C", Value: 250, RankChange: 0},
		}
		if req.Top > 0 && req.Top < len(rankings) {
			rankings = rankings[:req.Top]
		}

		var avg float64
		for _, r := range rankings {
			avg += r.Value
		}
		if len(rankings) > 0 {
			avg = avg / float64(len(rankings))
		}
		start, end := resolveTimeRange(req.TimeRange)

		res := new(teamRankingsResponse)
		res.Success = true
		res.Data.Rankings = rankings
		res.Data.Meta = analyticsMeta{Metric: req.Metric, Time: req.TimeRange, AverageValue: avg, DateRange: &dateSpan{Start: start, End: end}}
		c.Payload(res)
	}
}
