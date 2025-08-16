package cooperation_store

import (
	"fmt"
	"net/http"
	"strings"

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
	//
	// 参数格式说明：
	// - storeName: 门店名称 (必填，字符串)
	// - cooperationCity: 合作城市编码 (必填，字符串，如: "110100")
	// - cooperationType: 合作类型 (必填，字符串数组，如: ["直营店", "旗舰店"])
	// - storeShortName: 门店简称 (可选，字符串)
	// - companyName: 公司名称 (可选，字符串)
	// - cooperationMethod: 合作方式 (必填，字符串数组，如: ["线上合作", "线下合作"])
	// - cooperationStatus: 合作状态 (可选，字符串，默认: "active")
	// - businessLicense: 营业执照 (可选，字符串)
	// - storePhotos: 门店照片 (可选，字符串数组)
	// - actualBusinessAddress: 实际经营地址 (可选，字符串)
	// - contractPhotos: 合同照片 (可选，字符串数组)
	// - belongGroupId: 归属组ID (可选，字符串，mock字段)
	// - remark: 备注 (可选，字符串，mock字段)
	return func(ctx core.Context) {
		req := new(createRequest)

		// 添加调试日志
		fmt.Printf("DEBUG: 开始处理创建门店请求\n")

		if err := ctx.ShouldBindJSON(req); err != nil {
			fmt.Printf("DEBUG: JSON参数绑定失败: %v\n", err)
			ctx.AbortWithError(core.Error(
				http.StatusBadRequest,
				code.ParamBindError,
				"JSON参数绑定失败: "+err.Error()))
			return
		}

		// 打印解析后的参数
		fmt.Printf("DEBUG: 参数绑定成功，解析后的参数:\n")
		fmt.Printf("DEBUG: - storeName: %q\n", req.StoreName)
		fmt.Printf("DEBUG: - cooperationCity: %q\n", req.CooperationCity)
		fmt.Printf("DEBUG: - cooperationType: %v\n", req.CooperationType)
		fmt.Printf("DEBUG: - cooperationMethod: %v\n", req.CooperationMethod)
		fmt.Printf("DEBUG: - cooperationStatus: %q\n", req.CooperationStatus)
		fmt.Printf("DEBUG: - businessLicense: %v\n", req.BusinessLicense)
		fmt.Printf("DEBUG: - storePhotos: %v\n", req.StorePhotos)
		fmt.Printf("DEBUG: - actualBusinessAddress: %v\n", req.ActualBusinessAddress)
		fmt.Printf("DEBUG: - contractPhotos: %v\n", req.ContractPhotos)

		// 验证必填字段
		var validationErrors []string

		fmt.Printf("DEBUG: 开始参数验证...\n")

		if req.StoreName == "" {
			fmt.Printf("DEBUG: storeName验证失败: 为空\n")
			validationErrors = append(validationErrors, "storeName is required (门店名称不能为空)")
		} else {
			fmt.Printf("DEBUG: storeName验证通过: %q\n", req.StoreName)
		}

		if req.CooperationCity == "" {
			fmt.Printf("DEBUG: cooperationCity验证失败: 为空\n")
			validationErrors = append(validationErrors, "cooperationCity is required (合作城市不能为空)")
		} else {
			fmt.Printf("DEBUG: cooperationCity验证通过: %q\n", req.CooperationCity)
		}

		if len(req.CooperationType) == 0 {
			fmt.Printf("DEBUG: cooperationType验证失败: 数组为空\n")
			validationErrors = append(validationErrors, "cooperationType is required (合作类型不能为空，应为数组格式)")
		} else {
			fmt.Printf("DEBUG: cooperationType验证通过: %v\n", req.CooperationType)
		}

		if len(req.CooperationMethod) == 0 {
			fmt.Printf("DEBUG: cooperationMethod验证失败: 数组为空\n")
			validationErrors = append(validationErrors, "cooperationMethod is required (合作方式不能为空，应为数组格式)")
		} else {
			fmt.Printf("DEBUG: cooperationMethod验证通过: %v\n", req.CooperationMethod)
		}

		// 验证参数格式
		if req.BusinessLicense != nil && *req.BusinessLicense == "" {
			fmt.Printf("DEBUG: businessLicense验证失败: 空字符串\n")
			validationErrors = append(validationErrors, "businessLicense cannot be empty string (营业执照不能为空字符串)")
		}

		if req.StoreShortName != nil && *req.StoreShortName == "" {
			fmt.Printf("DEBUG: storeShortName验证失败: 空字符串\n")
			validationErrors = append(validationErrors, "storeShortName cannot be empty string (门店简称不能为空字符串)")
		}

		if req.CompanyName != nil && *req.CompanyName == "" {
			fmt.Printf("DEBUG: companyName验证失败: 空字符串\n")
			validationErrors = append(validationErrors, "companyName cannot be empty string (公司名称不能为空字符串)")
		}

		if req.ActualBusinessAddress != nil && *req.ActualBusinessAddress == "" {
			fmt.Printf("DEBUG: actualBusinessAddress验证失败: 空字符串\n")
			validationErrors = append(validationErrors, "actualBusinessAddress cannot be empty string (实际经营地址不能为空字符串)")
		}

		// 如果有验证错误，返回详细的错误信息
		if len(validationErrors) > 0 {
			fmt.Printf("DEBUG: 发现 %d 个验证错误\n", len(validationErrors))
			errorMsg := "参数验证失败:\n" + strings.Join(validationErrors, "\n")
			fmt.Printf("DEBUG: 返回错误信息: %s\n", errorMsg)
			ctx.AbortWithError(core.Error(
				http.StatusBadRequest,
				code.ParamBindError,
				errorMsg))
			return
		}

		fmt.Printf("DEBUG: 所有参数验证通过\n")

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
