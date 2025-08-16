package account

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/xinliangnote/go-gin-api/internal/code"
	"github.com/xinliangnote/go-gin-api/internal/pkg/core"
	"github.com/xinliangnote/go-gin-api/internal/services/account"
)

type createRequest struct {
	Username      string `json:"username" binding:"required"`      // 用户名
	Name          string `json:"name" binding:"required"`          // 姓名
	Phone         string `json:"phone"`                            // 手机号
	Password      string `json:"password" binding:"required"`      // 密码
	BelongGroupId string `json:"belongGroupId" binding:"required"` // 所属组ID（必填）
}

type createResponse struct {
	Data    accountData `json:"data"`
	Success bool        `json:"success"`
	Message string      `json:"message"`
}

// CreateAccount 创建账户
// @Summary 创建账户
// @Description 创建新账户（归属组必填，不允许创建公司管理员，默认角色为普通员工，不允许填写归属小队）
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
			// 提供详细的参数绑定错误信息
			var detailedMsg string
			switch {
			case strings.Contains(err.Error(), "belongGroupId"):
				detailedMsg = "归属组ID格式错误：必须是数字，不能是字符串"
			case strings.Contains(err.Error(), "username"):
				detailedMsg = "用户名格式错误：必须是字符串类型"
			case strings.Contains(err.Error(), "name"):
				detailedMsg = "姓名格式错误：必须是字符串类型"
			case strings.Contains(err.Error(), "password"):
				detailedMsg = "密码格式错误：必须是字符串类型"
			case strings.Contains(err.Error(), "phone"):
				detailedMsg = "手机号格式错误：必须是字符串类型"
			default:
				detailedMsg = fmt.Sprintf("参数绑定失败：%v", err)
			}

			c.AbortWithError(core.Error(
				http.StatusBadRequest,
				code.ParamBindError,
				detailedMsg).WithError(err),
			)
			return
		}

		// 验证必填字段
		if req.Username == "" {
			c.AbortWithError(core.Error(
				http.StatusBadRequest,
				code.ParamBindError,
				"用户名不能为空，请输入用户名"),
			)
			return
		}

		// 验证用户名格式：只能包含字母、数字、下划线，长度3-20
		if !regexp.MustCompile(`^[a-zA-Z0-9_]{3,20}$`).MatchString(req.Username) {
			c.AbortWithError(core.Error(
				http.StatusBadRequest,
				code.ParamBindError,
				fmt.Sprintf("用户名格式错误：'%s'，只能包含字母、数字、下划线，长度3-20位", req.Username)),
			)
			return
		}

		if req.Name == "" {
			c.AbortWithError(core.Error(
				http.StatusBadRequest,
				code.ParamBindError,
				"姓名不能为空，请输入姓名"),
			)
			return
		}

		// 验证姓名格式：中文、英文、数字，长度2-20
		if !regexp.MustCompile(`^[\p{Han}a-zA-Z0-9]{2,20}$`).MatchString(req.Name) {
			c.AbortWithError(core.Error(
				http.StatusBadRequest,
				code.ParamBindError,
				fmt.Sprintf("姓名格式错误：'%s'，只能包含中文、英文、数字，长度2-20位", req.Name)),
			)
			return
		}

		if req.Password == "" {
			c.AbortWithError(core.Error(
				http.StatusBadRequest,
				code.ParamBindError,
				"密码不能为空，请输入密码"),
			)
			return
		}

		// 验证密码格式：至少8位，包含字母和数字
		if len(req.Password) < 8 {
			c.AbortWithError(core.Error(
				http.StatusBadRequest,
				code.ParamBindError,
				"密码长度不足，至少需要8位"),
			)
			return
		}

		// 检查是否包含字母
		hasLetter := regexp.MustCompile(`[A-Za-z]`).MatchString(req.Password)
		if !hasLetter {
			c.AbortWithError(core.Error(
				http.StatusBadRequest,
				code.ParamBindError,
				"密码必须包含至少一个字母"),
			)
			return
		}

		// 检查是否包含数字
		hasDigit := regexp.MustCompile(`\d`).MatchString(req.Password)
		if !hasDigit {
			c.AbortWithError(core.Error(
				http.StatusBadRequest,
				code.ParamBindError,
				"密码必须包含至少一个数字"),
			)
			return
		}

		// 检查是否只包含允许的字符
		if !regexp.MustCompile(`^[A-Za-z\d@$!%*?&]+$`).MatchString(req.Password) {
			c.AbortWithError(core.Error(
				http.StatusBadRequest,
				code.ParamBindError,
				"密码只能包含字母、数字和特殊字符 @$!%*?&"),
			)
			return
		}

		if req.BelongGroupId == "" {
			c.AbortWithError(core.Error(
				http.StatusBadRequest,
				code.ParamBindError,
				"归属组ID不能为空，请输入归属组ID"),
			)
			return
		}

		// 验证手机号格式（如果提供）
		if req.Phone != "" {
			if !regexp.MustCompile(`^1[3-9]\d{9}$`).MatchString(req.Phone) {
				c.AbortWithError(core.Error(
					http.StatusBadRequest,
					code.ParamBindError,
					fmt.Sprintf("手机号格式错误：'%s'，请输入11位有效手机号，以1开头", req.Phone)),
				)
				return
			}
		}

		// 构造创建账户数据
		createData := &account.CreateAccountData{
			Username: req.Username,
			Name:     req.Name,
			Phone:    req.Phone,
			Password: req.Password,
			RoleType: "employee", // 默认角色为普通员工
			Status:   "enabled",
			BelongGroup: &account.OrgInfo{
				ID:   req.BelongGroupId,
				Name: fmt.Sprintf("组%s", req.BelongGroupId), // 这里可以根据ID查询实际的组名
			},
			BelongTeam: nil, // 不允许填写归属小队
		}

		// 调用服务层创建账户
		accountID, err := h.accountService.Create(c, createData)
		if err != nil {
			c.AbortWithError(core.Error(
				http.StatusInternalServerError,
				code.ServerError,
				fmt.Sprintf("创建账户失败: %v", err)),
			)
			return
		}

		// 构造响应数据
		newAccount := accountData{
			ID:          fmt.Sprintf("%d", accountID),
			Username:    req.Username,
			Name:        req.Name,
			RoleType:    "employee", // 默认角色为普通员工
			Status:      "enabled",
			Phone:       req.Phone,
			BelongGroup: &org{ID: uint64(accountID), Username: fmt.Sprintf("group_%s", req.BelongGroupId), Name: fmt.Sprintf("组%s", req.BelongGroupId), CreatedAt: time.Now().Unix(), UpdatedAt: time.Now().Unix(), CurrentCnt: 0},
			BelongTeam:  nil, // 不允许填写归属小队
			CreatedAt:   time.Now().Unix(),
			UpdatedAt:   time.Now().Unix(),
			LastLoginAt: 0,
		}

		res.Data = newAccount
		res.Success = true
		res.Message = "账户创建成功"

		c.Payload(res)
	}
}
