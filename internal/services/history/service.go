package history

import (
	"encoding/json"
	"fmt"
	"reflect"
	"time"

	"go.uber.org/zap"
)

// Database 数据库接口
type Database interface {
	GetDbW() interface {
		WithContext(ctx interface{}) interface {
			Table(name string) interface {
				Create(value interface{}) interface {
					Error() error
				}
			}
		}
	}
}

// Context 上下文接口
type Context interface {
	RequestContext() interface{}
}

// HistoryService 通用历史记录服务
type HistoryService struct {
	db interface {
		GetDbW() interface {
			WithContext(ctx interface{}) interface {
				Table(name string) interface {
					Create(value interface{}) interface {
						Error() error
					}
				}
			}
		}
	}
}

// NewHistoryService 创建历史记录服务
func NewHistoryService(db interface {
	GetDbW() interface {
		WithContext(ctx interface{}) interface {
			Table(name string) interface {
				Create(value interface{}) interface {
					Error() error
				}
			}
		}
	}
}) *HistoryService {
	return &HistoryService{db: db}
}

// HistoryRecord 历史记录结构
type HistoryRecord struct {
	EntityID         uint32                 `json:"entity_id"`          // 实体ID
	EntityType       EntityType             `json:"entity_type"`        // 实体类型 (account, org, cooperation_store, etc.)
	OperateType      OperateType            `json:"operate_type"`       // 操作类型 (created, updated, deleted)
	OperatedAt       time.Time              `json:"operated_at"`        // 操作时间
	Content          map[string]interface{} `json:"content"`            // 操作内容
	Operator         string                 `json:"operator"`           // 操作人
	OperatorRoleType string                 `json:"operator_role_type"` // 操作人角色类型
}

// CreateHistory 创建历史记录
func (s *HistoryService) CreateHistory(ctx interface {
	RequestContext() interface{}
}, record *HistoryRecord) error {
	// 验证枚举值
	if !record.EntityType.IsValid() {
		return fmt.Errorf("invalid entity type: %s", record.EntityType)
	}
	if !record.OperateType.IsValid() {
		return fmt.Errorf("invalid operate type: %s", record.OperateType)
	}

	// 根据实体类型确定表名
	tableName := record.EntityType.GetHistoryTableName()
	if tableName == "" {
		return fmt.Errorf("unsupported entity type: %s", record.EntityType)
	}

	// 构造统一的历史记录数据
	historyData := map[string]interface{}{
		"operate_type":       record.OperateType.String(),
		"operated_at":        record.OperatedAt,
		"content":            s.formatContent(record.Content),
		"operator_username":  record.Operator,
		"operator_name":      record.Operator,
		"operator_role_type": record.OperatorRoleType,
		"created_at":         record.OperatedAt,
		"updated_at":         record.OperatedAt,
		"created_user":       record.Operator,
		"updated_user":       record.Operator,
	}

	// 根据实体类型设置对应的ID字段
	switch record.EntityType {
	case EntityTypeAccount:
		historyData["account_id"] = record.EntityID
	case EntityTypeOrg, EntityTypeOrganization:
		historyData["org_id"] = record.EntityID
	case EntityTypeCooperationStore:
		historyData["store_id"] = record.EntityID
	case EntityTypeCustomerAuthorizationRecord:
		historyData["customer_authorization_record_id"] = record.EntityID
	default:
		return fmt.Errorf("unsupported entity type: %s", record.EntityType)
	}

	// 创建历史记录
	result := s.db.GetDbW().WithContext(ctx.RequestContext()).Table(tableName).Create(historyData)
	if result.Error() != nil {
		zap.L().Error("创建历史记录失败",
			zap.String("entityType", record.EntityType.String()),
			zap.Uint32("entityID", record.EntityID),
			zap.Error(result.Error()),
		)
		return result.Error()
	}

	return nil
}

// CompareAndCreateHistory 比较新旧对象并创建历史记录
func (s *HistoryService) CompareAndCreateHistory(ctx interface {
	RequestContext() interface{}
}, entityType EntityType, entityID uint32, oldObj, newObj interface{}, operateType OperateType, operator, operatorRoleType string) error {
	// 比较新旧对象的差异
	content := s.compareObjects(oldObj, newObj, operateType)

	// 创建历史记录
	record := &HistoryRecord{
		EntityID:         entityID,
		EntityType:       entityType,
		OperateType:      operateType,
		OperatedAt:       time.Now(),
		Content:          content,
		Operator:         operator,
		OperatorRoleType: operatorRoleType,
	}

	return s.CreateHistory(ctx, record)
}

