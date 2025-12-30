package task

import (
	"fmt"
	"strconv"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/emag"
	emagService "github.com/flipped-aurora/gin-vue-admin/server/service/emag"
	"go.uber.org/zap"
)

var processFlag = false

// UpdateCategoryStat 更新品类统计数据（定时任务入口）
func UpdateCategoryStat() error {
	if processFlag {
		global.GVA_LOG.Info("品类统计数据更新任务已在运行中，跳过本次执行")
		return nil
	}
	processFlag = true
	startTime := time.Now()
	global.GVA_LOG.Info("========== 开始更新品类统计数据 ==========")

	// 获取配置
	config := global.GVA_CONFIG.Emag
	maxFailCount := config.MaxFailCount
	if maxFailCount <= 0 {
		maxFailCount = 3
	}
	snapshotDayGap := config.SnapshotDayGap
	if snapshotDayGap <= 0 {
		snapshotDayGap = 15
	}
	batchSize := config.BatchSize
	if batchSize <= 0 {
		batchSize = 50
	}

	// 初始化服务
	categoryService := emagService.EmagCategoryServiceApp
	categoryStatService := emagService.EmagCategoryStatServiceApp
	apiClient := emagService.NewEmagAPIClient()

	// Step 1: 确定 snapshot_date
	snapshotDate, err := determineSnapshotDate(categoryStatService, snapshotDayGap)
	if err != nil {
		global.GVA_LOG.Error("确定快照日期失败", zap.Error(err))
		return err
	}
	global.GVA_LOG.Info("使用快照日期", zap.Time("snapshotDate", snapshotDate))

	// Step 2: 获取待处理列表
	pendingCategoryIds, err := getPendingCategoryIds(categoryService, categoryStatService, snapshotDate)
	if err != nil {
		global.GVA_LOG.Error("获取待处理分类列表失败", zap.Error(err))
		return err
	}

	if len(pendingCategoryIds) == 0 {
		global.GVA_LOG.Info("🎉 所有分类统计数据都已是最新的！")
		return nil
	}

	global.GVA_LOG.Info("待处理分类数量", zap.Int("count", len(pendingCategoryIds)))

	// Step 3: 分批处理
	successCount := 0
	failCount := 0
	totalCount := len(pendingCategoryIds)

	for i, categoryId := range pendingCategoryIds {
		progress := fmt.Sprintf("[%d/%d]", i+1, totalCount)

		// 处理单个分类
		err := processSingleCategory(
			categoryId,
			snapshotDate,
			apiClient,
			categoryService,
			categoryStatService,
			maxFailCount,
		)

		if err != nil {
			failCount++
			global.GVA_LOG.Warn(progress+" 处理分类失败",
				zap.String("categoryId", categoryId),
				zap.Error(err))
		} else {
			successCount++
			global.GVA_LOG.Info(progress+" 处理分类成功",
				zap.String("categoryId", categoryId))
		}

		// 每批处理完成后休息一下（每 batchSize 个）
		if (i+1)%batchSize == 0 && i+1 < totalCount {
			global.GVA_LOG.Info("批次处理完成，休息中...",
				zap.Int("processed", i+1),
				zap.Int("total", totalCount))
			time.Sleep(5 * time.Second)
		}
	}

	// Step 4: 完成，记录统计信息
	duration := time.Since(startTime)
	global.GVA_LOG.Info("========== 品类统计数据更新完成 ==========",
		zap.Int("successCount", successCount),
		zap.Int("failCount", failCount),
		zap.Int("totalCount", totalCount),
		zap.Duration("duration", duration))

	if failCount > 0 {
		global.GVA_LOG.Warn(fmt.Sprintf("⚠️ %d 个分类处理失败", failCount))
	}
	processFlag = false
	return nil
}

