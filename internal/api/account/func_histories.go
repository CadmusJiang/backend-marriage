package account

import (
	"net/http"

	"github.com/xinliangnote/go-gin-api/internal/authz"
	"github.com/xinliangnote/go-gin-api/internal/code"
	"github.com/xinliangnote/go-gin-api/internal/pkg/core"
)

type historiesRequest struct {
	AccountID       string `form:"accountId"`       // 账户ID
	AccountUsername string `form:"accountUsername"` // 账户用户名
	Current         int    `form:"current"`         // 当前页码
	PageSize        int    `form:"pageSize"`        // 每页数量
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
// @Param accountId query string false "账户ID"
// @Param accountUsername query string false "账户用户名"
// @Param current query int false "当前页码" default(1)
// @Param pageSize query int false "每页数量" default(10)
// @Success 200 {object} historiesResponse
// @Failure 400 {object} code.Failure
// @Router /api/v1/account-histories [get]
func (h *handler) GetAccountHistories() core.HandlerFunc {
	return func(c core.Context) {
		req := new(historiesRequest)
		res := new(historiesResponse)

		if err := c.ShouldBindForm(req); err != nil {
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

		// Mock数据 - 与前端数据结构保持一致
		mockHistories := []accountHistoryData{
			// accountId=1 的历史记录
			{
				ID:          "101",
				AccountID:   "1",
				OperateType: "created",
				OperatedAt:  1705123200, // 2024-01-13 10:00:00
				Content: map[string]interface{}{
					"username": map[string]string{
						"old": "",
						"new": "admin",
					},
					"name": map[string]string{
						"old": "",
						"new": "系统管理员",
					},
					"roleType": map[string]string{
						"old": "",
						"new": "company_manager",
					},
					"phone": map[string]string{
						"old": "",
						"new": "13800138000",
					},
					"status": map[string]string{
						"old": "",
						"new": "enabled",
					},
				},
				OperatorUsername: "system",
				OperatorName:     "系统",
				OperatorRoleType: "system",
			},
			{
				ID:          "102",
				AccountID:   "1",
				OperateType: "modified",
				OperatedAt:  1705209600, // 2024-01-14 10:00:00
				Content: map[string]interface{}{
					"phone": map[string]string{
						"old": "13800138000",
						"new": "13800138001",
					},
				},
				OperatorUsername: "admin",
				OperatorName:     "系统管理员",
				OperatorRoleType: "company_manager",
			},
			{
				ID:          "103",
				AccountID:   "1",
				OperateType: "modified",
				OperatedAt:  1705296000, // 2024-01-15 10:00:00
				Content: map[string]interface{}{
					"name": map[string]string{
						"old": "系统管理员",
						"new": "超级管理员",
					},
				},
				OperatorUsername: "admin",
				OperatorName:     "系统管理员",
				OperatorRoleType: "company_manager",
			},
			{
				ID:          "104",
				AccountID:   "1",
				OperateType: "modified",
				OperatedAt:  1705382400, // 2024-01-16 10:00:00
				Content: map[string]interface{}{
					"status": map[string]string{
						"old": "enabled",
						"new": "disabled",
					},
					"reason": map[string]string{
						"old": "",
						"new": "临时禁用账户",
					},
				},
				OperatorUsername: "zhangwei",
				OperatorName:     "张伟",
				OperatorRoleType: "company_manager",
			},
			{
				ID:          "105",
				AccountID:   "1",
				OperateType: "modified",
				OperatedAt:  1705468800, // 2024-01-17 10:00:00
				Content: map[string]interface{}{
					"status": map[string]string{
						"old": "disabled",
						"new": "enabled",
					},
					"reason": map[string]string{
						"old": "临时禁用账户",
						"new": "问题已解决，恢复账户",
					},
				},
				OperatorUsername: "zhangwei",
				OperatorName:     "张伟",
				OperatorRoleType: "company_manager",
			},
			{
				ID:          "106",
				AccountID:   "1",
				OperateType: "modified",
				OperatedAt:  1705555200, // 2024-01-18 10:00:00
				Content: map[string]interface{}{
					"password": map[string]string{
						"old": "****",
						"new": "******",
					},
				},
				OperatorUsername: "admin",
				OperatorName:     "超级管理员",
				OperatorRoleType: "company_manager",
			},
			{
				ID:          "107",
				AccountID:   "1",
				OperateType: "modified",
				OperatedAt:  1705641600, // 2024-01-19 10:00:00
				Content: map[string]interface{}{
					"name": map[string]string{
						"old": "超级管理员",
						"new": "系统管理员",
					},
				},
				OperatorUsername: "admin",
				OperatorName:     "超级管理员",
				OperatorRoleType: "company_manager",
			},
			{
				ID:          "108",
				AccountID:   "1",
				OperateType: "modified",
				OperatedAt:  1705728000, // 2024-01-20 10:00:00
				Content: map[string]interface{}{
					"phone": map[string]string{
						"old": "13800138001",
						"new": "13800138000",
					},
				},
				OperatorUsername: "admin",
				OperatorName:     "系统管理员",
				OperatorRoleType: "company_manager",
			},
			// accountId=6 的原有历史记录
			{
				ID:          "1",
				AccountID:   "6",
				OperateType: "modified",
				OperatedAt:  1705923000, // 2024-01-22 14:30:00
				Content: map[string]interface{}{
					"roleType": map[string]string{
						"old": "员工",
						"new": "小队管理员",
					},
					"belongTeam": map[string]string{
						"old": "无",
						"new": "营销团队A",
					},
					"status": map[string]string{
						"old": "enabled",
						"new": "disabled",
					},
				},
				OperatorUsername: "zhangwei",
				OperatorName:     "张伟",
				OperatorRoleType: "company_manager",
			},
			{
				ID:          "2",
				AccountID:   "6",
				OperateType: "modified",
				OperatedAt:  1705754700, // 2024-01-20 16:45:00
				Content: map[string]interface{}{
					"belongGroup": map[string]string{
						"old": "南京-天元大厦组",
						"new": "南京-南京南站组",
					},
					"belongTeam": map[string]string{
						"old": "营销团队A",
						"new": "营销团队C",
					},
				},
				OperatorUsername: "liming",
				OperatorName:     "李明",
				OperatorRoleType: "team_manager",
			},
			{
				ID:          "3",
				AccountID:   "6",
				OperateType: "modified",
				OperatedAt:  1705565700, // 2024-01-18 09:15:00
				Content: map[string]interface{}{
					"status": map[string]string{
						"old": "enabled",
						"new": "disabled",
					},
				},
				OperatorUsername: "wangfang",
				OperatorName:     "王芳",
				OperatorRoleType: "team_manager",
			},
			{
				ID:          "4",
				AccountID:   "6",
				OperateType: "modified",
				OperatedAt:  1705303200, // 2024-01-15 11:20:00
				Content: map[string]interface{}{
					"phone": map[string]string{
						"old": "13800138000",
						"new": "13900139000",
					},
					"belongTeam": map[string]string{
						"old": "营销团队A",
						"new": "营销团队B",
					},
				},
				OperatorUsername: "liuqiang",
				OperatorName:     "刘强",
				OperatorRoleType: "team_manager",
			},
			{
				ID:          "5",
				AccountID:   "6",
				OperateType: "modified",
				OperatedAt:  1705032000, // 2024-01-12 10:00:00
				Content: map[string]interface{}{
					"name": map[string]string{
						"old": "张三",
						"new": "张明",
					},
					"roleType": map[string]string{
						"old": "员工",
						"new": "小队管理员",
					},
				},
				OperatorUsername: "chenjing",
				OperatorName:     "陈静",
				OperatorRoleType: "team_manager",
			},
			{
				ID:          "6",
				AccountID:   "6",
				OperateType: "created",
				OperatedAt:  1704877800, // 2024-01-10 14:30:00
				Content: map[string]interface{}{
					"username": map[string]string{
						"old": "",
						"new": "employee001",
					},
					"name": map[string]string{
						"old": "",
						"new": "张三",
					},
					"roleType": map[string]string{
						"old": "",
						"new": "员工",
					},
					"phone": map[string]string{
						"old": "",
						"new": "13800138000",
					},
					"belongGroup": map[string]string{
						"old": "",
						"new": "南京-天元大厦组",
					},
					"belongTeam": map[string]string{
						"old": "",
						"new": "营销团队A",
					},
				},
				OperatorUsername: "admin",
				OperatorName:     "系统管理员",
				OperatorRoleType: "company_manager",
			},
			{
				ID:          "7",
				AccountID:   "6",
				OperateType: "modified",
				OperatedAt:  1704709200, // 2024-01-08 16:20:00
				Content: map[string]interface{}{
					"belongGroup": map[string]string{
						"old": "南京-天元大厦组",
						"new": "南京-夫子庙组",
					},
					"belongTeam": map[string]string{
						"old": "营销团队A",
						"new": "营销团队D",
					},
				},
				OperatorUsername: "zhaomin",
				OperatorName:     "赵敏",
				OperatorRoleType: "team_manager",
			},
			{
				ID:          "8",
				AccountID:   "6",
				OperateType: "modified",
				OperatedAt:  1704457500, // 2024-01-05 13:45:00
				Content: map[string]interface{}{
					"status": map[string]string{
						"old": "disabled",
						"new": "enabled",
					},
					"reason": map[string]string{
						"old": "",
						"new": "问题已解决，恢复账户",
					},
				},
				OperatorUsername: "liming",
				OperatorName:     "李明",
				OperatorRoleType: "team_manager",
			},
			{
				ID:          "9",
				AccountID:   "6",
				OperateType: "modified",
				OperatedAt:  1704265800, // 2024-01-03 10:30:00
				Content: map[string]interface{}{
					"password": map[string]string{
						"old": "****",
						"new": "******",
					},
				},
				OperatorUsername: "zhangwei",
				OperatorName:     "张伟",
				OperatorRoleType: "company_manager",
			},
			{
				ID:          "10",
				AccountID:   "6",
				OperateType: "modified",
				OperatedAt:  1704081600, // 2024-01-01 09:00:00
				Content: map[string]interface{}{
					"password": map[string]string{
						"old": "****",
						"new": "******",
					},
				},
				OperatorUsername: "wangfang",
				OperatorName:     "王芳",
				OperatorRoleType: "team_manager",
			},
		}

		// 根据accountId过滤历史记录
		filteredHistories := []accountHistoryData{}
		for _, hist := range mockHistories {
			if hist.AccountID == req.AccountID {
				filteredHistories = append(filteredHistories, hist)
			}
		}

		// 根据accountUsername过滤操作人
		if req.AccountUsername != "" {
			var filtered []accountHistoryData
			for _, hist := range filteredHistories {
				if hist.OperatorName == req.AccountUsername {
					filtered = append(filtered, hist)
				}
			}
			filteredHistories = filtered
		}

		// 分页逻辑
		total := len(filteredHistories)
		start := (req.Current - 1) * req.PageSize
		end := start + req.PageSize

		if start >= total {
			res.Data = []accountHistoryData{}
		} else if end > total {
			res.Data = filteredHistories[start:total]
		} else {
			res.Data = filteredHistories[start:end]
		}

		res.Total = total
		res.Success = true
		res.PageSize = req.PageSize
		res.Current = req.Current

		c.Payload(res)
	}
}
