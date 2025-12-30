package emag

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"go.uber.org/zap"
)

// CategoryStats 分类统计数据
type CategoryStats struct {
	Total             int // 总数
	SupperHotTotal    int // 超热销总数
	OemTotal          int // OEM总数
	OemSupperHotTotal int // OEM超热销总数
}

// APIResponse Emag API 响应结构
type APIResponse struct {
	Data struct {
		Meta struct {
			TotalCount int `json:"total_count"`
		} `json:"meta"`
	} `json:"data"`
}

// EmagAPIClient Emag API 客户端
type EmagAPIClient struct {
	apiUrl     string
	headers    map[string]string
	interval   time.Duration
	retryCount int
	httpClient *http.Client
}

// NewEmagAPIClient 创建 Emag API 客户端
func NewEmagAPIClient() *EmagAPIClient {
	config := global.GVA_CONFIG.Emag

	// 设置默认值
	interval := config.RequestInterval
	if interval <= 0 {
		interval = 1.5
	}
	retryCount := config.RetryCount
	if retryCount <= 0 {
		retryCount = 3
	}

	return &EmagAPIClient{
		apiUrl:     config.ApiUrl,
		interval:   time.Duration(interval * float64(time.Second)),
		retryCount: retryCount,
		headers: map[string]string{
			"Accept":           "application/json, text/plain, */*",
			"Accept-Language":  "zh-CN,zh;q=0.9,en-US;q=0.8,en;q=0.7",
			"Cache-Control":    "no-cache",
			"Connection":       "keep-alive",
			"Content-Type":     "application/json",
			"Origin":           "https://marketplace.emag.ro",
			"Pragma":           "no-cache",
			"Referer":          "https://marketplace.emag.ro/opportunities/list",
			"Sec-Fetch-Dest":   "empty",
			"Sec-Fetch-Mode":   "cors",
			"Sec-Fetch-Site":   "same-origin",
			"User-Agent":       "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/137.0.0.0 Safari/537.36",
			"X-Requested-With": "XMLHttpRequest",
			"Cookie":           config.Cookie,
		},
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// callAPIWithRetry 带重试的 API 调用
func (c *EmagAPIClient) callAPIWithRetry(payload map[string]interface{}) (*APIResponse, error) {
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("marshal payload error: %w", err)
	}

	var lastErr error
	for attempt := 0; attempt < c.retryCount; attempt++ {
		req, err := http.NewRequest("POST", c.apiUrl, bytes.NewBuffer(jsonData))
		if err != nil {
			lastErr = fmt.Errorf("create request error: %w", err)
			continue
		}

		// 设置请求头
		for key, value := range c.headers {
			req.Header.Set(key, value)
		}

		resp, err := c.httpClient.Do(req)
		if err != nil {
			lastErr = fmt.Errorf("request error: %w", err)
			global.GVA_LOG.Warn("API调用异常", zap.Int("attempt", attempt+1), zap.Error(err))
			// 重试前等待
			if attempt < c.retryCount-1 {
				time.Sleep(c.interval * time.Duration(attempt+1))
			}
			continue
		}

		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()

		if err != nil {
			lastErr = fmt.Errorf("read response error: %w", err)
			continue
		}

		if resp.StatusCode == 200 {
			var apiResp APIResponse
			if err := json.Unmarshal(body, &apiResp); err != nil {
				lastErr = fmt.Errorf("unmarshal response error: %w", err)
				continue
			}
			return &apiResp, nil
		}

		// 处理限流
		if resp.StatusCode == 429 {
			lastErr = errors.New("rate limited (429)")
			global.GVA_LOG.Warn("被限流，等待后重试", zap.Int("attempt", attempt+1))
			// 限流时等待更长时间
			time.Sleep(c.interval * time.Duration((attempt+1)*2))
			continue
		}

		lastErr = fmt.Errorf("API返回状态码: %d, body: %s", resp.StatusCode, string(body[:min(len(body), 200)]))
		global.GVA_LOG.Warn("API调用失败", zap.Int("statusCode", resp.StatusCode))

		// 重试前等待
		if attempt < c.retryCount-1 {
			time.Sleep(c.interval * time.Duration(attempt+1))
		}
	}

	return nil, fmt.Errorf("API调用失败，已达到最大重试次数: %w", lastErr)
}

// GetCategoryStatistics 获取单个分类的统计数据
func (c *EmagAPIClient) GetCategoryStatistics(categoryId int) (*CategoryStats, error) {
	stats := &CategoryStats{}

	// 基础请求参数
	basePayload := map[string]interface{}{
		"per_page": 10,
		"category": categoryId,
		"sort": []map[string]interface{}{
			{
				"field":     "performance",
				"direction": "asc",
			},
		},
	}

	// 1. 获取总数
	resp, err := c.callAPIWithRetry(basePayload)
	if err != nil {
		return nil, fmt.Errorf("获取总数失败: %w", err)
	}
	stats.Total = resp.Data.Meta.TotalCount

	// 如果总数为0，直接返回
	if stats.Total == 0 {
		global.GVA_LOG.Debug("分类总数为0，跳过后续统计", zap.Int("categoryId", categoryId))
		return stats, nil
	}

	// 等待间隔（添加随机抖动）
	c.waitWithJitter()

	// 2. 获取热销总数 (performance: [1])
	hotPayload := copyMap(basePayload)
	hotPayload["performance"] = []int{1}
	resp, err = c.callAPIWithRetry(hotPayload)
	if err != nil {
		return nil, fmt.Errorf("获取热销总数失败: %w", err)
	}
	stats.SupperHotTotal = resp.Data.Meta.TotalCount

	c.waitWithJitter()

	// 3. 获取OEM总数 (brand: 35)
	oemPayload := copyMap(basePayload)
	oemPayload["brand"] = 35
	resp, err = c.callAPIWithRetry(oemPayload)
	if err != nil {
		return nil, fmt.Errorf("获取OEM总数失败: %w", err)
	}
	stats.OemTotal = resp.Data.Meta.TotalCount

	c.waitWithJitter()

	// 4. 获取OEM热销总数 (brand: 35, performance: [1])
	oemHotPayload := copyMap(basePayload)
	oemHotPayload["brand"] = 35
	oemHotPayload["performance"] = []int{1}
	resp, err = c.callAPIWithRetry(oemHotPayload)
	if err != nil {
		return nil, fmt.Errorf("获取OEM热销总数失败: %w", err)
	}
	stats.OemSupperHotTotal = resp.Data.Meta.TotalCount

	return stats, nil
}

// waitWithJitter 等待间隔（添加随机抖动）
func (c *EmagAPIClient) waitWithJitter() {
	// 添加 0-50% 的随机抖动
	jitter := time.Duration(rand.Float64() * 0.5 * float64(c.interval))
	time.Sleep(c.interval + jitter)
}

// copyMap 复制 map
func copyMap(m map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	for k, v := range m {
		result[k] = v
	}
	return result
}

// min 返回两个整数中的较小值
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
