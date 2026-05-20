import request from './request'

export const inventoryTransferApi = {
  createInventoryTransfer: (data) => {
    return request({
      url: '/inventory/transfer/create',
      method: 'post',
      data
    })
  },

  getInventoryTransfer: (id) => {
    return request({
      url: `/inventory/transfer/get/${id}`,
      method: 'get'
    })
  },

  getInventoryTransferList: (params) => {
    return request({
      url: '/inventory/transfer/list',
      method: 'get',
      params
    })
  },

  updateInventoryTransfer: (data) => {
    return request({
      url: '/inventory/transfer/update',
      method: 'post',
      data
    })
  },

  deleteInventoryTransfer: (data) => {
    return request({
      url: '/inventory/transfer/delete',
      method: 'post',
      data
    })
  },

  auditInventoryTransfer: (data) => {
    return request({
      url: '/inventory/transfer/audit',
      method: 'post',
      data
    })
  },

  executeInventoryTransfer: (data) => {
    return request({
      url: '/inventory/transfer/execute',
      method: 'post',
      data
    })
  }
}