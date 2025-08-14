package customer_authorization_record_history

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

// CustomerAuthorizationRecordHistory 客户授权记录历史记录表
type CustomerAuthorizationRecordHistory struct {
	Id                            int32       `db:"id"`                               // 主键
	HistoryId                     string      `db:"history_id"`                       // 历史记录ID
	CustomerAuthorizationRecordId string      `db:"customer_authorization_record_id"` // 客户授权记录ID
	OperateType                   string      `db:"operate_type"`                     // 操作类型: created, modified, deleted
	OperateTimestamp              int64       `db:"operate_timestamp"`                // 操作时间戳
	Content                       *JSONString `db:"content"`                          // 操作内容
	OperatorUsername              string      `db:"operator_username"`                // 操作人用户名
	OperatorNickname              string      `db:"operator_nickname"`                // 操作人姓名
	OperatorRoleType              string      `db:"operator_role_type"`               // 操作人角色类型
	IsDeleted                     int8        `db:"is_deleted"`                       // 是否删除 1:是  -1:否
	CreatedUser                   string      `db:"created_user"`                     // 创建人
	UpdatedUser                   string      `db:"updated_user"`                     // 更新人
}

// TableName 表名
func (CustomerAuthorizationRecordHistory) TableName() string {
	return "customer_authorization_record_history"
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
