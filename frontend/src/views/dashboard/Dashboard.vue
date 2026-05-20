<template>
  <div class="dashboard">
    <el-card class="welcome-card">
      <template #header>
        <div class="card-header">
          <span>欢迎回来</span>
          <span class="current-time">{{ currentTime }}</span>
        </div>
      </template>
      <div class="welcome-content">
        <div class="avatar-wrapper">
          <el-avatar :size="100" :src="userAvatar"></el-avatar>
          <div class="status-dot"></div>
        </div>
        <h2>{{ username }}</h2>
        <p class="welcome-text">{{ greeting }}，今天是 {{ currentDate }}</p>
      </div>
    </el-card>

    <div class="stats-grid">
      <div class="stat-card">
        <div class="stat-card-header">
          <div class="stat-icon-box stat-icon-purchase">
            <el-icon class="stat-icon"><ShoppingCart /></el-icon>
          </div>
          <span class="stat-title">本月采购金额</span>
        </div>
        <div class="stat-number">¥{{ formatAmount(overviewData.purchase_amount) }}</div>
        <div v-if="purchaseTrend !== undefined" class="stat-change" :class="purchaseTrend >= 0 ? 'positive' : 'negative'">
          {{ purchaseTrend >= 0 ? '↑' : '↓' }} {{ Math.abs(purchaseTrend) }}%
        </div>
      </div>

      <div class="stat-card">
        <div class="stat-card-header">
          <div class="stat-icon-box stat-icon-sales">
            <el-icon class="stat-icon"><PieChart /></el-icon>
          </div>
          <span class="stat-title">本月销售金额</span>
        </div>
        <div class="stat-number">¥{{ formatAmount(overviewData.sales_amount) }}</div>
        <div v-if="salesTrend !== undefined" class="stat-change" :class="salesTrend >= 0 ? 'positive' : 'negative'">
          {{ salesTrend >= 0 ? '↑' : '↓' }} {{ Math.abs(salesTrend) }}%
        </div>
      </div>

      <div class="stat-card">
        <div class="stat-card-header">
          <div class="stat-icon-box stat-icon-inventory">
            <el-icon class="stat-icon"><Folder /></el-icon>
          </div>
          <span class="stat-title">当前库存总量</span>
        </div>
        <div class="stat-number">{{ overviewData.total_inventory }}<span class="stat-unit">件</span></div>
      </div>

      <div class="stat-card">
        <div class="stat-card-header">
          <div class="stat-icon-box stat-icon-order">
            <el-icon class="stat-icon"><Document /></el-icon>
          </div>
          <span class="stat-title">本月订单总数</span>
        </div>
        <div class="stat-number">{{ overviewData.total_orders }}</div>
        <div v-if="orderTrend !== undefined" class="stat-change" :class="orderTrend >= 0 ? 'positive' : 'negative'">
          {{ orderTrend >= 0 ? '↑' : '↓' }} {{ Math.abs(orderTrend) }}%
        </div>
      </div>

      <div class="stat-card">
        <div class="stat-card-header">
          <div class="stat-icon-box stat-icon-user">
            <el-icon class="stat-icon"><User /></el-icon>
          </div>
          <span class="stat-title">系统用户</span>
        </div>
        <div class="stat-number">{{ userCount }}</div>
      </div>

      <div class="stat-card">
        <div class="stat-card-header">
          <div class="stat-icon-box stat-icon-product">
            <el-icon class="stat-icon"><Box /></el-icon>
          </div>
          <span class="stat-title">产品数量</span>
        </div>
        <div class="stat-number">{{ productCount }}</div>
      </div>
    </div>

    <div class="main-content">
      <div class="left-panel">
        <el-card class="chart-card">
          <template #header>
            <div class="card-header">
              <span>采购/销售趋势</span>
              <el-select v-model="trendPeriod" class="period-select" @change="loadTrend">
                <el-option label="近7天" :value="7"></el-option>
                <el-option label="近30天" :value="30"></el-option>
                <el-option label="近90天" :value="90"></el-option>
              </el-select>
            </div>
          </template>
          <div ref="trendChart" class="chart-container"></div>
        </el-card>

        <el-card class="activity-card">
          <template #header>
            <div class="card-header">
              <span>最近活动</span>
              <el-button link @click="loadRecentActivities">刷新</el-button>
            </div>
          </template>
          <div class="activity-list">
            <div class="activity-list-header">
              <span class="col-time">时间</span>
              <span class="col-user">用户</span>
              <span class="col-action">操作</span>
              <span class="col-url">URL</span>
              <span class="col-ip">IP</span>
            </div>
            <div
              v-for="activity in recentActivities"
              :key="activity.id"
              class="activity-list-item"
            >
              <span class="col-time">{{ activity.time }}</span>
              <span class="col-user">{{ activity.username }}</span>
              <span class="col-action">{{ activity.action }}</span>
              <span class="col-url">{{ activity.url || '-' }}</span>
              <span class="col-ip">{{ activity.ip }}</span>
            </div>
            <div v-if="recentActivities.length === 0" class="empty-activity">暂无活动记录</div>
          </div>
        </el-card>
      </div>

      <div class="right-panel">
        <el-card class="alert-card" v-if="inventoryAlerts.length > 0">
          <template #header>
            <div class="card-header">
              <span class="alert-icon">
                <el-icon><Document /></el-icon>
              </span>
              <span>库存预警</span>
              <el-badge :value="inventoryAlertCount" class="badge" />
            </div>
          </template>
          <div class="alert-list">
            <div
              v-for="alert in inventoryAlerts"
              :key="alert.product_id + '-' + alert.warehouse_id"
              class="alert-item"
            >
              <div class="alert-info">
                <span class="alert-product">{{ alert.product_name }}</span>
                <span class="alert-quantity">库存: <span :class="getAlertQuantityClass(alert)">{{ alert.quantity }}</span></span>
              </div>
              <div class="alert-threshold">{{ getAlertTypeText(alert) }}</div>
              <div class="alert-gap">缺口: <span class="gap-value">{{ getAlertGap(alert) }}</span></div>
            </div>
          </div>
        </el-card>

        <el-card class="quick-actions-card">
          <template #header>
            <span>快捷操作</span>
          </template>
          <div class="quick-actions">
            <el-button
              v-for="action in quickActions"
              :key="action.path"
              class="action-btn"
              :type="action.type"
              @click="navigate(action.path)"
            >
              <el-icon><component :is="action.icon"></component></el-icon>
              <span>{{ action.label }}</span>
            </el-button>
          </div>
        </el-card>

        <el-card class="order-summary-card">
          <template #header>
            <span>订单状态概览</span>
          </template>
          <div class="order-summary">
            <div class="order-item pending">
              <div class="order-count">{{ orderStatus.pending }}</div>
              <div class="order-label">待审核</div>
            </div>
            <div class="order-item approved">
              <div class="order-count">{{ orderStatus.approved }}</div>
              <div class="order-label">已审核</div>
            </div>
            <div class="order-item completed">
              <div class="order-count">{{ orderStatus.completed }}</div>
              <div class="order-label">已完成</div>
            </div>
            <div class="order-item cancelled">
              <div class="order-count">{{ orderStatus.cancelled }}</div>
              <div class="order-label">已取消</div>
            </div>
          </div>
        </el-card>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onBeforeUnmount, nextTick } from 'vue'
