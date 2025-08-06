package code

// 组织类型枚举
const (
	OrgTypeGroup = 1 // 组
	OrgTypeTeam  = 2 // 团队
)

// 组织状态枚举
const (
	OrgStatusEnabled  = 1 // 启用
	OrgStatusDisabled = 2 // 禁用
)

// 关联类型枚举
const (
	RelationTypeBelong = 1 // 属于
	RelationTypeManage = 2 // 管理
)

// 关联状态枚举
const (
	RelationStatusActive   = 1 // 活跃
	RelationStatusInactive = 2 // 非活跃
)

// 角色类型枚举
const (
	RoleTypeCompanyManager = 1 // 公司管理员
	RoleTypeGroupManager   = 2 // 组管理员
	RoleTypeTeamManager    = 3 // 小队管理员
	RoleTypeEmployee       = 4 // 员工
)

// 账户状态枚举
const (
	AccountStatusEnabled  = 1 // 启用
	AccountStatusDisabled = 2 // 禁用
)

// 获取组织类型名称
func GetOrgTypeName(orgType int) string {
	switch orgType {
	case OrgTypeGroup:
		return "group"
	case OrgTypeTeam:
		return "team"
	default:
		return "unknown"
	}
}

// 获取组织状态名称
func GetOrgStatusName(status int) string {
	switch status {
	case OrgStatusEnabled:
		return "enabled"
	case OrgStatusDisabled:
		return "disabled"
	default:
		return "unknown"
	}
}

// 获取关联类型名称
func GetRelationTypeName(relationType int) string {
	switch relationType {
	case RelationTypeBelong:
		return "belong"
	case RelationTypeManage:
		return "manage"
	default:
		return "unknown"
	}
}

// 获取关联状态名称
func GetRelationStatusName(status int) string {
	switch status {
	case RelationStatusActive:
		return "active"
	case RelationStatusInactive:
		return "inactive"
	default:
		return "unknown"
	}
}

// 获取角色类型名称
func GetRoleTypeName(roleType int) string {
	switch roleType {
	case RoleTypeCompanyManager:
		return "company_manager"
	case RoleTypeGroupManager:
		return "group_manager"
	case RoleTypeTeamManager:
		return "team_manager"
	case RoleTypeEmployee:
		return "employee"
	default:
		return "unknown"
	}
}

// 获取账户状态名称
func GetAccountStatusName(status int) string {
	switch status {
	case AccountStatusEnabled:
		return "enabled"
	case AccountStatusDisabled:
		return "disabled"
	default:
		return "unknown"
	}
}
