package emag

import (
	"time"
)

// EmagOrder Emag订单结构体
type EmagOrder struct {
	ID             uint       `gorm:"primarykey;autoIncrement" json:"id"`                                                      // 主键ID
	OrderID        string     `gorm:"column:order_id;type:varchar(32);uniqueIndex;not null;comment:订单ID" json:"orderId"`       // 订单ID (对应mktId)
	OrderDateLocal string     `gorm:"column:order_date_local;type:varchar(32);comment:订单日期" json:"orderDateLocal"`             // 订单日期 (对应date)
	Country        string     `gorm:"column:country;type:varchar(20);comment:国家" json:"country"`                               // 国家 (对应customerBillingCountry)
	Currency       string     `gorm:"column:currency;type:varchar(20);comment:货币" json:"currency"`                             // 货币 (对应currency)
	Price          float64    `gorm:"column:price;type:decimal(10,2);not null;comment:金额" json:"price"`                        // 金额 (对应totalValue)
	PriceCny       float64    `gorm:"column:price_cny;type:decimal(10,2);not null;comment:金额(人民币)" json:"priceCny"`            // 金额(人民币)
	Status         string     `gorm:"column:status;type:varchar(32);comment:订单状态" json:"status"`                               // 订单状态
	CustomerName   string     `gorm:"column:customer_name;type:varchar(128);comment:客户名称" json:"customerName"`                 // 客户名称
	CreatedAt      *time.Time `gorm:"column:created_at;type:datetime;default:CURRENT_TIMESTAMP;comment:创建时间" json:"createdAt"` // 创建时间
	UpdatedAt      *time.Time `gorm:"column:updated_at;type:datetime;default:CURRENT_TIMESTAMP;comment:更新时间" json:"updatedAt"` // 更新时间
}

// TableName 指定表名
func (EmagOrder) TableName() string {
	return "emag_order"
}

// EmagOrderProduct Emag订单产品结构体
type EmagOrderProduct struct {
	ID           uint       `gorm:"primarykey;autoIncrement" json:"id"`                                                      // 主键ID
	OrderID      string     `gorm:"column:order_id;type:varchar(32);index;not null;comment:订单ID" json:"orderId"`             // 订单ID (对应mktId)
	ProductID    string     `gorm:"column:product_id;type:varchar(32);comment:产品ID" json:"productId"`                        // 产品ID (对应partNumberKey)
	ProductUrl   string     `gorm:"column:product_url;type:varchar(512);comment:产品链接" json:"productUrl"`                     // 产品链接 (对应actions里面的productUrl)
	ProductName  string     `gorm:"column:product_name;type:varchar(512);comment:产品名称" json:"productName"`                   // 产品名称
	SalePrice    float64    `gorm:"column:sale_price;type:decimal(10,2);comment:销售价格" json:"salePrice"`                      // 销售价格 (对应salePrice)
	SalePriceCny float64    `gorm:"column:sale_price_cny;type:decimal(10,2);comment:销售价格(人民币)" json:"salePriceCny"`          // 销售价格(人民币)
	Quantity     float64    `gorm:"column:quantity;type:float;comment:数量" json:"quantity"`                                   // 数量 (对应quantity)
	VatRate      float64    `gorm:"column:vat_rate;type:float;comment:增值税率" json:"vatRate"`                                  // 增值税率 (对应vatRate)
	CreatedAt    *time.Time `gorm:"column:created_at;type:datetime;default:CURRENT_TIMESTAMP;comment:创建时间" json:"createdAt"` // 创建时间
	UpdatedAt    *time.Time `gorm:"column:updated_at;type:datetime;default:CURRENT_TIMESTAMP;comment:更新时间" json:"updatedAt"` // 更新时间
}

// TableName 指定表名
func (EmagOrderProduct) TableName() string {
	return "emag_order_product"
}

// EmagOrderWithProducts 订单带产品列表结构体
type EmagOrderWithProducts struct {
	EmagOrder
	Products []EmagOrderProduct `json:"products"` // 产品列表
}

// EmagOrderProductFlat 订单产品平铺结构体（用于前端展示）
type EmagOrderProductFlat struct {
	// 订单信息
	ID             uint    `json:"id"`
	OrderID        string  `json:"orderId"`
	OrderDateLocal string  `json:"orderDateLocal"`
	Country        string  `json:"country"`
	Currency       string  `json:"currency"`
	Price          float64 `json:"price"`
	PriceCny       float64 `json:"priceCny"`
	Status         string  `json:"status"`
	// 产品信息
	ProductID    string  `json:"productId"`
	ProductUrl   string  `json:"productUrl"`
	ProductName  string  `json:"productName"`
	SalePrice    float64 `json:"salePrice"`
	SalePriceCny float64 `json:"salePriceCny"`
	Quantity     float64 `json:"quantity"`
	VatRate      float64 `json:"vatRate"`
}
