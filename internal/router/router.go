package router

import (
	"fmt"

	"github.com/xinliangnote/go-gin-api/configs"
	"github.com/xinliangnote/go-gin-api/internal/alert"
	"github.com/xinliangnote/go-gin-api/internal/metrics"
	"github.com/xinliangnote/go-gin-api/internal/pkg/core"
	"github.com/xinliangnote/go-gin-api/internal/repository/cron"
	"github.com/xinliangnote/go-gin-api/internal/repository/mysql"
	"github.com/xinliangnote/go-gin-api/internal/repository/redis"
	"github.com/xinliangnote/go-gin-api/internal/router/interceptor"
	"github.com/xinliangnote/go-gin-api/pkg/errors"
	"github.com/xinliangnote/go-gin-api/pkg/file"

	"go.uber.org/zap"
)

type resource struct {
	mux          core.Mux
	logger       *zap.Logger
	db           mysql.Repo
	cache        redis.Repo
	interceptors interceptor.Interceptor
	cronServer   cron.Server
}

type Server struct {
	Mux        core.Mux
	Db         mysql.Repo
	Cache      redis.Repo
	CronServer cron.Server
}

func NewHTTPServer(logger *zap.Logger, cronLogger *zap.Logger) (*Server, error) {
	if logger == nil {
		return nil, errors.New("logger required")
	}

	r := new(resource)
	r.logger = logger

	openBrowserUri := configs.ProjectDomain + configs.ProjectPort

	_, ok := file.IsExists(configs.ProjectInstallMark)
	if !ok { // 未安装
		openBrowserUri += "/install"
	} else { // 已安装

		// 初始化 DB
		logger.Info("正在初始化数据库连接...")
		dbRepo, err := mysql.New()
		if err != nil {
			logger.Error("数据库连接失败",
				zap.String("error", err.Error()),
				zap.String("config", fmt.Sprintf("Read: %s@%s/%s, Write: %s@%s/%s",
					configs.Get().MySQL.Read.User, configs.Get().MySQL.Read.Addr, configs.Get().MySQL.Read.Name,
					configs.Get().MySQL.Write.User, configs.Get().MySQL.Write.Addr, configs.Get().MySQL.Write.Name)),
			)
			return nil, fmt.Errorf("数据库连接失败: %v", err)
		}

		// 测试数据库连接
		logger.Info("正在测试数据库连接...")
		if err := testDatabaseConnection(dbRepo); err != nil {
			logger.Error("数据库连接测试失败", zap.Error(err))
			return nil, fmt.Errorf("数据库连接测试失败: %v", err)
		}

		logger.Info("数据库连接成功")
		r.db = dbRepo

		// 初始化 Cache
		logger.Info("正在初始化Redis缓存...")
		cacheRepo, err := redis.New()
		if err != nil {
			logger.Error("Redis连接失败", zap.Error(err))
			return nil, fmt.Errorf("Redis连接失败: %v", err)
		}
		logger.Info("Redis连接成功")
		r.cache = cacheRepo

		// 初始化 CRON Server
		logger.Info("正在初始化定时任务服务...")
		cronServer, err := cron.New(cronLogger, dbRepo, cacheRepo)
		if err != nil {
			logger.Error("定时任务服务初始化失败", zap.Error(err))
			return nil, fmt.Errorf("定时任务服务初始化失败: %v", err)
		}
		cronServer.Start()
		logger.Info("定时任务服务启动成功")
		r.cronServer = cronServer
	}

	mux, err := core.New(logger,
		// core.WithEnableOpenBrowser(openBrowserUri), // 注释掉自动打开浏览器功能
		core.WithEnableCors(),
		core.WithEnableRate(),
		core.WithAlertNotify(alert.NotifyHandler(logger)),
		core.WithRecordMetrics(metrics.RecordHandler(logger)),
	)

	if err != nil {
		panic(err)
	}

	r.mux = mux
	r.interceptors = interceptor.New(logger, r.cache, r.db)

	// 设置 Render 路由
	setRenderRouter(r)

	// 设置 API 路由
	setApiRouter(r)

	// 设置 GraphQL 路由
	setGraphQLRouter(r)

	// 设置 Socket 路由
	setSocketRouter(r)

	s := new(Server)
	s.Mux = mux
	s.Db = r.db
	s.Cache = r.cache
	s.CronServer = r.cronServer

	return s, nil
}

// testDatabaseConnection 测试数据库连接
func testDatabaseConnection(db mysql.Repo) error {
	// 测试读连接
	if err := db.GetDbR().Raw("SELECT 1").Error; err != nil {
		return fmt.Errorf("读数据库连接测试失败: %v", err)
	}

	// 测试写连接
	if err := db.GetDbW().Raw("SELECT 1").Error; err != nil {
		return fmt.Errorf("写数据库连接测试失败: %v", err)
	}

	return nil
}
