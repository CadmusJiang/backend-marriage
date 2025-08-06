package account_history

// AccountHistory 账户历史记录表
//
//go:generate gormgen -structs AccountHistory -input .
type AccountHistory struct {
	Id uint64 `gorm:"primaryKey"` // 主键

	// 基本信息
	AccountId        uint64 `gorm:"not null"`  // 账户ID
	OperateType      string `gorm:"size:20"`   // 操作类型: created, modified
	OperateTimestamp uint64 `gorm:"not null"`  // 操作时间戳
	Content          string `gorm:"type:json"` // 操作内容 (JSON格式)
	Operator         string `gorm:"size:60"`   // 操作人
	OperatorRoleType string `gorm:"size:20"`   // 操作人角色

	// 审计字段
	CreatedTimestamp  uint64 `gorm:"not null"` // 创建时间戳
	ModifiedTimestamp uint64 `gorm:"not null"` // 修改时间戳
	CreatedUser       string `gorm:"size:60"`  // 创建人
	UpdatedUser       string `gorm:"size:60"`  // 更新人
}
