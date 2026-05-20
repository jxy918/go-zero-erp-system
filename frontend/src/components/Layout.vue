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
            <ThemeSwitcher />
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
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import * as ElIcons from '@element-plus/icons-vue'
import { useUserStore, useMenuStore } from '../store'
import ThemeSwitcher from './ThemeSwitcher.vue'

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
  
  // 添加窗口大小变化监听
  window.addEventListener('resize', handleResize)
})

onUnmounted(() => {
  // 移除窗口大小变化监听
  window.removeEventListener('resize', handleResize)
})
</script>

<style scoped>
.layout {
  height: 100vh;
  overflow: hidden;
}

.el-container {
  height: 100%;
  display: flex;
}

.el-main {
  flex: 1;
  overflow-y: auto;
}

.aside {
  height: 100vh;
  transition: width 0.3s ease;
  box-shadow: 2px 0 10px rgba(0, 0, 0, 0.05);
}

.aside.collapsed {
  width: 60px !important;
}

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

.mobile-close-btn {
  display: none;
  margin-right: 8px;
}

.menu-text {
  transition: opacity 0.3s ease, transform 0.3s ease;
}

.aside.collapsed .menu-text {
  opacity: 0;
  transform: translateX(-10px);
  position: absolute;
  white-space: nowrap;
  padding: 4px 8px;
  border-radius: 4px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  z-index: 1000;
  margin-left: 10px;
  pointer-events: none;
  font-size: 14px;
}

.aside.collapsed .el-menu-item:hover .menu-text,
.aside.collapsed .el-sub-menu__title:hover .menu-text {
  opacity: 1;
  transform: translateX(0);
}

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
  gap: 10px;
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

.logo h2 {
  font-size: 18px;
  font-weight: bold;
  margin: 0;
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
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0 20px;
  height: 60px;
}

.header-left {
  display: flex;
  align-items: center;
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
}

.username {
  margin-left: 10px;
  display: block;
}

.main {
  padding: 20px;
}

.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.2s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}

@media screen and (max-width: 768px) {
  .menu-btn {
    width: 40px;
    height: 40px;
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
  
  .mobile-close-btn {
    display: flex;
  }
  
  .logo {
    justify-content: flex-start;
    padding: 0 12px;
  }
  
  .el-menu-item,
  .el-sub-menu__title {
    height: 48px !important;
    line-height: 48px !important;
  }
  
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