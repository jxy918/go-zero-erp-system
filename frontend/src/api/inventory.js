import request from './request'

export const inventoryApi = {
  getInventoryList: (params) => {
    return request({
      url: '/inventory/list',
      method: 'get',
      params
    })
  },

  adjustInventory: (data) => {
    return request({
      url: '/inventory/adjust',
      method: 'post',
      data
    })
  },

  getInventoryHistory: (params) => {
    return request({
      url: '/inventory/history',
      method: 'get',
      params
    })
  },

  getCurrentStock: (params) => {
    return request({
      url: '/inventory/current-stock',
      method: 'get',
      params
    })
  },

  createInventoryAdjustRequest: (data) => {
    return request({
      url: '/inventory/adjust/create',
      method: 'post',
      data
    })
  },

  getInventoryAdjustRequestList: (params) => {
    return request({
      url: '/inventory/adjust/list',
      method: 'get',
      params
    })
  },

  approveInventoryAdjustRequest: (data) => {
    return request({
      url: '/inventory/adjust/approve',
      method: 'post',
      data
    })
  },

  rejectInventoryAdjustRequest: (data) => {
    return request({
      url: '/inventory/adjust/reject',
      method: 'post',
      data
    })
  },

  listInventoryAlert: (params) => {
    return request({
      url: '/inventory/alert/list',
      method: 'get',
      params
    })
  },

  checkInventoryAlert: () => {
    return request({
      url: '/inventory/alert/check',
      method: 'get'
    })
  }
}