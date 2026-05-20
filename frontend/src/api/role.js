import request from './request'

// 角色API
export const roleApi = {
  // 获取角色列表
  getRoleList: (params) => {
    return request({
      url: '/role/list',
      method: 'get',
      params
    })
  },

  // 获取角色详情
  getRoleDetail: (id) => {
    return request({
      url: `/role/get/${id}`,
      method: 'get'
    })
  },

  // 创建角色
  createRole: (data) => {
    return request({
      url: '/role/create',
      method: 'post',
      data
    })
  },

  // 更新角色
  updateRole: (data) => {
    return request({
      url: '/role/update',
      method: 'post',
      data
    })
  },

  // 删除角色
  deleteRole: (id) => {
    return request({
      url: '/role/delete',
      method: 'post',
      data: { id: id }
    })
  },

  // 为角色分配权限
  assignPermissions: (roleId, permissionIds) => {
    return request({
      url: '/role/assign-permissions',
      method: 'post',
      data: {
        role_id: roleId,
        permission_ids: permissionIds
      }
    })
  },

  // 为角色分配菜单
  assignMenus: (roleId, menuIds) => {
    return request({
      url: '/role/assign-menus',
      method: 'post',
      data: {
        role_id: roleId,
        menu_ids: menuIds
      }
    })
  }
}