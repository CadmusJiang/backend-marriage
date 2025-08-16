package org

import "time"

// Org 组织表
//
//go:generate gormgen -structs Org -input .
type Org struct {
	Id       uint32 `gorm:"primaryKey"`                // 主键
	OrgType  string `gorm:"column:org_type;not null"`  // 组织类型 group team
	ParentId uint32 `gorm:"column:parent_id;not null"` // 父级ID，根为0
	Path     string `gorm:"column:path;size:255;not null"`
	Username string `gorm:"column:username;size:32;not null"`
	Name     string `gorm:"column:name;size:60;not null"`
	Status   string `gorm:"column:status;not null"` // enabled disabled

	Version     int32     `gorm:"column:version;not null"`
	CreatedAt   time.Time `gorm:"column:created_at;not null"`
	UpdatedAt   time.Time `gorm:"column:updated_at;not null"`
	CreatedUser string    `gorm:"column:created_user;size:60"`
	UpdatedUser string    `gorm:"column:updated_user;size:60"`
}
