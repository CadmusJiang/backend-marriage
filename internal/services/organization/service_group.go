package organization

import (
	"fmt"
	"strconv"
	"time"

	"github.com/xinliangnote/go-gin-api/internal/authz"
	"github.com/xinliangnote/go-gin-api/internal/pkg/core"
)

func (s *service) ListGroups(ctx core.Context, current, pageSize int, keyword string) (list interface{}, total int64, err error) {
	db := s.db.GetDbR().WithContext(ctx.RequestContext()).Table("org").Where("org_type = ?", 1)
	// 范围：公司管理员可见全部；组管仅可见本组
	if scope, e := authz.ComputeScope(ctx, s.db); e == nil && !scope.ScopeAll {
		if len(scope.AllowedGroupIDs) > 0 {
			db = db.Where("id IN (?)", scope.AllowedGroupIDs)
		} else {
			// 没有组范围的（队管/员工），不返回任何组
			db = db.Where("1 = 0")
		}
	}
	if keyword != "" {
		db = db.Where("nickname LIKE ?", "%"+keyword+"%")
	}
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var rows []struct {
		Id                uint32
		Username          string
		Nickname          string
		Status            int32
		CreatedTimestamp  int64
		ModifiedTimestamp int64
	}
	if err := db.Select("id, username, nickname, status, created_at, updated_at").
		Order("id DESC").Limit(pageSize).Offset((current - 1) * pageSize).Find(&rows).Error; err != nil {
		return nil, 0, err
	}
	type groupItem struct {
		Id        string `json:"id"`
		Username  string `json:"username"`
		Nickname  string `json:"nickname"`
		Status    int32  `json:"status"`
		Version   int    `json:"version"`
		CreatedAt string `json:"createdAt"`
		UpdatedAt string `json:"updatedAt"`
	}
	out := make([]groupItem, 0, len(rows))
	for _, r := range rows {
		out = append(out, groupItem{
			Id:        strconv.Itoa(int(r.Id)),
			Username:  r.Username,
			Nickname:  r.Nickname,
			Status:    r.Status,
			Version:   0,
			CreatedAt: time.Unix(r.CreatedTimestamp, 0).Format("2006-01-02T15:04:05.000Z"),
			UpdatedAt: time.Unix(r.ModifiedTimestamp, 0).Format("2006-01-02T15:04:05.000Z"),
		})
	}
	return out, total, nil
}

func (s *service) CreateGroup(ctx core.Context, payload *CreateGroupPayload) (id uint32, err error) {
	now := time.Now().Unix()
	m := map[string]interface{}{
		"org_type":     1,
		"parent_id":    0,
		"path":         "/",
		"username":     payload.Username,
		"nickname":     payload.Nickname,
		"status":       1,
		"created_at":   now,
		"updated_at":   now,
		"created_user": ctx.SessionUserInfo().UserName,
		"updated_user": ctx.SessionUserInfo().UserName,
	}
	if err := s.db.GetDbW().WithContext(ctx.RequestContext()).Table("org").Create(m).Error; err != nil {
		return 0, err
	}
	var newId uint32
	if err := s.db.GetDbR().Raw("SELECT LAST_INSERT_ID()").Scan(&newId).Error; err != nil {
		return 0, err
	}
	if err := s.db.GetDbW().WithContext(ctx.RequestContext()).Table("org").Where("id = ?", newId).
		Update("path", fmt.Sprintf("/%d/", newId)).Error; err != nil {
		return newId, err
	}
	return newId, nil
}

func (s *service) GetGroup(ctx core.Context, orgId string) (info interface{}, err error) {
	id, _ := strconv.Atoi(orgId)
	var row struct {
		Id                uint32
		Username          string
		Nickname          string
		Status            int32
		CreatedTimestamp  int64
		ModifiedTimestamp int64
	}
	if err := s.db.GetDbR().WithContext(ctx.RequestContext()).Table("org").
		Select("id, username, nickname, status, created_at, updated_at").
		Where("id = ? AND org_type = 1", id).Take(&row).Error; err != nil {
		return nil, err
	}
	if scope, e := authz.ComputeScope(ctx, s.db); e == nil && !scope.ScopeAll {
		allowed := false
		for _, gid := range scope.AllowedGroupIDs {
			if int(gid) == int(row.Id) {
				allowed = true
				break
			}
		}
		if !allowed {
			return nil, fmt.Errorf("forbidden")
		}
	}
	return map[string]interface{}{
		"id":        strconv.Itoa(int(row.Id)),
		"username":  row.Username,
		"nickname":  row.Nickname,
		"status":    row.Status,
		"version":   0,
		"createdAt": time.Unix(row.CreatedTimestamp, 0).Format("2006-01-02T15:04:05.000Z"),
		"updatedAt": time.Unix(row.ModifiedTimestamp, 0).Format("2006-01-02T15:04:05.000Z"),
	}, nil
}

func (s *service) UpdateGroup(ctx core.Context, orgId string, payload *UpdateGroupPayload) (newVersion int, err error) {
	id, _ := strconv.Atoi(orgId)
	updates := map[string]interface{}{
		"nickname":     payload.Nickname,
		"status":       payload.Status,
		"updated_at":   time.Now().Unix(),
		"updated_user": ctx.SessionUserInfo().UserName,
		"version":      payload.Version + 1,
	}
	tx := s.db.GetDbW().WithContext(ctx.RequestContext()).Table("org").
		Where("id = ? AND org_type = 1 AND version = ?", id, payload.Version).Updates(updates)
	if tx.Error != nil {
		return 0, tx.Error
	}
	if tx.RowsAffected == 0 {
		return 0, fmt.Errorf("version_conflict")
	}
	return payload.Version + 1, nil
}

func (s *service) ListGroupHistory(ctx core.Context, orgId string, current, pageSize int) (items interface{}, total int64, err error) {
	return []interface{}{}, 0, nil
}
