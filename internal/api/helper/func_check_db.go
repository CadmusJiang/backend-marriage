package helper

import (
	"github.com/xinliangnote/go-gin-api/internal/pkg/core"
)

// CheckDatabase 检查数据库表结构
func (h *handler) CheckDatabase() core.HandlerFunc {
	return func(c core.Context) {
		// 检查数据库连接
		if h.db == nil {
			c.AbortWithError(core.Error(500, 10101, "数据库连接未初始化"))
			return
		}

		db := h.db.GetDbR()
		if db == nil {
			c.AbortWithError(core.Error(500, 10101, "数据库读连接未初始化"))
			return
		}

		// 检查表是否存在
		tables := []string{"account", "account_history", "org"}
		result := make(map[string]interface{})

		for _, table := range tables {
			// 检查表是否存在
			var tableExists bool
			err := db.Raw("SELECT COUNT(*) > 0 FROM information_schema.tables WHERE table_schema = DATABASE() AND table_name = ?", table).Scan(&tableExists).Error
			if err != nil {
				result[table] = map[string]interface{}{
					"exists": false,
					"error":  err.Error(),
				}
				continue
			}

			if !tableExists {
				result[table] = map[string]interface{}{
					"exists": false,
					"error":  "表不存在",
				}
				continue
			}

			// 获取表结构
			var columns []map[string]interface{}
			err = db.Raw("DESCRIBE " + table).Scan(&columns).Error
			if err != nil {
				result[table] = map[string]interface{}{
					"exists": true,
					"error":  err.Error(),
				}
				continue
			}

			// 获取记录数
			var count int64
			err = db.Raw("SELECT COUNT(*) FROM " + table).Scan(&count).Error
			if err != nil {
				result[table] = map[string]interface{}{
					"exists":      true,
					"columns":     columns,
					"count":       count,
					"count_error": err.Error(),
				}
				continue
			}

			result[table] = map[string]interface{}{
				"exists":  true,
				"columns": columns,
				"count":   count,
			}
		}

		c.Payload(result)
	}
}
