package customer_authorization_record

// CreateRequest 服务层创建用的请求结构
type CreateRequest struct {
	Name                string
	Gender              *string
	BirthYear           *int32
	Height              *int32
	City                *string
	Education           *string
	Profession          *string
	Income              *string
	Phone               *string
	Wechat              *string
	DrainageAccount     *string
	DrainageId          *string
	DrainageChannel     *string
	Remark              *string
	AuthorizationStatus string
	AssignmentStatus    string
	CompletionStatus    string
	PaymentStatus       string
	PaymentAmount       float64
	RefundAmount        float64
	AuthorizedStore     *string
	BelongTeam          *string
	BelongAccount       *string
	AuthorizationPhotos []string
}

// UpdateRequest 服务层更新用的请求结构（与创建一致，字段可选）
type UpdateRequest struct {
	Name                string
	Gender              *string
	BirthYear           *int32
	Height              *int32
	City                *string
	Education           *string
	Profession          *string
	Income              *string
	Phone               *string
	Wechat              *string
	DrainageAccount     *string
	DrainageId          *string
	DrainageChannel     *string
	Remark              *string
	AuthorizationStatus string
	AssignmentStatus    string
	CompletionStatus    string
	PaymentStatus       string
	PaymentAmount       float64
	RefundAmount        float64
	AuthorizedStore     *string
	BelongTeam          *string
	BelongAccount       *string
	AuthorizationPhotos []string
}
