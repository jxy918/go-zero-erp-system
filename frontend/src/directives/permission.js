import { useUserStore } from '../store'
import { watch } from 'vue'

/**
 * v-has-permission 自定义指令
 * 用于控制元素的显示/隐藏，基于用户权限
 * 
 * 使用方式：
 * <button v-has-permission="'btn_user_create'">创建用户</button>
 * 
 * 工作原理：
 * 1. 组件挂载时检查用户是否有指定权限
 * 2. 管理员拥有所有权限，元素始终显示
 * 3. 普通用户需要匹配具体的权限code
 * 4. 监听权限变化，动态更新元素显示状态
 */
export const hasPermission = {
  /**
   * 指令挂载时的钩子函数
   * @param {HTMLElement} el - 绑定指令的DOM元素
   * @param {Object} binding - 指令绑定信息
   * @param {string} binding.value - 权限code
   */
  mounted(el, binding) {
    const userStore = useUserStore()
    const permissionCode = binding.value
    
    /**
     * 检查权限并更新元素显示状态
     */
    const checkAndUpdate = () => {
      // userStore不可用时隐藏元素
      if (!userStore) {
        el.style.display = 'none'
        return
      }
      
      // 管理员拥有所有权限，直接显示
      if (userStore.isAdmin) {
        el.style.display = ''
        return
      }
      
      // 获取用户权限code列表（兼容大小写字段名）
      const permissions = userStore.permissions || []
      const permissionCodes = permissions.map(p => p.code || p.Code)
      
      // 检查是否有权限
      const hasPermission = permissionCodes.includes(permissionCode)
      if (hasPermission) {
        el.style.display = ''
        el.classList.remove('permission-hidden')
      } else {
        el.style.display = 'none'
        el.classList.add('permission-hidden')
      }
    }
    
    // 立即执行权限检查
    checkAndUpdate()
    
    // 监听权限变化，动态更新
    const unwatch = watch(
      () => userStore.permissions,
      () => {
        checkAndUpdate()
      },
      { deep: true }
    )
    
    // 保存watch停止函数，用于组件卸载时清理
    el._permissionWatch = unwatch
  },
  
  /**
   * 指令更新时的钩子函数
   * @param {HTMLElement} el - 绑定指令的DOM元素
   * @param {Object} binding - 指令绑定信息
   */
  updated(el, binding) {
    const userStore = useUserStore()
    const permissionCode = binding.value
    
    if (!userStore) {
      el.style.display = 'none'
      return
    }
    
    // 管理员显示所有元素
    if (userStore.isAdmin) {
      el.style.display = ''
      return
    }
    
    // 检查权限并更新显示状态
    const permissions = userStore.permissions || []
    const codes = permissions.map(p => p.code || p.Code)
    if (codes.includes(permissionCode)) {
      el.style.display = ''
      el.classList.remove('permission-hidden')
    } else {
      el.style.display = 'none'
      el.classList.add('permission-hidden')
    }
  },
  
  /**
   * 指令卸载时的钩子函数
   * 清理watch监听器，防止内存泄漏
   * @param {HTMLElement} el - 绑定指令的DOM元素
   */
  unmounted(el) {
    if (el._permissionWatch && typeof el._permissionWatch === 'function') {
      el._permissionWatch()
    }
  }
}

export default hasPermission