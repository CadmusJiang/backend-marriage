package orgsvc

import (
	"strconv"
	"time"

	"github.com/xinliangnote/go-gin-api/internal/repository/mysql"
)

type service struct {
	db mysql.Repo
}

func New(db mysql.Repo) Service {
	return &service{db: db}
}

func (s *service) ListGroups(ctx Context, current, pageSize int, keyword string) ([]map[string]interface{}, int64, error) {
	db := s.db.GetDbR().WithContext(ctx.RequestContext()).Table("org").Where("org_type = 1")
	if keyword != "" {
		db = db.Where("name LIKE ?", "%"+keyword+"%")
	}

	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var rows []struct {
		Id        uint32    `db:"id"`
		Username  string    `db:"username"`
		Name      string    `db:"name"`
		Status    int32     `db:"status"`
		CreatedAt time.Time `db:"created_at"`
		UpdatedAt time.Time `db:"updated_at"`
	}

	if err := db.Select("id, username, name, status, created_at, updated_at").
		Order("id DESC").Limit(pageSize).Offset((current - 1) * pageSize).Find(&rows).Error; err != nil {
		return nil, 0, err
	}

	out := make([]map[string]interface{}, 0, len(rows))
	for _, r := range rows {
		out = append(out, map[string]interface{}{
			"id":        strconv.Itoa(int(r.Id)),
			"username":  r.Username,
			"name":      r.Name,
			"status":    r.Status,
			"version":   0,
			"createdAt": r.CreatedAt.Format("2006-01-02T15:04:05.000Z"),
			"updatedAt": r.UpdatedAt.Format("2006-01-02T15:04:05.000Z"),
		})
	}

	return out, total, nil
}

func (s *service) CreateGroup(ctx Context, payload *CreateGroupPayload) (uint32, error) {
	now := time.Now()
	m := map[string]interface{}{
		"org_type":     1,
		"parent_id":    0,
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

func (s *service) GetGroup(ctx Context, id uint32) (map[string]interface{}, error) {
	var row struct {
		Id        uint32    `db:"id"`
		Username  string    `db:"username"`
		Name      string    `db:"name"`
		Status    int32     `db:"status"`
		CreatedAt time.Time `db:"created_at"`
		UpdatedAt time.Time `db:"updated_at"`
	}

	if err := s.db.GetDbR().WithContext(ctx.RequestContext()).
		Select("id, username, name, status, created_at, updated_at").
		Table("org").Where("id = ? AND org_type = 1", id).Take(&row).Error; err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"id":        strconv.Itoa(int(row.Id)),
		"username":  row.Username,
		"name":      row.Name,
		"status":    row.Status,
		"version":   0,
		"createdAt": row.CreatedAt.Format("2006-01-02T15:04:05.000Z"),
		"updatedAt": row.UpdatedAt.Format("2006-01-02T15:04:05.000Z"),
	}, nil
}

func (s *service) UpdateGroup(ctx Context, id uint32, payload *UpdateGroupPayload) (uint32, error) {
	updates := map[string]interface{}{
		"name":         payload.Name,
		"status":       payload.Status,
		"updated_at":   time.Now(),
		"updated_user": ctx.SessionUserInfo().UserName,
		"version":      payload.Version + 1,
	}

	result := s.db.GetDbW().WithContext(ctx.RequestContext()).
		Table("org").Where("id = ? AND org_type = 1 AND version = ?", id, payload.Version).
		Updates(updates)

	if result.Error != nil {
		return 0, result.Error
	}

	return uint32(result.RowsAffected), nil
}

func (s *service) DeleteGroup(ctx Context, id uint32) error {
	return s.db.GetDbW().WithContext(ctx.RequestContext()).
		Table("org").Where("id = ? AND org_type = 1", id).Delete(&struct{}{}).Error
}
