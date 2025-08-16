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

	// ==================== Groups 相关接口 ====================
	// GetGroups 获取组列表
	// @Tags Organization
	// @Router /api/v1/organizations/groups [get]
	GetGroups() core.HandlerFunc

	// CreateGroup 创建组
	// @Tags Organization
	// @Router /api/v1/organizations/groups [post]
	CreateGroup() core.HandlerFunc

	// GetOrgInfoDetail 获取组织详情
	// @Tags Organization
	// @Router /api/v1/organizations/{orgId} [get]
	GetOrgInfoDetail() core.HandlerFunc

	// UpdateGroup 更新组
	// @Tags Organization
	// @Router /api/v1/organizations/groups/{groupId} [put]
	UpdateGroup() core.HandlerFunc

	// GetGroupHistory 获取组历史
	// @Tags Organization
	// @Router /api/v1/organizations/groups/{groupId}/history [get]
	GetGroupHistory() core.HandlerFunc

	// ==================== Teams 相关接口 ====================
	// ListTeams 获取团队列表
	// @Tags Organization
	// @Router /api/v1/organizations/teams [get]
	ListTeams() core.HandlerFunc

	// CreateTeam 创建团队
	// @Tags Organization
	// @Router /api/v1/organizations/teams [post]
	CreateTeam() core.HandlerFunc

	// GetTeam 获取团队详情
	// @Tags Organization
	// @Router /api/v1/organizations/teams/{teamId} [get]
	GetTeam() core.HandlerFunc

	// UpdateTeam 更新团队
	// @Tags Organization
	// @Router /api/v1/organizations/teams/{teamId} [put]
	UpdateTeam() core.HandlerFunc

	// GetTeamHistory 获取团队历史
	// @Tags Organization
	// @Router /api/v1/organizations/teams/{teamId}/history [get]
	GetTeamHistory() core.HandlerFunc

	// ==================== Team Members 相关接口 ====================
	// ListTeamMembers 获取团队成员列表
	// @Tags Organization
	// @Router /api/v1/organizations/{orgId}/members [get]
	ListTeamMembers() core.HandlerFunc

	// AddTeamMember 添加团队成员
	// @Tags Organization
	// @Router /api/v1/organizations/{orgId}/members [post]
	AddTeamMember() core.HandlerFunc

	// RemoveTeamMember 移除团队成员
	// @Tags Organization
	// @Router /api/v1/organizations/{orgId}/members/{accountId} [delete]
	RemoveTeamMember() core.HandlerFunc

	// UpdateTeamMemberRole 更新成员角色
	// @Tags Organization
	// @Router /api/v1/organizations/{orgId}/members/{accountId}/role [put]
	UpdateTeamMemberRole() core.HandlerFunc

	// UpdateTeamMember 更新团队成员信息（包括角色等）
	// @Tags Organization
	// @Router /api/v1/organizations/{orgId}/members/{accountId} [patch]
	UpdateTeamMember() core.HandlerFunc

	// ==================== 其他接口 ====================
	// ListUnassignedAccounts 获取未分配账户列表
	// @Tags Organization
	// @Router /api/v1/organizations/unassigned-accounts [get]
	ListUnassignedAccounts() core.HandlerFunc
}

type handler struct {
	logger     *zap.Logger
	db         mysql.Repo
	cache      redis.Repo
	orgService orgsvc.Service
}

func New(logger *zap.Logger, db mysql.Repo, cache redis.Repo) Handler {
	return &handler{
		logger:     logger,
		db:         db,
		cache:      cache,
		orgService: orgsvc.New(db),
	}
}

func (h *handler) i() {}
