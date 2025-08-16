package account_history

import "time"

// AccountHistory 账户历史记录表
//
//go:generate gormgen -structs AccountHistory -input .
type AccountHistory struct {
	Id uint64 `gorm:"primaryKey"` // 主键

	// 基本信息
	AccountId        uint64    `gorm:"column:account_id;not null"`        // 账户ID
	OperateType      string    `gorm:"column:operate_type;size:20"`       // 操作类型: created, updated
	OperatedAt       time.Time `gorm:"column:operated_at;not null"`       // 操作时间
	Content          string    `gorm:"column:content;type:json"`          // 操作内容 (JSON格式)
	OperatorUsername string    `gorm:"column:operator_username;size:32"`  // 操作人用户名
	OperatorName     string    `gorm:"column:operator_name;size:60"`      // 操作人姓名
	OperatorRoleType string    `gorm:"column:operator_role_type;size:20"` // 操作人角色类型

	// 审计字段
	CreatedAt   time.Time `gorm:"column:created_at;not null"`  // 创建时间
	UpdatedAt   time.Time `gorm:"column:updated_at;not null"`  // 修改时间
	CreatedUser string    `gorm:"column:created_user;size:60"` // 创建人
	UpdatedUser string    `gorm:"column:updated_user;size:60"` // 更新人
}
