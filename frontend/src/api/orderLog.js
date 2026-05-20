import request from './request'

export const orderLogApi = {
  getOrderLogList: (params) => {
    return request({
      url: '/inventory/order-log/list',
      method: 'get',
      params
    })
  }
}