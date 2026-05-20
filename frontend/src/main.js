import { createApp } from 'vue'
import { createPinia } from 'pinia'
import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'
import App from './App.vue'
import router from './router'
import axios from 'axios'
import { useUserStore } from './store'
import { useThemeStore } from './store/theme'
import { hasPermission } from './directives/permission'

// 配置axios全局拦截器
axios.interceptors.request.use(config => {
  // 从localStorage获取token
  const token = localStorage.getItem('token')
  if (token) {
    // 添加Authorization请求头
    config.headers.Authorization = `Bearer ${token}`
  }
  // 从localStorage获取用户信息，添加user_id请求头
  const userStr = localStorage.getItem('user')
  if (userStr) {
    try {
      const user = JSON.parse(userStr)
      if (user.id) {
        config.headers['user_id'] = user.id.toString()
      }
    } catch (e) {
      console.error('解析用户信息失败:', e)
    }
  }
  // 确保设置Content-Type为application/json; charset=utf-8
  config.headers['Content-Type'] = 'application/json; charset=utf-8'
  return config
}, error => {
  return Promise.reject(error)
})

const app = createApp(App)
const pinia = createPinia()

app.use(pinia)
app.use(router)
app.use(ElementPlus)

// 注册权限指令
app.directive('has-permission', hasPermission)


// 全局错误处理
app.config.errorHandler = (err, instance, info) => {
  console.error('[Vue Error]', err)
  console.error('[Vue Error] Component:', instance?.$options?.name || 'Unknown')
  console.error('[Vue Error] Info:', info)

// 开发环境显示详细错误，生产环境显示友好提示
  if (import.meta.env.DEV) {
    console.error('[Vue Error] Stack:', err.stack)
  }
}

// 捕获未处理的 Promise 拒绝
window.addEventListener('unhandledrejection', event => {
  console.error('[Unhandled Rejection]', event.reason)
})

// 捕获未捕获的 JavaScript 错误
window.addEventListener('error', event => {
  console.error('[Global Error]', event.error)
})



// 初始化用户信息
const userStore = useUserStore()
userStore.loadUserInfo()

// 初始化主题
const themeStore = useThemeStore()
themeStore.initTheme()

app.mount('#app')