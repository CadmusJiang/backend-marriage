package account

import (
	"fmt"
	"net/http"

	"github.com/xinliangnote/go-gin-api/internal/code"
	"github.com/xinliangnote/go-gin-api/internal/pkg/core"
	accountService "github.com/xinliangnote/go-gin-api/internal/services/account"
	"go.uber.org/zap"
)

type listRequest struct {
	Current      int    `form:"current"`      // 当前页码
	PageSize     int    `form:"pageSize"`     // 每页数量
	Username     string `form:"username"`     // 用户名搜索
	Name         string `form:"name"`         // 姓名搜索
	RoleType     string `form:"roleType"`     // 角色类型筛选
	Status       string `form:"status"`       // 状态筛选
	Phone        string `form:"phone"`        // 手机号搜索
	BelongGroup  string `form:"belongGroup"`  // 所属组筛选
	BelongTeam   string `form:"belongTeam"`   // 所属团队筛选
	IncludeGroup string `form:"includeGroup"` // 是否包含组信息
	IncludeTeam  string `form:"includeTeam"`  // 是否包含团队信息
}

type accountData struct {
	ID          string `json:"id"`          // 账户ID
	Username    string `json:"username"`    // 用户名
	Name        string `json:"name"`        // 姓名
	Phone       string `json:"phone"`       // 手机号
	RoleType    string `json:"roleType"`    // 角色类型
	Status      string `json:"status"`      // 状态
	CreatedAt   int64  `json:"createdAt"`   // 创建时间戳
	UpdatedAt   int64  `json:"updatedAt"`   // 修改时间戳
	LastLoginAt int64  `json:"lastLoginAt"` // 最后登录时间戳
	BelongGroup *org   `json:"belongGroup"` // 所属组
	BelongTeam  *org   `json:"belongTeam"`  // 所属团队
}

type org struct {
	ID         uint64 `json:"id"`         // 组织ID
	Username   string `json:"username"`   // 组织用户名
	Name       string `json:"name"`       // 组织名称
	CreatedAt  int64  `json:"createdAt"`  // 创建时间戳
	UpdatedAt  int64  `json:"updatedAt"`  // 修改时间戳
	CurrentCnt int32  `json:"currentCnt"` // 当前成员数量
}

type listResponse struct {
	Data     []accountData `json:"data"`
	Total    int64         `json:"total"`
	Success  bool          `json:"success"`
	PageSize int           `json:"pageSize"`
	Current  int           `json:"current"`
}

