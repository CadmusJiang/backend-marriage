package orgsvc

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/xinliangnote/go-gin-api/internal/authz"
	"github.com/xinliangnote/go-gin-api/internal/repository/mysql"
)

type service struct {
	db mysql.Repo
}

func New(db mysql.Repo) Service {
	return &service{db: db}
}

func (s *service) i() {}

func (s *service) ListGroups(ctx Context, current, pageSize int, keyword string, scope *authz.AccessScope) ([]map[string]interface{}, int64, error) {
	db := s.db.GetDbR().WithContext(ctx.RequestContext()).Table("org").Where("org_type = 'group'")

	// 根据权限范围添加过滤条件
	if scope != nil {
		if scope.ScopeAll {
			// company_manager: 可以查看所有组
			// 不需要额外过滤
		} else if len(scope.AllowedGroupIDs) > 0 {
			// group_manager: 只能查看自己管理的组
			db = db.Where("id IN ?", scope.AllowedGroupIDs)
		} else {
			// 其他角色没有组权限，返回空结果
			return []map[string]interface{}{}, 0, nil
		}
	}

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
		Status    string    `db:"status"`
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
		"org_type":     "group",
		"parent_id":    0,
		"path":         "/",
		"username":     payload.Username,
		"name":         payload.Name,
		"status":       "enabled",
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

func (s *service) GetGroup(ctx Context, id string) (map[string]interface{}, error) {
	// 验证组ID格式
	groupIDUint, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		return nil, fmt.Errorf("组ID格式错误：'%s'，必须是正整数", id)
	}
	if groupIDUint <= 0 {
		return nil, fmt.Errorf("组ID错误：%d，必须大于0", groupIDUint)
	}
	groupID := uint32(groupIDUint)

	var row struct {
		Id        uint32    `db:"id"`
		Username  string    `db:"username"`
		Name      string    `db:"name"`
		Address   *string   `db:"address"`
		Extra     *string   `db:"extra"`
		Status    string    `db:"status"`
		CreatedAt time.Time `db:"created_at"`
		UpdatedAt time.Time `db:"updated_at"`
	}

	if err := s.db.GetDbR().WithContext(ctx.RequestContext()).Table("org").Select("id, username, name, address, extra, status, created_at, updated_at").Where("id = ? AND org_type = 'group'", groupID).Take(&row).Error; err != nil {
		return nil, err
	}

	result := map[string]interface{}{
		"id":        strconv.Itoa(int(row.Id)),
		"username":  row.Username,
		"name":      row.Name,
		"status":    row.Status,
		"version":   0,
		"createdAt": row.CreatedAt.Format("2006-01-02T15:04:05.000Z"),
		"updatedAt": row.UpdatedAt.Format("2006-01-02T15:04:05.000Z"),
	}

	// 添加可选字段
	if row.Address != nil {
		result["address"] = *row.Address
	}
	if row.Extra != nil {
		// 解析JSON字段
		var extraData map[string]interface{}
		if err := json.Unmarshal([]byte(*row.Extra), &extraData); err == nil {
			result["extra"] = extraData
		}
	}

	return result, nil
}

func (s *service) UpdateGroup(ctx Context, id string, payload *UpdateGroupPayload) (uint32, error) {
	// 验证组ID格式
	groupIDUint, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		return 0, fmt.Errorf("组ID格式错误：'%s'，必须是正整数", id)
	}
	if groupIDUint <= 0 {
		return 0, fmt.Errorf("组ID错误：%d，必须大于0", groupIDUint)
	}
	groupID := uint32(groupIDUint)

	updates := map[string]interface{}{
		"updated_at":   time.Now(),
		"updated_user": ctx.SessionUserInfo().UserName,
		"version":      payload.Version + 1,
	}

	// 处理name字段：如果payload中提供了name，则使用；否则保持原值不变
	if payload.Name != "" {
		updates["name"] = payload.Name
	}

	// 处理status字段：如果payload中提供了status，则使用；否则保持原值不变
	if payload.Status != "" {
		// 验证status值是否有效
		if payload.Status != "enabled" && payload.Status != "disabled" {
			return 0, fmt.Errorf("无效的status值：'%s'，必须是 'enabled' 或 'disabled'", payload.Status)
		}
		updates["status"] = payload.Status
	}
	// 如果payload.Status为空字符串，则不更新status字段，保持数据库中的原值

	// 添加可选字段
	if payload.Address != nil {
		updates["address"] = *payload.Address
	}
	if payload.Extra != nil {
		// 将map序列化为JSON字符串
		if extraJSON, err := json.Marshal(payload.Extra); err == nil {
			updates["extra"] = string(extraJSON)
		}
	}

	result := s.db.GetDbW().WithContext(ctx.RequestContext()).
		Table("org").Where("id = ? AND org_type = 'group' AND version = ?", groupID, payload.Version).
		Updates(updates)

	if result.Error != nil {
		return 0, result.Error
	}

	return uint32(result.RowsAffected), nil
}

func (s *service) DeleteGroup(ctx Context, id string) error {
	// 验证组ID格式
	groupIDUint, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		return fmt.Errorf("组ID格式错误：'%s'，必须是正整数", id)
	}
	if groupIDUint <= 0 {
		return fmt.Errorf("组ID错误：%d，必须大于0", groupIDUint)
	}
	groupID := uint32(groupIDUint)

	return s.db.GetDbW().WithContext(ctx.RequestContext()).
		Table("org").Where("id = ? AND org_type = 'group'", groupID).Delete(&struct{}{}).Error
}
