package router

import (
	"github.com/xinliangnote/go-gin-api/internal/pkg/core"
	"github.com/xinliangnote/go-gin-api/internal/render/index"
	"github.com/xinliangnote/go-gin-api/internal/render/install"
	"github.com/xinliangnote/go-gin-api/internal/render/logs"
)

func setRenderRouter(r *resource) {
	renderIndex := index.New(r.logger, r.db, r.cache)
	renderInstall := install.New(r.logger)
	renderLogs := logs.New(r.logger, r.db, r.cache)

	// 无需记录日志，无需 RBAC 权限验证
	notRBAC := r.mux.Group("", core.DisableTraceLog, core.DisableRecordMetrics)
	{
		// 首页
		notRBAC.GET("", renderIndex.Index())

		// 安装页面
		notRBAC.GET("/install", renderInstall.View())
		notRBAC.POST("/install/execute", renderInstall.Execute())

		// 日志查看页面
		notRBAC.GET("/logs", renderLogs.Viewer())
	}
}
