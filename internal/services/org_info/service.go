package org_info

import (
	"github.com/xinliangnote/go-gin-api/internal/pkg/core"
	"github.com/xinliangnote/go-gin-api/internal/repository/mysql"
	"github.com/xinliangnote/go-gin-api/internal/repository/mysql/org_info"
	"github.com/xinliangnote/go-gin-api/internal/repository/redis"
)

var _ Service = (*service)(nil)

// CreateOrgInfoData 创建组织数据
type CreateOrgInfoData struct {
	OrgName  string `json:"orgName" binding:"required"`
	OrgType  string `json:"orgType" binding:"required"`
	OrgPath  string `json:"orgPath"`
	OrgLevel int32  `json:"orgLevel"`
	MaxCnt   int32  `json:"maxCnt"`
}

// UpdateOrgInfoData 更新组织数据
type UpdateOrgInfoData struct {
	OrgName  string `json:"orgName"`
	OrgType  string `json:"orgType"`
	OrgPath  string `json:"orgPath"`
	OrgLevel int32  `json:"orgLevel"`
	MaxCnt   int32  `json:"maxCnt"`
}

// SearchOrgInfoData 搜索组织数据
type SearchOrgInfoData struct {
	OrgName  string `form:"orgName"`
	OrgType  string `form:"orgType"`
	OrgLevel int32  `form:"orgLevel"`
	Current  int    `form:"current,default=1"`
	PageSize int    `form:"pageSize,default=10"`
}

type Service interface {
	i()

	// 组织管理
	Create(ctx core.Context, orgData *CreateOrgInfoData) (id int32, err error)
	PageList(ctx core.Context, searchData *SearchOrgInfoData) (listData []*org_info.OrgInfo, err error)
	PageListCount(ctx core.Context, searchData *SearchOrgInfoData) (total int64, err error)
	Detail(ctx core.Context, orgId string) (info *org_info.OrgInfo, err error)
	Update(ctx core.Context, orgId string, updateData *UpdateOrgInfoData) (err error)
	Delete(ctx core.Context, orgId string) (err error)

	// 组织层级管理
	GetChildren(ctx core.Context, parentOrgId string) (children []*org_info.OrgInfo, err error)
	GetParent(ctx core.Context, orgId string) (parent *org_info.OrgInfo, err error)

	// 成员统计
	UpdateMemberCount(ctx core.Context, orgId string, count int32) (err error)
}

type service struct {
	db    mysql.Repo
	cache redis.Repo
}

func New(db mysql.Repo, cache redis.Repo) Service {
	return &service{
		db:    db,
		cache: cache,
	}
}

func (s *service) i() {}
