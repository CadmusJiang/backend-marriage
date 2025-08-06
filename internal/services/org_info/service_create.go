package org_info

import (
	"fmt"
	"time"

	"github.com/xinliangnote/go-gin-api/internal/pkg/core"
	"github.com/xinliangnote/go-gin-api/internal/repository/mysql"
	"github.com/xinliangnote/go-gin-api/internal/repository/mysql/org_info"
)

// Create 创建组织信息
func (s *service) Create(ctx core.Context, orgData *CreateOrgInfoData) (id int32, err error) {
	// 检查组织名称是否已存在
	orgQueryBuilder := org_info.NewQueryBuilder()
	orgQueryBuilder.WhereOrgName(mysql.EqualPredicate, orgData.OrgName)

	existingOrg, err := orgQueryBuilder.QueryOne(s.db.GetDbR())
	if err != nil {
		return 0, fmt.Errorf("查询组织失败: %v", err)
	}
	if existingOrg != nil {
		return 0, fmt.Errorf("组织名称已存在")
	}

	// 设置默认值
	if orgData.OrgLevel == 0 {
		if orgData.OrgType == "1" {
			orgData.OrgLevel = 1
		} else if orgData.OrgType == "2" {
			orgData.OrgLevel = 2
		}
	}

	// 转换OrgType为int32
	var orgType int32
	if orgData.OrgType == "1" {
		orgType = 1
	} else if orgData.OrgType == "2" {
		orgType = 2
	} else {
		orgType = 1 // 默认为group
	}

	// 生成默认路径
	orgPath := orgData.OrgPath
	if orgPath == "" {
		orgPath = fmt.Sprintf("/%s", orgData.OrgName)
	}

	// 创建组织记录
	now := time.Now().Unix()
	newOrg := &org_info.OrgInfo{
		OrgName:     orgData.OrgName,
		OrgType:     orgType,
		OrgPath:     orgPath,
		OrgLevel:    orgData.OrgLevel,
		CurrentCnt:  0,
		MaxCnt:      orgData.MaxCnt,
		Status:      1,    // 默认启用
		ExtData:     "{}", // 默认空JSON
		CreatedAt:   now,
		UpdatedAt:   now,
		CreatedUser: ctx.SessionUserInfo().UserName,
		UpdatedUser: ctx.SessionUserInfo().UserName,
	}

	// 保存到数据库
	id, err = newOrg.Create(s.db.GetDbW())
	if err != nil {
		return 0, fmt.Errorf("创建组织失败: %v", err)
	}

	return id, nil
}
