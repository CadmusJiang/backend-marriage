package account

import (
	"net/http"
	"time"

	"github.com/xinliangnote/go-gin-api/internal/code"
	"github.com/xinliangnote/go-gin-api/internal/pkg/core"
)

type createRequest struct {
	Username    string `json:"username" binding:"required"` // 用户名
	Name        string `json:"name" binding:"required"`     // 姓名
	Phone       string `json:"phone"`                       // 手机号
	Password    string `json:"password" binding:"required"` // 密码
	BelongGroup string `json:"belongGroup"`                 // 所属组
	RoleType    string `json:"roleType"`                    // 角色类型
}

type createResponse struct {
	Data    accountData `json:"data"`
	Success bool        `json:"success"`
	Message string      `json:"message"`
}

// CreateAccount 创建账户
// @Summary 创建账户
// @Description 创建新账户
// @Tags Account
// @Accept application/json
// @Produce json
// @Param request body createRequest true "创建账户请求"
// @Success 200 {object} createResponse
// @Failure 400 {object} code.Failure
// @Router /api/v1/accounts [post]
func (h *handler) CreateAccount() core.HandlerFunc {
	return func(c core.Context) {
		req := new(createRequest)
		res := new(createResponse)

		if err := c.ShouldBindJSON(req); err != nil {
			c.AbortWithError(core.Error(
				http.StatusBadRequest,
				code.ParamBindError,
				code.Text(code.ParamBindError)).WithError(err),
			)
			return
		}

		// 生成新的账户ID
		newID := "acc_" + time.Now().Format("20060102150405")
		currentTimestamp := time.Now().Unix()

		// 创建新账户数据
		newAccount := accountData{
			ID:          newID,
			Username:    req.Username,
			Name:        req.Name,
			RoleType:    req.RoleType,
			Status:      "enabled",
			Phone:       req.Phone,
			BelongGroup: &org{ID: 0, Username: "default_group", Name: "默认组", CreatedAt: currentTimestamp, UpdatedAt: currentTimestamp, CurrentCnt: 0},
			BelongTeam:  &org{ID: 0, Username: "default_team", Name: "默认团队", CreatedAt: currentTimestamp, UpdatedAt: currentTimestamp, CurrentCnt: 0},
			CreatedAt:   currentTimestamp,
			UpdatedAt:   currentTimestamp,
			LastLoginAt: 0,
		}

		res.Data = newAccount
		res.Success = true
		res.Message = "账户创建成功"

		c.Payload(res)
	}
}
