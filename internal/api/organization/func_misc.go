package organization

import (
	"github.com/xinliangnote/go-gin-api/internal/pkg/core"
)

// ListUnassignedAccounts 获取未分配账户列表
func (h *handler) ListUnassignedAccounts() core.HandlerFunc {
	return func(c core.Context) {
		// TODO: 实现获取未分配账户列表逻辑
		c.Payload(map[string]interface{}{
			"data": []interface{}{},
			"meta": map[string]interface{}{
				"total":    0,
				"pageSize": 10,
				"current":  1,
			},
		})
	}
}
