package authz

import (
	"fmt"

	"github.com/xinliangnote/go-gin-api/internal/code"
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

	fmt.Printf("DEBUG: ComputeScope - userID: %d\n", userID)

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

	fmt.Printf("DEBUG: Account info - RoleType: %s, BelongGroupId: %d, BelongTeamId: %d\n",
		acc.RoleType, acc.BelongGroupId, acc.BelongTeamId)

	scope.RoleType = acc.RoleType
	scope.SelfAccountID = int32(acc.Id)

	switch acc.RoleType {
	case code.RoleTypeCompanyManager:
		scope.ScopeAll = true
		fmt.Printf("DEBUG: Set ScopeAll = true for company_manager\n")
	case code.RoleTypeGroupManager:
		if acc.BelongGroupId > 0 {
			scope.AllowedGroupIDs = []int32{int32(acc.BelongGroupId)}
			fmt.Printf("DEBUG: Set AllowedGroupIDs = %v for group_manager\n", scope.AllowedGroupIDs)
		}
	case code.RoleTypeTeamManager:
		if acc.BelongTeamId > 0 {
			scope.AllowedTeamIDs = []int32{int32(acc.BelongTeamId)}
			fmt.Printf("DEBUG: Set AllowedTeamIDs = %v for team_manager\n", scope.AllowedTeamIDs)
		}
	case code.RoleTypeEmployee:
		scope.AllowedAccountIDs = []int32{int32(acc.Id)}
		fmt.Printf("DEBUG: Set AllowedAccountIDs = %v for employee\n", scope.AllowedAccountIDs)
	default:
		// 未知角色，默认最小权限
		scope.AllowedAccountIDs = []int32{int32(acc.Id)}
		fmt.Printf("DEBUG: Unknown role type: %s, set AllowedAccountIDs = %v\n", acc.RoleType, scope.AllowedAccountIDs)
	}

	fmt.Printf("DEBUG: Final scope: %+v\n", scope)
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
		if int32(id) == int32(target.Id) {
			return true
		}
	}
	// 组
	for _, gid := range scope.AllowedGroupIDs {
		if gid > 0 && gid == int32(target.BelongGroupId) {
			return true
		}
	}
	// 队
	for _, tid := range scope.AllowedTeamIDs {
		if tid > 0 && tid == int32(target.BelongTeamId) {
			return true
		}
	}
	return false
}

// FilterScopeByRequest 根据请求参数过滤权限范围
// 如果请求的范围小于用户权限范围，则使用请求的范围
func FilterScopeByRequest(scope AccessScope, requestedGroupID, requestedTeamID uint32) AccessScope {
	filteredScope := scope

	// 如果请求指定了组ID，且用户有权限访问该组
	if requestedGroupID > 0 {
		if scope.ScopeAll {
			// company_manager: 可以访问任何组，但限制为请求的组
			filteredScope.ScopeAll = false
			filteredScope.AllowedGroupIDs = []int32{int32(requestedGroupID)}
			filteredScope.AllowedTeamIDs = nil
			filteredScope.AllowedAccountIDs = nil
		} else if len(scope.AllowedGroupIDs) > 0 {
			// group_manager: 检查是否有权限访问请求的组
			hasPermission := false
			for _, allowedGroupID := range scope.AllowedGroupIDs {
				if allowedGroupID == int32(requestedGroupID) {
					hasPermission = true
					break
				}
			}
			if hasPermission {
				// 限制为请求的组
				filteredScope.AllowedGroupIDs = []int32{int32(requestedGroupID)}
				filteredScope.AllowedTeamIDs = nil
				filteredScope.AllowedAccountIDs = nil
			}
			// 如果没有权限，保持原有范围（会返回空结果）
		}
		// 其他角色没有组权限，保持原有范围
	}

	// 如果请求指定了团队ID，且用户有权限访问该团队
	if requestedTeamID > 0 {
		if scope.ScopeAll {
			// company_manager: 可以访问任何团队，但限制为请求的团队
			filteredScope.ScopeAll = false
			filteredScope.AllowedGroupIDs = nil
			filteredScope.AllowedTeamIDs = []int32{int32(requestedTeamID)}
			filteredScope.AllowedAccountIDs = nil
		} else if len(scope.AllowedTeamIDs) > 0 {
			// team_manager: 检查是否有权限访问请求的团队
			hasPermission := false
			for _, allowedTeamID := range scope.AllowedTeamIDs {
				if allowedTeamID == int32(requestedTeamID) {
					hasPermission = true
					break
				}
			}
			if hasPermission {
				// 限制为请求的团队
				filteredScope.AllowedGroupIDs = nil
				filteredScope.AllowedTeamIDs = []int32{int32(requestedTeamID)}
				filteredScope.AllowedAccountIDs = nil
			}
			// 如果没有权限，保持原有范围（会返回空结果）
		}
		// 其他角色没有团队权限，保持原有范围
	}

	return filteredScope
}

// GetEffectiveScope 获取有效的权限范围（考虑请求参数和用户权限的交集）
func GetEffectiveScope(scope AccessScope, requestedGroupID, requestedTeamID uint32) AccessScope {
	// 如果请求的范围小于用户权限范围，使用请求的范围
	if requestedGroupID > 0 || requestedTeamID > 0 {
		return FilterScopeByRequest(scope, requestedGroupID, requestedTeamID)
	}

	// 否则使用用户的完整权限范围
	return scope
}
