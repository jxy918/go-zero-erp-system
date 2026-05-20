import request from './request'

export const supplierApi = {
  getSupplierList: (params) => {
    return request({
      url: '/supplier/list',
      method: 'get',
      params
    })
  },

  getActiveSupplierList: (params) => {
    return request({
      url: '/supplier/active-list',
      method: 'get',
      params
    })
  },

  createSupplier: (data) => {
    return request({
      url: '/supplier/create',
      method: 'post',
      data
    })
  },

  updateSupplier: (data) => {
    return request({
      url: '/supplier/update',
      method: 'post',
      data
    })
  },

  deleteSupplier: (data) => {
    return request({
      url: '/supplier/delete',
      method: 'post',
      data
    })
  },

  getSupplier: (id) => {
    return request({
      url: `/supplier/get/${id}`,
      method: 'get'
    })
  }
}