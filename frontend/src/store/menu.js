import { defineStore } from 'pinia'
import { menuApi } from '../api'

// 标准化菜单响应数据
function normalizeMenuResponse(response) {
  if (!response) {
    return []
  }
  
  // 响应格式: {code, data, message}
  const responseData = response.data
  if (!responseData) {
    return []
  }
  
  // 后端返回的数据格式: {Menus: [...], Total: N}
  // 注意：Go返回的字段名是大写开头的
  const menus = responseData.Menus || responseData.menus
  
  if (!menus) {
    return []
  }
  
  if (!Array.isArray(menus)) {
    return []
  }
  
  return menus
}

// 菜单状态管理
export const useMenuStore = defineStore('menu', {
  state: () => ({
    menuTree: [],
    menuList: [],
    loading: false,
    error: null,
  }),
  getters: {
    getMenuTree: (state) => state.menuTree,
    getMenuList: (state) => state.menuList,
    getLoading: (state) => state.loading,
    getError: (state) => state.error,
  },
  actions: {
    // 清除菜单缓存
    clearMenuCache() {
      localStorage.removeItem('menuTree')
      localStorage.removeItem('menuList')
    },
    
    // 加载菜单树（优先使用缓存）
    async loadMenuTree(forceRefresh = false) {
      // 如果不是强制刷新且有缓存，则使用缓存
      if (!forceRefresh) {
        const cachedTree = localStorage.getItem('menuTree')
        if (cachedTree) {
          try {
            this.menuTree = JSON.parse(cachedTree)
            this.generateMenuList(this.menuTree)
            return
          } catch (e) {
            console.error('解析缓存菜单树失败:', e)
          }
        }
      }
      
      this.loading = true
      this.error = null
      try {
        // 清除旧缓存
        this.clearMenuCache()
        
        const response = await menuApi.getMenuTree()
        
        // 使用标准化函数解析数据
        const data = normalizeMenuResponse(response)
        
        this.menuTree = data
        this.generateMenuList(this.menuTree)
        localStorage.setItem('menuTree', JSON.stringify(this.menuTree))
      } catch (error) {
        this.error = '获取菜单树失败'
        console.error('获取菜单树失败:', error)
        this.menuTree = []
        this.generateMenuList(this.menuTree)
      } finally {
        this.loading = false
      }
    },
    // 生成菜单列表
    generateMenuList(menuTree) {
      const list = []
      
      const traverse = (menu) => {
        list.push(menu)
        if (menu.children && menu.children.length > 0) {
          menu.children.forEach(child => {
            traverse(child)
          })
        }
      }
      
      menuTree.forEach(menu => {
        traverse(menu)
      })
      
      this.menuList = list
    },
    // 创建菜单
    async createMenu(menuData) {
      this.loading = true
      this.error = null
      try {
        await menuApi.createMenu(menuData)
        await this.loadMenuTree(true)
        return true
      } catch (error) {
        this.error = '创建菜单失败'
        console.error('创建菜单失败:', error)
        return false
      } finally {
        this.loading = false
      }
    },
    // 更新菜单
    async updateMenu(menuData) {
      this.loading = true
      this.error = null
      try {
        await menuApi.updateMenu(menuData)
        await this.loadMenuTree(true)
        return true
      } catch (error) {
        this.error = '更新菜单失败'
        console.error('更新菜单失败:', error)
        return false
      } finally {
        this.loading = false
      }
    },
    // 删除菜单
    async deleteMenu(menuId) {
      this.loading = true
      this.error = null
      try {
        await menuApi.deleteMenu(menuId)
        await this.loadMenuTree(true)
        return true
      } catch (error) {
        this.error = '删除菜单失败'
        console.error('删除菜单失败:', error)
        return false
      } finally {
        this.loading = false
      }
    },
    // 为菜单分配权限
    async assignPermissions(menuId, permissionIds) {
      this.loading = true
      this.error = null
      try {
        await menuApi.assignPermissions(menuId, permissionIds)
        return true
      } catch (error) {
        this.error = '分配权限失败'
        console.error('分配权限失败:', error)
        return false
      } finally {
        this.loading = false
      }
    },
  },
})
