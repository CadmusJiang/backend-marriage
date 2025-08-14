package analytics

import (
	"net/http"

	"github.com/xinliangnote/go-gin-api/internal/code"
	"github.com/xinliangnote/go-gin-api/internal/pkg/core"
)

type customerCompositionRequest struct {
	TimeRange     string `form:"time"`
	BelongGroupId string `form:"belongGroupId"`
}

type distributionItem struct {
	Type  string  `json:"type"`
	Value float64 `json:"value"`
}

type customerComposition struct {
	Gender    []distributionItem `json:"gender"`
	Store     []distributionItem `json:"store"`
	Income    []distributionItem `json:"income"`
	BirthYear []distributionItem `json:"birthYear"`
	Education []distributionItem `json:"education"`
	City      []distributionItem `json:"city"`
}

type customerCompositionResponse struct {
	Success bool `json:"success"`
	Data    struct {
		Compositions customerComposition `json:"compositions"`
		Meta         struct {
			BelongGroupId string  `json:"belongGroupId"`
			Total         float64 `json:"total"`
		} `json:"meta"`
	} `json:"data"`
}

// GetCustomerComposition 获取客资授权记录构成（包含城市分布）
// @Summary 获取客资授权记录构成（包含城市分布）
// @Description 返回按多个维度的客资构成占比
// @Tags Analytics
// @Accept application/x-www-form-urlencoded
// @Produce json
// @Param time query string false "统计时间范围" Enums(today, week, last7days, month) default(month)
// @Param belongGroupId query string false "归属组ID"
// @Success 200 {object} customerCompositionResponse
// @Failure 400 {object} code.Failure
// @Router /api/v1/analytics/customer-authorization-record/composition [get]
func (h *handler) GetCustomerComposition() core.HandlerFunc {
	return func(c core.Context) {
		req := new(customerCompositionRequest)
		if err := c.ShouldBindForm(req); err != nil {
			c.AbortWithError(core.Error(http.StatusBadRequest, code.ParamBindError, code.Text(code.ParamBindError)).WithError(err))
			return
		}
		if req.TimeRange == "" {
			req.TimeRange = "month"
		}

		compositions := customerComposition{
			Gender:    []distributionItem{{Type: "male", Value: 55}, {Type: "female", Value: 45}},
			Store:     []distributionItem{{Type: "门店A", Value: 30}, {Type: "门店B", Value: 70}},
			Income:    []distributionItem{{Type: "<30w", Value: 20}, {Type: "30-60w", Value: 50}, {Type: ">60w", Value: 30}},
			BirthYear: []distributionItem{{Type: "1980s", Value: 25}, {Type: "1990s", Value: 50}, {Type: "2000s", Value: 25}},
			Education: []distributionItem{{Type: "本科", Value: 60}, {Type: "硕士", Value: 30}, {Type: "其他", Value: 10}},
			City:      []distributionItem{{Type: "北京", Value: 40}, {Type: "上海", Value: 35}, {Type: "深圳", Value: 25}},
		}

		res := new(customerCompositionResponse)
		res.Success = true
		res.Data.Compositions = compositions
		res.Data.Meta.BelongGroupId = req.BelongGroupId
		res.Data.Meta.Total = 100
		c.Payload(res)
	}
}
