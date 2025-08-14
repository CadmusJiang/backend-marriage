package customer_authorization_record_history

import (
	"net/http"
	"time"

	"github.com/xinliangnote/go-gin-api/internal/code"
	"github.com/xinliangnote/go-gin-api/internal/pkg/core"
)

type detailResponse struct {
	Data    historyData `json:"data"`    // 数据
	Success bool        `json:"success"` // 是否成功
}

func (h *handler) GetCustomerAuthorizationRecordHistoryDetail() core.HandlerFunc {
	return func(ctx core.Context) {
		historyId := ctx.Param("historyId")
		if historyId == "" {
			ctx.AbortWithError(core.Error(
				http.StatusBadRequest,
				code.ParamBindError,
				"historyId is required"))
			return
		}

		// Mock数据 - 实际项目中应该从数据库查询
		mockData := map[string]historyData{
			"hist_car_001": {
				Key:                           1,
				HistoryId:                     "hist_car_001",
				CustomerAuthorizationRecordId: "1",
				OperateType:                   "created",
				OccurredAt:                    1705123200,
				Content: map[string]interface{}{
					"name": map[string]interface{}{
						"old": "",
						"new": "用户1",
					},
					"birthYear": map[string]interface{}{
						"old": "",
						"new": 1985,
					},
					"gender": map[string]interface{}{
						"old": "",
						"new": "male",
					},
					"height": map[string]interface{}{
						"old": "",
						"new": 175,
					},
					"city": map[string]interface{}{
						"old": "",
						"new": "北京",
					},
					"authStore": map[string]interface{}{
						"old": "",
						"new": "朝阳门店***",
					},
					"education": map[string]interface{}{
						"old": "",
						"new": "本科",
					},
					"profession": map[string]interface{}{
						"old": "",
						"new": "工程师",
					},
					"income": map[string]interface{}{
						"old": "",
						"new": "50w",
					},
					"phone": map[string]interface{}{
						"old": "",
						"new": "138****1234",
					},
					"wechat": map[string]interface{}{
						"old": "",
						"new": "wx_12345****",
					},
					"drainageAccount": map[string]interface{}{
						"old": "",
						"new": "drainage_001",
					},
					"drainageId": map[string]interface{}{
						"old": "",
						"new": "D12345",
					},
					"drainageChannel": map[string]interface{}{
						"old": "",
						"new": "小红书",
					},
					"remark": map[string]interface{}{
						"old": "",
						"new": "备注信息1",
					},
					"isAuthorized": map[string]interface{}{
						"old": "",
						"new": true,
					},
					"authPhotos": map[string]interface{}{
						"old": "",
						"new": []string{
							"https://picsum.photos/300/200?random=0",
							"https://picsum.photos/300/200?random=1",
							"https://picsum.photos/300/200?random=2",
						},
					},
					"isProfileComplete": map[string]interface{}{
						"old": "",
						"new": true,
					},
					"isAssigned": map[string]interface{}{
						"old": "",
						"new": true,
					},
					"isPaid": map[string]interface{}{
						"old": "",
						"new": true,
					},
					"paymentAmount": map[string]interface{}{
						"old": "",
						"new": 25000.00,
					},
					"refundAmount": map[string]interface{}{
						"old": "",
						"new": 0.00,
					},
					"group": map[string]interface{}{
						"old": "",
						"new": "南京-天元大厦组",
					},
					"team": map[string]interface{}{
						"old": "",
						"new": "营销团队A",
					},
					"account": map[string]interface{}{
						"old": "",
						"new": "张伟",
					},
				},
				OperatorUsername: "admin",
				OperatorNickname: "系统管理员",
				OperatorRoleType: "company_manager",
				CreatedAt:        time.Unix(1705123200, 0).Format("2006-01-02T15:04:05Z"),
				UpdatedAt:        time.Unix(1705123200, 0).Format("2006-01-02T15:04:05Z"),
			},
			"hist_car_002": {
				Key:                           2,
				HistoryId:                     "hist_car_002",
				CustomerAuthorizationRecordId: "1",
				OperateType:                   "modified",
				OccurredAt:                    1705209600,
				Content: map[string]interface{}{
					"isAuthorized": map[string]interface{}{
						"old": false,
						"new": true,
					},
					"authPhotos": map[string]interface{}{
						"old": []string{},
						"new": []string{
							"https://picsum.photos/300/200?random=0",
							"https://picsum.photos/300/200?random=1",
							"https://picsum.photos/300/200?random=2",
						},
					},
				},
				OperatorUsername: "zhangwei",
				OperatorNickname: "张伟",
				OperatorRoleType: "company_manager",
				CreatedAt:        time.Unix(1705209600, 0).Format("2006-01-02T15:04:05Z"),
				UpdatedAt:        time.Unix(1705209600, 0).Format("2006-01-02T15:04:05Z"),
			},
		}

		data, exists := mockData[historyId]
		if !exists {
			ctx.AbortWithError(core.Error(
				http.StatusNotFound,
				code.ServerError,
				"history record not found"))
			return
		}

		ctx.Payload(detailResponse{
			Data:    data,
			Success: true,
		})
	}
}
