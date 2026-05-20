<template>
  <div class="login-container">
    <el-card class="login-card">
      <template #header>
        <div class="login-header">
          <h2>后台管理系统</h2>
          <p>请登录您的账户</p>
        </div>
      </template>
      <el-form :model="loginForm" :rules="rules" ref="loginFormRef" label-width="80px">
        <el-form-item label="用户名" prop="username">
          <el-input v-model="loginForm.username" placeholder="请输入用户名"></el-input>
        </el-form-item>
        <el-form-item label="密码" prop="password">
          <el-input
            v-model="loginForm.password"
            type="password"
            placeholder="请输入密码"
            show-password
          ></el-input>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleLogin" style="width: 100%">登录</el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script setup>
import { ref } from 'vue'
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
          
          window.location.href = '/'
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
  background-color: #f5f7fa;
}

.login-card {
  width: 400px;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
}

.login-header {
  text-align: center;
  margin-bottom: 20px;
}

.login-header h2 {
  font-size: 24px;
  font-weight: bold;
  color: #303133;
  margin-bottom: 10px;
}

.login-header p {
  font-size: 14px;
  color: #606266;
}
</style>