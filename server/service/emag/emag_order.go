package emag

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/emag"
	emagReq "github.com/flipped-aurora/gin-vue-admin/server/model/emag/request"
	"go.uber.org/zap"
	"gorm.io/gorm/clause"
)

type EmagOrderService struct{}

var EmagOrderServiceApp = new(EmagOrderService)

// 固定汇率映射表 (货币 -> 人民币)
var currencyRates = map[string]float64{
	"RON": 1.6,  // 罗马尼亚列伊
	"BGN": 3.7,  // 保加利亚列弗
	"HUF": 0.02, // 匈牙利福林
	"EUR": 7.8,  // 欧元
	"USD": 7.2,  // 美元
}

// GraphQL 请求和响应结构体
type graphQLRequest struct {
	OperationName string      `json:"operationName"`
	Variables     interface{} `json:"variables"`
	Query         string      `json:"query"`
}

type ordersVariables struct {
	Filters ordersFilters `json:"filters"`
}

type ordersFilters struct {
	Pagination ordersPagination `json:"pagination"`
	Sort       []ordersSort     `json:"sort"`
	Status     []string         `json:"status"`
}

type ordersPagination struct {
	Page  int `json:"page"`
	Limit int `json:"limit"`
}

type ordersSort struct {
	Field     string `json:"field"`
	Direction string `json:"direction"`
}

// Emag API 响应结构
type emagOrdersResponse struct {
	Data struct {
		Orders struct {
			Items              []emagOrderItem `json:"items"`
			TotalNumberOfItems int             `json:"totalNumberOfItems"`
		} `json:"orders"`
	} `json:"data"`
}

type emagOrderItem struct {
	ID                     string            `json:"id"`
	MktId                  string            `json:"mktId"`
	Date                   string            `json:"date"`
	Type                   int               `json:"type"`
	Status                 int               `json:"status"`
	TotalValue             string            `json:"totalValue"`
	CustomerName           string            `json:"customerName"`
	CustomerBillingCountry string            `json:"customerBillingCountry"`
	Currency               string            `json:"currency"`
	SubStatus              string            `json:"subStatus"`
	Products               []emagProductItem `json:"products"`
}

type emagProductItem struct {
	ID            int                   `json:"id"`
	Name          string                `json:"name"`
	PartNumberKey string                `json:"partNumberKey"`
	Quantity      string                `json:"quantity"`
	VatRate       string                `json:"vatRate"`
	SalePrice     string                `json:"salePrice"`
	Actions       map[string]emagAction `json:"actions"`
}

type emagAction struct {
	Name   string `json:"name"`
	Type   string `json:"type"`
	Status string `json:"status"`
	Action string `json:"action"`
}

// GraphQL 查询语句
const ordersQuery = `query orders($filters: OrderFilterInput!) {
  orders(filters: $filters) {
    items {
      id
      mktId
      date
      type
      status
      isComplete
      paymentMode
      paymentModeMktId
      deliveryMode
      totalValue
      shipmentDeadline
      hasEmagClub
      lateShipment
      cancellationRequest
      sellerId
      sellerName
      deliveryNumber
      courierExternalOfficeId
      customerName
      customerCode
      customerIsJuridical
      customerPhone
      customerBillingName
      customerBillingCity
      customerBillingPhone
      hasLockedCourier
      confirmLockedCourier
      unlockLockedCourier
      hasAttachment
      hasInvoice
      hasAllowedVendorCourierAccounts
      customerBillingStreet
      customerBillingSuburb
      customerBillingCountry
      customerBillingPostalCode
      customerCompany
      customerShippingContact
      customerShippingPhone
      customerShippingCity
      customerShippingStreet
      customerShippingCountry
      customerShippingPostalCode
      customerShippingIsOffice
      customerShippingSuburb
      deliveryAwbStatusCode
      deliveryStatus
      deliveryKpiSellerPickupDeadline
      deliveryKpiSellerLatePickupAh
      deliveryDmsCourierName
      deliveryReservationId
      subPaymentStatus
      incompleteReason
      lockerName
      suggestedCourierDisplayName
      platformId
      observation
      actions
      currency
      subStatus
      wouldBeSupplier
      products {
        id
        extId
        serialNumbers
        name
        partNumber
        partNumberKey
        campaignId
        campaignName
        quantity
        unitSymbol
        vatRate
        salePrice
        status
        vendorProductId
        actions
        discounts {
          referenceId
          vcrId
          extName
          vatRate
          salePrice
          extCurrency
        }
        type
        productFees {
          id
          name
          quantity
          vatRate
          salePrice
          orderLineId
          unitQuantity
          mktId
          docId
          categoryId
          discounts {
            referenceId
            vcrId
            extName
            vatRate
            salePrice
            extCurrency
          }
        }
      }
      shippingTax {
        id
        value
        valueNoVat
        discounts {
          referenceId
          vcrId
          extName
          vatRate
          salePrice
          extCurrency
        }
      }
      stornoTransactions {
        id
        sourceId
        sourceType
        tax {
          id
          status
          value
        }
        shippingTax {
          id
          value
          status
          discounts {
            referenceId
            vcrId
            extName
            vatRate
            salePrice
            extCurrency
          }
        }
        lines {
          id
          status
          productId
          stornoProductId
          retainedAmount
        }
      }
      isStale
      customerGdpr
      isEligibleForDropoff
      deliveryFlagDropoffLocker
      hasPmt
    }
    totalNumberOfItems
  }
}`

