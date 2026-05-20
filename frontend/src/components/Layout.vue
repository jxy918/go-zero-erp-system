<template>
  <div class="layout" :class="{ 'dark-mode': isDarkMode }">
    <!-- 移动端遮罩层 -->
    <div 
      class="sidebar-overlay" 
      :class="{ 'visible': showSidebar }" 
      @click="closeSidebar"
    ></div>
    <el-container>
      <el-aside 
        :width="isCollapse ? '60px' : '200px'" 
        class="aside" 
        :class="{ 'collapsed': isCollapse, 'mobile-show': showSidebar }"
      >
        <div class="logo">
          <!-- 移动端关闭按钮 - 使用与打开按钮相同的风格 -->
          <button class="menu-btn mobile-close-btn" @click="closeSidebar">
            <svg class="menu-icon" viewBox="0 0 24 24" width="20" height="20">
              <path d="M18 6L6 18M6 6l12 12" stroke="currentColor" stroke-width="2" stroke-linecap="round" fill="none"/>
            </svg>
          </button>
          <div class="logo-img" v-if="!isCollapse">
            <svg viewBox="0 0 32 32" class="logo-svg">
              <defs>
                <linearGradient id="logoGradient" x1="0%" y1="0%" x2="100%" y2="100%">
                  <stop offset="0%" style="stop-color:#2c7d34"/>
                  <stop offset="100%" style="stop-color:#4CAF50"/>
                </linearGradient>
              </defs>
              <rect x="2" y="2" width="28" height="28" rx="6" fill="url(#logoGradient)"/>
              <path d="M10 10h12M10 16h12M10 22h8" stroke="white" stroke-width="2" stroke-linecap="round"/>
            </svg>
          </div>
          <div class="logo-img logo-collapsed" v-else>
            <svg viewBox="0 0 32 32" class="logo-svg">
              <rect x="2" y="2" width="28" height="28" rx="6" fill="url(#logoGradient)"/>
              <path d="M10 10h12M10 16h12M10 22h8" stroke="white" stroke-width="2" stroke-linecap="round"/>
            </svg>
          </div>
          <h2 v-if="!isCollapse">后台管理系统</h2>
          <h2 v-else class="logo-text-collapsed">后台</h2>
        </div>
        <el-menu
          :default-active="activeMenu"
          class="el-menu-vertical-demo"
          :background-color="isDarkMode ? 'rgba(30, 30, 30, 0.8)' : 'rgba(255, 255, 255, 0.8)'"
          :text-color="isDarkMode ? '#e0e0e0' : '#333'"
          :active-text-color="isDarkMode ? '#4CAF50' : '#2c7d34'"
          :collapse="isCollapse"
        >
          <!-- 动态渲染菜单 -->
          <template v-for="menu in menuStore.getMenuTree" :key="menu.id">
            <!-- 有子菜单的情况 -->
            <el-sub-menu v-if="menu.children && menu.children.length > 0" :index="menu.path">
              <template #title>
                <el-icon><component :is="getMenuIcon(menu.icon)"></component></el-icon>
                <span class="menu-text">{{ menu.name }}</span>
              </template>
              <el-menu-item
                v-for="child in menu.children"
                :key="child.id"
                :index="child.path"
                @click="navigate(child.path)"
              >
                <el-icon><component :is="getMenuIcon(child.icon)"></component></el-icon>
                <span class="menu-text">{{ child.name }}</span>
              </el-menu-item>
            </el-sub-menu>
            <!-- 没有子菜单的情况 -->
            <el-menu-item
              v-else
              :index="menu.path"
              @click="navigate(menu.path)"
            >
              <el-icon><component :is="getMenuIcon(menu.icon)"></component></el-icon>
              <span class="menu-text">{{ menu.name }}</span>
            </el-menu-item>
          </template>
        </el-menu>
      </el-aside>
      <el-container>
        <el-header class="header">
          <div class="header-left">
            <button type="button" class="menu-btn menu-toggle" @click="toggleCollapse">
              <svg v-if="isMobile" class="menu-icon" viewBox="0 0 24 24" width="22" height="22">
                <path d="M3 12h18M3 6h18M3 18h18" stroke="currentColor" stroke-width="2" stroke-linecap="round" fill="none"/>
              </svg>
              <svg v-else-if="isCollapse" class="menu-icon" viewBox="0 0 24 24" width="20" height="20">
                <path d="M10 18l6-6-6-6" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" fill="none"/>
              </svg>
              <svg v-else class="menu-icon" viewBox="0 0 24 24" width="20" height="20">
                <path d="M14 6l-6 6 6 6" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" fill="none"/>
              </svg>
            </button>
          </div>
          <div class="header-right">
            <el-button link @click="toggleTheme" class="theme-toggle">
              <el-icon v-if="isDarkMode"><Sunny /></el-icon>
              <el-icon v-else><Moon /></el-icon>
            </el-button>
            <el-dropdown>
              <span class="user">
                <el-avatar>{{ username }}</el-avatar>
                <span class="username">{{ username }}</span>
              </span>
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item @click="handleLogout">退出登录</el-dropdown-item>
                </el-dropdown-menu>
              </template>
            </el-dropdown>
          </div>
        </el-header>
        <el-main class="main">
          <router-view v-slot="{ Component }">
            <transition name="fade" mode="out-in">
              <component :is="Component" />
            </transition>
          </router-view>
        </el-main>
      </el-container>
    </el-container>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted, watch } from 'vue'
