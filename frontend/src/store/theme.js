/**
 * 主题状态管理
 */
import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

export const useThemeStore = defineStore('theme', () => {
  // 可用主题列表
  const themes = ref([
    {
      id: 'default',
      name: '默认风格',
      cssFile: 'default-theme.css'
    },
    {
      id: 'business',
      name: '商务专业',
      cssFile: 'business-theme.css'
    },
    {
      id: 'dark',
      name: '深色模式',
      cssFile: 'dark-theme.css'
    },
    {
      id: 'modern',
      name: '现代简约',
      cssFile: 'modern-theme.css'
    }
  ])

  // 当前主题ID
  const currentThemeId = ref(localStorage.getItem('theme') || 'default')

  // 当前主题
  const currentTheme = computed(() => {
    return themes.value.find(t => t.id === currentThemeId.value) || themes.value[0]
  })

  // 设置主题
  function setTheme(themeId) {
    if (themes.value.find(t => t.id === themeId)) {
      currentThemeId.value = themeId
      localStorage.setItem('theme', themeId)
      applyTheme(themeId)
    }
  }

  // 应用主题
  function applyTheme(themeId) {
    const theme = themes.value.find(t => t.id === themeId)
    if (!theme) return

    // 移除所有主题类名
    document.body.classList.remove('theme-default', 'theme-business', 'theme-dark', 'theme-modern')
    
    // 添加当前主题类名
    document.body.classList.add(`theme-${themeId}`)
    
    // 动态加载主题CSS
    loadThemeCSS(themeId)
  }

  // 动态加载主题CSS文件
  function loadThemeCSS(themeId) {
    // 移除旧的主题CSS
    const oldLink = document.getElementById('theme-css')
    if (oldLink) {
      oldLink.remove()
    }

    // 创建新的link
    const link = document.createElement('link')
    link.id = 'theme-css'
    link.rel = 'stylesheet'
    link.href = `/src/styles/${themeId}.css`
    document.head.appendChild(link)
  }

  // 初始化主题
  function initTheme() {
    applyTheme(currentThemeId.value)
  }

  return {
    themes,
    currentThemeId,
    currentTheme,
    setTheme,
    initTheme
  }
})