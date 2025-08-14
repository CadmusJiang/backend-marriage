package cooperation_store

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

// CooperationStore 合作门店表
type CooperationStore struct {
	Id                    int32       `db:"id"`                      // 主键
	StoreName             string      `db:"store_name"`              // 门店名称
	CooperationCity       string      `db:"cooperation_city"`        // 合作城市
	CooperationType       *JSONString `db:"cooperation_type"`        // 合作类型
	StoreShortName        *string     `db:"store_short_name"`        // 门店简称
	CompanyName           *string     `db:"company_name"`            // 公司名称
	CooperationMethod     *JSONString `db:"cooperation_method"`      // 合作方式
	CooperationStatus     string      `db:"cooperation_status"`      // 合作状态
	BusinessLicense       *string     `db:"business_license"`        // 营业执照
	StorePhotos           *JSONString `db:"store_photos"`            // 门店照片
	ActualBusinessAddress *string     `db:"actual_business_address"` // 实际经营地址
	ContractPhotos        *JSONString `db:"contract_photos"`         // 合同照片
	CreatedAt             int64       `db:"created_at"`              // 创建时间
	UpdatedAt             int64       `db:"updated_at"`              // 修改时间
	CreatedUser           string      `db:"created_user"`            // 创建人
	UpdatedUser           string      `db:"updated_user"`            // 更新人
}

// TableName 表名
func (CooperationStore) TableName() string {
	return "cooperation_store"
}

// JSONString 用于处理JSON字段
type JSONString []string

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