import { useRouter } from 'vue-router'
import * as ElIcons from '@element-plus/icons-vue'
import { Sunny, Moon } from '@element-plus/icons-vue'
import { useUserStore, useMenuStore } from '../store'

const router = useRouter()
const userStore = useUserStore()
const menuStore = useMenuStore()

const isCollapse = ref(false)
const username = ref('admin')
const isDarkMode = ref(false)
const showSidebar = ref(false)

const activeMenu = computed(() => {
  return router.currentRoute.value.path
})



const toggleCollapse = () => {
  // 移动端切换侧边栏显示/隐藏
  if (window.innerWidth <= 768) {
    showSidebar.value = !showSidebar.value
  } else {
    isCollapse.value = !isCollapse.value
  }
}

const isMobile = ref(window.innerWidth <= 768)

const handleResize = () => {
  isMobile.value = window.innerWidth <= 768
}

const closeSidebar = () => {
  showSidebar.value = false
}

const toggleTheme = () => {
  isDarkMode.value = !isDarkMode.value
  localStorage.setItem('darkMode', isDarkMode.value.toString())
  updateBodyClass()
}

const updateBodyClass = () => {
  if (isDarkMode.value) {
    document.body.classList.add('dark-mode')
  } else {
    document.body.classList.remove('dark-mode')
  }
}

const navigate = (path) => {
  router.push(path)
  // 移动端点击菜单项后关闭侧边栏
  if (window.innerWidth <= 768) {
    showSidebar.value = false
  }
}

const handleLogout = () => {
  userStore.logout()
  router.push('/login')
}

// 获取菜单图标
const getMenuIcon = (iconName) => {
  const icon = ElIcons[iconName] || ElIcons['House']
  return icon || ElIcons['House']
}

onMounted(async () => {
  // 从localStorage加载用户信息
  userStore.loadUserInfo()
  
  const userInfo = userStore.getUserInfo
  if (userInfo) {
    username.value = userInfo.username
  }
  
  // 加载菜单数据
  await menuStore.loadMenuTree()
  
  const savedDarkMode = localStorage.getItem('darkMode')
  if (savedDarkMode !== null) {
    isDarkMode.value = savedDarkMode === 'true'
  } else {
    // 检测系统偏好
    isDarkMode.value = window.matchMedia('(prefers-color-scheme: dark)').matches
  }
  
  updateBodyClass()
  
  // 添加窗口大小变化监听
  window.addEventListener('resize', handleResize)
})

onUnmounted(() => {
  // 移除窗口大小变化监听
  window.removeEventListener('resize', handleResize)
})

// 监听系统主题变化
watch(isDarkMode, () => {
  updateBodyClass()
})
</script>

<style>
/* 全局样式 */
:root {
  --primary-color: #2c7d34;
  --primary-light: #4CAF50;
  --background-light: #f5f7fa;
  --background-dark: #1a1a1a;
  --card-light: #ffffff;
  --card-dark: #2d2d2d;
  --text-light: #333333;
  --text-dark: #e0e0e0;
  --border-light: #e0e0e0;
  --border-dark: #404040;
  --shadow-light: 0 4px 6px rgba(0, 0, 0, 0.05), 0 1px 3px rgba(0, 0, 0, 0.1);
  --shadow-dark: 0 4px 6px rgba(0, 0, 0, 0.3), 0 1px 3px rgba(0, 0, 0, 0.4);
  --shadow-hover-light: 0 10px 15px rgba(0, 0, 0, 0.1), 0 4px 6px rgba(0, 0, 0, 0.05);
  --shadow-hover-dark: 0 10px 15px rgba(0, 0, 0, 0.4), 0 4px 6px rgba(0, 0, 0, 0.3);
  --border-radius: 8px;
}

