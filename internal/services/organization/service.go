package orgsvc

import (
	"github.com/xinliangnote/go-gin-api/internal/authz"
	"github.com/xinliangnote/go-gin-api/internal/pkg/core"
	"github.com/xinliangnote/go-gin-api/internal/proposal"
)

type Service interface {
	i()

	// 组管理
	ListGroups(ctx Context, current, pageSize int, keyword string, scope *authz.AccessScope) ([]map[string]interface{}, int64, error)
	GetGroup(ctx Context, id string) (map[string]interface{}, error)
	CreateGroup(ctx Context, payload *CreateGroupPayload) (uint32, error)
	UpdateGroup(ctx Context, id string, payload *UpdateGroupPayload) (uint32, error)
	DeleteGroup(ctx Context, id string) error

	// 团队管理
	ListTeams(ctx Context, belongGroupId string, current, pageSize int, keyword string, scope *authz.AccessScope) ([]map[string]interface{}, int64, error)
	GetTeam(ctx Context, id string) (map[string]interface{}, error)
	CreateTeam(ctx Context, payload *CreateTeamPayload) (uint32, error)
	UpdateTeam(ctx Context, id string, payload *UpdateTeamPayload) (uint32, error)
	DeleteTeam(ctx Context, id string) error

	// 成员管理
	ListMembers(ctx Context, orgId string, current, pageSize int, keyword string) ([]map[string]interface{}, int64, error)
	AddMember(ctx Context, orgId string, accountId string) error
	RemoveMember(ctx Context, orgId string, accountId string) error
	UpdateMemberRole(ctx Context, orgId string, accountId string, roleType string) error
}

type Context interface {
	RequestContext() core.StdContext
	SessionUserInfo() proposal.SessionUserInfo
	// 直接使用core.Context的GetTraceID方法
	GetTraceID() string
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
	Name    string                 `json:"name,omitempty"`
	Address *string                `json:"address,omitempty"`
	Extra   map[string]interface{} `json:"extra,omitempty"`
	Status  string                 `json:"status,omitempty"`
	Version int32                  `json:"version"`
}

type CreateTeamPayload struct {
	BelongGroupId uint32 `json:"belongGroupId"`
	Username      string `json:"username"`
	Name          string `json:"name"`
}

type UpdateTeamPayload struct {
	Name    string                 `json:"name,omitempty"`
	Address *string                `json:"address,omitempty"`
	Extra   map[string]interface{} `json:"extra,omitempty"`
	Status  string                 `json:"status,omitempty"`
	Version int32                  `json:"version"`
}

// MemberOperationError 成员操作错误
type MemberOperationError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details string `json:"details"`
}

func (e *MemberOperationError) Error() string {
	return e.Message
}

// NewMemberOperationError 创建成员操作错误
func NewMemberOperationError(code, message, details string) *MemberOperationError {
	return &MemberOperationError{
		Code:    code,
		Message: message,
		Details: details,
	}
}
