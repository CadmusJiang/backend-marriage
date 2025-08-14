package customer_authorization_record

import (
	"github.com/xinliangnote/go-gin-api/internal/pkg/core"
	"github.com/xinliangnote/go-gin-api/internal/repository/mysql"
	svc "github.com/xinliangnote/go-gin-api/internal/services/customer_authorization_record"
	"go.uber.org/zap"
)

type handler struct {
	logger *zap.Logger
	db     mysql.Repo
	cache  mysql.Repo
	svc    svc.Service
}

func New(logger *zap.Logger, db mysql.Repo, cache mysql.Repo) *handler {
	return &handler{
		logger: logger,
		db:     db,
		cache:  cache,
		svc:    svc.New(db),
	}
}

type checkPhoneRequest struct {
	Phone string `form:"phone"` // 手机号
}

type checkPhoneResponse struct {
	Success bool          `json:"success"`
	Exists  bool          `json:"exists"`
	Data    *customerData `json:"data,omitempty"`
}

// CheckPhoneExistence 检查手机号是否已存在（mock实现）
// @Summary 检查手机号是否已存在
// @Tags Customer
// @Accept application/x-www-form-urlencoded
// @Produce json
// @Param phone query string true "手机号"
// @Success 200 {object} checkPhoneResponse
// @Router /api/v1/customer-authorization-records/check-phone [get]
func (h *handler) CheckPhoneExistence() core.HandlerFunc {
	return func(c core.Context) {
		req := new(checkPhoneRequest)
		if err := c.ShouldBindQuery(req); err != nil || req.Phone == "" {
			// 即使缺失参数，也返回success=false? OpenAPI未定义失败结构，这里保持success=true但exists=false
			c.Payload(&checkPhoneResponse{Success: true, Exists: false})
			return
		}

		exists := false
		var found *customerData

		// 简单mock：匹配一个示例手机号或其脱敏形式
		if req.Phone == "13800000000" || req.Phone == "138****1234" || len(req.Phone) >= 4 && req.Phone[:4] == "138*" {
			exists = true
			found = &customerData{
				Key:       1,
				Name:      "用户1",
				Phone:     stringPtr("138****1234"),
				CreatedAt: "2024-01-13T10:00:00Z",
				UpdatedAt: "2024-01-13T10:00:00Z",
			}
		}

		c.Payload(&checkPhoneResponse{Success: true, Exists: exists, Data: found})
	}
}
