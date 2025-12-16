package request

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
)

// EmagCategorySearch 品类搜索请求参数
type EmagCategorySearch struct {
	request.PageInfo
	ID           uint   `json:"id" form:"id"`                     // 主键ID
	CategoryId   string `json:"categoryId" form:"categoryId"`     // 品类ID
	CategoryName string `json:"categoryName" form:"categoryName"` // 品类名称
}
