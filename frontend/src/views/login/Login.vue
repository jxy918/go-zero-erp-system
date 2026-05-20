<template>
  <div 
    class="login-container" 
    :style="containerStyle"
  >
    <el-card class="login-card" :style="cardStyle">
      <template #header>
        <div class="login-header">
          <h2 :style="{ color: themeConfig.titleColor }">后台管理系统</h2>
          <p :style="{ color: themeConfig.subtitleColor }">请登录您的账户</p>
        </div>
      </template>
      <el-form :model="loginForm" :rules="rules" ref="loginFormRef" label-width="80px">
        <el-form-item label="用户名" prop="username">
          <el-input 
            v-model="loginForm.username" 
            placeholder="请输入用户名"
            :style="inputStyle"
          ></el-input>
        </el-form-item>
        <el-form-item label="密码" prop="password">
          <el-input
            v-model="loginForm.password"
            type="password"
            placeholder="请输入密码"
            show-password
            :style="inputStyle"
          ></el-input>
        </el-form-item>
        <el-form-item>
          <el-button 
            type="primary" 
            @click="handleLogin" 
            style="width: 100%"
            :style="buttonStyle"
          >登录</el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { useUserStore, useMenuStore } from '../../store'
import { authApi } from '../../api'

const router = useRouter()
const userStore = useUserStore()
const menuStore = useMenuStore()
const loginForm = ref({
  username: '',
  password: ''
})
const rules = {
  username: [{ required: true, message: '请输入用户名', trigger: 'blur' }],
  password: [{ required: true, message: '请输入密码', trigger: 'blur' }]
}
const loginFormRef = ref(null)

const currentTheme = ref(localStorage.getItem('theme') || 'default')

const themeConfigs = {
  default: {
    containerBg: 'linear-gradient(135deg, #f5f7fa 0%, #fff 50%, #f5f7fa 100%)',
    cardBg: '#ffffff',
    cardBorder: '#e4e7ed',
    cardShadow: '0 4px 12px rgba(0, 0, 0, 0.1)',
    titleColor: '#303133',
    subtitleColor: '#909399',
    cardHeaderBg: '#fff',
    cardHeaderBorder: '#e4e7ed',
    labelColor: '#606266',
    inputBg: '#fff',
    inputBorder: '#dcdfe6',
    inputHoverBorder: '#c0c4cc',
    inputFocusBorder: '#2c7d34',
    inputFocusShadow: '0 0 0 4px rgba(44, 125, 52, 0.1)',
    inputColor: '#303133',
    placeholderColor: '#c0c4cc',
    buttonBg: '#2c7d34',
    buttonBorder: '#2c7d34',
    buttonHoverBg: '#4a9d56',
    buttonHoverBorder: '#4a9d56'
  },
  business: {
    containerBg: 'linear-gradient(135deg, #f7fafc 0%, #fff 50%, #f7fafc 100%)',
    cardBg: '#ffffff',
    cardBorder: '#edf2f7',
    cardShadow: '0 4px 12px rgba(30, 58, 95, 0.08)',
    titleColor: '#1a202c',
    subtitleColor: '#718096',
    cardHeaderBg: 'linear-gradient(135deg, #f7fafc, #ffffff)',
    cardHeaderBorder: '#edf2f7',
    labelColor: '#4a5568',
    inputBg: '#fff',
    inputBorder: '#e2e8f0',
    inputHoverBorder: '#cbd5e0',
    inputFocusBorder: '#1e3a5f',
    inputFocusShadow: '0 0 0 4px rgba(30, 58, 95, 0.12)',
    inputColor: '#1a202c',
    placeholderColor: '#a0aec0',
    buttonBg: 'linear-gradient(135deg, #1e3a5f, #2d5a87)',
    buttonBorder: 'none',
    buttonHoverBg: 'linear-gradient(135deg, #2d5a87, #1e3a5f)',
    buttonHoverBorder: 'none'
  },
  dark: {
    containerBg: 'linear-gradient(135deg, #0f172a 0%, #1e293b 50%, #0f172a 100%)',
    cardBg: 'linear-gradient(145deg, #1e293b, #1a2332)',
    cardBorder: 'rgba(51, 65, 85, 0.6)',
    cardShadow: '0 12px 40px rgba(0, 0, 0, 0.5)',
    titleColor: '#e2e8f0',
    subtitleColor: '#94a3b8',
    cardHeaderBg: 'linear-gradient(135deg, #1e293b, #273448)',
    cardHeaderBorder: 'rgba(51, 65, 85, 0.5)',
    labelColor: '#94a3b8',
    inputBg: 'linear-gradient(145deg, #1e293b, #243045)',
    inputBorder: 'rgba(51, 65, 85, 0.7)',
    inputHoverBorder: 'rgba(99, 102, 241, 0.6)',
    inputFocusBorder: '#6366f1',
    inputFocusShadow: '0 0 0 4px rgba(99, 102, 241, 0.15)',
    inputColor: '#e2e8f0',
    placeholderColor: '#475569',
    buttonBg: 'linear-gradient(135deg, #6366f1, #818cf8)',
    buttonBorder: 'none',
    buttonHoverBg: 'linear-gradient(135deg, #818cf8, #6366f1)',
    buttonHoverBorder: 'none'
  },
  modern: {
    containerBg: 'linear-gradient(135deg, #f8fafc 0%, #fff 50%, #f8fafc 100%)',
    cardBg: '#ffffff',
    cardBorder: 'none',
    cardShadow: '0 1px 3px rgba(0, 0, 0, 0.1)',
    titleColor: '#0f172a',
    subtitleColor: '#94a3b8',
    cardHeaderBg: '#fff',
    cardHeaderBorder: '#f1f5f9',
    labelColor: '#475569',
    inputBg: '#fff',
    inputBorder: '#e2e8f0',
    inputHoverBorder: '#94a3b8',
    inputFocusBorder: '#2563eb',
    inputFocusShadow: '0 0 0 3px rgba(37, 99, 235, 0.1)',
    inputColor: '#0f172a',
    placeholderColor: '#94a3b8',
    buttonBg: '#2563eb',
    buttonBorder: '#2563eb',
    buttonHoverBg: '#3b82f6',
    buttonHoverBorder: '#3b82f6'
  }
}

