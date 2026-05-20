import { createApp } from 'vue'
import { createPinia } from 'pinia'
import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'
import App from './App.vue'
import router from './router'
import axios from 'axios'
import { useUserStore } from './store'
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
app.use(ElementPlus, {
  theme: {
    primary: '#2c7d34',
  },
})

// 注册权限指令
app.directive('has-permission', hasPermission)

// 初始化用户信息
const userStore = useUserStore()
userStore.loadUserInfo()

app.mount('#app')