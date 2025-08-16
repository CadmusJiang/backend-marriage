package organization

import (
	"net/http"
	"strconv"

	"github.com/xinliangnote/go-gin-api/internal/code"
	"github.com/xinliangnote/go-gin-api/internal/pkg/core"
	"go.uber.org/zap"
)

// GetOrgInfoDetail 获取组织详情
func (h *handler) GetOrgInfoDetail() core.HandlerFunc {
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

		// 根据组织类型调用不同的service方法
		params := c.RequestInputParams()
		orgType := params.Get("orgType")
		var orgInfo map[string]interface{}
		var serviceErr error

		if orgType == "team" {
			orgInfo, serviceErr = h.orgService.GetTeam(h.createContext(c), uint32(id))
		} else {
			orgInfo, serviceErr = h.orgService.GetGroup(h.createContext(c), uint32(id))
		}

		if serviceErr != nil {
			h.logger.Error("获取组织详情失败", zap.Error(serviceErr))
			c.AbortWithError(core.Error(
				http.StatusInternalServerError,
				code.ServerError,
				"获取组织详情失败").WithError(serviceErr),
			)
			return
		}

		c.Payload(map[string]interface{}{
			"data": orgInfo,
		})
	}
}

// GetTeam 获取团队详情
func (h *handler) GetTeam() core.HandlerFunc {
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

		// 调用service层
		teamInfo, err := h.orgService.GetTeam(h.createContext(c), uint32(id))
		if err != nil {
			h.logger.Error("获取团队详情失败", zap.Error(err))
			c.AbortWithError(core.Error(
				http.StatusInternalServerError,
				code.ServerError,
				"获取团队详情失败").WithError(err),
			)
			return
		}

		c.Payload(map[string]interface{}{
			"data": teamInfo,
		})
	}
}