import { useRouter } from 'vue-router'
import {
  ShoppingCart, Document, User, Box, View, PieChart, Folder,
  Goods, Shop, OfficeBuilding, Warning, Edit, RefreshRight,
  DataAnalysis, Setting, Tools, Avatar
} from '@element-plus/icons-vue'
import * as echarts from 'echarts'
import { erpApi } from '../../api/erp'
import { activityApi } from '../../api/activity'

let trendChartInstance = null
let resizeHandler = null

const router = useRouter()

const username = ref('admin')
const userAvatar = ref('https://cube.elemecdn.com/3/7c/3ea6beec64369c2642b92c6726f1epng.png')

const currentDate = computed(() => {
  return new Date().toLocaleDateString('zh-CN', {
    year: 'numeric',
    month: 'long',
    day: 'numeric',
    weekday: 'long'
  })
})

const currentTime = computed(() => {
  return new Date().toLocaleTimeString('zh-CN', {
    hour: '2-digit',
    minute: '2-digit',
    second: '2-digit'
  })
})

const greeting = computed(() => {
  const hour = new Date().getHours()
  if (hour < 6) return '夜深了'
  if (hour < 12) return '早上好'
  if (hour < 14) return '中午好'
  if (hour < 18) return '下午好'
  return '晚上好'
})

