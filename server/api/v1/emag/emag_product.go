package emag

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	emagReq "github.com/flipped-aurora/gin-vue-admin/server/model/emag/request"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type EmagProductApi struct{}

// GetProductList 分页获取产品列表
// @Tags      EmagProduct
// @Summary   分页获取产品列表
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  query     emagReq.EmagProductSearch                         true  "查询参数"
// @Success   200   {object}  response.Response{data=response.PageResult,msg=string}  "分页获取产品列表"
// @Router    /emagProduct/getList [get]
func (e *EmagProductApi) GetProductList(c *gin.Context) {
	var info emagReq.EmagProductSearch
	err := c.ShouldBindQuery(&info)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	list, total, err := emagProductService.GetProductList(info)
	if err != nil {
		global.GVA_LOG.Error("获取产品列表失败!", zap.Error(err))
		response.FailWithMessage("获取产品列表失败: "+err.Error(), c)
		return
	}

	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     info.Page,
		PageSize: info.PageSize,
	}, "获取成功", c)
}

// SyncProducts 同步产品
// @Tags      EmagProduct
// @Summary   从Emag平台同步产品
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      emagReq.EmagProductSyncRequest  true  "同步参数"
// @Success   200   {object}  response.Response{msg=string}  "同步成功"
// @Router    /emagProduct/sync [post]
func (e *EmagProductApi) SyncProducts(c *gin.Context) {
	var req emagReq.EmagProductSyncRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		// 如果没有传参数，使用默认值
		req.Page = 0
		req.Limit = 25
	}

	if req.Limit <= 0 {
		req.Limit = 25
	}

	// 异步执行同步任务
	go func() {
		global.GVA_LOG.Info("开始同步产品")
		count, err := emagProductService.SyncProductsFromEmag(req.Page, req.Limit)
		if err != nil {
			global.GVA_LOG.Error("同步产品失败", zap.Error(err))
		} else {
			global.GVA_LOG.Info("同步产品完成", zap.Int("count", count))
		}
	}()

	response.OkWithMessage("同步任务已启动，请稍后刷新页面查看结果", c)
}

// GetStatusList 获取状态列表
// @Tags      EmagProduct
// @Summary   获取状态列表
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Success   200   {object}  response.Response{data=[]string,msg=string}  "获取状态列表"
// @Router    /emagProduct/getStatusList [get]
func (e *EmagProductApi) GetStatusList(c *gin.Context) {
	list, err := emagProductService.GetStatusList()
	if err != nil {
		global.GVA_LOG.Error("获取状态列表失败!", zap.Error(err))
		response.FailWithMessage("获取状态列表失败: "+err.Error(), c)
		return
	}
	response.OkWithDetailed(list, "获取成功", c)
}

// GetCountryList 获取国家列表
// @Tags      EmagProduct
// @Summary   获取国家列表
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Success   200   {object}  response.Response{data=[]string,msg=string}  "获取国家列表"
// @Router    /emagProduct/getCountryList [get]
func (e *EmagProductApi) GetCountryList(c *gin.Context) {
	list, err := emagProductService.GetCountryList()
	if err != nil {
		global.GVA_LOG.Error("获取国家列表失败!", zap.Error(err))
		response.FailWithMessage("获取国家列表失败: "+err.Error(), c)
		return
	}
	response.OkWithDetailed(list, "获取成功", c)
}
