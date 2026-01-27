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

type EmagProductService struct{}

var EmagProductServiceApp = new(EmagProductService)

// 货币到国家的映射
var currencyToCountry = map[string]string{
	"RON": "RO", // 罗马尼亚列伊 -> 罗马尼亚
	"EUR": "BG", // 保加利亚列弗 -> 保加利亚
	"HUF": "HU", // 匈牙利福林 -> 匈牙利
}

// GetCountryFromCurrency 根据货币获取国家代码
func (e *EmagProductService) GetCountryFromCurrency(currency string) string {
	currency = strings.ToUpper(currency)
	if country, ok := currencyToCountry[currency]; ok {
		return country
	}
	return ""
}

// GraphQL 请求和响应结构体
type productGraphQLRequest struct {
	OperationName string      `json:"operationName"`
	Variables     interface{} `json:"variables"`
	Query         string      `json:"query"`
}

type offersVariables struct {
	Filters offersFilters `json:"filters"`
}

type offersFilters struct {
	Pagination         offersPagination `json:"pagination"`
	Sort               []offersSort     `json:"sort"`
	OfferStatus        []string         `json:"offerStatus"`
	InactivationReason []string         `json:"inactivationReason"`
}

type offersPagination struct {
	Page  int `json:"page"`
	Limit int `json:"limit"`
}

type offersSort struct {
	Field     string `json:"field"`
	Direction string `json:"direction"`
}

// Emag API 响应结构
type emagOffersResponse struct {
	Data struct {
		Offers struct {
			Items              []emagOfferItem `json:"items"`
			TotalNumberOfItems int             `json:"totalNumberOfItems"`
		} `json:"offers"`
	} `json:"data"`
}

type emagOfferItem struct {
	ID                                      string            `json:"id"`
	ExtID                                   string            `json:"extId"`
	ExtPartNumber                           string            `json:"extPartNumber"`
	DocProductName                          string            `json:"docProductName"`
	DocProductPartNumberKey                 string            `json:"docProductPartNumberKey"`
	CategoryID                              int               `json:"categoryId"`
	Eans                                    []string          `json:"eans"`
	ExtStock                                int               `json:"extStock"`
	ExtSalePrice                            float64           `json:"extSalePrice"`
	Currency                                string            `json:"currency"`
	Links                                   map[string]string `json:"links"`
	OfferPrice                              emagOfferPrice    `json:"offerPrice"`
	ProductPerformanceBuyButtonRank         int               `json:"productPerformanceBuyButtonRank"`
	ProductPerformanceMultiofferOffersCount int               `json:"productPerformanceMultiofferOffersCount"`
}

type emagOfferPrice struct {
	StatusName string         `json:"statusName"`
	VatValues  map[string]int `json:"vatValues"`
}

