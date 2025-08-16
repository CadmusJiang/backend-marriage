package customer_authorization_record

import (
	"fmt"
	"strings"
	"time"

	"github.com/xinliangnote/go-gin-api/internal/pkg/core"
	"github.com/xinliangnote/go-gin-api/internal/repository/mysql"
)

type Service interface {
	i()
	PageList(ctx core.Context, req *ListRequest) (items []RecordItem, total int64, err error)
	GetByID(ctx core.Context, id uint64) (*RecordItem, error)
}

type service struct {
	db mysql.Repo
}

func New(db mysql.Repo) Service {
	return &service{db: db}
}

func (s *service) i() {}

type ListRequest struct {
	Current             int
	PageSize            int
	Name                string
	City                []string
	Phone               string
	AuthorizationStatus []string
	AssignmentStatus    []string
	CompletionStatus    []string
	PaymentStatus       []string
	BirthYearMin        int
	BirthYearMax        int
	HeightMin           int
	HeightMax           int
	IncomeMin           int
	IncomeMax           int
	BelongGroup         []string
	BelongTeamId        []string
	BelongAccountId     []string
	BelongGroupIDs      []uint64
	BelongTeamIDs       []uint64
	BelongAccountIDs    []uint64
}

type RecordItem struct {
	ID                  uint64
	Name                string
	BirthYear           *int
	Gender              *string
	Height              *int
	City                *string
	AuthStore           *string
	Education           *string
	Profession          *string
	Income              *string
	Phone               *string
	Wechat              *string
	DrainageAccount     *string
	DrainageId          *string
	DrainageChannel     *string
	Remark              *string
	AuthorizationStatus string
	AuthPhotos          *string
	CompletionStatus    string
	AssignmentStatus    string
	PaymentStatus       string
	PaymentAmount       float64
	RefundAmount        float64
	BelongGroupID       *uint64
	BelongTeamID        *uint64
	BelongAccountID     *uint64
	CreatedAt           time.Time
	UpdatedAt           time.Time
}

func (s *service) PageList(ctx core.Context, req *ListRequest) (items []RecordItem, total int64, err error) {
	db := s.db.GetDbR().WithContext(ctx.RequestContext()).Table("customer_authorization_record")
	if req.Name != "" {
		db = db.Where("name LIKE ?", "%"+req.Name+"%")
	}
	if len(req.City) > 0 {
		db = db.Where("city_code IN (?)", req.City)
	}
	if len(req.AuthorizationStatus) > 0 {
		var vals []string
		set := strings.Join(req.AuthorizationStatus, ",")
		if strings.Contains(set, "authorized") {
			vals = append(vals, "authorized")
		}
		if strings.Contains(set, "unauthorized") {
			vals = append(vals, "unauthorized")
		}
		if len(vals) == 1 {
			db = db.Where("authorization_status = ?", vals[0])
		} else if len(vals) > 1 {
			db = db.Where("authorization_status IN (?)", vals)
		}
	}
	if len(req.PaymentStatus) > 0 {
		var vals []string
		set := strings.Join(req.PaymentStatus, ",")
		if strings.Contains(set, "paid") {
			vals = append(vals, "paid")
		}
		if strings.Contains(set, "unpaid") {
			vals = append(vals, "unpaid")
		}
		if len(vals) == 1 {
			db = db.Where("payment_status = ?", vals[0])
		} else if len(vals) > 1 {
			db = db.Where("payment_status IN (?)", vals)
		}
	}
	if len(req.AssignmentStatus) > 0 {
		var vals []string
		set := strings.Join(req.AssignmentStatus, ",")
		if strings.Contains(set, "assigned") {
			vals = append(vals, "assigned")
		}
		if strings.Contains(set, "unassigned") {
			vals = append(vals, "unassigned")
		}
		if len(vals) == 1 {
			db = db.Where("assignment_status = ?", vals[0])
		} else if len(vals) > 1 {
			db = db.Where("assignment_status IN (?)", vals)
		}
	}
	if len(req.CompletionStatus) > 0 {
		var vals []string
		set := strings.Join(req.CompletionStatus, ",")
		if strings.Contains(set, "complete") {
			vals = append(vals, "complete")
		}
		if strings.Contains(set, "incomplete") {
			vals = append(vals, "incomplete")
		}
		if len(vals) == 1 {
			db = db.Where("completion_status = ?", vals[0])
		} else if len(vals) > 1 {
			db = db.Where("completion_status IN (?)", vals)
		}
	}
	if len(req.BelongGroupIDs) > 0 {
		db = db.Where("belong_group_id IN (?)", req.BelongGroupIDs)
	}
	if len(req.BelongTeamIDs) > 0 {
		db = db.Where("belong_team_id IN (?)", req.BelongTeamIDs)
	}
	if len(req.BelongAccountIDs) > 0 {
		db = db.Where("belong_account_id IN (?)", req.BelongAccountIDs)
	}
	if err = db.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count err: %v", err)
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}
	if req.Current <= 0 {
		req.Current = 1
	}
	var rows []RecordItem
	if err = db.Select("id, name, birth_year, gender, height, city_code, auth_store, education, profession, income, phone, wechat, drainage_account, drainage_id, drainage_channel, remark, authorization_status, auth_photos, completion_status, assignment_status, payment_status, payment_amount, refund_amount, belong_group_id, belong_team_id, belong_account_id, created_at, updated_at").
		Order("created_at DESC").Limit(req.PageSize).Offset((req.Current - 1) * req.PageSize).Find(&rows).Error; err != nil {
		return nil, 0, fmt.Errorf("query err: %v", err)
	}
	return rows, total, nil
}

func (s *service) GetByID(ctx core.Context, id uint64) (*RecordItem, error) {
	var row RecordItem
	if err := s.db.GetDbR().WithContext(ctx.RequestContext()).Table("customer_authorization_record").
		Select("id, name, birth_year, gender, height, city_code, auth_store, education, profession, income, phone, wechat, drainage_account, drainage_id, drainage_channel, remark, authorization_status, auth_photos, completion_status, assignment_status, payment_status, payment_amount, refund_amount, belong_group_id, belong_team_id, belong_account_id, created_at, updated_at").
		Where("id = ?", id).Take(&row).Error; err != nil {
		return nil, err
	}
	return &row, nil
}
