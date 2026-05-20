import { defineStore } from 'pinia'

// 用户状态管理
export const useUserStore = defineStore('user', {
  state: () => ({
    userInfo: null,
    token: localStorage.getItem('token') || '',
    isAuthenticated: !!localStorage.getItem('token'),
    permissions: [], // 用户权限列表
    permissionVersion: 0, // 添加版本号用于强制刷新
  }),
  getters: {
    getUserInfo: (state) => state.userInfo,
    getToken: (state) => state.token,
    getIsAuthenticated: (state) => state.isAuthenticated,
    getPermissions: (state) => state.permissions,
    // 判断是否为管理员
    isAdmin: (state) => {
      if (!state.userInfo) {
        return false
      }
      // 兼容不同的数据结构
      const roles = state.userInfo.roles || state.userInfo.Roles || []
      if (!Array.isArray(roles) || roles.length === 0) {
        return false
      }
      return roles.some(role => {
        const code = role.code || role.Code || ''
        return code.toLowerCase() === 'admin'
      })
    },
    // 获取权限code列表
    getPermissionCodes: (state) => {
      const codes = state.permissions.map(p => p.code || p.Code)
      console.log('[userStore] getPermissionCodes called, version:', state.permissionVersion, 'codes:', codes)
      return codes
    },
  },
  actions: {
    // 设置用户信息
    setUserInfo(userInfo) {
      this.userInfo = userInfo
      localStorage.setItem('user', JSON.stringify(userInfo))
    },
    // 设置token
    setToken(token) {
      this.token = token
      this.isAuthenticated = !!token
      if (token) {
        localStorage.setItem('token', token)
      } else {
        localStorage.removeItem('token')
      }
    },
    // 登录
    login(userInfo, token) {
      this.setUserInfo(userInfo)
      this.setToken(token)
      // 从角色中提取权限
      this.extractPermissions()
    },
    
    // 从角色中提取权限
    extractPermissions() {
      this.permissions = []
      console.log('[extractPermissions] userInfo:', JSON.stringify(this.userInfo))
      
      // 空值检查：如果 userInfo 为 null/undefined，直接返回
      if (!this.userInfo) {
        console.log('[extractPermissions] userInfo is null or undefined')
        this.permissionVersion++
        return
      }
      
      // 兼容大小写的 roles 字段
      const roles = this.userInfo.roles || this.userInfo.Roles || []
      console.log('[extractPermissions] roles found:', roles.length)
      if (roles.length > 0) {
        console.log('[extractPermissions] roles count:', roles.length)
        roles.forEach((role, idx) => {
          console.log('[extractPermissions] processing role', idx, ':', role.name || role.Name)
          // 兼容大小写字段名（后端返回可能是 Permissions，前端可能是 permissions）
          // 确保转换为普通数组，避免 Proxy 对象问题
          const rawPermissions = role.permissions || role.Permissions || []
          const rolePermissions = Array.isArray(rawPermissions) ? rawPermissions : Object.values(rawPermissions)
          console.log('[extractPermissions] rolePermissions:', rolePermissions, 'length:', rolePermissions.length)
          
          // 遍历所有权限
          for (const perm of rolePermissions) {
            console.log('[extractPermissions] checking permission:', perm.code || perm.Code)
            // 避免重复添加
            const permId = perm.id || perm.ID
            const permCode = perm.code || perm.Code
            const exists = this.permissions.some(p => (p.id || p.ID) === permId)
            if (!exists && permCode) {
              this.permissions.push(perm)
              console.log('[extractPermissions] added permission:', permCode)
            }
          }
        })
      } else {
        console.log('[extractPermissions] No roles found in userInfo')
      }
      // 递增版本号以强制刷新所有依赖权限的计算属性
      this.permissionVersion++
      console.log('[extractPermissions] Final permissions count:', this.permissions.length)
      console.log('[extractPermissions] Final permissionCodes:', this.permissions.map(p => p.code || p.Code))
      console.log('[extractPermissions] permissionVersion incremented to:', this.permissionVersion)
    },
    // 登出
    logout() {
      this.userInfo = null
      this.setToken('')
      localStorage.removeItem('user')
    },
    // 从localStorage加载用户信息
    loadUserInfo() {
      const userStr = localStorage.getItem('user')
      if (userStr) {
        try {
          this.userInfo = JSON.parse(userStr)
        } catch (e) {
          console.error('解析用户信息失败:', e)
          this.userInfo = null
        }
      }
      
      // 从localStorage加载token
      const token = localStorage.getItem('token')
      this.setToken(token)
      
      // 从用户信息中提取权限
      this.extractPermissions()
    },
  },
})
