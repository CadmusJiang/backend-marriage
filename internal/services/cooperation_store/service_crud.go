package cooperation_store

import (
	"fmt"
	"strconv"
	"time"

	"github.com/xinliangnote/go-gin-api/internal/pkg/core"
	model "github.com/xinliangnote/go-gin-api/internal/repository/mysql/cooperation_store"
)

func (s *service) Get(ctx core.Context, id string) (*StoreItem, error) {
	var row model.CooperationStore
	if err := s.db.GetDbR().WithContext(ctx.RequestContext()).Where("id = ?", id).First(&row).Error; err != nil {
		return nil, err
	}
	return convertRow(&row), nil
}

func (s *service) Create(ctx core.Context, req *CreateRequest) (*StoreItem, error) {

	m := &model.CooperationStore{
		StoreName:             req.StoreName,
		CooperationCityCode:   req.CooperationCity,
		StoreShortName:        req.StoreShortName,
		CompanyName:           req.CompanyName,
		CooperationStatus:     req.CooperationStatus,
		BusinessLicense:       req.BusinessLicense,
		ActualBusinessAddress: req.ActualBusinessAddress,
		CreatedAt:             time.Now(),
		UpdatedAt:             time.Now(),
		CreatedUser:           ctx.SessionUserInfo().UserName,
		UpdatedUser:           ctx.SessionUserInfo().UserName,
	}
	if len(req.CooperationType) > 0 {
		ct := model.JSONString(req.CooperationType)
		m.CooperationType = &ct
	}
	if len(req.CooperationMethod) > 0 {
		cm := model.JSONString(req.CooperationMethod)
		m.CooperationMethod = &cm
	}
	if len(req.StorePhotos) > 0 {
		sp := model.JSONString(req.StorePhotos)
		m.StorePhotos = &sp
	}
	if len(req.ContractPhotos) > 0 {
		cp := model.JSONString(req.ContractPhotos)
		m.ContractPhotos = &cp
	}

	if err := s.db.GetDbW().WithContext(ctx.RequestContext()).Create(m).Error; err != nil {
		return nil, err
	}
	return convertRow(m), nil
}

func (s *service) Update(ctx core.Context, id string, req *UpdateRequest) (*StoreItem, error) {
	var exists int64
	if err := s.db.GetDbR().WithContext(ctx.RequestContext()).Model(&model.CooperationStore{}).Where("id = ?", id).Count(&exists).Error; err != nil {
		return nil, err
	}
	if exists == 0 {
		return nil, fmt.Errorf("not_found")
	}

	updates := map[string]interface{}{
		"updated_at":   time.Now(),
		"updated_user": ctx.SessionUserInfo().UserName,
	}
	if req.StoreName != "" {
		updates["store_name"] = req.StoreName
	}
	if req.CooperationCity != "" {
		updates["cooperation_city_code"] = req.CooperationCity
	}
	if req.CooperationStatus != "" {
		updates["cooperation_status"] = req.CooperationStatus
	}
	if req.StoreShortName != nil {
		updates["store_short_name"] = req.StoreShortName
	}
	if req.CompanyName != nil {
		updates["company_name"] = req.CompanyName
	}
	if req.ActualBusinessAddress != nil {
		updates["actual_business_address"] = req.ActualBusinessAddress
	}
	if req.BusinessLicense != nil {
		updates["business_license"] = req.BusinessLicense
	}
	if req.CooperationType != nil {
		jt := model.JSONString(req.CooperationType)
		updates["cooperation_type"] = jt
	}
	if req.CooperationMethod != nil {
		jm := model.JSONString(req.CooperationMethod)
		updates["cooperation_method"] = jm
	}
	if req.StorePhotos != nil {
		js := model.JSONString(req.StorePhotos)
		updates["store_photos"] = js
	}
	if req.ContractPhotos != nil {
		jc := model.JSONString(req.ContractPhotos)
		updates["contract_photos"] = jc
	}

	if err := s.db.GetDbW().WithContext(ctx.RequestContext()).Model(&model.CooperationStore{}).Where("id = ?", id).Updates(updates).Error; err != nil {
		return nil, err
	}

	var out model.CooperationStore
	if err := s.db.GetDbR().WithContext(ctx.RequestContext()).Where("id = ?", id).First(&out).Error; err != nil {
		return nil, err
	}
	return convertRow(&out), nil
}

func convertRow(r *model.CooperationStore) *StoreItem {
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

	return &StoreItem{
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
	}
}
