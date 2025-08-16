package customer_authorization_record

import (
	"net/http"
	"strconv"
	"time"

	"github.com/xinliangnote/go-gin-api/internal/code"
	"github.com/xinliangnote/go-gin-api/internal/pkg/core"
	"go.uber.org/zap"
)

type updateRequest struct {
	Name            string  `json:"name"`
	BirthYear       *int    `json:"birthYear"`
	Gender          *string `json:"gender"`
	Height          *int    `json:"height"`
	City            *string `json:"city"`
	Education       *string `json:"education"`
	Profession      *string `json:"profession"`
	Income          *int    `json:"income"`
	Phone           *string `json:"phone"`
	Wechat          *string `json:"wechat"`
	DrainageAccount *string `json:"drainageAccount"`
	DrainageId      *string `json:"drainageId"`
	DrainageChannel *string `json:"drainageChannel"`
	Remark          *string `json:"remark"`
}

type updateResponse struct {
	Data    customerData `json:"data"`
	Success bool         `json:"success"`
}

// UpdateCustomerAuthorizationRecord 更新客资授权记录
// @Summary 更新客资授权记录
// @Description 根据ID更新客资授权记录
// @Tags Customer
// @Accept json
// @Produce json
// @Param id path int true "客资授权记录ID"
// @Param customer body updateRequest true "客资授权记录更新信息"
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

		// 更新客资信息
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
			AuthorizationStatus: foundCustomer.AuthorizationStatus,
			AssignmentStatus:    foundCustomer.AssignmentStatus,
			CompletionStatus:    foundCustomer.CompletionStatus,
			PaymentStatus:       foundCustomer.PaymentStatus,
			PaymentAmount:       foundCustomer.PaymentAmount,
			RefundAmount:        foundCustomer.RefundAmount,
			BelongGroup:         foundCustomer.BelongGroup,
			BelongTeam:          foundCustomer.BelongTeam,
			BelongAccount:       foundCustomer.BelongAccount,
			AuthorizationPhotos: foundCustomer.AuthorizationPhotos,
			CreatedAt:           foundCustomer.CreatedAt,
			UpdatedAt:           now.Format(time.RFC3339),
		}

		res.Data = updatedCustomer
		res.Success = true

		c.Payload(res)
	}
}