// compareObjects 比较新旧对象的差异
func (s *HistoryService) compareObjects(oldObj, newObj interface{}, operateType OperateType) map[string]interface{} {
	content := make(map[string]interface{})

	if operateType == OperateTypeCreated {
		// 创建操作：记录所有新字段
		if newObj != nil {
			newVal := reflect.ValueOf(newObj)
			if newVal.Kind() == reflect.Ptr {
				newVal = newVal.Elem()
			}
			if newVal.Kind() == reflect.Struct {
				for i := 0; i < newVal.NumField(); i++ {
					field := newVal.Type().Field(i)
					value := newVal.Field(i)

					// 跳过一些不需要记录的字段
					if s.shouldSkipField(field.Name) {
						continue
					}

					// 处理不同类型的字段
					content[field.Name] = map[string]interface{}{
						"old": "",
						"new": s.formatFieldValue(value),
					}
				}
			}
		}
	} else if operateType == OperateTypeUpdated {
		// 更新操作：比较新旧对象的差异
		if oldObj != nil && newObj != nil {
			oldVal := reflect.ValueOf(oldObj)
			newVal := reflect.ValueOf(newObj)

			if oldVal.Kind() == reflect.Ptr {
				oldVal = oldVal.Elem()
			}
			if newVal.Kind() == reflect.Ptr {
				newVal = newVal.Elem()
			}

			if oldVal.Kind() == reflect.Struct && newVal.Kind() == reflect.Struct {
				for i := 0; i < newVal.NumField(); i++ {
					field := newVal.Type().Field(i)
					oldField := oldVal.Field(i)
					newField := newVal.Field(i)

					// 跳过一些不需要记录的字段
					if s.shouldSkipField(field.Name) {
						continue
					}

					// 比较字段值
					if !reflect.DeepEqual(oldField.Interface(), newField.Interface()) {
						content[field.Name] = map[string]interface{}{
							"old": s.formatFieldValue(oldField),
							"new": s.formatFieldValue(newField),
						}
					}
				}
			}
		}
	}

	return content
}

// formatFieldValue 格式化字段值
func (s *HistoryService) formatFieldValue(value reflect.Value) interface{} {
	switch value.Kind() {
	case reflect.String:
		return value.String()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return value.Int()
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return value.Uint()
	case reflect.Float32, reflect.Float64:
		return value.Float()
	case reflect.Bool:
		return value.Bool()
	case reflect.Slice, reflect.Array:
		// 处理数组字段
		result := make([]interface{}, value.Len())
		for i := 0; i < value.Len(); i++ {
			result[i] = s.formatFieldValue(value.Index(i))
		}
		return result
	case reflect.Map:
		// 处理Map字段
		result := make(map[string]interface{})
		for _, key := range value.MapKeys() {
			result[fmt.Sprintf("%v", key.Interface())] = s.formatFieldValue(value.MapIndex(key))
		}
		return result
	case reflect.Struct:
		// 处理结构体字段
		if value.Type().String() == "time.Time" {
			return value.Interface().(time.Time).Format("2006-01-02 15:04:05")
		}
		return fmt.Sprintf("%v", value.Interface())
	default:
		return fmt.Sprintf("%v", value.Interface())
	}
}

// shouldSkipField 判断是否应该跳过某个字段
func (s *HistoryService) shouldSkipField(fieldName string) bool {
	skipFields := []string{
		"ID", "Id", "id",
		"CreatedAt", "UpdatedAt", "created_at", "updated_at",
		"CreatedUser", "UpdatedUser", "created_user", "updated_user",
		"Version", "version",
	}

	for _, skipField := range skipFields {
		if fieldName == skipField {
			return true
		}
	}
	return false
}

// formatContent 格式化内容为JSON字符串
func (s *HistoryService) formatContent(content map[string]interface{}) string {
	if len(content) == 0 {
		return "{}"
	}

	jsonBytes, err := json.Marshal(content)
	if err != nil {
		zap.L().Error("序列化历史记录内容失败", zap.Error(err))
		return "{}"
	}

	return string(jsonBytes)
}
