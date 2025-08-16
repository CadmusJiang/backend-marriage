package orgsvc

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/xinliangnote/go-gin-api/internal/authz"
	"github.com/xinliangnote/go-gin-api/internal/pkg/core"
	"github.com/xinliangnote/go-gin-api/internal/services/history"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// ListTeams 获取团队列表
func (s *service) ListTeams(ctx Context, belongGroupId string, current, pageSize int, keyword string, scope *authz.AccessScope) ([]map[string]interface{}, int64, error) {
	// 验证归属组ID格式
	belongGroupIDUint, err := strconv.ParseUint(belongGroupId, 10, 32)
	if err != nil {
		return nil, 0, fmt.Errorf("归属组ID格式错误：'%s'，必须是正整数", belongGroupId)
	}
	if belongGroupIDUint <= 0 {
		return nil, 0, fmt.Errorf("归属组ID错误：%d，必须大于0", belongGroupIDUint)
	}

	// 检查权限：只有同组用户或上级用户可以查看
	if scope.RoleType != "company_manager" {
		allowedGroupID := int32(belongGroupIDUint)
		hasPermission := false
		for _, gid := range scope.AllowedGroupIDs {
			if gid == allowedGroupID {
				hasPermission = true
				break
			}
		}
		if !hasPermission {
			return nil, 0, fmt.Errorf("权限不足，无法查看该组的团队信息")
		}
	}

	db := s.db.GetDbR().WithContext(ctx.RequestContext()).Table("org").Where("org_type = 'team'")

	// 获取trace_id用于日志关联
	traceID := ctx.GetTraceID()

	// 创建日志元数据
	meta := map[string]interface{}{
		"trace_id":      traceID,
		"service":       "organization-service",
		"operation":     "ListTeams",
		"belongGroupId": belongGroupIDUint,
		"current":       current,
		"pageSize":      pageSize,
		"keyword":       keyword,
	}

	// 添加调试日志
	zap.L().Debug("开始获取团队列表",
		zap.Uint32("belongGroupId", uint32(belongGroupIDUint)),
		zap.Any("scope", scope),
		zap.Any("meta", meta),
	)

	// 根据权限范围添加过滤条件
	if scope != nil {
		if scope.ScopeAll {
			// company_manager: 可以查看所有团队
			zap.L().Debug("使用company_manager权限范围",
				zap.Uint32("belongGroupId", uint32(belongGroupIDUint)),
				zap.Any("meta", meta),
			)
			if belongGroupIDUint > 0 {
				db = db.Where("parent_id = ?", belongGroupIDUint)
				zap.L().Debug("添加parent_id过滤条件",
					zap.Uint32("parentId", uint32(belongGroupIDUint)),
					zap.Any("meta", meta),
				)
			} else {
				zap.L().Debug("未添加parent_id过滤条件(belongGroupId = 0)",
					zap.Any("meta", meta),
				)
			}
		} else if len(scope.AllowedGroupIDs) > 0 {
			// group_manager: 只能查看自己管理的组下的团队
			zap.L().Debug("使用group_manager权限范围",
				zap.Any("allowedGroups", scope.AllowedGroupIDs),
				zap.Any("meta", meta),
			)
			if belongGroupIDUint > 0 {
				// 如果指定了组ID，验证权限
				hasPermission := false
				for _, allowedGroupID := range scope.AllowedGroupIDs {
					if allowedGroupID == int32(belongGroupIDUint) {
						hasPermission = true
						break
					}
				}
				if !hasPermission {
					// 无权限，返回空结果
					zap.L().Warn("用户无权限访问指定组",
						zap.Uint32("requestedGroupId", uint32(belongGroupIDUint)),
						zap.Any("allowedGroups", scope.AllowedGroupIDs),
						zap.Any("meta", meta),
					)
					return []map[string]interface{}{}, 0, nil
				}
			}
			// 限制为允许的组
			db = db.Where("parent_id IN ?", scope.AllowedGroupIDs)
			zap.L().Debug("添加parent_id IN过滤条件",
				zap.Any("allowedGroups", scope.AllowedGroupIDs),
				zap.Any("meta", meta),
			)
		} else if len(scope.AllowedTeamIDs) > 0 {
			// team_manager/employee: 只能查看自己所在的团队
			// 这里需要特殊处理，因为我们需要先知道用户所在的组
			// 暂时返回空结果，或者可以通过其他方式获取
			zap.L().Debug("使用team_manager/employee权限范围，暂时返回空结果",
				zap.Any("allowedTeams", scope.AllowedTeamIDs),
				zap.Any("meta", meta),
			)
			return []map[string]interface{}{}, 0, nil
		} else {
			// 其他情况，返回空结果
			zap.L().Debug("其他权限范围，返回空结果",
				zap.Any("meta", meta),
			)
			return []map[string]interface{}{}, 0, nil
		}
	} else {
		// 如果没有权限范围，使用传入的组ID
		zap.L().Debug("未提供权限范围，使用传入的belongGroupId",
			zap.Uint32("belongGroupId", uint32(belongGroupIDUint)),
			zap.Any("meta", meta),
		)
		if belongGroupIDUint > 0 {
			db = db.Where("parent_id = ?", belongGroupIDUint)
		}
	}

	// 添加关键词搜索
	if keyword != "" {
		db = db.Where("name LIKE ? OR username LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}

	zap.L().Debug("已应用所有查询条件",
		zap.Any("meta", meta),
	)

	// 统计总数
	var total int64
	if err := db.Count(&total).Error; err != nil {
		zap.L().Error("统计查询失败",
			zap.Uint32("belongGroupId", uint32(belongGroupIDUint)),
			zap.Error(err),
			zap.Any("meta", meta),
		)
		return nil, 0, err
	}

	zap.L().Debug("查询总数统计完成",
		zap.Uint32("belongGroupId", uint32(belongGroupIDUint)),
		zap.Int64("total", total),
		zap.Any("meta", meta),
	)

	// 查询数据
	var rows []struct {
		Id       uint32 `db:"id"`
		Name     string `db:"name"`
		Username string `db:"username"`
		OrgType  string `db:"org_type"`
		ParentId uint32 `db:"parent_id"`
		Status   string `db:"status"`
		Path     string `db:"path"`
	}

	if err := db.Select("id, name, username, org_type, parent_id, status, path").
		Order("id DESC").
		Limit(pageSize).
		Offset((current - 1) * pageSize).
		Find(&rows).Error; err != nil {
		zap.L().Error("数据查询失败",
			zap.Uint32("belongGroupId", uint32(belongGroupIDUint)),
			zap.Error(err),
			zap.Any("meta", meta),
		)
		return nil, 0, err
	}

	zap.L().Debug("数据查询完成",
		zap.Uint32("belongGroupId", uint32(belongGroupIDUint)),
		zap.Int("rowsCount", len(rows)),
		zap.Any("meta", meta),
	)

	// 转换为返回格式
	out := make([]map[string]interface{}, 0, len(rows))
	for _, r := range rows {
		team := map[string]interface{}{
			"id":       strconv.Itoa(int(r.Id)),
			"name":     r.Name,
			"username": r.Username,
			"orgType":  r.OrgType,
			"parentId": strconv.Itoa(int(r.ParentId)),
			"status":   r.Status,
			"path":     r.Path,
		}
		out = append(out, team)
	}

	zap.L().Info("团队列表获取成功",
		zap.Uint32("belongGroupId", uint32(belongGroupIDUint)),
		zap.Int("resultCount", len(out)),
		zap.Int64("totalCount", total),
		zap.Any("meta", meta),
	)

	return out, total, nil
}

func (s *service) CreateTeam(ctx Context, payload *CreateTeamPayload) (uint32, error) {
	// 获取trace_id用于日志关联
	traceID := ctx.GetTraceID()

	// 创建日志元数据
	meta := map[string]interface{}{
		"trace_id":      traceID,
		"service":       "organization-service",
		"operation":     "CreateTeam",
		"belongGroupId": payload.BelongGroupId,
		"username":      payload.Username,
		"name":          payload.Name,
	}

	zap.L().Info("开始创建团队",
		zap.String("username", payload.Username),
		zap.String("name", payload.Name),
		zap.Uint32("belongGroupId", payload.BelongGroupId),
		zap.Any("meta", meta),
	)

	// 开始事务
	tx := s.db.GetDbW().WithContext(ctx.RequestContext()).Begin()
	if tx.Error != nil {
		zap.L().Error("开始事务失败",
			zap.Error(tx.Error),
			zap.Any("meta", meta),
		)
		return 0, tx.Error
	}

	// 确保事务回滚（如果出错）
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	now := time.Now()
	gid := payload.BelongGroupId

	// 构造团队数据
	teamData := map[string]interface{}{
		"org_type":     "team",
		"parent_id":    gid,
		"path":         "/",
		"username":     payload.Username,
		"name":         payload.Name,
		"status":       "enabled",
		"created_at":   now,
		"updated_at":   now,
		"created_user": ctx.SessionUserInfo().UserName,
		"updated_user": ctx.SessionUserInfo().UserName,
	}

	// 1. 创建团队
	result := tx.Table("org").Create(teamData)
	if result.Error != nil {
		tx.Rollback()
		zap.L().Error("创建团队失败",
			zap.Error(result.Error),
			zap.Any("meta", meta),
		)
		return 0, result.Error
	}

	// 获取创建的团队ID
	var teamID uint32
	if err := tx.Table("org").Where("username = ? AND org_type = 'team'", payload.Username).Select("id").Scan(&teamID).Error; err != nil {
		tx.Rollback()
		zap.L().Error("获取团队ID失败",
			zap.Error(err),
			zap.Any("meta", meta),
		)
		return 0, err
	}

	// 2. 记录历史（使用通用history服务）
	// 构造新对象用于历史记录
	newTeam := map[string]interface{}{
		"org_type":  "team",
		"parent_id": gid,
		"username":  payload.Username,
		"name":      payload.Name,
		"status":    "enabled",
	}

	// 使用通用history服务记录历史
	// 创建适配器来满足history服务的接口要求
	historyService := history.NewHistoryService(&historyAdapter{repo: s.db})
	err := historyService.CompareAndCreateHistory(
		&contextAdapter{ctx: ctx},
		history.EntityTypeOrg,          // 实体类型
		teamID,                         // 实体ID
		nil,                            // 旧对象（创建时为nil）
		newTeam,                        // 新对象
		history.OperateTypeCreated,     // 操作类型
		ctx.SessionUserInfo().UserName, // 操作人
		ctx.SessionUserInfo().RoleType, // 操作人角色
	)

	if err != nil {
		tx.Rollback()
		zap.L().Error("记录团队创建历史失败",
			zap.Error(err),
			zap.Uint32("teamId", teamID),
			zap.Any("meta", meta),
		)
		return 0, err
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		zap.L().Error("提交事务失败",
			zap.Error(err),
			zap.Uint32("teamId", teamID),
			zap.Any("meta", meta),
		)
		return 0, err
	}

	zap.L().Info("团队创建成功",
		zap.Uint32("teamId", teamID),
		zap.String("username", payload.Username),
		zap.Any("meta", meta),
	)

	return teamID, nil
}

// historyAdapter 适配器，用于满足history服务的数据库接口要求
type historyAdapter struct {
	repo interface {
		GetDbW() *gorm.DB
	}
}

func (a *historyAdapter) GetDbW() interface {
	WithContext(ctx interface{}) interface {
		Table(name string) interface {
			Create(value interface{}) interface {
				Error() error
			}
		}
	}
} {
	return &gormAdapter{db: a.repo.GetDbW()}
}

// gormAdapter 适配GORM数据库接口
type gormAdapter struct {
	db *gorm.DB
}

func (g *gormAdapter) WithContext(ctx interface{}) interface {
	Table(name string) interface {
		Create(value interface{}) interface {
			Error() error
		}
	}
} {
	// 类型断言，确保ctx是context.Context类型
	if contextCtx, ok := ctx.(core.StdContext); ok {
		return &tableAdapter{db: g.db.WithContext(contextCtx)}
	}
	// 如果不是正确的类型，返回一个空的适配器
	return &tableAdapter{db: g.db}
}

// tableAdapter 适配表操作接口
type tableAdapter struct {
	db *gorm.DB
}

func (t *tableAdapter) Table(name string) interface {
	Create(value interface{}) interface {
		Error() error
	}
} {
	return &createAdapter{db: t.db.Table(name)}
}

// createAdapter 适配创建操作接口
type createAdapter struct {
	db *gorm.DB
}

func (c *createAdapter) Create(value interface{}) interface {
	Error() error
} {
	result := c.db.Create(value)
	return &resultAdapter{result: result}
}

// resultAdapter 适配结果接口
type resultAdapter struct {
	result *gorm.DB
}

func (r *resultAdapter) Error() error {
	return r.result.Error
}

// contextAdapter 适配器，用于满足history服务的上下文接口要求
type contextAdapter struct {
	ctx interface {
		RequestContext() core.StdContext
	}
}

func (a *contextAdapter) RequestContext() interface{} {
	return a.ctx.RequestContext()
}

// GetTeam 获取团队详情
func (s *service) GetTeam(ctx Context, id string) (map[string]interface{}, error) {
	// 验证团队ID格式
	teamIDUint, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		return nil, fmt.Errorf("团队ID格式错误：'%s'，必须是正整数", id)
	}
	if teamIDUint <= 0 {
		return nil, fmt.Errorf("团队ID错误：%d，必须大于0", teamIDUint)
	}
	teamID := uint32(teamIDUint)

	db := s.db.GetDbR().WithContext(ctx.RequestContext()).Table("org").Where("org_type = 'team'")

	// 获取trace_id用于日志关联
	traceID := ctx.GetTraceID()

	// 创建日志元数据
	meta := map[string]interface{}{
		"trace_id":  traceID,
		"service":   "organization-service",
		"operation": "GetTeam",
		"teamId":    teamID,
	}

	zap.L().Debug("开始获取团队详情",
		zap.Uint32("teamId", teamID),
		zap.Any("meta", meta),
	)

	var row struct {
		Id        uint32    `db:"id"`
		Username  string    `db:"username"`
		Name      string    `db:"name"`
		Address   *string   `db:"address"`
		Extra     *string   `db:"extra"`
		ParentId  uint32    `db:"parent_id"`
		Status    string    `db:"status"`
		CreatedAt time.Time `db:"created_at"`
		UpdatedAt time.Time `db:"updated_at"`
	}

	if err := db.Select("id, username, name, address, extra, parent_id, status, created_at, updated_at").Where("id = ? AND org_type = 'team'", teamID).Take(&row).Error; err != nil {
		zap.L().Error("获取团队详情失败",
			zap.Uint32("teamId", teamID),
			zap.Error(err),
			zap.Any("meta", meta),
		)
		return nil, err
	}

	teamInfo := map[string]interface{}{
		"id":            strconv.Itoa(int(row.Id)),
		"username":      row.Username,
		"name":          row.Name,
		"belongGroupId": strconv.Itoa(int(row.ParentId)),
		"status":        row.Status,
		"version":       0,
		"createdAt":     row.CreatedAt.Format("2006-01-02T15:04:05.000Z"),
		"updatedAt":     row.UpdatedAt.Format("2006-01-02T15:04:05.000Z"),
	}

	// 添加可选字段
	if row.Address != nil {
		teamInfo["address"] = *row.Address
	}
	if row.Extra != nil {
		// 解析JSON字段
		var extraData map[string]interface{}
		if err := json.Unmarshal([]byte(*row.Extra), &extraData); err == nil {
			teamInfo["extra"] = extraData
		}
	}

	zap.L().Debug("团队详情获取成功",
		zap.Uint32("teamId", teamID),
		zap.String("name", row.Name),
		zap.Any("meta", meta),
	)

	return teamInfo, nil
}

func (s *service) UpdateTeam(ctx Context, id string, payload *UpdateTeamPayload) (uint32, error) {
	// 验证团队ID格式
	teamIDUint, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		return 0, fmt.Errorf("团队ID格式错误：'%s'，必须是正整数", id)
	}
	if teamIDUint <= 0 {
		return 0, fmt.Errorf("团队ID错误：%d，必须大于0", teamIDUint)
	}
	teamID := uint32(teamIDUint)

	// 获取trace_id用于日志关联
	traceID := ctx.GetTraceID()

	// 创建日志元数据
	meta := map[string]interface{}{
		"trace_id":  traceID,
		"service":   "organization-service",
		"operation": "UpdateTeam",
		"teamId":    teamID,
		"name":      payload.Name,
		"status":    payload.Status,
		"version":   payload.Version,
	}

	zap.L().Info("开始更新团队",
		zap.Uint32("teamId", teamID),
		zap.String("name", payload.Name),
		zap.String("status", payload.Status),
		zap.Int32("version", payload.Version),
		zap.Any("meta", meta),
	)

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
		Table("org").Where("id = ? AND org_type = 'team' AND version = ?", teamID, payload.Version).
		Updates(updates)

	if result.Error != nil {
		zap.L().Error("更新团队失败",
			zap.Uint32("teamId", teamID),
			zap.Error(result.Error),
			zap.Any("meta", meta),
		)
		return 0, result.Error
	}

	if result.RowsAffected == 0 {
		zap.L().Warn("更新团队失败，可能是版本不匹配或记录不存在",
			zap.Uint32("teamId", teamID),
			zap.Int32("expectedVersion", payload.Version),
			zap.Any("meta", meta),
		)
		return 0, fmt.Errorf("team not found or version mismatch")
	}

	affectedRows := uint32(result.RowsAffected)
	zap.L().Info("团队更新成功",
		zap.Uint32("teamId", teamID),
		zap.Uint32("affectedRows", affectedRows),
		zap.Any("meta", meta),
	)

	return affectedRows, nil
}

