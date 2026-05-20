import request from './request'

// 权限API
export const permissionApi = {
  // 获取权限列表
  getPermissionList: (params) => {
    return request({
      url: '/permission/list',
      method: 'get',
      params
    })
  },

  // 获取权限详情
  getPermissionDetail: (id) => {
    return request({
      url: `/permission/get/${id}`,
      method: 'get'
    })
  },

  // 创建权限
  createPermission: (data) => {
    return request({
      url: '/permission/create',
      method: 'post',
      data
    })
  },

  // 更新权限
  updatePermission: (data) => {
    return request({
      url: '/permission/update',
      method: 'post',
      data
    })
  },

  // 删除权限
  deletePermission: (id) => {
    return request({
      url: '/permission/delete',
      method: 'post',
      data: { id }
    })
  }
}