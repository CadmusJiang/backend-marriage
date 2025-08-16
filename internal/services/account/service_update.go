package account

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/xinliangnote/go-gin-api/internal/authz"
	"github.com/xinliangnote/go-gin-api/internal/pkg/core"
	"github.com/xinliangnote/go-gin-api/internal/pkg/password"
	"github.com/xinliangnote/go-gin-api/internal/repository/mysql"
	"github.com/xinliangnote/go-gin-api/internal/repository/mysql/account"
	"github.com/xinliangnote/go-gin-api/internal/repository/mysql/account_history"
	"github.com/xinliangnote/go-gin-api/internal/repository/mysql/account_org_relation"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// Update 更新账户
func (s *service) Update(ctx core.Context, accountId string, updateData *UpdateAccountData) (err error) {
	// 转换accountId为int32
	var id int32
	fmt.Sscanf(accountId, "%d", &id)

	// 查询账户是否存在
	accountQueryBuilder := account.NewQueryBuilder()
	accountQueryBuilder.WhereId(mysql.EqualPredicate, id)

	existingAccount, err := accountQueryBuilder.QueryOne(s.db.GetDbR())
	if err != nil {
		return fmt.Errorf("查询账户失败: %v", err)
	}
	if existingAccount == nil {
		return fmt.Errorf("账户不存在")
	}

	// 范围校验：无权限则拒绝
	scope, scopeErr := authz.ComputeScope(ctx, s.db)
	if scopeErr == nil && !authz.CanAccessAccount(scope, existingAccount) {
		return fmt.Errorf("无权限更新该账户")
	}

	// 准备更新数据
	updateFields := make(map[string]interface{})
	updateFields["updated_at"] = uint64(time.Now().Unix())
	updateFields["updated_user"] = ctx.SessionUserInfo().UserName

	// 记录变更内容用于历史记录
	changes := make(map[string]interface{})

	if updateData.Name != "" && updateData.Name != existingAccount.Name {
		updateFields["name"] = updateData.Name
		changes["name"] = map[string]string{
			"old": existingAccount.Name,
			"new": updateData.Name,
		}
	}

	if updateData.Phone != "" && updateData.Phone != existingAccount.Phone {
		updateFields["phone"] = updateData.Phone
		changes["phone"] = map[string]string{
			"old": existingAccount.Phone,
			"new": updateData.Phone,
		}
	}

	// 状态变更
	if updateData.Status != "" {
		var newStatus string
		if updateData.Status == "enabled" {
			newStatus = "enabled"
		} else {
			newStatus = "disabled"
		}
		if newStatus != existingAccount.Status {
			updateFields["status"] = newStatus
			changes["status"] = map[string]string{
				"old": existingAccount.Status,
				"new": updateData.Status,
			}
		}
	}

	// 更新组织信息
	if updateData.BelongGroup != nil {
		updateFields["belong_group_id"] = updateData.BelongGroup.ID

		// 更新账户-组关联：先删除原 belong 关联，再插入新的
		qbRel := account_org_relation.NewQueryBuilder()
		qbRel.WhereAccountId(mysql.EqualPredicate, uint64(id))
		qbRel.WhereRelationType(mysql.EqualPredicate, "belong")
		_ = qbRel.Delete(s.db.GetDbW())

		rel := &account_org_relation.AccountOrgRelation{
			AccountId:    uint32(id),
			OrgId:        uint32(updateData.BelongGroup.ID),
			RelationType: "belong",
			Status:       "active",
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
			CreatedUser:  ctx.SessionUserInfo().UserName,
			UpdatedUser:  ctx.SessionUserInfo().UserName,
		}
		if _, err := rel.Create(s.db.GetDbW()); err != nil {
			ctx.Logger().Error("更新账户-组关联失败", zap.Error(err))
		}
	}

	if updateData.BelongTeam != nil {
		updateFields["belong_team_id"] = updateData.BelongTeam.ID

		// 更新账户-团队关联：先删除原 belong 关联，再插入新的
		qbRel := account_org_relation.NewQueryBuilder()
		qbRel.WhereAccountId(mysql.EqualPredicate, uint64(id))
		qbRel.WhereRelationType(mysql.EqualPredicate, "belong")
		_ = qbRel.Delete(s.db.GetDbW())

		rel := &account_org_relation.AccountOrgRelation{
			AccountId:    uint32(id),
			OrgId:        uint32(updateData.BelongTeam.ID),
			RelationType: "belong",
			Status:       "active",
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
			CreatedUser:  ctx.SessionUserInfo().UserName,
			UpdatedUser:  ctx.SessionUserInfo().UserName,
		}
		if _, err := rel.Create(s.db.GetDbW()); err != nil {
			ctx.Logger().Error("更新账户-团队关联失败", zap.Error(err))
		}
	}

	// 在一个事务中执行：更新账户 + 写入历史
	err = s.db.GetDbW().WithContext(ctx.RequestContext()).Transaction(func(tx *gorm.DB) error {
		// 乐观锁：如果表有 version，可在 updateFields 中自增并带 where version
		// 这里先直接更新
		if err := accountQueryBuilder.Updates(tx, updateFields); err != nil {
			return fmt.Errorf("更新账户失败: %v", err)
		}

		if len(changes) > 0 {
			contentBytes, _ := json.Marshal(changes)
			now := time.Now()
			accIdUint, _ := fmt.Sscanf(accountId, "%d", &id)
			_ = accIdUint
			hist := &account_history.AccountHistory{
				AccountId:        uint64(id),
				OperateType:      "modified",
				OperatedAt:       now,
				Content:          string(contentBytes),
				Operator:         ctx.SessionUserInfo().UserName,
				OperatorRoleType: scope.RoleType,
				CreatedAt:        now,
				UpdatedAt:        now,
				CreatedUser:      ctx.SessionUserInfo().UserName,
				UpdatedUser:      ctx.SessionUserInfo().UserName,
			}
			if _, err := hist.Create(tx); err != nil {
				return fmt.Errorf("写入账户历史失败: %v", err)
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (s *service) UpdatePassword(ctx core.Context, accountId string, newPassword string) (err error) {
	var id int32
	fmt.Sscanf(accountId, "%d", &id)
	qb := account.NewQueryBuilder()
	qb.WhereId(mysql.EqualPredicate, id)
	hashed := password.GeneratePassword(newPassword)
	return qb.Updates(s.db.GetDbW().WithContext(ctx.RequestContext()), map[string]interface{}{"password": hashed})
}

// Delete 删除账户
func (s *service) Delete(ctx core.Context, accountId string) (err error) {
	// 转换accountId为int32
	var id int32
	fmt.Sscanf(accountId, "%d", &id)

	// 查询账户是否存在
	accountQueryBuilder := account.NewQueryBuilder()
	accountQueryBuilder.WhereId(mysql.EqualPredicate, id)

	existingAccount, err := accountQueryBuilder.QueryOne(s.db.GetDbR())
	if err != nil {
		return fmt.Errorf("查询账户失败: %v", err)
	}
	if existingAccount == nil {
		return fmt.Errorf("账户不存在")
	}

	// 范围校验
	scope, scopeErr := authz.ComputeScope(ctx, s.db)
	if scopeErr == nil && !authz.CanAccessAccount(scope, existingAccount) {
		return fmt.Errorf("无权限删除该账户")
	}

	// 软删除账户
	updateFields := map[string]interface{}{
		"is_deleted":         1,
		"modified_timestamp": time.Now().Unix(),
		"updated_user":       ctx.SessionUserInfo().UserName,
	}

	err = accountQueryBuilder.Updates(s.db.GetDbW(), updateFields)
	if err != nil {
		return fmt.Errorf("删除账户失败: %v", err)
	}

	// 创建删除历史记录
	historyData := &CreateHistoryData{
		AccountId:        accountId,
		OperateType:      "deleted",
		Content:          `{"action": "deleted"}`,
		Operator:         ctx.SessionUserInfo().UserName,
		OperatorRoleType: "admin", // 暂时硬编码
	}

	err = s.CreateHistory(ctx, historyData)
	if err != nil {
		// 记录日志但不影响主流程
		ctx.Logger().Error("创建历史记录失败", zap.Error(err))
	}

	return nil
}
