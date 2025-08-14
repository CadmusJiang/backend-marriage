package cooperation_store

import (
	"net/http"

	"github.com/xinliangnote/go-gin-api/internal/authz"
	"github.com/xinliangnote/go-gin-api/internal/code"
	"github.com/xinliangnote/go-gin-api/internal/pkg/core"
	svc "github.com/xinliangnote/go-gin-api/internal/services/cooperation_store"
)

type createRequest struct {
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
}

type createResponse struct {
	Data    storeData `json:"data"`
	Success bool      `json:"success"`
}

func (h *handler) CreateCooperationStore() core.HandlerFunc {
	// @Summary 创建合作门店
	// @Description 创建合作门店
	// @Tags CooperationStore
	// @Accept json
	// @Produce json
	// @Param request body createRequest true "创建门店请求"
	// @Success 200 {object} createResponse
	// @Failure 400 {object} code.Failure
	// @Router /api/v1/cooperation-stores [post]
	return func(ctx core.Context) {
		req := new(createRequest)
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

		// 设置默认状态（根据前端规范默认 active）
		if req.CooperationStatus == "" {
			req.CooperationStatus = "active"
		}

		// 范围校验（短期：仅允许具备组范围的用户创建；并将响应中的 belongGroupId 约束为用户本组）
		scope, _ := authz.ComputeScope(ctx, h.db)
		if !scope.ScopeAll && len(scope.AllowedGroupIDs) == 0 {
			ctx.AbortWithError(core.Error(
				http.StatusForbidden,
				code.RBACError,
				code.Text(code.RBACError)))
			return
		}

		// 写入通过 service
		item, err := h.svc.Create(ctx, &svc.CreateRequest{
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

		ctx.Payload(createResponse{
			Data:    *item,
			Success: true,
		})
	}
}
