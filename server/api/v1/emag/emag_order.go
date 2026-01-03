package emag

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	emagReq "github.com/flipped-aurora/gin-vue-admin/server/model/emag/request"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type EmagOrderApi struct{}

// GetOrderList 分页获取订单列表
// @Tags      EmagOrder
// @Summary   分页获取订单列表
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  query     emagReq.EmagOrderSearch                         true  "查询参数"
// @Success   200   {object}  response.Response{data=response.PageResult,msg=string}  "分页获取订单列表"
// @Router    /emagOrder/getList [get]
func (e *EmagOrderApi) GetOrderList(c *gin.Context) {
	var info emagReq.EmagOrderSearch
	err := c.ShouldBindQuery(&info)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	list, total, err := emagOrderService.GetOrderList(info)
	if err != nil {
		global.GVA_LOG.Error("获取订单列表失败!", zap.Error(err))
		response.FailWithMessage("获取订单列表失败: "+err.Error(), c)
		return
	}

	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     info.Page,
		PageSize: info.PageSize,
	}, "获取成功", c)
}

// SyncOrders 同步订单
// @Tags      EmagOrder
// @Summary   从Emag平台同步订单
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      emagReq.EmagOrderSyncRequest  true  "同步参数"
// @Success   200   {object}  response.Response{msg=string}  "同步成功"
// @Router    /emagOrder/sync [post]
func (e *EmagOrderApi) SyncOrders(c *gin.Context) {
	var req emagReq.EmagOrderSyncRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		// 如果没有传参数，使用默认值
		req.Status = "STATUS_FINALIZED"
		req.Page = 0
		req.Limit = 50
	}

	if req.Status == "" {
		req.Status = "STATUS_FINALIZED"
	}
	if req.Limit <= 0 {
		req.Limit = 50
	}

	// 异步执行同步任务
	go func() {
		global.GVA_LOG.Info("开始同步订单")
		count, err := emagOrderService.SyncOrdersFromEmag(req.Status, req.Page, req.Limit)
		if err != nil {
			global.GVA_LOG.Error("同步订单失败", zap.Error(err))
		} else {
			global.GVA_LOG.Info("同步订单完成", zap.Int("count", count))
		}
	}()

	response.OkWithMessage("同步任务已启动，请稍后刷新页面查看结果", c)
}

// GetCountryList 获取国家列表
// @Tags      EmagOrder
// @Summary   获取国家列表
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Success   200   {object}  response.Response{data=[]string,msg=string}  "获取国家列表"
// @Router    /emagOrder/getCountryList [get]
func (e *EmagOrderApi) GetCountryList(c *gin.Context) {
	list, err := emagOrderService.GetCountryList()
	if err != nil {
		global.GVA_LOG.Error("获取国家列表失败!", zap.Error(err))
		response.FailWithMessage("获取国家列表失败: "+err.Error(), c)
		return
	}
	response.OkWithDetailed(list, "获取成功", c)
}
