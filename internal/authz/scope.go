package authz

import (
	"fmt"

	"github.com/xinliangnote/go-gin-api/internal/pkg/core"
	"github.com/xinliangnote/go-gin-api/internal/repository/mysql"
	mysqlAccount "github.com/xinliangnote/go-gin-api/internal/repository/mysql/account"
)

// AccessScope 表示当前请求可访问的数据范围
type AccessScope struct {
	ScopeAll          bool
	AllowedGroupIDs   []int32
	AllowedTeamIDs    []int32
	AllowedAccountIDs []int32
	RoleType          string
	SelfAccountID     int32
}

// ComputeScope 基于当前登录用户，计算其可访问范围
func ComputeScope(ctx core.Context, db mysql.Repo) (AccessScope, error) {
	scope := AccessScope{}

	userID := ctx.SessionUserInfo().UserID
	if userID == 0 {
		return scope, fmt.Errorf("未登录或会话无效")
	}

	// 读取账户信息
	qb := mysqlAccount.NewQueryBuilder()
	qb.WhereId(mysql.EqualPredicate, userID)
	acc, err := qb.QueryOne(db.GetDbR())
	if err != nil {
		return scope, fmt.Errorf("查询当前用户失败: %v", err)
	}
	if acc == nil {
		return scope, fmt.Errorf("当前用户不存在")
	}

	scope.RoleType = acc.RoleType
	scope.SelfAccountID = acc.Id

	switch acc.RoleType {
	case "company_manager":
		scope.ScopeAll = true
	case "group_manager":
		if acc.BelongGroupId > 0 {
			scope.AllowedGroupIDs = []int32{acc.BelongGroupId}
		}
	case "team_manager":
		if acc.BelongTeamId > 0 {
			scope.AllowedTeamIDs = []int32{acc.BelongTeamId}
		}
	case "employee":
		scope.AllowedAccountIDs = []int32{acc.Id}
	default:
		// 未知角色，默认最小权限
		scope.AllowedAccountIDs = []int32{acc.Id}
	}

	return scope, nil
}

// CanAccessAccount 判断是否可访问目标账户
func CanAccessAccount(scope AccessScope, target *mysqlAccount.Account) bool {
	if target == nil {
		return false
	}
	if scope.ScopeAll {
		return true
	}
	// 自己
	for _, id := range scope.AllowedAccountIDs {
		if id == target.Id {
			return true
		}
	}
	// 组
	for _, gid := range scope.AllowedGroupIDs {
		if gid > 0 && gid == target.BelongGroupId {
			return true
		}
	}
	// 队
	for _, tid := range scope.AllowedTeamIDs {
		if tid > 0 && tid == target.BelongTeamId {
			return true
		}
	}
	return false
}
