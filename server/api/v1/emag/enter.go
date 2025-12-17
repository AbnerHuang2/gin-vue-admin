package emag

import "github.com/flipped-aurora/gin-vue-admin/server/service"

type ApiGroup struct {
	EmagCategoryApi
	EmagCategoryStatApi
}

var (
	emagCategoryService     = service.ServiceGroupApp.EmagServiceGroup.EmagCategoryService
	emagCategoryStatService = service.ServiceGroupApp.EmagServiceGroup.EmagCategoryStatService
)
