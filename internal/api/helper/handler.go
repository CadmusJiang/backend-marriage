package helper

import (
	"github.com/xinliangnote/go-gin-api/internal/pkg/core"
	"github.com/xinliangnote/go-gin-api/internal/repository/mysql"
	"github.com/xinliangnote/go-gin-api/internal/repository/redis"
	"github.com/xinliangnote/go-gin-api/internal/services/authorized"

	"go.uber.org/zap"
)

var _ Handler = (*handler)(nil)

type Handler interface {
	i()

	// Md5 加密
	// @Tags Helper
	// @Router /helper/md5/{str} [get]
	Md5() core.HandlerFunc

	// Sign 签名
	// @Tags Helper
	// @Router /helper/sign [post]
	Sign() core.HandlerFunc

	// CheckDatabase 检查数据库表结构
	// @Tags Helper
	// @Router /api/v1/check-db [get]
	CheckDatabase() core.HandlerFunc

	// GetLogs 获取日志列表
	// @Tags API.helper
	// @Router /api/v1/logs [get]
	GetLogs() core.HandlerFunc

	// GetLogsRealtime 获取实时日志
	// @Tags API.helper
	// @Router /api/v1/logs/realtime [get]
	GetLogsRealtime() core.HandlerFunc
}

type handler struct {
	logger            *zap.Logger
	db                mysql.Repo
	authorizedService authorized.Service
}

func New(logger *zap.Logger, db mysql.Repo, cache redis.Repo) Handler {
	return &handler{
		logger:            logger,
		db:                db,
		authorizedService: authorized.New(db, cache),
	}
}

func (h *handler) i() {}
