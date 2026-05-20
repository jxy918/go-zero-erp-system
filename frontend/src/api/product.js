import request from './request'

export const productApi = {
  getProductList: (params) => {
    return request({
      url: '/product/list',
      method: 'get',
      params
    })
  },

  getActiveProductList: (params) => {
    return request({
      url: '/product/active-list',
      method: 'get',
      params
    })
  },

  createProduct: (data) => {
    return request({
      url: '/product/create',
      method: 'post',
      data
    })
  },

  updateProduct: (data) => {
    return request({
      url: '/product/update',
      method: 'post',
      data
    })
  },

  deleteProduct: (data) => {
    return request({
      url: '/product/delete',
      method: 'post',
      data
    })
  },

  getProduct: (id) => {
    return request({
      url: `/product/get/${id}`,
      method: 'get'
    })
  },

  getCategoryList: (params) => {
    return request({
      url: '/product/category/list',
      method: 'get',
      params
    })
  },

  createCategory: (data) => {
    return request({
      url: '/product/category/create',
      method: 'post',
      data
    })
  },

  updateCategory: (data) => {
    return request({
      url: '/product/category/update',
      method: 'post',
      data
    })
  },

  deleteCategory: (data) => {
    return request({
      url: '/product/category/delete',
      method: 'post',
      data
    })
  },

  listProductUnit: (params) => {
    return request({
      url: '/product/unit/list',
      method: 'get',
      params
    })
  },

  createProductUnit: (data) => {
    return request({
      url: '/product/unit/create',
      method: 'post',
      data
    })
  },

  updateProductUnit: (data) => {
    return request({
      url: '/product/unit/update',
      method: 'post',
      data
    })
  },

  deleteProductUnit: (data) => {
    return request({
      url: '/product/unit/delete',
      method: 'post',
      data
    })
  },

  getProductUnit: (id) => {
    return request({
      url: `/product/unit/get/${id}`,
      method: 'get'
    })
  }
}