const overviewData = ref({
  purchase_amount: 0,
  sales_amount: 0,
  total_inventory: 0,
  total_orders: 0,
  user_count: 0,
  product_count: 0
})

const purchaseTrend = ref(12.5)
const salesTrend = ref(-2.3)
const orderTrend = ref(8.7)

const userCount = ref(0)
const productCount = ref(0)
const recentActivities = ref([])
const inventoryAlerts = ref([])
const inventoryAlertCount = ref(0)
const orderStatus = ref({
  pending: 0,
  approved: 0,
  completed: 0,
  cancelled: 0
})

const trendPeriod = ref(30)
const trendChart = ref(null)

const quickActions = [
  { label: '采购管理', icon: ShoppingCart, path: '/purchase', type: 'default' },
  { label: '销售管理', icon: PieChart, path: '/sales', type: 'default' },
  { label: '产品管理', icon: Goods, path: '/product', type: 'default' },
  { label: '库存盘点', icon: Edit, path: '/inventory/check', type: 'default' },
  { label: '库存调拨', icon: RefreshRight, path: '/inventory/transfer', type: 'default' },
  { label: '库存预警', icon: Warning, path: '/inventory/alert', type: 'default' },
  { label: '供应商', icon: Shop, path: '/supplier', type: 'default' },
  { label: '客户管理', icon: Avatar, path: '/customer', type: 'default' },
  { label: '仓库管理', icon: OfficeBuilding, path: '/warehouse', type: 'default' },
  { label: '统计报表', icon: DataAnalysis, path: '/erp', type: 'default' },
  { label: '用户管理', icon: User, path: '/user', type: 'default' },
  { label: '系统设置', icon: Tools, path: '/menu', type: 'default' }
]

const formatAmount = (amount) => {
  const value = amount || 0
  if (value >= 10000) {
    return (value / 10000).toFixed(1) + '万'
  }
  return value.toLocaleString()
}

const loadOverview = async () => {
  try {
    const response = await erpApi.getOverview()
    if (response.code === 0) {
      overviewData.value = response.data
      // 从概览接口获取用户数量和产品数量
      userCount.value = response.data.user_count || 0
      productCount.value = response.data.product_count || 0
    }
  } catch (error) {
    console.error('加载概览数据失败:', error)
  }
}

const loadRecentActivities = async () => {
  try {
    const response = await activityApi.getActivityList({ page: 1, page_size: 8 })
    if (response.code === 0) {
      const activities = response.data.activities || []
      recentActivities.value = activities.map((activity, index) => ({
        id: index,
        time: new Date(activity.created_at || activity.CreatedAt).toLocaleString('zh-CN', {
          month: '2-digit',
          day: '2-digit',
          hour: '2-digit',
          minute: '2-digit'
        }),
        username: activity.username || activity.Username || '未知用户',
        action: activity.action || activity.Action || '',
        url: activity.url || activity.URL || '',
        ip: activity.ip || activity.IP || '-'
      }))
    }
  } catch (error) {
    console.error('加载活动失败:', error)
  }
}

const loadInventoryAlerts = async () => {
  try {
    const response = await erpApi.getInventoryAlert()
    if (response.code === 0) {
      inventoryAlertCount.value = response.data.length
      inventoryAlerts.value = response.data.slice(0, 5)
    }
  } catch (error) {
    console.error('加载库存预警失败:', error)
  }
}

const getAlertQuantityClass = (alert) => {
  if (alert.alert_type === 2) return 'low-quantity critical'
  if (alert.alert_type === 1) return 'low-quantity warning'
  if (alert.alert_type === 3) return 'high-quantity'
  return 'normal-quantity'
}

const getAlertTypeText = (alert) => {
  if (alert.alert_type === 1) return `低于安全库存(${alert.safety_stock})`
  if (alert.alert_type === 2) return `低于最低库存(${alert.min_stock})`
  if (alert.alert_type === 3) return `高于最高库存(${alert.max_stock})`
  return '库存异常'
}

