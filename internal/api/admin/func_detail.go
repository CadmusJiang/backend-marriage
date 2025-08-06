package admin

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/xinliangnote/go-gin-api/configs"
	"github.com/xinliangnote/go-gin-api/internal/code"
	"github.com/xinliangnote/go-gin-api/internal/pkg/core"
	"github.com/xinliangnote/go-gin-api/internal/pkg/password"
	"github.com/xinliangnote/go-gin-api/internal/repository/redis"
	"github.com/xinliangnote/go-gin-api/internal/services/admin"
)

type detailResponse struct {
	Username string                 `json:"username"` // 用户名
	Nickname string                 `json:"nickname"` // 昵称
	Mobile   string                 `json:"mobile"`   // 手机号
	Menu     []admin.ListMyMenuData `json:"menu"`     // 菜单栏
}

// Detail 管理员详情
// @Summary 管理员详情
// @Description 管理员详情
// @Tags API.admin
// @Accept application/x-www-form-urlencoded
// @Produce json
// @Success 200 {object} detailResponse
// @Failure 400 {object} code.Failure
// @Router /api/admin/info [get]
// @Security LoginToken
func (h *handler) Detail() core.HandlerFunc {
	return func(ctx core.Context) {
		res := new(detailResponse)

		searchOneData := new(admin.SearchOneData)
		searchOneData.Id = ctx.SessionUserInfo().UserID
		searchOneData.IsUsed = 1

		info, err := h.adminService.Detail(ctx, searchOneData)
		if err != nil {
			ctx.AbortWithError(core.Error(
				http.StatusBadRequest,
				code.AdminDetailError,
				code.Text(code.AdminDetailError)).WithError(err),
			)
			return
		}

		menuCacheData, err := h.cache.Get(configs.RedisKeyPrefixLoginUser+password.GenerateLoginToken(searchOneData.Id)+":menu", redis.WithTrace(ctx.Trace()))
		if err != nil {
			ctx.AbortWithError(core.Error(
				http.StatusBadRequest,
				code.AdminDetailError,
				code.Text(code.AdminDetailError)).WithError(err),
			)
			return
		}

		var menuData []admin.ListMyMenuData
		err = json.Unmarshal([]byte(menuCacheData), &menuData)
		if err != nil {
			ctx.AbortWithError(core.Error(
				http.StatusBadRequest,
				code.AdminDetailError,
				code.Text(code.AdminDetailError)).WithError(err),
			)
			return
		}

		res.Username = info.Username
		res.Nickname = info.Nickname
		res.Mobile = info.Mobile
		res.Menu = menuData
		ctx.Payload(res)
	}
}

// DetailByUsername 通过用户名或用户ID获取管理员详情
// @Summary 通过用户名或用户ID获取管理员详情
// @Description 通过用户名或用户ID获取管理员详情
// @Tags API.admin
// @Accept application/x-www-form-urlencoded
// @Produce json
// @Param username path string true "用户名或用户ID"
// @Success 200 {object} detailResponse
// @Failure 400 {object} code.Failure
// @Router /api/user/{username} [get]
// @Security LoginToken
func (h *handler) DetailByUsername() core.HandlerFunc {
	return func(ctx core.Context) {
		res := new(detailResponse)

		// 从URL参数中获取用户名或用户ID
		param := ctx.Param("username")
		fmt.Printf("DEBUG DetailByUsername: Raw param from URL = '%s'\n", param)
		fmt.Printf("DEBUG DetailByUsername: Param length = %d\n", len(param))
		
		if param == "" {
			fmt.Printf("DEBUG DetailByUsername: Param is empty, returning error\n")
			ctx.AbortWithError(core.Error(
				http.StatusBadRequest,
				code.ParamBindError,
				code.Text(code.ParamBindError)).WithError(errors.New("用户名或用户ID不能为空")),
			)
			return
		}

		searchOneData := new(admin.SearchOneData)
		searchOneData.IsUsed = 1

		// 判断参数是数字（用户ID）还是字符串（用户名）
		if id, err := strconv.ParseInt(param, 10, 32); err == nil {
			// 参数是数字，按用户ID查询
			fmt.Printf("DEBUG DetailByUsername: Param is numeric, searching by ID = %d\n", id)
			searchOneData.Id = int32(id)
		} else {
			// 参数是字符串，按用户名查询
			fmt.Printf("DEBUG DetailByUsername: Param is string, searching by Username = '%s'\n", param)
			searchOneData.Username = param
		}

		fmt.Printf("DEBUG DetailByUsername: About to query with searchOneData: %+v\n", searchOneData)
		
		info, err := h.adminService.Detail(ctx, searchOneData)
		if err != nil {
			fmt.Printf("DEBUG DetailByUsername: adminService.Detail error: %v\n", err)
			ctx.AbortWithError(core.Error(
				http.StatusBadRequest,
				code.AdminDetailError,
				code.Text(code.AdminDetailError)).WithError(err),
			)
			return
		}

		if info == nil {
			fmt.Printf("DEBUG DetailByUsername: adminService.Detail returned nil info\n")
			ctx.AbortWithError(core.Error(
				http.StatusBadRequest,
				code.AdminDetailError,
				code.Text(code.AdminDetailError)).WithError(errors.New("未找到指定用户")),
			)
			return
		}

		fmt.Printf("DEBUG DetailByUsername: Successfully found user info: %+v\n", info)

		menuCacheData, err := h.cache.Get(configs.RedisKeyPrefixLoginUser+password.GenerateLoginToken(info.Id)+":menu", redis.WithTrace(ctx.Trace()))
		if err != nil {
			ctx.AbortWithError(core.Error(
				http.StatusBadRequest,
				code.AdminDetailError,
				code.Text(code.AdminDetailError)).WithError(err),
			)
			return
		}

		var menuData []admin.ListMyMenuData
		err = json.Unmarshal([]byte(menuCacheData), &menuData)
		if err != nil {
			ctx.AbortWithError(core.Error(
				http.StatusBadRequest,
				code.AdminDetailError,
				code.Text(code.AdminDetailError)).WithError(err),
			)
			return
		}

		res.Username = info.Username
		res.Nickname = info.Nickname
		res.Mobile = info.Mobile
		res.Menu = menuData
		ctx.Payload(res)
	}
}
