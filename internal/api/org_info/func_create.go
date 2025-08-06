package org_info

import (
	"net/http"

	"github.com/xinliangnote/go-gin-api/internal/code"
	"github.com/xinliangnote/go-gin-api/internal/pkg/core"
	"github.com/xinliangnote/go-gin-api/internal/services/org_info"
	"go.uber.org/zap"
)

type createRequest struct {
	OrgName  string `json:"orgName" binding:"required"` // 组织名称
	OrgType  string `json:"orgType" binding:"required"` // 组织类型: 1-group, 2-team
	OrgPath  string `json:"orgPath"`                    // 组织路径
	OrgLevel int32  `json:"orgLevel"`                   // 组织层级: 1-组, 2-团队
	MaxCnt   int32  `json:"maxCnt"`                     // 最大成员数量
}

type createResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	ID      int32  `json:"id"`
}

// CreateOrgInfo 创建组织信息
// @Summary 创建组织信息
// @Description 创建新的组织信息
// @Tags API.org_info
// @Accept application/json
// @Produce json
// @Param request body createRequest true "创建组织请求"
// @Success 200 {object} createResponse
// @Failure 400 {object} code.Failure
// @Router /api/v1/org-infos [post]
func (h *handler) CreateOrgInfo() core.HandlerFunc {
	return func(c core.Context) {
		req := new(createRequest)
		res := new(createResponse)

		if err := c.ShouldBindJSON(req); err != nil {
			c.AbortWithError(core.Error(
				http.StatusBadRequest,
				code.ParamBindError,
				code.Text(code.ParamBindError)).WithError(err),
			)
			return
		}

		// 调用服务层创建组织信息
		orgData := &org_info.CreateOrgInfoData{
			OrgName:  req.OrgName,
			OrgType:  req.OrgType,
			OrgPath:  req.OrgPath,
			OrgLevel: req.OrgLevel,
			MaxCnt:   req.MaxCnt,
		}

		id, err := h.orgInfoService.Create(c, orgData)
		if err != nil {
			h.logger.Error("创建组织信息失败", zap.Error(err))
			c.AbortWithError(core.Error(
				http.StatusInternalServerError,
				code.ServerError,
				code.Text(code.ServerError)).WithError(err),
			)
			return
		}

		res.Success = true
		res.Message = "组织信息创建成功"
		res.ID = id

		c.Payload(res)
	}
}
