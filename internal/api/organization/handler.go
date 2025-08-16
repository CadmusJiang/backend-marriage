package organization

import (
	"github.com/xinliangnote/go-gin-api/internal/pkg/core"
	"github.com/xinliangnote/go-gin-api/internal/repository/mysql"
	"github.com/xinliangnote/go-gin-api/internal/repository/redis"
	orgsvc "github.com/xinliangnote/go-gin-api/internal/services/organization"
	"go.uber.org/zap"
)

var _ Handler = (*handler)(nil)

type Handler interface {
	i()

	// groups
	GetGroups() core.HandlerFunc
	CreateGroup() core.HandlerFunc
	GetOrgInfoDetail() core.HandlerFunc
	UpdateGroup() core.HandlerFunc
	GetGroupHistory() core.HandlerFunc

	// teams
	ListTeams() core.HandlerFunc
	CreateTeam() core.HandlerFunc
	GetTeam() core.HandlerFunc
	UpdateTeam() core.HandlerFunc
	GetTeamHistory() core.HandlerFunc

	// team members
	ListTeamMembers() core.HandlerFunc
	AddTeamMember() core.HandlerFunc
	RemoveTeamMember() core.HandlerFunc
	UpdateTeamMemberRole() core.HandlerFunc

	// misc
	ListUnassignedAccounts() core.HandlerFunc
}

type handler struct {
	logger     *zap.Logger
	db         mysql.Repo
	cache      redis.Repo
	orgService orgsvc.Service
}

func New(logger *zap.Logger, db mysql.Repo, cache redis.Repo) Handler {
	return &handler{logger: logger, db: db, cache: cache, orgService: orgsvc.New(db)}
}

func (h *handler) i() {}

// 临时实现，返回空结果
func (h *handler) GetGroups() core.HandlerFunc {
	return func(c core.Context) {
		c.Payload(map[string]interface{}{
			"success": true,
			"data":    []interface{}{},
			"meta": map[string]interface{}{
				"total":    0,
				"pageSize": 10,
				"current":  1,
			},
		})
	}
}

func (h *handler) CreateGroup() core.HandlerFunc {
	return func(c core.Context) {
		c.Payload(map[string]interface{}{
			"success": true,
			"data":    map[string]interface{}{"id": "1"},
		})
	}
}

func (h *handler) GetOrgInfoDetail() core.HandlerFunc {
	return func(c core.Context) {
		c.Payload(map[string]interface{}{
			"success": true,
			"data":    map[string]interface{}{},
		})
	}
}

func (h *handler) UpdateGroup() core.HandlerFunc {
	return func(c core.Context) {
		c.Payload(map[string]interface{}{
			"success": true,
		})
	}
}

func (h *handler) GetGroupHistory() core.HandlerFunc {
	return func(c core.Context) {
		c.Payload(map[string]interface{}{
			"success": true,
			"data":    []interface{}{},
			"meta": map[string]interface{}{
				"total":    0,
				"pageSize": 10,
				"current":  1,
			},
		})
	}
}

func (h *handler) ListTeams() core.HandlerFunc {
	return func(c core.Context) {
		c.Payload(map[string]interface{}{
			"success": true,
			"data":    []interface{}{},
			"meta": map[string]interface{}{
				"total":    0,
				"pageSize": 10,
				"current":  1,
			},
		})
	}
}

func (h *handler) CreateTeam() core.HandlerFunc {
	return func(c core.Context) {
		c.Payload(map[string]interface{}{
			"success": true,
			"data":    map[string]interface{}{"id": "1"},
		})
	}
}

func (h *handler) GetTeam() core.HandlerFunc {
	return func(c core.Context) {
		c.Payload(map[string]interface{}{
			"success": true,
			"data":    map[string]interface{}{},
		})
	}
}

func (h *handler) UpdateTeam() core.HandlerFunc {
	return func(c core.Context) {
		c.Payload(map[string]interface{}{
			"success": true,
		})
	}
}

func (h *handler) GetTeamHistory() core.HandlerFunc {
	return func(c core.Context) {
		c.Payload(map[string]interface{}{
			"success": true,
			"data":    []interface{}{},
			"meta": map[string]interface{}{
				"total":    0,
				"pageSize": 10,
				"current":  1,
			},
		})
	}
}

func (h *handler) ListTeamMembers() core.HandlerFunc {
	return func(c core.Context) {
		c.Payload(map[string]interface{}{
			"success": true,
			"data":    []interface{}{},
			"meta": map[string]interface{}{
				"total":    0,
				"pageSize": 10,
				"current":  1,
			},
		})
	}
}

func (h *handler) AddTeamMember() core.HandlerFunc {
	return func(c core.Context) {
		c.Payload(map[string]interface{}{
			"success": true,
		})
	}
}

func (h *handler) RemoveTeamMember() core.HandlerFunc {
	return func(c core.Context) {
		c.Payload(map[string]interface{}{
			"success": true,
		})
	}
}

func (h *handler) UpdateTeamMemberRole() core.HandlerFunc {
	return func(c core.Context) {
		c.Payload(map[string]interface{}{
			"success": true,
		})
	}
}

func (h *handler) ListUnassignedAccounts() core.HandlerFunc {
	return func(c core.Context) {
		c.Payload(map[string]interface{}{
			"success": true,
			"data":    []interface{}{},
			"meta": map[string]interface{}{
				"total":    0,
				"pageSize": 10,
				"current":  1,
			},
		})
	}
}
