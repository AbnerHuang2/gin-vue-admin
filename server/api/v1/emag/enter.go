package emag

import "github.com/flipped-aurora/gin-vue-admin/server/service"

type ApiGroup struct {
	EmagCategoryApi
}

var (
	emagCategoryService = service.ServiceGroupApp.EmagServiceGroup.EmagCategoryService
)
