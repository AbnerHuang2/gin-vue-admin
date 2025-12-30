import service from '@/utils/request'

// 获取快照日期列表
export const getSnapshotDateList = () => {
  return service({
    url: '/emagCategoryStat/getSnapshotDateList',
    method: 'get'
  })
}

// 分页获取品类指标列表
export const getCategoryStatList = (params) => {
  return service({
    url: '/emagCategoryStat/getList',
    method: 'get',
    params
  })
}

// 获取品类指标同比增长排名
export const getCategoryStatGrowthRank = (params) => {
  return service({
    url: '/emagCategoryStat/getGrowthRank',
    method: 'get',
    params
  })
}

// 手动触发更新品类统计任务
export const triggerUpdateTask = () => {
  return service({
    url: '/emagCategoryStat/triggerUpdate',
    method: 'post'
  })
}