body {
  transition: background-color 0.3s ease, color 0.3s ease;
  background-color: var(--background-light);
  color: var(--text-light);
  font-family: 'Helvetica Neue', Helvetica, 'PingFang SC', 'Hiragino Sans GB', 'Microsoft YaHei', Arial, sans-serif;
}

body.dark-mode {
  background-color: var(--background-dark);
  color: var(--text-dark);
}

/* 卡片样式 */
.el-card {
  border-radius: var(--border-radius) !important;
  border: 1px solid var(--border-light);
  box-shadow: var(--shadow-light);
  transition: all 0.3s ease;
  background-color: var(--card-light);
  overflow: hidden;
}

body.dark-mode .el-card {
  border: 1px solid var(--border-dark);
  box-shadow: var(--shadow-dark);
  background-color: var(--card-dark);
}

.el-card:hover {
  transform: translateY(-5px);
  box-shadow: var(--shadow-hover-light);
}

body.dark-mode .el-card:hover {
  box-shadow: var(--shadow-hover-dark);
}

/* 按钮样式 */
.el-button {
  border-radius: var(--border-radius) !important;
  transition: all 0.3s ease;
}

.el-button:hover {
  transform: translateY(-1px);
}

/* 表格样式 */
.el-table {
  border-radius: var(--border-radius) !important;
  overflow: hidden;
}

.el-table th {
  border-radius: var(--border-radius) var(--border-radius) 0 0 !important;
}

/* 对话框样式 */
.el-dialog {
  border-radius: var(--border-radius) !important;
  overflow: hidden;
}

.el-dialog__header {
  background-color: var(--card-light);
  border-bottom: 1px solid var(--border-light);
}

body.dark-mode .el-dialog__header {
  background-color: var(--card-dark);
  border-bottom: 1px solid var(--border-dark);
  color: var(--text-dark);
}

.el-dialog__body {
  background-color: var(--card-light);
  color: var(--text-light);
}

body.dark-mode .el-dialog__body {
  background-color: var(--card-dark);
  color: var(--text-dark);
}

.el-dialog__footer {
  background-color: var(--card-light);
  border-top: 1px solid var(--border-light);
}

body.dark-mode .el-dialog__footer {
  background-color: var(--card-dark);
  border-top: 1px solid var(--border-dark);
}

/* 表单样式 */
.el-form-item__label {
  color: var(--text-light);
}

body.dark-mode .el-form-item__label {
  color: var(--text-dark);
}

.el-input__wrapper {
  border-radius: var(--border-radius) !important;
}

body.dark-mode .el-input__wrapper {
  background-color: var(--card-dark) !important;
  border-color: var(--border-dark) !important;
}

body.dark-mode .el-input__inner {
  color: var(--text-dark) !important;
}

/* 树组件样式 */
.el-tree {
  background-color: transparent;
}

body.dark-mode .el-tree-node__content:hover {
  background-color: rgba(76, 175, 80, 0.15) !important;
}

/* 下拉菜单样式 */
.el-dropdown-menu {
  border-radius: var(--border-radius) !important;
  box-shadow: var(--shadow-light) !important;
}

body.dark-mode .el-dropdown-menu {
  background-color: var(--card-dark) !important;
  border: 1px solid var(--border-dark) !important;
  box-shadow: var(--shadow-dark) !important;
}

body.dark-mode .el-dropdown-menu__item {
  color: var(--text-dark) !important;
}

body.dark-mode .el-dropdown-menu__item:hover {
  background-color: rgba(76, 175, 80, 0.15) !important;
}

/* 标签样式 */
body.dark-mode .el-tag {
  background-color: rgba(76, 175, 80, 0.2) !important;
  border-color: rgba(76, 175, 80, 0.3) !important;
  color: var(--text-dark) !important;
}

body.dark-mode .el-tag--success {
  background-color: rgba(76, 175, 80, 0.3) !important;
  border-color: rgba(76, 175, 80, 0.4) !important;
  color: #4CAF50 !important;
}

body.dark-mode .el-tag--danger {
  background-color: rgba(244, 67, 54, 0.3) !important;
  border-color: rgba(244, 67, 54, 0.4) !important;
  color: #f44336 !important;
}

/* 分页样式 */
body.dark-mode .el-pagination__total {
  color: var(--text-dark) !important;
}

body.dark-mode .el-pagination__sizes .el-select .el-input__inner {
  color: var(--text-dark) !important;
}

body.dark-mode .el-pagination button {
  color: var(--text-dark) !important;
}

