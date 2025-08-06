package org_info

import (
	"net/http"

	"github.com/xinliangnote/go-gin-api/internal/code"
	"github.com/xinliangnote/go-gin-api/internal/pkg/core"
	"go.uber.org/zap"
)

type childrenRequest struct {
	OrgId string `uri:"orgId" binding:"required"` // 组织ID
}

type childrenResponse struct {
	Data    []orgData `json:"data"`
	Success bool      `json:"success"`
	Total   int       `json:"total"`
	OrgId   string    `json:"orgId"`
}

type parentRequest struct {
	OrgId string `uri:"orgId" binding:"required"` // 组织ID
}

type parentResponse struct {
	Data    orgData `json:"data"`
	Success bool    `json:"success"`
	OrgId   string  `json:"orgId"`
}

// GetOrgInfoChildren 获取子组织
// @Summary 获取子组织
// @Description 根据组织ID获取子组织列表
// @Tags API.org_info
// @Accept application/x-www-form-urlencoded
// @Produce json
// @Param orgId path string true "组织ID"
// @Success 200 {object} childrenResponse
// @Failure 400 {object} code.Failure
// @Router /api/v1/org-infos/{orgId}/children [get]
func (h *handler) GetOrgInfoChildren() core.HandlerFunc {
	return func(c core.Context) {
		req := new(childrenRequest)
		res := new(childrenResponse)

		if err := c.ShouldBindURI(req); err != nil {
			c.AbortWithError(core.Error(
				http.StatusBadRequest,
				code.ParamBindError,
				code.Text(code.ParamBindError)).WithError(err),
			)
			return
		}

		// 调用服务层获取子组织
		children, err := h.orgInfoService.GetChildren(c, req.OrgId)
		if err != nil {
			h.logger.Error("查询子组织失败", zap.Error(err))
			c.AbortWithError(core.Error(
				http.StatusInternalServerError,
				code.ServerError,
				code.Text(code.ServerError)).WithError(err),
			)
			return
		}

		// 转换为API响应格式
		var orgDataList []orgData
		for _, org := range children {
			orgData := orgData{
				ID:         org.Id,
				OrgName:    org.OrgName,
				OrgType:    org.OrgType,
				OrgPath:    org.OrgPath,
				OrgLevel:   org.OrgLevel,
				CurrentCnt: org.CurrentCnt,
				MaxCnt:     org.MaxCnt,
				Status:     org.Status,
				CreatedAt:  org.CreatedAt,
				UpdatedAt:  org.UpdatedAt,
			}
			orgDataList = append(orgDataList, orgData)
		}

		res.Success = true
		res.Data = orgDataList

		c.Payload(res)
	}
}

// GetOrgInfoParent 获取父组织
// @Summary 获取父组织
// @Description 根据组织ID获取父组织信息
// @Tags API.org_info
// @Accept application/x-www-form-urlencoded
// @Produce json
// @Param orgId path string true "组织ID"
// @Success 200 {object} parentResponse
// @Failure 400 {object} code.Failure
// @Router /api/v1/org-infos/{orgId}/parent [get]
func (h *handler) GetOrgInfoParent() core.HandlerFunc {
	return func(c core.Context) {
		req := new(parentRequest)
		res := new(parentResponse)

		if err := c.ShouldBindURI(req); err != nil {
			c.AbortWithError(core.Error(
				http.StatusBadRequest,
				code.ParamBindError,
				code.Text(code.ParamBindError)).WithError(err),
			)
			return
		}

		// 调用服务层获取父组织
		parent, err := h.orgInfoService.GetParent(c, req.OrgId)
		if err != nil {
			h.logger.Error("查询父组织失败", zap.Error(err))
			c.AbortWithError(core.Error(
				http.StatusInternalServerError,
				code.ServerError,
				code.Text(code.ServerError)).WithError(err),
			)
			return
		}

		if parent == nil {
			res.Success = true
			res.Data = orgData{}
			c.Payload(res)
			return
		}

		// 转换为API响应格式
		orgData := orgData{
			ID:         parent.Id,
			OrgName:    parent.OrgName,
			OrgType:    parent.OrgType,
			OrgPath:    parent.OrgPath,
			OrgLevel:   parent.OrgLevel,
			CurrentCnt: parent.CurrentCnt,
			MaxCnt:     parent.MaxCnt,
			Status:     parent.Status,
			CreatedAt:  parent.CreatedAt,
			UpdatedAt:  parent.UpdatedAt,
		}

		res.Success = true
		res.Data = orgData

		c.Payload(res)
	}
}
