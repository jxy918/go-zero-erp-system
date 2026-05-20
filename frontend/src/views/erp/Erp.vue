<template>
  <div class="erp-statistics">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>ERP管理统计报表</span>
          <el-button type="primary" size="small" @click="loadAllData">
            <el-icon><Refresh /></el-icon>
            刷新数据
          </el-button>
        </div>
      </template>
      <div class="description-box">
        <span class="description-icon">ℹ️</span>
        <span class="description-text">统计报表：展示采购、销售、库存、财务等业务数据的统计汇总，包括金额趋势、订单状态、库存预警等分析。</span>
      </div>

      <div class="section-title">
        <el-icon><DataAnalysis /></el-icon>
        <span>核心指标</span>
      </div>
      <div class="overview-cards">
        <div class="overview-card" @click="goToModule('/purchase')">
          <div class="card-icon purchase-icon">
            <el-icon><ShoppingCart /></el-icon>
          </div>
          <div class="card-content">
            <div class="card-label">本月采购金额</div>
            <div class="card-value">¥{{ formatAmount(overviewData.purchase_amount) }}</div>
          </div>
        </div>
        <div class="overview-card" @click="goToModule('/sales')">
          <div class="card-icon sales-icon">
            <el-icon><Sell /></el-icon>
          </div>
          <div class="card-content">
            <div class="card-label">本月销售金额</div>
            <div class="card-value">¥{{ formatAmount(overviewData.sales_amount) }}</div>
          </div>
        </div>
        <div class="overview-card" @click="goToModule('/inventory')">
          <div class="card-icon inventory-icon">
            <el-icon><Box /></el-icon>
          </div>
          <div class="card-content">
            <div class="card-label">当前库存总量</div>
            <div class="card-value">{{ overviewData.total_inventory }} 件</div>
          </div>
        </div>
        <div class="overview-card" @click="goToModule('/inventory/alert')">
          <div class="card-icon alert-icon">
            <el-icon><Warning /></el-icon>
          </div>
          <div class="card-content">
            <div class="card-label">库存预警数</div>
            <div class="card-value warning-value">{{ inventoryAlertCount }} 条</div>
          </div>
        </div>
      </div>

      <div class="section-title">
        <el-icon><TrendCharts /></el-icon>
        <span>采购销售分析</span>
      </div>
      <div class="charts-row">
        <el-card class="chart-card" style="width: 100%">
          <template #header>
            <div class="card-header-with-action">
              <span>采购/销售趋势（近30天）</span>
              <el-radio-group v-model="trendPeriod" size="small" @change="loadTrend">
                <el-radio-button :value="7" label="近7天">近7天</el-radio-button>
                <el-radio-button :value="30" label="近30天">近30天</el-radio-button>
              </el-radio-group>
            </div>
          </template>
          <div ref="trendChart" class="chart-container"></div>
        </el-card>
      </div>

      <div class="section-title">
        <el-icon><DataLine /></el-icon>
        <span>订单分析</span>
      </div>
      <div class="two-col-row">
        <el-card class="chart-card" style="width: 48%">
          <template #header>
            <div class="card-header-with-action">
              <span>订单状态分布</span>
              <el-radio-group v-model="orderType" size="small" @change="loadOrderStatus">
                <el-radio-button :value="'purchase'" label="采购订单">采购订单</el-radio-button>
                <el-radio-button :value="'sales'" label="销售订单">销售订单</el-radio-button>
              </el-radio-group>
            </div>
          </template>
          <div ref="statusChart" class="chart-container"></div>
        </el-card>

        <el-card class="chart-card" style="width: 48%">
          <template #header>
            <span>热销商品排行 TOP10</span>
          </template>
          <el-table :data="topProducts" style="width: 100%" border size="small">
            <el-table-column prop="rank" label="排名" width="60" align="center">
              <template #default="scope">
                <el-tag v-if="scope.row.rank <= 3" :type="getRankType(scope.row.rank)" size="small">
                  {{ scope.row.rank }}
                </el-tag>
                <span v-else>{{ scope.row.rank }}</span>
              </template>
            </el-table-column>
            <el-table-column prop="product_name" label="商品名称" min-width="120"></el-table-column>
            <el-table-column prop="quantity" label="销售数量" width="90" align="right"></el-table-column>
            <el-table-column prop="amount" label="销售额" width="100" align="right">
              <template #default="scope">
                <span class="amount-text">¥{{ formatAmount(scope.row.amount) }}</span>
              </template>
            </el-table-column>
          </el-table>
          <el-empty v-if="topProducts.length === 0" description="暂无销售数据" :image-size="60"></el-empty>
        </el-card>
      </div>

      <div class="section-title">
        <el-icon><Warning /></el-icon>
        <span>库存预警</span>
      </div>
      <div class="two-col-row">
        <el-card class="chart-card" style="width: 58%">
          <template #header>
            <div class="card-header-with-action">
              <span>库存预警列表</span>
              <el-button type="primary" size="small" link @click="goToModule('/inventory/alert')">
                查看全部预警
              </el-button>
            </div>
          </template>
          <el-table :data="inventoryAlerts" style="width: 100%" border size="small">
            <el-table-column prop="product_name" label="商品名称" min-width="120"></el-table-column>
            <el-table-column prop="product_code" label="商品编码" width="100"></el-table-column>
            <el-table-column prop="warehouse" label="仓库" width="80"></el-table-column>
            <el-table-column prop="quantity" label="当前库存" width="90" align="right">
              <template #default="scope">
                <span :class="getStockClass(scope.row)">{{ scope.row.quantity }}</span>
              </template>
            </el-table-column>
            <el-table-column label="预警类型" width="100">
              <template #default="scope">
                <el-tag v-if="scope.row.alert_type === 1" type="warning" size="small">低于安全库存</el-tag>
                <el-tag v-else-if="scope.row.alert_type === 2" type="danger" size="small">低于最低库存</el-tag>
                <el-tag v-else-if="scope.row.alert_type === 3" type="info" size="small">高于最高库存</el-tag>
              </template>
            </el-table-column>
            <el-table-column label="预警级别" width="80" align="center">
              <template #default="scope">
                <el-tag v-if="scope.row.alert_level === 3" type="danger" size="small">高</el-tag>
                <el-tag v-else-if="scope.row.alert_level === 2" type="warning" size="small">中</el-tag>
                <el-tag v-else type="info" size="small">低</el-tag>
              </template>
            </el-table-column>
          </el-table>
          <el-empty v-if="inventoryAlerts.length === 0" description="暂无库存预警" :image-size="60"></el-empty>
        </el-card>

        <el-card class="chart-card" style="width: 38%">
          <template #header>
            <span>库存预警统计</span>
          </template>
          <div class="alert-stats">
            <div class="alert-stat-item danger">
              <div class="alert-stat-icon"><el-icon><WarningFilled /></el-icon></div>
              <div class="alert-stat-content">
                <div class="alert-stat-label">紧急预警</div>
                <div class="alert-stat-value">{{ alertStats.critical }}</div>
              </div>
            </div>
            <div class="alert-stat-item warning">
              <div class="alert-stat-icon"><el-icon><Warning /></el-icon></div>
              <div class="alert-stat-content">
                <div class="alert-stat-label">中度预警</div>
                <div class="alert-stat-value">{{ alertStats.warning }}</div>
              </div>
            </div>
            <div class="alert-stat-item info">
              <div class="alert-stat-icon"><el-icon><InfoFilled /></el-icon></div>
              <div class="alert-stat-content">
                <div class="alert-stat-label">轻度预警</div>
                <div class="alert-stat-value">{{ alertStats.info }}</div>
              </div>
            </div>
          </div>
          <div ref="alertChart" class="chart-container-small"></div>
        </el-card>
      </div>

      <div class="section-title">
        <el-icon><OfficeBuilding /></el-icon>
        <span>业务实体统计</span>
      </div>
      <div class="two-col-row">
        <el-card class="chart-card" style="width: 48%">
          <template #header>
            <div class="card-header-with-action">
              <span>供应商/客户统计</span>
              <el-button type="primary" size="small" link @click="goToModule('/supplier')">
                供应商管理
              </el-button>
              <el-button type="primary" size="small" link @click="goToModule('/customer')">
                客户管理
              </el-button>
            </div>
          </template>
          <el-row :gutter="20">
            <el-col :span="12">
              <div class="business-item" @click="goToModule('/supplier')">
                <div class="business-icon supplier-icon">
                  <el-icon><Shop /></el-icon>
                </div>
                <div class="business-info">
                  <div class="business-label">供应商总数</div>
                  <div class="business-value">{{ businessData.supplier_count }}</div>
                </div>
              </div>
            </el-col>
            <el-col :span="12">
              <div class="business-item" @click="goToModule('/customer')">
                <div class="business-icon customer-icon">
                  <el-icon><User /></el-icon>
                </div>
                <div class="business-info">
                  <div class="business-label">客户总数</div>
                  <div class="business-value">{{ businessData.customer_count }}</div>
                </div>
              </div>
            </el-col>
            <el-col :span="12">
              <div class="business-item">
                <div class="business-icon active-supplier-icon">
                  <el-icon><DataBoard /></el-icon>
                </div>
                <div class="business-info">
                  <div class="business-label">本月活跃供应商</div>
                  <div class="business-value">{{ businessData.active_supplier_count }}</div>
                </div>
              </div>
            </el-col>
            <el-col :span="12">
              <div class="business-item">
                <div class="business-icon active-customer-icon">
                  <el-icon><TrendCharts /></el-icon>
                </div>
                <div class="business-info">
                  <div class="business-label">本月活跃客户</div>
                  <div class="business-value">{{ businessData.active_customer_count }}</div>
                </div>
              </div>
            </el-col>
          </el-row>
        </el-card>

        <el-card class="chart-card" style="width: 48%">
          <template #header>
            <div class="card-header-with-action">
              <span>产品/仓库统计</span>
              <el-button type="primary" size="small" link @click="goToModule('/product')">
                产品管理
              </el-button>
              <el-button type="primary" size="small" link @click="goToModule('/warehouse')">
                仓库管理
              </el-button>
            </div>
          </template>
          <el-row :gutter="20">
            <el-col :span="12">
              <div class="business-item" @click="goToModule('/product')">
                <div class="business-icon product-icon">
                  <el-icon><Goods /></el-icon>
                </div>
                <div class="business-info">
                  <div class="business-label">产品总数</div>
                  <div class="business-value">{{ businessData.product_count }}</div>
                </div>
              </div>
            </el-col>
            <el-col :span="12">
              <div class="business-item" @click="goToModule('/warehouse')">
                <div class="business-icon warehouse-icon">
                  <el-icon><OfficeBuilding /></el-icon>
                </div>
                <div class="business-info">
                  <div class="business-label">仓库总数</div>
                  <div class="business-value">{{ businessData.warehouse_count }}</div>
                </div>
              </div>
            </el-col>
            <el-col :span="12">
              <div class="business-item">
                <div class="business-icon low-stock-icon">
                  <el-icon><Warning /></el-icon>
                </div>
                <div class="business-info">
                  <div class="business-label">库存不足产品</div>
                  <div class="business-value warning-value">{{ businessData.low_stock_product_count }}</div>
                </div>
              </div>
            </el-col>
            <el-col :span="12">
              <div class="business-item">
                <div class="business-icon order-icon">
                  <el-icon><Document /></el-icon>
                </div>
                <div class="business-info">
                  <div class="business-label">本月订单总数</div>
                  <div class="business-value">{{ overviewData.total_orders }}</div>
                </div>
              </div>
            </el-col>
          </el-row>
        </el-card>
      </div>

      <div class="section-title">
        <el-icon><Clock /></el-icon>
        <span>待处理事项</span>
      </div>
      <div class="todo-row">
        <el-card class="todo-card" @click="goToModule('/inventory/check')">
          <div class="todo-item">
            <div class="todo-icon check-icon">
              <el-icon><Edit /></el-icon>
            </div>
            <div class="todo-content">
              <div class="todo-label">待盘点单</div>
              <div class="todo-value">{{ todoData.pending_check_orders }} 单</div>
            </div>
          </div>
        </el-card>
        <el-card class="todo-card" @click="goToModule('/inventory/adjust-request')">
          <div class="todo-item">
            <div class="todo-icon adjust-icon">
              <el-icon><Setting /></el-icon>
            </div>
            <div class="todo-content">
              <div class="todo-label">待处理调整</div>
              <div class="todo-value">{{ todoData.pending_adjust_requests }} 单</div>
            </div>
          </div>
        </el-card>
        <el-card class="todo-card" @click="goToModule('/inventory/transfer')">
          <div class="todo-item">
            <div class="todo-icon transfer-icon">
              <el-icon><RefreshRight /></el-icon>
            </div>
            <div class="todo-content">
              <div class="todo-label">调拨待审核</div>
              <div class="todo-value">{{ todoData.pending_transfers }} 单</div>
            </div>
          </div>
        </el-card>
        <el-card class="todo-card" @click="goToModule('/purchase')">
          <div class="todo-item">
            <div class="todo-icon purchase-icon">
              <el-icon><ShoppingCart /></el-icon>
            </div>
            <div class="todo-content">
              <div class="todo-label">采购待审核</div>
              <div class="todo-value">{{ todoData.pending_purchases }} 单</div>
            </div>
          </div>
        </el-card>
        <el-card class="todo-card" @click="goToModule('/sales')">
          <div class="todo-item">
            <div class="todo-icon sales-icon">
              <el-icon><Sell /></el-icon>
            </div>
            <div class="todo-content">
              <div class="todo-label">销售待审核</div>
              <div class="todo-value">{{ todoData.pending_sales }} 单</div>
            </div>
          </div>
        </el-card>
      </div>
    </el-card>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted, nextTick } from 'vue'
