package cooperation_store

import (
	"net/http"

	"github.com/xinliangnote/go-gin-api/internal/code"
	"github.com/xinliangnote/go-gin-api/internal/pkg/core"
)

type detailResponse struct {
	Data    storeData `json:"data"`    // 数据
	Success bool      `json:"success"` // 是否成功
}

func (h *handler) GetCooperationStoreDetail() core.HandlerFunc {
	// @Summary 获取合作门店详情
	// @Description 根据ID获取合作门店详情
	// @Tags CooperationStore
	// @Accept application/x-www-form-urlencoded
	// @Produce json
	// @Param id path string true "门店ID"
	// @Success 200 {object} detailResponse
	// @Failure 400 {object} code.Failure
	// @Failure 404 {object} code.Failure
	// @Router /api/v1/cooperation-stores/{id} [get]
	return func(ctx core.Context) {
		idStr := ctx.Param("id")
		if idStr == "" {
			ctx.AbortWithError(core.Error(
				http.StatusBadRequest,
				code.ParamBindError,
				"id is required"))
			return
		}

		item, err := h.svc.Get(ctx, idStr)
		if err != nil {
			ctx.AbortWithError(core.Error(
				http.StatusNotFound,
				code.ServerError,
				"store not found").WithError(err))
			return
		}
		ctx.Payload(detailResponse{Data: *item, Success: true})
	}
}
