package customer_authorization_record

import (
	"fmt"
	"net/http"

	"github.com/xinliangnote/go-gin-api/internal/authz"
	"github.com/xinliangnote/go-gin-api/internal/code"
	"github.com/xinliangnote/go-gin-api/internal/pkg/core"
)

type listRequest struct {
	Current             int      `form:"current"`             // 当前页码
	PageSize            int      `form:"pageSize"`            // 每页数量
	Name                string   `form:"name"`                // 姓名模糊搜索
	City                []string `form:"city"`                // 城市ID，数组
	Phone               string   `form:"phone"`               // 手机号（支持带*的脱敏搜索）
	AuthorizationStatus []string `form:"authorizationStatus"` // 授权状态数组（authorized/unauthorized）
	AssignmentStatus    []string `form:"assignmentStatus"`    // 分配状态数组（assigned/unassigned）
	CompletionStatus    []string `form:"completionStatus"`    // 完善状态数组（complete/incomplete）
	PaymentStatus       []string `form:"paymentStatus"`       // 付费状态数组（paid/unpaid）
	BirthYearMin        int      `form:"birthYearMin"`        // 出生年份最小值
	BirthYearMax        int      `form:"birthYearMax"`        // 出生年份最大值
	HeightMin           int      `form:"heightMin"`           // 身高最小值
	HeightMax           int      `form:"heightMax"`           // 身高最大值
	IncomeMin           int      `form:"incomeMin"`           // 收入最小值（万）
	IncomeMax           int      `form:"incomeMax"`           // 收入最大值（万）
	Traffic             []string `form:"traffic"`             // 流量类型（natural/paid）
	BelongGroup         []string `form:"belongGroup"`         // 归属组ID数组
	BelongTeamId        []string `form:"belongTeamId"`        // 归属小队ID数组
	BelongAccountId     []string `form:"belongAccountId"`     // 归属账户ID数组
}

type customerData struct {
	Key                 int        `json:"key"`                 // 记录ID
	Name                string     `json:"name"`                // 客户姓名
	BirthYear           *int       `json:"birthYear"`           // 出生年份
	Gender              *string    `json:"gender"`              // 性别
	Height              *int       `json:"height"`              // 身高(cm)
	City                *string    `json:"city"`                // 城市ID
	AuthStore           *string    `json:"authStore"`           // 授权门店名称（展示用）
	Education           *string    `json:"education"`           // 学历
	Profession          *string    `json:"profession"`          // 职业
	Income              *int       `json:"income"`              // 收入（万）
	Phone               *string    `json:"phone"`               // 手机号（可脱敏）
	Wechat              *string    `json:"wechat"`              // 微信号
	DrainageAccount     *string    `json:"drainageAccount"`     // 引流账户
	DrainageId          *string    `json:"drainageId"`          // 引流ID
	DrainageChannel     *string    `json:"drainageChannel"`     // 引流渠道
	Remark              *string    `json:"remark"`              // 备注
	AuthorizationStatus string     `json:"authorizationStatus"` // 授权状态
	AssignmentStatus    string     `json:"assignmentStatus"`    // 分配状态
	CompletionStatus    string     `json:"completionStatus"`    // 完善状态
	PaymentStatus       string     `json:"paymentStatus"`       // 付费状态
	PaymentAmount       float64    `json:"paymentAmount"`       // 支付金额
	RefundAmount        float64    `json:"refundAmount"`        // 退款金额
	BelongGroup         *refObject `json:"belongGroup"`         // 归属组
	BelongTeam          *refObject `json:"belongTeam"`          // 归属小队
	BelongAccount       *refObject `json:"belongAccount"`       // 归属账户
	AuthorizationPhotos []string   `json:"authorizationPhotos"` // 授权照片
	CreatedAt           string     `json:"createdAt"`           // 创建时间
	UpdatedAt           string     `json:"updatedAt"`           // 修改时间
}

type refObject struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type listMeta struct {
	Total    int `json:"total"`
	PageSize int `json:"pageSize"`
	Current  int `json:"current"`
}

