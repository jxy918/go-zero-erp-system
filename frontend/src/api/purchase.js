import request from './request'

export const purchaseApi = {
  getPurchaseOrderList: (params) => {
    return request({
      url: '/purchase/order/list',
      method: 'get',
      params
    })
  },

  createPurchaseOrder: (data) => {
    return request({
      url: '/purchase/order/create',
      method: 'post',
      data
    })
  },

  updatePurchaseOrder: (data) => {
    return request({
      url: '/purchase/order/update',
      method: 'post',
      data
    })
  },

  deletePurchaseOrder: (data) => {
    return request({
      url: '/purchase/order/delete',
      method: 'post',
      data
    })
  },

  getPurchaseOrder: (id) => {
    return request({
      url: `/purchase/order/get/${id}`,
      method: 'get'
    })
  },

  purchaseInbound: (data) => {
    return request({
      url: '/purchase/order/inbound',
      method: 'post',
      data
    })
  }
}