import { useRouter } from 'vue-router'
import {
  ShoppingCart, Sell, Box, Warning, TrendCharts, DataLine,
  WarningFilled, InfoFilled, Shop, User, DataBoard, Goods,
  OfficeBuilding, Document, Clock, Edit, Setting, RefreshRight,
  DataAnalysis, Refresh
} from '@element-plus/icons-vue'
import * as echarts from 'echarts'
import { erpApi } from '../../api/erp'

const router = useRouter()

const trendPeriod = ref('30')
const orderType = ref('purchase')

const overviewData = ref({
  purchase_amount: 0,
  sales_amount: 0,
  total_inventory: 0,
  total_orders: 0,
  user_count: 0,
  product_count: 0
})

const inventoryAlertCount = ref(0)
const inventoryAlerts = ref([])
const alertStats = ref({
  critical: 0,
  warning: 0,
  info: 0
})

const topProducts = ref([])
const businessData = ref({
  supplier_count: 0,
  customer_count: 0,
  active_supplier_count: 0,
  active_customer_count: 0,
  product_count: 0,
  warehouse_count: 0,
  low_stock_product_count: 0
})

const todoData = ref({
  pending_check_orders: 0,
  pending_adjust_requests: 0,
  pending_transfers: 0,
  pending_purchases: 0,
  pending_sales: 0
})