func (s *service) DeleteTeam(ctx Context, id string) error {
	// 验证团队ID格式
	teamIDUint, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		return fmt.Errorf("团队ID格式错误：'%s'，必须是正整数", id)
	}
	if teamIDUint <= 0 {
		return fmt.Errorf("团队ID错误：%d，必须大于0", teamIDUint)
	}
	teamID := uint32(teamIDUint)

	// 获取trace_id用于日志关联
	traceID := ctx.GetTraceID()

	// 创建日志元数据
	meta := map[string]interface{}{
		"trace_id":  traceID,
		"service":   "organization-service",
		"operation": "DeleteTeam",
		"teamId":    teamID,
	}

	zap.L().Info("开始删除团队",
		zap.Uint32("teamId", teamID),
		zap.Any("meta", meta),
	)

	err = s.db.GetDbW().WithContext(ctx.RequestContext()).
		Table("org").Where("id = ? AND org_type = 'team'", teamID).Delete(&struct{}{}).Error

	if err != nil {
		zap.L().Error("删除团队失败",
			zap.Uint32("teamId", teamID),
			zap.Error(err),
			zap.Any("meta", meta),
		)
		return err
	}

	zap.L().Info("团队删除成功",
		zap.Uint32("teamId", teamID),
		zap.Any("meta", meta),
	)

	return nil
}

