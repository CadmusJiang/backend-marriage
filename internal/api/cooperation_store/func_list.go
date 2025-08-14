package cooperation_store

import (
	"net/http"
	"strings"

	"github.com/xinliangnote/go-gin-api/internal/authz"
	"github.com/xinliangnote/go-gin-api/internal/code"
	"github.com/xinliangnote/go-gin-api/internal/pkg/core"
	svc "github.com/xinliangnote/go-gin-api/internal/services/cooperation_store"
)

type listRequest struct {
	Current               int    `form:"current"`               // 当前页码
	PageSize              int    `form:"pageSize"`              // 每页数量
	StoreName             string `form:"storeName"`             // 门店名称
	CooperationCity       string `form:"cooperationCity"`       // 合作城市（逗号分隔）
	CooperationType       string `form:"cooperationType"`       // 合作类型（逗号分隔）
	StoreShortName        string `form:"storeShortName"`        // 门店简称
	CompanyName           string `form:"companyName"`           // 公司名称
	CooperationMethod     string `form:"cooperationMethod"`     // 合作方式（逗号分隔）
	CooperationStatus     string `form:"cooperationStatus"`     // 合作状态（逗号分隔）
	ActualBusinessAddress string `form:"actualBusinessAddress"` // 实际经营地址
	ContractSignStatus    string `form:"contractSignStatus"`    // 合同签署状态（逗号分隔）- mock占位
	BelongGroupId         string `form:"belongGroupId"`         // 归属组ID - mock占位
}

type storeData = svc.StoreItem

type listResponse struct {
	Success bool        `json:"success"` // 是否成功
	Data    []storeData `json:"data"`    // 数据列表
	Meta    struct {
		Total    int `json:"total"`
		PageSize int `json:"pageSize"`
		Current  int `json:"current"`
	} `json:"meta"`
}

func (h *handler) GetCooperationStoreList() core.HandlerFunc {
	// @Summary 获取合作门店列表
	// @Description 分页筛选合作门店
	// @Tags CooperationStore
	// @Accept application/x-www-form-urlencoded
	// @Produce json
	// @Param current query int false "当前页码" default(1)
	// @Param pageSize query int false "每页数量" default(10)
	// @Param storeName query string false "门店名称"
	// @Param cooperationCity query string false "合作城市（逗号分隔）"
	// @Param cooperationType query string false "合作类型（逗号分隔）"
	// @Param storeShortName query string false "门店简称"
	// @Param companyName query string false "公司名称"
	// @Param cooperationMethod query string false "合作方式（逗号分隔）"
	// @Param cooperationStatus query string false "合作状态（逗号分隔）"
	// @Param actualBusinessAddress query string false "实际经营地址"
	// @Success 200 {object} listResponse
	// @Failure 400 {object} code.Failure
	// @Router /api/v1/cooperation-stores [get]
	return func(ctx core.Context) {
		req := new(listRequest)
		if err := ctx.ShouldBindQuery(req); err != nil {
			ctx.AbortWithError(core.Error(
				http.StatusBadRequest,
				code.ParamBindError,
				code.Text(code.ParamBindError)).WithError(err),
			)
			return
		}

		// 设置默认值
		if req.Current <= 0 {
			req.Current = 1
		}
		if req.PageSize <= 0 {
			req.PageSize = 10
		}

		// 组装 service 请求
		statuses := splitCSV(req.CooperationStatus)
		sreq := &svc.ListRequest{
			Current:               req.Current,
			PageSize:              req.PageSize,
			StoreName:             req.StoreName,
			CooperationCity:       req.CooperationCity,
			CooperationType:       splitCSV(req.CooperationType),
			StoreShortName:        req.StoreShortName,
			CompanyName:           req.CompanyName,
			CooperationMethod:     splitCSV(req.CooperationMethod),
			CooperationStatus:     statuses,
			ActualBusinessAddress: req.ActualBusinessAddress,
		}

		items, total, err := h.svc.List(ctx, sreq)
		if err != nil {
			ctx.AbortWithError(core.Error(
				http.StatusInternalServerError,
				code.ServerError,
				code.Text(code.ServerError)).WithError(err),
			)
			return
		}
		resp := make([]storeData, 0, len(items))
		// 范围控制（短期：无表字段，基于当前用户范围做内存过滤，默认本组）
		scope, _ := authz.ComputeScope(ctx, h.db)
		for _, item := range items {
			if scope.ScopeAll {
				resp = append(resp, item)
				continue
			}
			if len(scope.AllowedGroupIDs) > 0 {
				resp = append(resp, item)
			}
		}

		out := listResponse{Success: true, Data: resp}
		out.Meta.Total = int(total)
		out.Meta.PageSize = req.PageSize
		out.Meta.Current = req.Current
		ctx.Payload(out)
	}
}

// 辅助函数
func splitCSV(s string) []string {
	if s == "" {
		return nil
	}
	parts := strings.Split(s, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}
