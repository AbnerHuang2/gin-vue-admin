package emag

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/model/emag"
	emagReq "github.com/flipped-aurora/gin-vue-admin/server/model/emag/request"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type EmagCategoryApi struct{}

// CreateEmagCategory 创建品类
// @Tags      EmagCategory
// @Summary   创建品类
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      emag.EmagCategory              true  "品类信息"
// @Success   200   {object}  response.Response{msg=string}  "创建品类"
// @Router    /emagCategory/createEmagCategory [post]
func (e *EmagCategoryApi) CreateEmagCategory(c *gin.Context) {
	var category emag.EmagCategory
	err := c.ShouldBindJSON(&category)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if category.CategoryId == "" || category.CategoryName == "" {
		response.FailWithMessage("品类ID和品类名称不能为空", c)
		return
	}
	err = emagCategoryService.CreateEmagCategory(category)
	if err != nil {
		global.GVA_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败: "+err.Error(), c)
		return
	}
	response.OkWithMessage("创建成功", c)
}

// DeleteEmagCategory 删除品类
// @Tags      EmagCategory
// @Summary   删除品类
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      emag.EmagCategory              true  "品类ID"
// @Success   200   {object}  response.Response{msg=string}  "删除品类"
// @Router    /emagCategory/deleteEmagCategory [delete]
func (e *EmagCategoryApi) DeleteEmagCategory(c *gin.Context) {
	var category emag.EmagCategory
	err := c.ShouldBindJSON(&category)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if category.CategoryId == "" {
		response.FailWithMessage("品类ID不能为空", c)
		return
	}
	err = emagCategoryService.DeleteEmagCategory(category.CategoryId)
	if err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败: "+err.Error(), c)
		return
	}
	response.OkWithMessage("删除成功", c)
}

// DeleteEmagCategoryByIds 批量删除品类
// @Tags      EmagCategory
// @Summary   批量删除品类
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      request.IdsReq                 true  "品类ID列表"
// @Success   200   {object}  response.Response{msg=string}  "批量删除品类"
// @Router    /emagCategory/deleteEmagCategoryByIds [delete]
func (e *EmagCategoryApi) DeleteEmagCategoryByIds(c *gin.Context) {
	var IDS struct {
		Ids []uint `json:"ids"`
	}
	err := c.ShouldBindJSON(&IDS)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = emagCategoryService.DeleteEmagCategoryByIds(IDS.Ids)
	if err != nil {
		global.GVA_LOG.Error("批量删除失败!", zap.Error(err))
		response.FailWithMessage("批量删除失败: "+err.Error(), c)
		return
	}
	response.OkWithMessage("批量删除成功", c)
}

// UpdateEmagCategory 更新品类
// @Tags      EmagCategory
// @Summary   更新品类
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      emag.EmagCategory              true  "品类信息"
// @Success   200   {object}  response.Response{msg=string}  "更新品类"
// @Router    /emagCategory/updateEmagCategory [put]
func (e *EmagCategoryApi) UpdateEmagCategory(c *gin.Context) {
	var category emag.EmagCategory
	err := c.ShouldBindJSON(&category)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if category.CategoryId == "" {
		response.FailWithMessage("品类ID不能为空", c)
		return
	}
	err = emagCategoryService.UpdateEmagCategory(category)
	if err != nil {
		global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败: "+err.Error(), c)
		return
	}
	response.OkWithMessage("更新成功", c)
}

// FindEmagCategory 获取单个品类信息
// @Tags      EmagCategory
// @Summary   获取单个品类信息
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  query     emagReq.EmagCategorySearch                         true  "查询参数"
// @Success   200   {object}  response.Response{data=emag.EmagCategory,msg=string}  "获取单个品类信息"
// @Router    /emagCategory/findEmagCategory [get]
func (e *EmagCategoryApi) FindEmagCategory(c *gin.Context) {
	var info emagReq.EmagCategorySearch
	err := c.ShouldBindQuery(&info)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	category, err := emagCategoryService.GetEmagCategory(info)
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败: "+err.Error(), c)
		return
	}
	response.OkWithDetailed(category, "获取成功", c)
}

// GetEmagCategoryList 分页获取品类列表
// @Tags      EmagCategory
// @Summary   分页获取品类列表
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  query     emagReq.EmagCategorySearch                         true  "分页参数"
// @Success   200   {object}  response.Response{data=response.PageResult,msg=string}  "分页获取品类列表"
// @Router    /emagCategory/getEmagCategoryList [get]
func (e *EmagCategoryApi) GetEmagCategoryList(c *gin.Context) {
	var pageInfo emagReq.EmagCategorySearch
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := emagCategoryService.GetEmagCategoryList(pageInfo)
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败: "+err.Error(), c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
	}, "获取成功", c)
}
