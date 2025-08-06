package org_info

import (
	"fmt"
	"time"

	"github.com/xinliangnote/go-gin-api/internal/pkg/core"
	"github.com/xinliangnote/go-gin-api/internal/repository/mysql"
	"github.com/xinliangnote/go-gin-api/internal/repository/mysql/org_info"
)

// Update 更新组织信息
func (s *service) Update(ctx core.Context, orgId string, updateData *UpdateOrgInfoData) (err error) {
	// 转换orgId为uint64
	var id uint64
	fmt.Sscanf(orgId, "%d", &id)

	// 查询现有组织
	orgQueryBuilder := org_info.NewQueryBuilder()
	orgQueryBuilder.WhereId(mysql.EqualPredicate, id)

	existingOrg, err := orgQueryBuilder.QueryOne(s.db.GetDbR())
	if err != nil {
		return fmt.Errorf("查询组织失败: %v", err)
	}
	if existingOrg == nil {
		return fmt.Errorf("组织不存在")
	}

	// 构建更新字段
	updateFields := make(map[string]interface{})
	updateFields["updated_at"] = time.Now().Unix()
	updateFields["updated_user"] = ctx.SessionUserInfo().UserName

	if updateData.OrgName != "" && updateData.OrgName != existingOrg.OrgName {
		updateFields["org_name"] = updateData.OrgName
	}

	// 转换OrgType为int32
	if updateData.OrgType != "" {
		var orgType int32
		if updateData.OrgType == "1" {
			orgType = 1
		} else if updateData.OrgType == "2" {
			orgType = 2
		} else {
			orgType = existingOrg.OrgType
		}
		if orgType != existingOrg.OrgType {
			updateFields["org_type"] = orgType
		}
	}

	if updateData.OrgPath != "" && updateData.OrgPath != existingOrg.OrgPath {
		updateFields["org_path"] = updateData.OrgPath
	}

	if updateData.OrgLevel > 0 && updateData.OrgLevel != existingOrg.OrgLevel {
		updateFields["org_level"] = updateData.OrgLevel
	}

	if updateData.MaxCnt > 0 && updateData.MaxCnt != existingOrg.MaxCnt {
		updateFields["max_cnt"] = updateData.MaxCnt
	}

	// 执行更新
	updateQueryBuilder := org_info.NewQueryBuilder()
	updateQueryBuilder.WhereId(mysql.EqualPredicate, existingOrg.Id)

	err = updateQueryBuilder.Updates(s.db.GetDbW(), updateFields)
	if err != nil {
		return fmt.Errorf("更新组织失败: %v", err)
	}

	return nil
}

// Delete 删除组织信息
func (s *service) Delete(ctx core.Context, orgId string) (err error) {
	// 转换orgId为uint64
	var id uint64
	fmt.Sscanf(orgId, "%d", &id)

	// 查询现有组织
	orgQueryBuilder := org_info.NewQueryBuilder()
	orgQueryBuilder.WhereId(mysql.EqualPredicate, id)

	existingOrg, err := orgQueryBuilder.QueryOne(s.db.GetDbR())
	if err != nil {
		return fmt.Errorf("查询组织失败: %v", err)
	}
	if existingOrg == nil {
		return fmt.Errorf("组织不存在")
	}

	// 检查是否有子组织
	children, err := s.GetChildren(ctx, orgId)
	if err != nil {
		return fmt.Errorf("查询子组织失败: %v", err)
	}
	if len(children) > 0 {
		return fmt.Errorf("组织下还有子组织，无法删除")
	}

	// 执行删除
	deleteQueryBuilder := org_info.NewQueryBuilder()
	deleteQueryBuilder.WhereId(mysql.EqualPredicate, existingOrg.Id)

	err = deleteQueryBuilder.Delete(s.db.GetDbW())
	if err != nil {
		return fmt.Errorf("删除组织失败: %v", err)
	}

	return nil
}

// UpdateMemberCount 更新成员数量
func (s *service) UpdateMemberCount(ctx core.Context, orgId string, count int32) (err error) {
	// 转换orgId为uint64
	var id uint64
	fmt.Sscanf(orgId, "%d", &id)

	// 查询现有组织
	orgQueryBuilder := org_info.NewQueryBuilder()
	orgQueryBuilder.WhereId(mysql.EqualPredicate, id)

	existingOrg, err := orgQueryBuilder.QueryOne(s.db.GetDbR())
	if err != nil {
		return fmt.Errorf("查询组织失败: %v", err)
	}
	if existingOrg == nil {
		return fmt.Errorf("组织不存在")
	}

	// 更新成员数量
	updateFields := map[string]interface{}{
		"current_cnt":        count,
		"modified_timestamp": time.Now().Unix(),
		"updated_user":       ctx.SessionUserInfo().UserName,
	}

	updateQueryBuilder := org_info.NewQueryBuilder()
	updateQueryBuilder.WhereId(mysql.EqualPredicate, existingOrg.Id)

	err = updateQueryBuilder.Updates(s.db.GetDbW(), updateFields)
	if err != nil {
		return fmt.Errorf("更新成员数量失败: %v", err)
	}

	return nil
}
