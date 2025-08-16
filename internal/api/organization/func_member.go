package organization

import (
	"net/http"
	"strconv"

	"github.com/xinliangnote/go-gin-api/internal/code"
	"github.com/xinliangnote/go-gin-api/internal/pkg/core"
	orgsvc "github.com/xinliangnote/go-gin-api/internal/services/organization"
	"go.uber.org/zap"
)

// ListTeamMembers 获取团队成员列表
func (h *handler) ListTeamMembers() core.HandlerFunc {
	return func(c core.Context) {
		teamID := c.Param("teamId")
		if teamID == "" {
			c.AbortWithError(core.Error(
				http.StatusBadRequest,
				code.ParamBindError,
				"团队ID不能为空"),
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
		members, total, err := h.orgService.ListMembers(h.createContext(c), teamID, current, pageSize, keyword)
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
		teamID := c.Param("teamId") // 修复：使用正确的参数名
		if teamID == "" {
			c.AbortWithError(core.Error(
				http.StatusBadRequest,
				code.ParamBindError,
				"团队ID不能为空"),
			)
			return
		}

		var payload struct {
			AccountID string `json:"accountId" binding:"required"`
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
		err := h.orgService.AddMember(h.createContext(c), teamID, payload.AccountID)
		if err != nil {
			h.logger.Error("添加成员失败", zap.Error(err))

			// 检查是否是自定义错误类型
			if memberErr, ok := err.(*orgsvc.MemberOperationError); ok {
				// 返回详细的错误信息
				c.AbortWithError(core.Error(
					http.StatusBadRequest,
					code.ParamBindError,
					memberErr.Message).WithError(err),
				)
				return
			}

			// 其他错误返回通用错误信息
			c.AbortWithError(core.Error(
				http.StatusInternalServerError,
				code.ServerError,
				"添加成员失败").WithError(err),
			)
			return
		}

		// 获取更新后的团队信息
		updatedTeam, err := h.orgService.GetTeam(h.createContext(c), teamID)
		if err != nil {
			h.logger.Error("获取更新后的团队信息失败", zap.Error(err))
			c.AbortWithError(core.Error(
				http.StatusInternalServerError,
				code.ServerError,
				"获取更新后的团队信息失败").WithError(err),
			)
			return
		}

		c.Payload(map[string]interface{}{
			"message": "添加成员成功",
			"data":    updatedTeam,
		})
	}
}

// RemoveTeamMember 移除团队成员
func (h *handler) RemoveTeamMember() core.HandlerFunc {
	return func(c core.Context) {
		teamID := c.Param("teamId")
		accountID := c.Param("accountId") // 使用accountId参数
		if teamID == "" || accountID == "" {
			c.AbortWithError(core.Error(
				http.StatusBadRequest,
				code.ParamBindError,
				"团队ID和账户ID不能为空"),
			)
			return
		}

		// 调用service层
		err := h.orgService.RemoveMember(h.createContext(c), teamID, accountID)
		if err != nil {
			h.logger.Error("移除成员失败", zap.Error(err))

			// 检查是否是自定义错误类型
			if memberErr, ok := err.(*orgsvc.MemberOperationError); ok {
				// 返回详细的错误信息
				c.AbortWithError(core.Error(
					http.StatusBadRequest,
					code.ParamBindError,
					memberErr.Message).WithError(err),
				)
				return
			}

			// 其他错误返回通用错误信息
			c.AbortWithError(core.Error(
				http.StatusInternalServerError,
				code.ServerError,
				"移除成员失败").WithError(err),
			)
			return
		}

		// 获取更新后的团队信息
		updatedTeam, err := h.orgService.GetTeam(h.createContext(c), teamID)
		if err != nil {
			h.logger.Error("获取更新后的团队信息失败", zap.Error(err))
			c.AbortWithError(core.Error(
				http.StatusInternalServerError,
				code.ServerError,
				"获取更新后的团队信息失败").WithError(err),
			)
			return
		}

		c.Payload(map[string]interface{}{
			"message": "移除成员成功",
			"data":    updatedTeam,
		})
	}
}

// UpdateTeamMemberRole 更新成员角色
func (h *handler) UpdateTeamMemberRole() core.HandlerFunc {
	return func(c core.Context) {
		teamID := c.Param("teamId")
		accountID := c.Param("accountId") // 使用accountId参数
		if teamID == "" || accountID == "" {
			c.AbortWithError(core.Error(
				http.StatusBadRequest,
				code.ParamBindError,
				"团队ID和账户ID不能为空"),
			)
			return
		}

		var payload struct {
			RoleType string `json:"roleType" binding:"required"`
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
		err := h.orgService.UpdateMemberRole(h.createContext(c), teamID, accountID, payload.RoleType)
		if err != nil {
			h.logger.Error("更新成员角色失败", zap.Error(err))

			// 检查是否是自定义错误类型
			if memberErr, ok := err.(*orgsvc.MemberOperationError); ok {
				// 返回详细的错误信息
				c.AbortWithError(core.Error(
					http.StatusBadRequest,
					code.ParamBindError,
					memberErr.Message).WithError(err),
				)
				return
			}

			// 其他错误返回通用错误信息
			c.AbortWithError(core.Error(
				http.StatusInternalServerError,
				code.ServerError,
				"更新成员角色失败").WithError(err),
			)
			return
		}

		c.Payload(map[string]interface{}{
			"message": "更新成员角色成功",
		})
	}
}

// UpdateTeamMember 更新团队成员信息（包括角色等）
func (h *handler) UpdateTeamMember() core.HandlerFunc {
	return func(c core.Context) {
		teamID := c.Param("teamId")
		accountID := c.Param("accountId")
		if teamID == "" || accountID == "" {
			c.AbortWithError(core.Error(
				http.StatusBadRequest,
				code.ParamBindError,
				"团队ID和账户ID不能为空"),
			)
			return
		}

		var payload struct {
			RoleType string `json:"roleType,omitempty"`
			// 可以在这里添加其他需要更新的字段
		}
		if err := c.ShouldBindJSON(&payload); err != nil {
			c.AbortWithError(core.Error(
				http.StatusBadRequest,
				code.ParamBindError,
				"参数绑定失败").WithError(err),
			)
			return
		}

		// 如果包含角色更新，调用service层更新角色
		if payload.RoleType != "" {
			err := h.orgService.UpdateMemberRole(h.createContext(c), teamID, accountID, payload.RoleType)
			if err != nil {
				h.logger.Error("更新成员角色失败", zap.Error(err))

				// 检查是否是自定义错误类型
				if memberErr, ok := err.(*orgsvc.MemberOperationError); ok {
					// 返回详细的错误信息
					c.AbortWithError(core.Error(
						http.StatusBadRequest,
						code.ParamBindError,
						memberErr.Message).WithError(err),
					)
					return
				}

				// 其他错误返回通用错误信息
				c.AbortWithError(core.Error(
					http.StatusInternalServerError,
					code.ServerError,
					"更新成员角色失败").WithError(err),
				)
				return
			}
		}

		// 获取更新后的团队信息
		updatedTeam, err := h.orgService.GetTeam(h.createContext(c), teamID)
		if err != nil {
			h.logger.Error("获取更新后的团队信息失败", zap.Error(err))
			c.AbortWithError(core.Error(
				http.StatusInternalServerError,
				code.ServerError,
				"获取更新后的团队信息失败").WithError(err),
			)
			return
		}

		c.Payload(map[string]interface{}{
			"message": "更新团队成员成功",
			"data":    updatedTeam,
		})
	}
}
