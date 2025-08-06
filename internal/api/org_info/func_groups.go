package org_info

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/xinliangnote/go-gin-api/internal/code"
	"github.com/xinliangnote/go-gin-api/internal/pkg/core"
	"github.com/xinliangnote/go-gin-api/internal/services/org_info"
	"go.uber.org/zap"
)

type groupsRequest struct {
	Current      int    `form:"current"`      // 当前页码
	PageSize     int    `form:"pageSize"`     // 每页数量
	OrgName      string `form:"orgName"`      // 组织名称搜索
	OrgLevel     int32  `form:"orgLevel"`     // 组织层级筛选
	IncludeTeams string `form:"includeTeams"` // 是否包含teams信息: true/false
}

type teamInfo struct {
	ID         uint64 `json:"id"`         // 团队ID
	Name       string `json:"name"`       // 团队名称
	Type       int32  `json:"type"`       // 组织类型: 1-group, 2-team
	Path       string `json:"path"`       // 组织路径
	Level      int32  `json:"level"`      // 组织层级
	CurrentCnt int32  `json:"currentCnt"` // 当前成员数量
	MaxCnt     int32  `json:"maxCnt"`     // 最大成员数量
	Status     int32  `json:"status"`     // 状态: 1-启用, 0-禁用
	CreatedAt  int64  `json:"createdAt"`  // 创建时间戳
	UpdatedAt  int64  `json:"updatedAt"`  // 修改时间戳
}

type groupData struct {
	ID         uint64                 `json:"id"`              // 组织ID (自增ID)
	Name       string                 `json:"name"`            // 组织名称
	Type       int32                  `json:"type"`            // 组织类型: 1-group, 2-team
	Path       string                 `json:"path"`            // 组织路径
	Level      int32                  `json:"level"`           // 组织层级
	CurrentCnt int32                  `json:"currentCnt"`      // 当前成员数量
	MaxCnt     int32                  `json:"maxCnt"`          // 最大成员数量
	Status     int32                  `json:"status"`          // 状态: 1-启用, 0-禁用
	CreatedAt  int64                  `json:"createdAt"`       // 创建时间戳
	UpdatedAt  int64                  `json:"updatedAt"`       // 修改时间戳
	Extra      map[string]interface{} `json:"extra"`           // 扩展数据
	Teams      []teamInfo             `json:"teams,omitempty"` // 下属teams信息（可选）
}

type groupsResponse struct {
	Data     []groupData `json:"data"`
	Total    int64       `json:"total"`
	Success  bool        `json:"success"`
	PageSize int         `json:"pageSize"`
	Current  int         `json:"current"`
}

// GetGroups 获取所有groups信息
// @Summary 获取所有groups信息
// @Description 分页获取所有groups信息，支持搜索和筛选，可选择包含teams信息
// @Tags API.org_info
// @Accept application/x-www-form-urlencoded
// @Produce json
// @Param current query int false "当前页码" default(1)
// @Param pageSize query int false "每页数量" default(10)
// @Param orgName query string false "组织名称搜索"
// @Param orgLevel query int false "组织层级筛选"
// @Param includeTeams query string false "是否包含teams信息" default(false)
// @Success 200 {object} groupsResponse
// @Failure 400 {object} code.Failure
// @Router /api/v1/groups [get]
func (h *handler) GetGroups() core.HandlerFunc {
	return func(c core.Context) {
		req := new(groupsRequest)
		res := new(groupsResponse)

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

		// 调用服务层获取groups信息列表
		searchData := &org_info.SearchOrgInfoData{
			OrgName:  req.OrgName,
			OrgType:  "1", // 固定为group类型 (1-group)
			OrgLevel: req.OrgLevel,
			Current:  req.Current,
			PageSize: req.PageSize,
		}

		// 获取groups列表
		groupList, err := h.orgInfoService.PageList(c, searchData)
		if err != nil {
			h.logger.Error("查询groups列表失败", zap.Error(err))
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
			h.logger.Error("查询groups总数失败", zap.Error(err))
			c.AbortWithError(core.Error(
				http.StatusInternalServerError,
				code.ServerError,
				code.Text(code.ServerError)).WithError(err),
			)
			return
		}

		// 转换为API响应格式
		var groupDataList []groupData
		for _, group := range groupList {
			// 为group类型的组织只保留prefix信息
			var extraData map[string]interface{}
			if group.OrgType == 1 { // group类型
				extraData = map[string]interface{}{
					"prefix": "weilan_",
				}
			} else {
				extraData = make(map[string]interface{})
			}

			// 去掉组织名称中的"组"字
			orgName := group.OrgName
			if strings.HasSuffix(orgName, "组") {
				orgName = strings.TrimSuffix(orgName, "组")
			}

			groupData := groupData{
				ID:         group.Id,
				Name:       orgName,
				Type:       group.OrgType,
				Path:       group.OrgPath,
				Level:      group.OrgLevel,
				CurrentCnt: group.CurrentCnt,
				MaxCnt:     group.MaxCnt,
				Status:     group.Status,
				CreatedAt:  group.CreatedAt,
				UpdatedAt:  group.UpdatedAt,
				Extra:      extraData,
			}

			// 如果需要包含teams信息，则获取该group的teams
			if req.IncludeTeams == "true" {
				teams, err := h.getTeamsByGroupId(c, group.Id)
				if err != nil {
					h.logger.Error("获取group的teams失败", zap.Error(err), zap.Uint64("groupId", group.Id))
				} else {
					groupData.Teams = teams
				}
			}

			groupDataList = append(groupDataList, groupData)
		}

		res.Data = groupDataList
		res.Total = total
		res.Success = true
		res.PageSize = req.PageSize
		res.Current = req.Current

		c.Payload(res)
	}
}

// getTeamsByGroupId 根据group ID获取其下属的teams
func (h *handler) getTeamsByGroupId(c core.Context, groupId uint64) ([]teamInfo, error) {
	// 构建查询条件：查找parent_org_id等于groupId且org_type为2的teams
	searchData := &org_info.SearchOrgInfoData{
		OrgType:  "2", // team类型
		Current:  1,
		PageSize: 100, // 获取所有teams
	}

	// 获取所有teams
	teamList, err := h.orgInfoService.PageList(c, searchData)
	if err != nil {
		return nil, err
	}

	var teams []teamInfo
	for _, team := range teamList {
		// 检查team是否属于指定的group
		// 这里需要根据实际的数据库结构来判断父子关系
		// 假设通过path来判断：如果team的path包含group的path，则属于该group
		if strings.HasPrefix(team.OrgPath, fmt.Sprintf("/%d/", groupId)) ||
			strings.Contains(team.OrgPath, fmt.Sprintf("/%d/", groupId)) {

			// 去掉组织名称中的"团队"字
			orgName := team.OrgName
			if strings.HasSuffix(orgName, "团队") {
				orgName = strings.TrimSuffix(orgName, "团队")
			}

			teamInfo := teamInfo{
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
			}
			teams = append(teams, teamInfo)
		}
	}

	return teams, nil
}
