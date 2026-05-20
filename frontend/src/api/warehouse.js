import request from './request'

export const warehouseApi = {
  getWarehouseList: (params) => {
    return request({
      url: '/warehouse/list',
      method: 'get',
      params
    })
  },

  getActiveWarehouseList: (params) => {
    return request({
      url: '/warehouse/active-list',
      method: 'get',
      params
    })
  },

  createWarehouse: (data) => {
    return request({
      url: '/warehouse/create',
      method: 'post',
      data
    })
  },

  updateWarehouse: (data) => {
    return request({
      url: '/warehouse/update',
      method: 'post',
      data
    })
  },

  deleteWarehouse: (data) => {
    return request({
      url: '/warehouse/delete',
      method: 'post',
      data
    })
  },

  getWarehouse: (id) => {
    return request({
      url: `/warehouse/get/${id}`,
      method: 'get'
    })
  }
}