<template>
  <el-dropdown @command="handleCommand" trigger="click">
    <el-button text>
      <el-icon><component :is="currentIcon" /></el-icon>
      <span class="theme-name">{{ currentThemeName }}</span>
      <el-icon class="el-icon--right"><ArrowDown /></el-icon>
    </el-button>
    <template #dropdown>
      <el-dropdown-menu>
        <el-dropdown-item 
          v-for="theme in themes" 
          :key="theme.id" 
          :command="theme.id"
          :class="{ 'is-active': theme.id === currentThemeId }"
        >
          <div class="theme-item">
            <el-icon><component :is="getThemeIcon(theme.id)" /></el-icon>
            <span>{{ theme.name }}</span>
            <el-icon v-if="theme.id === currentThemeId" class="check-icon"><Check /></el-icon>
          </div>
        </el-dropdown-item>
      </el-dropdown-menu>
    </template>
  </el-dropdown>
</template>

<script setup>
import { computed } from 'vue'
import { ArrowDown, Check, Brush, Moon, Sunny, Grid } from '@element-plus/icons-vue'
import { useThemeStore } from '@/store/theme'

const themeStore = useThemeStore()

const themes = computed(() => themeStore.themes)
const currentThemeId = computed(() => themeStore.currentThemeId)

const currentThemeName = computed(() => {
  const theme = themes.value.find(t => t.id === currentThemeId.value)
  return theme ? theme.name : '默认风格'
})

const currentIcon = computed(() => getThemeIcon(currentThemeId.value))

function getThemeIcon(themeId) {
  const icons = {
    'default': Brush,
    'business': Grid,
    'dark': Moon,
    'modern': Sunny
  }
  return icons[themeId] || Brush
}

function handleCommand(themeId) {
  themeStore.setTheme(themeId)
}
</script>

<style scoped>
.theme-name {
  margin: 0 8px;
  font-size: 14px;
}

.theme-item {
  display: flex;
  align-items: center;
  gap: 8px;
  width: 100%;
}

.theme-item .el-icon {
  font-size: 16px;
}

.check-icon {
  margin-left: auto;
  color: var(--el-color-primary);
}

.el-dropdown-menu__item.is-active {
  background-color: var(--el-color-primary-light-9);
  color: var(--el-color-primary);
}

.el-dropdown-menu__item:hover {
  background-color: var(--el-color-fill-color-light);
}
</style>
