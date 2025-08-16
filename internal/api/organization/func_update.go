package organization

import (
	"net/http"

	"github.com/xinliangnote/go-gin-api/internal/code"
	"github.com/xinliangnote/go-gin-api/internal/pkg/core"
	orgsvc "github.com/xinliangnote/go-gin-api/internal/services/organization"
	"go.uber.org/zap"
)

// UpdateGroup 更新组
func (h *handler) UpdateGroup() core.HandlerFunc {
	return func(c core.Context) {
		groupID := c.Param("orgId")
		if groupID == "" {
			c.AbortWithError(core.Error(
				http.StatusBadRequest,
				code.ParamBindError,
				"组ID不能为空"),
			)
			return
		}

		var payload orgsvc.UpdateGroupPayload
		if err := c.ShouldBindJSON(&payload); err != nil {
			c.AbortWithError(core.Error(
				http.StatusBadRequest,
				code.ParamBindError,
				"参数绑定失败"),
			)
			return
		}

		// 调用service层
		_, updateErr := h.orgService.UpdateGroup(h.createContext(c), groupID, &payload)
		if updateErr != nil {
			h.logger.Error("更新组失败", zap.Error(updateErr))
			c.AbortWithError(core.Error(
				http.StatusInternalServerError,
				code.ServerError,
				"更新组失败").WithError(updateErr),
			)
			return
		}

		// 获取更新后的组信息
		updatedGroup, err := h.orgService.GetGroup(h.createContext(c), groupID)
		if err != nil {
			h.logger.Error("获取更新后的组信息失败", zap.Error(err))
			c.AbortWithError(core.Error(
				http.StatusInternalServerError,
				code.ServerError,
				"获取更新后的组信息失败").WithError(err),
			)
			return
		}

		c.Payload(map[string]interface{}{
			"message": "更新组成功",
			"data":    updatedGroup,
		})
	}
}

// UpdateTeam 更新团队
func (h *handler) UpdateTeam() core.HandlerFunc {
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

		var payload orgsvc.UpdateTeamPayload
		if err := c.ShouldBindJSON(&payload); err != nil {
			c.AbortWithError(core.Error(
				http.StatusBadRequest,
				code.ParamBindError,
				"参数绑定失败"),
			)
			return
		}

		// 调用service层
		_, updateErr := h.orgService.UpdateTeam(h.createContext(c), teamID, &payload)
		if updateErr != nil {
			h.logger.Error("更新团队失败", zap.Error(updateErr))
			c.AbortWithError(core.Error(
				http.StatusInternalServerError,
				code.ServerError,
				"更新团队失败").WithError(updateErr),
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
			"message": "更新团队成功",
			"data":    updatedTeam,
		})
	}
}
