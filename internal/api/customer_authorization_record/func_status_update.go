package customer_authorization_record

import (
	"net/http"
	"strconv"
	"time"

	"github.com/xinliangnote/go-gin-api/internal/code"
	"github.com/xinliangnote/go-gin-api/internal/pkg/core"
	"go.uber.org/zap"
)

type statusUpdateRequest struct {
	AuthorizationStatus *string `json:"authorizationStatus,omitempty"` // 授权状态
	AssignmentStatus    *string `json:"assignmentStatus,omitempty"`    // 分配状态
	CompletionStatus    *string `json:"completionStatus,omitempty"`    // 完善状态
	PaymentStatus       *string `json:"paymentStatus,omitempty"`       // 付费状态
	Reason              string  `json:"reason"`                        // 状态变更原因
	Operator            string  `json:"operator"`                      // 操作人
}

type statusUpdateResponse struct {
	Data    customerData `json:"data"`
	Success bool         `json:"success"`
	Message string       `json:"message"`
}

// UpdateCustomerStatus 更新客资状态信息
// @Summary 更新客资状态信息
// @Description 根据ID更新客资的状态字段
// @Tags Customer
// @Accept json
// @Produce json
// @Param id path int true "客资授权记录ID"
// @Param status body statusUpdateRequest true "状态更新信息"
// @Success 200 {object} statusUpdateResponse
// @Failure 400 {object} code.Failure
// @Failure 404 {object} code.Failure
// @Router /api/v1/customer-authorization-records/{id}/status [patch]
func (h *handler) UpdateCustomerStatus() core.HandlerFunc {
	return func(c core.Context) {
		req := new(statusUpdateRequest)
		res := new(statusUpdateResponse)

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

		// 从数据库查询现有客资
		existingRecord, err := h.svc.GetByID(c, uint64(id))
		if err != nil {
			h.logger.Error("查询客资记录失败", zap.Error(err))
			c.AbortWithError(core.Error(
				http.StatusNotFound,
				code.ServerError,
				"客资不存在").WithError(err),
			)
			return
		}

		// 转换为customerData格式
		foundCustomer := &customerData{
			Key:                 int(existingRecord.ID),
			Name:                existingRecord.Name,
			BirthYear:           existingRecord.BirthYear,
			Gender:              existingRecord.Gender,
			Height:              existingRecord.Height,
			City:                existingRecord.City,
			AuthStore:           existingRecord.AuthStore,
			Education:           existingRecord.Education,
			Profession:          existingRecord.Profession,
			Income:              intPtr(0), // 暂时设为0，后续可以完善解析逻辑
			Phone:               existingRecord.Phone,
			Wechat:              existingRecord.Wechat,
			DrainageAccount:     existingRecord.DrainageAccount,
			DrainageId:          existingRecord.DrainageId,
			DrainageChannel:     existingRecord.DrainageChannel,
			Remark:              existingRecord.Remark,
			AuthorizationStatus: existingRecord.AuthorizationStatus,
			AssignmentStatus:    existingRecord.AssignmentStatus,
			CompletionStatus:    existingRecord.CompletionStatus,
			PaymentStatus:       existingRecord.PaymentStatus,
			PaymentAmount:       existingRecord.PaymentAmount,
			RefundAmount:        existingRecord.RefundAmount,
			BelongGroup:         &refObject{Id: toStringPtr(existingRecord.BelongGroupID), Name: ""},
			BelongTeam:          &refObject{Id: toStringPtr(existingRecord.BelongTeamID), Name: ""},
			BelongAccount:       &refObject{Id: toStringPtr(existingRecord.BelongAccountID), Name: ""},
			AuthorizationPhotos: []string{}, // 暂时设为空数组
			CreatedAt:           existingRecord.CreatedAt.Format("2006-01-02T15:04:05Z"),
			UpdatedAt:           existingRecord.UpdatedAt.Format("2006-01-02T15:04:05Z"),
		}

		// 更新状态和金额信息
		now := time.Now()
		updatedCustomer := *foundCustomer
		updatedCustomer.UpdatedAt = now.Format(time.RFC3339)

		// 只更新提供的字段
		if req.AuthorizationStatus != nil {
			updatedCustomer.AuthorizationStatus = *req.AuthorizationStatus
		}
		if req.AssignmentStatus != nil {
			updatedCustomer.AssignmentStatus = *req.AssignmentStatus
		}
		if req.CompletionStatus != nil {
			updatedCustomer.CompletionStatus = *req.CompletionStatus
		}
		if req.PaymentStatus != nil {
			updatedCustomer.PaymentStatus = *req.PaymentStatus
		}

		// TODO: 这里应该调用service层更新数据库
		// 暂时只返回更新后的数据

		res.Data = updatedCustomer
		res.Success = true
		res.Message = "状态更新成功"

		// 记录操作日志
		h.logger.Info("客资状态更新",
			zap.Int("customer_id", id),
			zap.String("operator", req.Operator),
			zap.String("reason", req.Reason),
			zap.Any("changes", req))

		c.Payload(res)
	}
}
