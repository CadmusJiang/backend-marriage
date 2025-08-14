package router

import (
	"context"
	"fmt"

	goRedis "github.com/go-redis/redis/v7"
	"github.com/xinliangnote/go-gin-api/configs"
	"github.com/xinliangnote/go-gin-api/internal/alert"
	obrelay "github.com/xinliangnote/go-gin-api/internal/consumers/outbox"
	"github.com/xinliangnote/go-gin-api/internal/metrics"
	"github.com/xinliangnote/go-gin-api/internal/pkg/core"
	"github.com/xinliangnote/go-gin-api/internal/proposal/tablesqls"
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
}

type Server struct {
	Mux   core.Mux
	Db    mysql.Repo
	Cache redis.Repo
}

func NewHTTPServer(logger *zap.Logger, forceReseed bool) (*Server, error) {
	if logger == nil {
		return nil, errors.New("logger required")
	}

	r := new(resource)
	r.logger = logger

	_, ok := file.IsExists(configs.ProjectInstallMark)
	if ok { // 已安装

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

		// 启动即确保表结构并插入 mock 数据（逐表输出进度，便于中断场景下确认进度）
		logger.Info("正在检查并创建缺失的数据表...")
		if err := ensureTables(dbRepo, logger); err != nil {
			logger.Error("创建/校验数据表失败", zap.Error(err))
			return nil, fmt.Errorf("创建/校验数据表失败: %v", err)
		}
		logger.Info("数据表检查完成")

		if forceReseed {
			logger.Info("正在强制重置并插入 mock 数据...")
			if err := reinsertMockData(dbRepo, logger); err != nil {
				logger.Error("插入 mock 数据失败", zap.Error(err))
				return nil, fmt.Errorf("插入 mock 数据失败: %v", err)
			}
			logger.Info("mock 数据插入完成")
		}

		// 初始化 Cache
		logger.Info("正在初始化Redis缓存...")
		cacheRepo, err := redis.New()
		if err != nil {
			logger.Error("Redis连接失败", zap.Error(err))
			return nil, fmt.Errorf("Redis连接失败: %v", err)
		}
		logger.Info("Redis连接成功")
		r.cache = cacheRepo

		// 定时任务已移除
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

	// 页面渲染路由已移除

	// 设置 API 路由
	setApiRouter(r)

	// GraphQL 路由已移除

	// WebSocket 路由已移除

	// Ensure required Redis Streams exist
	ensureRedisStreams(logger)

	s := new(Server)
	s.Mux = mux
	s.Db = r.db
	s.Cache = r.cache

	// Start outbox relay after resources are ready
	obrelay.StartRelay(context.Background(), logger, r.db)

	return s, nil
}

// ensureRedisStreams ensures required streams exist; creates them if missing.
func ensureRedisStreams(logger *zap.Logger) {
	cfg := configs.Get().Redis
	rdb := goRedis.NewClient(&goRedis.Options{
		Addr:         cfg.Addr,
		Password:     cfg.Pass,
		DB:           cfg.Db,
		MaxRetries:   cfg.MaxRetries,
		PoolSize:     cfg.PoolSize,
		MinIdleConns: cfg.MinIdleConns,
	})
	defer func() { _ = rdb.Close() }()

	const streamCustomer = "stream.customer.record"
	const bootstrapGroup = "bootstrap"

	if err := rdb.XGroupCreateMkStream(streamCustomer, bootstrapGroup, "$").Err(); err != nil {
		if err.Error() != "BUSYGROUP Consumer Group name already exists" {
			logger.Warn("ensure stream/group", zap.String("stream", streamCustomer), zap.Error(err))
		} else {
			logger.Debug("stream/group exists", zap.String("stream", streamCustomer))
		}
	} else {
		logger.Info("created redis stream+group", zap.String("stream", streamCustomer), zap.String("group", bootstrapGroup))
	}
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

// ensureTables 创建核心业务相关的数据表（若不存在）
func ensureTables(db mysql.Repo, logger *zap.Logger) error {
	type table struct {
		name string
		sql  string
	}
	tables := []table{
		{name: "org", sql: tablesqls.CreateOrgTableSql()},
		{name: "org_history", sql: tablesqls.CreateOrgHistoryTableSql()},
		{name: "account", sql: tablesqls.CreateAccountTableSql()},
		{name: "account_history", sql: tablesqls.CreateAccountHistoryTableSql()},
		{name: "account_org_relation", sql: tablesqls.CreateAccountOrgRelationTableSql()},
		{name: "customer_authorization_record", sql: tablesqls.CreateCustomerAuthorizationRecordTableSql()},
		{name: "customer_authorization_record_history", sql: tablesqls.CreateCustomerAuthorizationRecordHistoryTableSql()},
		{name: "cooperation_store", sql: tablesqls.CreateCooperationStoreTableSql()},
		{name: "outbox_events", sql: tablesqls.CreateOutboxTableSql()},
	}

	for _, t := range tables {
		logger.Info("检查/创建表", zap.String("table", t.name))
		// 是否存在
		var exists bool
		if err := db.GetDbW().Raw("SELECT COUNT(*) > 0 FROM information_schema.tables WHERE table_schema = DATABASE() AND table_name = ?", t.name).Scan(&exists).Error; err != nil {
			return fmt.Errorf("检查表 %s 是否存在失败: %v", t.name, err)
		}
		if !exists {
			logger.Info("表不存在，正在创建...", zap.String("table", t.name))
			if err := db.GetDbW().Exec(t.sql).Error; err != nil {
				return fmt.Errorf("创建表 %s 失败: %v", t.name, err)
			}
			logger.Info("创建完成", zap.String("table", t.name))
		} else {
			logger.Info("表已存在，跳过创建", zap.String("table", t.name))
		}
	}

	return nil
}

// reinsertMockData 重新插入mock数据
func reinsertMockData(db mysql.Repo, logger *zap.Logger) error {
	// 定义需要重新插入的表和对应的SQL（按依赖顺序执行）
	type tdata struct {
		name   string
		insert string
	}
	items := []tdata{
		{name: "org", insert: tablesqls.CreateOrgTableDataSql()},
		{name: "org_history", insert: tablesqls.CreateOrgHistoryTableDataSql()},
		{name: "account", insert: tablesqls.CreateAccountTableDataSql()},
		{name: "account_history", insert: tablesqls.CreateAccountHistoryTableDataSql()},
		{name: "account_org_relation", insert: tablesqls.CreateAccountOrgRelationTableDataSql()},
		{name: "customer_authorization_record", insert: tablesqls.CreateCustomerAuthorizationRecordTableDataSql()},
		{name: "customer_authorization_record_history", insert: tablesqls.CreateCustomerAuthorizationRecordHistoryTableDataSql()},
		{name: "cooperation_store", insert: tablesqls.CreateCooperationStoreTableDataSql()},
	}

	for _, it := range items {
		logger.Info("清空表数据", zap.String("table", it.name))
		res := db.GetDbW().Exec(fmt.Sprintf("DELETE FROM %s", it.name))
		if err := res.Error; err != nil {
			return fmt.Errorf("清空表 %s 失败: %v", it.name, err)
		}
		logger.Info("已清空", zap.String("table", it.name), zap.Int64("rows_affected", res.RowsAffected))
		if it.insert == "" {
			logger.Info("无插入SQL，跳过", zap.String("table", it.name))
			continue
		}
		logger.Info("插入 mock 数据", zap.String("table", it.name))
		res = db.GetDbW().Exec(it.insert)
		if err := res.Error; err != nil {
			return fmt.Errorf("插入表 %s 的mock数据失败: %v", it.name, err)
		}
		logger.Info("插入完成", zap.String("table", it.name), zap.Int64("rows_affected", res.RowsAffected))
	}

	return nil
}
