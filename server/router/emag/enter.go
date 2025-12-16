package emag

import (
	api "github.com/flipped-aurora/gin-vue-admin/server/api/v1"
)

type RouterGroup struct {
	EmagCategoryRouter
}

var (
	emagCategoryApi = api.ApiGroupApp.EmagApiGroup.EmagCategoryApi
)
