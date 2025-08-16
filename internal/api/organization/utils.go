package organization

import (
	"github.com/xinliangnote/go-gin-api/internal/pkg/core"
	"github.com/xinliangnote/go-gin-api/internal/proposal"
	"github.com/xinliangnote/go-gin-api/internal/repository/mysql"
	"github.com/xinliangnote/go-gin-api/internal/repository/mysql/account"
	orgsvc "github.com/xinliangnote/go-gin-api/internal/services/organization"
)

// createContext 创建service层的Context
func (h *handler) createContext(c core.Context) orgsvc.Context {
	return &serviceContext{
		coreContext: c,
	}
}

// serviceContext 实现orgsvc.Context接口
type serviceContext struct {
	coreContext core.Context
}

func (sc *serviceContext) RequestContext() core.StdContext {
	return sc.coreContext.RequestContext()
}

func (sc *serviceContext) SessionUserInfo() proposal.SessionUserInfo {
	return sc.coreContext.SessionUserInfo()
}

// getUserAccountInfo 获取用户账户信息
func (h *handler) getUserAccountInfo(userID int) (*account.Account, error) {
	qb := account.NewQueryBuilder()
	qb.WhereId(mysql.EqualPredicate, int32(userID))
	return qb.QueryOne(h.db.GetDbR())
}
