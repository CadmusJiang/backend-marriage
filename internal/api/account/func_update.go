package account

import (
	"fmt"
	"net/http"

	"github.com/xinliangnote/go-gin-api/internal/code"
	"github.com/xinliangnote/go-gin-api/internal/pkg/core"
	"github.com/xinliangnote/go-gin-api/internal/services/account"
	"go.uber.org/zap"
)

type updateRequest struct {
	AccountID string `uri:"accountId" binding:"required"` // 账户ID
	Name      string `json:"name"`                        // 姓名
	Phone     string `json:"phone"`                       // 手机号
	Status    string `json:"status"`                      // 状态
}

type updateResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    *accountData `json:"data,omitempty"` // 更新后的账户信息
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

		// 调用服务层更新账户
		err := h.accountService.Update(c, req.AccountID, &account.UpdateAccountData{
			Name:   req.Name,
			Phone:  req.Phone,
			Status: req.Status,
		})

		if err != nil {
			h.logger.Error("更新账户失败", zap.Error(err))
			c.AbortWithError(core.Error(
				http.StatusInternalServerError,
				code.ServerError,
				fmt.Sprintf("更新账户失败: %v", err)).WithError(err),
			)
			return
		}

		// 更新成功后，查询并返回更新后的账户信息
		updatedAccount, err := h.accountService.Detail(c, req.AccountID)
		if err != nil {
			h.logger.Error("查询更新后的账户信息失败", zap.Error(err))
			// 即使查询失败，更新操作本身是成功的，所以仍然返回成功
			res.Success = true
			res.Message = "账户更新成功，但获取更新后信息失败"
			c.Payload(res)
			return
		}

		// 构建返回的账户数据
		accountData := accountData{
			ID:        fmt.Sprintf("%d", updatedAccount.Id),
			Username:  updatedAccount.Username,
			Name:      updatedAccount.Name,
			Phone:     updatedAccount.Phone,
			RoleType:  updatedAccount.RoleType,
			Status:    updatedAccount.Status,
			CreatedAt: updatedAccount.CreatedAt.Unix(),
			UpdatedAt: updatedAccount.UpdatedAt.Unix(),
			LastLoginAt: func() int64 {
				if updatedAccount.LastLoginAt != nil {
					return updatedAccount.LastLoginAt.Unix()
				}
				return 0
			}(),
		}

		// 根据includeGroup和includeTeam参数决定是否包含组织信息
		// 这里默认包含组织信息，因为更新后用户通常需要看到完整信息
		belongGroup, belongTeam := h.getAccountOrgInfo(int(updatedAccount.Id))
		accountData.BelongGroup = belongGroup
		accountData.BelongTeam = belongTeam

		res.Success = true
		res.Message = "账户更新成功"
		res.Data = &accountData

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
