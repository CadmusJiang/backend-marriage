package org_info

import (
	"fmt"
	"strings"

	"github.com/xinliangnote/go-gin-api/internal/pkg/core"
	"github.com/xinliangnote/go-gin-api/internal/repository/mysql"
	"github.com/xinliangnote/go-gin-api/internal/repository/mysql/org_info"
)

// PageList 分页获取组织信息列表
func (s *service) PageList(ctx core.Context, searchData *SearchOrgInfoData) (listData []*org_info.OrgInfo, err error) {
	// 检查数据库连接
	if s.db == nil {
		return nil, fmt.Errorf("数据库连接为空")
	}

	orgQueryBuilder := org_info.NewQueryBuilder()

	// 添加搜索条件
	if searchData.OrgName != "" {
		orgQueryBuilder.WhereOrgName(mysql.LikePredicate, "%"+searchData.OrgName+"%")
	}
	if searchData.OrgType != "" {
		// 转换OrgType为int32
		var orgType int32
		if searchData.OrgType == "1" {
			orgType = 1
		} else if searchData.OrgType == "2" {
			orgType = 2
		} else {
			orgType = 1 // 默认为group
		}
		orgQueryBuilder.WhereOrgType(mysql.EqualPredicate, orgType)
	}
	if searchData.OrgLevel > 0 {
		orgQueryBuilder.WhereOrgLevel(mysql.EqualPredicate, searchData.OrgLevel)
	}

	// 设置分页
	offset := (searchData.Current - 1) * searchData.PageSize
	orgQueryBuilder.Limit(searchData.PageSize).Offset(offset)

	// 按创建时间倒序排列
	orgQueryBuilder.OrderByCreatedAt(false)

	// 查询数据
	listData, err = orgQueryBuilder.QueryAll(s.db.GetDbR())
	if err != nil {
		// 检查是否是表不存在的错误
		if strings.Contains(err.Error(), "doesn't exist") || strings.Contains(err.Error(), "Table") {
			return nil, fmt.Errorf("org表不存在，请先创建数据库表: %v", err)
		}
		return nil, fmt.Errorf("查询组织列表失败: %v", err)
	}

	return listData, nil
}

// PageListCount 获取组织信息列表总数
func (s *service) PageListCount(ctx core.Context, searchData *SearchOrgInfoData) (total int64, err error) {
	orgQueryBuilder := org_info.NewQueryBuilder()

	// 添加搜索条件
	if searchData.OrgName != "" {
		orgQueryBuilder.WhereOrgName(mysql.LikePredicate, "%"+searchData.OrgName+"%")
	}
	if searchData.OrgType != "" {
		// 转换OrgType为int32
		var orgType int32
		if searchData.OrgType == "1" {
			orgType = 1
		} else if searchData.OrgType == "2" {
			orgType = 2
		} else {
			orgType = 1 // 默认为group
		}
		orgQueryBuilder.WhereOrgType(mysql.EqualPredicate, orgType)
	}
	if searchData.OrgLevel > 0 {
		orgQueryBuilder.WhereOrgLevel(mysql.EqualPredicate, searchData.OrgLevel)
	}

	// 查询总数
	total, err = orgQueryBuilder.Count(s.db.GetDbR())
	if err != nil {
		return 0, fmt.Errorf("查询组织总数失败: %v", err)
	}

	return total, nil
}

// Detail 获取组织信息详情
func (s *service) Detail(ctx core.Context, orgId string) (info *org_info.OrgInfo, err error) {
	// 转换orgId为uint64
	var id uint64
	fmt.Sscanf(orgId, "%d", &id)

	orgQueryBuilder := org_info.NewQueryBuilder()
	orgQueryBuilder.WhereId(mysql.EqualPredicate, id)

	// 查询组织详情
	info, err = orgQueryBuilder.QueryOne(s.db.GetDbR())
	if err != nil {
		return nil, fmt.Errorf("查询组织详情失败: %v", err)
	}
	if info == nil {
		return nil, fmt.Errorf("组织不存在")
	}

	return info, nil
}

// GetChildren 获取子组织
func (s *service) GetChildren(ctx core.Context, parentOrgId string) (children []*org_info.OrgInfo, err error) {
	// 由于没有parent_org_id字段，这里暂时返回空
	// 后续可以通过org_path来实现层级关系
	return []*org_info.OrgInfo{}, nil
}

// GetParent 获取父组织
func (s *service) GetParent(ctx core.Context, orgId string) (parent *org_info.OrgInfo, err error) {
	// 由于没有parent_org_id字段，这里暂时返回nil
	// 后续可以通过org_path来实现层级关系
	return nil, nil
}
