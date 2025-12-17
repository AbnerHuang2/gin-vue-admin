package emag

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/emag"
	emagReq "github.com/flipped-aurora/gin-vue-admin/server/model/emag/request"
)

type EmagCategoryStatService struct{}

var EmagCategoryStatServiceApp = new(EmagCategoryStatService)

// GetSnapshotDateList 获取快照日期列表
func (e *EmagCategoryStatService) GetSnapshotDateList() (list []string, err error) {
	err = global.GVA_DB.Model(&emag.EmagCategoryStat{}).
		Select("DISTINCT DATE_FORMAT(snapshot_date, '%Y-%m-%d') as snapshot_date").
		Order("snapshot_date DESC").
		Pluck("snapshot_date", &list).Error
	return list, err
}

// GetCategoryStatTop20 根据条件查询品类指标Top20（关联品类表获取品类名称）
func (e *EmagCategoryStatService) GetCategoryStatTop20(info emagReq.CategoryStatTop20Search) (list []emag.EmagCategoryStatWithName, total int64, err error) {
	limit := 20
	if info.PageSize > 0 && info.PageSize < 20 {
		limit = info.PageSize
	}

	// 构建基础查询条件
	whereClause := "1=1"
	args := []interface{}{}

	// 根据 supper_hot_rate 大于等于查询
	if info.SupperHotRate > 0 {
		whereClause += " AND s.supper_hot_rate >= ?"
		args = append(args, info.SupperHotRate)
	}

	// 根据 oem_supper_hot_rate 大于等于查询
	if info.OemSupperHotRate > 0 {
		whereClause += " AND s.oem_supper_hot_rate >= ?"
		args = append(args, info.OemSupperHotRate)
	}

	// 根据 snapshot_date 查询
	if info.SnapshotDate != "" {
		whereClause += " AND DATE_FORMAT(s.snapshot_date, '%Y-%m-%d') = ?"
		args = append(args, info.SnapshotDate)
	}

	// 查询总数
	countSQL := `SELECT COUNT(*) FROM emag_category_stat s WHERE ` + whereClause
	err = global.GVA_DB.Raw(countSQL, args...).Scan(&total).Error
	if err != nil {
		return
	}

	// 使用 LEFT JOIN 关联 emag_category 表获取品类名称
	sql := `
		SELECT 
			s.id, s.category_id, s.total, s.supper_hot_total, s.oem_total, s.oem_supper_hot_total,
			s.supper_hot_rate, s.oem_supper_hot_rate, s.snapshot_date, s.tags, s.create_at, s.update_at,
			COALESCE(c.category_name, '') as category_name,
			COALESCE(c.subcategory_name, '') as subcategory_name
		FROM emag_category_stat s
		LEFT JOIN emag_category c ON s.category_id = c.category_id
		WHERE ` + whereClause + `
		ORDER BY s.supper_hot_rate DESC, s.oem_supper_hot_rate DESC
		LIMIT ?
	`
	args = append(args, limit)

	err = global.GVA_DB.Raw(sql, args...).Scan(&list).Error

	return list, total, err
}

// GetCategoryStatGrowthRank 获取品类指标同比增长排名（关联品类表获取品类名称）
func (e *EmagCategoryStatService) GetCategoryStatGrowthRank(info emagReq.CategoryStatGrowthSearch) (list []emag.CategoryStatGrowth, currentDate string, previousDate string, err error) {
	limit := 50
	if info.Limit > 0 {
		limit = info.Limit
	}

	// 1. 获取最新的两个 snapshot_date
	var dates []string
	err = global.GVA_DB.Model(&emag.EmagCategoryStat{}).
		Select("DISTINCT DATE_FORMAT(snapshot_date, '%Y-%m-%d') as snapshot_date").
		Order("snapshot_date DESC").
		Limit(2).
		Pluck("snapshot_date", &dates).Error
	if err != nil {
		return nil, "", "", err
	}

	if len(dates) < 2 {
		return nil, "", "", nil // 数据不足，无法计算同比
	}

	currentDate = dates[0]
	previousDate = dates[1]

	// 2. 使用原生SQL进行JOIN查询计算增长率，同时关联 emag_category 表获取品类名称
	sql := `
		SELECT 
			cur.category_id,
			COALESCE(cat.category_name, '') as category_name,
			COALESCE(cat.subcategory_name, '') as subcategory_name,
			cur.total as current_total,
			COALESCE(prev.total, 0) as previous_total,
			CASE WHEN COALESCE(prev.total, 0) > 0 THEN (cur.total - prev.total) * 100.0 / prev.total ELSE 0 END as total_growth_rate,
			cur.supper_hot_total as current_supper_hot_total,
			COALESCE(prev.supper_hot_total, 0) as previous_supper_hot_total,
			CASE WHEN COALESCE(prev.supper_hot_total, 0) > 0 THEN (cur.supper_hot_total - prev.supper_hot_total) * 100.0 / prev.supper_hot_total ELSE 0 END as supper_hot_growth_rate,
			cur.supper_hot_rate as current_supper_hot_rate,
			COALESCE(prev.supper_hot_rate, 0) as previous_supper_hot_rate,
			cur.oem_supper_hot_rate as current_oem_supper_hot_rate,
			COALESCE(prev.oem_supper_hot_rate, 0) as previous_oem_supper_hot_rate,
			? as current_snapshot_date,
			? as previous_snapshot_date
		FROM emag_category_stat cur
		LEFT JOIN emag_category_stat prev ON cur.category_id = prev.category_id 
			AND DATE_FORMAT(prev.snapshot_date, '%Y-%m-%d') = ?
		LEFT JOIN emag_category cat ON cur.category_id = cat.category_id
		WHERE DATE_FORMAT(cur.snapshot_date, '%Y-%m-%d') = ?
		ORDER BY total_growth_rate DESC
		LIMIT ?
	`

	err = global.GVA_DB.Raw(sql, currentDate, previousDate, previousDate, currentDate, limit).Scan(&list).Error

	return list, currentDate, previousDate, err
}
