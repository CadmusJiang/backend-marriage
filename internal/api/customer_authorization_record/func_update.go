package customer_authorization_record

import (
	"net/http"
	"strconv"
	"time"

	"github.com/xinliangnote/go-gin-api/internal/code"
	"github.com/xinliangnote/go-gin-api/internal/pkg/core"
)

type updateRequest struct {
	Name                string   `json:"name"`
	BirthYear           *int     `json:"birthYear"`
	Gender              *string  `json:"gender"`
	Height              *int     `json:"height"`
	City                *string  `json:"city"`
	Education           *string  `json:"education"`
	Profession          *string  `json:"profession"`
	Income              *int     `json:"income"`
	Phone               *string  `json:"phone"`
	Wechat              *string  `json:"wechat"`
	DrainageAccount     *string  `json:"drainageAccount"`
	DrainageId          *string  `json:"drainageId"`
	DrainageChannel     *string  `json:"drainageChannel"`
	Remark              *string  `json:"remark"`
	AuthorizationStatus string   `json:"authorizationStatus"`
	AssignmentStatus    string   `json:"assignmentStatus"`
	CompletionStatus    string   `json:"completionStatus"`
	PaymentStatus       string   `json:"paymentStatus"`
	PaymentAmount       float64  `json:"paymentAmount"`
	RefundAmount        float64  `json:"refundAmount"`
	AuthorizationPhotos []string `json:"authorizationPhotos"`
}

type updateResponse struct {
	Data    customerData `json:"data"`
	Success bool         `json:"success"`
}

// UpdateCustomerAuthorizationRecord 更新客户授权记录
// @Summary 更新客户授权记录
// @Description 根据ID更新客户授权记录
// @Tags Customer
// @Accept json
// @Produce json
// @Param id path int true "客户授权记录ID"
// @Param customer body updateRequest true "客户授权记录更新信息"
// @Success 200 {object} updateResponse
// @Failure 400 {object} code.Failure
// @Failure 404 {object} code.Failure
// @Router /api/v1/customer-authorization-records/{id} [put]
func (h *handler) UpdateCustomerAuthorizationRecord() core.HandlerFunc {
	return func(c core.Context) {
		req := new(updateRequest)
		res := new(updateResponse)

		// 获取路径参数
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.AbortWithError(core.Error(
				http.StatusBadRequest,
				code.ParamBindError,
				"无效的ID参数").WithError(err),
			)
			return
		}

		if err := c.ShouldBindJSON(req); err != nil {
			c.AbortWithError(core.Error(
				http.StatusBadRequest,
				code.ParamBindError,
				code.Text(code.ParamBindError)).WithError(err),
			)
			return
		}

		// Mock数据 - 查找现有客户
		mockCustomers := []customerData{
			{
				Key:                 1,
				Name:                "用户1",
				BirthYear:           intPtr(1985),
				Gender:              stringPtr("male"),
				Height:              intPtr(175),
				City:                stringPtr("110000"),
				AuthStore:           stringPtr("朝阳门店***"),
				Education:           stringPtr("本科"),
				Profession:          stringPtr("工程师"),
				Income:              intPtr(50),
				Phone:               stringPtr("138****1234"),
				Wechat:              stringPtr("wx_12345****"),
				DrainageAccount:     stringPtr("drainage_001"),
				DrainageId:          stringPtr("D12345"),
				DrainageChannel:     stringPtr("小红书"),
				Remark:              stringPtr("备注信息1"),
				AuthorizationStatus: "authorized",
				AssignmentStatus:    "assigned",
				CompletionStatus:    "complete",
				PaymentStatus:       "paid",
				PaymentAmount:       25000.00,
				RefundAmount:        0.00,
				BelongGroup:         &refObject{Id: "g-a", Name: "归属组A"},
				BelongTeam:          &refObject{Id: "t-a", Name: "归属小队A"},
				BelongAccount:       &refObject{Id: "acc-1", Name: "账户1"},
				AuthorizationPhotos: []string{"https://picsum.photos/300/200?random=0", "https://picsum.photos/300/200?random=1", "https://picsum.photos/300/200?random=2"},
				CreatedAt:           "2024-01-13T10:00:00Z",
				UpdatedAt:           "2024-01-13T10:00:00Z",
			},
		}

		// 查找指定ID的客户
		var foundCustomer *customerData
		for _, customer := range mockCustomers {
			if customer.Key == id {
				foundCustomer = &customer
				break
			}
		}

		if foundCustomer == nil {
			c.AbortWithError(core.Error(
				http.StatusNotFound,
				code.ServerError,
				"客户不存在"),
			)
			return
		}

		// 更新客户信息
		now := time.Now()
		updatedCustomer := customerData{
			Key:                 foundCustomer.Key,
			Name:                req.Name,
			BirthYear:           req.BirthYear,
			Gender:              req.Gender,
			Height:              req.Height,
			City:                req.City,
			AuthStore:           foundCustomer.AuthStore,
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
			BelongGroup:         foundCustomer.BelongGroup,
			BelongTeam:          foundCustomer.BelongTeam,
			BelongAccount:       foundCustomer.BelongAccount,
			AuthorizationPhotos: req.AuthorizationPhotos,
			CreatedAt:           foundCustomer.CreatedAt,
			UpdatedAt:           now.Format(time.RFC3339),
		}

		res.Data = updatedCustomer
		res.Success = true

		c.Payload(res)
	}
}
