package customer_authorization_record

import (
	"strconv"
	"strings"
	"time"

	"github.com/xinliangnote/go-gin-api/internal/pkg/core"
	model "github.com/xinliangnote/go-gin-api/internal/repository/mysql/customer_authorization_record"
)

func (s *service) List(ctx core.Context, req *ListRequest) (items []RecordItem, total int64, err error) {
	if req.Current <= 0 {
		req.Current = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}

	db := s.db.GetDbR().WithContext(ctx.RequestContext()).Model(&model.CustomerAuthorizationRecord{})
	if req.Name != "" {
		db = db.Where("name LIKE ?", "%"+req.Name+"%")
	}
	if len(req.City) > 0 {
		db = db.Where("city IN (?)", req.City)
	}
	if req.Phone != "" {
		like := strings.ReplaceAll(req.Phone, "*", "%")
		db = db.Where("phone LIKE ?", "%"+like+"%")
	}
	if len(req.AuthorizationStatus) > 0 {
		// map to boolean
		var vals []int
		for _, v := range req.AuthorizationStatus {
			if strings.ToLower(v) == "authorized" || v == "1" || strings.ToLower(v) == "true" {
				vals = append(vals, 1)
			} else if strings.ToLower(v) == "unauthorized" || v == "0" || strings.ToLower(v) == "false" {
				vals = append(vals, 0)
			}
		}
		if len(vals) > 0 {
			db = db.Where("is_authorized IN (?)", vals)
		}
	}
	if len(req.AssignmentStatus) > 0 {
		var vals []int
		for _, v := range req.AssignmentStatus {
			if strings.ToLower(v) == "assigned" || v == "1" || strings.ToLower(v) == "true" {
				vals = append(vals, 1)
			} else if strings.ToLower(v) == "unassigned" || v == "0" || strings.ToLower(v) == "false" {
				vals = append(vals, 0)
			}
		}
		if len(vals) > 0 {
			db = db.Where("is_assigned IN (?)", vals)
		}
	}
	if len(req.CompletionStatus) > 0 {
		var vals []int
		for _, v := range req.CompletionStatus {
			if strings.ToLower(v) == "complete" || v == "1" || strings.ToLower(v) == "true" {
				vals = append(vals, 1)
			} else if strings.ToLower(v) == "incomplete" || v == "0" || strings.ToLower(v) == "false" {
				vals = append(vals, 0)
			}
		}
		if len(vals) > 0 {
			db = db.Where("is_profile_complete IN (?)", vals)
		}
	}
	if len(req.PaymentStatus) > 0 {
		var vals []int
		for _, v := range req.PaymentStatus {
			if strings.ToLower(v) == "paid" || v == "1" || strings.ToLower(v) == "true" {
				vals = append(vals, 1)
			} else if strings.ToLower(v) == "unpaid" || v == "0" || strings.ToLower(v) == "false" {
				vals = append(vals, 0)
			}
		}
		if len(vals) > 0 {
			db = db.Where("is_paid IN (?)", vals)
		}
	}
	if req.BirthYearMin > 0 {
		db = db.Where("birth_year >= ?", req.BirthYearMin)
	}
	if req.BirthYearMax > 0 {
		db = db.Where("birth_year <= ?", req.BirthYearMax)
	}
	if req.HeightMin > 0 {
		db = db.Where("height >= ?", req.HeightMin)
	}
	if req.HeightMax > 0 {
		db = db.Where("height <= ?", req.HeightMax)
	}
	if req.IncomeMin > 0 {
		// income stored as varchar like "50w"; simple filter: extract leading digits
		db = db.Where("CAST(SUBSTRING_INDEX(income, 'w', 1) AS SIGNED) >= ?", req.IncomeMin)
	}
	if req.IncomeMax > 0 {
		db = db.Where("CAST(SUBSTRING_INDEX(income, 'w', 1) AS SIGNED) <= ?", req.IncomeMax)
	}
	if len(req.BelongGroup) > 0 {
		db = db.Where("`group` IN (?)", req.BelongGroup)
	}
	if len(req.BelongTeamId) > 0 {
		db = db.Where("team IN (?)", req.BelongTeamId)
	}
	if len(req.BelongAccountId) > 0 {
		db = db.Where("account IN (?)", req.BelongAccountId)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var rows []model.CustomerAuthorizationRecord
	if err := db.Order("id DESC").Limit(req.PageSize).Offset((req.Current - 1) * req.PageSize).Find(&rows).Error; err != nil {
		return nil, 0, err
	}

	out := make([]RecordItem, 0, len(rows))
	for _, r := range rows {
		out = append(out, *convertRow(&r))
	}
	return out, total, nil
}

func (s *service) Get(ctx core.Context, id string) (*RecordItem, error) {
	var r model.CustomerAuthorizationRecord
	if err := s.db.GetDbR().WithContext(ctx.RequestContext()).Where("id = ?", id).First(&r).Error; err != nil {
		return nil, err
	}
	return convertRow(&r), nil
}

func convertRow(r *model.CustomerAuthorizationRecord) *RecordItem {
	var photos []string
	if r.AuthPhotos != nil {
		photos = []string(*r.AuthPhotos)
	}
	// normalize group/team/account to pointers
	toPtr := func(s *string) *string { return s }
	created := time.Unix(r.CreatedTimestamp, 0).Format(time.RFC3339)
	updated := time.Unix(r.ModifiedTimestamp, 0).Format(time.RFC3339)
	idStr := strconv.Itoa(int(r.Id))
	return &RecordItem{
		Id:                idStr,
		Name:              r.Name,
		BirthYear:         (*int32)(r.BirthYear),
		Gender:            r.Gender,
		Height:            (*int32)(r.Height),
		City:              r.City,
		AuthStore:         r.AuthStore,
		Education:         r.Education,
		Profession:        r.Profession,
		Income:            r.Income,
		Phone:             r.Phone,
		Wechat:            r.Wechat,
		DrainageAccount:   r.DrainageAccount,
		DrainageId:        r.DrainageId,
		DrainageChannel:   r.DrainageChannel,
		Remark:            r.Remark,
		IsAuthorized:      r.IsAuthorized,
		AuthPhotos:        photos,
		IsProfileComplete: r.IsProfileComplete,
		IsAssigned:        r.IsAssigned,
		IsPaid:            r.IsPaid,
		PaymentAmount:     r.PaymentAmount,
		RefundAmount:      r.RefundAmount,
		Group:             toPtr(r.Group),
		Team:              toPtr(r.Team),
		Account:           toPtr(r.Account),
		CreatedAt:         created,
		UpdatedAt:         updated,
	}
}
