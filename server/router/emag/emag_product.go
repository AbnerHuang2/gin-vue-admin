package emag

import (
	"github.com/gin-gonic/gin"
)

type EmagProductRouter struct{}

// InitEmagProductRouter 初始化产品路由
func (e *EmagProductRouter) InitEmagProductRouter(Router *gin.RouterGroup) {
	emagProductRouter := Router.Group("emagProduct")
	{
		emagProductRouter.GET("getList", emagProductApi.GetProductList)        // 分页获取产品列表
		emagProductRouter.POST("sync", emagProductApi.SyncProducts)            // 同步产品
		emagProductRouter.GET("getStatusList", emagProductApi.GetStatusList)   // 获取状态列表
		emagProductRouter.GET("getCountryList", emagProductApi.GetCountryList) // 获取国家列表
	}
}
