package org_history

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

// OrgHistory 组织历史记录表
type OrgHistory struct {
	Id               uint64      `db:"id"`                 // 主键
	OrgId            uint32      `db:"org_id"`             // 组织ID
	OrgType          int32       `db:"org_type"`           // 组织类型 1:group 2:team
	OperateType      string      `db:"operate_type"`       // 操作类型: created, updated
	OccurredAt       time.Time   `db:"occurred_at"`        // 操作发生时间
	Content          *JSONString `db:"content"`            // 操作内容
	Operator         string      `db:"operator"`           // 操作人
	OperatorRoleType string      `db:"operator_role_type"` // 操作人角色
	CreatedAt        time.Time   `db:"created_at"`         // 创建时间
	UpdatedAt        time.Time   `db:"updated_at"`         // 修改时间
	CreatedUser      string      `db:"created_user"`       // 创建人
	UpdatedUser      string      `db:"updated_user"`       // 更新人
}

// TableName 表名
func (OrgHistory) TableName() string {
	return "org_history"
}

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