// GraphQL 查询语句
const offersQuery = `query offers($filters: OfferFilterInput!) {
  offers(filters: $filters) {
    items {
      id
      hasSmartDealsBadge
      campaignId
      isUnsafeForSite
      sellerId
      sellerName
      extName
      extPartNumber
      eans
      categoryDocId
      categoryId
      commissionEstimation {
        value
        priority
      }
      extId
      valid
      inactivationCriticalities
      extStock
      extSalePrice
      invalidationReason
      extStatus
      extHandlingTime
      fullfilledByEmag
      extUrl
      brandName
      mktId
      type
      offerProperties
      productValidations
      pricingOpportunity
      bestPriceDiff
      documentationUpdateEligibility
      minPrice
      maxPrice
      extWarranty
      modified
      created
      offerDetailsUnfairPrice
      offerDetailsLockerEligibility
      offerDetailsSupplyLeadTime
      offerDetailsEmagClub
      offerDetailsEmagClubType
      greenTax
      docProductId
      docProductName
      docProductPartNumber
      docProductPartNumberKey
      docProductBrandId
      docProductBrandMktId
      docProductBrandName
      productPerformanceSuggestedPrice
      productPerformanceOrderValue1
      productPerformanceOrderValue2
      productPerformanceBuyButtonRank
      productPerformanceMultiofferOffersCount
      productPerformanceMultiofferBestPrice
      productPerformanceMultiofferNoOfReviews
      productPerformanceMultiofferReviewScore
      hasProductMeasurements
      rrp
      rrpSeller
      rrpType
      rrpBadge
      extWeight
      productPerformancePriceGmvUplift30
      productPerformanceLostOrderValue
      productPerformanceNeededStockNext7
      productPerformanceNeededStockNext30
      priceIndexProductCompetitionPrice
      priceIndexModified
      platformId
      currency
      productPerformanceMultiofferReviewScore
      measurements
      offerPrice {
        category
        isEmagClubEligible
        isUnfairPriceForced
        unfairPriceMessage
        statusType
        statusName
        inactivationReason
        hotnessType
        hotnessName
        rrpBadgeColor
        rrpGuideline
        rrpValue
        supplyLeadTimeValues
        superPriceValue
        superPrice
        priceIndex
        priceIndexModified
        handlingTimeValues
        vatId
        startDate
        propertyTypes
        isEanMandatory
        vatValues
        category
        invalidationReasonLabel
        statusDetails
      }
      links {
        edit
        seller
        details
      }
      recycleWarrantiesQuantity
      invalidationReasonDate
      invalidationReasonName
      invalidationSubReason
      invalidationSubReasonName
    }
    totalNumberOfItems
  }
}`

// ConvertToCNY 将指定货币金额转换为人民币 (复用订单管理的汇率)
func (e *EmagProductService) ConvertToCNY(amount float64, currency string) float64 {
	currency = strings.ToUpper(currency)
	if rate, ok := currencyRates[currency]; ok {
		return amount * rate
	}
	return amount
}

