package organization

import (
	"net/http"
	"strconv"

	"github.com/xinliangnote/go-gin-api/internal/code"
	"github.com/xinliangnote/go-gin-api/internal/pkg/core"
	orgsvc "github.com/xinliangnote/go-gin-api/internal/services/organization"
	"go.uber.org/zap"
)

// UpdateGroup 更新组
func (h *handler) UpdateGroup() core.HandlerFunc {
	return func(c core.Context) {
		groupID := c.Param("groupId")
		if groupID == "" {
			c.AbortWithError(core.Error(
				http.StatusBadRequest,
				code.ParamBindError,
				"组ID不能为空"),
			)
			return
		}

		id, err := strconv.ParseUint(groupID, 10, 32)
		if err != nil {
			c.AbortWithError(core.Error(
				http.StatusBadRequest,
				code.ParamBindError,
				"组ID格式错误").WithError(err),
			)
			return
		}

		var payload orgsvc.UpdateGroupPayload
		if err := c.ShouldBindJSON(&payload); err != nil {
			c.AbortWithError(core.Error(
				http.StatusBadRequest,
				code.ParamBindError,
				"参数绑定失败").WithError(err),
			)
			return
		}

		// 调用service层
		_, updateErr := h.orgService.UpdateGroup(h.createContext(c), uint32(id), &payload)
		if updateErr != nil {
			h.logger.Error("更新组失败", zap.Error(updateErr))
			c.AbortWithError(core.Error(
				http.StatusInternalServerError,
				code.ServerError,
				"更新组失败").WithError(updateErr),
			)
			return
		}

		c.Payload(map[string]interface{}{
			"message": "更新组成功",
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

		id, err := strconv.ParseUint(teamID, 10, 32)
		if err != nil {
			c.AbortWithError(core.Error(
				http.StatusBadRequest,
				code.ParamBindError,
				"团队ID格式错误").WithError(err),
			)
			return
		}

		var payload orgsvc.UpdateTeamPayload
		if err := c.ShouldBindJSON(&payload); err != nil {
			c.AbortWithError(core.Error(
				http.StatusBadRequest,
				code.ParamBindError,
				"参数绑定失败").WithError(err),
			)
			return
		}

		// 调用service层
		_, updateErr := h.orgService.UpdateTeam(h.createContext(c), uint32(id), &payload)
		if updateErr != nil {
			h.logger.Error("更新团队失败", zap.Error(updateErr))
			c.AbortWithError(core.Error(
				http.StatusInternalServerError,
				code.ServerError,
				"更新团队失败").WithError(updateErr),
			)
			return
		}

		c.Payload(map[string]interface{}{
			"message": "更新团队成功",
		})
	}
}
