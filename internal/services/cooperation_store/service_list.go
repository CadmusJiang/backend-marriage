package cooperation_store

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
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

	// 构建基础查询
	baseQuery := s.db.GetDbR().WithContext(ctx.RequestContext()).Model(&model.CooperationStore{})

	// 应用基本筛选条件
	if req.StoreName != "" {
		baseQuery = baseQuery.Where("store_name LIKE ?", "%"+req.StoreName+"%")
	}
	if req.CooperationCity != "" {
		baseQuery = baseQuery.Where("cooperation_city_code = ?", req.CooperationCity)
	}
	if len(req.CooperationStatus) > 0 {
		baseQuery = baseQuery.Where("cooperation_status IN (?)", req.CooperationStatus)
	}
	if req.CompanyName != "" {
		baseQuery = baseQuery.Where("company_name LIKE ?", "%"+req.CompanyName+"%")
	}
	if req.StoreShortName != "" {
		baseQuery = baseQuery.Where("store_short_name LIKE ?", "%"+req.StoreShortName+"%")
	}
	if req.ActualBusinessAddress != "" {
		baseQuery = baseQuery.Where("actual_business_address LIKE ?", "%"+req.ActualBusinessAddress+"%")
	}

	// 简化JSON查询，避免复杂的OR逻辑
	if len(req.CooperationType) > 0 {
		for _, t := range req.CooperationType {
			if t != "" {
				baseQuery = baseQuery.Where("JSON_CONTAINS(cooperation_type, JSON_QUOTE(?))", t)
			}
		}
	}
	if len(req.CooperationMethod) > 0 {
		for _, m := range req.CooperationMethod {
			if m != "" {
				baseQuery = baseQuery.Where("JSON_CONTAINS(cooperation_method, JSON_QUOTE(?))", m)
			}
		}
	}

	// 执行Count查询
	var count int64
	countSQL := "SELECT COUNT(*) FROM cooperation_store"
	countArgs := []interface{}{}

	// 构建WHERE条件
	whereConditions := []string{}

	if req.StoreName != "" {
		whereConditions = append(whereConditions, "store_name LIKE ?")
		countArgs = append(countArgs, "%"+req.StoreName+"%")
	}
	if req.CooperationCity != "" {
		whereConditions = append(whereConditions, "cooperation_city_code = ?")
		countArgs = append(countArgs, req.CooperationCity)
	}
	if len(req.CooperationStatus) > 0 {
		placeholders := make([]string, len(req.CooperationStatus))
		for i := range req.CooperationStatus {
			placeholders[i] = "?"
		}
		whereConditions = append(whereConditions, "cooperation_status IN ("+strings.Join(placeholders, ",")+")")
		for _, status := range req.CooperationStatus {
			countArgs = append(countArgs, status)
		}
	}
	if req.CompanyName != "" {
		whereConditions = append(whereConditions, "company_name LIKE ?")
		countArgs = append(countArgs, "%"+req.CompanyName+"%")
	}
	if req.StoreShortName != "" {
		whereConditions = append(whereConditions, "store_short_name LIKE ?")
		countArgs = append(countArgs, "%"+req.StoreShortName+"%")
	}
	if req.ActualBusinessAddress != "" {
		whereConditions = append(whereConditions, "actual_business_address LIKE ?")
		countArgs = append(countArgs, "%"+req.ActualBusinessAddress+"%")
	}

	// 添加JSON查询条件
	if len(req.CooperationType) > 0 {
		for _, t := range req.CooperationType {
			if t != "" {
				whereConditions = append(whereConditions, "JSON_CONTAINS(cooperation_type, JSON_QUOTE(?))")
				countArgs = append(countArgs, t)
			}
		}
	}
	if len(req.CooperationMethod) > 0 {
		for _, m := range req.CooperationMethod {
			if m != "" {
				whereConditions = append(whereConditions, "JSON_CONTAINS(cooperation_method, JSON_QUOTE(?))")
				countArgs = append(countArgs, m)
			}
		}
	}

	// 组合SQL
	if len(whereConditions) > 0 {
		countSQL += " WHERE " + strings.Join(whereConditions, " AND ")
	}

	// 执行Count查询
	if err := s.db.GetDbR().WithContext(ctx.RequestContext()).Raw(countSQL, countArgs...).Scan(&count).Error; err != nil {
		return nil, 0, err
	}
	total = count

	// 调试日志：打印Count查询结果
	fmt.Printf("DEBUG: Count查询成功，总数: %d\n", total)
	fmt.Printf("DEBUG: Count SQL: %s\n", countSQL)
	fmt.Printf("DEBUG: Count参数: %v\n", countArgs)

	// 执行数据查询
	// 构建数据查询SQL - 使用具体的字段名而不是*
	dataSQL := "SELECT id, store_name, cooperation_city_code, cooperation_type, store_short_name, company_name, cooperation_method, cooperation_status, business_license, store_photos, actual_business_address, contract_photos, created_at, updated_at FROM cooperation_store"
	if len(whereConditions) > 0 {
		dataSQL += " WHERE " + strings.Join(whereConditions, " AND ")
	}
	dataSQL += " ORDER BY id DESC LIMIT ? OFFSET ?"

	// 添加分页参数
	dataArgs := append(countArgs, req.PageSize, (req.Current-1)*req.PageSize)

	// 调试日志：打印数据查询SQL
	fmt.Printf("DEBUG: 数据查询SQL: %s\n", dataSQL)
	fmt.Printf("DEBUG: 数据查询参数: %v\n", dataArgs)

	// 使用Scan到map，然后手动转换，避免GORM模型问题
	var rawRows []map[string]interface{}
	if err := s.db.GetDbR().WithContext(ctx.RequestContext()).Raw(dataSQL, dataArgs...).Scan(&rawRows).Error; err != nil {
		return nil, 0, err
	}

	// 调试日志：打印查询到的原始数据行数
	fmt.Printf("DEBUG: 查询到原始数据行数: %d\n", len(rawRows))

	// 手动转换数据
	out := make([]StoreItem, 0, len(rawRows))
	for _, rawRow := range rawRows {
		// 安全地处理JSON字段
		var coopType []string
		if rawRow["cooperation_type"] != nil {
			if jsonBytes, ok := rawRow["cooperation_type"].([]byte); ok {
				if err := json.Unmarshal(jsonBytes, &coopType); err != nil {
					coopType = []string{} // 如果解析失败，使用空切片
				}
			} else {
				coopType = []string{}
			}
		} else {
			coopType = []string{}
		}

		var coopMethod []string
		if rawRow["cooperation_method"] != nil {
			if jsonBytes, ok := rawRow["cooperation_method"].([]byte); ok {
				if err := json.Unmarshal(jsonBytes, &coopMethod); err != nil {
					coopMethod = []string{} // 如果解析失败，使用空切片
				}
			} else {
				coopMethod = []string{}
			}
		} else {
			coopMethod = []string{}
		}

		var storePhotos []string
		if rawRow["store_photos"] != nil {
			if jsonBytes, ok := rawRow["store_photos"].([]byte); ok {
				if err := json.Unmarshal(jsonBytes, &storePhotos); err != nil {
					storePhotos = []string{} // 如果解析失败，使用空切片
				}
			} else {
				storePhotos = []string{}
			}
		} else {
			storePhotos = []string{}
		}

		var contractPhotos []string
		if rawRow["contract_photos"] != nil {
			if jsonBytes, ok := rawRow["contract_photos"].([]byte); ok {
				if err := json.Unmarshal(jsonBytes, &contractPhotos); err != nil {
					contractPhotos = []string{} // 如果解析失败，使用空切片
				}
			} else {
				contractPhotos = []string{}
			}
		} else {
			contractPhotos = []string{}
		}

		// 安全地处理时间字段
		var createdAt, updatedAt string
		if rawRow["created_at"] != nil {
			if t, ok := rawRow["created_at"].(time.Time); ok {
				createdAt = t.Format(time.RFC3339)
			} else {
				createdAt = time.Now().Format(time.RFC3339)
			}
		} else {
			createdAt = time.Now().Format(time.RFC3339)
		}

		if rawRow["updated_at"] != nil {
			if t, ok := rawRow["updated_at"].(time.Time); ok {
				updatedAt = t.Format(time.RFC3339)
			} else {
				updatedAt = time.Now().Format(time.RFC3339)
			}
		} else {
			updatedAt = time.Now().Format(time.RFC3339)
		}

		// 安全地处理ID字段
		var id string
		if rawRow["id"] != nil {
			// 调试：打印ID字段的类型和值
			fmt.Printf("DEBUG: ID字段类型: %T, 值: %v\n", rawRow["id"], rawRow["id"])

			switch v := rawRow["id"].(type) {
			case int32:
				id = strconv.Itoa(int(v))
			case int64:
				id = strconv.FormatInt(v, 10)
			case int:
				id = strconv.Itoa(v)
			case float64:
				id = strconv.Itoa(int(v))
			default:
				fmt.Printf("DEBUG: 未知的ID类型: %T, 值: %v\n", v, v)
				id = "0"
			}
		} else {
			fmt.Printf("DEBUG: ID字段为nil\n")
			id = "0"
		}

		// 调试：打印JSON字段的原始值
		fmt.Printf("DEBUG: cooperation_type原始值: %T, %v\n", rawRow["cooperation_type"], rawRow["cooperation_type"])
		fmt.Printf("DEBUG: cooperation_method原始值: %T, %v\n", rawRow["cooperation_method"], rawRow["cooperation_method"])

		out = append(out, StoreItem{
			Id:                    id,
			StoreName:             getStringValue(rawRow["store_name"]),
			CooperationCity:       getStringValue(rawRow["cooperation_city_code"]),
			CooperationType:       coopType,
			StoreShortName:        getStringPtr(rawRow["store_short_name"]),
			CompanyName:           getStringPtr(rawRow["company_name"]),
			CooperationMethod:     coopMethod,
			CooperationStatus:     getStringValue(rawRow["cooperation_status"]),
			BusinessLicense:       getStringPtr(rawRow["business_license"]),
			StorePhotos:           storePhotos,
			ActualBusinessAddress: getStringPtr(rawRow["actual_business_address"]),
			ContractPhotos:        contractPhotos,
			CreatedAt:             createdAt,
			UpdatedAt:             updatedAt,
		})
	}

	// 调试日志：打印最终转换后的数据行数
	fmt.Printf("DEBUG: 最终转换后的数据行数: %d\n", len(out))

	return out, total, nil
}

// 辅助函数：安全地获取字符串值
func getStringValue(v interface{}) string {
	if v == nil {
		return ""
	}
	if str, ok := v.(string); ok {
		return str
	}
	return ""
}

// 辅助函数：安全地获取字符串指针
func getStringPtr(v interface{}) *string {
	if v == nil {
		return nil
	}
	if str, ok := v.(string); ok {
		return &str
	}
	return nil
}
