package emag

import (
	"time"
)

// EmagCategoryStat Emag品类指标结构体
type EmagCategoryStat struct {
	ID                uint       `gorm:"primarykey;autoIncrement" json:"id"`                                                     // 主键ID
	CategoryId        string     `gorm:"column:category_id;type:varchar(64);not null;comment:品类ID" json:"categoryId"`            // 品类ID
	Total             int        `gorm:"column:total;type:int;comment:总数" json:"total"`                                          // 总数
	SupperHotTotal    int        `gorm:"column:supper_hot_total;type:int;comment:超热销总数" json:"supperHotTotal"`                   // 超热销总数
	OemTotal          int        `gorm:"column:oem_total;type:int;comment:OEM总数" json:"oemTotal"`                                // OEM总数
	OemSupperHotTotal int        `gorm:"column:oem_supper_hot_total;type:int;comment:OEM超热销总数" json:"oemSupperHotTotal"`         // OEM超热销总数
	SupperHotRate     float64    `gorm:"column:supper_hot_rate;type:float;comment:超热销率" json:"supperHotRate"`                    // 超热销率
	OemSupperHotRate  float64    `gorm:"column:oem_supper_hot_rate;type:float;comment:OEM超热销率" json:"oemSupperHotRate"`          // OEM超热销率
	SnapshotDate      *time.Time `gorm:"column:snapshot_date;type:date;not null;comment:快照日期" json:"snapshotDate"`               // 快照日期
	Tags              string     `gorm:"column:tags;type:varchar(30);comment:标签" json:"tags"`                                    // 标签
	CreatedAt         *time.Time `gorm:"column:create_at;type:datetime;default:CURRENT_TIMESTAMP;comment:创建时间" json:"createdAt"` // 创建时间
	UpdatedAt         *time.Time `gorm:"column:update_at;type:datetime;default:CURRENT_TIMESTAMP;comment:更新时间" json:"updatedAt"` // 更新时间
}

// TableName 指定表名
func (EmagCategoryStat) TableName() string {
	return "emag_category_stat"
}

// EmagCategoryStatWithName 品类指标带品类名称结构体（用于关联查询返回）
type EmagCategoryStatWithName struct {
	EmagCategoryStat
	CategoryName    string `json:"categoryName"`    // 品类名称（来自 emag_category 表）
	SubcategoryName string `json:"subcategoryName"` // 子品类名称（来自 emag_category 表）
}

// CategoryStatGrowth 品类指标同比增长结构体
type CategoryStatGrowth struct {
	CategoryId               string  `json:"categoryId"`               // 品类ID
	CategoryName             string  `json:"categoryName"`             // 品类名称
	SubcategoryName          string  `json:"subcategoryName"`          // 子品类名称
	CurrentTotal             int     `json:"currentTotal"`             // 当前总数
	PreviousTotal            int     `json:"previousTotal"`            // 上期总数
	TotalGrowthRate          float64 `json:"totalGrowthRate"`          // 总数增长率
	CurrentSupperHotTotal    int     `json:"currentSupperHotTotal"`    // 当前超热销总数
	PreviousSupperHotTotal   int     `json:"previousSupperHotTotal"`   // 上期超热销总数
	SupperHotGrowthRate      float64 `json:"supperHotGrowthRate"`      // 超热销增长率
	CurrentSupperHotRate     float64 `json:"currentSupperHotRate"`     // 当前超热销率
	PreviousSupperHotRate    float64 `json:"previousSupperHotRate"`    // 上期超热销率
	CurrentOemSupperHotRate  float64 `json:"currentOemSupperHotRate"`  // 当前OEM超热销率
	PreviousOemSupperHotRate float64 `json:"previousOemSupperHotRate"` // 上期OEM超热销率
	CurrentSnapshotDate      string  `json:"currentSnapshotDate"`      // 当前快照日期
	PreviousSnapshotDate     string  `json:"previousSnapshotDate"`     // 上期快照日期
}
