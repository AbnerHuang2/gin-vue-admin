package request

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
)

// EmagOrderSearch 订单搜索参数
type EmagOrderSearch struct {
	request.PageInfo
	OrderID   string `json:"orderId" form:"orderId"`     // 订单ID
	Country   string `json:"country" form:"country"`     // 国家
	Status    string `json:"status" form:"status"`       // 订单状态
	StartDate string `json:"startDate" form:"startDate"` // 开始日期
	EndDate   string `json:"endDate" form:"endDate"`     // 结束日期
}

// EmagOrderSyncRequest 订单同步请求参数
type EmagOrderSyncRequest struct {
	Status string `json:"status" form:"status"` // 订单状态筛选 (如 STATUS_FINALIZED)
	Page   int    `json:"page" form:"page"`     // 页码
	Limit  int    `json:"limit" form:"limit"`   // 每页数量
}
