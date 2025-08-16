package organization

import (
	"net/http"

	"github.com/xinliangnote/go-gin-api/internal/code"
	"github.com/xinliangnote/go-gin-api/internal/pkg/core"
	orgsvc "github.com/xinliangnote/go-gin-api/internal/services/organization"
	"go.uber.org/zap"
)

// CreateGroup 创建组
func (h *handler) CreateGroup() core.HandlerFunc {
	return func(c core.Context) {
		var payload orgsvc.CreateGroupPayload
		if err := c.ShouldBindJSON(&payload); err != nil {
			c.AbortWithError(core.Error(
				http.StatusBadRequest,
				code.ParamBindError,
				"参数绑定失败").WithError(err),
			)
			return
		}

		// 参数验证
		if payload.Username == "" || payload.Name == "" {
			c.AbortWithError(core.Error(
				http.StatusBadRequest,
				code.ParamBindError,
				"用户名和名称不能为空"),
			)
			return
		}

		// 调用service层
		id, err := h.orgService.CreateGroup(h.createContext(c), &payload)
		if err != nil {
			h.logger.Error("创建组失败", zap.Error(err))
			c.AbortWithError(core.Error(
				http.StatusInternalServerError,
				code.ServerError,
				"创建组失败").WithError(err),
			)
			return
		}

		c.Payload(map[string]interface{}{
			"data":    map[string]interface{}{"id": id},
			"message": "创建组成功",
		})
	}
}

// CreateTeam 创建团队
func (h *handler) CreateTeam() core.HandlerFunc {
	return func(c core.Context) {
		var payload orgsvc.CreateTeamPayload
		if err := c.ShouldBindJSON(&payload); err != nil {
			c.AbortWithError(core.Error(
				http.StatusBadRequest,
				code.ParamBindError,
				"参数绑定失败").WithError(err),
			)
			return
		}

		// 参数验证
		if payload.Username == "" || payload.Name == "" {
			c.AbortWithError(core.Error(
				http.StatusBadRequest,
				code.ParamBindError,
				"用户名和名称不能为空"),
			)
			return
		}

		if payload.BelongGroupId == 0 {
			c.AbortWithError(core.Error(
				http.StatusBadRequest,
				code.ParamBindError,
				"所属组ID不能为空"),
			)
			return
		}

		// 调用service层
		id, err := h.orgService.CreateTeam(h.createContext(c), &payload)
		if err != nil {
			h.logger.Error("创建团队失败", zap.Error(err))
			c.AbortWithError(core.Error(
				http.StatusInternalServerError,
				code.ServerError,
				"创建团队失败").WithError(err),
			)
			return
		}

		c.Payload(map[string]interface{}{
			"data":    map[string]interface{}{"id": id},
			"message": "创建团队成功",
		})
	}
}
