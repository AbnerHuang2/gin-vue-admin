package emag

import (
	"time"
)

// EmagCategory Emag品类管理结构体
type EmagCategory struct {
	ID              uint       `gorm:"primarykey;autoIncrement" json:"id"`                                                       // 主键ID
	CategoryId      string     `gorm:"column:category_id;type:varchar(64);not null;comment:品类ID" json:"categoryId"`              // 品类ID
	CategoryName    string     `gorm:"column:category_name;type:varchar(128);not null;comment:品类名称" json:"categoryName"`         // 品类名称
	SubcategoryName string     `gorm:"column:subcategory_name;type:varchar(128);comment:子品类名称" json:"subcategoryName"`           // 子品类名称
	CreatedAt       *time.Time `gorm:"column:created_at;type:timestamp;default:CURRENT_TIMESTAMP;comment:创建时间" json:"createdAt"` // 创建时间
	UpdatedAt       *time.Time `gorm:"column:updated_at;type:timestamp;default:CURRENT_TIMESTAMP;comment:更新时间" json:"updatedAt"` // 更新时间
}

// TableName 指定表名
func (EmagCategory) TableName() string {
	return "emag_category"
}
