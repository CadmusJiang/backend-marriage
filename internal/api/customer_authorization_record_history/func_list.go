package customer_authorization_record_history

import (
	"net/http"
	"time"

	"github.com/xinliangnote/go-gin-api/internal/code"
	"github.com/xinliangnote/go-gin-api/internal/pkg/core"
)

type listRequest struct {
	Current                       int    `form:"current"`                       // 当前页码
	PageSize                      int    `form:"pageSize"`                      // 每页数量
	CustomerAuthorizationRecordId string `form:"customerAuthorizationRecordId"` // 客户授权记录ID
	OperateType                   string `form:"operateType"`                   // 操作类型
	OperatorUsername              string `form:"operatorUsername"`              // 操作人用户名
	OperatorRoleType              string `form:"operatorRoleType"`              // 操作人角色类型
}

type historyData struct {
	Key                           int                    `json:"key"`                           // 主键
	HistoryId                     string                 `json:"historyId"`                     // 历史记录ID
	CustomerAuthorizationRecordId string                 `json:"customerAuthorizationRecordId"` // 客户授权记录ID
	OperateType                   string                 `json:"operateType"`                   // 操作类型
	OccurredAt                    int64                  `json:"occurredAt"`                    // 操作时间戳
	Content                       map[string]interface{} `json:"content"`                       // 操作内容
	OperatorUsername              string                 `json:"operatorUsername"`              // 操作人用户名
	OperatorNickname              string                 `json:"operatorNickname"`              // 操作人姓名
	OperatorRoleType              string                 `json:"operatorRoleType"`              // 操作人角色类型
	CreatedAt                     string                 `json:"createdAt"`                     // 创建时间
	UpdatedAt                     string                 `json:"updatedAt"`                     // 修改时间
}

type listResponse struct {
	Data     []historyData `json:"data"`     // 数据列表
	Total    int           `json:"total"`    // 总数
	Success  bool          `json:"success"`  // 是否成功
	PageSize int           `json:"pageSize"` // 每页数量
	Current  int           `json:"current"`  // 当前页码
}

