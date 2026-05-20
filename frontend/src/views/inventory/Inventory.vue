<template>
  <div class="inventory-management">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>库存管理</span>
          <div>
            <el-button type="primary" @click="goToAdjustRequest" v-has-permission="'btn_inventory_adjust'">库存调整</el-button>
          </div>
        </div>
      </template>
      <div class="description-box">
        <span class="description-icon">ℹ️</span>
        <span class="description-text">库存管理：展示各仓库中所有产品的实时库存数量，支持按产品名称和仓库筛选。点击"调整"按钮可发起库存调整申请。</span>
      </div>
      <div class="query-bar">
        <el-input
          v-model="searchQuery"
          placeholder="请输入产品名称"
          style="width: 300px"
          prefix-icon="el-icon-search"
        ></el-input>
        <el-select v-model="warehouseFilter" placeholder="请选择仓库" clearable style="width: 150px">
          <el-option :value="0" label="全部仓库"></el-option>
          <el-option v-for="wh in warehouses" :key="wh.id" :label="wh.name" :value="wh.id"></el-option>
        </el-select>
        <el-button type="primary" @click="handleSearch">搜索</el-button>
      </div>
      <el-table :data="inventoryList" border stripe v-loading="loading">
        <el-table-column prop="id" label="ID" width="60"></el-table-column>
        <el-table-column prop="product.name" label="产品名称" min-width="120">
          <template #default="{ row }">
            {{ row.product?.name || '-' }}
          </template>
        </el-table-column>
        <el-table-column prop="product.code" label="产品编码" min-width="120">
          <template #default="{ row }">
            {{ row.product?.code || '-' }}
          </template>
        </el-table-column>
        <el-table-column prop="warehouse.name" label="仓库" width="100">
          <template #default="{ row }">
            {{ row.warehouse?.name || '-' }}
          </template>
        </el-table-column>
        <el-table-column prop="quantity" label="库存数量" width="100">
          <template #default="{ row }">
            <span :class="getQuantityClass(row.quantity, row.type)">
              {{ row.quantity }}
            </span>
          </template>
        </el-table-column>
        <el-table-column prop="product.main_unit" label="单位" width="80"></el-table-column>
        <el-table-column prop="product.price" label="单价" width="100">
          <template #default="{ row }">
            ¥{{ row.product?.price?.toFixed(2) || '-' }}
          </template>
        </el-table-column>
        <el-table-column prop="type" label="类型" width="80">
          <template #default="{ row }">
            <el-tag :type="getTypeTagType(row.type)">{{ getTypeText(row.type) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="order_type" label="订单类型" width="100">
          <template #default="{ row }">
            {{ getOrderTypeText(row.order_type) }}
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="更新时间" width="160"></el-table-column>
        <el-table-column label="操作" width="130" fixed="right" align="center">
          <template #default="{ row }">
            <div class="action-buttons">
              <el-button size="small" @click="goToAdjustRequestWithParams(row)">调整</el-button>
              <el-button size="small" type="primary" @click="handleViewHistory(row)">查看记录</el-button>
            </div>
          </template>
        </el-table-column>
      </el-table>
      <el-pagination
        v-model:current-page="currentPage"
        v-model:page-size="pageSize"
        :total="total"
        :page-sizes="[10, 20, 50, 100]"
        layout="total, sizes, prev, pager, next, jumper"
        @size-change="handleSizeChange"
        @current-change="handleCurrentChange"
        style="margin-top: 20px"
      />
    </el-card>

    <el-dialog
      title="库存变动记录"
      v-model="historyVisible"
      width="800px"
      :fullscreen="isMobile"
    >
      <el-table :data="historyList" style="width: 100%" border>
        <el-table-column prop="id" label="ID" width="80"></el-table-column>
        <el-table-column prop="type" label="变动类型" width="100">
          <template #default="scope">
            <el-tag :type="scope.row.type === 1 ? 'success' : 'danger'">
              {{ scope.row.type === 1 ? '入库' : '出库' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="quantity" label="变动数量" width="120">
          <template #default="scope">
            <span :class="scope.row.quantity >= 0 ? 'text-success' : 'text-danger'">
              {{ scope.row.quantity >= 0 ? '+' : '' }}{{ scope.row.quantity }}
            </span>
          </template>
        </el-table-column>
        <el-table-column prop="before_quantity" label="变动前数量" width="120"></el-table-column>
        <el-table-column prop="after_quantity" label="变动后数量" width="120"></el-table-column>
        <el-table-column prop="remark" label="变动原因" min-width="150"></el-table-column>
        <el-table-column prop="created_at" label="变动时间" width="160"></el-table-column>
      </el-table>
      <div v-if="historyList.length === 0" class="empty-history">
        暂无变动记录
      </div>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { inventoryApi } from '../../api/inventory'
import { warehouseApi } from '../../api/warehouse'

const router = useRouter()

const loading = ref(false)
const inventoryList = ref([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(10)
const searchQuery = ref('')
const warehouseFilter = ref(0)

const historyVisible = ref(false)
const historyList = ref([])

const warehouses = ref([])

const isMobile = computed(() => {
  if (typeof window !== 'undefined') {
    return window.innerWidth < 768
  }
  return false
})

const goToAdjustRequest = () => {
  router.push({
    path: '/inventory/adjust-request',
    query: {
      back: 'true'
    }
  })
}

const goToAdjustRequestWithParams = (row) => {
  router.push({
    path: '/inventory/adjust-request',
    query: {
      productId: row.productId,
      warehouseId: row.warehouseId,
      back: 'true'
    }
  })
}

const getQuantityClass = (quantity, type) => {
  if (type === 1) return 'text-success'
  if (type === 2) return 'text-danger'
  if (quantity <= 10) return 'quantity-low'
  if (quantity <= 50) return 'quantity-medium'
  return 'quantity-normal'
}



const handleViewHistory = (item) => {
  inventoryApi.getInventoryHistory({
    product_id: item.product_id,
    warehouse_id: item.warehouse_id
  })
  .then(response => {
    if (response.code === 0) {
      historyList.value = response.data?.records || []
    } else {
      ElMessage.error('获取库存记录失败：' + response.message)
      historyList.value = []
    }
    historyVisible.value = true
  })
  .catch(error => {
    ElMessage.error('获取库存记录失败')
    historyList.value = []
    historyVisible.value = true
  })
}

const loadInventory = async () => {
  loading.value = true
  try {
    const response = await inventoryApi.getInventoryList({
    page: currentPage.value,
    page_size: pageSize.value,
    product_name: searchQuery.value,
    warehouse_id: warehouseFilter.value
  })
    if (response.code === 0) {
      inventoryList.value = response.data.inventory || []
      total.value = response.data.total || 0
    } else {
      ElMessage.error('获取库存失败：' + response.message)
    }
  } catch (error) {
    ElMessage.error('获取库存失败')
  } finally {
    loading.value = false
  }
}

const loadWarehouses = () => {
  warehouseApi.getActiveWarehouseList({})
  .then(response => {
    if (response.code === 0) {
      warehouses.value = response.data.warehouses || []
    }
  })
}

const getTypeTagType = (type) => {
  const types = {
    1: 'success',
    2: 'danger',
    3: 'warning'
  }
  return types[type] || 'info'
}

const getMainUnit = (product) => {
  if (!product || !product.units) return product?.unit || ''
  const mainUnit = product.units.find(u => u.isMain === 1)
  return mainUnit ? mainUnit.unitName : (product.units[0]?.unitName || product?.unit || '')
}

const getTypeText = (type) => {
  const texts = {
    1: '入库',
    2: '出库',
    3: '调整'
  }
  return texts[type] || '未知'
}

const getOrderTypeText = (orderType) => {
  const texts = {
    1: '采购订单',
    2: '销售订单',
    3: '库存调整',
    4: '调整申请'
  }
  return texts[orderType] || '-'
}

const handleSearch = () => {
  currentPage.value = 1
  loadInventory()
}

const handleSizeChange = (size) => {
  pageSize.value = size
  loadInventory()
}

const handleCurrentChange = (current) => {
  currentPage.value = current
  loadInventory()
}

onMounted(() => {
  loadInventory()
  loadWarehouses()
})
</script>

<style scoped>
.inventory-management {
  padding: 20px;
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

.query-bar {
  margin-bottom: 20px;
  display: flex;
  align-items: center;
  gap: 10px;
}

.quantity-low {
  color: #f56c6c;
  font-weight: 500;
}

.quantity-medium {
  color: #e6a23c;
  font-weight: 500;
}

.quantity-normal {
  color: #67c23a;
  font-weight: 500;
}

.empty-history {
  text-align: center;
  padding: 40px;
  color: #909399;
}

.text-success {
  color: #67c23a;
}

.text-danger {
  color: #f56c6c;
}

.action-buttons {
  display: flex;
  justify-content: center;
  align-items: center;
  gap: 12px;
}

.action-buttons .el-button {
  margin-bottom: 0;
}
</style>