// GetAccountList 获取账户列表
// @Summary 获取账户列表
// @Description 分页获取账户列表，支持搜索和筛选
// @Tags Account
// @Accept application/x-www-form-urlencoded
// @Produce json
// @Param current query int false "当前页码" default(1)
// @Param pageSize query int false "每页数量" default(10)
// @Param username query string false "用户名搜索"
// @Param name query string false "姓名搜索"
// @Param roleType query string false "角色类型筛选"
// @Param status query string false "状态筛选"
// @Param phone query string false "手机号搜索"
// @Param belongGroup query string false "所属组筛选"
// @Param belongTeam query string false "所属团队筛选"
// @Param includeGroup query string false "是否包含组信息" default(true)
// @Param includeTeam query string false "是否包含团队信息" default(true)
// @Success 200 {object} listResponse
// @Failure 400 {object} code.Failure
// @Router /api/v1/accounts [get]
func (h *handler) GetAccountList() core.HandlerFunc {
	return func(c core.Context) {
		req := new(listRequest)
		res := new(listResponse)

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

		// 调用服务层获取账户列表
		searchData := &accountService.SearchData{
			Username:    req.Username,
			Name:        req.Name,
			RoleType:    req.RoleType,
			Status:      req.Status,
			BelongGroup: req.BelongGroup,
			Current:     req.Current,
			PageSize:    req.PageSize,
		}

		// 获取账户列表
		accountList, err := h.accountService.PageList(c, searchData)
		if err != nil {
			h.logger.Error("查询账户列表失败", zap.Error(err))
			c.AbortWithError(core.Error(
				http.StatusInternalServerError,
				code.ServerError,
				code.Text(code.ServerError)).WithError(err),
			)
			return
		}

		// 获取总数
		total, err := h.accountService.PageListCount(c, searchData)
		if err != nil {
			h.logger.Error("查询账户总数失败", zap.Error(err))
			c.AbortWithError(core.Error(
				http.StatusInternalServerError,
				code.ServerError,
				code.Text(code.ServerError)).WithError(err),
			)
			return
		}

		// 转换为API响应格式
		var accountDataList []accountData
		for _, acc := range accountList {
			// 添加调试信息
			fmt.Printf("Debug: Account ID=%d, Username=%s, BelongGroupId=%d, BelongTeamId=%d\n",
				acc.Id, acc.Username, acc.BelongGroupId, acc.BelongTeamId)

			accountData := accountData{
				ID:        fmt.Sprintf("%d", acc.Id),
				Username:  acc.Username,
				Name:      acc.Name,
				Phone:     acc.Phone,
				RoleType:  acc.RoleType,
				Status:    acc.Status, // 直接赋值，不需要格式化
				CreatedAt: acc.CreatedAt.Unix(),
				UpdatedAt: acc.UpdatedAt.Unix(),
				LastLoginAt: func() int64 {
					if acc.LastLoginAt != nil {
						return acc.LastLoginAt.Unix()
					}
					return 0
				}(),
			}

			// 账户接口不做脱敏

			// 根据includeGroup参数决定是否包含组信息
			if req.IncludeGroup != "false" {
				// 通过组织关系表查询账户的归属组信息
				belongGroup, _ := h.getAccountOrgInfo(int(acc.Id))
				accountData.BelongGroup = belongGroup
			}

			// 根据includeTeam参数决定是否包含团队信息
			if req.IncludeTeam != "false" {
				fmt.Printf("Debug: IncludeTeam is not false, checking team info for %s\n", acc.Username)
				// 根据角色类型决定是否包含团队信息
				switch acc.RoleType {
				case "company_manager", "group_manager":
					// company_manager和group_manager不能有归属团队
					fmt.Printf("Debug: %s cannot have belongTeam for %s\n", acc.RoleType, acc.Username)
					accountData.BelongTeam = nil
				case "team_manager":
					// team_manager必须有归属团队
					if acc.BelongTeamId > 0 {
						fmt.Printf("Debug: BelongTeamId > 0, getting team info for ID %d\n", acc.BelongTeamId)
						teamInfo := h.getTeamInfo(int(acc.BelongTeamId))
						if teamInfo != nil {
							accountData.BelongTeam = teamInfo
							fmt.Printf("Debug: Added team info for %s: %s\n", acc.Username, teamInfo.Name)
						} else {
							fmt.Printf("Debug: Team info is nil for ID %d\n", acc.BelongTeamId)
						}
					} else {
						fmt.Printf("Debug: BelongTeamId is 0 for %s\n", acc.Username)
					}
				case "employee":
					// employee可以有归属团队（可选）
					if acc.BelongTeamId > 0 {
						fmt.Printf("Debug: BelongTeamId > 0, getting team info for ID %d\n", acc.BelongTeamId)
						teamInfo := h.getTeamInfo(int(acc.BelongTeamId))
						if teamInfo != nil {
							accountData.BelongTeam = teamInfo
							fmt.Printf("Debug: Added team info for %s: %s\n", acc.Username, teamInfo.Name)
						} else {
							fmt.Printf("Debug: Team info is nil for ID %d\n", acc.BelongTeamId)
						}
					} else {
						fmt.Printf("Debug: BelongTeamId is 0 for %s\n", acc.Username)
					}
				}
			} else {
				fmt.Printf("Debug: IncludeTeam is false for %s\n", acc.Username)
			}

			accountDataList = append(accountDataList, accountData)
		}

		res.Data = accountDataList
		res.Total = total
		res.Success = true
		res.PageSize = req.PageSize
		res.Current = req.Current

		c.Payload(res)
	}
}
