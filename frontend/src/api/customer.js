import request from './request'

export const customerApi = {
  getCustomerList: (params) => {
    return request({
      url: '/customer/list',
      method: 'get',
      params
    })
  },

  getActiveCustomerList: (params) => {
    return request({
      url: '/customer/active-list',
      method: 'get',
      params
    })
  },

  createCustomer: (data) => {
    return request({
      url: '/customer/create',
      method: 'post',
      data
    })
  },

  updateCustomer: (data) => {
    return request({
      url: '/customer/update',
      method: 'post',
      data
    })
  },

  deleteCustomer: (data) => {
    return request({
      url: '/customer/delete',
      method: 'post',
      data
    })
  },

  getCustomer: (id) => {
    return request({
      url: `/customer/get/${id}`,
      method: 'get'
    })
  }
}