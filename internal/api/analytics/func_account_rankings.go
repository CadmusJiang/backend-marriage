package analytics

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/xinliangnote/go-gin-api/internal/authz"
	"github.com/xinliangnote/go-gin-api/internal/code"
	"github.com/xinliangnote/go-gin-api/internal/pkg/core"
)

type accountRankingsRequest struct {
	Metric        string `form:"metric"`        // paid_count | lead_count
	TimeRange     string `form:"time"`          // today | week | last7days | month
	Top           int    `form:"top"`           // top N
	BelongGroupId string `form:"belongGroupId"` // group filter
	SortBy        string `form:"sortBy"`        // total | name
	SortOrder     string `form:"sortOrder"`     // asc | desc
}

type rankingItem struct {
	Name       string  `json:"name"`
	Value      float64 `json:"value"`
	Rank       int     `json:"rank,omitempty"`
	RankChange int     `json:"rankChange,omitempty"`
	Avatar     string  `json:"avatar,omitempty"`
}

type analyticsMeta struct {
	Metric        string    `json:"metric,omitempty"`
	Time          string    `json:"time,omitempty"`
	BelongGroupId string    `json:"belongGroupId,omitempty"`
	TotalUsers    int       `json:"totalUsers,omitempty"`
	TotalTeams    int       `json:"totalTeams,omitempty"`
	AverageValue  float64   `json:"averageValue,omitempty"`
	DateRange     *dateSpan `json:"dateRange,omitempty"`
}

type dateSpan struct {
	Start time.Time `json:"start"`
	End   time.Time `json:"end"`
}

type accountRankingsResponse struct {
	Success bool `json:"success"`
	Data    struct {
		Rankings []rankingItem `json:"rankings"`
		Meta     analyticsMeta `json:"meta"`
	} `json:"data"`
}

// GetAccountRankings 获取个人排名数据
// @Summary 获取个人排名数据
// @Description 按开单或留资获取个人排名
// @Tags Analytics
// @Accept application/x-www-form-urlencoded
// @Produce json
// @Param metric query string true "指标类型" Enums(paid_count, lead_count)
// @Param time query string false "时间范围" Enums(today, week, last7days, month) default(month)
// @Param top query int false "返回前N名"
// @Param belongGroupId query string false "归属组ID"
// @Param sortBy query string false "排序字段" Enums(total,name)
// @Param sortOrder query string false "排序方向" Enums(asc, desc) default(desc)
// @Success 200 {object} accountRankingsResponse
// @Failure 400 {object} code.Failure
// @Router /api/v1/analytics/account/rankings [get]
func (h *handler) GetAccountRankings() core.HandlerFunc {
	return func(c core.Context) {
		req := new(accountRankingsRequest)
		if err := c.ShouldBindForm(req); err != nil {
			c.AbortWithError(core.Error(http.StatusBadRequest, code.ParamBindError, code.Text(code.ParamBindError)).WithError(err))
			return
		}

		// defaults and validation
		if req.Metric == "" {
			c.AbortWithError(core.Error(http.StatusBadRequest, code.ParamBindError, "metric is required"))
			return
		}
		req.Metric = strings.ToLower(req.Metric)
		if req.TimeRange == "" {
			req.TimeRange = "month"
		}
		if req.SortOrder == "" {
			req.SortOrder = "desc"
		}

		// 范围上限：非公司管理员一律限定在“本组”
		scope, _ := authz.ComputeScope(c, h.db)
		if !scope.ScopeAll {
			if len(scope.AllowedGroupIDs) > 0 {
				// 强制覆盖到本组（防止越权指定其他组）
				req.BelongGroupId = formatFirst(scope.AllowedGroupIDs)
			} else {
				// 没有组范围的（比如 team/employee），仍然按其归属组做聚合，前端 Meta 中体现
				req.BelongGroupId = ""
			}
		}

		// TODO: replace with real aggregation queries. For now return stub consistent with schema
		rankings := []rankingItem{
			{Name: "张三", Value: 120, Rank: 1, RankChange: 2, Avatar: "https://example.com/a.png"},
			{Name: "李四", Value: 110, Rank: 2, RankChange: -1, Avatar: "https://example.com/b.png"},
			{Name: "王五", Value: 90, Rank: 3, RankChange: 0, Avatar: "https://example.com/c.png"},
		}

		// top N
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

		res := new(accountRankingsResponse)
		res.Success = true
		res.Data.Rankings = rankings
		res.Data.Meta = analyticsMeta{
			Metric:        req.Metric,
			Time:          req.TimeRange,
			BelongGroupId: req.BelongGroupId,
			AverageValue:  avg,
			DateRange:     &dateSpan{Start: start, End: end},
		}
		c.Payload(res)
	}
}

func formatFirst(vals []int32) string {
	if len(vals) == 0 {
		return ""
	}
	return fmt.Sprintf("%d", vals[0])
}

func resolveTimeRange(r string) (time.Time, time.Time) {
	now := time.Now().UTC()
	switch r {
	case "today":
		start := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
		end := start.Add(24 * time.Hour)
		return start, end
	case "week":
		// ISO week: start Monday
		weekday := int(now.Weekday())
		if weekday == 0 { // Sunday
			weekday = 7
		}
		start := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC).AddDate(0, 0, -(weekday - 1))
		end := start.AddDate(0, 0, 7)
		return start, end
	case "last7days":
		end := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, time.UTC)
		start := end.AddDate(0, 0, -7)
		return start, end
	case "month":
		start := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
		end := start.AddDate(0, 1, 0)
		return start, end
	default:
		start := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
		end := start.AddDate(0, 1, 0)
		return start, end
	}
}
