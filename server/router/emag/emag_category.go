package emag

import (
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type EmagCategoryRouter struct{}

// InitEmagCategoryRouter 初始化品类管理路由
func (e *EmagCategoryRouter) InitEmagCategoryRouter(Router *gin.RouterGroup) {
	emagCategoryRouter := Router.Group("emagCategory").Use(middleware.OperationRecord())
	emagCategoryRouterWithoutRecord := Router.Group("emagCategory")
	{
		emagCategoryRouter.POST("createEmagCategory", emagCategoryApi.CreateEmagCategory)             // 创建品类
		emagCategoryRouter.DELETE("deleteEmagCategory", emagCategoryApi.DeleteEmagCategory)           // 删除品类
		emagCategoryRouter.DELETE("deleteEmagCategoryByIds", emagCategoryApi.DeleteEmagCategoryByIds) // 批量删除品类
		emagCategoryRouter.PUT("updateEmagCategory", emagCategoryApi.UpdateEmagCategory)              // 更新品类
	}
	{
		emagCategoryRouterWithoutRecord.GET("findEmagCategory", emagCategoryApi.FindEmagCategory)       // 获取单个品类
		emagCategoryRouterWithoutRecord.GET("getEmagCategoryList", emagCategoryApi.GetEmagCategoryList) // 获取品类列表
	}
}
