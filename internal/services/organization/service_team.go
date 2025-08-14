package organization

import (
	"fmt"
	"strconv"
	"time"

	"github.com/xinliangnote/go-gin-api/internal/authz"
	"github.com/xinliangnote/go-gin-api/internal/pkg/core"
)

func (s *service) ListTeams(ctx core.Context, current, pageSize int, belongGroupId string, keyword string) (list interface{}, total int64, err error) {
	db := s.db.GetDbR().WithContext(ctx.RequestContext()).Table("org").Where("org_type = ?", 2)
	if keyword != "" {
		db = db.Where("nickname LIKE ?", "%"+keyword+"%")
	}
	if belongGroupId != "" {
		gid, _ := strconv.Atoi(belongGroupId)
		db = db.Where("parent_id = ?", gid)
	}
	// 范围：公司管理员可见全部；组/队/员工可见本组下团队
	if scope, e := authz.ComputeScope(ctx, s.db); e == nil && !scope.ScopeAll {
		if len(scope.AllowedGroupIDs) > 0 {
			db = db.Where("parent_id IN (?)", scope.AllowedGroupIDs)
		} else {
			// 没有组范围则不返回任何团队
			db = db.Where("1 = 0")
		}
	}
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var rows []struct {
		Id                uint32
		Nickname          string
		Username          string
		ParentId          uint32
		Status            int32
		CreatedTimestamp  int64
		ModifiedTimestamp int64
	}
	if err := db.Select("id, nickname, username, parent_id, status, created_at, updated_at").Order("id DESC").Limit(pageSize).Offset((current - 1) * pageSize).Find(&rows).Error; err != nil {
		return nil, 0, err
	}
	type teamItem struct {
		Id            string `json:"id"`
		Nickname      string `json:"nickname"`
		Username      string `json:"username"`
		BelongGroupId string `json:"belongGroupId"`
		Status        int32  `json:"status"`
		Version       int    `json:"version"`
		CreatedAt     string `json:"createdAt"`
		UpdatedAt     string `json:"updatedAt"`
	}
	out := make([]teamItem, 0, len(rows))
	for _, r := range rows {
		out = append(out, teamItem{Id: strconv.Itoa(int(r.Id)), Nickname: r.Nickname, Username: r.Username, BelongGroupId: strconv.Itoa(int(r.ParentId)), Status: r.Status, Version: 0, CreatedAt: time.Unix(r.CreatedTimestamp, 0).Format("2006-01-02T15:04:05.000Z"), UpdatedAt: time.Unix(r.ModifiedTimestamp, 0).Format("2006-01-02T15:04:05.000Z")})
	}
	return out, total, nil
}

func (s *service) CreateTeam(ctx core.Context, payload *CreateTeamPayload) (id uint32, err error) {
	now := time.Now().Unix()
	gid, _ := strconv.Atoi(payload.BelongGroupId)
	m := map[string]interface{}{"org_type": 2, "parent_id": gid, "path": "/", "username": payload.Username, "nickname": payload.Nickname, "status": 1, "created_at": now, "updated_at": now, "created_user": ctx.SessionUserInfo().UserName, "updated_user": ctx.SessionUserInfo().UserName}
	if err := s.db.GetDbW().WithContext(ctx.RequestContext()).Table("org").Create(m).Error; err != nil {
		return 0, err
	}
	var newId uint32
	if err := s.db.GetDbR().Raw("SELECT LAST_INSERT_ID()").Scan(&newId).Error; err != nil {
		return 0, err
	}
	if err := s.db.GetDbW().WithContext(ctx.RequestContext()).Table("org").Where("id = ?", newId).Update("path", fmt.Sprintf("/%d/%d/", gid, newId)).Error; err != nil {
		return newId, err
	}
	return newId, nil
}

func (s *service) GetTeam(ctx core.Context, teamId string) (info interface{}, err error) {
	id, _ := strconv.Atoi(teamId)
	var row struct {
		Id                uint32
		Username          string
		Nickname          string
		ParentId          uint32
		Status            int32
		CreatedTimestamp  int64
		ModifiedTimestamp int64
	}
	if err := s.db.GetDbR().WithContext(ctx.RequestContext()).Table("org").Select("id, username, nickname, parent_id, status, created_at, updated_at").Where("id = ? AND org_type = 2", id).Take(&row).Error; err != nil {
		return nil, err
	}
	if scope, e := authz.ComputeScope(ctx, s.db); e == nil && !scope.ScopeAll {
		allowed := false
		for _, gid := range scope.AllowedGroupIDs {
			if int(gid) == int(row.ParentId) {
				allowed = true
				break
			}
		}
		if !allowed {
			return nil, fmt.Errorf("forbidden")
		}
	}
	return map[string]interface{}{"id": strconv.Itoa(int(row.Id)), "username": row.Username, "nickname": row.Nickname, "belongGroupId": strconv.Itoa(int(row.ParentId)), "status": row.Status, "version": 0, "createdAt": time.Unix(row.CreatedTimestamp, 0).Format("2006-01-02T15:04:05.000Z"), "updatedAt": time.Unix(row.ModifiedTimestamp, 0).Format("2006-01-02T15:04:05.000Z")}, nil
}

func (s *service) UpdateTeam(ctx core.Context, teamId string, payload *UpdateTeamPayload) (newVersion int, err error) {
	id, _ := strconv.Atoi(teamId)
	// 乐观锁更新
	updates := map[string]interface{}{"nickname": payload.Nickname, "status": payload.Status, "updated_at": time.Now().Unix(), "updated_user": ctx.SessionUserInfo().UserName, "version": payload.Version + 1}
	tx := s.db.GetDbW().WithContext(ctx.RequestContext()).Table("org").Where("id = ? AND org_type = 2 AND version = ?", id, payload.Version).Updates(updates)
	if tx.Error != nil {
		return 0, tx.Error
	}
	if tx.RowsAffected == 0 {
		return 0, fmt.Errorf("version_conflict")
	}
	return payload.Version + 1, nil
}

func (s *service) ListTeamHistory(ctx core.Context, teamId string, current, pageSize int) (items interface{}, total int64, err error) {
	return []interface{}{}, 0, nil
}
