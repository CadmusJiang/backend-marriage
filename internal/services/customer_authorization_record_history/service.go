package customer_authorization_record_history

import (
	"github.com/xinliangnote/go-gin-api/internal/pkg/core"
	"github.com/xinliangnote/go-gin-api/internal/repository/mysql"
)

type Service interface {
	i()
	List(ctx core.Context, req *ListRequest) (items []HistoryItem, total int64, err error)
	Get(ctx core.Context, historyId string) (*HistoryItem, error)
}

type service struct{ db mysql.Repo }

func New(db mysql.Repo) Service { return &service{db: db} }
func (s *service) i()           {}

type ListRequest struct {
	Current, PageSize             int
	CustomerAuthorizationRecordId string
	OperateType                   string
	OperatorUsername              string
	OperatorRoleType              string
}

type HistoryItem struct {
	Id                            string
	HistoryId                     string
	CustomerAuthorizationRecordId string
	OperateType                   string
	OperateTimestamp              int64
	Content                       map[string]interface{}
	OperatorUsername              string
	OperatorNickname              string
	OperatorRoleType              string
	CreatedAt                     string
	UpdatedAt                     string
}

// Stub implementations; real SQL to be added when needed
func (s *service) List(ctx core.Context, req *ListRequest) ([]HistoryItem, int64, error) {
	return []HistoryItem{}, 0, nil
}

func (s *service) Get(ctx core.Context, historyId string) (*HistoryItem, error) {
	return nil, nil
}
