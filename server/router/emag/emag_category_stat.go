package emag

import (
	"github.com/gin-gonic/gin"
)

type EmagCategoryStatRouter struct{}

// InitEmagCategoryStatRouter 初始化品类指标路由
func (e *EmagCategoryStatRouter) InitEmagCategoryStatRouter(Router *gin.RouterGroup) {
	emagCategoryStatRouter := Router.Group("emagCategoryStat")
	{
		emagCategoryStatRouter.GET("getSnapshotDateList", emagCategoryStatApi.GetSnapshotDateList) // 获取快照日期列表
		emagCategoryStatRouter.GET("getTop20", emagCategoryStatApi.GetCategoryStatTop20)           // 获取品类指标Top20
		emagCategoryStatRouter.GET("getGrowthRank", emagCategoryStatApi.GetCategoryStatGrowthRank) // 获取品类指标同比增长排名
	}
}