const getAlertGap = (alert) => {
  if (alert.alert_type === 1) return alert.safety_stock - alert.quantity
  if (alert.alert_type === 2) return alert.min_stock - alert.quantity
  if (alert.alert_type === 3) return alert.quantity - alert.max_stock
  return 0
}

const loadOrderStatus = async () => {
  try {
    const response = await erpApi.getOrderStatus()
    if (response.code === 0) {
      orderStatus.value = response.data
    }
  } catch (error) {
    console.error('加载订单状态失败:', error)
  }
}

const loadTrend = async () => {
  try {
    const response = await erpApi.getTrend({ days: trendPeriod.value })
    if (response.code === 0) {
      renderTrendChart(response.data)
    }
  } catch (error) {
    console.error('加载趋势数据失败:', error)
  }
}

const renderTrendChart = (data) => {
  nextTick(() => {
    if (!trendChart.value) return
    
    // 如果已存在实例，先销毁
    if (trendChartInstance) {
      trendChartInstance.dispose()
      trendChartInstance = null
    }
    
    // 移除旧的 resize 监听器
    if (resizeHandler) {
      window.removeEventListener('resize', resizeHandler)
    }
    
    trendChartInstance = echarts.init(trendChart.value)
    const option = {
      tooltip: {
        trigger: 'axis',
        backgroundColor: 'rgba(255, 255, 255, 0.95)',
        borderColor: '#ebeef5',
        borderWidth: 1,
        textStyle: { color: '#606266' },
        formatter: (params) => {
          let result = `<div style="font-weight:bold;margin-bottom:8px;">${params[0].axisValue}</div>`
          params.forEach(item => {
            result += `<div style="display:flex;align-items:center;margin:4px 0;">
              <span style="display:inline-block;width:10px;height:10px;border-radius:50%;background:${item.color};margin-right:8px;"></span>
              <span>${item.seriesName}: ¥${item.value.toLocaleString()}</span>
            </div>`
          })
          return result
        }
      },
      legend: {
        data: ['采购金额', '销售金额'],
        top: 0,
        right: '5%',
        textStyle: { color: '#909399', fontSize: 12 }
      },
      grid: {
        left: '5%',
        right: '5%',
        top: '20%',
        bottom: '25%',
        containLabel: true
      },
      xAxis: {
        type: 'category',
        data: data.dates,
        axisLine: { lineStyle: { color: '#ebeef5' } },
        axisLabel: { color: '#909399', rotate: 30, fontSize: 10 }
      },
      yAxis: {
        type: 'value',
        axisLine: { show: false },
        axisTick: { show: false },
        splitLine: { lineStyle: { color: '#f0f0f0' } },
        axisLabel: { color: '#909399' }
      },
      series: [
        {
          name: '采购金额',
          type: 'line',
          data: data.purchase_data,
          smooth: true,
          areaStyle: {
            color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
              { offset: 0, color: 'rgba(64, 158, 255, 0.3)' },
              { offset: 1, color: 'rgba(64, 158, 255, 0.05)' }
            ])
          },
          lineStyle: { color: '#409EFF', width: 2 },
          itemStyle: { color: '#409EFF' }
        },
        {
          name: '销售金额',
          type: 'line',
          data: data.sales_data,
          smooth: true,
          areaStyle: {
            color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
              { offset: 0, color: 'rgba(103, 194, 58, 0.3)' },
              { offset: 1, color: 'rgba(103, 194, 58, 0.05)' }
            ])
          },
          lineStyle: { color: '#67C23A', width: 2 },
          itemStyle: { color: '#67C23A' }
        }
      ]
    }
    trendChartInstance.setOption(option)
    
    // 添加新的 resize 监听器
    resizeHandler = () => {
      if (trendChartInstance) {
        trendChartInstance.resize()
      }
    }
    window.addEventListener('resize', resizeHandler)
  })
}

// 组件卸载前销毁图表实例
onBeforeUnmount(() => {
  if (trendChartInstance) {
    trendChartInstance.dispose()
    trendChartInstance = null
  }
  if (resizeHandler) {
    window.removeEventListener('resize', resizeHandler)
    resizeHandler = null
  }
})

