package account

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/xinliangnote/go-gin-api/internal/pkg/core"
	"github.com/xinliangnote/go-gin-api/internal/repository/mysql"
	"github.com/xinliangnote/go-gin-api/internal/repository/mysql/account"
	"github.com/xinliangnote/go-gin-api/internal/repository/mysql/account_org_relation"
	"github.com/xinliangnote/go-gin-api/internal/repository/mysql/org"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// Create 创建账户
func (s *service) Create(ctx core.Context, accountData *CreateAccountData) (id int32, err error) {
	// 检查用户名是否已存在
	accountQueryBuilder := account.NewQueryBuilder()
	accountQueryBuilder.WhereUsername(mysql.EqualPredicate, accountData.Username)

	existingAccount, err := accountQueryBuilder.QueryOne(s.db.GetDbR())
	if err != nil {
		return 0, fmt.Errorf("查询账户失败，请稍后重试: %v", err)
	}
	if existingAccount != nil {
		return 0, fmt.Errorf("用户名 '%s' 已存在，请使用其他用户名", accountData.Username)
	}

	// 验证必填字段：归属组必须填写
	if accountData.BelongGroup == nil || accountData.BelongGroup.ID == "" {
		return 0, fmt.Errorf("归属组信息缺失，请提供有效的归属组ID")
	}

	// 验证归属组ID格式：必须是正整数
	belongGroupID, err := strconv.ParseUint(accountData.BelongGroup.ID, 10, 32)
	if err != nil {
		return 0, fmt.Errorf("归属组ID格式错误：'%s'，必须是正整数", accountData.BelongGroup.ID)
	}
	if belongGroupID <= 0 {
		return 0, fmt.Errorf("归属组ID错误：%d，必须大于0", belongGroupID)
	}

	// 检查归属组是否存在
	var existingOrg org.Org
	err = s.db.GetDbR().Where("id = ? AND org_type = 'group'", belongGroupID).First(&existingOrg).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return 0, fmt.Errorf("归属组不存在，ID: %d，请检查组ID是否正确或联系管理员创建该组", belongGroupID)
		}
		return 0, fmt.Errorf("查询归属组失败，请稍后重试: %v", err)
	}

	// 加密密码 - 使用bcrypt格式，与mock数据保持一致
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(accountData.Password), bcrypt.DefaultCost)
	if err != nil {
		return 0, fmt.Errorf("密码加密失败: %v", err)
	}

	// 设置默认值：角色类型默认为普通员工
	if accountData.RoleType == "" {
		accountData.RoleType = "employee"
	}

	if accountData.Status == "" {
		accountData.Status = "enabled"
	}

	// 在一个事务中执行：创建账户 + 创建关联关系 + 写入历史记录
	err = s.db.GetDbW().WithContext(ctx.RequestContext()).Transaction(func(tx *gorm.DB) error {
		// 1. 创建账户记录
		now := time.Now()
		newAccount := &account.Account{
			Username:      accountData.Username,
			Name:          accountData.Name,
			Password:      string(hashedPassword), // bcrypt生成的是[]byte，需要转换为string
			Phone:         accountData.Phone,
			RoleType:      accountData.RoleType,
			Status:        accountData.Status,
			BelongGroupId: uint32(belongGroupID), // 必填
			BelongTeamId:  0,                     // 不允许填写归属小队，默认为0
			CreatedAt:     now,
			UpdatedAt:     now,
			CreatedUser:   ctx.SessionUserInfo().UserName,
			UpdatedUser:   ctx.SessionUserInfo().UserName,
		}

		// 保存到数据库
		createdId, err := newAccount.Create(tx)
		if err != nil {
			return fmt.Errorf("创建账户记录失败，数据库操作异常: %v", err)
		}
		id = int32(createdId)

		// 2. 写入账户与组织的关联关系（belong）
		// 2.1 创建与组的关联关系
		if accountData.BelongGroup != nil && belongGroupID > 0 {
			rel := &account_org_relation.AccountOrgRelation{
				AccountId:    uint32(id),
				OrgId:        uint32(belongGroupID),
				RelationType: "belong",
				Status:       "active",
				CreatedAt:    now,
				UpdatedAt:    now,
				CreatedUser:  ctx.SessionUserInfo().UserName,
				UpdatedUser:  ctx.SessionUserInfo().UserName,
			}
			if _, err := rel.Create(tx); err != nil {
				return fmt.Errorf("创建账户与组 '%d' 的关联关系失败: %v", belongGroupID, err)
			}
		}

		// 2.2 创建与公司的关联关系（默认公司ID 1001）
		companyRel := &account_org_relation.AccountOrgRelation{
			AccountId:    uint32(id),
			OrgId:        1001, // 默认公司ID
			RelationType: "belong",
			Status:       "active",
			CreatedAt:    now,
			UpdatedAt:    now,
			CreatedUser:  ctx.SessionUserInfo().UserName,
			UpdatedUser:  ctx.SessionUserInfo().UserName,
		}
		if _, err := companyRel.Create(tx); err != nil {
			return fmt.Errorf("创建账户与公司 '1001' 的关联关系失败: %v", err)
		}

		// 3. 创建历史记录
		historyContent := map[string]interface{}{
			"username": map[string]string{"old": "", "new": accountData.Username},
			"name":     map[string]string{"old": "", "new": accountData.Name},
			"phone":    map[string]string{"old": "", "new": accountData.Phone},
			"roleType": map[string]string{"old": "", "new": accountData.RoleType},
			"status":   map[string]string{"old": "", "new": accountData.Status},
			"belongGroup": map[string]interface{}{"old": "", "new": map[string]interface{}{
				"id":   belongGroupID,
				"name": accountData.BelongGroup.Name,
			}},
			"belongCompany": map[string]interface{}{"old": "", "new": map[string]interface{}{
				"id":   1001,
				"name": "默认公司",
			}},
		}

		contentBytes, _ := json.Marshal(historyContent)

		historyData := &CreateHistoryData{
			AccountId:        fmt.Sprintf("%d", id), // 使用数据库生成的ID
			OperateType:      "created",
			Content:          string(contentBytes),
			OperatorUsername: ctx.SessionUserInfo().UserName, // 操作人用户名
			OperatorName:     ctx.SessionUserInfo().UserName, // 操作人姓名
			OperatorRoleType: ctx.SessionUserInfo().RoleType, // 操作人角色类型
		}

		if err := s.CreateHistory(ctx, historyData); err != nil {
			return fmt.Errorf("创建账户历史记录失败，无法记录操作日志: %v", err)
		}

		return nil
	})

	if err != nil {
		return 0, err
	}

	return id, nil
}
