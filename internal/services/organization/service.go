package organization

import (
	"github.com/xinliangnote/go-gin-api/internal/pkg/core"
	"github.com/xinliangnote/go-gin-api/internal/repository/mysql"
)

type Service interface {
	i()

	// groups
	ListGroups(ctx core.Context, current, pageSize int, keyword string) (list interface{}, total int64, err error)
	CreateGroup(ctx core.Context, payload *CreateGroupPayload) (id uint32, err error)
	GetGroup(ctx core.Context, orgId string) (info interface{}, err error)
	UpdateGroup(ctx core.Context, orgId string, payload *UpdateGroupPayload) (newVersion int, err error)
	ListGroupHistory(ctx core.Context, orgId string, current, pageSize int) (items interface{}, total int64, err error)

	// teams
	ListTeams(ctx core.Context, current, pageSize int, belongGroupId string, keyword string) (list interface{}, total int64, err error)
	CreateTeam(ctx core.Context, payload *CreateTeamPayload) (id uint32, err error)
	GetTeam(ctx core.Context, teamId string) (info interface{}, err error)
	UpdateTeam(ctx core.Context, teamId string, payload *UpdateTeamPayload) (newVersion int, err error)
	ListTeamHistory(ctx core.Context, teamId string, current, pageSize int) (items interface{}, total int64, err error)

	// members
	ListTeamMembers(ctx core.Context, teamId string, current, pageSize int) (list interface{}, total int64, err error)
	AddTeamMember(ctx core.Context, teamId string, accountId string, roleType string) (id string, err error)
	RemoveTeamMember(ctx core.Context, teamId string, accountId string) (err error)
	UpdateTeamMemberRole(ctx core.Context, teamId string, accountId string, roleType string) (err error)

	// misc
	ListUnassignedAccounts(ctx core.Context, current, pageSize int, keyword string) (list interface{}, total int64, err error)
}

type service struct {
	db mysql.Repo
}

func New(db mysql.Repo) Service {
	return &service{db: db}
}

func (s *service) i() {}

// Payloads
type CreateGroupPayload struct {
	Username string `json:"username"`
	Nickname string `json:"nickname"`
}

type UpdateGroupPayload struct {
	Nickname string `json:"nickname"`
	Status   int32  `json:"status"`
	Version  int    `json:"version"`
}

type CreateTeamPayload struct {
	BelongGroupId string `json:"belongGroupId"`
	Username      string `json:"username"`
	Nickname      string `json:"nickname"`
}

type UpdateTeamPayload struct {
	Nickname string `json:"nickname"`
	Status   int32  `json:"status"`
	Version  int    `json:"version"`
}
