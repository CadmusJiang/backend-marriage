package org_info

import (
	"net/http"

	"github.com/xinliangnote/go-gin-api/internal/code"
	"github.com/xinliangnote/go-gin-api/internal/pkg/core"
	"github.com/xinliangnote/go-gin-api/internal/services/org_info"
	"go.uber.org/zap"
)

type updateRequest struct {
	OrgName  string `json:"orgName"`  // 组织名称
	OrgType  string `json:"orgType"`  // 组织类型: 1-group, 2-team
	OrgPath  string `json:"orgPath"`  // 组织路径
	OrgLevel int32  `json:"orgLevel"` // 组织层级: 1-组, 2-团队
	MaxCnt   int32  `json:"maxCnt"`   // 最大成员数量
}

type updateResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

type deleteRequest struct {
	OrgId string `uri:"orgId" binding:"required"` // 组织ID
}

type deleteResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// UpdateOrgInfo 更新组织信息
// @Summary 更新组织信息
// @Description 更新组织信息
// @Tags API.org_info
// @Accept application/json
// @Produce json
// @Param orgId path string true "组织ID"
// @Param request body updateRequest true "更新组织请求"
// @Success 200 {object} updateResponse
// @Failure 400 {object} code.Failure
// @Router /api/v1/org-infos/{orgId} [put]
func (h *handler) UpdateOrgInfo() core.HandlerFunc {
	return func(c core.Context) {
		req := new(updateRequest)
		res := new(updateResponse)

		if err := c.ShouldBindJSON(req); err != nil {
			c.AbortWithError(core.Error(
				http.StatusBadRequest,
				code.ParamBindError,
				code.Text(code.ParamBindError)).WithError(err),
			)
			return
		}

		// 调用服务层更新组织信息
		updateData := &org_info.UpdateOrgInfoData{
			OrgName:  req.OrgName,
			OrgType:  req.OrgType,
			OrgPath:  req.OrgPath,
			OrgLevel: req.OrgLevel,
			MaxCnt:   req.MaxCnt,
		}

		err := h.orgInfoService.Update(c, "", updateData) // 这里需要根据实际需求修改
		if err != nil {
			h.logger.Error("更新组织信息失败", zap.Error(err))
			c.AbortWithError(core.Error(
				http.StatusInternalServerError,
				code.ServerError,
				code.Text(code.ServerError)).WithError(err),
			)
			return
		}

		res.Success = true
		res.Message = "组织信息更新成功"

		c.Payload(res)
	}
}

// DeleteOrgInfo 删除组织信息
// @Summary 删除组织信息
// @Description 删除组织信息
// @Tags API.org_info
// @Accept application/x-www-form-urlencoded
// @Produce json
// @Param orgId path string true "组织ID"
// @Success 200 {object} deleteResponse
// @Failure 400 {object} code.Failure
// @Router /api/v1/org-infos/{orgId} [delete]
func (h *handler) DeleteOrgInfo() core.HandlerFunc {
	return func(c core.Context) {
		req := new(deleteRequest)
		res := new(deleteResponse)

		if err := c.ShouldBindURI(req); err != nil {
			c.AbortWithError(core.Error(
				http.StatusBadRequest,
				code.ParamBindError,
				code.Text(code.ParamBindError)).WithError(err),
			)
			return
		}

		// 调用服务层删除组织信息
		err := h.orgInfoService.Delete(c, req.OrgId)
		if err != nil {
			h.logger.Error("删除组织信息失败", zap.Error(err))
			c.AbortWithError(core.Error(
				http.StatusInternalServerError,
				code.ServerError,
				code.Text(code.ServerError)).WithError(err),
			)
			return
		}

		res.Success = true
		res.Message = "组织信息删除成功"

		c.Payload(res)
	}
}
