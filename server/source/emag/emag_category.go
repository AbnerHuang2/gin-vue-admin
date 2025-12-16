package emag

import (
	"context"

	"github.com/flipped-aurora/gin-vue-admin/server/model/emag"
	"github.com/flipped-aurora/gin-vue-admin/server/service/system"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

const initOrderEmagCategory = system.InitOrderSystem + 10

type initEmagCategory struct{}

// auto run
func init() {
	system.RegisterInit(initOrderEmagCategory, &initEmagCategory{})
}

func (i *initEmagCategory) InitializerName() string {
	return emag.EmagCategory{}.TableName()
}

func (i *initEmagCategory) MigrateTable(ctx context.Context) (context.Context, error) {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return ctx, system.ErrMissingDBContext
	}
	return ctx, db.AutoMigrate(&emag.EmagCategory{})
}

func (i *initEmagCategory) TableCreated(ctx context.Context) bool {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return false
	}
	return db.Migrator().HasTable(&emag.EmagCategory{})
}

func (i *initEmagCategory) InitializeData(ctx context.Context) (context.Context, error) {
	// 可选：初始化一些默认数据
	return ctx, nil
}

func (i *initEmagCategory) DataInserted(ctx context.Context) bool {
	return true
}
