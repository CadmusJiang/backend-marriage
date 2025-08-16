package customer_authorization_record

import (
	"net/http"
	"time"

	"github.com/xinliangnote/go-gin-api/internal/code"
	"github.com/xinliangnote/go-gin-api/internal/pkg/core"
)

type createRequest struct {
	Name                string   `json:"name"`                // 用户称呼/姓名
	Gender              *string  `json:"gender"`              // male/female
	BirthYear           *int     `json:"birthYear"`           // 出生年份
	Height              *int     `json:"height"`              // 身高
	City                *string  `json:"city"`                // 城市ID
	Education           *string  `json:"education"`           // 学历
	Profession          *string  `json:"profession"`          // 职业
	Income              *int     `json:"income"`              // 收入（万）
	Phone               *string  `json:"phone"`               // 11位手机号
	Wechat              *string  `json:"wechat"`              // 微信
	DrainageAccount     *string  `json:"drainageAccount"`     // 引流账户
	DrainageId          *string  `json:"drainageId"`          // 引流ID
	DrainageChannel     *string  `json:"drainageChannel"`     // 引流渠道
	Remark              *string  `json:"remark"`              // 备注
	AuthorizationStatus string   `json:"authorizationStatus"` // authorized/unauthorized
	AssignmentStatus    string   `json:"assignmentStatus"`    // assigned/unassigned
	CompletionStatus    string   `json:"completionStatus"`    // complete/incomplete
	PaymentStatus       string   `json:"paymentStatus"`       // paid/unpaid
	PaymentAmount       float64  `json:"paymentAmount"`       // 支付金额
	RefundAmount        float64  `json:"refundAmount"`        // 退款金额
	AuthorizedStoreId   *string  `json:"authorizedStoreId"`   // 授权门店ID
	BelongTeamId        *string  `json:"belongTeamId"`        // 归属小队ID
	BelongAccountId     *string  `json:"belongAccountId"`     // 归属账户ID
	AuthorizationPhotos []string `json:"authorizationPhotos"` // 授权照片
}

type createResponse struct {
	Data customerData `json:"data"`
}

// CreateCustomerAuthorizationRecord 创建客户授权记录
// @Summary 创建客户授权记录
// @Description 创建新的客户授权记录
// @Tags Customer
// @Accept json
// @Produce json
// @Param customer body createRequest true "客户授权记录信息"
// @Success 200 {object} createResponse
// @Failure 400 {object} code.Failure
// @Router /api/v1/customer-authorization-records [post]
func (h *handler) CreateCustomerAuthorizationRecord() core.HandlerFunc {
	return func(c core.Context) {
		req := new(createRequest)
		res := new(createResponse)

		if err := c.ShouldBindJSON(req); err != nil {
			c.AbortWithError(core.Error(
				http.StatusBadRequest,
				code.ParamBindError,
				code.Text(code.ParamBindError)).WithError(err),
			)
			return
		}

		// 验证必填字段
		if req.Name == "" {
			c.AbortWithError(core.Error(
				http.StatusBadRequest,
				code.ParamBindError,
				"客户姓名不能为空"),
			)
			return
		}

		// 生成新的客户记录
		now := time.Now()
		newCustomer := customerData{
			Key:                 6,
			Name:                req.Name,
			BirthYear:           req.BirthYear,
			Gender:              req.Gender,
			Height:              req.Height,
			City:                req.City,
			AuthStore:           nil,
			Education:           req.Education,
			Profession:          req.Profession,
			Income:              req.Income,
			Phone:               req.Phone,
			Wechat:              req.Wechat,
			DrainageAccount:     req.DrainageAccount,
			DrainageId:          req.DrainageId,
			DrainageChannel:     req.DrainageChannel,
			Remark:              req.Remark,
			AuthorizationStatus: req.AuthorizationStatus,
			AssignmentStatus:    req.AssignmentStatus,
			CompletionStatus:    req.CompletionStatus,
			PaymentStatus:       req.PaymentStatus,
			PaymentAmount:       req.PaymentAmount,
			RefundAmount:        req.RefundAmount,
			BelongGroup:         nil,
			BelongTeam:          nil,
			BelongAccount:       nil,
			AuthorizationPhotos: req.AuthorizationPhotos,
			CreatedAt:           now.Format(time.RFC3339),
			UpdatedAt:           now.Format(time.RFC3339),
		}

		res.Data = newCustomer

		c.Payload(res)
	}
}