// SyncProductsFromEmag 从 Emag 平台同步产品数据（自动分页获取所有数据）
func (e *EmagProductService) SyncProductsFromEmag(page, limit int) (int, error) {
	if limit <= 0 {
		limit = 25
	}

	allProducts := make([]emag.EmagProduct, 0)
	currentPage := 0
	totalFetched := 0

	// 循环获取所有页的数据
	for {
		global.GVA_LOG.Info("正在同步产品", zap.Int("page", currentPage), zap.Int("limit", limit))

		// 构建请求体
		reqBody := productGraphQLRequest{
			OperationName: "offers",
			Variables: offersVariables{
				Filters: offersFilters{
					Pagination: offersPagination{
						Page:  currentPage,
						Limit: limit,
					},
					Sort: []offersSort{
						{Field: "gmv30Ron", Direction: "DESC"},
						{Field: "offerHotness", Direction: "DESC"},
						{Field: "extId", Direction: "ASC"},
					},
					OfferStatus: []string{"ACTIVE", "AUTO_INACTIVATED", "BLOCKED", "INACTIVE"},
					InactivationReason: []string{
						"INACTIVATION_REASON_BRAND_REJECTED",
						"INACTIVATION_REASON_BY_SELLER",
						"INACTIVATION_REASON_CATEGORY_CHANGED",
						"INACTIVATION_REASON_CROSS_BORDER_STATUS",
						"INACTIVATION_REASON_DOCUMENTATION_CHANGED",
						"INACTIVATION_REASON_DUPLICATE_PRODUCT",
						"INACTIVATION_REASON_EAN_BEST_PRICE_RULE",
						"INACTIVATION_REASON_EAN_CONFLICT",
						"INACTIVATION_REASON_INVALID_PRICE",
						"INACTIVATION_REASON_RECYCLE_WARRANTIES",
						"INACTIVATION_REASON_RECYCLE_WARRANTIES_EAN",
						"INACTIVATION_REASON_VAT_MISMATCH",
						"INACTIVATION_REASON_XB_MULTIPLE_OFFERS",
						"INACTIVATION_SOURCE_SPECIFIC_DOCUMENTS_REQUIRED",
						"INACTIVATION_SOURCE_VAT_CLASSIFICATION_CHANGE",
					},
				},
			},
			Query: offersQuery,
		}

		jsonBody, err := json.Marshal(reqBody)
		if err != nil {
			return totalFetched, fmt.Errorf("序列化请求体失败: %v", err)
		}

		// 创建 HTTP 请求
		req, err := http.NewRequest("POST", "https://marketplace.emag.ro/global-listing", bytes.NewBuffer(jsonBody))
		if err != nil {
			return totalFetched, fmt.Errorf("创建请求失败: %v", err)
		}

		// 设置请求头
		req.Header.Set("Accept", "application/json, text/plain, */*")
		req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9,en-US;q=0.8,en;q=0.7")
		req.Header.Set("Cache-Control", "no-cache")
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Origin", "https://marketplace.emag.ro")
		req.Header.Set("Pragma", "no-cache")
		req.Header.Set("Priority", "u=1, i")
		req.Header.Set("Referer", "https://marketplace.emag.ro/offers/list-xb")
		req.Header.Set("Sec-Ch-Ua", `"Not(A:Brand";v="8", "Chromium";v="144", "Google Chrome";v="144"`)
		req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
		req.Header.Set("Sec-Ch-Ua-Platform", `"macOS"`)
		req.Header.Set("Sec-Fetch-Dest", "empty")
		req.Header.Set("Sec-Fetch-Mode", "cors")
		req.Header.Set("Sec-Fetch-Site", "same-origin")
		req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/144.0.0.0 Safari/537.36")
		req.Header.Set("X-Requested-With", "XMLHttpRequest")
		req.Header.Set("Cookie", global.GVA_CONFIG.Emag.Cookie)

		// 发送请求
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return totalFetched, fmt.Errorf("发送请求失败: %v", err)
		}

		// 读取响应
		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			return totalFetched, fmt.Errorf("读取响应失败: %v", err)
		}

		if resp.StatusCode != http.StatusOK {
			return totalFetched, fmt.Errorf("请求失败, 状态码: %d, 响应: %s", resp.StatusCode, string(body))
		}

		// 解析响应
		var emagResp emagOffersResponse
		if err := json.Unmarshal(body, &emagResp); err != nil {
			return totalFetched, fmt.Errorf("解析响应失败: %v, 响应内容: %s", err, string(body))
		}

		// 处理当前页的产品数据
		products := make([]emag.EmagProduct, 0)
		for _, item := range emagResp.Data.Offers.Items {
			// 提取 EAN (取数组第一个)
			ean := ""
			if len(item.Eans) > 0 {
				ean = item.Eans[0]
			}

			// 提取 VAT (取 vatValues 的第一个 key)
			vat := ""
			vatRate := 0.0
			for k := range item.OfferPrice.VatValues {
				vat = k
				// 解析 VAT 率，例如 "21%" -> 0.21
				vatStr := strings.TrimSuffix(k, "%")
				if vatFloat, err := strconv.ParseFloat(vatStr, 64); err == nil {
					vatRate = vatFloat / 100
				}
				break
			}

			// 计算含税价格
			afterTaxPrice := item.ExtSalePrice * (1 + vatRate)

			// 转换为人民币
			salePriceCN := e.ConvertToCNY(item.ExtSalePrice, item.Currency)

			// 提取 URL
			url := ""
			if detailsUrl, ok := item.Links["details"]; ok {
				url = detailsUrl
			}

			// 根据货币获取国家
			country := e.GetCountryFromCurrency(item.Currency)

			product := emag.EmagProduct{
				ProductID:     item.ID,
				CategoryID:    item.CategoryID,
				ExtID:         item.ExtID,
				PN:            item.ExtPartNumber,
				PNK:           item.DocProductPartNumberKey,
				EAN:           ean,
				Title:         item.DocProductName,
				Status:        item.OfferPrice.StatusName,
				SalePrice:     item.ExtSalePrice,
				AfterTaxPrice: afterTaxPrice,
				Currency:      item.Currency,
				Country:       country,
				SalePriceCN:   salePriceCN,
				VAT:           vat,
				Stock:         item.ExtStock,
				URL:           url,
				BuyButtonRank: item.ProductPerformanceBuyButtonRank,
				BuyButtonCnt:  item.ProductPerformanceMultiofferOffersCount,
			}
			products = append(products, product)
		}

		allProducts = append(allProducts, products...)
		totalFetched += len(products)

		global.GVA_LOG.Info("当前页同步完成",
			zap.Int("page", currentPage),
			zap.Int("currentPageCount", len(products)),
			zap.Int("totalFetched", totalFetched),
			zap.Int("totalAvailable", emagResp.Data.Offers.TotalNumberOfItems))

		// 判断是否已获取所有数据
		if len(products) == 0 || totalFetched >= emagResp.Data.Offers.TotalNumberOfItems {
			global.GVA_LOG.Info("所有产品已获取完毕", zap.Int("totalFetched", totalFetched))
			break
		}

		currentPage++
	}

	// 批量保存所有产品
	if len(allProducts) > 0 {
		if err := global.GVA_DB.Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "product_id"}},
			DoUpdates: clause.AssignmentColumns([]string{
				"category_id", "ext_id", "pn", "pnk", "ean", "title", "status",
				"sale_price", "after_tax_price", "currency", "country", "sale_price_cn",
				"vat", "stock", "url", "buy_button_rank", "buy_button_cnt", "updated_at",
			}),
		}).CreateInBatches(allProducts, 100).Error; err != nil {
			global.GVA_LOG.Error("保存产品失败", zap.Error(err))
			return totalFetched, fmt.Errorf("保存产品失败: %v", err)
		}
	}

	global.GVA_LOG.Info("同步产品成功", zap.Int("count", totalFetched))
	return totalFetched, nil
}

