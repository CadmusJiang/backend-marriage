package account

import "time"

// Account 账户表
//
//go:generate gormgen -structs Account -input .
type Account struct {
	Id uint32 `gorm:"primaryKey"` // 主键

	// 基本信息
	Username string `gorm:"column:username;size:32;not null"` // 用户名
	Name     string // 姓名
	Password string `gorm:"column:password;size:255;not null"` // 密码
	Phone    string `gorm:"column:phone;size:20;not null"`     // 手机号
	RoleType string `gorm:"column:role_type;size:20;not null"` // 角色类型
	Status   string `gorm:"column:status;not null"`            // 状态 enabled disabled

	// 组织关系
	BelongGroupId uint32 `gorm:"column:belong_group_id;not null"` // 所属组ID
	BelongTeamId  uint32 `gorm:"column:belong_team_id;not null"`  // 所属团队ID

	// 登录信息
	LastLoginAt *time.Time `gorm:"column:last_login_at"` // 最后登录时间

	// 审计字段
	Version     int32     `gorm:"column:version;not null"`
	CreatedAt   time.Time `gorm:"column:created_at;not null"`
	UpdatedAt   time.Time `gorm:"column:updated_at;not null"`
	CreatedUser string    `gorm:"column:created_user;size:60"`
	UpdatedUser string    `gorm:"column:updated_user;size:60"`
}