// determineSnapshotDate 确定快照日期
func determineSnapshotDate(statService *emagService.EmagCategoryStatService, snapshotDayGap int) (time.Time, error) {
	today := time.Now().Truncate(24 * time.Hour)

	// 获取最新的快照日期
	latestDate, err := statService.GetLatestSnapshotDate()
	if err != nil {
		return today, err
	}

	// 如果没有历史数据，使用今天
	if latestDate == nil {
		global.GVA_LOG.Info("没有历史快照数据，使用当前日期")
		return today, nil
	}

	// 计算日期差
	daysDiff := int(today.Sub(*latestDate).Hours() / 24)
	global.GVA_LOG.Info("快照日期计算",
		zap.Time("latestDate", *latestDate),
		zap.Int("daysDiff", daysDiff),
		zap.Int("snapshotDayGap", snapshotDayGap))

	// 如果距离上次快照超过配置的天数，使用新日期
	if daysDiff > snapshotDayGap {
		global.GVA_LOG.Info("距离上次快照超过配置天数，使用新日期")
		return today, nil
	}

	// 否则继续使用上次的快照日期（续传模式）
	global.GVA_LOG.Info("继续使用上次快照日期（续传模式）")
	return *latestDate, nil
}

// getPendingCategoryIds 获取待处理的分类ID列表
func getPendingCategoryIds(
	categoryService *emagService.EmagCategoryService,
	statService *emagService.EmagCategoryStatService,
	snapshotDate time.Time,
) ([]string, error) {
	// 获取所有活跃的分类ID
	allCategoryIds, err := categoryService.GetActiveCategoryIds()
	if err != nil {
		return nil, fmt.Errorf("获取活跃分类失败: %w", err)
	}
	global.GVA_LOG.Info("活跃分类总数", zap.Int("count", len(allCategoryIds)))

	// 获取当前快照日期已处理的分类ID
	processedIds, err := statService.GetProcessedCategoryIds(snapshotDate)
	if err != nil {
		return nil, fmt.Errorf("获取已处理分类失败: %w", err)
	}
	global.GVA_LOG.Info("已处理分类数", zap.Int("count", len(processedIds)))

	// 计算差集：待处理 = 全部 - 已处理
	processedMap := make(map[string]bool)
	for _, id := range processedIds {
		processedMap[id] = true
	}

	var pendingIds []string
	for _, id := range allCategoryIds {
		if !processedMap[id] {
			pendingIds = append(pendingIds, id)
		}
	}

	return pendingIds, nil
}

// processSingleCategory 处理单个分类
func processSingleCategory(
	categoryId string,
	snapshotDate time.Time,
	apiClient *emagService.EmagAPIClient,
	categoryService *emagService.EmagCategoryService,
	statService *emagService.EmagCategoryStatService,
	maxFailCount int,
) error {
	// 转换 categoryId 为 int
	categoryIdInt, err := strconv.Atoi(categoryId)
	if err != nil {
		return fmt.Errorf("无效的分类ID: %w", err)
	}

	// 调用 API 获取统计数据
	stats, err := apiClient.GetCategoryStatistics(categoryIdInt)
	if err != nil {
		// 失败处理：增加失败计数
		failCount, _ := categoryService.IncrementFailCount(categoryId, err.Error())

		// 检查是否需要标记为 bad_request
		if failCount >= maxFailCount {
			categoryService.MarkAsBadRequest(categoryId, err.Error())
			global.GVA_LOG.Warn("分类已标记为 bad_request",
				zap.String("categoryId", categoryId),
				zap.Int("failCount", failCount))
		}

		return err
	}

	// 成功处理
	// 1. 重置失败计数
	if err := categoryService.ResetFailCount(categoryId); err != nil {
		global.GVA_LOG.Warn("重置失败计数失败", zap.String("categoryId", categoryId), zap.Error(err))
	}

	// 2. 计算占比率
	supperHotRate := 0.0
	if stats.Total > 0 {
		supperHotRate = float64(stats.SupperHotTotal) / float64(stats.Total)
	}
	oemSupperHotRate := 0.0
	if stats.SupperHotTotal > 0 {
		oemSupperHotRate = float64(stats.OemSupperHotTotal) / float64(stats.SupperHotTotal)
	}

	// 3. 写入统计数据
	categoryStat := &emag.EmagCategoryStat{
		CategoryId:        categoryId,
		Total:             stats.Total,
		SupperHotTotal:    stats.SupperHotTotal,
		OemTotal:          stats.OemTotal,
		OemSupperHotTotal: stats.OemSupperHotTotal,
		SupperHotRate:     supperHotRate,
		OemSupperHotRate:  oemSupperHotRate,
		SnapshotDate:      &snapshotDate,
	}

	if err := statService.CreateCategoryStat(categoryStat); err != nil {
		return fmt.Errorf("写入统计数据失败: %w", err)
	}

	return nil
}
