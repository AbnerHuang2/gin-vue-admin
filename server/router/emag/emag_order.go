package emag

import (
	"github.com/gin-gonic/gin"
)

type EmagOrderRouter struct{}

// InitEmagOrderRouter 初始化订单路由
func (e *EmagOrderRouter) InitEmagOrderRouter(Router *gin.RouterGroup) {
	emagOrderRouter := Router.Group("emagOrder")
	{
		emagOrderRouter.GET("getList", emagOrderApi.GetOrderList)          // 分页获取订单列表
		emagOrderRouter.POST("sync", emagOrderApi.SyncOrders)              // 同步订单
		emagOrderRouter.GET("getCountryList", emagOrderApi.GetCountryList) // 获取国家列表
	}
}
