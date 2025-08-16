package orgsvc

import (
	"fmt"
	"strconv"
	"time"

	"github.com/xinliangnote/go-gin-api/internal/authz"
)

func (s *service) ListTeams(ctx Context, belongGroupId uint32, current, pageSize int, keyword string, scope *authz.AccessScope) ([]map[string]interface{}, int64, error) {
	db := s.db.GetDbR().WithContext(ctx.RequestContext()).Table("org").Where("org_type = 'team'")

	// 添加调试日志
	fmt.Printf("DEBUG: ListTeams - belongGroupId: %d, scope: %+v\n", belongGroupId, scope)

	// 根据权限范围添加过滤条件
	if scope != nil {
		if scope.ScopeAll {
			// company_manager: 可以查看所有团队
			fmt.Printf("DEBUG: company_manager scope, belongGroupId: %d\n", belongGroupId)
			if belongGroupId > 0 {
				db = db.Where("parent_id = ?", belongGroupId)
				fmt.Printf("DEBUG: Added parent_id filter: %d\n", belongGroupId)
			} else {
				fmt.Printf("DEBUG: No parent_id filter added (belongGroupId = 0)\n")
			}
		} else if len(scope.AllowedGroupIDs) > 0 {
			// group_manager: 只能查看自己管理的组下的团队
			fmt.Printf("DEBUG: group_manager scope, allowed groups: %v\n", scope.AllowedGroupIDs)
			if belongGroupId > 0 {
				// 如果指定了组ID，验证权限
				hasPermission := false
				for _, allowedGroupID := range scope.AllowedGroupIDs {
					if allowedGroupID == int32(belongGroupId) {
						hasPermission = true
						break
					}
				}
				if !hasPermission {
					// 无权限，返回空结果
					fmt.Printf("DEBUG: No permission for group %d\n", belongGroupId)
					return []map[string]interface{}{}, 0, nil
				}
			}
			// 限制为允许的组
			db = db.Where("parent_id IN ?", scope.AllowedGroupIDs)
			fmt.Printf("DEBUG: Added parent_id IN filter: %v\n", scope.AllowedGroupIDs)
		} else if len(scope.AllowedTeamIDs) > 0 {
			// team_manager/employee: 只能查看自己所在的团队
			// 这里需要特殊处理，因为我们需要先知道用户所在的组
			// 暂时返回空结果，或者可以通过其他方式获取
			fmt.Printf("DEBUG: team_manager/employee scope, returning empty\n")
			return []map[string]interface{}{}, 0, nil
		} else {
			// 其他情况，返回空结果
			fmt.Printf("DEBUG: Other scope, returning empty\n")
			return []map[string]interface{}{}, 0, nil
		}
	} else {
		// 如果没有权限范围，使用传入的组ID
		fmt.Printf("DEBUG: No scope provided, using belongGroupId: %d\n", belongGroupId)
		if belongGroupId > 0 {
			db = db.Where("parent_id = ?", belongGroupId)
		}
	}

	if keyword != "" {
		db = db.Where("name LIKE ?", "%"+keyword+"%")
	}

	// 打印最终的SQL查询
	fmt.Printf("DEBUG: Final query conditions applied\n")

	var total int64
	if err := db.Count(&total).Error; err != nil {
		fmt.Printf("DEBUG: Count query failed: %v\n", err)
		return nil, 0, err
	}
	fmt.Printf("DEBUG: Total count: %d\n", total)

	var rows []struct {
		Id        uint32    `db:"id"`
		Name      string    `db:"name"`
		Username  string    `db:"username"`
		ParentId  uint32    `db:"parent_id"`
		Status    string    `db:"status"`
		CreatedAt time.Time `db:"created_at"`
		UpdatedAt time.Time `db:"updated_at"`
	}

	if err := db.Select("id, name, username, parent_id, status, created_at, updated_at").Order("id DESC").Limit(pageSize).Offset((current - 1) * pageSize).Find(&rows).Error; err != nil {
		fmt.Printf("DEBUG: Select query failed: %v\n", err)
		return nil, 0, err
	}
	fmt.Printf("DEBUG: Found %d rows\n", len(rows))

	out := make([]map[string]interface{}, 0, len(rows))
	for _, r := range rows {
		out = append(out, map[string]interface{}{
			"id":            strconv.Itoa(int(r.Id)),
			"name":          r.Name,
			"username":      r.Username,
			"belongGroupId": strconv.Itoa(int(r.ParentId)),
			"status":        r.Status,
			"version":       0,
			"createdAt":     r.CreatedAt.Format("2006-01-02T15:04:05.000Z"),
			"updatedAt":     r.UpdatedAt.Format("2006-01-02T15:04:05.000Z"),
		})
	}

	return out, total, nil
}

