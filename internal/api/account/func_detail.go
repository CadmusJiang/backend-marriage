package account

import (
	"fmt"
	"net/http"
	"strings"
	"time"

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
	Data accountData `json:"data"`
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
			ID:        fmt.Sprintf("%d", accountInfo.Id),
			Username:  accountInfo.Username,
			Name:      accountInfo.Name,
			Phone:     accountInfo.Phone,
			RoleType:  accountInfo.RoleType,
			Status:    accountInfo.Status, // 直接赋值，不需要格式化
			CreatedAt: accountInfo.CreatedAt.Unix(),
			UpdatedAt: accountInfo.UpdatedAt.Unix(),
			LastLoginAt: func() int64 {
				if accountInfo.LastLoginAt != nil {
					return accountInfo.LastLoginAt.Unix()
				}
				return 0
			}(),
		}

		// 账户接口不做脱敏

		// 根据includeGroup和includeTeam参数决定是否包含组织信息
		if strings.ToLower(req.IncludeGroup) != "false" || strings.ToLower(req.IncludeTeam) != "false" {
			// 通过组织关系表查询账户的归属信息
			belongGroup, belongTeam := h.getAccountOrgInfo(int(accountInfo.Id))

			if strings.ToLower(req.IncludeGroup) != "false" {
				accountData.BelongGroup = belongGroup
			}

			if strings.ToLower(req.IncludeTeam) != "false" {
				accountData.BelongTeam = belongTeam
			}
		}

		res.Data = accountData

		c.Payload(res)
	}
}

// Me 获取当前用户信息
// @Summary 获取当前用户信息
// @Description 获取当前登录用户的详细信息
// @Tags CoreAuth
// @Accept application/x-www-form-urlencoded
// @Produce json
// @Success 200 {object} detailResponse
// @Failure 400 {object} code.Failure
// @Failure 401 {object} code.Failure
// @Router /api/v1/me [get]
func (h *handler) Me() core.HandlerFunc {
	return func(c core.Context) {
		res := new(detailResponse)

		// 从会话中获取当前用户ID
		sessionUserInfo := c.SessionUserInfo()
		if sessionUserInfo.UserID == 0 {
			c.AbortWithError(core.Error(
				http.StatusUnauthorized,
				code.AuthorizationError,
				"用户未登录").WithError(fmt.Errorf("session user info is empty")),
			)
			return
		}

		// 查询当前用户信息
		accountQueryBuilder := account.NewQueryBuilder()
		accountQueryBuilder.WhereId(mysql.EqualPredicate, sessionUserInfo.UserID)

		accountInfo, err := accountQueryBuilder.QueryOne(h.db.GetDbR())
		if err != nil {
			h.logger.Error("查询当前用户失败", zap.Error(err))
			c.AbortWithError(core.Error(
				http.StatusInternalServerError,
				code.AdminDetailError,
				"查询用户信息失败").WithError(err),
			)
			return
		}

		if accountInfo == nil {
			c.AbortWithError(core.Error(
				http.StatusNotFound,
				code.AdminDetailError,
				"用户不存在"),
			)
			return
		}

		// 构建响应数据
		accountData := accountData{
			ID:        fmt.Sprintf("%d", accountInfo.Id),
			Username:  accountInfo.Username,
			Name:      accountInfo.Name,
			Phone:     accountInfo.Phone,
			RoleType:  accountInfo.RoleType,
			Status:    accountInfo.Status, // 直接赋值，不需要格式化
			CreatedAt: accountInfo.CreatedAt.Unix(),
			UpdatedAt: accountInfo.UpdatedAt.Unix(),
			LastLoginAt: func() int64 {
				if accountInfo.LastLoginAt != nil {
					return accountInfo.LastLoginAt.Unix()
				}
				return 0
			}(),
		}

		// 直接查询归属组信息（如果有的话）
		if accountInfo.BelongGroupId > 0 {
			groupInfo := h.getGroupInfo(int(accountInfo.BelongGroupId))
			if groupInfo != nil {
				accountData.BelongGroup = groupInfo
			}
		}

		// 直接查询归属团队信息（如果有的话）
		if accountInfo.BelongTeamId > 0 {
			teamInfo := h.getTeamInfo(int(accountInfo.BelongTeamId))
			if teamInfo != nil {
				accountData.BelongTeam = teamInfo
			}
		}

		res.Data = accountData

		c.Payload(res)
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

// getAccountOrgInfo 通过组织关系表查询账户的归属信息
func (h *handler) getAccountOrgInfo(accountID int) (belongGroup *org, belongTeam *org) {
	// 查询账户的所有组织关系
	var relations []struct {
		OrgID        uint32    `gorm:"column:org_id"`
		OrgType      string    `gorm:"column:org_type"`
		OrgUsername  string    `gorm:"column:org_username"`
		OrgName      string    `gorm:"column:org_name"`
		OrgCreatedAt time.Time `gorm:"column:org_created_at"`
		OrgUpdatedAt time.Time `gorm:"column:org_updated_at"`
	}

	// 联表查询：account_org_relation + org
	query := `
		SELECT 
			o.id as org_id,
			o.org_type,
			o.username as org_username,
			o.name as org_name,
			o.created_at as org_created_at,
			o.updated_at as org_updated_at
		FROM account_org_relation aor
		JOIN org o ON aor.org_id = o.id
		WHERE aor.account_id = ? 
		AND aor.relation_type = 'belong' 
		AND aor.status = 'active'
		AND o.status = 'enabled'
	`

	err := h.db.GetDbR().Raw(query, accountID).Scan(&relations).Error
	if err != nil {
		return nil, nil
	}

	// 处理查询结果
	for _, rel := range relations {
		orgInfo := &org{
			ID:         uint64(rel.OrgID),
			Username:   rel.OrgUsername,
			Name:       rel.OrgName,
			CreatedAt:  rel.OrgCreatedAt.Unix(),
			UpdatedAt:  rel.OrgUpdatedAt.Unix(),
			CurrentCnt: 0, // 暂时设为0，如果需要可以再查询
		}

		if rel.OrgType == "group" {
			belongGroup = orgInfo
		} else if rel.OrgType == "team" {
			belongTeam = orgInfo
		}
	}

	return belongGroup, belongTeam
}

// getGroupInfo 根据组ID获取组信息（保留兼容性）
func (h *handler) getGroupInfo(groupID int) *org {
	rec := new(orgrepo.Org)
	err := h.db.GetDbR().Where("id = ? AND org_type = ?", groupID, 1).First(rec).Error
	if err != nil {
		return nil
	}
	return &org{
		ID:         uint64(rec.Id),
		Username:   rec.Username,
		Name:       rec.Name,
		CreatedAt:  rec.CreatedAt.Unix(),
		UpdatedAt:  rec.UpdatedAt.Unix(),
		CurrentCnt: 0,
	}
}

// getTeamInfo 根据团队ID获取团队信息（保留兼容性）
func (h *handler) getTeamInfo(teamID int) *org {
	rec := new(orgrepo.Org)
	err := h.db.GetDbR().Where("id = ? AND org_type = ?", teamID, 2).First(rec).Error
	if err != nil {
		return nil
	}
	return &org{
		ID:         uint64(rec.Id),
		Username:   rec.Username,
		Name:       rec.Name,
		CreatedAt:  rec.CreatedAt.Unix(),
		UpdatedAt:  rec.UpdatedAt.Unix(),
		CurrentCnt: 0,
	}
}
