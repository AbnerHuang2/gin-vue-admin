import service from '@/utils/request'

// 分页获取产品列表
export const getProductList = (params) => {
  return service({
    url: '/emagProduct/getList',
    method: 'get',
    params
  })
}

// 同步产品
export const syncProducts = (data) => {
  return service({
    url: '/emagProduct/sync',
    method: 'post',
    data
  })
}

// 获取状态列表
export const getStatusList = () => {
  return service({
    url: '/emagProduct/getStatusList',
    method: 'get'
  })
}

// 获取国家列表
export const getCountryList = () => {
  return service({
    url: '/emagProduct/getCountryList',
    method: 'get'
  })
}
