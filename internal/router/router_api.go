package router

import (
	"github.com/xinliangnote/go-gin-api/internal/api/account"
	"github.com/xinliangnote/go-gin-api/internal/api/analytics"
	"github.com/xinliangnote/go-gin-api/internal/api/cooperation_store"
	"github.com/xinliangnote/go-gin-api/internal/api/customer_authorization_record"
	"github.com/xinliangnote/go-gin-api/internal/api/customer_authorization_record_history"
	"github.com/xinliangnote/go-gin-api/internal/api/logs"
	"github.com/xinliangnote/go-gin-api/internal/api/organization"
	"github.com/xinliangnote/go-gin-api/internal/pkg/core"
)

func setApiRouter(r *resource) {
	// account
	accountHandler := account.New(r.logger, r.db, r.cache)

	// cooperation_store
	cooperationStoreHandler := cooperation_store.New(r.logger, r.db, r.cache)

	// customer_authorization_record
	customerAuthRecordHandler := customer_authorization_record.New(r.logger, r.db, r.db)

	// customer_authorization_record_history
	customerAuthRecordHistoryHandler := customer_authorization_record_history.New(r.logger, r.db, r.cache)

	// helper 已移除

	// organization
	organizationHandler := organization.New(r.logger, r.db, r.cache)

	// analytics
	analyticsHandler := analytics.New(r.logger, r.db, r.db)

	// logs
	logsHandler := logs.New()

	// API v1 路由组 - 无需登录验证（CoreAuth）
	apiV1 := r.mux.Group("/api/v1")
	{
		// CoreAuth - 只有登录不需要认证
		authAPI := apiV1.Group("/auth")
		{
			authAPI.POST("/login", accountHandler.Login())
		}

		// 日志相关接口 - 无需认证
		apiV1.GET("/logs/latest", logsHandler.GetLatestLogs())
		apiV1.GET("/logs/unified", logsHandler.GetUnifiedLogs())
		apiV1.GET("/logs/paginated", logsHandler.GetPaginatedLogs())
		apiV1.GET("/logs/trace", logsHandler.GetTraceLogs())
		apiV1.GET("/logs/trace/range", logsHandler.GetTraceLogsByTimeRange())
	}

	// 需要签名验证、登录验证、RBAC 权限验证
	api := r.mux.Group("/api", core.WrapAuthHandler(r.interceptors.CheckLogin), r.interceptors.CheckRBAC())
	{
		// Organization
		api.GET("/v1/groups", organizationHandler.GetGroups())
		api.POST("/v1/groups", organizationHandler.CreateGroup())
		api.GET("/v1/groups/:orgId", organizationHandler.GetOrgInfoDetail())
		api.PUT("/v1/groups/:orgId", organizationHandler.UpdateGroup())
		api.GET("/v1/groups/:orgId/history", organizationHandler.GetGroupHistory())

		// teams endpoints
		api.GET("/v1/teams", organizationHandler.ListTeams())
		api.POST("/v1/teams", organizationHandler.CreateTeam())
		api.GET("/v1/teams/:teamId", organizationHandler.GetTeam())
		api.PUT("/v1/teams/:teamId", organizationHandler.UpdateTeam())
		api.GET("/v1/teams/:teamId/history", organizationHandler.GetTeamHistory())
		api.GET("/v1/teams/:teamId/members", organizationHandler.ListTeamMembers())
		api.POST("/v1/teams/:teamId/members", organizationHandler.AddTeamMember())
		api.DELETE("/v1/teams/:teamId/members/:accountId", organizationHandler.RemoveTeamMember())
		api.PATCH("/v1/teams/:teamId/members/:accountId", organizationHandler.UpdateTeamMember())
		api.GET("/v1/unassigned-account", organizationHandler.ListUnassignedAccounts())

		// Analytics
		api.GET("/v1/analytics/account/rankings", analyticsHandler.GetAccountRankings())
		api.GET("/v1/analytics/account/trends", analyticsHandler.GetAccountTrends())
		api.GET("/v1/analytics/teams/rankings", analyticsHandler.GetTeamsRankings())
		api.GET("/v1/analytics/teams/trends", analyticsHandler.GetTeamsTrends())
		api.GET("/v1/analytics/customer-authorization-record/composition", analyticsHandler.GetCustomerComposition())

		// Customer
		api.GET("/v1/customer-authorization-records", customerAuthRecordHandler.GetCustomerAuthorizationRecordList())
		api.POST("/v1/customer-authorization-records", customerAuthRecordHandler.CreateCustomerAuthorizationRecord())
		api.GET("/v1/customer-authorization-records/check-phone", customerAuthRecordHandler.CheckPhoneExistence())
		api.GET("/v1/customer-authorization-records/:id", customerAuthRecordHandler.GetCustomerAuthorizationRecordDetail())
		api.PATCH("/v1/customer-authorization-records/:id", customerAuthRecordHandler.UpdateCustomerAuthorizationRecord())
		api.PATCH("/v1/customer-authorization-records/:id/status", customerAuthRecordHandler.UpdateCustomerStatus())
		api.PATCH("/v1/customer-authorization-records/:id/financial", customerAuthRecordHandler.UpdateCustomerFinancial())
		api.GET("/v1/customer-authorization-records/:id/history", customerAuthRecordHistoryHandler.GetCustomerAuthorizationRecordHistoryList())

		// CooperationStore
		api.GET("/v1/cooperation-stores", cooperationStoreHandler.GetCooperationStoreList())
		api.POST("/v1/cooperation-stores", cooperationStoreHandler.CreateCooperationStore())
		api.GET("/v1/cooperation-stores/:id", cooperationStoreHandler.GetCooperationStoreDetail())
		api.PUT("/v1/cooperation-stores/:id", cooperationStoreHandler.UpdateCooperationStore())
		api.GET("/v1/cooperation-stores/:id/history", cooperationStoreHandler.GetCooperationStoreHistory())

		// Account
		api.GET("/v1/accounts", accountHandler.GetAccountList())
		api.POST("/v1/accounts", accountHandler.CreateAccount())
		api.GET("/v1/accounts/:accountId", accountHandler.GetAccountDetail())
		api.PUT("/v1/accounts/:accountId", accountHandler.UpdateAccount())
		api.PUT("/v1/accounts/:accountId/password", accountHandler.UpdateAccountPassword())
		api.GET("/v1/accounts/:accountId/history", accountHandler.GetAccountHistories())

		// Me (Account)
		api.GET("/v1/me", accountHandler.Me())

		// Logout - 需要登录验证
		api.POST("/v1/auth/logout", accountHandler.Logout())

		// org-infos 路由移除，统一使用 /v1/groups 系列
	}

	// system 工具接口：最近日志
	system := r.mux.Group("/system")
	{
		system.GET("/logs/latest", logsHandler.GetLatestLogs())

		system.GET("/logs/unified", logsHandler.GetUnifiedLogs())

		// 新增：分页日志接口，支持按需加载
		system.GET("/logs/paginated", logsHandler.GetPaginatedLogs())
	}
}
