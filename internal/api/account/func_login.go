package account

import (
	"fmt"
	"net/http"

	"github.com/xinliangnote/go-gin-api/internal/code"
	"github.com/xinliangnote/go-gin-api/internal/pkg/core"
	"go.uber.org/zap"
)

type loginRequest struct {
	Username string `json:"username" binding:"required"` // 用户名
	Password string `json:"password" binding:"required"` // 密码
}

type loginResponse struct {
	Status string `json:"status"`
	Data   struct {
		Token    string `json:"token"`
		Username string `json:"username"`
	} `json:"data"`
}

type logoutResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// Login 用户登录
// @Summary 用户登录
// @Description 用户登录接口
// @Tags API.account
// @Accept application/json
// @Produce json
// @Param request body loginRequest true "登录请求"
// @Success 200 {object} loginResponse
// @Failure 400 {object} code.Failure
// @Router /api/v1/auth/login [post]
func (h *handler) Login() core.HandlerFunc {
	return func(c core.Context) {
		req := new(loginRequest)
		res := new(loginResponse)

		if err := c.ShouldBindJSON(req); err != nil {
			h.logger.Error("参数绑定失败", zap.Error(err))
			c.AbortWithError(core.Error(
				http.StatusBadRequest,
				code.ParamBindError,
				code.Text(code.ParamBindError)).WithError(err),
			)
			return
		}

		// 添加调试日志
		h.logger.Info("登录请求",
			zap.String("username", req.Username),
			zap.String("password", req.Password),
		)

		// 调用服务层进行登录验证
		accountInfo, err := h.accountService.Login(c, req.Username, req.Password)
		if err != nil {
			h.logger.Error("登录失败",
				zap.String("username", req.Username),
				zap.Error(err),
			)
			c.AbortWithError(core.Error(
				http.StatusUnauthorized,
				code.AdminLoginError,
				"用户名或密码错误").WithError(err),
			)
			return
		}

		// 添加成功日志
		h.logger.Info("登录成功",
			zap.String("username", req.Username),
			zap.Int32("user_id", accountInfo.Id),
			zap.String("role_type", accountInfo.RoleType),
		)

		// 转换为API响应格式
		userInfo := accountData{
			ID:          fmt.Sprintf("%d", accountInfo.Id),
			Username:    accountInfo.Username,
			Name:        accountInfo.Nickname,
			Phone:       accountInfo.Phone,
			RoleType:    accountInfo.RoleType,
			Status:      accountInfo.Status,
			CreatedAt:   int64(accountInfo.CreatedTimestamp),
			UpdatedAt:   int64(accountInfo.ModifiedTimestamp),
			LastLoginAt: int64(accountInfo.LastLoginTimestamp),
		}

		// 暂时移除组织信息，因为新模型中没有相关字段
		// if accountInfo.BelongGroupId > 0 {
		// 	userInfo.BelongGroup = &orgInfo{
		// 		ID:                int(accountInfo.BelongGroupId),
		// 		Username:          accountInfo.BelongGroupUsername,
		// 		Name:              accountInfo.BelongGroupNickname,
		// 		CreatedAt:         accountInfo.BelongGroupCreatedTimestamp,
		// 		UpdatedAt:         accountInfo.BelongGroupModifiedTimestamp,
		// 		CurrentCnt:        int(accountInfo.BelongGroupCurrentCnt),
		// 	}
		// }

		// if accountInfo.BelongTeamId > 0 {
		// 	userInfo.BelongTeam = &orgInfo{
		// 		ID:                int(accountInfo.BelongTeamId),
		// 		Username:          accountInfo.BelongTeamUsername,
		// 		Name:              accountInfo.BelongTeamNickname,
		// 		CreatedAt:         accountInfo.BelongTeamCreatedTimestamp,
		// 		UpdatedAt:         accountInfo.BelongTeamModifiedTimestamp,
		// 		CurrentCnt:        int(accountInfo.BelongTeamCurrentCnt),
		// 	}
		// }

		// 生成token（实际应该使用JWT）
		token := userInfo.Username + "-token-123"

		res.Status = "success"
		res.Data.Token = token
		res.Data.Username = userInfo.Username

		c.Payload(res)
	}
}

// Logout 退出登录
// @Summary 退出登录
// @Description 用户退出登录
// @Tags API.account
// @Accept application/x-www-form-urlencoded
// @Produce json
// @Success 200 {object} logoutResponse
// @Failure 400 {object} code.Failure
// @Router /api/v1/auth/logout [post]
func (h *handler) Logout() core.HandlerFunc {
	return func(c core.Context) {
		res := new(logoutResponse)

		// 调用服务层退出登录
		err := h.accountService.Logout(c, "")
		if err != nil {
			h.logger.Error("退出登录失败", zap.Error(err))
		}

		res.Success = true
		res.Message = "退出登录成功"

		c.Payload(res)
	}
}
