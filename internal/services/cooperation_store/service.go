package cooperation_store

import (
	"github.com/xinliangnote/go-gin-api/internal/pkg/core"
	"github.com/xinliangnote/go-gin-api/internal/repository/mysql"
)

type Service interface {
	i()

	List(ctx core.Context, req *ListRequest) (items []StoreItem, total int64, err error)
	Get(ctx core.Context, id string) (item *StoreItem, err error)
	Create(ctx core.Context, req *CreateRequest) (item *StoreItem, err error)
	Update(ctx core.Context, id string, req *UpdateRequest) (item *StoreItem, err error)
}

type service struct {
	db mysql.Repo
}

func New(db mysql.Repo) Service {
	return &service{db: db}
}

func (s *service) i() {}

// DTOs used by handler↔service layer
type ListRequest struct {
	Current               int
	PageSize              int
	StoreName             string
	CooperationCity       string
	CooperationType       []string
	StoreShortName        string
	CompanyName           string
	CooperationMethod     []string
	CooperationStatus     []string
	ActualBusinessAddress string
}

type CreateRequest struct {
	StoreName             string
	CooperationCity       string
	CooperationType       []string
	StoreShortName        *string
	CompanyName           *string
	CooperationMethod     []string
	CooperationStatus     string
	BusinessLicense       *string
	StorePhotos           []string
	ActualBusinessAddress *string
	ContractPhotos        []string
}

type UpdateRequest = CreateRequest

type StoreItem struct {
	Id                    string
	StoreName             string
	CooperationCity       string
	CooperationType       []string
	StoreShortName        *string
	CompanyName           *string
	CooperationMethod     []string
	CooperationStatus     string
	BusinessLicense       *string
	StorePhotos           []string
	ActualBusinessAddress *string
	ContractPhotos        []string
	CreatedAt             string
	UpdatedAt             string
}
