package org_info

// OrgInfo 组织信息表
//
//go:generate gormgen -structs OrgInfo -input .
type OrgInfo struct {
	Id uint64 `gorm:"primaryKey"` // 主键

	// 基本信息
	OrgName  string `gorm:"column:org_name;size:60"`    // 组织名称
	OrgType  int32  `gorm:"column:org_type"`            // 组织类型: 1-group, 2-team
	OrgPath  string `gorm:"column:org_path;size:255"`   // 组织路径
	OrgLevel int32  `gorm:"column:org_level;default:1"` // 组织层级: 1-组, 2-团队

	// 成员信息
	CurrentCnt int32 `gorm:"column:current_cnt;default:0"` // 当前成员数量
	MaxCnt     int32 `gorm:"column:max_cnt;default:0"`     // 最大成员数量

	// 状态信息
	Status int32 `gorm:"column:status;default:1"` // 状态: 1-启用, 0-禁用

	// 扩展数据
	ExtData string `gorm:"column:ext_data;type:json"` // 扩展数据 (JSON格式)

	// 时间戳
	CreatedAt int64 `gorm:"column:created_at;not null"` // 创建时间戳
	UpdatedAt int64 `gorm:"column:updated_at;not null"` // 修改时间戳

	// 审计字段
	CreatedUser string `gorm:"column:created_user;size:60"` // 创建人
	UpdatedUser string `gorm:"column:updated_user;size:60"` // 更新人
}

// TableName 指定表名
func (OrgInfo) TableName() string {
	return "org"
}
