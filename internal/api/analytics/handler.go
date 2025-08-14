package analytics

import (
	"github.com/xinliangnote/go-gin-api/internal/repository/mysql"
	"go.uber.org/zap"
)

type handler struct {
	logger *zap.Logger
	db     mysql.Repo
	cache  mysql.Repo
}

func New(logger *zap.Logger, db mysql.Repo, cache mysql.Repo) *handler {
	return &handler{
		logger: logger,
		db:     db,
		cache:  cache,
	}
}

// internal helpers used by analytics handlers are placed here to share across files
