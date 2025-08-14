package customer_authorization_record

import (
	"encoding/json"
	"time"

	"github.com/xinliangnote/go-gin-api/internal/pkg/core"
	model "github.com/xinliangnote/go-gin-api/internal/repository/mysql/customer_authorization_record"
	outbox "github.com/xinliangnote/go-gin-api/internal/repository/mysql/outbox"
	"gorm.io/gorm"
)

func (s *service) Create(ctx core.Context, req *CreateRequest) (*RecordItem, error) {
	now := time.Now().Unix()
	m := &model.CustomerAuthorizationRecord{
		Name:              req.Name,
		BirthYear:         (*int32)(req.BirthYear),
		Gender:            req.Gender,
		Height:            (*int32)(req.Height),
		City:              req.City,
		Education:         req.Education,
		Profession:        req.Profession,
		Income:            (*string)(req.Income),
		Phone:             req.Phone,
		Wechat:            req.Wechat,
		DrainageAccount:   req.DrainageAccount,
		DrainageId:        req.DrainageId,
		DrainageChannel:   req.DrainageChannel,
		Remark:            req.Remark,
		IsAuthorized:      req.AuthorizationStatus,
		IsProfileComplete: req.CompletionStatus,
		IsAssigned:        req.AssignmentStatus,
		IsPaid:            req.PaymentStatus,
		PaymentAmount:     req.PaymentAmount,
		RefundAmount:      req.RefundAmount,
		Group:             req.AuthorizedStore, // mapping fields as needed
		Team:              req.BelongTeam,
		Account:           req.BelongAccount,
		CreatedAt:         now,
		UpdatedAt:         now,
		CreatedUser:       ctx.SessionUserInfo().UserName,
		UpdatedUser:       ctx.SessionUserInfo().UserName,
	}
	if len(req.AuthorizationPhotos) > 0 {
		jp := model.JSONString(req.AuthorizationPhotos)
		m.AuthPhotos = &jp
	}
	// write record + outbox within one transaction
	if err := s.db.GetDbW().WithContext(ctx.RequestContext()).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(m).Error; err != nil {
			return err
		}
		// after snapshot
		afterBytes, _ := json.Marshal(m)
		payload := map[string]interface{}{
			"type":       "customer.created",
			"recordId":   m.Id,
			"occurredAt": time.Now().Unix(),
			"before":     nil,
			"after":      json.RawMessage(afterBytes),
		}
		pbytes, _ := json.Marshal(payload)
		evt := &outbox.Event{Topic: "stream.customer.record", Payload: pbytes, Status: 0}
		return tx.Create(evt).Error
	}); err != nil {
		return nil, err
	}
	return convertRow(m), nil
}

func (s *service) Update(ctx core.Context, id string, req *UpdateRequest) (*RecordItem, error) {
	updates := map[string]interface{}{
		"updated_at":   time.Now().Unix(),
		"updated_user": ctx.SessionUserInfo().UserName,
	}
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.BirthYear != nil {
		updates["birth_year"] = *req.BirthYear
	}
	if req.Gender != nil {
		updates["gender"] = *req.Gender
	}
	if req.Height != nil {
		updates["height"] = *req.Height
	}
	if req.City != nil {
		updates["city"] = *req.City
	}
	if req.Education != nil {
		updates["education"] = *req.Education
	}
	if req.Profession != nil {
		updates["profession"] = *req.Profession
	}
	if req.Income != nil {
		updates["income"] = *req.Income
	}
	if req.Phone != nil {
		updates["phone"] = *req.Phone
	}
	if req.Wechat != nil {
		updates["wechat"] = *req.Wechat
	}
	if req.DrainageAccount != nil {
		updates["drainage_account"] = *req.DrainageAccount
	}
	if req.DrainageId != nil {
		updates["drainage_id"] = *req.DrainageId
	}
	if req.DrainageChannel != nil {
		updates["drainage_channel"] = *req.DrainageChannel
	}
	if req.Remark != nil {
		updates["remark"] = *req.Remark
	}
	updates["is_authorized"] = req.AuthorizationStatus
	updates["is_assigned"] = req.AssignmentStatus
	updates["is_profile_complete"] = req.CompletionStatus
	updates["is_paid"] = req.PaymentStatus
	updates["payment_amount"] = req.PaymentAmount
	updates["refund_amount"] = req.RefundAmount
	if req.AuthorizedStore != nil {
		updates["auth_store"] = *req.AuthorizedStore
	}
	if req.BelongTeam != nil {
		updates["team"] = *req.BelongTeam
	}
	if req.BelongAccount != nil {
		updates["account"] = *req.BelongAccount
	}
	if req.AuthorizationPhotos != nil {
		jp := model.JSONString(req.AuthorizationPhotos)
		updates["auth_photos"] = jp
	}

	// transactional: load before, update, load after, enqueue outbox
	var out *RecordItem
	if err := s.db.GetDbW().WithContext(ctx.RequestContext()).Transaction(func(tx *gorm.DB) error {
		var before model.CustomerAuthorizationRecord
		if err := tx.Where("id = ?", id).First(&before).Error; err != nil {
			return err
		}
		beforeBytes, _ := json.Marshal(&before)

		if err := tx.Model(&model.CustomerAuthorizationRecord{}).Where("id = ?", id).Updates(updates).Error; err != nil {
			return err
		}
		var after model.CustomerAuthorizationRecord
		if err := tx.Where("id = ?", id).First(&after).Error; err != nil {
			return err
		}
		afterBytes, _ := json.Marshal(&after)

		payload := map[string]interface{}{
			"type":       "customer.updated",
			"recordId":   after.Id,
			"occurredAt": time.Now().Unix(),
			"before":     json.RawMessage(beforeBytes),
			"after":      json.RawMessage(afterBytes),
		}
		pbytes, _ := json.Marshal(payload)
		evt := &outbox.Event{Topic: "stream.customer.record", Payload: pbytes, Status: 0}
		if err := tx.Create(evt).Error; err != nil {
			return err
		}
		r := convertRow(&after)
		out = r
		return nil
	}); err != nil {
		return nil, err
	}
	return out, nil
}
