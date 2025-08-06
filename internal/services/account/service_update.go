package account

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/xinliangnote/go-gin-api/internal/pkg/core"
	"github.com/xinliangnote/go-gin-api/internal/repository/mysql"
	"github.com/xinliangnote/go-gin-api/internal/repository/mysql/account"
	"go.uber.org/zap"
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

	// 准备更新数据
	updateFields := make(map[string]interface{})
	updateFields["modified_timestamp"] = uint64(time.Now().Unix())
	updateFields["updated_user"] = ctx.SessionUserInfo().UserName

	// 记录变更内容用于历史记录
	changes := make(map[string]interface{})

	if updateData.Nickname != "" && updateData.Nickname != existingAccount.Nickname {
		updateFields["nickname"] = updateData.Nickname
		changes["nickname"] = map[string]string{
			"old": existingAccount.Nickname,
			"new": updateData.Nickname,
		}
	}

	if updateData.Phone != "" && updateData.Phone != existingAccount.Phone {
		updateFields["phone"] = updateData.Phone
		changes["phone"] = map[string]string{
			"old": existingAccount.Phone,
			"new": updateData.Phone,
		}
	}

	if updateData.Status != "" && updateData.Status != existingAccount.Status {
		updateFields["status"] = updateData.Status
		changes["status"] = map[string]string{
			"old": existingAccount.Status,
			"new": updateData.Status,
		}
	}

	// 更新组织信息
	if updateData.BelongGroup != nil {
		updateFields["belong_group_id"] = updateData.BelongGroup.ID
	}

	if updateData.BelongTeam != nil {
		updateFields["belong_team_id"] = updateData.BelongTeam.ID
	}

	// 执行更新
	err = accountQueryBuilder.Updates(s.db.GetDbW(), updateFields)
	if err != nil {
		return fmt.Errorf("更新账户失败: %v", err)
	}

	// 如果有变更，创建历史记录
	if len(changes) > 0 {
		contentBytes, _ := json.Marshal(changes)

		historyData := &CreateHistoryData{
			AccountId:        accountId,
			OperateType:      "modified",
			Content:          string(contentBytes),
			Operator:         ctx.SessionUserInfo().UserName,
			OperatorRoleType: "admin", // 暂时硬编码
		}

		err = s.CreateHistory(ctx, historyData)
		if err != nil {
			// 记录日志但不影响主流程
			ctx.Logger().Error("创建历史记录失败", zap.Error(err))
		}
	}

	return nil
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