// GetProductList 分页查询产品列表
func (e *EmagProductService) GetProductList(info emagReq.EmagProductSearch) (list []emag.EmagProduct, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	if limit <= 0 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}

	db := global.GVA_DB.Model(&emag.EmagProduct{})

	// 条件查询
	if info.PNK != "" {
		db = db.Where("pnk LIKE ?", "%"+info.PNK+"%")
	}
	if info.ExtID != "" {
		db = db.Where("ext_id LIKE ?", "%"+info.ExtID+"%")
	}
	if info.StockMin > 0 {
		db = db.Where("stock >= ?", info.StockMin)
	}
	if info.Country != "" {
		db = db.Where("country = ?", info.Country)
	}
	// 保留旧的筛选条件（向后兼容）
	if info.PN != "" {
		db = db.Where("pn LIKE ?", "%"+info.PN+"%")
	}
	if info.Title != "" {
		db = db.Where("title LIKE ?", "%"+info.Title+"%")
	}
	if info.Status != "" {
		db = db.Where("status = ?", info.Status)
	}

	// 查询总数
	if err = db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 查询产品列表 - 按照 ext_id 排序
	if err = db.Order("ext_id ASC").Limit(limit).Offset(offset).Find(&list).Error; err != nil {
		return nil, 0, err
	}

	return list, total, nil
}

// GetStatusList 获取状态列表
func (e *EmagProductService) GetStatusList() (list []string, err error) {
	err = global.GVA_DB.Model(&emag.EmagProduct{}).
		Select("DISTINCT status").
		Where("status IS NOT NULL AND status != ''").
		Pluck("status", &list).Error
	return list, err
}

// GetCountryList 获取国家列表
func (e *EmagProductService) GetCountryList() (list []string, err error) {
	err = global.GVA_DB.Model(&emag.EmagProduct{}).
		Select("DISTINCT country").
		Where("country IS NOT NULL AND country != ''").
		Pluck("country", &list).Error
	return list, err
}
