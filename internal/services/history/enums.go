package history

// OperateType 操作类型枚举
type OperateType string

const (
	// OperateTypeCreated 创建操作
	OperateTypeCreated OperateType = "created"
	// OperateTypeUpdated 更新操作
	OperateTypeUpdated OperateType = "updated"
)

// String 返回操作类型的字符串表示
func (ot OperateType) String() string {
	return string(ot)
}

// IsValid 检查操作类型是否有效
func (ot OperateType) IsValid() bool {
	switch ot {
	case OperateTypeCreated, OperateTypeUpdated:
		return true
	default:
		return false
	}
}

// EntityType 实体类型枚举
type EntityType string

const (
	// EntityTypeAccount 账户实体
	EntityTypeAccount EntityType = "account"
	// EntityTypeOrg 组织实体
	EntityTypeOrg EntityType = "org"
	// EntityTypeCooperationStore 合作门店实体
	EntityTypeCooperationStore EntityType = "cooperation_store"
	// EntityTypeCustomerAuthorizationRecord 客户授权记录实体
	EntityTypeCustomerAuthorizationRecord EntityType = "customer_authorization_record"
	// EntityTypeOrganization 组织实体（别名）
	EntityTypeOrganization EntityType = "organization"
)

// String 返回实体类型的字符串表示
func (et EntityType) String() string {
	return string(et)
}

// IsValid 检查实体类型是否有效
func (et EntityType) IsValid() bool {
	switch et {
	case EntityTypeAccount, EntityTypeOrg, EntityTypeCooperationStore,
		EntityTypeCustomerAuthorizationRecord, EntityTypeOrganization:
		return true
	default:
		return false
	}
}

// GetHistoryTableName 根据实体类型获取对应的历史表名
func (et EntityType) GetHistoryTableName() string {
	switch et {
	case EntityTypeAccount:
		return "account_history"
	case EntityTypeOrg, EntityTypeOrganization:
		return "org_history"
	case EntityTypeCooperationStore:
		return "cooperation_store_history"
	case EntityTypeCustomerAuthorizationRecord:
		return "customer_authorization_record_history"
	default:
		return ""
	}
}
