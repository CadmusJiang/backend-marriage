package customer_authorization_record

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

// CustomerAuthorizationRecord 客户授权记录表
type CustomerAuthorizationRecord struct {
	Id                  int32       `db:"id"`                   // 主键
	Name                string      `db:"name"`                 // 客户姓名
	BirthYear           *int32      `db:"birth_year"`           // 出生年份
	Gender              *string     `db:"gender"`               // 性别
	Height              *int32      `db:"height"`               // 身高(cm)
	CityCode            *string     `db:"city_code"`            // 城市编码
	AuthStore           *string     `db:"auth_store"`           // 授权门店
	Education           *string     `db:"education"`            // 学历
	Profession          *string     `db:"profession"`           // 职业
	Income              *string     `db:"income"`               // 收入
	Phone               *string     `db:"phone"`                // 手机号
	Wechat              *string     `db:"wechat"`               // 微信号
	DrainageAccount     *string     `db:"drainage_account"`     // 引流账户
	DrainageId          *string     `db:"drainage_id"`          // 引流ID
	DrainageChannel     *string     `db:"drainage_channel"`     // 引流渠道
	Remark              *string     `db:"remark"`               // 备注
	AuthorizationStatus string      `db:"authorization_status"` // 授权状态 authorized unauthorized
	AuthPhotos          *JSONString `db:"auth_photos"`          // 授权照片
	CompletionStatus    string      `db:"completion_status"`    // 完善状态 complete incomplete
	AssignmentStatus    string      `db:"assignment_status"`    // 分配状态 assigned unassigned
	PaymentStatus       string      `db:"payment_status"`       // 付费状态 paid unpaid
	PaymentAmount       float64     `db:"payment_amount"`       // 支付金额
	RefundAmount        float64     `db:"refund_amount"`        // 退款金额
	Group               *string     `db:"group"`                // 归属组
	Team                *string     `db:"team"`                 // 归属团队
	Account             *string     `db:"account"`              // 归属账户
	CreatedAt           time.Time   `db:"created_at"`           // 创建时间
	UpdatedAt           time.Time   `db:"updated_at"`           // 修改时间
	CreatedUser         string      `db:"created_user"`         // 创建人
	UpdatedUser         string      `db:"updated_user"`         // 更新人
}

// TableName 表名
func (CustomerAuthorizationRecord) TableName() string {
	return "customer_authorization_record"
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