const navigate = (path) => {
  router.push(path)
}

onMounted(() => {
  const userStr = localStorage.getItem('user')
  if (userStr) {
    const user = JSON.parse(userStr)
    username.value = user.username
  }

  loadOverview()
  loadRecentActivities()
  loadInventoryAlerts()
  loadOrderStatus()
  loadTrend()
})
</script>

<style scoped>
.dashboard {
  padding: 20px;
  min-height: calc(100vh - 100px);
  background: #f5f7fa;
}

.welcome-card {
  margin-bottom: 20px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  border: none;
  border-radius: 12px;
}

.welcome-card .card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  color: rgba(255, 255, 255, 0.8);
}

.current-time {
  font-size: 14px;
  font-weight: 500;
}

.welcome-content {
  display: flex;
  flex-direction: row;
  align-items: center;
  padding: 20px;
  gap: 20px;
}

.avatar-wrapper {
  position: relative;
}

.status-dot {
  position: absolute;
  bottom: 5px;
  right: 5px;
  width: 14px;
  height: 14px;
  border-radius: 50%;
  background: #67C23A;
  border: 2px solid white;
}

.welcome-content h2 {
  margin: 0 0 8px 0;
  font-size: 20px;
  font-weight: bold;
}

.welcome-text {
  color: rgba(255, 255, 255, 0.8);
  font-size: 14px;
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(6, 1fr);
  gap: 16px;
  margin-bottom: 20px;
}

.stat-card {
  border-radius: 12px;
  border: 1px solid #ebeef5;
  background: #ffffff;
  transition: all 0.3s ease;
  padding: 20px;
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.stat-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 10px 20px rgba(0, 0, 0, 0.08);
}

.stat-card-header {
  display: flex;
  align-items: center;
  gap: 12px;
}