const trendChart = ref(null)
const statusChart = ref(null)
const alertChart = ref(null)

// ECharts 实例引用
let trendChartInstance = null
let statusChartInstance = null
let alertChartInstance = null

const formatAmount = (amount) => {
  const value = amount || 0
  if (value >= 100000000) {
    return (value / 100000000).toFixed(2) + '亿'
  } else if (value >= 10000) {
    return (value / 10000).toFixed(2) + '万'
  }
  return value.toFixed(2)
}

const getRankType = (rank) => {
  switch (rank) {
    case 1: return 'danger'
    case 2: return 'warning'
    case 3: return 'success'
    default: return ''
  }
}

const getStockClass = (row) => {
  if (row.alert_level === 3) return 'stock-danger'
  if (row.alert_level === 2) return 'stock-warning'
  return ''
}

const goToModule = (path) => {
  router.push(path)
}

const loadAllData = async () => {
  await Promise.all([
    loadOverview(),
    loadTrend(),
    loadInventoryAlerts(),
    loadTopProducts(),
    loadOrderStatus(),
    loadBusiness(),
    loadTodoData()
  ])
}

const loadOverview = async () => {
  try {
    const response = await erpApi.getOverview()
    if (response.code === 0) {
      overviewData.value = response.data
    }
  } catch (error) {
    console.error('加载概览数据失败:', error)
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

const loadInventoryAlerts = async () => {
  try {
    const response = await erpApi.getInventoryAlert()
    if (response.code === 0) {
      inventoryAlerts.value = response.data.slice(0, 10)
      inventoryAlertCount.value = response.data.length
      alertStats.value = {
        critical: response.data.filter(item => item.alert_level === 3).length,
        warning: response.data.filter(item => item.alert_level === 2).length,
        info: response.data.filter(item => item.alert_level === 1).length
      }
      renderAlertChart(alertStats.value)
    }
  } catch (error) {
    console.error('加载库存预警失败:', error)
  }
}

const loadTopProducts = async () => {
  try {
    const response = await erpApi.getTopProducts()
    if (response.code === 0) {
      topProducts.value = response.data.slice(0, 10)
    }
  } catch (error) {
    console.error('加载热销商品失败:', error)
  }
}

const loadOrderStatus = async () => {
  try {
    const response = await erpApi.getOrderStatus({ type: orderType.value })
    if (response.code === 0) {
      renderStatusChart(response.data)
    }
  } catch (error) {
    console.error('加载订单状态失败:', error)
  }
}

const loadBusiness = async () => {
  try {
    const response = await erpApi.getBusiness()
    if (response.code === 0) {
      businessData.value = {
        ...response.data,
        product_count: response.data.product_count || 0,
        warehouse_count: response.data.warehouse_count || 0,
        low_stock_product_count: response.data.low_stock_product_count || 0
      }
    }
  } catch (error) {
    console.error('加载业务统计失败:', error)
  }
}

const loadTodoData = async () => {
  try {
    const response = await erpApi.getTodoData()
    if (response.code === 0) {
      todoData.value = response.data
    }
  } catch (error) {
    console.error('加载待处理数据失败:', error)
  }
}

const renderTrendChart = (data) => {
  nextTick(() => {
    if (!trendChart.value) return
    // 如果已有实例，先销毁
    if (trendChartInstance) {
      trendChartInstance.dispose()
      trendChartInstance = null
    }
    trendChartInstance = echarts.init(trendChart.value)
    const option = {
      tooltip: {
        trigger: 'axis'
      },
      legend: {
        data: ['采购金额', '销售金额'],
        top: 0,
        right: '5%'
      },
      grid: {
        left: '3%',
        right: '4%',
        bottom: '15%',
        top: '15%',
        containLabel: true
      },
      xAxis: {
        type: 'category',
        data: data.dates,
        axisLabel: {
          rotate: 45,
          fontSize: 10
        }
      },
      yAxis: {
        type: 'value',
        axisLabel: {
          formatter: (value) => {
            if (value >= 100000000) return (value / 100000000).toFixed(1) + '亿'
            if (value >= 10000) return (value / 10000).toFixed(1) + '万'
            return value
          }
        }
      },
      series: [
        {
          name: '采购金额',
          type: 'line',
          data: data.purchase_data,
          smooth: true,
          itemStyle: { color: '#409EFF' },
          areaStyle: { color: 'rgba(64, 158, 255, 0.1)' }
        },
        {
          name: '销售金额',
          type: 'line',
          data: data.sales_data,
          smooth: true,
          itemStyle: { color: '#67C23A' },
          areaStyle: { color: 'rgba(103, 194, 58, 0.1)' }
        }
      ]
    }
    trendChartInstance.setOption(option)
    window.addEventListener('resize', () => { if (trendChartInstance) trendChartInstance.resize() })
  })
}

const renderStatusChart = (data) => {
  nextTick(() => {
    if (!statusChart.value) return
    // 如果已有实例，先销毁
    if (statusChartInstance) {
      statusChartInstance.dispose()
      statusChartInstance = null
    }
    statusChartInstance = echarts.init(statusChart.value)
    const option = {
      tooltip: { trigger: 'item' },
      legend: { orient: 'vertical', left: 'left' },
      series: [{
        name: orderType.value === 'purchase' ? '采购订单' : '销售订单',
        type: 'pie',
        radius: ['40%', '70%'],
        avoidLabelOverlap: false,
        itemStyle: { borderRadius: 10, borderColor: '#fff', borderWidth: 2 },
        label: { show: true, formatter: '{b}: {c} ({d}%)' },
        emphasis: {
          label: { show: true, fontSize: 16, fontWeight: 'bold' }
        },
        data: [
          { value: data.pending, name: '待审核', itemStyle: { color: '#E6A23C' } },
          { value: data.approved, name: '已审核', itemStyle: { color: '#409EFF' } },
          { value: data.completed, name: '已完成', itemStyle: { color: '#67C23A' } },
          { value: data.cancelled, name: '已取消', itemStyle: { color: '#F56C6C' } }
        ]
      }]
    }
    statusChartInstance.setOption(option)
    window.addEventListener('resize', () => { if (statusChartInstance) statusChartInstance.resize() })
  })
}

const renderAlertChart = (data) => {
  nextTick(() => {
    if (!alertChart.value) return
    // 如果已有实例，先销毁
    if (alertChartInstance) {
      alertChartInstance.dispose()
      alertChartInstance = null
    }
    alertChartInstance = echarts.init(alertChart.value)
    const option = {
      tooltip: { trigger: 'item' },
      series: [{
        type: 'pie',
        radius: ['50%', '80%'],
        center: ['50%', '50%'],
        avoidLabelOverlap: false,
        label: { show: false },
        data: [
          { value: data.critical, name: '紧急', itemStyle: { color: '#F56C6C' } },
          { value: data.warning, name: '中度', itemStyle: { color: '#E6A23C' } },
          { value: data.info, name: '轻度', itemStyle: { color: '#909399' } }
        ]
      }]
    }
    alertChartInstance.setOption(option)
    window.addEventListener('resize', () => { if (alertChartInstance) alertChartInstance.resize() })
  })
}

onMounted(() => {
  loadAllData()
})

onUnmounted(() => {
  // 清理 ECharts 实例
  if (trendChartInstance) {
    trendChartInstance.dispose()
    trendChartInstance = null
  }
  if (statusChartInstance) {
    statusChartInstance.dispose()
    statusChartInstance = null
  }
  if (alertChartInstance) {
    alertChartInstance.dispose()
    alertChartInstance = null
  }
})
</script>

<style scoped>
.erp-statistics {
  padding: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.card-header-with-action {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.section-title {
  display: flex;
  align-items: center;
  font-size: 16px;
  font-weight: 600;
  color: #303133;
  margin: 24px 0 16px 0;
  padding-left: 10px;
  border-left: 4px solid #409EFF;
}

.section-title:first-of-type {
  margin-top: 0;
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

.overview-cards {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 20px;
  margin-bottom: 10px;
}

.overview-card {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 24px 16px;
  border-radius: 12px;
  border: 1px solid #ebeef5;
  background: #ffffff;
  cursor: pointer;
  transition: all 0.3s;
}

.overview-card:hover {
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
  transform: translateY(-2px);
}

.card-icon {
  width: 56px;
  height: 56px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-bottom: 12px;
  font-size: 24px;
}

.purchase-icon { background: linear-gradient(135deg, #667eea 0%, #764ba2 100%); color: white; }
.sales-icon { background: linear-gradient(135deg, #11998e 0%, #38ef7d 100%); color: white; }
.inventory-icon { background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%); color: white; }
.alert-icon { background: linear-gradient(135deg, #f5365c 0%, #f99c33 100%); color: white; }

.card-content {
  display: flex;
  flex-direction: column;
  align-items: center;
}

.card-label {
  font-size: 13px;
  color: #909399;
  margin-bottom: 6px;
}

.card-value {
  font-size: 22px;
  font-weight: 600;
  color: #303133;
}

.warning-value {
  color: #E6A23C;
}

.charts-row {
  margin-bottom: 10px;
}

.two-col-row {
  display: flex;
  justify-content: space-between;
  margin-bottom: 10px;
}

.chart-card {
  min-height: 320px;
}

.chart-container {
  height: 280px;
}

.chart-container-small {
  height: 160px;
}

.low-stock {
  color: #F56C6C;
  font-weight: bold;
}

.stock-danger {
  color: #F56C6C;
  font-weight: bold;
}

.stock-warning {
  color: #E6A23C;
  font-weight: bold;
}

.amount-text {
  color: #67C23A;
  font-weight: 500;
}

.alert-stats {
  display: flex;
  justify-content: space-around;
  margin-bottom: 15px;
}

.alert-stat-item {
  display: flex;
  align-items: center;
  padding: 10px 15px;
  border-radius: 8px;
  background: #f5f7fa;
}

.alert-stat-item.danger { border-left: 3px solid #F56C6C; }
.alert-stat-item.warning { border-left: 3px solid #E6A23C; }
.alert-stat-item.info { border-left: 3px solid #909399; }

.alert-stat-icon {
  font-size: 20px;
  margin-right: 10px;
}

.alert-stat-item.danger .alert-stat-icon { color: #F56C6C; }
.alert-stat-item.warning .alert-stat-icon { color: #E6A23C; }
.alert-stat-item.info .alert-stat-icon { color: #909399; }

.alert-stat-label {
  font-size: 12px;
  color: #909399;
}

.alert-stat-value {
  font-size: 20px;
  font-weight: bold;
  color: #303133;
}

.business-item {
  display: flex;
  align-items: center;
  padding: 15px;
  border-radius: 8px;
  background: #f5f7fa;
  margin-bottom: 15px;
  cursor: pointer;
  transition: all 0.3s;
}

.business-item:hover {
  background: #eef1f5;
  transform: translateX(5px);
}

.business-icon {
  width: 45px;
  height: 45px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-right: 12px;
  font-size: 20px;
  color: white;
}

.supplier-icon { background: linear-gradient(135deg, #667eea 0%, #764ba2 100%); }
.customer-icon { background: linear-gradient(135deg, #11998e 0%, #38ef7d 100%); }
.active-supplier-icon { background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%); }
.active-customer-icon { background: linear-gradient(135deg, #4facfe 0%, #00f2fe 100%); }
.product-icon { background: linear-gradient(135deg, #a18cd1 0%, #fbc2eb 100%); }
.warehouse-icon { background: linear-gradient(135deg, #ff9a9e 0%, #fecfef 100%); }
.low-stock-icon { background: linear-gradient(135deg, #ff6b35 0%, #f7b733 100%); }
.order-icon { background: linear-gradient(135deg, #4facfe 0%, #00f2fe 100%); }

.business-label {
  font-size: 12px;
  color: #999;
  margin-bottom: 3px;
}

.business-value {
  font-size: 20px;
  font-weight: bold;
  color: #303133;
}

.todo-row {
  display: grid;
  grid-template-columns: repeat(5, 1fr);
  gap: 15px;
  margin-bottom: 10px;
}

.todo-card {
  cursor: pointer;
  transition: all 0.3s;
}

.todo-card:hover {
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
  transform: translateY(-2px);
}

.todo-item {
  display: flex;
  align-items: center;
}

.todo-icon {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-right: 12px;
  font-size: 18px;
  color: white;
}

.check-icon { background: linear-gradient(135deg, #667eea 0%, #764ba2 100%); }
.adjust-icon { background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%); }
.transfer-icon { background: linear-gradient(135deg, #11998e 0%, #38ef7d 100%); }
.purchase-icon { background: linear-gradient(135deg, #4facfe 0%, #00f2fe 100%); }
.sales-icon { background: linear-gradient(135deg, #ff6b35 0%, #f7b733 100%); }

.todo-label {
  font-size: 12px;
  color: #909399;
}

.todo-value {
  font-size: 16px;
  font-weight: bold;
  color: #303133;
}
</style>
