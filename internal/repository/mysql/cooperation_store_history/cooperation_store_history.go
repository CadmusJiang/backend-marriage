package cooperation_store_history

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

// CooperationStoreHistory 合作门店历史记录表
type CooperationStoreHistory struct {
	Id               uint64      `db:"id"`                 // 主键
	StoreId          uint32      `db:"store_id"`           // 合作门店ID
	OperateType      string      `db:"operate_type"`       // 操作类型: created, modified, deleted
	OccurredAt       int64       `db:"occurred_at"`        // 操作发生时间戳
	Content          *JSONString `db:"content"`            // 操作内容
	OperatorUsername string      `db:"operator_username"`  // 操作人用户名
	OperatorRoleType string      `db:"operator_role_type"` // 操作人角色类型
	CreatedAt        int64       `db:"created_at"`         // 创建时间戳
	UpdatedAt        int64       `db:"updated_at"`         // 修改时间戳
	CreatedUser      string      `db:"created_user"`       // 创建人
	UpdatedUser      string      `db:"updated_user"`       // 更新人
}

// TableName 表名
func (CooperationStoreHistory) TableName() string { return "cooperation_store_history" }

// JSONString 用于处理JSON字段
type JSONString map[string]interface{}

// Value 实现driver.Valuer接口
func (j JSONString) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}
	return json.Marshal(j)
}

// Scan 实现sql.Scanner接口
func (j *JSONString) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}
	switch v := value.(type) {
	case []byte:
		return json.Unmarshal(v, j)
	case string:
		return json.Unmarshal([]byte(v), j)
	default:
		return errors.New("cannot scan non-string value into JSONString")
	}
}
