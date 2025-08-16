package organization

import (
	"net/http"
	"strconv"

	"github.com/xinliangnote/go-gin-api/internal/authz"
	"github.com/xinliangnote/go-gin-api/internal/code"
	"github.com/xinliangnote/go-gin-api/internal/pkg/core"
	"go.uber.org/zap"
)

// GetGroups 获取组列表
func (h *handler) GetGroups() core.HandlerFunc {
	return func(c core.Context) {
		// 获取查询参数
		params := c.RequestInputParams()
		groupID, _ := strconv.ParseUint(params.Get("groupId"), 10, 32)
		current, _ := strconv.Atoi(params.Get("current"))
		if current == 0 {
			current = 1
		}
		pageSize, _ := strconv.Atoi(params.Get("pageSize"))
		if pageSize == 0 {
			pageSize = 10
		}
		keyword := params.Get("keyword")

		// 计算用户访问范围
		scope, scopeErr := authz.ComputeScope(c, h.db)
		if scopeErr != nil {
			h.logger.Error("计算访问范围失败", zap.Error(scopeErr))
			c.AbortWithError(core.Error(
				http.StatusInternalServerError,
				code.ServerError,
				"权限验证失败").WithError(scopeErr),
			)
			return
		}

		// 获取有效的权限范围（考虑请求参数和用户权限的交集）
		effectiveScope := authz.GetEffectiveScope(scope, uint32(groupID), 0)

		// 声明变量
		var groups []map[string]interface{}
		var total int64
		var err error

		// 直接调用service层，传递权限范围参数
		groups, total, err = h.orgService.ListGroups(h.createContext(c), current, pageSize, keyword, &effectiveScope)

		if err != nil {
			h.logger.Error("获取组列表失败", zap.Error(err))
			c.AbortWithError(core.Error(
				http.StatusInternalServerError,
				code.ServerError,
				"获取组列表失败").WithError(err),
			)
			return
		}

		c.Payload(map[string]interface{}{
			"data": groups,
			"meta": map[string]interface{}{
				"total":    total,
				"pageSize": pageSize,
				"current":  current,
			},
		})
	}
}

// ListTeams 获取团队列表
func (h *handler) ListTeams() core.HandlerFunc {
	return func(c core.Context) {
		// 获取查询参数
		params := c.RequestInputParams()
		belongGroupID, _ := strconv.ParseUint(params.Get("belongGroupId"), 10, 32)
		teamID, _ := strconv.ParseUint(params.Get("teamId"), 10, 32)
		current, _ := strconv.Atoi(params.Get("current"))
		if current == 0 {
			current = 1
		}
		pageSize, _ := strconv.Atoi(params.Get("pageSize"))
		if pageSize == 0 {
			pageSize = 10
		}
		keyword := params.Get("keyword")

		// 计算用户访问范围
		scope, scopeErr := authz.ComputeScope(c, h.db)
		if scopeErr != nil {
			h.logger.Error("计算访问范围失败", zap.Error(scopeErr))
			c.AbortWithError(core.Error(
				http.StatusInternalServerError,
				code.ServerError,
				"权限验证失败").WithError(scopeErr),
			)
			return
		}

		// 获取有效的权限范围（考虑请求参数和用户权限的交集）
		effectiveScope := authz.GetEffectiveScope(scope, uint32(belongGroupID), uint32(teamID))

		// 声明变量
		var teams []map[string]interface{}
		var total int64
		var err error

		// 直接调用service层，传递权限范围参数
		teams, total, err = h.orgService.ListTeams(h.createContext(c), uint32(belongGroupID), current, pageSize, keyword, &effectiveScope)

		if err != nil {
			h.logger.Error("获取团队列表失败", zap.Error(err))
			c.AbortWithError(core.Error(
				http.StatusInternalServerError,
				code.ServerError,
				"获取团队列表失败").WithError(err),
			)
			return
		}

		c.Payload(map[string]interface{}{
			"data": teams,
			"meta": map[string]interface{}{
				"total":    total,
				"pageSize": pageSize,
				"current":  current,
			},
		})
	}
}
