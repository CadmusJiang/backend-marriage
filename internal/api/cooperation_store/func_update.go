package cooperation_store

import (
	"net/http"
	"strconv"

	"github.com/xinliangnote/go-gin-api/internal/authz"
	"github.com/xinliangnote/go-gin-api/internal/code"
	"github.com/xinliangnote/go-gin-api/internal/pkg/core"
	svc "github.com/xinliangnote/go-gin-api/internal/services/cooperation_store"
)

type updateRequest struct {
	StoreName             string   `json:"storeName"`             // 门店名称
	CooperationCity       string   `json:"cooperationCity"`       // 合作城市
	CooperationType       []string `json:"cooperationType"`       // 合作类型
	StoreShortName        *string  `json:"storeShortName"`        // 门店简称
	CompanyName           *string  `json:"companyName"`           // 公司名称
	CooperationMethod     []string `json:"cooperationMethod"`     // 合作方式
	CooperationStatus     string   `json:"cooperationStatus"`     // 合作状态
	BusinessLicense       *string  `json:"businessLicense"`       // 营业执照
	StorePhotos           []string `json:"storePhotos"`           // 门店照片
	ActualBusinessAddress *string  `json:"actualBusinessAddress"` // 实际经营地址
	ContractPhotos        []string `json:"contractPhotos"`        // 合同照片
	BelongGroupId         string   `json:"belongGroupId"`         // 归属组ID（mock）
	Remark                string   `json:"remark"`                // 备注（mock）
	Version               *int     `json:"version"`               // 资源版本号（mock忽略）
}

type updateResponse struct {
	Data    storeData `json:"data"`    // 数据
	Success bool      `json:"success"` // 是否成功
}

func (h *handler) UpdateCooperationStore() core.HandlerFunc {
	// @Summary 更新合作门店
	// @Description 更新合作门店
	// @Tags CooperationStore
	// @Accept json
	// @Produce json
	// @Param id path string true "门店ID"
	// @Param request body updateRequest true "更新门店请求"
	// @Success 200 {object} updateResponse
	// @Failure 400 {object} code.Failure
	// @Failure 403 {object} code.Failure
	// @Failure 404 {object} code.Failure
	// @Router /api/v1/cooperation-stores/{id} [put]
	return func(ctx core.Context) {
		idStr := ctx.Param("id")
		if idStr == "" {
			ctx.AbortWithError(core.Error(
				http.StatusBadRequest,
				code.ParamBindError,
				"id is required"))
			return
		}

		if _, err := strconv.Atoi(idStr); err != nil {
			ctx.AbortWithError(core.Error(
				http.StatusBadRequest,
				code.ParamBindError,
				"invalid id format"))
			return
		}

		req := new(updateRequest)
		if err := ctx.ShouldBindJSON(req); err != nil {
			ctx.AbortWithError(core.Error(
				http.StatusBadRequest,
				code.ParamBindError,
				code.Text(code.ParamBindError)).WithError(err),
			)
			return
		}

		// 验证必填字段
		if req.StoreName == "" {
			ctx.AbortWithError(core.Error(
				http.StatusBadRequest,
				code.ParamBindError,
				"storeName is required"))
			return
		}

		if req.CooperationCity == "" {
			ctx.AbortWithError(core.Error(
				http.StatusBadRequest,
				code.ParamBindError,
				"cooperationCity is required"))
			return
		}

		if len(req.CooperationType) == 0 {
			ctx.AbortWithError(core.Error(
				http.StatusBadRequest,
				code.ParamBindError,
				"cooperationType is required"))
			return
		}

		if len(req.CooperationMethod) == 0 {
			ctx.AbortWithError(core.Error(
				http.StatusBadRequest,
				code.ParamBindError,
				"cooperationMethod is required"))
			return
		}

		// 范围校验（短期：仅允许具备组范围的用户更新；无表字段，无法逐条校验归属，只做是否具有组范围的校验）
		scope, _ := authz.ComputeScope(ctx, h.db)
		if !scope.ScopeAll && len(scope.AllowedGroupIDs) == 0 {
			ctx.AbortWithError(core.Error(
				http.StatusForbidden,
				code.RBACError,
				code.Text(code.RBACError)))
			return
		}
		item, err := h.svc.Update(ctx, idStr, &svc.UpdateRequest{
			StoreName:             req.StoreName,
			CooperationCity:       req.CooperationCity,
			CooperationType:       req.CooperationType,
			StoreShortName:        req.StoreShortName,
			CompanyName:           req.CompanyName,
			CooperationMethod:     req.CooperationMethod,
			CooperationStatus:     req.CooperationStatus,
			BusinessLicense:       req.BusinessLicense,
			StorePhotos:           req.StorePhotos,
			ActualBusinessAddress: req.ActualBusinessAddress,
			ContractPhotos:        req.ContractPhotos,
		})
		if err != nil {
			ctx.AbortWithError(core.Error(
				http.StatusInternalServerError,
				code.ServerError,
				code.Text(code.ServerError)).WithError(err),
			)
			return
		}
		ctx.Payload(updateResponse{Data: *item, Success: true})
	}
}
