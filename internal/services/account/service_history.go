package account

import (
	"fmt"
	"strconv"
	"time"

	"github.com/xinliangnote/go-gin-api/internal/pkg/core"
	"github.com/xinliangnote/go-gin-api/internal/repository/mysql"
	"github.com/xinliangnote/go-gin-api/internal/repository/mysql/account_history"
)

// CreateHistory 创建历史记录
func (s *service) CreateHistory(ctx core.Context, historyData *CreateHistoryData) (err error) {
	// 创建历史记录
	now := uint64(time.Now().Unix())
	accountId, _ := strconv.ParseUint(historyData.AccountId, 10, 64)
	newHistory := &account_history.AccountHistory{
		AccountId:         accountId,
		OperateType:       historyData.OperateType,
		OperateTimestamp:  now,
		Content:           historyData.Content,
		Operator:          historyData.Operator,
		OperatorRoleType:  historyData.OperatorRoleType,
		CreatedTimestamp:  now,
		ModifiedTimestamp: now,
		CreatedUser:       ctx.SessionUserInfo().UserName,
		UpdatedUser:       ctx.SessionUserInfo().UserName,
	}

	// 保存到数据库
	_, err = newHistory.Create(s.db.GetDbW())
	if err != nil {
		return fmt.Errorf("创建历史记录失败: %v", err)
	}

	return nil
}

// PageListHistory 分页获取历史记录列表
func (s *service) PageListHistory(ctx core.Context, searchData *SearchHistoryData) (listData []*account_history.AccountHistory, err error) {
	historyQueryBuilder := account_history.NewQueryBuilder()

	// 添加搜索条件
	if searchData.AccountId != "" {
		accountId, _ := strconv.ParseUint(searchData.AccountId, 10, 64)
		historyQueryBuilder.WhereAccountId(mysql.EqualPredicate, accountId)
	}
	if searchData.OperateType != "" {
		historyQueryBuilder.WhereOperateType(mysql.EqualPredicate, searchData.OperateType)
	}

	// 设置分页
	offset := (searchData.Current - 1) * searchData.PageSize
	historyQueryBuilder.Limit(searchData.PageSize).Offset(offset)

	// 按时间倒序排列
	historyQueryBuilder.OrderByOperateTimestamp(false)

	// 查询数据
	listData, err = historyQueryBuilder.QueryAll(s.db.GetDbR())
	if err != nil {
		return nil, fmt.Errorf("查询历史记录失败: %v", err)
	}

	return listData, nil
}

// PageListHistoryCount 获取历史记录总数
func (s *service) PageListHistoryCount(ctx core.Context, searchData *SearchHistoryData) (total int64, err error) {
	historyQueryBuilder := account_history.NewQueryBuilder()

	// 添加搜索条件
	if searchData.AccountId != "" {
		accountId, _ := strconv.ParseUint(searchData.AccountId, 10, 64)
		historyQueryBuilder.WhereAccountId(mysql.EqualPredicate, accountId)
	}
	if searchData.OperateType != "" {
		historyQueryBuilder.WhereOperateType(mysql.EqualPredicate, searchData.OperateType)
	}

	// 查询总数
	total, err = historyQueryBuilder.Count(s.db.GetDbR())
	if err != nil {
		return 0, fmt.Errorf("查询历史记录总数失败: %v", err)
	}

	return total, nil
}
