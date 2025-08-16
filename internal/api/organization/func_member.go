package organization

import (
	"net/http"
	"strconv"

	"github.com/xinliangnote/go-gin-api/internal/code"
	"github.com/xinliangnote/go-gin-api/internal/pkg/core"
	"go.uber.org/zap"
)

// ListTeamMembers 获取团队成员列表
func (h *handler) ListTeamMembers() core.HandlerFunc {
	return func(c core.Context) {
		orgID := c.Param("orgId")
		if orgID == "" {
			c.AbortWithError(core.Error(
				http.StatusBadRequest,
				code.ParamBindError,
				"组织ID不能为空"),
			)
			return
		}

		id, err := strconv.ParseUint(orgID, 10, 32)
		if err != nil {
			c.AbortWithError(core.Error(
				http.StatusBadRequest,
				code.ParamBindError,
				"组织ID格式错误").WithError(err),
			)
			return
		}

		params := c.RequestInputParams()
		current, _ := strconv.Atoi(params.Get("current"))
		if current == 0 {
			current = 1
		}
		pageSize, _ := strconv.Atoi(params.Get("pageSize"))
		if pageSize == 0 {
			pageSize = 10
		}
		keyword := params.Get("keyword")

		// 调用service层
		members, total, err := h.orgService.ListMembers(h.createContext(c), uint32(id), current, pageSize, keyword)
		if err != nil {
			h.logger.Error("获取成员列表失败", zap.Error(err))
			c.AbortWithError(core.Error(
				http.StatusInternalServerError,
				code.ServerError,
				"获取成员列表失败").WithError(err),
			)
			return
		}

		c.Payload(map[string]interface{}{
			"data": members,
			"meta": map[string]interface{}{
				"total":    total,
				"pageSize": pageSize,
				"current":  current,
			},
		})
	}
}

// AddTeamMember 添加团队成员
func (h *handler) AddTeamMember() core.HandlerFunc {
	return func(c core.Context) {
		orgID := c.Param("orgId")
		if orgID == "" {
			c.AbortWithError(core.Error(
				http.StatusBadRequest,
				code.ParamBindError,
				"组织ID不能为空"),
			)
			return
		}

		orgIDUint, err := strconv.ParseUint(orgID, 10, 32)
		if err != nil {
			c.AbortWithError(core.Error(
				http.StatusBadRequest,
				code.ParamBindError,
				"组织ID格式错误").WithError(err),
			)
			return
		}

		var payload struct {
			AccountID uint32 `json:"accountId" binding:"required"`
		}
		if err := c.ShouldBindJSON(&payload); err != nil {
			c.AbortWithError(core.Error(
				http.StatusBadRequest,
				code.ParamBindError,
				"参数绑定失败").WithError(err),
			)
			return
		}

		// 调用service层
		err = h.orgService.AddMember(h.createContext(c), uint32(orgIDUint), payload.AccountID)
		if err != nil {
			h.logger.Error("添加成员失败", zap.Error(err))
			c.AbortWithError(core.Error(
				http.StatusInternalServerError,
				code.ServerError,
				"添加成员失败").WithError(err),
			)
			return
		}

		c.Payload(map[string]interface{}{
			"message": "添加成员成功",
		})
	}
}

// RemoveTeamMember 移除团队成员
func (h *handler) RemoveTeamMember() core.HandlerFunc {
	return func(c core.Context) {
		orgID := c.Param("orgId")
		accountID := c.Param("accountId")
		if orgID == "" || accountID == "" {
			c.AbortWithError(core.Error(
				http.StatusBadRequest,
				code.ParamBindError,
				"组织ID和账户ID不能为空"),
			)
			return
		}

		orgIDUint, err := strconv.ParseUint(orgID, 10, 32)
		if err != nil {
			c.AbortWithError(core.Error(
				http.StatusBadRequest,
				code.ParamBindError,
				"组织ID格式错误").WithError(err),
			)
			return
		}

		accountIDUint, err := strconv.ParseUint(accountID, 10, 32)
		if err != nil {
			c.AbortWithError(core.Error(
				http.StatusBadRequest,
				code.ParamBindError,
				"账户ID格式错误").WithError(err),
			)
			return
		}

		// 调用service层
		err = h.orgService.RemoveMember(h.createContext(c), uint32(orgIDUint), uint32(accountIDUint))
		if err != nil {
			h.logger.Error("移除成员失败", zap.Error(err))
			c.AbortWithError(core.Error(
				http.StatusInternalServerError,
				code.ServerError,
				"移除成员失败").WithError(err),
			)
			return
		}

		c.Payload(map[string]interface{}{
			"message": "移除成员成功",
		})
	}
}

// UpdateTeamMemberRole 更新成员角色
func (h *handler) UpdateTeamMemberRole() core.HandlerFunc {
	return func(c core.Context) {
		// TODO: 实现更新成员角色逻辑
		c.Payload(map[string]interface{}{
			"message": "更新成员角色成功",
		})
	}
}
