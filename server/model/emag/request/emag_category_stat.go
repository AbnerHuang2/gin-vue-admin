package request

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
)

// CategoryStatSearch 品类指标查询请求参数
type CategoryStatSearch struct {
	request.PageInfo
	SupperHotRate    float64 `json:"supperHotRate" form:"supperHotRate"`       // 超热销率 (大于等于)
	OemSupperHotRate float64 `json:"oemSupperHotRate" form:"oemSupperHotRate"` // OEM超热销率 (大于等于)
	SnapshotDate     string  `json:"snapshotDate" form:"snapshotDate"`         // 快照日期 (格式: 2006-01-02)
}

// CategoryStatGrowthSearch 品类指标同比增长查询请求参数
type CategoryStatGrowthSearch struct {
	request.PageInfo
}
