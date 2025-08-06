package org_info

import (
	"github.com/xinliangnote/go-gin-api/internal/pkg/core"
	"github.com/xinliangnote/go-gin-api/internal/repository/mysql"
	"github.com/xinliangnote/go-gin-api/internal/repository/redis"
	"github.com/xinliangnote/go-gin-api/internal/services/org_info"

	"go.uber.org/zap"
)

var _ Handler = (*handler)(nil)

type Handler interface {
	i()

	// GetGroups 获取所有groups信息
	// @Tags API.org_info
	// @Router /api/v1/groups [get]
	GetGroups() core.HandlerFunc

	// GetTeams 获取所有teams信息
	// @Tags API.org_info
	// @Router /api/v1/teams [get]
	GetTeams() core.HandlerFunc

	// GetOrgInfoList 获取组织信息列表
	// @Tags API.org_info
	// @Router /api/v1/org-infos [get]
	GetOrgInfoList() core.HandlerFunc

	// CreateOrgInfo 创建组织信息
	// @Tags API.org_info
	// @Router /api/v1/org-infos [post]
	CreateOrgInfo() core.HandlerFunc

	// GetOrgInfoDetail 获取组织信息详情
	// @Tags API.org_info
	// @Router /api/v1/org-infos/{orgId} [get]
	GetOrgInfoDetail() core.HandlerFunc

	// UpdateOrgInfo 更新组织信息
	// @Tags API.org_info
	// @Router /api/v1/org-infos/{orgId} [put]
	UpdateOrgInfo() core.HandlerFunc

	// DeleteOrgInfo 删除组织信息
	// @Tags API.org_info
	// @Router /api/v1/org-infos/{orgId} [delete]
	DeleteOrgInfo() core.HandlerFunc

	// GetOrgInfoChildren 获取子组织
	// @Tags API.org_info
	// @Router /api/v1/org-infos/{orgId}/children [get]
	GetOrgInfoChildren() core.HandlerFunc

	// GetOrgInfoParent 获取父组织
	// @Tags API.org_info
	// @Router /api/v1/org-infos/{orgId}/parent [get]
	GetOrgInfoParent() core.HandlerFunc
}

type handler struct {
	logger         *zap.Logger
	cache          redis.Repo
	db             mysql.Repo
	orgInfoService org_info.Service
}

func New(logger *zap.Logger, db mysql.Repo, cache redis.Repo) Handler {
	return &handler{
		logger:         logger,
		cache:          cache,
		db:             db,
		orgInfoService: org_info.New(db, cache),
	}
}

func (h *handler) i() {}
