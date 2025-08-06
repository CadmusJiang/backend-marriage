package org_info

import (
	"net/http"

	"github.com/xinliangnote/go-gin-api/internal/code"
	"github.com/xinliangnote/go-gin-api/internal/pkg/core"
	"github.com/xinliangnote/go-gin-api/internal/services/org_info"
	"go.uber.org/zap"
)

type listRequest struct {
	Current  int    `form:"current"`  // 当前页码
	PageSize int    `form:"pageSize"` // 每页数量
	OrgName  string `form:"orgName"`  // 组织名称搜索
	OrgType  string `form:"orgType"`  // 组织类型筛选
	OrgLevel int32  `form:"orgLevel"` // 组织层级筛选
}

type orgData struct {
	ID         uint64                 `json:"id"`         // 组织ID
	OrgName    string                 `json:"orgName"`    // 组织名称
	OrgType    int32                  `json:"orgType"`    // 组织类型
	OrgPath    string                 `json:"orgPath"`    // 组织路径
	OrgLevel   int32                  `json:"orgLevel"`   // 组织层级
	CurrentCnt int32                  `json:"currentCnt"` // 当前成员数量
	MaxCnt     int32                  `json:"maxCnt"`     // 最大成员数量
	Status     int32                  `json:"status"`     // 状态
	CreatedAt  int64                  `json:"createdAt"`  // 创建时间戳
	UpdatedAt  int64                  `json:"updatedAt"`  // 修改时间戳
	Extra      map[string]interface{} `json:"extra"`      // 扩展数据
}

type orgInfoListResponse struct {
	Data     []orgData `json:"data"`
	Total    int64     `json:"total"`
	Success  bool      `json:"success"`
	PageSize int       `json:"pageSize"`
	Current  int       `json:"current"`
}

// GetOrgInfoList 获取组织信息列表
// @Summary 获取组织信息列表
// @Description 分页获取组织信息列表，支持搜索和筛选
// @Tags API.org_info
// @Accept application/x-www-form-urlencoded
// @Produce json
// @Param current query int false "当前页码" default(1)
// @Param pageSize query int false "每页数量" default(10)
// @Param orgName query string false "组织名称搜索"
// @Param orgType query string false "组织类型筛选"
// @Param orgLevel query int false "组织层级筛选"
// @Success 200 {object} orgInfoListResponse
// @Failure 400 {object} code.Failure
// @Router /api/v1/org-infos [get]
func (h *handler) GetOrgInfoList() core.HandlerFunc {
	return func(c core.Context) {
		req := new(listRequest)
		res := new(orgInfoListResponse)

		if err := c.ShouldBindForm(req); err != nil {
			c.AbortWithError(core.Error(
				http.StatusBadRequest,
				code.ParamBindError,
				code.Text(code.ParamBindError)).WithError(err),
			)
			return
		}

		// 设置默认值
		if req.Current == 0 {
			req.Current = 1
		}
		if req.PageSize == 0 {
			req.PageSize = 10
		}

		// 调用服务层获取组织信息列表
		searchData := &org_info.SearchOrgInfoData{
			OrgName:  req.OrgName,
			OrgType:  req.OrgType,
			OrgLevel: req.OrgLevel,
			Current:  req.Current,
			PageSize: req.PageSize,
		}

		// 获取组织信息列表
		orgList, err := h.orgInfoService.PageList(c, searchData)
		if err != nil {
			h.logger.Error("查询组织信息列表失败", zap.Error(err))
			c.AbortWithError(core.Error(
				http.StatusInternalServerError,
				code.ServerError,
				code.Text(code.ServerError)).WithError(err),
			)
			return
		}

		// 获取总数
		total, err := h.orgInfoService.PageListCount(c, searchData)
		if err != nil {
			h.logger.Error("查询组织信息总数失败", zap.Error(err))
			c.AbortWithError(core.Error(
				http.StatusInternalServerError,
				code.ServerError,
				code.Text(code.ServerError)).WithError(err),
			)
			return
		}

		// 转换为API响应格式
		var orgDataList []orgData
		for _, org := range orgList {
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
				Extra:      make(map[string]interface{}), // 初始化为空map
			}
			orgDataList = append(orgDataList, orgData)
		}

		res.Data = orgDataList
		res.Total = total
		res.Success = true
		res.PageSize = req.PageSize
		res.Current = req.Current

		c.Payload(res)
	}
}
