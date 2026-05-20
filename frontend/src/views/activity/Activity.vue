<template>
  <div class="activity-management">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>活动日志</span>
        </div>
      </template>
      <div class="description-box">
        <span class="description-icon">ℹ️</span>
        <span class="description-text">活动日志：记录用户在系统中的所有操作行为，包括登录、增删改查等操作，便于审计追踪。</span>
      </div>
      <div class="search-container">
        <el-input
          v-model="searchQuery"
          placeholder="请输入用户名"
          style="width: 300px"
          prefix-icon="el-icon-search"
        ></el-input>
        <el-button type="primary" @click="handleSearch">搜索</el-button>
      </div>
      <el-table 
        :data="activities" 
        style="width: 100%"
        stripe
        :header-cell-style="{ backgroundColor: 'var(--card-light)', fontWeight: 'bold' }"
        :cell-style="{ transition: 'background-color 0.3s ease' }"
        :row-style="{ transition: 'background-color 0.3s ease' }"
      >
        <el-table-column prop="id" label="ID" width="80" fixed></el-table-column>
        <el-table-column prop="username" label="用户名"></el-table-column>
        <el-table-column prop="action" label="操作"></el-table-column>
        <el-table-column prop="url" label="URL" show-overflow-tooltip></el-table-column>
        <el-table-column prop="ip" label="IP地址"></el-table-column>
        <el-table-column prop="created_at" label="时间" width="180">
          <template #default="scope">
            {{ formatTime(scope.row.created_at) }}
          </template>
        </el-table-column>
      </el-table>
      <div class="pagination-container">
        <el-pagination
          v-model:current-page="currentPage"
          v-model:page-size="pageSize"
          :page-sizes="[10, 20, 50, 100]"
          layout="total, sizes, prev, pager, next, jumper"
          :total="total"
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
        ></el-pagination>
      </div>
    </el-card>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useActivityStore } from '../../store'

const activityStore = useActivityStore()
const searchQuery = ref('')

const activities = computed(() => {
  return activityStore.getActivityList
})

const total = computed(() => {
  return activityStore.getTotal
})

const currentPage = computed({
  get: () => activityStore.getSearchParams.page,
  set: (value) => activityStore.setSearchParams({ page: value })
})

const pageSize = computed({
  get: () => activityStore.getSearchParams.pageSize,
  set: (value) => activityStore.setSearchParams({ pageSize: value })
})

const handleSearch = async () => {
  await activityStore.searchActivity(searchQuery.value)
}

const handleSizeChange = async (size) => {
  await activityStore.changePageSize(size)
}

const handleCurrentChange = async (current) => {
  await activityStore.changePage(current)
}

const formatTime = (time) => {
  if (!time) return ''
  const date = new Date(time)
  return date.toLocaleString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
    second: '2-digit'
  })
}

onMounted(async () => {
  await activityStore.loadActivityList()
})
</script>

<style scoped>
.activity-management {
  padding: 0;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.description-box {
  display: flex;
  align-items: center;
  padding: 12px 16px;
  background-color: #fff7f0;
  border-left: 4px solid #ff6b35;
  margin-bottom: 16px;
  border-radius: 4px;
}

.description-icon {
  font-size: 18px;
  margin-right: 10px;
}

.description-text {
  color: #d93026;
  font-size: 13px;
  line-height: 1.5;
}

.search-container {
  margin: 20px 0;
  display: flex;
  gap: 10px;
}

.pagination-container {
  margin-top: 20px;
  margin-bottom: 20px;
  display: flex;
  justify-content: flex-end;
  width: 100%;
  padding: 0 10px;
}

/* 表格样式 */
:deep(.el-table__header th) {
  border-radius: var(--border-radius) var(--border-radius) 0 0 !important;
}

:deep(.el-table__row:hover) {
  background-color: rgba(76, 175, 80, 0.05) !important;
}

body.dark-mode :deep(.el-table__row:hover) {
  background-color: rgba(76, 175, 80, 0.15) !important;
}

body.dark-mode :deep(.el-table__header th) {
  background-color: var(--card-dark) !important;
  color: var(--text-dark) !important;
}

/* 响应式设计 */
@media screen and (max-width: 768px) {
  .card-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 10px;
  }
  
  .search-container {
    flex-direction: column;
    gap: 10px;
  }
  
  .search-container .el-input {
    width: 100% !important;
  }
  
  .pagination-container {
    justify-content: center;
  }
  
  .el-table {
    font-size: 12px;
  }
  
  .el-table th,
  .el-table td {
    padding: 8px 4px;
  }
  
  .el-button {
    font-size: 12px;
    padding: 4px 8px;
  }
}
</style>