package emag

import (
	"time"

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

// GetActiveCategoryIds 获取状态为 normal 的所有 category_id
func (e *EmagCategoryService) GetActiveCategoryIds() ([]string, error) {
	var categoryIds []string
	err := global.GVA_DB.Model(&emag.EmagCategory{}).
		Where("status = ? OR status IS NULL OR status = ''", emag.CategoryStatusNormal).
		Pluck("category_id", &categoryIds).Error
	return categoryIds, err
}

// GetAllCategoryIdsWithInfo 获取所有活跃分类的ID和信息（用于定时任务）
func (e *EmagCategoryService) GetAllCategoryIdsWithInfo() ([]emag.EmagCategory, error) {
	var categories []emag.EmagCategory
	err := global.GVA_DB.Model(&emag.EmagCategory{}).
		Where("status = ? OR status IS NULL OR status = ''", emag.CategoryStatusNormal).
		Select("id", "category_id", "category_name", "fail_count").
		Order("id ASC").
		Find(&categories).Error
	return categories, err
}

// IncrementFailCount 增加失败计数，返回更新后的失败次数
func (e *EmagCategoryService) IncrementFailCount(categoryId string, reason string) (int, error) {
	now := time.Now()

	// 先获取当前的 fail_count
	var category emag.EmagCategory
	err := global.GVA_DB.Model(&emag.EmagCategory{}).
		Where("category_id = ?", categoryId).
		Select("fail_count").
		First(&category).Error
	if err != nil {
		return 0, err
	}

	newFailCount := category.FailCount + 1

	// 更新失败计数和相关信息
	err = global.GVA_DB.Model(&emag.EmagCategory{}).
		Where("category_id = ?", categoryId).
		Updates(map[string]interface{}{
			"fail_count":       newFailCount,
			"last_fail_reason": reason,
			"last_fail_at":     &now,
		}).Error

	return newFailCount, err
}

// MarkAsBadRequest 标记为 bad_request 状态
func (e *EmagCategoryService) MarkAsBadRequest(categoryId string, reason string) error {
	now := time.Now()
	return global.GVA_DB.Model(&emag.EmagCategory{}).
		Where("category_id = ?", categoryId).
		Updates(map[string]interface{}{
			"status":           emag.CategoryStatusBadRequest,
			"last_fail_reason": reason,
			"last_fail_at":     &now,
		}).Error
}

// ResetFailCount 重置失败计数（成功时调用）
func (e *EmagCategoryService) ResetFailCount(categoryId string) error {
	return global.GVA_DB.Model(&emag.EmagCategory{}).
		Where("category_id = ?", categoryId).
		Updates(map[string]interface{}{
			"fail_count": 0,
			"status":     emag.CategoryStatusNormal,
		}).Error
}

// ResetBadRequestStatus 手动重置 bad_request 状态为 normal（提供 API 供管理员使用）
func (e *EmagCategoryService) ResetBadRequestStatus(categoryId string) error {
	return global.GVA_DB.Model(&emag.EmagCategory{}).
		Where("category_id = ?", categoryId).
		Updates(map[string]interface{}{
			"status":           emag.CategoryStatusNormal,
			"fail_count":       0,
			"last_fail_reason": "",
			"last_fail_at":     nil,
		}).Error
}

// ResetAllBadRequestStatus 重置所有 bad_request 状态为 normal
func (e *EmagCategoryService) ResetAllBadRequestStatus() (int64, error) {
	result := global.GVA_DB.Model(&emag.EmagCategory{}).
		Where("status = ?", emag.CategoryStatusBadRequest).
		Updates(map[string]interface{}{
			"status":           emag.CategoryStatusNormal,
			"fail_count":       0,
			"last_fail_reason": "",
			"last_fail_at":     nil,
		})
	return result.RowsAffected, result.Error
}

// GetBadRequestCategories 获取所有 bad_request 状态的分类
func (e *EmagCategoryService) GetBadRequestCategories() ([]emag.EmagCategory, error) {
	var categories []emag.EmagCategory
	err := global.GVA_DB.Model(&emag.EmagCategory{}).
		Where("status = ?", emag.CategoryStatusBadRequest).
		Find(&categories).Error
	return categories, err
}

// MarkAsNotCare 标记为 not_care 状态（不关注）
func (e *EmagCategoryService) MarkAsNotCare(categoryId string) error {
	return global.GVA_DB.Model(&emag.EmagCategory{}).
		Where("category_id = ?", categoryId).
		Update("status", emag.CategoryStatusNotCare).Error
}

// UpdateCategoryStatus 更新品类状态
func (e *EmagCategoryService) UpdateCategoryStatus(categoryId string, status string) error {
	return global.GVA_DB.Model(&emag.EmagCategory{}).
		Where("category_id = ?", categoryId).
		Update("status", status).Error
}