func (s *service) CreateTeam(ctx Context, payload *CreateTeamPayload) (uint32, error) {
	now := time.Now()
	gid := payload.BelongGroupId
	m := map[string]interface{}{
		"org_type":     "team",
		"parent_id":    gid,
		"path":         "/",
		"username":     payload.Username,
		"name":         payload.Name,
		"status":       1,
		"created_at":   now,
		"updated_at":   now,
		"created_user": ctx.SessionUserInfo().UserName,
		"updated_user": ctx.SessionUserInfo().UserName,
	}

	result := s.db.GetDbW().WithContext(ctx.RequestContext()).Table("org").Create(m)
	if result.Error != nil {
		return 0, result.Error
	}

	return uint32(result.RowsAffected), nil
}

func (s *service) GetTeam(ctx Context, id uint32) (map[string]interface{}, error) {
	var row struct {
		Id        uint32    `db:"id"`
		Username  string    `db:"username"`
		Name      string    `db:"name"`
		ParentId  uint32    `db:"parent_id"`
		Status    string    `db:"status"`
		CreatedAt time.Time `db:"created_at"`
		UpdatedAt time.Time `db:"updated_at"`
	}

	if err := s.db.GetDbR().WithContext(ctx.RequestContext()).Table("org").Select("id, username, name, parent_id, status, created_at, updated_at").Where("id = ? AND org_type = 'team'", id).Take(&row).Error; err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"id":            strconv.Itoa(int(row.Id)),
		"username":      row.Username,
		"name":          row.Name,
		"belongGroupId": strconv.Itoa(int(row.ParentId)),
		"status":        row.Status,
		"version":       0,
		"createdAt":     row.CreatedAt.Format("2006-01-02T15:04:05.000Z"),
		"updatedAt":     row.UpdatedAt.Format("2006-01-02T15:04:05.000Z"),
	}, nil
}

func (s *service) UpdateTeam(ctx Context, id uint32, payload *UpdateTeamPayload) (uint32, error) {
	updates := map[string]interface{}{
		"name":         payload.Name,
		"status":       payload.Status,
		"updated_at":   time.Now(),
		"updated_user": ctx.SessionUserInfo().UserName,
		"version":      payload.Version + 1,
	}

	result := s.db.GetDbW().WithContext(ctx.RequestContext()).
		Table("org").Where("id = ? AND org_type = 'team' AND version = ?", id, payload.Version).
		Updates(updates)

	if result.Error != nil {
		return 0, result.Error
	}

	return uint32(result.RowsAffected), nil
}

func (s *service) DeleteTeam(ctx Context, id uint32) error {
	return s.db.GetDbW().WithContext(ctx.RequestContext()).
		Table("org").Where("id = ? AND org_type = 'team'", id).Delete(&struct{}{}).Error
}

func (s *service) ListMembers(ctx Context, orgId uint32, current, pageSize int, keyword string) ([]map[string]interface{}, int64, error) {
	// Implementation for listing members
	return []map[string]interface{}{}, 0, nil
}

func (s *service) AddMember(ctx Context, orgId uint32, accountId uint32) error {
	// Implementation for adding member
	return nil
}

func (s *service) RemoveMember(ctx Context, orgId uint32, accountId uint32) error {
	// Implementation for removing member
	return nil
}
