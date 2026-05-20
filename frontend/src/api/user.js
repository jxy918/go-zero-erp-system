import request from './request'

// 用户API
export const userApi = {
  // 获取用户列表
  getUserList: (params) => {
    return request({
      url: '/user/list',
      method: 'get',
      params
    })
  },

  // 获取用户详情
  getUserDetail: (id) => {
    return request({
      url: `/user/get/${id}`,
      method: 'get'
    })
  },

  // 创建用户
  createUser: (data) => {
    return request({
      url: '/user/create',
      method: 'post',
      data
    })
  },

  // 更新用户
  updateUser: (data) => {
    return request({
      url: '/user/update',
      method: 'post',
      data
    })
  },

  // 删除用户
  deleteUser: (id) => {
    return request({
      url: '/user/delete',
      method: 'post',
      data: { id }
    })
  },

  // 为用户分配角色
  assignRoles: (userId, roleIds) => {
    return request({
      url: '/user/assign-roles',
      method: 'post',
      data: {
        user_id: userId,
        role_ids: roleIds
      }
    })
  }
}