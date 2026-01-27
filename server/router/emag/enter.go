package emag

import (
	api "github.com/flipped-aurora/gin-vue-admin/server/api/v1"
)

type RouterGroup struct {
	EmagCategoryRouter
	EmagCategoryStatRouter
	EmagOrderRouter
	EmagProductRouter
}

var (
	emagCategoryApi     = api.ApiGroupApp.EmagApiGroup.EmagCategoryApi
	emagCategoryStatApi = api.ApiGroupApp.EmagApiGroup.EmagCategoryStatApi
	emagOrderApi        = api.ApiGroupApp.EmagApiGroup.EmagOrderApi
	emagProductApi      = api.ApiGroupApp.EmagApiGroup.EmagProductApi
)
