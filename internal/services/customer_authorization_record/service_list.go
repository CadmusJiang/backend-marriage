package customer_authorization_record

import (
	"strconv"
	"strings"

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
		// map to enum values
		var vals []string
		for _, v := range req.AuthorizationStatus {
			if strings.ToLower(v) == "authorized" {
				vals = append(vals, "authorized")
			} else if strings.ToLower(v) == "unauthorized" {
				vals = append(vals, "unauthorized")
			}
		}
		if len(vals) > 0 {
			db = db.Where("authorization_status IN (?)", vals)
		}
	}
	if len(req.AssignmentStatus) > 0 {
		var vals []string
		for _, v := range req.AssignmentStatus {
			if strings.ToLower(v) == "assigned" {
				vals = append(vals, "assigned")
			} else if strings.ToLower(v) == "unassigned" {
				vals = append(vals, "unassigned")
			}
		}
		if len(vals) > 0 {
			db = db.Where("assignment_status IN (?)", vals)
		}
	}
	if len(req.CompletionStatus) > 0 {
		var vals []string
		for _, v := range req.CompletionStatus {
			if strings.ToLower(v) == "complete" {
				vals = append(vals, "complete")
			} else if strings.ToLower(v) == "incomplete" {
				vals = append(vals, "incomplete")
			}
		}
		if len(vals) > 0 {
			db = db.Where("completion_status IN (?)", vals)
		}
	}
	if len(req.PaymentStatus) > 0 {
		var vals []string
		for _, v := range req.PaymentStatus {
			if strings.ToLower(v) == "paid" {
				vals = append(vals, "paid")
			} else if strings.ToLower(v) == "unpaid" {
				vals = append(vals, "unpaid")
			}
		}
		if len(vals) > 0 {
			db = db.Where("payment_status IN (?)", vals)
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
	// Convert group/team/account strings to uint64 pointers
	// Convert int32 to int for BirthYear and Height
	var birthYear *int
	if r.BirthYear != nil {
		year := int(*r.BirthYear)
		birthYear = &year
	}
	var height *int
	if r.Height != nil {
		h := int(*r.Height)
		height = &h
	}

	// Convert photos array to string
	var authPhotosStr *string
	if len(photos) > 0 {
		photosStr := strings.Join(photos, ",")
		authPhotosStr = &photosStr
	}

	// Convert group/team/account strings to uint64 pointers
	var belongGroupID *uint64
	if r.Group != nil {
		if groupID, err := strconv.ParseUint(*r.Group, 10, 64); err == nil {
			belongGroupID = &groupID
		}
	}
	var belongTeamID *uint64
	if r.Team != nil {
		if teamID, err := strconv.ParseUint(*r.Team, 10, 64); err == nil {
			belongTeamID = &teamID
		}
	}
	var belongAccountID *uint64
	if r.Account != nil {
		if accountID, err := strconv.ParseUint(*r.Account, 10, 64); err == nil {
			belongAccountID = &accountID
		}
	}

	return &RecordItem{
		ID:                  uint64(r.Id),
		Name:                r.Name,
		BirthYear:           birthYear,
		Gender:              r.Gender,
		Height:              height,
		City:                r.CityCode,
		AuthStore:           r.AuthStore,
		Education:           r.Education,
		Profession:          r.Profession,
		Income:              r.Income,
		Phone:               r.Phone,
		Wechat:              r.Wechat,
		DrainageAccount:     r.DrainageAccount,
		DrainageId:          r.DrainageId,
		DrainageChannel:     r.DrainageChannel,
		Remark:              r.Remark,
		AuthorizationStatus: r.AuthorizationStatus,
		AuthPhotos:          authPhotosStr,
		CompletionStatus:    r.CompletionStatus,
		AssignmentStatus:    r.AssignmentStatus,
		PaymentStatus:       r.PaymentStatus,
		PaymentAmount:       r.PaymentAmount,
		RefundAmount:        r.RefundAmount,
		BelongGroupID:       belongGroupID,
		BelongTeamID:        belongTeamID,
		BelongAccountID:     belongAccountID,
		CreatedAt:           r.CreatedAt,
		UpdatedAt:           r.UpdatedAt,
	}
}