body.dark-mode .el-pagination__item {
  background-color: var(--card-dark) !important;
  border-color: var(--border-dark) !important;
  color: var(--text-dark) !important;
}

body.dark-mode .el-pagination__item:hover {
  border-color: var(--primary-color) !important;
  color: var(--primary-color) !important;
}

body.dark-mode .el-pagination__item.active {
  background-color: var(--primary-color) !important;
  border-color: var(--primary-color) !important;
  color: white !important;
}

/* 页面切换动画 */
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.3s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}

/* 涟漪效果 */
.el-button {
  position: relative;
  overflow: hidden;
}

.el-button::after {
  content: '';
  position: absolute;
  top: 50%;
  left: 50%;
  width: 0;
  height: 0;
  border-radius: 50%;
  background-color: rgba(255, 255, 255, 0.5);
  transform: translate(-50%, -50%);
  transition: width 0.6s, height 0.6s;
}

.el-button:active::after {
  width: 300px;
  height: 300px;
}

body.dark-mode .el-button::after {
  background-color: rgba(255, 255, 255, 0.2);
}
</style>

<style scoped>
.layout {
  height: 100vh;
  overflow: hidden;
  transition: background-color 0.3s ease;
}

.layout.dark-mode {
  background-color: var(--background-dark);
}

.el-container {
  height: 100%;
  display: flex;
}

.el-main {
  flex: 1;
  overflow-y: auto;
  transition: background-color 0.3s ease;
}

.aside {
  height: 100vh;
  transition: width 0.3s ease, background-color 0.3s ease;
  backdrop-filter: blur(10px);
  -webkit-backdrop-filter: blur(10px);
  border-right: 1px solid rgba(255, 255, 255, 0.1);
  box-shadow: 2px 0 10px rgba(0, 0, 0, 0.05);
}

.aside.collapsed {
  width: 60px !important;
}

/* 统一的菜单按钮样式 */
.menu-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 36px;
  height: 36px;
  border: 1px solid rgba(44, 125, 52, 0.3);
  background: rgba(44, 125, 52, 0.1);
  border-radius: 8px;
  cursor: pointer;
  color: #2c7d34;
  transition: all 0.3s ease;
  flex-shrink: 0;
}

.menu-btn:hover {
  background: rgba(44, 125, 52, 0.2);
  transform: translateY(-1px);
  box-shadow: 0 2px 4px rgba(44, 125, 52, 0.15);
}

.menu-btn .menu-icon {
  color: inherit;
}

/* 移动端关闭按钮 */
.mobile-close-btn {
  display: none;
  margin-right: 8px;
}

.layout.dark-mode .menu-btn {
  border-color: rgba(76, 175, 80, 0.4);
  background: rgba(76, 175, 80, 0.15);
  color: #4CAF50;
}

.layout.dark-mode .menu-btn:hover {
  background: rgba(76, 175, 80, 0.25);
  box-shadow: 0 2px 4px rgba(76, 175, 80, 0.2);
}

/* 菜单文本样式 */
.menu-text {
  transition: opacity 0.3s ease, transform 0.3s ease;
}

/* 折叠状态下的菜单文本 */
.aside.collapsed .menu-text {
  opacity: 0;
  transform: translateX(-10px);
  position: absolute;
  white-space: nowrap;
  background-color: var(--card-light);
  color: var(--text-light);
  padding: 4px 8px;
  border-radius: 4px;
  box-shadow: var(--shadow-light);
  z-index: 1000;
  margin-left: 10px;
  pointer-events: none;
}

.layout.dark-mode .aside.collapsed .menu-text {
  background-color: var(--card-dark);
  color: var(--text-dark);
  box-shadow: var(--shadow-dark);
}

/* 折叠状态下鼠标悬停显示菜单文本 */
.aside.collapsed .el-menu-item:hover .menu-text,
.aside.collapsed .el-sub-menu__title:hover .menu-text {
  opacity: 1;
  transform: translateX(0);
}

/* 非折叠状态下的菜单文本 */
.aside:not(.collapsed) .menu-text {
  opacity: 1;
  transform: translateX(0);
  position: static;
  background-color: transparent;
  padding: 0;
  box-shadow: none;
  pointer-events: auto;
}

.logo {
  height: 60px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--text-light);
  border-bottom: 1px solid rgba(255, 255, 255, 0.1);
  transition: color 0.3s ease;
  gap: 10px;
}

.layout.dark-mode .logo {
  color: var(--text-dark);
}

