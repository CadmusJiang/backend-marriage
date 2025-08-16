package cooperation_store

import (
	"strconv"
	"time"

	"github.com/xinliangnote/go-gin-api/internal/pkg/core"
	model "github.com/xinliangnote/go-gin-api/internal/repository/mysql/cooperation_store"
)

func (s *service) List(ctx core.Context, req *ListRequest) (items []StoreItem, total int64, err error) {
	if req.Current <= 0 {
		req.Current = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}

	db := s.db.GetDbR().WithContext(ctx.RequestContext()).Model(&model.CooperationStore{})
	if req.StoreName != "" {
		db = db.Where("store_name LIKE ?", "%"+req.StoreName+"%")
	}
	if req.CooperationCity != "" {
		db = db.Where("cooperation_city_code = ?", req.CooperationCity)
	}
	if len(req.CooperationStatus) > 0 {
		db = db.Where("cooperation_status IN (?)", req.CooperationStatus)
	}
	if req.CompanyName != "" {
		db = db.Where("company_name LIKE ?", "%"+req.CompanyName+"%")
	}
	if req.StoreShortName != "" {
		db = db.Where("store_short_name LIKE ?", "%"+req.StoreShortName+"%")
	}
	if req.ActualBusinessAddress != "" {
		db = db.Where("actual_business_address LIKE ?", "%"+req.ActualBusinessAddress+"%")
	}
	// JSON 查询（包含任一）
	for i, t := range req.CooperationType {
		cond := "JSON_CONTAINS(cooperation_type, JSON_QUOTE(?))"
		if i == 0 {
			db = db.Where(cond, t)
		} else {
			db = db.Or(cond, t)
		}
	}
	for i, m := range req.CooperationMethod {
		cond := "JSON_CONTAINS(cooperation_method, JSON_QUOTE(?))"
		if i == 0 {
			db = db.Where(cond, m)
		} else {
			db = db.Or(cond, m)
		}
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var rows []model.CooperationStore
	if err := db.Order("id DESC").Limit(req.PageSize).Offset((req.Current - 1) * req.PageSize).Find(&rows).Error; err != nil {
		return nil, 0, err
	}

	out := make([]StoreItem, 0, len(rows))
	for _, r := range rows {
		var coopType []string
		if r.CooperationType != nil {
			coopType = []string(*r.CooperationType)
		}
		var coopMethod []string
		if r.CooperationMethod != nil {
			coopMethod = []string(*r.CooperationMethod)
		}
		var storePhotos []string
		if r.StorePhotos != nil {
			storePhotos = []string(*r.StorePhotos)
		}
		var contractPhotos []string
		if r.ContractPhotos != nil {
			contractPhotos = []string(*r.ContractPhotos)
		}

		out = append(out, StoreItem{
			Id:                    strconv.Itoa(int(r.Id)),
			StoreName:             r.StoreName,
			CooperationCity:       r.CooperationCityCode,
			CooperationType:       coopType,
			StoreShortName:        r.StoreShortName,
			CompanyName:           r.CompanyName,
			CooperationMethod:     coopMethod,
			CooperationStatus:     r.CooperationStatus,
			BusinessLicense:       r.BusinessLicense,
			StorePhotos:           storePhotos,
			ActualBusinessAddress: r.ActualBusinessAddress,
			ContractPhotos:        contractPhotos,
			CreatedAt:             r.CreatedAt.Format(time.RFC3339),
			UpdatedAt:             r.UpdatedAt.Format(time.RFC3339),
		})
	}
	return out, total, nil
}