.stat-icon-box {
  width: 40px;
  height: 40px;
  border-radius: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.stat-icon-purchase {
  background: rgba(44, 125, 52, 0.1);
}
.stat-icon-purchase .stat-icon {
  color: #2c7d34;
}

.stat-icon-sales {
  background: rgba(24, 144, 255, 0.1);
}
.stat-icon-sales .stat-icon {
  color: #1890ff;
}

.stat-icon-inventory {
  background: rgba(114, 46, 209, 0.1);
}
.stat-icon-inventory .stat-icon {
  color: #722ed1;
}

.stat-icon-order {
  background: rgba(250, 140, 22, 0.1);
}
.stat-icon-order .stat-icon {
  color: #fa8c16;
}

.stat-icon-user {
  background: rgba(235, 47, 150, 0.1);
}
.stat-icon-user .stat-icon {
  color: #eb2f96;
}

.stat-icon-product {
  background: rgba(19, 194, 194, 0.1);
}
.stat-icon-product .stat-icon {
  color: #13c2c2;
}

.stat-icon {
  font-size: 18px;
}

.stat-title {
  font-size: 13px;
  color: #909399;
}

.stat-number {
  font-size: 24px;
  font-weight: 600;
  color: #303133;
  line-height: 1.2;
}

.stat-number .stat-unit {
  font-size: 12px;
  font-weight: 400;
  color: #909399;
  margin-left: 4px;
}

.stat-change {
  font-size: 12px;
  margin-top: 4px;
}

.stat-change.positive {
  color: #67c23a;
}

.stat-change.negative {
  color: #f56c6c;
}

.main-content {
  display: flex;
  gap: 20px;
}

.left-panel {
  flex: 2;
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.right-panel {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.chart-card, .activity-card, .alert-card, .quick-actions-card, .order-summary-card {
  border-radius: 12px;
  border: 1px solid #ebeef5;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.period-select {
  width: 100px;
}

.chart-container {
  height: 320px;
}

.activity-card {
  max-height: 350px;
  overflow-y: auto;
}

.activity-card::-webkit-scrollbar {
  width: 6px;
}

.activity-card::-webkit-scrollbar-track {
  background: #f1f1f1;
  border-radius: 3px;
}

.activity-card::-webkit-scrollbar-thumb {
  background: #c1c1c1;
  border-radius: 3px;
}

.activity-card::-webkit-scrollbar-thumb:hover {
  background: #a1a1a1;
}

.activity-list {
  padding: 8px;
}

.activity-list-header {
  display: flex;
  padding: 10px 8px;
  background: #f5f7fa;
  border-radius: 4px;
  font-weight: 600;
  font-size: 13px;
  color: #606266;
  margin-bottom: 6px;
}

.activity-list-item {
  display: flex;
  padding: 10px 8px;
  border-bottom: 1px solid #ebeef5;
  font-size: 13px;
  align-items: center;
}

.activity-list-item:last-child {
  border-bottom: none;
}

.activity-list-item:hover {
  background: #fafafa;
}

.col-time {
  flex: 0 0 18%;
  max-width: 100px;
  min-width: 70px;
  color: #909399;
}

.col-user {
  flex: 0 0 15%;
  max-width: 90px;
  min-width: 60px;
  color: #303133;
  font-weight: 500;
}

.col-action {
  flex: 1;
  min-width: 0;
  color: #606266;
  word-break: break-all;
  padding: 0 12px;
}

.col-url {
  flex: 0 0 20%;
  max-width: 140px;
  min-width: 80px;
  color: #409EFF;
  word-break: break-all;
}

.col-ip {
  flex: 0 0 15%;
  max-width: 90px;
  min-width: 70px;
  color: #909399;
}

.empty-activity {
  text-align: center;
  color: #909399;
  padding: 20px;
  font-size: 13px;
}

.alert-card {
  border-left: 4px solid #F56C6C;
}

.alert-icon {
  color: #F56C6C;
  font-size: 18px;
  margin-right: 8px;
}

.badge {
  background: #F56C6C;
}

.alert-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.alert-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 12px;
  background: #fef0f0;
  border-radius: 8px;
}

.alert-product {
  font-weight: 600;
  color: #303133;
}

.alert-quantity {
  font-size: 13px;
  color: #606266;
}

.low-quantity {
  color: #F56C6C;
  font-weight: bold;
}

.low-quantity.critical {
  color: #F56C6C;
}

.low-quantity.warning {
  color: #E6A23C;
}

.high-quantity {
  color: #67C23A;
  font-weight: bold;
}

.normal-quantity {
  color: #909399;
  font-weight: bold;
}

.alert-threshold {
  font-size: 12px;
  color: #909399;
}

.gap-value {
  color: #E6A23C;
  font-weight: bold;
}

.quick-actions {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 12px;
}

.quick-actions :deep(.el-button) {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 6px;
  padding: 16px 12px;
  border-radius: 8px;
  transition: all 0.3s ease;
  width: 100%;
  min-width: 0;
  margin: 0;
  box-sizing: border-box;
  border: 1px solid #ebeef5;
  background: #ffffff;
  color: #606266;
}

.quick-actions :deep(.el-button:hover) {
  background: #f5f7fa;
  border-color: #d9d9d9;
  transform: translateY(-1px);
}

.quick-actions :deep(.el-button) i {
  font-size: 18px;
  color: #409eff;
}

.order-summary {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 8px;
}

.order-item {
  text-align: center;
  padding: 14px 8px;
  border-radius: 8px;
  background: #fafafa;
}

.order-item.pending .order-count {
  color: #fa8c16;
}

.order-item.approved .order-count {
  color: #1890ff;
}

.order-item.completed .order-count {
  color: #67c23a;
}

.order-item.cancelled .order-count {
  color: #f56c6c;
}

.order-count {
  font-size: 22px;
  font-weight: bold;
  color: #303133;
}

.order-label {
  font-size: 12px;
  color: #909399;
  margin-top: 4px;
}

@media (max-width: 1200px) {
  .stats-grid {
    grid-template-columns: repeat(3, 1fr);
  }
  
  .main-content {
    flex-direction: column;
  }
}

@media (max-width: 768px) {
  .stats-grid {
    grid-template-columns: repeat(2, 1fr);
  }
  
  .quick-actions {
    grid-template-columns: repeat(2, 1fr);
  }
  
  .order-summary {
    grid-template-columns: repeat(2, 1fr);
  }
}
</style>