func (s *service) ListMembers(ctx Context, orgId string, current, pageSize int, keyword string) ([]map[string]interface{}, int64, error) {
	// 验证组织ID格式
	orgIDUint, err := strconv.ParseUint(orgId, 10, 32)
	if err != nil {
		return nil, 0, fmt.Errorf("组织ID格式错误：'%s'，必须是正整数", orgId)
	}
	if orgIDUint <= 0 {
		return nil, 0, fmt.Errorf("组织ID错误：%d，必须大于0", orgIDUint)
	}
	orgID := uint32(orgIDUint)

	// 获取trace_id用于日志关联
	traceID := ctx.GetTraceID()

	// 创建日志元数据
	meta := map[string]interface{}{
		"trace_id":  traceID,
		"service":   "organization-service",
		"operation": "ListMembers",
		"orgId":     orgID,
		"current":   current,
		"pageSize":  pageSize,
		"keyword":   keyword,
	}

	zap.L().Debug("开始获取团队成员列表",
		zap.Uint32("orgId", orgID),
		zap.Int("current", current),
		zap.Int("pageSize", pageSize),
		zap.String("keyword", keyword),
		zap.Any("meta", meta),
	)

	// 使用子查询方式获取团队成员 - 避免JOIN数据不一致问题
	// 1. 先获取属于该团队的账户ID列表
	var accountIDs []uint32
	subQuery := s.db.GetDbR().WithContext(ctx.RequestContext()).
		Table("account_org_relation").
		Select("account_id").
		Where("org_id = ? AND relation_type = 'belong' AND status = 'active'", orgID)

	if err := subQuery.Find(&accountIDs).Error; err != nil {
		zap.L().Error("获取团队成员账户ID列表失败",
			zap.Uint32("orgId", orgID),
			zap.Error(err),
			zap.Any("meta", meta),
		)
		return nil, 0, err
	}

	if len(accountIDs) == 0 {
		zap.L().Info("团队没有成员",
			zap.Uint32("orgId", orgID),
			zap.Any("meta", meta),
		)
		return []map[string]interface{}{}, 0, nil
	}

	zap.L().Debug("找到团队成员账户ID",
		zap.Uint32("orgId", orgID),
		zap.Int("accountCount", len(accountIDs)),
		zap.Any("accountIDs", accountIDs),
		zap.Any("meta", meta),
	)

	// 2. 根据账户ID列表查询账户信息
	db := s.db.GetDbR().WithContext(ctx.RequestContext()).
		Table("account").
		Where("id IN ?", accountIDs)

	// 添加关键词搜索
	if keyword != "" {
		db = db.Where("name LIKE ? OR username LIKE ? OR phone LIKE ?",
			"%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}

	// 统计总数
	var total int64
	if err := db.Count(&total).Error; err != nil {
		zap.L().Error("统计团队成员总数失败",
			zap.Uint32("orgId", orgID),
			zap.Error(err),
			zap.Any("meta", meta),
		)
		return nil, 0, err
	}

	zap.L().Debug("团队成员总数统计完成",
		zap.Uint32("orgId", orgID),
		zap.Int64("total", total),
		zap.Any("meta", meta),
	)

	// 查询成员数据
	var rows []struct {
		Id          uint32     `db:"id"`
		Username    string     `db:"username"`
		Name        string     `db:"name"`
		Phone       string     `db:"phone"`
		RoleType    string     `db:"role_type"`
		Status      string     `db:"status"`
		LastLoginAt *time.Time `db:"last_login_at"`
		CreatedAt   time.Time  `db:"created_at"`
		UpdatedAt   time.Time  `db:"updated_at"`
	}

	if err := db.Select("id, username, name, phone, role_type, status, last_login_at, created_at, updated_at").
		Order("id DESC").
		Limit(pageSize).
		Offset((current - 1) * pageSize).
		Find(&rows).Error; err != nil {
		zap.L().Error("查询团队成员数据失败",
			zap.Uint32("orgId", orgID),
			zap.Error(err),
			zap.Any("meta", meta),
		)
		return nil, 0, err
	}

	zap.L().Debug("团队成员数据查询完成",
		zap.Uint32("orgId", orgID),
		zap.Int("rowsCount", len(rows)),
		zap.Any("meta", meta),
	)

	// 转换为返回格式
	out := make([]map[string]interface{}, 0, len(rows))
	for _, r := range rows {
		member := map[string]interface{}{
			"id":        strconv.Itoa(int(r.Id)),
			"username":  r.Username,
			"name":      r.Name,
			"phone":     r.Phone,
			"roleType":  r.RoleType,
			"status":    r.Status,
			"createdAt": r.CreatedAt.Format("2006-01-02T15:04:05.000Z"),
			"updatedAt": r.UpdatedAt.Format("2006-01-02T15:04:05.000Z"),
		}

		// 添加最后登录时间
		if r.LastLoginAt != nil {
			member["lastLoginAt"] = r.LastLoginAt.Format("2006-01-02T15:04:05.000Z")
		} else {
			member["lastLoginAt"] = nil
		}

		out = append(out, member)
	}

	zap.L().Info("团队成员列表获取成功",
		zap.Uint32("orgId", orgID),
		zap.Int("resultCount", len(out)),
		zap.Int64("totalCount", total),
		zap.Any("meta", meta),
	)

	return out, total, nil
}

func (s *service) AddMember(ctx Context, orgId string, accountId string) error {
	// 验证组织ID格式
	orgIDUint, err := strconv.ParseUint(orgId, 10, 32)
	if err != nil {
		return fmt.Errorf("组织ID格式错误：'%s'，必须是正整数", orgId)
	}
	if orgIDUint <= 0 {
		return fmt.Errorf("组织ID错误：%d，必须大于0", orgIDUint)
	}
	orgID := uint32(orgIDUint)

	// 验证账户ID格式
	accountIDUint, err := strconv.ParseUint(accountId, 10, 32)
	if err != nil {
		return fmt.Errorf("账户ID格式错误：'%s'，必须是正整数", accountId)
	}
	if accountIDUint <= 0 {
		return fmt.Errorf("账户ID错误：%d，必须大于0", accountIDUint)
	}
	accountID := uint32(accountIDUint)

	// 获取trace_id用于日志关联
	traceID := ctx.GetTraceID()

	// 创建日志元数据
	meta := map[string]interface{}{
		"trace_id":  traceID,
		"service":   "organization-service",
		"operation": "AddMember",
		"orgId":     orgID,
		"accountId": accountID,
	}

	zap.L().Info("开始添加团队成员",
		zap.Uint32("orgId", orgID),
		zap.Uint32("accountId", accountID),
		zap.Any("meta", meta),
	)

	// 1. 验证团队ID的合法性
	var teamInfo struct {
		ID      uint32  `db:"id"`
		OrgType string  `db:"org_type"`
		Status  string  `db:"status"`
		Extra   *string `db:"extra"`
	}
	err = s.db.GetDbR().WithContext(ctx.RequestContext()).
		Table("org").
		Where("id = ? AND org_type = 'team'", orgID).
		Select("id, org_type, status, extra").
		First(&teamInfo).Error

	if err != nil {
		return NewMemberOperationError(
			"TEAM_NOT_FOUND",
			"团队不存在或不是有效的团队",
			fmt.Sprintf("团队ID '%d' 不存在或不是有效的团队", orgID),
		)
	}

	if teamInfo.Status != "enabled" {
		return NewMemberOperationError(
			"TEAM_DISABLED",
			"团队状态异常，无法添加成员",
			fmt.Sprintf("团队 '%d' 状态为 '%s'，无法添加成员", orgID, teamInfo.Status),
		)
	}

	// 2. 检查团队人数限制（从extra字段读取，如果有的话）
	var maxMembers *int32
	var currentCount int32

	if teamInfo.Extra != nil && *teamInfo.Extra != "" {
		var extraData map[string]interface{}
		if err := json.Unmarshal([]byte(*teamInfo.Extra), &extraData); err == nil {
			// 尝试读取 max_members
			if maxMembersVal, exists := extraData["max_members"]; exists {
				if maxMembersFloat, ok := maxMembersVal.(float64); ok {
					maxMembersInt := int32(maxMembersFloat)
					maxMembers = &maxMembersInt
				}
			}

			// 尝试读取 current_count
			if currentCountVal, exists := extraData["current_count"]; exists {
				if currentCountFloat, ok := currentCountVal.(float64); ok {
					currentCount = int32(currentCountFloat)
				}
			}
		}
	}

	// 如果extra中有max_members，则进行人数限制检查
	if maxMembers != nil && currentCount >= *maxMembers {
		return NewMemberOperationError(
			"TEAM_FULL",
			"团队人数已满，无法添加新成员",
			fmt.Sprintf("团队 '%d' 当前人数 %d，已达到最大人数限制 %d", orgID, currentCount, *maxMembers),
		)
	}

	// 2. 验证账户ID的合法性
	var accountInfo struct {
		ID       uint32 `db:"id"`
		Username string `db:"username"`
		Status   string `db:"status"`
	}
	err = s.db.GetDbR().WithContext(ctx.RequestContext()).
		Table("account").
		Where("id = ?", accountID).
		Select("id, username, status").
		First(&accountInfo).Error

	if err != nil {
		return NewMemberOperationError(
			"ACCOUNT_NOT_FOUND",
			"账户不存在",
			fmt.Sprintf("账户ID '%d' 不存在", accountID),
		)
	}

	if accountInfo.Status != "enabled" {
		return NewMemberOperationError(
			"ACCOUNT_DISABLED",
			"账户状态异常，无法添加到团队",
			fmt.Sprintf("账户 '%s' 状态为 '%s'，无法添加到团队", accountInfo.Username, accountInfo.Status),
		)
	}

	// 3. 检查是否已经是团队成员
	var existingRelation struct {
		ID     uint32 `db:"id"`
		Status string `db:"status"`
	}
	err = s.db.GetDbR().WithContext(ctx.RequestContext()).
		Table("account_org_relation").
		Where("account_id = ? AND org_id = ? AND relation_type = 'belong'", accountID, orgID).
		Select("id, status").
		First(&existingRelation).Error

	if err == nil {
		// 关系已存在
		if existingRelation.Status == "active" {
			return NewMemberOperationError(
				"MEMBER_ALREADY_EXISTS",
				"账户已经是团队成员",
				fmt.Sprintf("账户 '%s' 已经是团队 '%d' 的成员", accountInfo.Username, orgID),
			)
		} else if existingRelation.Status == "inactive" {
			// 关系存在但状态为inactive，可以重新激活
			zap.L().Info("重新激活团队成员关系",
				zap.Uint32("orgId", orgID),
				zap.Uint32("accountId", accountID),
				zap.String("username", accountInfo.Username),
				zap.Any("meta", meta),
			)
		}
	}

	// 开始事务
	tx := s.db.GetDbW().WithContext(ctx.RequestContext()).Begin()
	if tx.Error != nil {
		zap.L().Error("开始事务失败",
			zap.Error(tx.Error),
			zap.Any("meta", meta),
		)
		return tx.Error
	}

	// 确保事务回滚（如果出错）
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	// 4. 创建或更新团队成员关系
	var isNewMember bool
	if err == nil && existingRelation.Status == "inactive" {
		// 重新激活现有关系
		if err := tx.Table("account_org_relation").
			Where("account_id = ? AND org_id = ? AND relation_type = 'belong'", accountID, orgID).
			Update("status", "active").Error; err != nil {
			tx.Rollback()
			zap.L().Error("重新激活团队成员关系失败",
				zap.Uint32("orgId", orgID),
				zap.Uint32("accountId", accountID),
				zap.Error(err),
				zap.Any("meta", meta),
			)
			return err
		}
		isNewMember = false // 重新激活不算新成员
	} else {
		// 创建新的团队成员关系
		now := time.Now()
		relation := map[string]interface{}{
			"account_id":    accountID,
			"org_id":        orgID,
			"relation_type": "belong",
			"status":        "active",
			"created_at":    now,
			"updated_at":    now,
			"created_user":  ctx.SessionUserInfo().UserName,
			"updated_user":  ctx.SessionUserInfo().UserName,
		}

		if err := tx.Table("account_org_relation").Create(relation).Error; err != nil {
			tx.Rollback()
			zap.L().Error("创建团队成员关系失败",
				zap.Uint32("orgId", orgID),
				zap.Uint32("accountId", accountID),
				zap.Error(err),
				zap.Any("meta", meta),
			)
			return err
		}
		isNewMember = true // 新创建的成员
	}

	// 5. 如果是新成员，更新团队当前人数（从extra字段读取，如果有的话）
	if isNewMember {
		// 尝试更新extra字段中的current_count
		if teamInfo.Extra != nil && *teamInfo.Extra != "" {
			var extraData map[string]interface{}
			if err := json.Unmarshal([]byte(*teamInfo.Extra), &extraData); err == nil {
				// 如果extra中有current_count，则更新
				if _, exists := extraData["current_count"]; exists {
					extraData["current_count"] = currentCount + 1

					// 序列化回JSON
					if newExtraBytes, err := json.Marshal(extraData); err == nil {
						newExtra := string(newExtraBytes)

						// 更新数据库中的extra字段
						if err := tx.Table("org").
							Where("id = ?", orgID).
							Update("extra", newExtra).Error; err != nil {
							tx.Rollback()
							zap.L().Error("更新团队extra字段失败",
								zap.Uint32("orgId", orgID),
								zap.Error(err),
								zap.Any("meta", meta),
							)
							return err
						}

						zap.L().Info("团队人数更新成功",
							zap.Uint32("orgId", orgID),
							zap.Int32("oldCount", currentCount),
							zap.Int32("newCount", currentCount+1),
							zap.Any("meta", meta),
						)
					}
				}
			}
		}

		zap.L().Info("新成员添加成功",
			zap.Uint32("orgId", orgID),
			zap.Uint32("accountId", accountID),
			zap.Any("meta", meta),
		)
	}

	// 记录团队历史：添加成员操作
	// 构造操作内容
	operationContent := map[string]interface{}{
		"operation":  "add_member",
		"account_id": accountID,
		"action":     "added",
		"timestamp":  time.Now().Format("2006-01-02 15:04:05"),
	}

	// 使用通用history服务记录历史
	historyService := history.NewHistoryService(&historyAdapter{repo: s.db})
	err = historyService.CompareAndCreateHistory(
		&contextAdapter{ctx: ctx},
		history.EntityTypeOrg,          // 实体类型
		orgID,                          // 实体ID（团队ID）
		nil,                            // 旧对象（添加成员时不需要）
		operationContent,               // 新对象（操作内容）
		history.OperateTypeUpdated,     // 操作类型
		ctx.SessionUserInfo().UserName, // 操作人
		ctx.SessionUserInfo().RoleType, // 操作人角色
	)

	if err != nil {
		tx.Rollback()
		zap.L().Error("记录添加成员历史失败",
			zap.Error(err),
			zap.Uint32("orgId", orgID),
			zap.Uint32("accountId", accountID),
			zap.Any("meta", meta),
		)
		return err
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		zap.L().Error("提交事务失败",
			zap.Error(err),
			zap.Uint32("orgId", orgID),
			zap.Uint32("accountId", accountID),
			zap.Any("meta", meta),
		)
		return err
	}

	zap.L().Info("团队成员添加完成",
		zap.Uint32("orgId", orgID),
		zap.Uint32("accountId", accountID),
		zap.Any("meta", meta),
	)

	return nil
}

func (s *service) RemoveMember(ctx Context, orgId string, accountId string) error {
	// 验证组织ID格式
	orgIDUint, err := strconv.ParseUint(orgId, 10, 32)
	if err != nil {
		return fmt.Errorf("组织ID格式错误：'%s'，必须是正整数", orgId)
	}
	if orgIDUint <= 0 {
		return fmt.Errorf("组织ID错误：%d，必须大于0", orgIDUint)
	}
	orgID := uint32(orgIDUint)

	// 验证账户ID格式
	accountIDUint, err := strconv.ParseUint(accountId, 10, 32)
	if err != nil {
		return fmt.Errorf("账户ID格式错误：'%s'，必须是正整数", accountId)
	}
	if accountIDUint <= 0 {
		return fmt.Errorf("账户ID错误：%d，必须大于0", accountIDUint)
	}
	accountID := uint32(accountIDUint)

	// 获取trace_id用于日志关联
	traceID := ctx.GetTraceID()

	// 创建日志元数据
	meta := map[string]interface{}{
		"trace_id":  traceID,
		"service":   "organization-service",
		"operation": "RemoveMember",
		"orgId":     orgID,
		"accountId": accountID,
	}

	zap.L().Info("开始移除团队成员",
		zap.Uint32("orgId", orgID),
		zap.Uint32("accountId", accountID),
		zap.Any("meta", meta),
	)

	// 1. 验证团队ID的合法性
	var teamInfo struct {
		ID      uint32  `db:"id"`
		OrgType string  `db:"org_type"`
		Status  string  `db:"status"`
		Extra   *string `db:"extra"`
	}
	err = s.db.GetDbR().WithContext(ctx.RequestContext()).
		Table("org").
		Where("id = ? AND org_type = 'team'", orgID).
		Select("id, org_type, status, extra").
		First(&teamInfo).Error

	if err != nil {
		return NewMemberOperationError(
			"TEAM_NOT_FOUND",
			"团队不存在或不是有效的团队",
			fmt.Sprintf("团队ID '%d' 不存在或不是有效的团队", orgID),
		)
	}

	if teamInfo.Status != "enabled" {
		return NewMemberOperationError(
			"TEAM_DISABLED",
			"团队状态异常，无法移除成员",
			fmt.Sprintf("团队 '%d' 状态为 '%s'，无法移除成员", orgID, teamInfo.Status),
		)
	}

	// 2. 验证账户ID的合法性
	var accountInfo struct {
		ID       uint32 `db:"id"`
		Username string `db:"username"`
		Status   string `db:"status"`
	}
	err = s.db.GetDbR().WithContext(ctx.RequestContext()).
		Table("account").
		Where("id = ?", accountID).
		Select("id, username, status").
		First(&accountInfo).Error

	if err != nil {
		return NewMemberOperationError(
			"ACCOUNT_NOT_FOUND",
			"账户不存在",
			fmt.Sprintf("账户ID '%d' 不存在", accountID),
		)
	}

	// 3. 检查是否已经是团队成员
	var existingRelation struct {
		ID     uint32 `db:"id"`
		Status string `db:"status"`
	}
	err = s.db.GetDbR().WithContext(ctx.RequestContext()).
		Table("account_org_relation").
		Where("account_id = ? AND org_id = ? AND relation_type = 'belong'", accountID, orgID).
		Select("id, status").
		First(&existingRelation).Error

	if err != nil {
		return NewMemberOperationError(
			"MEMBER_NOT_FOUND",
			"账户不是团队成员",
			fmt.Sprintf("账户 '%s' 不是团队 '%d' 的成员", accountInfo.Username, orgID),
		)
	}

	if existingRelation.Status == "inactive" {
		return NewMemberOperationError(
			"MEMBER_ALREADY_REMOVED",
			"账户已经不是团队成员",
			fmt.Sprintf("账户 '%s' 已经不是团队 '%d' 的成员", accountInfo.Username, orgID),
		)
	}

	// 开始事务
	tx := s.db.GetDbW().WithContext(ctx.RequestContext()).Begin()
	if tx.Error != nil {
		zap.L().Error("开始事务失败",
			zap.Error(tx.Error),
			zap.Any("meta", meta),
		)
		return tx.Error
	}

	// 确保事务回滚（如果出错）
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	// 4. 将团队成员关系状态设置为inactive（软删除）
	result := tx.Table("account_org_relation").
		Where("account_id = ? AND org_id = ? AND relation_type = 'belong'", accountID, orgID).
		Update("status", "inactive")

	if result.Error != nil {
		tx.Rollback()
		zap.L().Error("移除团队成员失败",
			zap.Uint32("orgId", orgID),
			zap.Uint32("accountId", accountID),
			zap.Error(result.Error),
			zap.Any("meta", meta),
		)
		return result.Error
	}

	if result.RowsAffected == 0 {
		tx.Rollback()
		zap.L().Warn("未找到要移除的团队成员关系",
			zap.Uint32("orgId", orgID),
			zap.Uint32("accountId", accountID),
			zap.Any("meta", meta),
		)
		return NewMemberOperationError(
			"MEMBER_RELATION_NOT_FOUND",
			"团队成员关系不存在",
			fmt.Sprintf("团队 '%d' 与账户 '%d' 的成员关系不存在", orgID, accountID),
		)
	}

	// 5. 更新团队当前人数（从extra字段读取，如果有的话）
	// 尝试更新extra字段中的current_count
	if teamInfo.Extra != nil && *teamInfo.Extra != "" {
		var extraData map[string]interface{}
		if err := json.Unmarshal([]byte(*teamInfo.Extra), &extraData); err == nil {
			// 如果extra中有current_count，则更新
			if currentCountVal, exists := extraData["current_count"]; exists {
				if currentCountFloat, ok := currentCountVal.(float64); ok {
					currentCount := int32(currentCountFloat)
					extraData["current_count"] = currentCount - 1

					// 序列化回JSON
					if newExtraBytes, err := json.Marshal(extraData); err == nil {
						newExtra := string(newExtraBytes)

						// 更新数据库中的extra字段
						if err := tx.Table("org").
							Where("id = ?", orgID).
							Update("extra", newExtra).Error; err != nil {
							tx.Rollback()
							zap.L().Error("更新团队extra字段失败",
								zap.Uint32("orgId", orgID),
								zap.Error(err),
								zap.Any("meta", meta),
							)
							return err
						}

						zap.L().Info("团队人数更新成功",
							zap.Uint32("orgId", orgID),
							zap.Int32("oldCount", currentCount),
							zap.Int32("newCount", currentCount-1),
							zap.Any("meta", meta),
						)
					}
				}
			}
		}
	}

	zap.L().Info("成员移除成功",
		zap.Uint32("orgId", orgID),
		zap.Uint32("accountId", accountID),
		zap.Any("meta", meta),
	)

	// 记录团队历史：移除成员操作
	// 构造操作内容
	operationContent := map[string]interface{}{
		"operation":  "remove_member",
		"account_id": accountID,
		"action":     "removed",
		"timestamp":  time.Now().Format("2006-01-02 15:04:05"),
	}

	// 使用通用history服务记录历史
	historyService := history.NewHistoryService(&historyAdapter{repo: s.db})
	err = historyService.CompareAndCreateHistory(
		&contextAdapter{ctx: ctx},
		history.EntityTypeOrg,          // 实体类型
		orgID,                          // 实体ID（团队ID）
		nil,                            // 旧对象（移除成员时不需要）
		operationContent,               // 新对象（操作内容）
		history.OperateTypeUpdated,     // 操作类型
		ctx.SessionUserInfo().UserName, // 操作人
		ctx.SessionUserInfo().RoleType, // 操作人角色
	)

	if err != nil {
		tx.Rollback()
		zap.L().Error("记录移除成员历史失败",
			zap.Error(err),
			zap.Uint32("orgId", orgID),
			zap.Uint32("accountId", accountID),
			zap.Any("meta", meta),
		)
		return err
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		zap.L().Error("提交事务失败",
			zap.Error(err),
			zap.Uint32("orgId", orgID),
			zap.Uint32("accountId", accountID),
			zap.Any("meta", meta),
		)
		return err
	}

	zap.L().Info("团队成员移除完成",
		zap.Uint32("orgId", orgID),
		zap.Uint32("accountId", accountID),
		zap.Int64("rowsAffected", result.RowsAffected),
		zap.Any("meta", meta),
	)

	return nil
}

// UpdateMemberRole 更新团队成员角色
func (s *service) UpdateMemberRole(ctx Context, orgId string, accountId string, roleType string) error {
	// 验证组织ID格式
	orgIDUint, err := strconv.ParseUint(orgId, 10, 32)
	if err != nil {
		return fmt.Errorf("组织ID格式错误：'%s'，必须是正整数", orgId)
	}
	if orgIDUint <= 0 {
		return fmt.Errorf("组织ID错误：%d，必须大于0", orgIDUint)
	}
	orgID := uint32(orgIDUint)

	// 验证账户ID格式
	accountIDUint, err := strconv.ParseUint(accountId, 10, 32)
	if err != nil {
		return fmt.Errorf("账户ID格式错误：'%s'，必须是正整数", accountId)
	}
	if accountIDUint <= 0 {
		return fmt.Errorf("账户ID错误：%d，必须大于0", accountIDUint)
	}
	accountID := uint32(accountIDUint)

	// 验证角色类型
	if roleType != "employee" && roleType != "team_manager" {
		return fmt.Errorf("无效的角色类型：'%s'，必须是 'employee' 或 'team_manager'", roleType)
	}

	// 获取trace_id用于日志关联
	traceID := ctx.GetTraceID()

	// 创建日志元数据
	meta := map[string]interface{}{
		"trace_id":  traceID,
		"service":   "organization-service",
		"operation": "UpdateMemberRole",
		"orgId":     orgID,
		"accountId": accountID,
		"roleType":  roleType,
	}

	zap.L().Info("开始更新团队成员角色",
		zap.Uint32("orgId", orgID),
		zap.Uint32("accountId", accountID),
		zap.String("roleType", roleType),
		zap.Any("meta", meta),
	)

	// 1. 验证团队ID的合法性
	var teamInfo struct {
		ID      uint32 `db:"id"`
		OrgType string `db:"org_type"`
		Status  string `db:"status"`
	}
	err = s.db.GetDbR().WithContext(ctx.RequestContext()).
		Table("org").
		Where("id = ? AND org_type = 'team'", orgID).
		Select("id, org_type, status").
		First(&teamInfo).Error

	if err != nil {
		return NewMemberOperationError(
			"TEAM_NOT_FOUND",
			"团队不存在或不是有效的团队",
			fmt.Sprintf("团队ID '%d' 不存在或不是有效的团队", orgID),
		)
	}

	if teamInfo.Status != "enabled" {
		return NewMemberOperationError(
			"TEAM_DISABLED",
			"团队状态异常，无法更新成员角色",
			fmt.Sprintf("团队 '%d' 状态为 '%s'，无法更新成员角色", orgID, teamInfo.Status),
		)
	}

	// 2. 验证账户ID的合法性
	var accountInfo struct {
		ID       uint32 `db:"id"`
		Username string `db:"username"`
		Status   string `db:"status"`
	}
	err = s.db.GetDbR().WithContext(ctx.RequestContext()).
		Table("account").
		Where("id = ?", accountID).
		Select("id, username, status").
		First(&accountInfo).Error

	if err != nil {
		return NewMemberOperationError(
			"ACCOUNT_NOT_FOUND",
			"账户不存在",
			fmt.Sprintf("账户ID '%d' 不存在", accountID),
		)
	}

	if accountInfo.Status != "enabled" {
		return NewMemberOperationError(
			"ACCOUNT_DISABLED",
			"账户状态异常，无法更新角色",
			fmt.Sprintf("账户 '%s' 状态为 '%s'，无法更新角色", accountInfo.Username, accountInfo.Status),
		)
	}

	// 3. 检查是否已经是团队成员
	var existingRelation struct {
		ID     uint32 `db:"id"`
		Status string `db:"status"`
	}
	err = s.db.GetDbR().WithContext(ctx.RequestContext()).
		Table("account_org_relation").
		Where("account_id = ? AND org_id = ? AND relation_type = 'belong'", accountID, orgID).
		Select("id, status").
		First(&existingRelation).Error

	if err != nil {
		return NewMemberOperationError(
			"MEMBER_NOT_FOUND",
			"账户不是团队成员",
			fmt.Sprintf("账户 '%s' 不是团队 '%d' 的成员", accountInfo.Username, orgID),
		)
	}

	if existingRelation.Status == "inactive" {
		return NewMemberOperationError(
			"MEMBER_ALREADY_REMOVED",
			"账户已经不是团队成员",
			fmt.Sprintf("账户 '%s' 已经不是团队 '%d' 的成员", accountInfo.Username, orgID),
		)
	}

	// 4. 获取当前角色信息
	var currentRole struct {
		RoleType string `db:"role_type"`
	}
	err = s.db.GetDbR().WithContext(ctx.RequestContext()).
		Table("account").
		Where("id = ?", accountID).
		Select("role_type").
		First(&currentRole).Error

	if err != nil {
		return NewMemberOperationError(
			"ACCOUNT_ROLE_QUERY_FAILED",
			"获取账户角色信息失败",
			fmt.Sprintf("获取账户角色信息失败：%v", err),
		)
	}

	// 5. 检查角色是否已经符合预期
	if currentRole.RoleType == roleType {
		return NewMemberOperationError(
			"ROLE_ALREADY_MATCHES",
			"角色已经是目标角色，无需更新",
			fmt.Sprintf("账户 '%s' 的角色已经是 '%s'，无需更新", accountInfo.Username, roleType),
		)
	}

	// 开始事务
	tx := s.db.GetDbW().WithContext(ctx.RequestContext()).Begin()
	if tx.Error != nil {
		zap.L().Error("开始事务失败",
			zap.Error(tx.Error),
			zap.Any("meta", meta),
		)
		return tx.Error
	}

	// 确保事务回滚（如果出错）
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	// 更新账户角色
	result := tx.Table("account").
		Where("id = ?", accountID).
		Update("role_type", roleType)

	if result.Error != nil {
		tx.Rollback()
		zap.L().Error("更新账户角色失败",
			zap.Uint32("accountId", accountID),
			zap.Error(result.Error),
			zap.Any("meta", meta),
		)
		return result.Error
	}

	if result.RowsAffected == 0 {
		tx.Rollback()
		zap.L().Warn("未找到要更新的账户",
			zap.Uint32("accountId", accountID),
			zap.Any("meta", meta),
		)
		return fmt.Errorf("账户不存在")
	}

	// 记录团队历史：更新成员角色操作
	// 构造操作内容
	operationContent := map[string]interface{}{
		"operation":  "update_member_role",
		"account_id": accountID,
		"old_role":   currentRole.RoleType,
		"new_role":   roleType,
		"action":     "role_changed",
		"timestamp":  time.Now().Format("2006-01-02 15:04:05"),
	}

	// 使用通用history服务记录历史
	historyService := history.NewHistoryService(&historyAdapter{repo: s.db})
	err = historyService.CompareAndCreateHistory(
		&contextAdapter{ctx: ctx},
		history.EntityTypeOrg,          // 实体类型
		orgID,                          // 实体ID（团队ID）
		nil,                            // 旧对象（更新角色时不需要）
		operationContent,               // 新对象（操作内容）
		history.OperateTypeUpdated,     // 操作类型
		ctx.SessionUserInfo().UserName, // 操作人
		ctx.SessionUserInfo().RoleType, // 操作人角色
	)

	if err != nil {
		tx.Rollback()
		zap.L().Error("记录更新成员角色历史失败",
			zap.Error(err),
			zap.Uint32("orgId", orgID),
			zap.Uint32("accountId", accountID),
			zap.Any("meta", meta),
		)
		return err
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		zap.L().Error("提交事务失败",
			zap.Error(err),
			zap.Uint32("orgId", orgID),
			zap.Uint32("accountId", accountID),
			zap.Any("meta", meta),
		)
		return err
	}

	zap.L().Info("团队成员角色更新完成",
		zap.Uint32("orgId", orgID),
		zap.Uint32("accountId", accountID),
		zap.String("oldRole", currentRole.RoleType),
		zap.String("newRole", roleType),
		zap.Int64("rowsAffected", result.RowsAffected),
		zap.Any("meta", meta),
	)

	return nil
}
