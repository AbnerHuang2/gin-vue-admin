package emag

import (
	"time"
)

// 品类状态常量
const (
	CategoryStatusNormal     = "normal"      // 正常状态
	CategoryStatusBadRequest = "bad_request" // 请求失败状态
)

// EmagCategory Emag品类管理结构体
type EmagCategory struct {
	ID              uint       `gorm:"primarykey;autoIncrement" json:"id"`                                                       // 主键ID
	CategoryId      string     `gorm:"column:category_id;type:varchar(64);not null;comment:品类ID" json:"categoryId"`              // 品类ID
	CategoryName    string     `gorm:"column:category_name;type:varchar(128);not null;comment:品类名称" json:"categoryName"`         // 品类名称
	SubcategoryName string     `gorm:"column:subcategory_name;type:varchar(128);comment:子品类名称" json:"subcategoryName"`           // 子品类名称
	Status          string     `gorm:"column:status;type:varchar(32);default:normal;comment:状态" json:"status"`                   // 状态: normal/bad_request
	FailCount       int        `gorm:"column:fail_count;type:int;default:0;comment:连续失败次数" json:"failCount"`                     // 连续失败次数
	LastFailReason  string     `gorm:"column:last_fail_reason;type:varchar(255);comment:最后失败原因" json:"lastFailReason"`           // 最后失败原因
	LastFailAt      *time.Time `gorm:"column:last_fail_at;type:datetime;comment:最后失败时间" json:"lastFailAt"`                       // 最后失败时间
	CreatedAt       *time.Time `gorm:"column:created_at;type:timestamp;default:CURRENT_TIMESTAMP;comment:创建时间" json:"createdAt"` // 创建时间
	UpdatedAt       *time.Time `gorm:"column:updated_at;type:timestamp;default:CURRENT_TIMESTAMP;comment:更新时间" json:"updatedAt"` // 更新时间
}

// TableName 指定表名
func (EmagCategory) TableName() string {
	return "emag_category"
}
