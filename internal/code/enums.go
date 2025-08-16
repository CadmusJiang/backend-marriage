package code

// 数据库ENUM约束定义
// 这些常量用于数据库表结构中的ENUM字段约束

// 状态枚举
const (
	StatusEnabled  = "enabled"  // 启用
	StatusDisabled = "disabled" // 禁用
	StatusActive   = "active"   // 活跃
	StatusInactive = "inactive" // 非活跃
)

// 组织类型枚举
const (
	OrgTypeGroup = "group" // 组
	OrgTypeTeam  = "team"  // 团队
)

// 关联类型枚举
const (
	RelationTypeBelong = "belong" // 属于
	RelationTypeManage = "manage" // 管理
)

// 角色类型枚举
const (
	RoleTypeCompanyManager = "company_manager" // 公司管理员
	RoleTypeGroupManager   = "group_manager"   // 组管理员
	RoleTypeTeamManager    = "team_manager"    // 小队管理员
	RoleTypeEmployee       = "employee"        // 员工
)

// 布尔状态枚举
const (
	BoolTrue  = "true"  // 是
	BoolFalse = "false" // 否
)

// 操作类型枚举
const (
	OperateTypeCreated = "created" // 创建
	OperateTypeUpdated = "updated" // 修改
	OperateTypeDeleted = "deleted" // 删除
)

// 发布状态枚举
const (
	PublishStatusUnpublished = "unpublished" // 未发布
	PublishStatusPublished   = "published"   // 已发布
)

// 合作状态枚举
const (
	CooperationStatusActive   = "active"   // 活跃
	CooperationStatusInactive = "inactive" // 非活跃
	CooperationStatusPending  = "pending"  // 待定
)

// 客资状态枚举
const (
	// 授权状态
	AuthorizationStatusAuthorized   = "authorized"   // 已授权
	AuthorizationStatusUnauthorized = "unauthorized" // 未授权

	// 分配状态
	AssignmentStatusAssigned   = "assigned"   // 已分配
	AssignmentStatusUnassigned = "unassigned" // 未分配

	// 完善状态
	CompletionStatusComplete   = "complete"   // 已完善
	CompletionStatusIncomplete = "incomplete" // 未完善

	// 付费状态
	PaymentStatusPaid   = "paid"   // 已付费
	PaymentStatusUnpaid = "unpaid" // 未付费
)
