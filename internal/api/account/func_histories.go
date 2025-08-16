package account

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"go.uber.org/zap"

	"github.com/xinliangnote/go-gin-api/internal/authz"
	"github.com/xinliangnote/go-gin-api/internal/code"
	"github.com/xinliangnote/go-gin-api/internal/pkg/core"
)

type historiesRequest struct {
	AccountID string `form:"accountId"` // 账户ID（从路径参数获取）
	Current   int    `form:"current"`   // 当前页码
	PageSize  int    `form:"pageSize"`  // 每页数量
}

type accountHistoryData struct {
	ID               string                 `json:"id"`                // 历史记录ID
	AccountID        string                 `json:"accountId"`         // 账户ID
	OperateType      string                 `json:"operateType"`       // 操作类型
	OperatedAt       int64                  `json:"operatedAt"`        // 操作时间戳
	Content          map[string]interface{} `json:"content"`           // 操作内容
	OperatorUsername string                 `json:"operator_username"` // 操作人用户名
	OperatorName     string                 `json:"operator_name"`     // 操作人姓名
	OperatorRoleType string                 `json:"operatorRoleType"`  // 操作人角色类型
}

type historiesResponse struct {
	Data     []accountHistoryData `json:"data"`
	Total    int                  `json:"total"`
	Success  bool                 `json:"success"`
	PageSize int                  `json:"pageSize"`
	Current  int                  `json:"current"`
}

// GetAccountHistories 获取账户历史记录
// @Summary 获取账户历史记录
// @Description 分页获取账户操作历史记录
// @Tags Account
// @Accept application/x-www-form-urlencoded
// @Produce json
// @Param accountId path string true "账户ID"
// @Param current query int false "当前页码" default(1)
// @Param pageSize query int false "每页数量" default(10)
// @Success 200 {object} historiesResponse
// @Failure 400 {object} code.Failure
// @Router /api/v1/accounts/{accountId}/history [get]
func (h *handler) GetAccountHistories() core.HandlerFunc {
	return func(c core.Context) {
		req := new(historiesRequest)
		res := new(historiesResponse)

		// 从URL路径参数中获取accountId
		req.AccountID = c.Param("accountId")

		// 从查询参数中获取分页参数
		if err := c.ShouldBindQuery(req); err != nil {
			c.AbortWithError(core.Error(
				http.StatusBadRequest,
				code.ParamBindError,
				code.Text(code.ParamBindError)).WithError(err),
			)
			return
		}

		// 设置默认值
		if req.Current == 0 {
			req.Current = 1
		}
		if req.PageSize == 0 {
			req.PageSize = 10
		}

		// 检查必要参数
		if req.AccountID == "" {
			c.AbortWithError(core.Error(
				http.StatusBadRequest,
				code.ParamBindError,
				"缺少必要参数: accountId"),
			)
			return
		}

		// 范围校验：仅允许在可见范围内查看该账户的历史
		// 使用 service.Detail 获取账户信息，再进行范围判定
		if h.accountService != nil {
			info, _ := h.accountService.Detail(c, req.AccountID)
			if info != nil {
				scope, _ := authz.ComputeScope(c, h.db)
				if !authz.CanAccessAccount(scope, info) {
					c.AbortWithError(core.Error(
						http.StatusForbidden,
						code.RBACError,
						code.Text(code.RBACError)),
					)
					return
				}
			}
		}

		// 从数据库查询历史记录
		var histories []struct {
			ID               uint32    `json:"id" gorm:"column:id"`
			AccountID        uint32    `json:"account_id" gorm:"column:account_id"`
			OperateType      string    `json:"operate_type" gorm:"column:operate_type"`
			OperatedAt       time.Time `json:"operated_at" gorm:"column:operated_at"`
			Content          []byte    `json:"content" gorm:"column:content;type:json"`
			OperatorUsername string    `json:"operator_username" gorm:"column:operator_username"`
			OperatorName     string    `json:"operator_name" gorm:"column:operator_name"`
			OperatorRoleType string    `json:"operator_role_type" gorm:"column:operator_role_type"`
		}

		// 查询指定账户的历史记录
		if err := h.db.GetDbR().Table("account_history").
			Where("account_id = ?", req.AccountID).
			Order("operated_at DESC").
			Find(&histories).Error; err != nil {
			c.AbortWithError(core.Error(
				http.StatusInternalServerError,
				code.ServerError,
				"查询历史记录失败").WithError(err),
			)
			return
		}

		// 转换为响应格式
		var historyDataList []accountHistoryData
		for _, hist := range histories {
			// 解析JSON内容
			var content map[string]interface{}
			if len(hist.Content) > 0 {
				if err := json.Unmarshal(hist.Content, &content); err != nil {
					// 如果JSON解析失败，记录错误但继续处理
					zap.L().Warn("JSON解析失败",
						zap.Uint32("historyID", hist.ID),
						zap.Error(err),
						zap.String("rawContent", string(hist.Content)))
					content = nil
				}
			}

			historyDataList = append(historyDataList, accountHistoryData{
				ID:               strconv.FormatUint(uint64(hist.ID), 10),
				AccountID:        strconv.FormatUint(uint64(hist.AccountID), 10),
				OperateType:      hist.OperateType,
				OperatedAt:       hist.OperatedAt.Unix(),
				Content:          content,
				OperatorUsername: hist.OperatorUsername,
				OperatorName:     hist.OperatorName,
				OperatorRoleType: hist.OperatorRoleType,
			})
		}

		// 分页逻辑
		total := len(historyDataList)
		start := (req.Current - 1) * req.PageSize
		end := start + req.PageSize

		if start >= total {
			res.Data = []accountHistoryData{}
		} else if end > total {
			res.Data = historyDataList[start:total]
		} else {
			res.Data = historyDataList[start:end]
		}

		res.Total = total
		res.Success = true
		res.PageSize = req.PageSize
		res.Current = req.Current

		c.Payload(res)
	}
}
