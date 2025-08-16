package account

import (
	"fmt"

	"github.com/xinliangnote/go-gin-api/internal/authz"
	"github.com/xinliangnote/go-gin-api/internal/pkg/core"
	"github.com/xinliangnote/go-gin-api/internal/repository/mysql"
	"github.com/xinliangnote/go-gin-api/internal/repository/mysql/account"
)

// PageList 分页获取账户列表
func (s *service) PageList(ctx core.Context, searchData *SearchData) (listData []*account.Account, err error) {
	accountQueryBuilder := account.NewQueryBuilder()

	// 添加搜索条件
	if searchData.Username != "" {
		accountQueryBuilder.WhereUsername(mysql.LikePredicate, "%"+searchData.Username+"%")
	}
	if searchData.Name != "" {
		accountQueryBuilder.WhereName(mysql.LikePredicate, "%"+searchData.Name+"%")
	}
	if searchData.RoleType != "" {
		accountQueryBuilder.WhereRoleType(mysql.EqualPredicate, searchData.RoleType)
	}
	if searchData.Status != "" {
		accountQueryBuilder.WhereStatus(mysql.EqualPredicate, searchData.Status)
	}
	// 暂时移除组织搜索，因为新模型中没有相关字段
	// if searchData.BelongGroup != "" {
	// 	accountQueryBuilder.WhereBelongGroupName(mysql.LikePredicate, "%"+searchData.BelongGroup+"%")
	// }

	// 访问范围过滤
	scope, scopeErr := authz.ComputeScope(ctx, s.db)
	if scopeErr == nil && !scope.ScopeAll {
		if len(scope.AllowedAccountIDs) > 0 {
			accountQueryBuilder.WhereIdIn(scope.AllowedAccountIDs)
		} else if len(scope.AllowedTeamIDs) > 0 {
			accountQueryBuilder.WhereBelongTeamIdIn(scope.AllowedTeamIDs)
		} else if len(scope.AllowedGroupIDs) > 0 {
			accountQueryBuilder.WhereBelongGroupIdIn(scope.AllowedGroupIDs)
		}
	}

	// 设置分页
	offset := (searchData.Current - 1) * searchData.PageSize
	accountQueryBuilder.Limit(searchData.PageSize).Offset(offset)

	// 按创建时间倒序排列
	accountQueryBuilder.OrderByCreatedAt(false)

	// 查询数据
	listData, err = accountQueryBuilder.QueryAll(s.db.GetDbR())
	if err != nil {
		return nil, fmt.Errorf("查询账户列表失败: %v", err)
	}

	return listData, nil
}

// PageListCount 获取账户列表总数
func (s *service) PageListCount(ctx core.Context, searchData *SearchData) (total int64, err error) {
	accountQueryBuilder := account.NewQueryBuilder()

	// 添加搜索条件
	if searchData.Username != "" {
		accountQueryBuilder.WhereUsername(mysql.LikePredicate, "%"+searchData.Username+"%")
	}
	if searchData.Name != "" {
		accountQueryBuilder.WhereName(mysql.LikePredicate, "%"+searchData.Name+"%")
	}
	if searchData.RoleType != "" {
		accountQueryBuilder.WhereRoleType(mysql.EqualPredicate, searchData.RoleType)
	}
	if searchData.Status != "" {
		accountQueryBuilder.WhereStatus(mysql.EqualPredicate, searchData.Status)
	}
	// 暂时移除组织搜索，因为新模型中没有相关字段
	// if searchData.BelongGroup != "" {
	// 	accountQueryBuilder.WhereBelongGroupName(mysql.LikePredicate, "%"+searchData.BelongGroup+"%")
	// }

	// 访问范围过滤
	scope, scopeErr := authz.ComputeScope(ctx, s.db)
	if scopeErr == nil && !scope.ScopeAll {
		if len(scope.AllowedAccountIDs) > 0 {
			accountQueryBuilder.WhereIdIn(scope.AllowedAccountIDs)
		} else if len(scope.AllowedTeamIDs) > 0 {
			accountQueryBuilder.WhereBelongTeamIdIn(scope.AllowedTeamIDs)
		} else if len(scope.AllowedGroupIDs) > 0 {
			accountQueryBuilder.WhereBelongGroupIdIn(scope.AllowedGroupIDs)
		}
	}

	// 查询总数
	total, err = accountQueryBuilder.Count(s.db.GetDbR())
	if err != nil {
		return 0, fmt.Errorf("查询账户总数失败: %v", err)
	}

	return total, nil
}

// Detail 获取账户详情
func (s *service) Detail(ctx core.Context, accountId string) (info *account.Account, err error) {
	// 转换accountId为int32
	var id int32
	fmt.Sscanf(accountId, "%d", &id)

	accountQueryBuilder := account.NewQueryBuilder()
	accountQueryBuilder.WhereId(mysql.EqualPredicate, id)

	// 查询账户详情
	info, err = accountQueryBuilder.QueryOne(s.db.GetDbR())
	if err != nil {
		return nil, fmt.Errorf("查询账户详情失败: %v", err)
	}
	if info == nil {
		return nil, fmt.Errorf("账户不存在")
	}

	// 范围校验
	scope, scopeErr := authz.ComputeScope(ctx, s.db)
	if scopeErr == nil && !authz.CanAccessAccount(scope, info) {
		return nil, fmt.Errorf("无权限访问该账户")
	}

	return info, nil
}