type listResponse struct {
	Success bool           `json:"success"`
	Data    []customerData `json:"data"`
	Meta    listMeta       `json:"meta"`
}

// GetCustomerAuthorizationRecordList 获取客户授权记录列表
// @Summary 获取客户授权记录列表
// @Description 分页获取客户授权记录列表，支持搜索和筛选
// @Tags Customer
// @Accept application/x-www-form-urlencoded
// @Produce json
// @Param current query int false "当前页码" default(1)
// @Param pageSize query int false "每页数量" default(10)
// @Param name query string false "姓名搜索"
// @Param city query string false "城市搜索"
// @Param phone query string false "手机号搜索"
// @Param group query string false "归属组搜索"
// @Param team query string false "归属团队搜索"
// @Param account query string false "归属账户搜索"
// @Param isAuthorized query string false "是否已授权筛选"
// @Param isPaid query string false "是否已买单筛选"
// @Param birthYearMin query int false "出生年份最小值"
// @Param birthYearMax query int false "出生年份最大值"
// @Param incomeMin query int false "收入最小值"
// @Param incomeMax query int false "收入最大值"
// @Success 200 {object} listResponse
// @Failure 400 {object} code.Failure
// @Router /api/v1/customer-authorization-records [get]
func (h *handler) GetCustomerAuthorizationRecordList() core.HandlerFunc {
	return func(c core.Context) {
		req := new(listRequest)
		res := new(listResponse)

		if err := c.ShouldBindForm(req); err != nil {
			c.AbortWithError(core.Error(
				http.StatusBadRequest,
				code.ParamBindError,
				code.Text(code.ParamBindError)).WithError(err),
			)
			return
		}

		// 设置默认值
		if req.Current == 0 {
			req.Current = 1
		}
		if req.PageSize == 0 {
			req.PageSize = 10
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

		// 访问范围过滤（基于 Mock 数据）
		scope, _ := authz.ComputeScope(c, h.db)
		filteredData := filterCustomersWithScope(mockCustomers, req, scope)

		// 脱敏：非公司/组管理员对 phone/wechat 做掩码
		if authz.ShouldMask(scope) {
			for i := range filteredData {
				if filteredData[i].Phone != nil {
					masked := authz.MaskPhone(*filteredData[i].Phone)
					filteredData[i].Phone = &masked
				}
				if filteredData[i].Wechat != nil {
					maskedW := authz.MaskWechat(*filteredData[i].Wechat)
					filteredData[i].Wechat = &maskedW
				}
			}
		}

		// 分页逻辑
		total := len(filteredData)
		start := (req.Current - 1) * req.PageSize
		end := start + req.PageSize

		if start >= total {
			res.Data = []customerData{}
		} else if end > total {
			res.Data = filteredData[start:total]
		} else {
			res.Data = filteredData[start:end]
		}

		res.Success = true
		res.Meta = listMeta{Total: total, PageSize: req.PageSize, Current: req.Current}

		c.Payload(res)
	}
}

// 辅助函数
func intPtr(i int) *int {
	return &i
}

func stringPtr(s string) *string {
	return &s
}

// 过滤客户数据
func filterCustomers(customers []customerData, req *listRequest) []customerData {
	var filtered []customerData

	for _, customer := range customers {
		// 姓名搜索
		if req.Name != "" && !contains(customer.Name, req.Name) {
			continue
		}

		// 城市ID筛选（数组）
		if len(req.City) > 0 {
			if customer.City == nil || !stringInSlice(*customer.City, req.City) {
				continue
			}
		}

		// 手机号搜索
		if req.Phone != "" && (customer.Phone == nil || !contains(*customer.Phone, req.Phone)) {
			continue
		}

		// 归属组ID筛选
		if len(req.BelongGroup) > 0 {
			if customer.BelongGroup == nil || !stringInSlice(customer.BelongGroup.Id, req.BelongGroup) {
				continue
			}
		}

		// 归属小队ID筛选
		if len(req.BelongTeamId) > 0 {
			if customer.BelongTeam == nil || !stringInSlice(customer.BelongTeam.Id, req.BelongTeamId) {
				continue
			}
		}

		// 归属账户ID筛选
		if len(req.BelongAccountId) > 0 {
			if customer.BelongAccount == nil || !stringInSlice(customer.BelongAccount.Id, req.BelongAccountId) {
				continue
			}
		}

		// 出生年份范围搜索
		if req.BirthYearMin > 0 || req.BirthYearMax > 0 {
			if customer.BirthYear == nil {
				continue
			}
			if req.BirthYearMin > 0 && *customer.BirthYear < req.BirthYearMin {
				continue
			}
			if req.BirthYearMax > 0 && *customer.BirthYear > req.BirthYearMax {
				continue
			}
		}

		// 身高范围搜索
		if req.HeightMin > 0 || req.HeightMax > 0 {
			if customer.Height == nil {
				continue
			}
			if req.HeightMin > 0 && *customer.Height < req.HeightMin {
				continue
			}
			if req.HeightMax > 0 && *customer.Height > req.HeightMax {
				continue
			}
		}

		// 收入范围搜索（单位：万）
		if req.IncomeMin > 0 || req.IncomeMax > 0 {
			if customer.Income == nil {
				continue
			}
			if req.IncomeMin > 0 && *customer.Income < req.IncomeMin {
				continue
			}
			if req.IncomeMax > 0 && *customer.Income > req.IncomeMax {
				continue
			}
		}

		// 授权状态筛选
		if len(req.AuthorizationStatus) > 0 && !stringInSlice(customer.AuthorizationStatus, req.AuthorizationStatus) {
			continue
		}

		// 分配状态筛选
		if len(req.AssignmentStatus) > 0 && !stringInSlice(customer.AssignmentStatus, req.AssignmentStatus) {
			continue
		}

		// 完善状态筛选
		if len(req.CompletionStatus) > 0 && !stringInSlice(customer.CompletionStatus, req.CompletionStatus) {
			continue
		}

		// 付费状态筛选
		if len(req.PaymentStatus) > 0 && !stringInSlice(customer.PaymentStatus, req.PaymentStatus) {
			continue
		}

		filtered = append(filtered, customer)
	}

	return filtered
}

// 带范围控制的过滤（Mock数据版）
func filterCustomersWithScope(customers []customerData, req *listRequest, scope authz.AccessScope) []customerData {
	// 先按原有筛选过滤
	base := filterCustomers(customers, req)

	// 再按范围过滤
	if scope.ScopeAll {
		return base
	}

	var ret []customerData
	for _, c := range base {
		allowed := false

		// 员工：只能看自己的
		if len(scope.AllowedAccountIDs) > 0 && c.BelongAccount != nil {
			for _, id := range scope.AllowedAccountIDs {
				if c.BelongAccount.Id == toString(id) {
					allowed = true
					break
				}
			}
		}

		// 小队：看本队
		if !allowed && len(scope.AllowedTeamIDs) > 0 && c.BelongTeam != nil {
			for _, id := range scope.AllowedTeamIDs {
				if c.BelongTeam.Id == toString(id) {
					allowed = true
					break
				}
			}
		}

		// 组：看本组
		if !allowed && len(scope.AllowedGroupIDs) > 0 && c.BelongGroup != nil {
			for _, id := range scope.AllowedGroupIDs {
				if c.BelongGroup.Id == toString(id) {
					allowed = true
					break
				}
			}
		}

		if allowed {
			ret = append(ret, c)
		}
	}
	return ret
}

func toString(i int32) string {
	return fmt.Sprintf("%d", i)
}

// 辅助函数：检查字符串是否包含子字符串（不区分大小写）
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr ||
		(len(s) > len(substr) && (s[:len(substr)] == substr ||
			s[len(s)-len(substr):] == substr ||
			containsSubstring(s, substr))))
}

func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func stringInSlice(s string, arr []string) bool {
	for _, v := range arr {
		if v == s {
			return true
		}
	}
	return false
}
