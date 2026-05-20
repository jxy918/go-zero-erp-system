import request from './request'

export const erpApi = {
  getOverview: () => {
    return request({
      url: '/erp/statistics/overview',
      method: 'get'
    })
  },

  getTrend: (params) => {
    return request({
      url: '/erp/statistics/trend',
      method: 'get',
      params
    })
  },

  getInventoryAlert: () => {
    return request({
      url: '/erp/statistics/inventory-alert',
      method: 'get'
    })
  },

  getTopProducts: () => {
    return request({
      url: '/erp/statistics/top-products',
      method: 'get'
    })
  },

  getOrderStatus: (params) => {
    return request({
      url: '/erp/statistics/order-status',
      method: 'get',
      params
    })
  },

  getBusiness: () => {
    return request({
      url: '/erp/statistics/business',
      method: 'get'
    })
  },

  getTodoData: () => {
    return request({
      url: '/erp/statistics/todo-data',
      method: 'get'
    })
  }
}
