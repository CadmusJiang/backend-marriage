package organization

import (
	"net/http"
	"strconv"

	"github.com/xinliangnote/go-gin-api/internal/code"
	"github.com/xinliangnote/go-gin-api/internal/pkg/core"
	"github.com/xinliangnote/go-gin-api/internal/repository/mysql"
	"github.com/xinliangnote/go-gin-api/internal/repository/redis"
	orgsvc "github.com/xinliangnote/go-gin-api/internal/services/organization"
	"go.uber.org/zap"
)

var _ Handler = (*handler)(nil)

type Handler interface {
	i()

	// groups
	GetGroups() core.HandlerFunc
	CreateGroup() core.HandlerFunc
	GetOrgInfoDetail() core.HandlerFunc
	UpdateGroup() core.HandlerFunc
	GetGroupHistory() core.HandlerFunc

	// teams
	ListTeams() core.HandlerFunc
	CreateTeam() core.HandlerFunc
	GetTeam() core.HandlerFunc
	UpdateTeam() core.HandlerFunc
	GetTeamHistory() core.HandlerFunc

	// team members
	ListTeamMembers() core.HandlerFunc
	AddTeamMember() core.HandlerFunc
	RemoveTeamMember() core.HandlerFunc
	UpdateTeamMemberRole() core.HandlerFunc

	// misc
	ListUnassignedAccounts() core.HandlerFunc
}

type handler struct {
	logger     *zap.Logger
	db         mysql.Repo
	cache      redis.Repo
	orgService orgsvc.Service
}

func New(logger *zap.Logger, db mysql.Repo, cache redis.Repo) Handler {
	return &handler{logger: logger, db: db, cache: cache, orgService: orgsvc.New(db)}
}

func (h *handler) i() {}

// Below are placeholder implementations that return empty results.
// They can be incrementally filled with real logic.

func (h *handler) GetGroups() core.HandlerFunc {
	type req struct {
		Current  int    `form:"current,default=1"`
		PageSize int    `form:"pageSize,default=10"`
		Keyword  string `form:"keyword"`
	}
	type meta struct {
		Total    int `json:"total"`
		PageSize int `json:"pageSize"`
		Current  int `json:"current"`
	}
	type resp struct {
		Success bool        `json:"success"`
		Data    interface{} `json:"data"`
		Meta    meta        `json:"meta"`
	}
	// @Summary 获取组列表
	// @Description 分页获取组织组列表
	// @Tags Organization
	// @Accept application/x-www-form-urlencoded
	// @Produce json
	// @Param current query int false "当前页码" default(1)
	// @Param pageSize query int false "每页数量" default(10)
	// @Param keyword query string false "关键词"
	// @Success 200 {object} resp
	// @Failure 400 {object} code.Failure
	// @Router /api/v1/groups [get]
	return func(c core.Context) {
		r := new(req)
		if err := c.ShouldBindQuery(r); err != nil {
			c.AbortWithError(core.Error(http.StatusBadRequest, code.ParamBindError, code.Text(code.ParamBindError)).WithError(err))
			return
		}
		list, total, err := h.orgService.ListGroups(c, r.Current, r.PageSize, r.Keyword)
		if err != nil {
			c.AbortWithError(core.Error(http.StatusInternalServerError, code.ServerError, code.Text(code.ServerError)).WithError(err))
			return
		}
		c.Payload(resp{Success: true, Data: list, Meta: meta{Total: int(total), PageSize: r.PageSize, Current: r.Current}})
	}
}

func (h *handler) CreateGroup() core.HandlerFunc {
	type body struct {
		Username string `json:"username" binding:"required"`
		Nickname string `json:"nickname" binding:"required"`
	}
	type resp struct {
		Success bool        `json:"success"`
		Data    interface{} `json:"data"`
	}
	// @Summary 创建组
	// @Description 创建组织组
	// @Tags Organization
	// @Accept json
	// @Produce json
	// @Param request body body true "创建组请求"
	// @Success 200 {object} resp
	// @Failure 400 {object} code.Failure
	// @Router /api/v1/groups [post]
	return func(c core.Context) {
		var b body
		if err := c.ShouldBindJSON(&b); err != nil {
			c.AbortWithError(core.Error(http.StatusBadRequest, code.ParamBindError, code.Text(code.ParamBindError)).WithError(err))
			return
		}
		id, err := h.orgService.CreateGroup(c, &orgsvc.CreateGroupPayload{Username: b.Username, Nickname: b.Nickname})
		if err != nil {
			c.AbortWithError(core.Error(http.StatusInternalServerError, code.ServerError, code.Text(code.ServerError)).WithError(err))
			return
		}
		data, err := h.orgService.GetGroup(c, strconv.Itoa(int(id)))
		if err != nil {
			c.Payload(resp{Success: true, Data: map[string]interface{}{"id": strconv.Itoa(int(id))}})
			return
		}
		c.Payload(resp{Success: true, Data: data})
	}
}

