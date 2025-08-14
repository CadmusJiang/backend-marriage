package account

// Account 账户表
//
//go:generate gormgen -structs Account -input .
type Account struct {
	Id       int32  // 主键
	Username string // 用户名
	Nickname string // 姓名
	Password string // 密码
	Phone    string // 手机号
	RoleType string // 角色类型
	Status   string // 状态
	Version  int32  // 乐观锁版本号

	// 所属组信息
	BelongGroupId int32 // 所属组ID
	BelongTeamId  int32 // 所属团队ID

	LastLoginAt int64  // 最后登录时间
	CreatedAt   int64  // 创建时间
	UpdatedAt   int64  // 修改时间
	CreatedUser string // 创建人
	UpdatedUser string // 更新人
}