func (h *handler) GetCustomerAuthorizationRecordHistoryList() core.HandlerFunc {
	return func(ctx core.Context) {
		req := new(listRequest)
		if err := ctx.ShouldBindQuery(req); err != nil {
			ctx.AbortWithError(core.Error(
				http.StatusBadRequest,
				code.ParamBindError,
				code.Text(code.ParamBindError)).WithError(err),
			)
			return
		}

		// 设置默认值
		if req.Current <= 0 {
			req.Current = 1
		}
		if req.PageSize <= 0 {
			req.PageSize = 10
		}

		// Mock数据 - 实际项目中应该从数据库查询
		mockData := []historyData{
			{
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
			{
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
			{
				Key:                           3,
				HistoryId:                     "hist_car_003",
				CustomerAuthorizationRecordId: "1",
				OperateType:                   "modified",
				OccurredAt:                    1705296000,
				Content: map[string]interface{}{
					"isAssigned": map[string]interface{}{
						"old": false,
						"new": true,
					},
					"team": map[string]interface{}{
						"old": "",
						"new": "营销团队A",
					},
				},
				OperatorUsername: "liming",
				OperatorNickname: "李明",
				OperatorRoleType: "team_manager",
				CreatedAt:        time.Unix(1705296000, 0).Format("2006-01-02T15:04:05Z"),
				UpdatedAt:        time.Unix(1705296000, 0).Format("2006-01-02T15:04:05Z"),
			},
			{
				Key:                           4,
				HistoryId:                     "hist_car_004",
				CustomerAuthorizationRecordId: "1",
				OperateType:                   "modified",
				OccurredAt:                    1705382400,
				Content: map[string]interface{}{
					"isPaid": map[string]interface{}{
						"old": false,
						"new": true,
					},
					"paymentAmount": map[string]interface{}{
						"old": 0.00,
						"new": 25000.00,
					},
				},
				OperatorUsername: "wangfang",
				OperatorNickname: "王芳",
				OperatorRoleType: "team_manager",
				CreatedAt:        time.Unix(1705382400, 0).Format("2006-01-02T15:04:05Z"),
				UpdatedAt:        time.Unix(1705382400, 0).Format("2006-01-02T15:04:05Z"),
			},
			{
				Key:                           5,
				HistoryId:                     "hist_car_005",
				CustomerAuthorizationRecordId: "2",
				OperateType:                   "created",
				OccurredAt:                    1705209600,
				Content: map[string]interface{}{
					"name": map[string]interface{}{
						"old": "",
						"new": "用户2",
					},
					"birthYear": map[string]interface{}{
						"old": "",
						"new": 1990,
					},
					"gender": map[string]interface{}{
						"old": "",
						"new": "female",
					},
					"height": map[string]interface{}{
						"old": "",
						"new": 165,
					},
					"city": map[string]interface{}{
						"old": "",
						"new": "上海",
					},
					"authStore": map[string]interface{}{
						"old": "",
						"new": "浦东门店***",
					},
					"education": map[string]interface{}{
						"old": "",
						"new": "硕士",
					},
					"profession": map[string]interface{}{
						"old": "",
						"new": "设计师",
					},
					"income": map[string]interface{}{
						"old": "",
						"new": "80w",
					},
					"phone": map[string]interface{}{
						"old": "",
						"new": "139****5678",
					},
					"wechat": map[string]interface{}{
						"old": "",
						"new": "wx_67890****",
					},
					"drainageAccount": map[string]interface{}{
						"old": "",
						"new": "drainage_002",
					},
					"drainageId": map[string]interface{}{
						"old": "",
						"new": "D67890",
					},
					"drainageChannel": map[string]interface{}{
						"old": "",
						"new": "小红书",
					},
					"remark": map[string]interface{}{
						"old": "",
						"new": "",
					},
					"isAuthorized": map[string]interface{}{
						"old": "",
						"new": true,
					},
					"authPhotos": map[string]interface{}{
						"old": "",
						"new": []string{
							"https://picsum.photos/300/200?random=3",
							"https://picsum.photos/300/200?random=4",
							"https://picsum.photos/300/200?random=5",
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
						"new": false,
					},
					"paymentAmount": map[string]interface{}{
						"old": "",
						"new": 0.00,
					},
					"refundAmount": map[string]interface{}{
						"old": "",
						"new": 0.00,
					},
					"group": map[string]interface{}{
						"old": "",
						"new": "南京-南京南站组",
					},
					"team": map[string]interface{}{
						"old": "",
						"new": "营销团队B",
					},
					"account": map[string]interface{}{
						"old": "",
						"new": "刘强",
					},
				},
				OperatorUsername: "admin",
				OperatorNickname: "系统管理员",
				OperatorRoleType: "company_manager",
				CreatedAt:        time.Unix(1705209600, 0).Format("2006-01-02T15:04:05Z"),
				UpdatedAt:        time.Unix(1705209600, 0).Format("2006-01-02T15:04:05Z"),
			},
		}

		// 过滤数据
		filteredData := make([]historyData, 0)
		for _, item := range mockData {
			// 按客户授权记录ID过滤
			if req.CustomerAuthorizationRecordId != "" && item.CustomerAuthorizationRecordId != req.CustomerAuthorizationRecordId {
				continue
			}
			// 按操作类型过滤
			if req.OperateType != "" && item.OperateType != req.OperateType {
				continue
			}
			// 按操作人用户名过滤
			if req.OperatorUsername != "" && item.OperatorUsername != req.OperatorUsername {
				continue
			}
			// 按操作人角色类型过滤
			if req.OperatorRoleType != "" && item.OperatorRoleType != req.OperatorRoleType {
				continue
			}
			filteredData = append(filteredData, item)
		}

		// 计算分页
		total := len(filteredData)
		start := (req.Current - 1) * req.PageSize
		end := start + req.PageSize
		if end > total {
			end = total
		}
		if start > total {
			start = total
		}

		var pagedData []historyData
		if start < total {
			pagedData = filteredData[start:end]
		}

		ctx.Payload(listResponse{
			Data:     pagedData,
			Total:    total,
			Success:  true,
			PageSize: req.PageSize,
			Current:  req.Current,
		})
	}
}
