import request from './request'

// 菜单API
export const menuApi = {
  // 获取菜单树
  getMenuTree: () => {
    return request({
      url: '/menu/tree',
      method: 'get'
    })
  },

  // 获取菜单列表
  getMenuList: () => {
    return request({
      url: '/menu/list',
      method: 'get'
    })
  },

  // 获取菜单详情
  getMenuDetail: (id) => {
    return request({
      url: `/menu/get/${id}`,
      method: 'get'
    })
  },

  // 创建菜单
  createMenu: (data) => {
    return request({
      url: '/menu/create',
      method: 'post',
      data
    })
  },

  // 更新菜单
  updateMenu: (data) => {
    return request({
      url: '/menu/update',
      method: 'post',
      data
    })
  },

  // 删除菜单
  deleteMenu: (id) => {
    return request({
      url: '/menu/delete',
      method: 'post',
      data: { id }
    })
  },

  // 为菜单分配权限
  assignPermissions: (menuId, permissionIds) => {
    return request({
      url: '/menu/assign-permissions',
      method: 'post',
      data: {
        menu_id: menuId,
        permission_ids: permissionIds
      }
    })
  }
}