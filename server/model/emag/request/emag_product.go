package request

import "github.com/flipped-aurora/gin-vue-admin/server/model/common/request"

// EmagProductSearch 产品搜索条件
type EmagProductSearch struct {
	request.PageInfo
	PNK      string `json:"pnk" form:"pnk"`           // 产品编号Key
	ExtID    string `json:"extId" form:"extId"`       // 外部ID
	StockMin int    `json:"stockMin" form:"stockMin"` // 最小库存
	Country  string `json:"country" form:"country"`   // 国家
	PN       string `json:"pn" form:"pn"`             // 产品编号（保留向后兼容）
	Title    string `json:"title" form:"title"`       // 产品标题（保留向后兼容）
	Status   string `json:"status" form:"status"`     // 产品状态（保留向后兼容）
}

// EmagProductSyncRequest 同步产品请求
type EmagProductSyncRequest struct {
	Page  int `json:"page"`  // 页码
	Limit int `json:"limit"` // 每页数量
}
