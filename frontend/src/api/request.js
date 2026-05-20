import axios from 'axios'
import { ElMessage } from 'element-plus'
import { useUserStore } from '../store'

// 创建axios实例
// baseURL 使用 '/api' 通过Vite代理转发到后端服务
const service = axios.create({
  baseURL: '/api', // 基础URL，通过Vite代理转发到后端
  timeout: 10000,  // 请求超时时间（10秒）
  headers: {
    'Content-Type': 'application/json; charset=utf-8'
  }
})

/**
 * 请求拦截器
 * 负责在发送请求前添加认证信息（token、用户ID）
 * 注意：拦截器中不能使用Pinia store（因为不在Vue组件上下文中），直接从localStorage读取
 */
service.interceptors.request.use(
  config => {
    // 从localStorage获取JWT token
    const token = localStorage.getItem('token')
    
    // 添加Authorization请求头（Bearer token格式）
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    
    // 添加用户ID到请求头（用于后端权限验证）
    const userStr = localStorage.getItem('user')
    if (userStr) {
      try {
        const userInfo = JSON.parse(userStr)
        const userId = userInfo.id || userInfo.ID // 兼容大小写
        if (userId) {
          config.headers['user_id'] = userId.toString()
        }
      } catch (e) {
        console.error('[Request Interceptor] 解析用户信息失败:', e)
      }
    }
    
    return config
  },
  error => {
    console.error('[Request Interceptor] 请求错误:', error)
    return Promise.reject(error)
  }
)

/**
 * 响应拦截器
 * 负责统一处理后端响应：
 * 1. 统一响应格式
 * 2. 处理错误状态码
 * 3. 处理token过期（401）
 */
service.interceptors.response.use(
  response => {
    const res = response.data
    
    // 兼容后端返回的不同字段名格式（code/Code, message/Message）
    const code = res.code || res.Code || 0
    const message = res.message || res.Message || 'success'
    
    // 判断code是否为有效的数字状态码
    // 如果code是字符串且不是"0"或"200"，说明不是状态码而是业务字段（如分类编码）
    const isStatusCode = typeof code === 'number' || 
                        (typeof code === 'string' && (code === '0' || code === '200'))
    
    // 处理业务错误（只有当code是有效状态码时才判断）
    const hasError = isStatusCode && code !== 0 && code !== 200
    
    if (hasError) {
      ElMessage.error(message || '请求失败')
      
      // 处理未授权/Token过期（401）
      if (code === 401) {
        const userStore = useUserStore()
        userStore.logout()
        window.location.href = '/login' // 强制跳转到登录页
      }
      
      return Promise.reject(new Error(message || '请求失败'))
    }
    
    // 统一响应数据格式
    let responseData = res
    
    // 情况1：后端返回标准格式 {code, data: {...}, message}
    if (res.data !== undefined) {
      responseData = res.data
    } else if (res.Data !== undefined) {
      responseData = res.Data // 兼容大写Data
    }
    // 情况2：后端直接返回业务数据（没有外层data包装）- 包含menus/categories等字段
    else if (res.menus !== undefined || res.Menus !== undefined || 
             res.categories !== undefined || res.Categories !== undefined ||
             res.products !== undefined || res.Products !== undefined ||
             res.total !== undefined || res.Total !== undefined) {
      responseData = res
    }
    // 情况3：go-zero httpx.ErrorCtx返回的格式（只有code和message，没有data）
    else if (isStatusCode && res.data === undefined && res.Data === undefined) {
      // 这种情况是错误响应，应该已经在上面的hasError中处理了
      responseData = null
    }
    
    // 仅在开发模式下输出调试日志
    if (import.meta.env.DEV) {
      console.log('[Response Interceptor] 原始响应:', JSON.stringify(res, null, 2))
      console.log('[Response Interceptor] 处理后数据:', JSON.stringify(responseData, null, 2))
    }
    
    // 返回统一格式的响应对象
    return {
      code: code,
      data: responseData,
      message: message
    }
  },
  error => {
    console.error('[Response Interceptor] 响应错误:', error)
    
    let errorMessage = '网络错误'
    if (error.response) {
      // 根据HTTP状态码处理不同错误
      switch (error.response.status) {
        case 401:
          errorMessage = '登录已失效，请重新登录'
          const userStore = useUserStore()
          userStore.logout()
          window.location.href = '/login'
          break
        case 403:
          // 获取后端返回的具体错误信息，如果没有则使用默认提示
          let detailMessage = ''
          const responseData = error.response.data
          if (typeof responseData === 'string') {
            // 如果返回的是纯文本
            detailMessage = responseData
          } else if (responseData && (responseData.message || responseData.Message)) {
            // 如果返回的是JSON对象
            detailMessage = responseData.message || responseData.Message
          } else if (responseData && responseData.error) {
            // 如果返回的是包含error字段的对象
            detailMessage = responseData.error
          }
          // 如果后端返回了具体的权限相关错误信息，使用它；否则使用默认提示
          if (detailMessage && (detailMessage.includes('权限') || detailMessage.includes('forbidden') || detailMessage.includes('Forbidden'))) {
            // 清理可能的HTML标签
            errorMessage = detailMessage.replace(/<[^>]*>/g, '').trim()
            // 如果清理后是空的或者只是"forbidden"，使用友好提示
            if (!errorMessage || errorMessage.toLowerCase() === 'forbidden') {
              errorMessage = '当前功能您暂无权限访问，请联系管理员'
            }
          } else {
            errorMessage = '当前功能您暂无权限访问，请联系管理员'
          }
          break
        case 404:
          errorMessage = '访问的页面或接口不存在'
          break
        case 500:
          errorMessage = '服务器繁忙，请稍后重试'
          break
        default:
          errorMessage = `操作失败 (${error.response.status})`
      }
    } else if (error.message.includes('timeout')) {
      errorMessage = '请求超时，请稍后重试'
    }
    
    ElMessage.error(errorMessage)
    return Promise.reject(error)
  }
)

export default service
