import { defineStore } from 'pinia'
import { activityApi } from '../api'

// 活动日志状态管理
export const useActivityStore = defineStore('activity', {
  state: () => ({
    activityList: [],
    total: 0,
    loading: false,
    error: null,
    searchParams: {
      username: '',
      page: 1,
      page_size: 10
    }
  }),
  getters: {
    getActivityList: (state) => state.activityList,
    getTotal: (state) => state.total,
    getLoading: (state) => state.loading,
    getError: (state) => state.error,
    getSearchParams: (state) => state.searchParams
  },
  actions: {
    // 设置搜索参数
    setSearchParams(params) {
      // 转换pageSize为page_size
      if (params.pageSize !== undefined) {
        params.page_size = params.pageSize
        delete params.pageSize
      }
      this.searchParams = { ...this.searchParams, ...params }
    },
    
    // 重置搜索参数
    resetSearchParams() {
      this.searchParams = {
        username: '',
        page: 1,
        page_size: 10
      }
    },
    
    // 加载活动日志列表（后端根据登录会话自动筛选）
    async loadActivityList() {
      this.loading = true
      this.error = null
      try {
        const response = await activityApi.getActivityList(this.searchParams)
        console.log('活动日志响应:', response)
        console.log('响应数据:', response.data)
        // 适配后端返回的数据格式
        this.activityList = response.data.activities || []
        this.total = response.data.total || 0
        console.log('活动列表:', this.activityList)
        console.log('总数:', this.total)
      } catch (error) {
        this.error = '获取活动日志失败'
        console.error('获取活动日志失败:', error)
      } finally {
        this.loading = false
      }
    },
    
    // 搜索活动日志
    async searchActivity(username) {
      this.setSearchParams({ username, page: 1 })
      await this.loadActivityList()
    },
    
    // 分页
    async changePage(page) {
      this.setSearchParams({ page })
      await this.loadActivityList()
    },
    
    // 改变每页条数
    async changePageSize(pageSize) {
      this.setSearchParams({ pageSize, page: 1 })
      await this.loadActivityList()
    }
  }
})