func (h *handler) GetOrgInfoDetail() core.HandlerFunc {
	type uri struct {
		OrgId string `uri:"orgId"`
	}
	type resp struct {
		Success bool        `json:"success"`
		Data    interface{} `json:"data"`
	}
	// @Summary 获取组详情
	// @Description 根据组ID获取详情
	// @Tags Organization
	// @Accept application/x-www-form-urlencoded
	// @Produce json
	// @Param orgId path string true "组织ID"
	// @Success 200 {object} resp
	// @Failure 400 {object} code.Failure
	// @Failure 404 {object} code.Failure
	// @Router /api/v1/groups/{orgId} [get]
	return func(c core.Context) {
		var u uri
		if err := c.ShouldBindURI(&u); err != nil {
			c.AbortWithError(core.Error(http.StatusBadRequest, code.ParamBindError, code.Text(code.ParamBindError)).WithError(err))
			return
		}
		info, err := h.orgService.GetGroup(c, u.OrgId)
		if err != nil {
			c.AbortWithError(core.Error(http.StatusNotFound, code.ServerError, code.Text(code.ServerError)).WithError(err))
			return
		}
		c.Payload(resp{Success: true, Data: info})
	}
}

func (h *handler) UpdateGroup() core.HandlerFunc {
	type uri struct {
		OrgId string `uri:"orgId"`
	}
	type body struct {
		Nickname string `json:"nickname"`
		Status   int32  `json:"status"`
		Version  int    `json:"version"`
	}
	type resp struct {
		Success bool `json:"success"`
	}
	// @Summary 更新组
	// @Description 更新组织组
	// @Tags Organization
	// @Accept json
	// @Produce json
	// @Param orgId path string true "组织ID"
	// @Param request body body true "更新组请求"
	// @Success 200 {object} resp
	// @Failure 400 {object} code.Failure
	// @Router /api/v1/groups/{orgId} [put]
	return func(c core.Context) {
		var u uri
		if err := c.ShouldBindURI(&u); err != nil {
			c.Payload(resp{Success: false})
			return
		}
		var b body
		if err := c.ShouldBindJSON(&b); err != nil {
			c.Payload(resp{Success: false})
			return
		}
		if _, err := h.orgService.UpdateGroup(c, u.OrgId, &orgsvc.UpdateGroupPayload{Nickname: b.Nickname, Status: b.Status, Version: b.Version}); err != nil {
			c.Payload(resp{Success: false})
			return
		}
		c.Payload(resp{Success: true})
	}
}

func (h *handler) GetGroupHistory() core.HandlerFunc {
	type uri struct {
		OrgId string `uri:"orgId"`
	}
	type query struct {
		Current  int `form:"current,default=1"`
		PageSize int `form:"pageSize,default=10"`
	}
	type meta struct {
		Total    int `json:"total"`
		PageSize int `json:"pageSize"`
		Current  int `json:"current"`
	}
	type resp struct {
		Success bool        `json:"success"`
		Data    interface{} `json:"data"`
		Meta    meta        `json:"meta"`
	}
	// @Summary 获取组历史
	// @Description 分页获取组的历史记录
	// @Tags Organization
	// @Accept application/x-www-form-urlencoded
	// @Produce json
	// @Param orgId path string true "组织ID"
	// @Param current query int false "当前页码" default(1)
	// @Param pageSize query int false "每页数量" default(10)
	// @Success 200 {object} resp
	// @Failure 400 {object} code.Failure
	// @Router /api/v1/groups/{orgId}/history [get]
	return func(c core.Context) {
		var u uri
		if err := c.ShouldBindURI(&u); err != nil {
			c.Payload(resp{Success: false, Data: []interface{}{}, Meta: meta{Total: 0, PageSize: 10, Current: 1}})
			return
		}
		var q query
		_ = c.ShouldBindQuery(&q)
		items, total, _ := h.orgService.ListGroupHistory(c, u.OrgId, q.Current, q.PageSize)
		c.Payload(resp{Success: true, Data: items, Meta: meta{Total: int(total), PageSize: q.PageSize, Current: q.Current}})
	}
}

