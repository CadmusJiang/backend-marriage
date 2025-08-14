package cooperation_store

import (
	"net/http"
	"strconv"
	"time"

	"github.com/xinliangnote/go-gin-api/internal/code"
	"github.com/xinliangnote/go-gin-api/internal/pkg/core"
	mysqlModel "github.com/xinliangnote/go-gin-api/internal/repository/mysql/cooperation_store"
)

type historyListRequest struct {
	Current  int `form:"current"`
	PageSize int `form:"pageSize"`
}

type historyOperator struct {
	ID       string  `json:"id"`
	Username *string `json:"username"`
	Name     string  `json:"name"`
	Role     string  `json:"role"`
	RoleType string  `json:"roleType"`
}

type historyChange struct {
	Field string      `json:"field"`
	Old   interface{} `json:"old"`
	New   interface{} `json:"new"`
}

type historyItem struct {
	ID         string          `json:"id"`
	OccurredAt string          `json:"occurredAt"`
	Operator   historyOperator `json:"operator"`
	Action     string          `json:"action"`
	Changes    []historyChange `json:"changes"`
}

type historyListResponse struct {
	Success bool          `json:"success"`
	Data    []historyItem `json:"data"`
	Meta    struct {
		Total    int `json:"total"`
		PageSize int `json:"pageSize"`
		Current  int `json:"current"`
	} `json:"meta"`
}

// GetCooperationStoreHistory GET /v1/cooperation-stores/:id/history
// Mock 历史记录，用于前端对接联调。
func (h *handler) GetCooperationStoreHistory() core.HandlerFunc {
	// @Summary 获取门店历史
	// @Description 获取合作门店历史记录（Mock）
	// @Tags CooperationStore
	// @Accept application/x-www-form-urlencoded
	// @Produce json
	// @Param id path string true "门店ID"
	// @Param current query int false "当前页码" default(1)
	// @Param pageSize query int false "每页数量" default(10)
	// @Success 200 {object} historyListResponse
	// @Failure 400 {object} code.Failure
	// @Failure 404 {object} code.Failure
	// @Router /api/v1/cooperation-stores/{id}/history [get]
	return func(ctx core.Context) {
		idStr := ctx.Param("id")
		if idStr == "" {
			ctx.AbortWithError(core.Error(
				http.StatusBadRequest,
				code.ParamBindError,
				"id is required"))
			return
		}

		// 若存在则返回 mock 历史，不存在返回 404
		var exists mysqlModel.CooperationStore
		if err := h.db.GetDbR().Where("id = ?", idStr).Take(&exists).Error; err != nil {
			ctx.AbortWithError(core.Error(
				http.StatusNotFound,
				code.ServerError,
				"store not found").WithError(err))
			return
		}

		req := new(historyListRequest)
		if err := ctx.ShouldBindQuery(req); err != nil {
			ctx.AbortWithError(core.Error(
				http.StatusBadRequest,
				code.ParamBindError,
				code.Text(code.ParamBindError)).WithError(err))
			return
		}
		if req.Current <= 0 {
			req.Current = 1
		}
		if req.PageSize <= 0 {
			req.PageSize = 10
		}

		// 组装一些静态的 mock 历史数据
		now := time.Now()
		items := []historyItem{
			{
				ID:         "h-" + idStr + "-1",
				OccurredAt: now.Add(-72 * time.Hour).Format("2006-01-02T15:04:05.000Z"),
				Operator:   historyOperator{ID: "1", Username: strPtr("admin"), Name: "系统管理员", Role: "admin", RoleType: "system"},
				Action:     "created",
				Changes: []historyChange{
					{Field: "storeName", Old: nil, New: exists.StoreName},
					{Field: "cooperationCity", Old: nil, New: exists.CooperationCity},
				},
			},
			{
				ID:         "h-" + idStr + "-2",
				OccurredAt: now.Add(-24 * time.Hour).Format("2006-01-02T15:04:05.000Z"),
				Operator:   historyOperator{ID: "2", Username: strPtr("ops"), Name: "运营同学", Role: "operator", RoleType: "employee"},
				Action:     "updated",
				Changes: []historyChange{
					{Field: "cooperationStatus", Old: "pending", New: exists.CooperationStatus},
				},
			},
		}

		// 简单分页（内存分页）
		total := len(items)
		start := (req.Current - 1) * req.PageSize
		end := start + req.PageSize
		if start > total {
			start = total
		}
		if end > total {
			end = total
		}
		pageItems := items[start:end]

		// 规范化 ID 为字符串
		for i := range pageItems {
			// ensure ID present
			if pageItems[i].ID == "" {
				pageItems[i].ID = "h-" + idStr + "-" + strconv.Itoa(i+1)
			}
		}

		resp := historyListResponse{Success: true, Data: pageItems}
		resp.Meta.Total = total
		resp.Meta.PageSize = req.PageSize
		resp.Meta.Current = req.Current
		ctx.Payload(resp)
	}
}

func strPtr(s string) *string { return &s }
