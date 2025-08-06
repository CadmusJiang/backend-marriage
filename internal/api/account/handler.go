package account

import (
	"github.com/xinliangnote/go-gin-api/internal/pkg/core"
	"github.com/xinliangnote/go-gin-api/internal/repository/mysql"
	"github.com/xinliangnote/go-gin-api/internal/repository/redis"
	"github.com/xinliangnote/go-gin-api/internal/services/account"

	"go.uber.org/zap"
)

var _ Handler = (*handler)(nil)

type Handler interface {
	i()

	// GetAccountList 获取账户列表
	// @Tags API.account
	// @Router /api/v1/accounts [get]
	GetAccountList() core.HandlerFunc

	// CreateAccount 创建账户
	// @Tags API.account
	// @Router /api/v1/accounts [post]
	CreateAccount() core.HandlerFunc

	// GetAccountDetail 获取账户详情
	// @Tags API.account
	// @Router /api/v1/accounts/{accountId} [get]
	GetAccountDetail() core.HandlerFunc

	// UpdateAccount 更新账户
	// @Tags API.account
	// @Router /api/v1/accounts/{accountId} [put]
	UpdateAccount() core.HandlerFunc

	// GetAccountHistories 获取账户历史记录
	// @Tags API.account
	// @Router /api/v1/account-histories [get]
	GetAccountHistories() core.HandlerFunc

	// Login 用户登录
	// @Tags API.account
	// @Router /api/v1/auth/login [post]
	Login() core.HandlerFunc

	// Logout 退出登录
	// @Tags API.account
	// @Router /api/v1/auth/logout [post]
	Logout() core.HandlerFunc
}

type handler struct {
	logger         *zap.Logger
	cache          redis.Repo
	db             mysql.Repo
	accountService account.Service
}

func New(logger *zap.Logger, db mysql.Repo, cache redis.Repo) Handler {
	return &handler{
		logger:         logger,
		cache:          cache,
		db:             db,
		accountService: account.New(db, cache),
	}
}

func (h *handler) i() {}
