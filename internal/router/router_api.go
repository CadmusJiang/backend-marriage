package router

import (
	"github.com/xinliangnote/go-gin-api/internal/api/account"
	"github.com/xinliangnote/go-gin-api/internal/api/helper"
	"github.com/xinliangnote/go-gin-api/internal/api/org_info"
	"github.com/xinliangnote/go-gin-api/internal/pkg/core"
)

func setApiRouter(r *resource) {
	// account
	accountHandler := account.New(r.logger, r.db, r.cache)

	// helper
	helperHandler := helper.New(r.logger, r.db, r.cache)

	// org_info
	orgInfoHandler := org_info.New(r.logger, r.db, r.cache)

	// API v1 路由组 - 无需登录验证
	apiV1 := r.mux.Group("/api/v1")
	{
		// 认证相关接口
		authAPI := apiV1.Group("/auth")
		{
			authAPI.POST("/login", accountHandler.Login())
			authAPI.POST("/logout", accountHandler.Logout())
		}

		// 数据库检查接口
		apiV1.GET("/check-db", helperHandler.CheckDatabase())

		// 日志查看接口 - 无需认证
		apiV1.GET("/logs", helperHandler.GetLogs())
		apiV1.GET("/logs/realtime", helperHandler.GetLogsRealtime())
	}

	// 需要签名验证、登录验证、RBAC 权限验证
	api := r.mux.Group("/api", core.WrapAuthHandler(r.interceptors.CheckLogin), r.interceptors.CheckRBAC())
	{
		// account management
		api.GET("/v1/accounts", accountHandler.GetAccountList())
		api.POST("/v1/accounts", accountHandler.CreateAccount())
		api.GET("/v1/accounts/:accountId", accountHandler.GetAccountDetail())
		api.PUT("/v1/accounts/:accountId", accountHandler.UpdateAccount())
		api.GET("/v1/account-histories", accountHandler.GetAccountHistories())

		// groups management - 专门获取groups信息
		api.GET("/v1/groups", orgInfoHandler.GetGroups())

		// teams management - 专门获取teams信息
		api.GET("/v1/teams", orgInfoHandler.GetTeams())

		// org_info management
		api.GET("/v1/org-infos", orgInfoHandler.GetOrgInfoList())
		api.POST("/v1/org-infos", orgInfoHandler.CreateOrgInfo())
		api.GET("/v1/org-infos/:orgId", orgInfoHandler.GetOrgInfoDetail())
		api.PUT("/v1/org-infos/:orgId", orgInfoHandler.UpdateOrgInfo())
		api.DELETE("/v1/org-infos/:orgId", orgInfoHandler.DeleteOrgInfo())
		api.GET("/v1/org-infos/:orgId/children", orgInfoHandler.GetOrgInfoChildren())
		api.GET("/v1/org-infos/:orgId/parent", orgInfoHandler.GetOrgInfoParent())
	}
}
