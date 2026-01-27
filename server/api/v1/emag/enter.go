package emag

import "github.com/flipped-aurora/gin-vue-admin/server/service"

type ApiGroup struct {
	EmagCategoryApi
	EmagCategoryStatApi
	EmagOrderApi
	EmagProductApi
}

var (
	emagCategoryService     = service.ServiceGroupApp.EmagServiceGroup.EmagCategoryService
	emagCategoryStatService = service.ServiceGroupApp.EmagServiceGroup.EmagCategoryStatService
	emagOrderService        = service.ServiceGroupApp.EmagServiceGroup.EmagOrderService
	emagProductService      = service.ServiceGroupApp.EmagServiceGroup.EmagProductService
)
