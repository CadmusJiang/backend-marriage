package account

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/xinliangnote/go-gin-api/internal/code"
	"github.com/xinliangnote/go-gin-api/internal/pkg/core"
	"github.com/xinliangnote/go-gin-api/internal/repository/mysql"
	"github.com/xinliangnote/go-gin-api/internal/repository/mysql/account"
	orgrepo "github.com/xinliangnote/go-gin-api/internal/repository/mysql/org"
	"go.uber.org/zap"
)

type detailRequest struct {
	AccountID    string `uri:"accountId" binding:"required"` // 账户ID
	IncludeGroup string `form:"includeGroup"`                // 是否包含组信息
	IncludeTeam  string `form:"includeTeam"`                 // 是否包含团队信息
}

type detailResponse struct {
	Data    accountData `json:"data"`
	Success bool        `json:"success"`
}

// GetAccountDetail 获取账户详情
// @Summary 获取账户详情
// @Description 根据账户ID获取详细信息
// @Tags Account
// @Accept application/x-www-form-urlencoded
// @Produce json
// @Param accountId path string true "账户ID"
// @Param includeGroup query string false "是否包含组信息" default(true)
// @Param includeTeam query string false "是否包含团队信息" default(true)
// @Success 200 {object} detailResponse
// @Failure 400 {object} code.Failure
// @Failure 404 {object} code.Failure
// @Router /api/v1/accounts/{accountId} [get]
func (h *handler) GetAccountDetail() core.HandlerFunc {
	return func(c core.Context) {
		req := new(detailRequest)
		res := new(detailResponse)

		if err := c.ShouldBindURI(req); err != nil {
			c.AbortWithError(core.Error(
				http.StatusBadRequest,
				code.ParamBindError,
				code.Text(code.ParamBindError)).WithError(err),
			)
			return
		}

		// 绑定查询参数
		if err := c.ShouldBindForm(req); err != nil {
			c.AbortWithError(core.Error(
				http.StatusBadRequest,
				code.ParamBindError,
				code.Text(code.ParamBindError)).WithError(err),
			)
			return
		}

		// 查询账户信息
		accountQueryBuilder := account.NewQueryBuilder()

		// 支持通过 ID 或 username 查找账户
		if isNumeric(req.AccountID) {
			// 如果是数字，按ID查询
			var id int32
			fmt.Sscanf(req.AccountID, "%d", &id)
			accountQueryBuilder.WhereId(mysql.EqualPredicate, id)
		} else {
			// 如果是字符串，按username查询
			accountQueryBuilder.WhereUsername(mysql.EqualPredicate, req.AccountID)
		}

		accountInfo, err := accountQueryBuilder.QueryOne(h.db.GetDbR())
		if err != nil {
			h.logger.Error("查询账户失败", zap.Error(err))
			c.AbortWithError(core.Error(
				http.StatusInternalServerError,
				code.AdminDetailError,
				"查询账户失败").WithError(err),
			)
			return
		}

		if accountInfo == nil {
			c.AbortWithError(core.Error(
				http.StatusNotFound,
				code.AdminDetailError,
				"账户不存在"),
			)
			return
		}

		// 构建响应数据
		accountData := accountData{
			ID:          fmt.Sprintf("%d", accountInfo.Id),
			Username:    accountInfo.Username,
			Name:        accountInfo.Nickname,
			Phone:       accountInfo.Phone,
			RoleType:    accountInfo.RoleType,
			Status:      accountInfo.Status,
			CreatedAt:   int64(accountInfo.CreatedTimestamp),
			UpdatedAt:   int64(accountInfo.ModifiedTimestamp),
			LastLoginAt: int64(accountInfo.LastLoginTimestamp), // 改为lastLoginAt
		}

		// 账户接口不做脱敏

		// 根据角色类型和includeGroup参数决定是否包含组信息
		if strings.ToLower(req.IncludeGroup) != "false" {
			switch accountInfo.RoleType {
			case "company_manager":
				// company_manager不能有归属组
				accountData.BelongGroup = nil
			case "group_manager", "team_manager", "employee":
				// 这些角色必须有归属组
				if accountInfo.BelongGroupId > 0 {
					groupInfo := h.getGroupInfo(int(accountInfo.BelongGroupId))
					if groupInfo != nil {
						accountData.BelongGroup = groupInfo
					}
				}
			}
		}

		// 根据角色类型和includeTeam参数决定是否包含团队信息
		if strings.ToLower(req.IncludeTeam) != "false" {
			switch accountInfo.RoleType {
			case "company_manager", "group_manager":
				// company_manager和group_manager不能有归属团队
				accountData.BelongTeam = nil
			case "team_manager":
				// team_manager必须有归属团队
				if accountInfo.BelongTeamId > 0 {
					teamInfo := h.getTeamInfo(int(accountInfo.BelongTeamId))
					if teamInfo != nil {
						accountData.BelongTeam = teamInfo
					}
				}
			case "employee":
				// employee可以有归属团队（可选）
				if accountInfo.BelongTeamId > 0 {
					teamInfo := h.getTeamInfo(int(accountInfo.BelongTeamId))
					if teamInfo != nil {
						accountData.BelongTeam = teamInfo
					}
				}
			}
		}

		res.Data = accountData
		res.Success = true

		c.Payload(res)
	}
}

// Me 返回当前用户信息（示例使用 sessionUserInfo 或简单查询）
func (h *handler) Me() core.HandlerFunc {
	return func(c core.Context) {
		// 简化：从 accountId 查询或返回固定用户
		// 这里可以结合登录态信息返回当前用户
		c.Payload(map[string]interface{}{"success": true, "data": map[string]interface{}{"id": "u_1", "username": "admin", "name": "系统管理员", "roleType": "company_manager", "status": "enabled", "phone": "13800000000"}})
	}
}

// isNumeric 检查字符串是否为数字
func isNumeric(s string) bool {
	for _, r := range s {
		if r < '0' || r > '9' {
			return false
		}
	}
	return len(s) > 0
}

// getGroupInfo 根据组ID获取组信息
func (h *handler) getGroupInfo(groupID int) *org {
	rec := new(orgrepo.Org)
	err := h.db.GetDbR().Where("id = ? AND org_type = ?", groupID, 1).First(rec).Error
	if err != nil {
		return nil
	}
	return &org{
		ID:         uint64(rec.Id),
		Username:   rec.Username,
		Name:       rec.Nickname,
		CreatedAt:  rec.CreatedTimestamp,
		UpdatedAt:  rec.ModifiedTimestamp,
		CurrentCnt: 0,
	}
}

// getTeamInfo 根据团队ID获取团队信息
func (h *handler) getTeamInfo(teamID int) *org {
	rec := new(orgrepo.Org)
	err := h.db.GetDbR().Where("id = ? AND org_type = ?", teamID, 2).First(rec).Error
	if err != nil {
		return nil
	}
	return &org{
		ID:         uint64(rec.Id),
		Username:   rec.Username,
		Name:       rec.Nickname,
		CreatedAt:  rec.CreatedTimestamp,
		UpdatedAt:  rec.ModifiedTimestamp,
		CurrentCnt: 0,
	}
}
