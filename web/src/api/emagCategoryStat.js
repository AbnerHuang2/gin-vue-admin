import service from '@/utils/request'

// 获取快照日期列表
export const getSnapshotDateList = () => {
  return service({
    url: '/emagCategoryStat/getSnapshotDateList',
    method: 'get'
  })
}

// 获取品类指标Top20
export const getCategoryStatTop20 = (params) => {
  return service({
    url: '/emagCategoryStat/getTop20',
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