// ConvertToCNY 将指定货币金额转换为人民币
func (e *EmagOrderService) ConvertToCNY(amount float64, currency string) float64 {
	currency = strings.ToUpper(currency)
	if rate, ok := currencyRates[currency]; ok {
		return amount * rate
	}
	// 默认返回原始金额（假设是人民币）
	return amount
}

// SyncOrdersFromEmag 从 Emag 平台同步订单数据
func (e *EmagOrderService) SyncOrdersFromEmag(status string, page, limit int) (int, error) {
	if status == "" {
		status = "STATUS_FINALIZED"
	}
	if page < 0 {
		page = 0
	}
	if limit <= 0 {
		limit = 25
	}

	// 构建请求体
	reqBody := graphQLRequest{
		OperationName: "orders",
		Variables: ordersVariables{
			Filters: ordersFilters{
				Pagination: ordersPagination{
					Page:  page,
					Limit: limit,
				},
				Sort: []ordersSort{
					{Field: "type", Direction: "DESC"},
					{Field: "date", Direction: "DESC"},
				},
				Status: []string{status},
			},
		},
		Query: ordersQuery,
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return 0, fmt.Errorf("序列化请求体失败: %v", err)
	}

	// 创建 HTTP 请求
	req, err := http.NewRequest("POST", "https://marketplace.emag.ro/global-listing", bytes.NewBuffer(jsonBody))
	if err != nil {
		return 0, fmt.Errorf("创建请求失败: %v", err)
	}

	// 设置请求头
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9,en-US;q=0.8,en;q=0.7")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Origin", "https://marketplace.emag.ro")
	req.Header.Set("Pragma", "no-cache")
	req.Header.Set("Referer", "https://marketplace.emag.ro/order/list-xb")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/143.0.0.0 Safari/537.36")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	req.Header.Set("Cookie", global.GVA_CONFIG.Emag.Cookie)

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return 0, fmt.Errorf("发送请求失败: %v", err)
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("读取响应失败: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("请求失败, 状态码: %d, 响应: %s", resp.StatusCode, string(body))
	}

	// 解析响应
	var emagResp emagOrdersResponse
	if err := json.Unmarshal(body, &emagResp); err != nil {
		return 0, fmt.Errorf("解析响应失败: %v, 响应内容: %s", err, string(body))
	}

	// 处理订单数据
	orders := make([]emag.EmagOrder, 0)
	products := make([]emag.EmagOrderProduct, 0)

	for _, item := range emagResp.Data.Orders.Items {
		// 解析金额
		price, _ := strconv.ParseFloat(item.TotalValue, 64)
		priceCny := e.ConvertToCNY(price, item.Currency)

		order := emag.EmagOrder{
			OrderID:        item.MktId,
			OrderDateLocal: item.Date,
			Country:        item.CustomerBillingCountry,
			Currency:       item.Currency,
			Price:          price,
			PriceCny:       priceCny,
			Status:         item.SubStatus,
			CustomerName:   item.CustomerName,
		}
		orders = append(orders, order)

		// 处理产品数据
		for _, prod := range item.Products {
			quantity, _ := strconv.ParseFloat(prod.Quantity, 64)
			vatRate, _ := strconv.ParseFloat(prod.VatRate, 64)
			salePrice, _ := strconv.ParseFloat(prod.SalePrice, 64)
			salePriceCny := e.ConvertToCNY(salePrice, item.Currency)

			// 获取产品链接
			productUrl := ""
			if prodUrlAction, ok := prod.Actions["productUrl"]; ok {
				productUrl = prodUrlAction.Action
			}

			product := emag.EmagOrderProduct{
				OrderID:      item.MktId,
				ProductID:    prod.PartNumberKey,
				ProductUrl:   productUrl,
				ProductName:  prod.Name,
				SalePrice:    salePrice,
				SalePriceCny: salePriceCny,
				Quantity:     quantity,
				VatRate:      vatRate,
			}
			products = append(products, product)
		}
	}

	// 使用 Upsert 保存订单
	if len(orders) > 0 {
		if err := global.GVA_DB.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "order_id"}},
			DoUpdates: clause.AssignmentColumns([]string{"order_date_local", "country", "currency", "price", "price_cny", "status", "customer_name", "updated_at"}),
		}).CreateInBatches(orders, 100).Error; err != nil {
			global.GVA_LOG.Error("保存订单失败", zap.Error(err))
			return 0, fmt.Errorf("保存订单失败: %v", err)
		}
	}

	// 保存产品（先删除旧的，再插入新的）
	if len(products) > 0 {
		// 获取所有订单ID
		orderIDs := make([]string, 0)
		for _, order := range orders {
			orderIDs = append(orderIDs, order.OrderID)
		}
		// 删除旧的产品记录
		if err := global.GVA_DB.Where("order_id IN ?", orderIDs).Delete(&emag.EmagOrderProduct{}).Error; err != nil {
			global.GVA_LOG.Error("删除旧产品记录失败", zap.Error(err))
		}
		// 插入新的产品记录
		if err := global.GVA_DB.CreateInBatches(products, 100).Error; err != nil {
			global.GVA_LOG.Error("保存产品失败", zap.Error(err))
			return 0, fmt.Errorf("保存产品失败: %v", err)
		}
	}

	global.GVA_LOG.Info("同步订单成功", zap.Int("count", len(orders)))
	return len(orders), nil
}

