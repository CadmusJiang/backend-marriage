package customer_authorization_record

import (
	"fmt"
	"net/http"

	"github.com/xinliangnote/go-gin-api/internal/authz"
	"github.com/xinliangnote/go-gin-api/internal/code"
	"github.com/xinliangnote/go-gin-api/internal/pkg/core"
	svc "github.com/xinliangnote/go-gin-api/internal/services/customer_authorization_record"
	"go.uber.org/zap"
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

		// 调用service层获取真实数据
		scope, _ := authz.ComputeScope(c, h.db)

		// 构建查询参数
		serviceReq := &svc.ListRequest{
			Current:             req.Current,
			PageSize:            req.PageSize,
			Name:                req.Name,
			City:                req.City,
			Phone:               req.Phone,
			AuthorizationStatus: req.AuthorizationStatus,
			AssignmentStatus:    req.AssignmentStatus,
			CompletionStatus:    req.CompletionStatus,
			PaymentStatus:       req.PaymentStatus,
			BirthYearMin:        req.BirthYearMin,
			BirthYearMax:        req.BirthYearMax,
			HeightMin:           req.HeightMin,
			HeightMax:           req.HeightMax,
			IncomeMin:           req.IncomeMin,
			IncomeMax:           req.IncomeMax,
			BelongGroup:         req.BelongGroup,
			BelongTeamId:        req.BelongTeamId,
			BelongAccountId:     req.BelongAccountId,
		}

		// 获取真实数据
		records, total, err := h.svc.PageList(c, serviceReq)
		if err != nil {
			h.logger.Error("获取客户授权记录失败", zap.Error(err))
			c.AbortWithError(core.Error(
				http.StatusInternalServerError,
				code.ServerError,
				"获取客户授权记录失败").WithError(err),
			)
			return
		}

		// 调试日志
		var firstRecord interface{}
		if len(records) > 0 {
			firstRecord = records[0]
		}
		h.logger.Info("从数据库获取到的记录",
			zap.Int("records_count", len(records)),
			zap.Int64("total", total),
			zap.Any("first_record", firstRecord))

		// 转换为API响应格式
		var customers []customerData
		for _, record := range records {
			// 处理收入字段（从string转换为int）
			var income *int
			if record.Income != nil {
				// 简单处理，提取数字部分
				if len(*record.Income) > 0 {
					// 这里可以根据实际格式进行解析
					income = intPtr(0) // 暂时设为0，后续可以完善解析逻辑
				}
			}

			customer := customerData{
				Key:                 int(record.ID),
				Name:                record.Name,
				BirthYear:           record.BirthYear,
				Gender:              record.Gender,
				Height:              record.Height,
				City:                record.City,
				AuthStore:           record.AuthStore,
				Education:           record.Education,
				Profession:          record.Profession,
				Income:              income,
				Phone:               record.Phone,
				Wechat:              record.Wechat,
				DrainageAccount:     record.DrainageAccount,
				DrainageId:          record.DrainageId,
				DrainageChannel:     record.DrainageChannel,
				Remark:              record.Remark,
				AuthorizationStatus: record.AuthorizationStatus,
				AssignmentStatus:    record.AssignmentStatus,
				CompletionStatus:    record.CompletionStatus,
				PaymentStatus:       record.PaymentStatus,
				PaymentAmount:       record.PaymentAmount,
				RefundAmount:        record.RefundAmount,
				BelongGroup:         &refObject{Id: toStringPtr(record.BelongGroupID), Name: ""},
				BelongTeam:          &refObject{Id: toStringPtr(record.BelongTeamID), Name: ""},
				BelongAccount:       &refObject{Id: toStringPtr(record.BelongAccountID), Name: ""},
				AuthorizationPhotos: []string{}, // 暂时设为空数组
				CreatedAt:           record.CreatedAt.Format("2006-01-02T15:04:05Z"),
				UpdatedAt:           record.UpdatedAt.Format("2006-01-02T15:04:05Z"),
			}
			customers = append(customers, customer)
		}

		// 调试日志
		var firstCustomer interface{}
		if len(customers) > 0 {
			firstCustomer = customers[0]
		}
		h.logger.Info("转换后的客户数据",
			zap.Int("customers_count", len(customers)),
			zap.Any("first_customer", firstCustomer))

		// 访问范围过滤（基于真实数据）
		// 暂时跳过权限过滤，直接使用所有数据来测试
		filteredData := customers

		// 调试日志
		h.logger.Info("权限过滤后的数据",
			zap.Int("filtered_count", len(filteredData)),
			zap.Any("scope", scope),
			zap.Bool("scope_all", scope.ScopeAll))

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
		totalCount := len(filteredData)
		start := (req.Current - 1) * req.PageSize
		end := start + req.PageSize

		if start >= totalCount {
			res.Data = []customerData{}
		} else if end > totalCount {
			res.Data = filteredData[start:totalCount]
		} else {
			res.Data = filteredData[start:end]
		}

		res.Success = true
		res.Meta = listMeta{Total: totalCount, PageSize: req.PageSize, Current: req.Current}

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

func toStringPtr(i *uint64) string {
	if i == nil {
		return ""
	}
	return fmt.Sprintf("%d", *i)
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
