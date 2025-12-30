package emag

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	emagReq "github.com/flipped-aurora/gin-vue-admin/server/model/emag/request"
	"github.com/flipped-aurora/gin-vue-admin/server/task"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type EmagCategoryStatApi struct{}

// GetSnapshotDateList 获取快照日期列表
// @Tags      EmagCategoryStat
// @Summary   获取快照日期列表
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Success   200   {object}  response.Response{data=[]string,msg=string}  "获取快照日期列表"
// @Router    /emagCategoryStat/getSnapshotDateList [get]
func (e *EmagCategoryStatApi) GetSnapshotDateList(c *gin.Context) {
	list, err := emagCategoryStatService.GetSnapshotDateList()
	if err != nil {
		global.GVA_LOG.Error("获取快照日期列表失败!", zap.Error(err))
		response.FailWithMessage("获取快照日期列表失败: "+err.Error(), c)
		return
	}
	response.OkWithDetailed(list, "获取成功", c)
}

// GetCategoryStatList 分页获取品类指标列表
// @Tags      EmagCategoryStat
// @Summary   分页获取品类指标列表
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  query     emagReq.CategoryStatSearch                         true  "查询参数"
// @Success   200   {object}  response.Response{data=response.PageResult,msg=string}  "分页获取品类指标列表"
// @Router    /emagCategoryStat/getList [get]
func (e *EmagCategoryStatApi) GetCategoryStatList(c *gin.Context) {
	var info emagReq.CategoryStatSearch
	err := c.ShouldBindQuery(&info)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	list, total, err := emagCategoryStatService.GetCategoryStatList(info)
	if err != nil {
		global.GVA_LOG.Error("获取品类指标列表失败!", zap.Error(err))
		response.FailWithMessage("获取品类指标列表失败: "+err.Error(), c)
		return
	}

	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     info.Page,
		PageSize: info.PageSize,
	}, "获取成功", c)
}

// GetCategoryStatGrowthRank 分页获取品类指标同比增长排名
// @Tags      EmagCategoryStat
// @Summary   分页获取品类指标同比增长排名
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  query     emagReq.CategoryStatGrowthSearch                         true  "查询参数"
// @Success   200   {object}  response.Response{data=map[string]interface{},msg=string}  "分页获取品类指标同比增长排名"
// @Router    /emagCategoryStat/getGrowthRank [get]
func (e *EmagCategoryStatApi) GetCategoryStatGrowthRank(c *gin.Context) {
	var info emagReq.CategoryStatGrowthSearch
	err := c.ShouldBindQuery(&info)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	list, total, currentDate, previousDate, err := emagCategoryStatService.GetCategoryStatGrowthRank(info)
	if err != nil {
		global.GVA_LOG.Error("获取品类指标同比增长排名失败!", zap.Error(err))
		response.FailWithMessage("获取品类指标同比增长排名失败: "+err.Error(), c)
		return
	}

	response.OkWithDetailed(map[string]interface{}{
		"list":         list,
		"total":        total,
		"page":         info.Page,
		"pageSize":     info.PageSize,
		"currentDate":  currentDate,
		"previousDate": previousDate,
	}, "获取成功", c)
}

// TriggerUpdateTask 手动触发更新品类统计任务
// @Tags      EmagCategoryStat
// @Summary   手动触发更新品类统计任务
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Success   200   {object}  response.Response{msg=string}  "触发成功"
// @Router    /emagCategoryStat/triggerUpdate [post]
func (e *EmagCategoryStatApi) TriggerUpdateTask(c *gin.Context) {
	// 异步执行任务，避免接口超时
	go func() {
		global.GVA_LOG.Info("手动触发品类统计更新任务")
		err := task.UpdateCategoryStat()
		if err != nil {
			global.GVA_LOG.Error("品类统计更新任务执行失败", zap.Error(err))
		} else {
			global.GVA_LOG.Info("品类统计更新任务执行完成")
		}
	}()

	response.OkWithMessage("任务已触发，请稍后刷新页面查看结果", c)
}
