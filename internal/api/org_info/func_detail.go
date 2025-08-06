package org_info

import (
	"net/http"

	"github.com/xinliangnote/go-gin-api/internal/code"
	"github.com/xinliangnote/go-gin-api/internal/pkg/core"
	"go.uber.org/zap"
)

type detailRequest struct {
	OrgId string `uri:"orgId" binding:"required"` // 组织ID
}

type detailResponse struct {
	Data    orgData `json:"data"`
	Success bool    `json:"success"`
}

// GetOrgInfoDetail 获取组织信息详情
// @Summary 获取组织信息详情
// @Description 根据组织ID获取组织信息详情
// @Tags API.org_info
// @Accept application/x-www-form-urlencoded
// @Produce json
// @Param orgId path string true "组织ID"
// @Success 200 {object} detailResponse
// @Failure 400 {object} code.Failure
// @Failure 404 {object} code.Failure
// @Router /api/v1/org-infos/{orgId} [get]
func (h *handler) GetOrgInfoDetail() core.HandlerFunc {
	return func(c core.Context) {
		req := new(detailRequest)
		res := new(detailResponse)

		if err := c.ShouldBindURI(req); err != nil {
			c.AbortWithError(core.Error(
				http.StatusBadRequest,
				code.ParamBindError,
				code.Text(code.ParamBindError)).WithError(err),
			)
			return
		}

		// 调用服务层获取组织信息详情
		orgInfo, err := h.orgInfoService.Detail(c, req.OrgId)
		if err != nil {
			h.logger.Error("查询组织信息详情失败", zap.Error(err))
			c.AbortWithError(core.Error(
				http.StatusInternalServerError,
				code.ServerError,
				code.Text(code.ServerError)).WithError(err),
			)
			return
		}

		// 转换为API响应格式
		orgData := orgData{
			ID:         orgInfo.Id,
			OrgName:    orgInfo.OrgName,
			OrgType:    orgInfo.OrgType,
			OrgPath:    orgInfo.OrgPath,
			OrgLevel:   orgInfo.OrgLevel,
			CurrentCnt: orgInfo.CurrentCnt,
			MaxCnt:     orgInfo.MaxCnt,
			Status:     orgInfo.Status,
			CreatedAt:  orgInfo.CreatedAt,
			UpdatedAt:  orgInfo.UpdatedAt,
		}

		res.Success = true
		res.Data = orgData

		c.Payload(res)
	}
}