.logo-img {
  width: 32px;
  height: 32px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.logo-svg {
  width: 100%;
  height: 100%;
}

.logo-img img {
  width: 100%;
  height: 100%;
  object-fit: contain;
}

.logo-collapsed {
  width: 32px;
  height: 32px;
}

.logo h2 {
  font-size: 18px;
  font-weight: bold;
  margin: 0;
  transition: font-size 0.3s ease;
  white-space: nowrap;
}

.aside.collapsed .logo h2 {
  font-size: 14px;
}

.logo-text-collapsed {
  display: none;
}

.aside.collapsed .logo-text-collapsed {
  display: block;
}

.header {
  backdrop-filter: blur(10px);
  -webkit-backdrop-filter: blur(10px);
  background-color: rgba(255, 255, 255, 0.8);
  color: var(--text-light);
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.05);
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0 20px;
  height: 60px;
  transition: background-color 0.3s ease, color 0.3s ease, box-shadow 0.3s ease;
}

.layout.dark-mode .header {
  background-color: rgba(30, 30, 30, 0.8);
  color: var(--text-dark);
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.3);
}

.header-left {
  display: flex;
  align-items: center;
}

.menu-toggle,
.theme-toggle {
  color: var(--text-light);
  transition: color 0.3s ease;
}

.layout.dark-mode .menu-toggle,
.layout.dark-mode .theme-toggle {
  color: var(--text-dark);
}

.header-right {
  display: flex;
  align-items: center;
  gap: 15px;
}

.user {
  display: flex;
  align-items: center;
  cursor: pointer;
  transition: color 0.3s ease;
}

.username {
  margin-left: 10px;
  display: block;
  transition: color 0.3s ease;
}

.layout.dark-mode .username {
  color: var(--text-dark);
}

.main {
  padding: 20px;
  background-color: var(--background-light);
  overflow-y: auto;
  transition: background-color 0.3s ease;
}

.layout.dark-mode .main {
  background-color: var(--background-dark);
}

/* 响应式设计 */
@media screen and (max-width: 768px) {
  /* 移动端菜单按钮样式优化 */
  .menu-btn {
    width: 40px;
    height: 40px;
    border: 1px solid rgba(44, 125, 52, 0.3);
    background: rgba(44, 125, 52, 0.15);
    border-radius: 8px;
  }
  
  .menu-btn:hover {
    background: rgba(44, 125, 52, 0.25);
  }
  
  .layout.dark-mode .menu-btn {
    border-color: rgba(76, 175, 80, 0.4);
    background: rgba(76, 175, 80, 0.2);
  }
  
  .layout.dark-mode .menu-btn:hover {
    background: rgba(76, 175, 80, 0.3);
  }
  
  .menu-btn .menu-icon {
    width: 24px;
    height: 24px;
  }
  
  .aside {
    position: fixed;
    left: 0;
    top: 0;
    z-index: 1000;
    width: 240px !important;
    transform: translateX(-100%);
    transition: transform 0.3s ease;
    height: 100vh;
    box-shadow: 2px 0 20px rgba(0, 0, 0, 0.2);
  }
  
  .aside.mobile-show {
    transform: translateX(0);
  }
  
  /* 移动端遮罩层 */
  .sidebar-overlay {
    position: fixed;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background-color: rgba(0, 0, 0, 0.5);
    z-index: 999;
    opacity: 0;
    visibility: hidden;
    transition: opacity 0.3s ease, visibility 0.3s ease;
  }
  
  .sidebar-overlay.visible {
    opacity: 1;
    visibility: visible;
  }
  
  .header {
    padding: 0 12px;
  }
  
  .username {
    display: none;
  }
  
  .main {
    padding: 12px;
  }
  
  .el-container {
    flex-direction: column;
  }
  
  /* 移动端菜单图标样式优化 */
  .el-menu-item, .el-sub-menu__title {
    padding: 0 15px !important;
  }
  
  .el-menu-item .el-icon, .el-sub-menu__title .el-icon {
    font-size: 18px !important;
  }
  
  /* 移动端显示关闭按钮 */
  .mobile-close-btn {
    display: flex;
  }
  
  /* 移动端logo区域调整 */
  .logo {
    justify-content: flex-start;
    padding: 0 12px;
    gap: 10px;
  }
  
  /* 移动端菜单项间距优化 */
  .el-menu-item, .el-sub-menu__title {
    height: 48px !important;
    line-height: 48px !important;
  }
  
  /* 移动端隐藏折叠状态的样式 */
  .aside.collapsed {
    width: 240px !important;
  }
}

@media screen and (min-width: 769px) {
  .aside {
    transform: translateX(0) !important;
  }
}
</style>