// GetOrderList 分页查询订单列表（一个订单一行，产品列表附带）
func (e *EmagOrderService) GetOrderList(info emagReq.EmagOrderSearch) (list []emag.EmagOrderWithProducts, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	if limit <= 0 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}

	db := global.GVA_DB.Model(&emag.EmagOrder{})

	// 条件查询
	if info.OrderID != "" {
		db = db.Where("order_id LIKE ?", "%"+info.OrderID+"%")
	}
	if info.Country != "" {
		db = db.Where("country = ?", info.Country)
	}
	if info.Status != "" {
		db = db.Where("status = ?", info.Status)
	}
	if info.StartDate != "" {
		db = db.Where("order_date_local >= ?", info.StartDate)
	}
	if info.EndDate != "" {
		db = db.Where("order_date_local <= ?", info.EndDate+" 23:59:59")
	}

	// 查询总数
	if err = db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 查询订单列表
	var orders []emag.EmagOrder
	if err = db.Order("order_date_local DESC").Limit(limit).Offset(offset).Find(&orders).Error; err != nil {
		return nil, 0, err
	}

	// 查询每个订单的产品
	list = make([]emag.EmagOrderWithProducts, len(orders))
	for i, order := range orders {
		list[i].EmagOrder = order
		var products []emag.EmagOrderProduct
		if err = global.GVA_DB.Where("order_id = ?", order.OrderID).Find(&products).Error; err != nil {
			global.GVA_LOG.Error("查询订单产品失败", zap.String("orderId", order.OrderID), zap.Error(err))
			continue
		}
		list[i].Products = products
	}

	return list, total, nil
}

// GetCountryList 获取国家列表
func (e *EmagOrderService) GetCountryList() (list []string, err error) {
	err = global.GVA_DB.Model(&emag.EmagOrder{}).
		Select("DISTINCT country").
		Where("country IS NOT NULL AND country != ''").
		Pluck("country", &list).Error
	return list, err
}
