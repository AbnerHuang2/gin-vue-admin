package emag

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	emagReq "github.com/flipped-aurora/gin-vue-admin/server/model/emag/request"
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

// GetCategoryStatTop20 获取品类指标Top20
// @Tags      EmagCategoryStat
// @Summary   获取品类指标Top20
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  query     emagReq.CategoryStatTop20Search                         true  "查询参数"
// @Success   200   {object}  response.Response{data=response.PageResult,msg=string}  "获取品类指标Top20"
// @Router    /emagCategoryStat/getTop20 [get]
func (e *EmagCategoryStatApi) GetCategoryStatTop20(c *gin.Context) {
	var info emagReq.CategoryStatTop20Search
	err := c.ShouldBindQuery(&info)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	list, total, err := emagCategoryStatService.GetCategoryStatTop20(info)
	if err != nil {
		global.GVA_LOG.Error("获取品类指标Top20失败!", zap.Error(err))
		response.FailWithMessage("获取品类指标Top20失败: "+err.Error(), c)
		return
	}

	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     info.Page,
		PageSize: 20,
	}, "获取成功", c)
}

// GetCategoryStatGrowthRank 获取品类指标同比增长排名
// @Tags      EmagCategoryStat
// @Summary   获取品类指标同比增长排名
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  query     emagReq.CategoryStatGrowthSearch                         true  "查询参数"
// @Success   200   {object}  response.Response{data=map[string]interface{},msg=string}  "获取品类指标同比增长排名"
// @Router    /emagCategoryStat/getGrowthRank [get]
func (e *EmagCategoryStatApi) GetCategoryStatGrowthRank(c *gin.Context) {
	var info emagReq.CategoryStatGrowthSearch
	err := c.ShouldBindQuery(&info)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	list, currentDate, previousDate, err := emagCategoryStatService.GetCategoryStatGrowthRank(info)
	if err != nil {
		global.GVA_LOG.Error("获取品类指标同比增长排名失败!", zap.Error(err))
		response.FailWithMessage("获取品类指标同比增长排名失败: "+err.Error(), c)
		return
	}

	response.OkWithDetailed(map[string]interface{}{
		"list":         list,
		"currentDate":  currentDate,
		"previousDate": previousDate,
		"total":        len(list),
	}, "获取成功", c)
}
