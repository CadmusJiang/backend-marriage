package account_org_relation

import "time"

// AccountOrgRelation 账户组织关联表
//
//go:generate gormgen -structs AccountOrgRelation -input .
type AccountOrgRelation struct {
	Id uint64 `gorm:"primaryKey"` // 主键

	// 关联信息
	AccountId    uint32 `gorm:"column:account_id;not null"`    // 账户ID
	OrgId        uint32 `gorm:"column:org_id;not null"`        // 组织ID
	RelationType string `gorm:"column:relation_type;not null"` // 关联类型: belong, manage
	Status       string `gorm:"column:status;not null"`        // 状态: active, inactive

	// 时间戳
	Version   int32     `gorm:"column:version;not null"`
	CreatedAt time.Time `gorm:"column:created_at;not null"` // 创建时间
	UpdatedAt time.Time `gorm:"column:updated_at;not null"` // 修改时间

	// 审计字段
	CreatedUser string `gorm:"column:created_user;size:60"` // 创建人
	UpdatedUser string `gorm:"column:updated_user;size:60"` // 更新人
}
