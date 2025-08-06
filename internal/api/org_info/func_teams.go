package org_info

import (
	"net/http"
	"strings"

	"github.com/xinliangnote/go-gin-api/internal/code"
	"github.com/xinliangnote/go-gin-api/internal/pkg/core"
	"github.com/xinliangnote/go-gin-api/internal/services/org_info"
	"go.uber.org/zap"
)

type teamsRequest struct {
	Current  int    `form:"current"`  // 当前页码
	PageSize int    `form:"pageSize"` // 每页数量
	OrgName  string `form:"orgName"`  // 组织名称搜索
	OrgLevel int32  `form:"orgLevel"` // 组织层级筛选
}

type teamData struct {
	ID         uint64                 `json:"id"`         // 组织ID (自增ID)
	Name       string                 `json:"name"`       // 组织名称
	Type       int32                  `json:"type"`       // 组织类型: 1-group, 2-team
	Path       string                 `json:"path"`       // 组织路径
	Level      int32                  `json:"level"`      // 组织层级
	CurrentCnt int32                  `json:"currentCnt"` // 当前成员数量
	MaxCnt     int32                  `json:"maxCnt"`     // 最大成员数量
	Status     int32                  `json:"status"`     // 状态: 1-启用, 0-禁用
	CreatedAt  int64                  `json:"createdAt"`  // 创建时间戳
	UpdatedAt  int64                  `json:"updatedAt"`  // 修改时间戳
	Extra      map[string]interface{} `json:"extra"`      // 扩展数据
}

type teamsResponse struct {
	Data     []teamData `json:"data"`
	Total    int64      `json:"total"`
	Success  bool       `json:"success"`
	PageSize int        `json:"pageSize"`
	Current  int        `json:"current"`
}

// GetTeams 获取所有teams信息
// @Summary 获取所有teams信息
// @Description 分页获取所有teams信息，支持搜索和筛选
// @Tags API.org_info
// @Accept application/x-www-form-urlencoded
// @Produce json
// @Param current query int false "当前页码" default(1)
// @Param pageSize query int false "每页数量" default(10)
// @Param orgName query string false "组织名称搜索"
// @Param orgLevel query int false "组织层级筛选"
// @Success 200 {object} teamsResponse
// @Failure 400 {object} code.Failure
// @Router /api/v1/teams [get]
func (h *handler) GetTeams() core.HandlerFunc {
	return func(c core.Context) {
		req := new(teamsRequest)
		res := new(teamsResponse)

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

		// 调用服务层获取teams信息列表
		searchData := &org_info.SearchOrgInfoData{
			OrgName:  req.OrgName,
			OrgType:  "2", // 固定为team类型 (2-team)
			OrgLevel: req.OrgLevel,
			Current:  req.Current,
			PageSize: req.PageSize,
		}

		// 获取teams列表
		teamList, err := h.orgInfoService.PageList(c, searchData)
		if err != nil {
			h.logger.Error("查询teams列表失败", zap.Error(err))
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
			h.logger.Error("查询teams总数失败", zap.Error(err))
			c.AbortWithError(core.Error(
				http.StatusInternalServerError,
				code.ServerError,
				code.Text(code.ServerError)).WithError(err),
			)
			return
		}

		// 转换为API响应格式
		var teamDataList []teamData
		for _, team := range teamList {
			// 为team类型的组织不返回extra信息
			extraData := make(map[string]interface{})

			// 去掉组织名称中的"团队"字
			orgName := team.OrgName
			if strings.HasSuffix(orgName, "团队") {
				orgName = strings.TrimSuffix(orgName, "团队")
			}

			teamData := teamData{
				ID:         team.Id,
				Name:       orgName,
				Type:       team.OrgType,
				Path:       team.OrgPath,
				Level:      team.OrgLevel,
				CurrentCnt: team.CurrentCnt,
				MaxCnt:     team.MaxCnt,
				Status:     team.Status,
				CreatedAt:  team.CreatedAt,
				UpdatedAt:  team.UpdatedAt,
				Extra:      extraData,
			}
			teamDataList = append(teamDataList, teamData)
		}

		res.Data = teamDataList
		res.Total = total
		res.Success = true
		res.PageSize = req.PageSize
		res.Current = req.Current

		c.Payload(res)
	}
}
