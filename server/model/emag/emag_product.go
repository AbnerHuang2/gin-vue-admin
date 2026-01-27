package emag

import (
	"time"
)

// EmagProduct Emag产品结构体
type EmagProduct struct {
	ID            uint       `gorm:"primarykey;autoIncrement" json:"id"`                                                      // 主键ID
	ProductID     string     `gorm:"column:product_id;type:varchar(64);uniqueIndex;not null;comment:产品ID" json:"productId"`   // 产品ID (对应id)
	CategoryID    int        `gorm:"column:category_id;type:int;comment:类目ID" json:"categoryId"`                              // 类目ID (对应categoryId)
	ExtID         string     `gorm:"column:ext_id;type:varchar(64);comment:外部ID" json:"extId"`                                // 外部ID (对应extId)
	PN            string     `gorm:"column:pn;type:varchar(64);comment:产品编号" json:"pn"`                                       // 产品编号 (对应extPartNumber)
	PNK           string     `gorm:"column:pnk;type:varchar(64);comment:产品编号Key" json:"pnk"`                                  // 产品编号Key (对应docProductPartNumberKey)
	EAN           string     `gorm:"column:ean;type:varchar(64);comment:EAN码" json:"ean"`                                     // EAN码 (对应eans数组第一个)
	Title         string     `gorm:"column:title;type:varchar(256);comment:产品标题" json:"title"`                                // 产品标题 (对应docProductName)
	Status        string     `gorm:"column:status;type:varchar(32);comment:产品状态" json:"status"`                               // 产品状态 (对应statusName)
	SalePrice     float64    `gorm:"column:sale_price;type:decimal(10,2);comment:销售价格" json:"salePrice"`                      // 销售价格 (对应extSalePrice)
	AfterTaxPrice float64    `gorm:"column:after_tax_price;type:decimal(10,2);comment:含税价格" json:"afterTaxPrice"`             // 含税价格 (extSalePrice * (1 + vat))
	Currency      string     `gorm:"column:currency;type:varchar(32);comment:货币" json:"currency"`                             // 货币 (对应currency)
	Country       string     `gorm:"column:country;type:varchar(20);comment:国家" json:"country"`                               // 国家 (根据currency转换)
	SalePriceCN   float64    `gorm:"column:sale_price_cn;type:decimal(10,2);comment:销售价格(人民币)" json:"salePriceCn"`            // 销售价格(人民币)
	CostPriceCN   float64    `gorm:"column:cost_price_cn;type:decimal(10,2);comment:成本价格(人民币)" json:"costPriceCn"`            // 成本价格(人民币) - 后续手动补充
	VAT           string     `gorm:"column:vat;type:varchar(64);comment:增值税率" json:"vat"`                                     // 增值税率 (对应vatValues的key)
	Stock         int        `gorm:"column:stock;type:int;comment:库存" json:"stock"`                                           // 库存 (对应extStock)
	URL           string     `gorm:"column:url;type:varchar(256);comment:产品链接" json:"url"`                                    // 产品链接 (对应links.details)
	BuyButtonRank int        `gorm:"column:buy_button_rank;type:tinyint;comment:购买按钮排名" json:"buyButtonRank"`                 // 购买按钮排名 (对应productPerformanceBuyButtonRank)
	BuyButtonCnt  int        `gorm:"column:buy_button_cnt;type:tinyint;comment:购买按钮数量" json:"buyButtonCnt"`                   // 购买按钮数量 (对应productPerformanceMultiofferOffersCount)
	CreatedAt     *time.Time `gorm:"column:created_at;type:datetime;default:CURRENT_TIMESTAMP;comment:创建时间" json:"createdAt"` // 创建时间
	UpdatedAt     *time.Time `gorm:"column:updated_at;type:datetime;default:CURRENT_TIMESTAMP;comment:更新时间" json:"updatedAt"` // 更新时间
}

// TableName 指定表名
func (EmagProduct) TableName() string {
	return "emag_product"
}
