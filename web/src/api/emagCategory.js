import service from '@/utils/request'

// 创建品类
export const createEmagCategory = (data) => {
  return service({
    url: '/emagCategory/createEmagCategory',
    method: 'post',
    data
  })
}

// 删除品类
export const deleteEmagCategory = (data) => {
  return service({
    url: '/emagCategory/deleteEmagCategory',
    method: 'delete',
    data
  })
}

// 批量删除品类
export const deleteEmagCategoryByIds = (data) => {
  return service({
    url: '/emagCategory/deleteEmagCategoryByIds',
    method: 'delete',
    data
  })
}

// 更新品类
export const updateEmagCategory = (data) => {
  return service({
    url: '/emagCategory/updateEmagCategory',
    method: 'put',
    data
  })
}

// 获取单个品类
export const findEmagCategory = (params) => {
  return service({
    url: '/emagCategory/findEmagCategory',
    method: 'get',
    params
  })
}

// 获取品类列表
export const getEmagCategoryList = (params) => {
  return service({
    url: '/emagCategory/getEmagCategoryList',
    method: 'get',
    params
  })
}

