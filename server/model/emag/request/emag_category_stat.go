package request

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
)

// CategoryStatTop20Search 品类指标Top20查询请求参数
type CategoryStatTop20Search struct {
	request.PageInfo
	SupperHotRate    float64 `json:"supperHotRate" form:"supperHotRate"`       // 超热销率 (大于等于)
	OemSupperHotRate float64 `json:"oemSupperHotRate" form:"oemSupperHotRate"` // OEM超热销率 (大于等于)
	SnapshotDate     string  `json:"snapshotDate" form:"snapshotDate"`         // 快照日期 (格式: 2006-01-02)
}

// CategoryStatGrowthSearch 品类指标同比增长查询请求参数
type CategoryStatGrowthSearch struct {
	request.PageInfo
	Limit int `json:"limit" form:"limit"` // 返回数量限制
}
