package account

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/xinliangnote/go-gin-api/internal/pkg/core"
	"github.com/xinliangnote/go-gin-api/internal/repository/mysql"
	"github.com/xinliangnote/go-gin-api/internal/repository/mysql/account"
	"github.com/xinliangnote/go-gin-api/internal/repository/mysql/account_org_relation"
	"go.uber.org/zap"
)

// Create 创建账户
func (s *service) Create(ctx core.Context, accountData *CreateAccountData) (id int32, err error) {
	// 检查用户名是否已存在
	accountQueryBuilder := account.NewQueryBuilder()
	accountQueryBuilder.WhereUsername(mysql.EqualPredicate, accountData.Username)

	existingAccount, err := accountQueryBuilder.QueryOne(s.db.GetDbR())
	if err != nil {
		return 0, fmt.Errorf("查询账户失败: %v", err)
	}
	if existingAccount != nil {
		return 0, fmt.Errorf("用户名已存在")
	}

	// 加密密码
	hashedPassword := s.generateMD5(accountData.Password)

	// 设置默认值
	if accountData.RoleType == "" {
		accountData.RoleType = "employee"
	}
	if accountData.Status == "" {
		accountData.Status = "enabled"
	}

	// 创建账户记录
	now := int64(time.Now().Unix())
	newAccount := &account.Account{
		Username:    accountData.Username,
		Nickname:    accountData.Nickname,
		Password:    hashedPassword,
		Phone:       accountData.Phone,
		RoleType:    accountData.RoleType,
		Status:      accountData.Status,
		CreatedAt:   now,
		UpdatedAt:   now,
		CreatedUser: ctx.SessionUserInfo().UserName,
		UpdatedUser: ctx.SessionUserInfo().UserName,
	}

	// 设置组织信息
	if accountData.BelongGroup != nil {
		newAccount.BelongGroupId = int32(accountData.BelongGroup.ID)
	}

	if accountData.BelongTeam != nil {
		newAccount.BelongTeamId = int32(accountData.BelongTeam.ID)
	}

	// 保存到数据库
	id, err = newAccount.Create(s.db.GetDbW())
	if err != nil {
		return 0, fmt.Errorf("创建账户失败: %v", err)
	}

	// 写入账户与组织的关联关系（belong）
	// 关联类型: 1 belong, 2 manage；状态: 1 active
	if accountData.BelongGroup != nil && accountData.BelongGroup.ID > 0 {
		rel := &account_org_relation.AccountOrgRelation{
			AccountId:         uint32(id),
			OrgId:             uint32(accountData.BelongGroup.ID),
			RelationType:      1,
			Status:            1,
			CreatedTimestamp:  now,
			ModifiedTimestamp: now,
			CreatedUser:       ctx.SessionUserInfo().UserName,
			UpdatedUser:       ctx.SessionUserInfo().UserName,
		}
		if _, err := rel.Create(s.db.GetDbW()); err != nil {
			ctx.Logger().Error("创建账户-组关联失败", zap.Error(err))
		}
	}

	if accountData.BelongTeam != nil && accountData.BelongTeam.ID > 0 {
		rel := &account_org_relation.AccountOrgRelation{
			AccountId:         uint32(id),
			OrgId:             uint32(accountData.BelongTeam.ID),
			RelationType:      1,
			Status:            1,
			CreatedTimestamp:  now,
			ModifiedTimestamp: now,
			CreatedUser:       ctx.SessionUserInfo().UserName,
			UpdatedUser:       ctx.SessionUserInfo().UserName,
		}
		if _, err := rel.Create(s.db.GetDbW()); err != nil {
			ctx.Logger().Error("创建账户-团队关联失败", zap.Error(err))
		}
	}

	// 创建历史记录
	historyContent := map[string]interface{}{
		"username": map[string]string{"old": "", "new": accountData.Username},
		"nickname": map[string]string{"old": "", "new": accountData.Nickname},
		"phone":    map[string]string{"old": "", "new": accountData.Phone},
		"roleType": map[string]string{"old": "", "new": accountData.RoleType},
		"status":   map[string]string{"old": "", "new": accountData.Status},
	}

	contentBytes, _ := json.Marshal(historyContent)

	historyData := &CreateHistoryData{
		AccountId:        fmt.Sprintf("%d", id), // 使用数据库生成的ID
		OperateType:      "created",
		Content:          string(contentBytes),
		Operator:         ctx.SessionUserInfo().UserName,
		OperatorRoleType: "admin", // 暂时硬编码
	}

	err = s.CreateHistory(ctx, historyData)
	if err != nil {
		// 记录日志但不影响主流程
		ctx.Logger().Error("创建历史记录失败", zap.Error(err))
	}

	return id, nil
}
