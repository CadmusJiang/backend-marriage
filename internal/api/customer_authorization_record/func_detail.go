package customer_authorization_record

import (
	"net/http"
	"strconv"

	"github.com/xinliangnote/go-gin-api/internal/authz"
	"github.com/xinliangnote/go-gin-api/internal/code"
	"github.com/xinliangnote/go-gin-api/internal/pkg/core"
)

type detailResponse struct {
	Success bool         `json:"success"`
	Data    customerData `json:"data"`
}

// GetCustomerAuthorizationRecordDetail 获取客户授权记录详情
// @Summary 获取客户授权记录详情
// @Description 根据ID获取客户授权记录详情
// @Tags Customer
// @Accept application/x-www-form-urlencoded
// @Produce json
// @Param id path int true "客户授权记录ID"
// @Success 200 {object} detailResponse
// @Failure 400 {object} code.Failure
// @Failure 404 {object} code.Failure
// @Router /api/v1/customer-authorization-records/{id} [get]
func (h *handler) GetCustomerAuthorizationRecordDetail() core.HandlerFunc {
	return func(c core.Context) {
		res := new(detailResponse)

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

		// Mock数据 - 与OpenAPI定义的数据结构保持一致
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
			{
				Key:                 2,
				Name:                "用户2",
				BirthYear:           intPtr(1990),
				Gender:              stringPtr("female"),
				Height:              intPtr(165),
				City:                stringPtr("310000"),
				AuthStore:           stringPtr("浦东门店***"),
				Education:           stringPtr("硕士"),
				Profession:          stringPtr("设计师"),
				Income:              intPtr(80),
				Phone:               stringPtr("139****5678"),
				Wechat:              stringPtr("wx_67890****"),
				DrainageAccount:     stringPtr("drainage_002"),
				DrainageId:          stringPtr("D67890"),
				DrainageChannel:     stringPtr("小红书"),
				Remark:              stringPtr(""),
				AuthorizationStatus: "authorized",
				AssignmentStatus:    "assigned",
				CompletionStatus:    "complete",
				PaymentStatus:       "unpaid",
				PaymentAmount:       0.00,
				RefundAmount:        0.00,
				BelongGroup:         &refObject{Id: "g-b", Name: "归属组B"},
				BelongTeam:          &refObject{Id: "t-c", Name: "归属小队C"},
				BelongAccount:       &refObject{Id: "acc-2", Name: "账户2"},
				AuthorizationPhotos: []string{"https://picsum.photos/300/200?random=3", "https://picsum.photos/300/200?random=4", "https://picsum.photos/300/200?random=5"},
				CreatedAt:           "2024-01-14T10:00:00Z",
				UpdatedAt:           "2024-01-14T10:00:00Z",
			},
			{
				Key:                 3,
				Name:                "用户3",
				BirthYear:           intPtr(1988),
				Gender:              stringPtr("male"),
				Height:              intPtr(180),
				City:                stringPtr("440100"),
				AuthStore:           stringPtr("天河门店***"),
				Education:           stringPtr("大专"),
				Profession:          stringPtr("销售"),
				Income:              intPtr(30),
				Phone:               stringPtr("137****9012"),
				Wechat:              stringPtr("wx_34567****"),
				DrainageAccount:     stringPtr("drainage_003"),
				DrainageId:          stringPtr("D34567"),
				DrainageChannel:     stringPtr("小红书"),
				Remark:              stringPtr("备注信息3"),
				AuthorizationStatus: "unauthorized",
				AssignmentStatus:    "unassigned",
				CompletionStatus:    "incomplete",
				PaymentStatus:       "unpaid",
				PaymentAmount:       0.00,
				RefundAmount:        0.00,
				BelongGroup:         &refObject{Id: "g-c", Name: "归属组C"},
				BelongTeam:          nil,
				BelongAccount:       &refObject{Id: "acc-3", Name: "账户3"},
				AuthorizationPhotos: []string{},
				CreatedAt:           "2024-01-15T10:00:00Z",
				UpdatedAt:           "2024-01-15T10:00:00Z",
			},
			{
				Key:                 4,
				Name:                "用户4",
				BirthYear:           intPtr(1992),
				Gender:              stringPtr("female"),
				Height:              intPtr(160),
				City:                stringPtr("440300"),
				AuthStore:           stringPtr("南山门店***"),
				Education:           stringPtr("本科"),
				Profession:          stringPtr("教师"),
				Income:              intPtr(40),
				Phone:               stringPtr("136****3456"),
				Wechat:              stringPtr("wx_78901****"),
				DrainageAccount:     stringPtr("drainage_004"),
				DrainageId:          stringPtr("D78901"),
				DrainageChannel:     stringPtr("小红书"),
				Remark:              stringPtr(""),
				AuthorizationStatus: "authorized",
				AssignmentStatus:    "unassigned",
				CompletionStatus:    "complete",
				PaymentStatus:       "unpaid",
				PaymentAmount:       0.00,
				RefundAmount:        0.00,
				BelongGroup:         &refObject{Id: "g-d", Name: "归属组D"},
				BelongTeam:          &refObject{Id: "t-g", Name: "归属小队G"},
				BelongAccount:       &refObject{Id: "acc-4", Name: "账户4"},
				AuthorizationPhotos: []string{"https://picsum.photos/300/200?random=6", "https://picsum.photos/300/200?random=7", "https://picsum.photos/300/200?random=8"},
				CreatedAt:           "2024-01-16T10:00:00Z",
				UpdatedAt:           "2024-01-16T10:00:00Z",
			},
			{
				Key:                 5,
				Name:                "用户5",
				BirthYear:           intPtr(1987),
				Gender:              stringPtr("male"),
				Height:              intPtr(178),
				City:                stringPtr("330100"),
				AuthStore:           stringPtr("西湖门店***"),
				Education:           stringPtr("博士"),
				Profession:          stringPtr("医生"),
				Income:              intPtr(120),
				Phone:               stringPtr("135****7890"),
				Wechat:              stringPtr("wx_23456****"),
				DrainageAccount:     stringPtr("drainage_005"),
				DrainageId:          stringPtr("D23456"),
				DrainageChannel:     stringPtr("小红书"),
				Remark:              stringPtr("备注信息5"),
				AuthorizationStatus: "authorized",
				AssignmentStatus:    "assigned",
				CompletionStatus:    "complete",
				PaymentStatus:       "paid",
				PaymentAmount:       35000.00,
				RefundAmount:        2000.00,
				BelongGroup:         &refObject{Id: "g-e", Name: "归属组E"},
				BelongTeam:          &refObject{Id: "t-h", Name: "归属小队H"},
				BelongAccount:       &refObject{Id: "acc-5", Name: "账户5"},
				AuthorizationPhotos: []string{"https://picsum.photos/300/200?random=9", "https://picsum.photos/300/200?random=10", "https://picsum.photos/300/200?random=11"},
				CreatedAt:           "2024-01-17T10:00:00Z",
				UpdatedAt:           "2024-01-17T10:00:00Z",
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

		// 非公司/组管理员需要对 phone、wechat 进行脱敏
		scope, _ := authz.ComputeScope(c, h.db)
		if authz.ShouldMask(scope) {
			if foundCustomer.Phone != nil {
				masked := authz.MaskPhone(*foundCustomer.Phone)
				foundCustomer.Phone = &masked
			}
			if foundCustomer.Wechat != nil {
				maskedW := authz.MaskWechat(*foundCustomer.Wechat)
				foundCustomer.Wechat = &maskedW
			}
		}

		res.Data = *foundCustomer
		res.Success = true

		c.Payload(res)
	}
}
