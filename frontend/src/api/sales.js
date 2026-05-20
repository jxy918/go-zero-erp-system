import request from './request'

export const salesApi = {
  getSalesOrderList: (params) => {
    return request({
      url: '/sales/order/list',
      method: 'get',
      params
    })
  },

  createSalesOrder: (data) => {
    return request({
      url: '/sales/order/create',
      method: 'post',
      data
    })
  },

  updateSalesOrder: (data) => {
    return request({
      url: '/sales/order/update',
      method: 'post',
      data
    })
  },

  deleteSalesOrder: (data) => {
    return request({
      url: '/sales/order/delete',
      method: 'post',
      data
    })
  },

  getSalesOrder: (id) => {
    return request({
      url: `/sales/order/get/${id}`,
      method: 'get'
    })
  },

  salesOutbound: (data) => {
    return request({
      url: '/inventory/sales-outbound',
      method: 'post',
      data
    })
  }
}