import service from '@/utils/request'

// 分页获取订单列表
export const getOrderList = (params) => {
  return service({
    url: '/emagOrder/getList',
    method: 'get',
    params
  })
}

// 同步订单
export const syncOrders = (data) => {
  return service({
    url: '/emagOrder/sync',
    method: 'post',
    data
  })
}

// 获取国家列表
export const getCountryList = () => {
  return service({
    url: '/emagOrder/getCountryList',
    method: 'get'
  })
}

