import request from './request'

export const inventoryCheckApi = {
  createInventoryCheck: (data) => {
    return request({
      url: '/inventory/check/create',
      method: 'post',
      data
    })
  },

  getInventoryCheck: (id) => {
    return request({
      url: `/inventory/check/get/${id}`,
      method: 'get'
    })
  },

  getInventoryCheckList: (params) => {
    return request({
      url: '/inventory/check/list',
      method: 'get',
      params
    })
  },

  updateInventoryCheck: (data) => {
    return request({
      url: '/inventory/check/update',
      method: 'post',
      data
    })
  },

  deleteInventoryCheck: (data) => {
    return request({
      url: '/inventory/check/delete',
      method: 'post',
      data
    })
  },

  submitInventoryCheck: (data) => {
    return request({
      url: '/inventory/check/submit',
      method: 'post',
      data
    })
  },

  generateInventoryAdjust: (data) => {
    return request({
      url: '/inventory/generate-adjust-from-check',
      method: 'post',
      data
    })
  }
}