const themeConfig = computed(() => {
  return themeConfigs[currentTheme.value] || themeConfigs.default
})

const containerStyle = computed(() => ({
  background: themeConfig.value.containerBg,
  transition: 'background 0.35s ease'
}))

const cardStyle = computed(() => ({
  background: themeConfig.value.cardBg,
  border: `1px solid ${themeConfig.value.cardBorder}`,
  boxShadow: themeConfig.value.cardShadow,
  borderRadius: '12px',
  transition: 'all 0.35s ease'
}))

const inputStyle = computed(() => ({
  '--input-bg': themeConfig.value.inputBg,
  '--input-border': themeConfig.value.inputBorder,
  '--input-hover-border': themeConfig.value.inputHoverBorder,
  '--input-focus-border': themeConfig.value.inputFocusBorder,
  '--input-focus-shadow': themeConfig.value.inputFocusShadow,
  '--input-color': themeConfig.value.inputColor,
  '--placeholder-color': themeConfig.value.placeholderColor,
  '--label-color': themeConfig.value.labelColor,
  '--card-header-bg': themeConfig.value.cardHeaderBg,
  '--card-header-border': themeConfig.value.cardHeaderBorder
}))

const buttonStyle = computed(() => ({
  background: themeConfig.value.buttonBg,
  borderColor: themeConfig.value.buttonBorder,
  fontWeight: '600',
  transition: 'all 0.25s cubic-bezier(0.4, 0, 0.2, 1)'
}))

const handleLogin = () => {
  loginFormRef.value.validate((valid) => {
    if (valid) {
      // 移除前后空格
      loginForm.value.username = loginForm.value.username.trim()
      loginForm.value.password = loginForm.value.password.trim()
      
      authApi.login(loginForm.value)
        .then(response => {
          const data = response.data
          const token = data.token
          const userData = data.user
          
          localStorage.removeItem('user')
          localStorage.removeItem('token')
          
          userStore.login(userData, token)
          menuStore.clearMenuCache()
          
          ElMessage.success('登录成功')
          
          router.push('/')
        })
        .catch(error => {
          console.error('登录失败:', error)
          ElMessage.error('登录失败：' + (error.message || '网络错误'))
        })
    }
  })
}
</script>

<style scoped>
.login-container {
  height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
}

.login-card {
  width: 400px;
}

.login-header {
  text-align: center;
  margin-bottom: 24px;
}

.login-header h2 {
  font-size: 26px;
  font-weight: 700;
  margin-bottom: 8px;
  letter-spacing: 1px;
}

.login-header p {
  font-size: 14px;
}

<:deep(.el-card__header) {
  background: var(--card-header-bg, #fff);
  border-bottom: 1px solid var(--card-header-border, #e4e7ed);
  padding: 18px 24px;
}

:deep(.el-form-item__label) {
  color: var(--label-color, #606266);
  font-weight: 500;
}

:deep(.el-input__wrapper) {
  background: var(--input-bg, #fff);
  border: 1px solid var(--input-border, #dcdfe6);
  border-radius: 8px;
  transition: all 0.25s cubic-bezier(0.4, 0, 0.2, 1);
}

:deep(.el-input__wrapper:hover) {
  border-color: var(--input-hover-border, #c0c4cc);
}

:deep(.el-input__wrapper.is-focus) {
  border-color: var(--input-focus-border, #2c7d34);
  box-shadow: var(--input-focus-shadow, '0 0 0 4px rgba(44, 125, 52, 0.1)');
}

:deep(.el-input__inner) {
  color: var(--input-color, #303133);
}

:deep(.el-input__inner::placeholder) {
  color: var(--placeholder-color, #c0c4cc);
}

:deep(.el-button--primary) {
  border: none;
  font-weight: 600;
  transition: all 0.25s cubic-bezier(0.4, 0, 0.2, 1);
}

:deep(.el-form-item__error) {
  color: #fca5a5;
}
</style>