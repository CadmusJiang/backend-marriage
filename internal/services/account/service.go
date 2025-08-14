package account

import (
	"github.com/xinliangnote/go-gin-api/internal/pkg/core"
	"github.com/xinliangnote/go-gin-api/internal/repository/mysql"
	"github.com/xinliangnote/go-gin-api/internal/repository/mysql/account"
	"github.com/xinliangnote/go-gin-api/internal/repository/mysql/account_history"
	"github.com/xinliangnote/go-gin-api/internal/repository/redis"
)

var _ Service = (*service)(nil)

// CreateAccountData 创建账户数据
type CreateAccountData struct {
	Username    string   `json:"username" binding:"required"`
	Nickname    string   `json:"nickname" binding:"required"`
	Password    string   `json:"password" binding:"required"`
	Phone       string   `json:"phone"`
	RoleType    string   `json:"roleType"`
	Status      string   `json:"status"`
	BelongGroup *OrgInfo `json:"belongGroup"`
	BelongTeam  *OrgInfo `json:"belongTeam"`
}

// UpdateAccountData 更新账户数据
type UpdateAccountData struct {
	Nickname    string   `json:"nickname"`
	Phone       string   `json:"phone"`
	Status      string   `json:"status"`
	BelongGroup *OrgInfo `json:"belongGroup"`
	BelongTeam  *OrgInfo `json:"belongTeam"`
}

// OrgInfo 组织信息
type OrgInfo struct {
	ID                int32  `json:"id"`
	Username          string `json:"username"`
	Nickname          string `json:"nickname"`
	CreatedTimestamp  int64  `json:"createdTimestamp"`
	ModifiedTimestamp int64  `json:"modifiedTimestamp"`
	CurrentCnt        int32  `json:"currentCnt"`
}

// SearchData 搜索数据
type SearchData struct {
	Username    string `form:"username"`
	Nickname    string `form:"nickname"`
	RoleType    string `form:"roleType"`
	Status      string `form:"status"`
	BelongGroup string `form:"belongGroup"`
	Current     int    `form:"current,default=1"`
	PageSize    int    `form:"pageSize,default=10"`
}

// CreateHistoryData 创建历史记录数据
type CreateHistoryData struct {
	AccountId        string `json:"accountId"`
	OperateType      string `json:"operateType"`
	Content          string `json:"content"`
	Operator         string `json:"operator"`
	OperatorRoleType string `json:"operatorRoleType"`
}

// SearchHistoryData 搜索历史记录数据
type SearchHistoryData struct {
	AccountId       string `form:"accountId"`
	AccountUsername string `form:"accountUsername"`
	OperateType     string `form:"operateType"`
	Current         int    `form:"current,default=1"`
	PageSize        int    `form:"pageSize,default=10"`
}

type Service interface {
	i()

	// 账户管理
	Create(ctx core.Context, accountData *CreateAccountData) (id int32, err error)
	PageList(ctx core.Context, searchData *SearchData) (listData []*account.Account, err error)
	PageListCount(ctx core.Context, searchData *SearchData) (total int64, err error)
	Detail(ctx core.Context, accountId string) (info *account.Account, err error)
	Update(ctx core.Context, accountId string, updateData *UpdateAccountData) (err error)
	Delete(ctx core.Context, accountId string) (err error)

	// 认证相关
	Login(ctx core.Context, username, password string) (accountInfo *account.Account, err error)
	Logout(ctx core.Context, token string) (err error)

	// 历史记录
	CreateHistory(ctx core.Context, historyData *CreateHistoryData) (err error)
	PageListHistory(ctx core.Context, searchData *SearchHistoryData) (listData []*account_history.AccountHistory, err error)
	PageListHistoryCount(ctx core.Context, searchData *SearchHistoryData) (total int64, err error)

	// 仅更新密码
	UpdatePassword(ctx core.Context, accountId string, newPassword string) (err error)
}

type service struct {
	db    mysql.Repo
	cache redis.Repo
}

func New(db mysql.Repo, cache redis.Repo) Service {
	return &service{
		db:    db,
		cache: cache,
	}
}

func (s *service) i() {}
