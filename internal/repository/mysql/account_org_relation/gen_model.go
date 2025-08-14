package account_org_relation

// AccountOrgRelation 账户组织关联表
//
//go:generate gormgen -structs AccountOrgRelation -input .
type AccountOrgRelation struct {
	Id uint64 `gorm:"primaryKey"` // 主键

	// 关联信息
	AccountId    uint32 `gorm:"column:account_id;not null"`    // 账户ID
	OrgId        uint32 `gorm:"column:org_id;not null"`        // 组织ID
	RelationType int32  `gorm:"column:relation_type;not null"` // 关联类型: 1 belong, 2 manage
	Status       int32  `gorm:"column:status;not null"`        // 状态: 1 active, 2 inactive

	// 时间戳
	Version   int32 `gorm:"column:version;not null"`
	CreatedAt int64 `gorm:"column:created_at;not null"` // 创建时间
	UpdatedAt int64 `gorm:"column:updated_at;not null"` // 修改时间

	// 审计字段
	CreatedUser string `gorm:"column:created_user;size:60"` // 创建人
	UpdatedUser string `gorm:"column:updated_user;size:60"` // 更新人
}
