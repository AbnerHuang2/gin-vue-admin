package emag

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/emag"
	emagReq "github.com/flipped-aurora/gin-vue-admin/server/model/emag/request"
)

type EmagCategoryService struct{}

var EmagCategoryServiceApp = new(EmagCategoryService)

// CreateEmagCategory 创建品类记录
func (e *EmagCategoryService) CreateEmagCategory(category emag.EmagCategory) (err error) {
	err = global.GVA_DB.Create(&category).Error
	return err
}

// DeleteEmagCategory 删除品类记录（通过category_id）
func (e *EmagCategoryService) DeleteEmagCategory(categoryId string) (err error) {
	err = global.GVA_DB.Where("category_id = ?", categoryId).Delete(&emag.EmagCategory{}).Error
	return err
}

// DeleteEmagCategoryByIds 批量删除品类记录
func (e *EmagCategoryService) DeleteEmagCategoryByIds(ids []uint) (err error) {
	err = global.GVA_DB.Delete(&[]emag.EmagCategory{}, "id in ?", ids).Error
	return err
}

// UpdateEmagCategory 更新品类记录（通过category_id）
func (e *EmagCategoryService) UpdateEmagCategory(category emag.EmagCategory) (err error) {
	err = global.GVA_DB.Model(&emag.EmagCategory{}).Where("category_id = ?", category.CategoryId).Updates(map[string]interface{}{
		"category_name":    category.CategoryName,
		"subcategory_name": category.SubcategoryName,
	}).Error
	return err
}

// GetEmagCategory 根据id或category_name获取品类记录
func (e *EmagCategoryService) GetEmagCategory(info emagReq.EmagCategorySearch) (category emag.EmagCategory, err error) {
	db := global.GVA_DB.Model(&emag.EmagCategory{})
	if info.ID != 0 {
		db = db.Where("id = ?", info.ID)
	}
	if info.CategoryName != "" {
		db = db.Where("category_name LIKE ?", "%"+info.CategoryName+"%")
	}
	if info.CategoryId != "" {
		db = db.Where("category_id = ?", info.CategoryId)
	}
	err = db.First(&category).Error
	return
}

// GetEmagCategoryList 分页获取品类列表
func (e *EmagCategoryService) GetEmagCategoryList(info emagReq.EmagCategorySearch) (list []emag.EmagCategory, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	db := global.GVA_DB.Model(&emag.EmagCategory{})

	// 根据id查询
	if info.ID != 0 {
		db = db.Where("id = ?", info.ID)
	}
	// 根据category_name模糊查询
	if info.CategoryName != "" {
		db = db.Where("category_name LIKE ?", "%"+info.CategoryName+"%")
	}
	// 根据category_id查询
	if info.CategoryId != "" {
		db = db.Where("category_id = ?", info.CategoryId)
	}

	err = db.Count(&total).Error
	if err != nil {
		return
	}
	err = db.Limit(limit).Offset(offset).Order("id desc").Find(&list).Error
	return list, total, err
}
