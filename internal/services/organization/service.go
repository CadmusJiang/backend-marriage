package orgsvc

import (
	"github.com/xinliangnote/go-gin-api/internal/pkg/core"
	"github.com/xinliangnote/go-gin-api/internal/proposal"
)

type Service interface {
	// Group operations
	ListGroups(ctx Context, current, pageSize int, keyword string) ([]map[string]interface{}, int64, error)
	GetGroup(ctx Context, id uint32) (map[string]interface{}, error)
	CreateGroup(ctx Context, payload *CreateGroupPayload) (uint32, error)
	UpdateGroup(ctx Context, id uint32, payload *UpdateGroupPayload) (uint32, error)
	DeleteGroup(ctx Context, id uint32) error

	// Team operations
	ListTeams(ctx Context, belongGroupId uint32, current, pageSize int, keyword string) ([]map[string]interface{}, int64, error)
	GetTeam(ctx Context, id uint32) (map[string]interface{}, error)
	CreateTeam(ctx Context, payload *CreateTeamPayload) (uint32, error)
	UpdateTeam(ctx Context, id uint32, payload *UpdateTeamPayload) (uint32, error)
	DeleteTeam(ctx Context, id uint32) error

	// Member operations
	ListMembers(ctx Context, orgId uint32, current, pageSize int, keyword string) ([]map[string]interface{}, int64, error)
	AddMember(ctx Context, orgId uint32, accountId uint32) error
	RemoveMember(ctx Context, orgId uint32, accountId uint32) error
}

type Context interface {
	RequestContext() core.StdContext
	SessionUserInfo() proposal.SessionUserInfo
}

type SessionUserInfo struct {
	UserID   int32
	UserName string
	RoleType string
}

type CreateGroupPayload struct {
	Username string `json:"username"`
	Name     string `json:"name"`
}

type UpdateGroupPayload struct {
	Name    string `json:"name"`
	Status  int32  `json:"status"`
	Version int32  `json:"version"`
}

type CreateTeamPayload struct {
	BelongGroupId uint32 `json:"belongGroupId"`
	Username      string `json:"username"`
	Name          string `json:"name"`
}

type UpdateTeamPayload struct {
	Name    string `json:"name"`
	Status  int32  `json:"status"`
	Version int32  `json:"version"`
}
