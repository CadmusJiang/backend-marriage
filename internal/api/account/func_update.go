package account

import (
	"net/http"
	"time"

	"github.com/xinliangnote/go-gin-api/internal/code"
	"github.com/xinliangnote/go-gin-api/internal/pkg/core"
)

type updateRequest struct {
	AccountID   string `uri:"accountId" binding:"required"` // 账户ID
	Name        string `json:"name"`                        // 姓名
	Phone       string `json:"phone"`                       // 手机号
	BelongGroup string `json:"belongGroup"`                 // 所属组
	BelongTeam  string `json:"belongTeam"`                  // 所属团队
	Status      string `json:"status"`                      // 状态
}

type updateResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// UpdateAccount 更新账户
// @Summary 更新账户
// @Description 更新账户信息
// @Tags Account
// @Accept application/json
// @Produce json
// @Param accountId path string true "账户ID"
// @Param request body updateRequest true "更新账户请求"
// @Success 200 {object} updateResponse
// @Failure 400 {object} code.Failure
// @Failure 404 {object} code.Failure
// @Router /api/v1/accounts/{accountId} [put]
func (h *handler) UpdateAccount() core.HandlerFunc {
	return func(c core.Context) {
		req := new(updateRequest)
		res := new(updateResponse)

		if err := c.ShouldBindURI(req); err != nil {
			c.AbortWithError(core.Error(
				http.StatusBadRequest,
				code.ParamBindError,
				code.Text(code.ParamBindError)).WithError(err),
			)
			return
		}

		if err := c.ShouldBindJSON(req); err != nil {
			c.AbortWithError(core.Error(
				http.StatusBadRequest,
				code.ParamBindError,
				code.Text(code.ParamBindError)).WithError(err),
			)
			return
		}

		// Mock数据 - 检查账户是否存在
		mockAccounts := map[string]accountData{
			"0": {
				ID:          "0",
				Username:    "admin",
				Name:        "系统管理员",
				RoleType:    "company_manager",
				Status:      "enabled",
				Phone:       "13800138000",
				BelongGroup: &org{ID: 0, Username: "admin_group", Name: "系统管理组", CreatedAt: 1705123200, UpdatedAt: 1705123200, CurrentCnt: 1},
				BelongTeam:  &org{ID: 0, Username: "admin_team", Name: "系统管理团队", CreatedAt: 1705123200, UpdatedAt: 1705123200, CurrentCnt: 1},
				CreatedAt:   1705123200,
				UpdatedAt:   1705123200,
				LastLoginAt: time.Now().Unix(),
			},
			"1": {
				ID:          "1",
				Username:    "company_manager",
				Name:        "张伟",
				RoleType:    "company_manager",
				Status:      "enabled",
				Phone:       "13800138001",
				BelongGroup: &org{ID: 1, Username: "nanjing_tianyuan", Name: "南京-天元大厦组", CreatedAt: 1705123200, UpdatedAt: 1705123200, CurrentCnt: 15},
				BelongTeam:  &org{ID: 1, Username: "marketing_team_a", Name: "营销团队A", CreatedAt: 1705123200, UpdatedAt: 1705123200, CurrentCnt: 8},
				CreatedAt:   1705123200,
				UpdatedAt:   1705123200,
				LastLoginAt: time.Now().Unix(),
			},
		}

		_, exists := mockAccounts[req.AccountID]
		if !exists {
			c.AbortWithError(core.Error(
				http.StatusNotFound,
				code.AdminDetailError,
				"账户不存在"),
			)
			return
		}

		res.Success = true
		res.Message = "账户更新成功"

		c.Payload(res)
	}
}

// UpdateAccountPassword 仅更新密码
func (h *handler) UpdateAccountPassword() core.HandlerFunc {
	type response struct {
		Success bool `json:"success"`
	}
	return func(c core.Context) {
		var uri struct {
			AccountId string `uri:"accountId"`
		}
		if err := c.ShouldBindURI(&uri); err != nil {
			c.Payload(response{Success: false})
			return
		}
		var body struct {
			Password string `json:"password" binding:"required,min=6,max=64"`
		}
		if err := c.ShouldBindJSON(&body); err != nil {
			c.Payload(response{Success: false})
			return
		}
		if err := h.accountService.UpdatePassword(c, uri.AccountId, body.Password); err != nil {
			c.Payload(response{Success: false})
			return
		}
		c.Payload(response{Success: true})
	}
}
