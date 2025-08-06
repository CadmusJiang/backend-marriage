package account_org_relation

// AccountOrgRelation 账户组织关联表
//
//go:generate gormgen -structs AccountOrgRelation -input .
type AccountOrgRelation struct {
	Id uint64 `gorm:"primaryKey"` // 主键

	// 关联信息
	AccountId    string `gorm:"column:account_id;size:32;not null"`    // 账户ID
	OrgId        string `gorm:"column:org_id;size:32;not null"`        // 组织ID
	RelationType string `gorm:"column:relation_type;size:20;not null"` // 关联类型: group, team

	// 时间戳
	CreatedTimestamp  int64 `gorm:"column:created_timestamp;not null"`  // 创建时间戳
	ModifiedTimestamp int64 `gorm:"column:modified_timestamp;not null"` // 修改时间戳

	// 审计字段
	CreatedUser string `gorm:"column:created_user;size:60"` // 创建人
	UpdatedUser string `gorm:"column:updated_user;size:60"` // 更新人
}
