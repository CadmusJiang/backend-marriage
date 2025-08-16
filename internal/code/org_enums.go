package code

// 获取组织类型名称
func GetOrgTypeName(orgType string) string {
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
func GetOrgStatusName(status string) string {
	switch status {
	case StatusEnabled:
		return "enabled"
	case StatusDisabled:
		return "disabled"
	default:
		return "unknown"
	}
}

// 获取关联类型名称
func GetRelationTypeName(relationType string) string {
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
func GetRelationStatusName(status string) string {
	switch status {
	case StatusActive:
		return "active"
	case StatusInactive:
		return "inactive"
	default:
		return "unknown"
	}
}

// 获取角色类型名称
func GetRoleTypeName(roleType string) string {
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
func GetAccountStatusName(status string) string {
	switch status {
	case StatusEnabled:
		return "enabled"
	case StatusDisabled:
		return "disabled"
	default:
		return "unknown"
	}
}