func (h *handler) ListTeams() core.HandlerFunc {
	type req struct {
		Current       int    `form:"current,default=1"`
		PageSize      int    `form:"pageSize,default=10"`
		Keyword       string `form:"keyword"`
		BelongGroupId string `form:"belongGroupId"`
	}
	type meta struct {
		Total    int `json:"total"`
		PageSize int `json:"pageSize"`
		Current  int `json:"current"`
	}
	type resp struct {
		Success bool        `json:"success"`
		Data    interface{} `json:"data"`
		Meta    meta        `json:"meta"`
	}
	// @Summary 获取团队列表
	// @Description 分页获取团队列表
	// @Tags Organization
	// @Accept application/x-www-form-urlencoded
	// @Produce json
	// @Param current query int false "当前页码" default(1)
	// @Param pageSize query int false "每页数量" default(10)
	// @Param keyword query string false "关键词"
	// @Param belongGroupId query string false "归属组ID"
	// @Success 200 {object} resp
	// @Failure 400 {object} code.Failure
	// @Router /api/v1/teams [get]
	return func(c core.Context) {
		r := new(req)
		_ = c.ShouldBindQuery(r)
		list, total, err := h.orgService.ListTeams(c, r.Current, r.PageSize, r.BelongGroupId, r.Keyword)
		if err != nil {
			c.AbortWithError(core.Error(http.StatusInternalServerError, code.ServerError, code.Text(code.ServerError)).WithError(err))
			return
		}
		c.Payload(resp{Success: true, Data: list, Meta: meta{Total: int(total), PageSize: r.PageSize, Current: r.Current}})
	}
}

func (h *handler) CreateTeam() core.HandlerFunc {
	type body struct {
		BelongGroupId string `json:"belongGroupId" binding:"required"`
		Username      string `json:"username" binding:"required"`
		Nickname      string `json:"nickname" binding:"required"`
	}
	type resp struct {
		Success bool        `json:"success"`
		Data    interface{} `json:"data"`
	}
	// @Summary 创建团队
	// @Description 创建团队
	// @Tags Organization
	// @Accept json
	// @Produce json
	// @Param request body body true "创建团队请求"
	// @Success 200 {object} resp
	// @Failure 400 {object} code.Failure
	// @Router /api/v1/teams [post]
	return func(c core.Context) {
		var b body
		if err := c.ShouldBindJSON(&b); err != nil {
			c.AbortWithError(core.Error(http.StatusBadRequest, code.ParamBindError, code.Text(code.ParamBindError)).WithError(err))
			return
		}
		id, err := h.orgService.CreateTeam(c, &orgsvc.CreateTeamPayload{BelongGroupId: b.BelongGroupId, Username: b.Username, Nickname: b.Nickname})
		if err != nil {
			c.AbortWithError(core.Error(http.StatusInternalServerError, code.ServerError, code.Text(code.ServerError)).WithError(err))
			return
		}
		data, err := h.orgService.GetTeam(c, strconv.Itoa(int(id)))
		if err != nil {
			c.Payload(resp{Success: true, Data: map[string]interface{}{"id": strconv.Itoa(int(id))}})
			return
		}
		c.Payload(resp{Success: true, Data: data})
	}
}

func (h *handler) GetTeam() core.HandlerFunc {
	type uri struct {
		TeamId string `uri:"teamId"`
	}
	type resp struct {
		Success bool        `json:"success"`
		Data    interface{} `json:"data"`
	}
	// @Summary 获取团队详情
	// @Description 根据团队ID获取详情
	// @Tags Organization
	// @Accept application/x-www-form-urlencoded
	// @Produce json
	// @Param teamId path string true "团队ID"
	// @Success 200 {object} resp
	// @Failure 400 {object} code.Failure
	// @Failure 404 {object} code.Failure
	// @Router /api/v1/teams/{teamId} [get]
	return func(c core.Context) {
		var u uri
		if err := c.ShouldBindURI(&u); err != nil {
			c.AbortWithError(core.Error(http.StatusBadRequest, code.ParamBindError, code.Text(code.ParamBindError)).WithError(err))
			return
		}
		info, err := h.orgService.GetTeam(c, u.TeamId)
		if err != nil {
			c.AbortWithError(core.Error(http.StatusNotFound, code.ServerError, code.Text(code.ServerError)).WithError(err))
			return
		}
		c.Payload(resp{Success: true, Data: info})
	}
}

