import { useUserStore } from '../store'
import { ElMessage } from 'element-plus'
import { PermissionCodes, hasPermission as checkPermission } from './permissionCodes'

/**
 * 权限检查工具类
 * 提供权限验证相关的方法，用于控制页面元素的显示和操作权限
 * 管理员用户拥有所有权限，普通用户需要匹配具体的权限code
 */
export const permission = {
  /**
   * 检查用户是否有指定权限
   * @param {string|string[]} permissionCodes - 权限code，可以是单个字符串或字符串数组
   * @returns {boolean} - 有权限返回true，无权限返回false
   */
  has(permissionCodes) {
    const userStore = useUserStore()
    
    // 管理员拥有所有权限，直接返回true
    if (userStore.isAdmin) {
      return true
    }
    
    // 获取用户的权限列表
    const userPermissions = userStore.permissions
    
    // 确保权限码是数组格式
    if (!Array.isArray(permissionCodes)) {
      permissionCodes = [permissionCodes]
    }
    
    // 检查是否有任一权限匹配（满足任一权限即可）
    return permissionCodes.some(code => checkPermission(code, userPermissions))
  },
  
  /**
   * 检查用户是否有指定权限，如果没有则显示警告提示
   * @param {string|string[]} permissionCodes - 权限code
   * @param {string} message - 无权限时的提示消息，默认"无权限操作"
   * @returns {boolean} - 有权限返回true，无权限返回false并显示提示
   */
  check(permissionCodes, message = '无权限操作') {
    if (this.has(permissionCodes)) {
      return true
    }
    ElMessage.warning(message)
    return false
  },
  
  /**
   * 检查用户是否有指定权限，如果没有则抛出错误
   * @param {string|string[]} permissionCodes - 权限code
   * @param {string} message - 无权限时的错误消息
   * @throws {Error} - 无权限时抛出错误
   */
  require(permissionCodes, message = '无权限操作') {
    if (!this.has(permissionCodes)) {
      throw new Error(message)
    }
  },
}

// 导出统一的权限码常量
export { PermissionCodes }

// 为保持向后兼容，导出别名 PermissionCode
export const PermissionCode = PermissionCodes

export default permission