func (h *handler) UpdateTeam() core.HandlerFunc {
	type uri struct {
		TeamId string `uri:"teamId"`
	}
	type body struct {
		Nickname string `json:"nickname"`
		Status   int32  `json:"status"`
		Version  int    `json:"version"`
	}
	type resp struct {
		Success bool `json:"success"`
	}
	// @Summary 更新团队
	// @Description 更新团队
	// @Tags Organization
	// @Accept json
	// @Produce json
	// @Param teamId path string true "团队ID"
	// @Param request body body true "更新团队请求"
	// @Success 200 {object} resp
	// @Failure 400 {object} code.Failure
	// @Router /api/v1/teams/{teamId} [put]
	return func(c core.Context) {
		var u uri
		if err := c.ShouldBindURI(&u); err != nil {
			c.Payload(resp{Success: false})
			return
		}
		var b body
		if err := c.ShouldBindJSON(&b); err != nil {
			c.Payload(resp{Success: false})
			return
		}
		if _, err := h.orgService.UpdateTeam(c, u.TeamId, &orgsvc.UpdateTeamPayload{Nickname: b.Nickname, Status: b.Status, Version: b.Version}); err != nil {
			c.Payload(resp{Success: false})
			return
		}
		c.Payload(resp{Success: true})
	}
}

func (h *handler) GetTeamHistory() core.HandlerFunc {
	type resp struct {
		Success bool           `json:"success"`
		Data    []interface{}  `json:"data"`
		Meta    map[string]int `json:"meta"`
	}
	// @Summary 获取团队历史
	// @Description 获取团队历史（占位）
	// @Tags Organization
	// @Produce json
	// @Param teamId path string true "团队ID"
	// @Success 200 {object} resp
	// @Router /api/v1/teams/{teamId}/history [get]
	return func(c core.Context) {
		c.Payload(resp{Success: true, Data: []interface{}{}, Meta: map[string]int{"total": 0, "pageSize": 10, "current": 1}})
	}
}

func (h *handler) ListTeamMembers() core.HandlerFunc {
	type resp struct {
		Success bool           `json:"success"`
		Data    []interface{}  `json:"data"`
		Meta    map[string]int `json:"meta"`
	}
	// @Summary 获取团队成员
	// @Description 获取团队成员（占位）
	// @Tags Organization
	// @Produce json
	// @Param teamId path string true "团队ID"
	// @Success 200 {object} resp
	// @Router /api/v1/teams/{teamId}/members [get]
	return func(c core.Context) {
		c.Payload(resp{Success: true, Data: []interface{}{}, Meta: map[string]int{"total": 0, "pageSize": 10, "current": 1}})
	}
}

func (h *handler) AddTeamMember() core.HandlerFunc {
	type resp struct {
		Success bool                   `json:"success"`
		Data    map[string]interface{} `json:"data"`
	}
	// @Summary 添加团队成员
	// @Description 添加团队成员（占位）
	// @Tags Organization
	// @Accept json
	// @Produce json
	// @Param teamId path string true "团队ID"
	// @Success 200 {object} resp
	// @Router /api/v1/teams/{teamId}/members [post]
	return func(c core.Context) {
		c.Payload(resp{Success: true, Data: map[string]interface{}{"id": "u_new", "name": "新成员"}})
	}
}

func (h *handler) RemoveTeamMember() core.HandlerFunc {
	// @Summary 移除团队成员
	// @Description 移除团队成员（占位）
	// @Tags Organization
	// @Param teamId path string true "团队ID"
	// @Param memberId path string true "成员ID"
	// @Success 200
	// @Router /api/v1/teams/{teamId}/members/{memberId} [delete]
	return func(c core.Context) { c.SetHeader("Content-Length", "0"); c.Payload(nil) }
}

func (h *handler) UpdateTeamMemberRole() core.HandlerFunc {
	type resp struct {
		Success bool                   `json:"success"`
		Data    map[string]interface{} `json:"data"`
	}
	// @Summary 更新团队成员角色
	// @Description 更新团队成员角色（占位）
	// @Tags Organization
	// @Accept json
	// @Produce json
	// @Param teamId path string true "团队ID"
	// @Param memberId path string true "成员ID"
	// @Success 200 {object} resp
	// @Router /api/v1/teams/{teamId}/members/{memberId}/role [put]
	return func(c core.Context) {
		c.Payload(resp{Success: true, Data: map[string]interface{}{"id": "u_1", "roleType": "team_manager"}})
	}
}

func (h *handler) ListUnassignedAccounts() core.HandlerFunc {
	type resp struct {
		Success bool           `json:"success"`
		Data    []interface{}  `json:"data"`
		Meta    map[string]int `json:"meta"`
	}
	// @Summary 获取未分配账户
	// @Description 获取未分配账户（占位）
	// @Tags Organization
	// @Produce json
	// @Success 200 {object} resp
	// @Router /api/v1/unassigned-account [get]
	return func(c core.Context) {
		c.Payload(resp{Success: true, Data: []interface{}{}, Meta: map[string]int{"total": 0, "pageSize": 10, "current": 1}})
	}
}
