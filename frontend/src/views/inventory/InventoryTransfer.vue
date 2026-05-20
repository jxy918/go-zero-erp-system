<template>
  <div class="inventory-transfer">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>库存调拨</span>
          <el-button type="primary" @click="handleCreate" v-has-permission="'inventory_transfer:create'">新建调拨单</el-button>
        </div>
      </template>
      <div class="description-box">
        <span class="description-icon">ℹ️</span>
        <span class="description-text">库存调拨：将产品从一个仓库转移到另一个仓库。流程：新建调拨单（选择源仓库、目标仓库、产品、数量）→ 提交审核 → 审核通过后自动执行调拨。</span>
      </div>
      <div class="query-bar">
        <el-input
          v-model="searchQuery"
          placeholder="请输入调拨单号"
          style="width: 200px"
          prefix-icon="el-icon-search"
          clearable
        ></el-input>
        <el-select v-model="fromWarehouseFilter" placeholder="源仓库" clearable style="width: 150px">
          <el-option v-for="wh in warehouses" :key="wh.id" :label="wh.name" :value="wh.id"></el-option>
        </el-select>
        <el-select v-model="toWarehouseFilter" placeholder="目标仓库" clearable style="width: 150px">
          <el-option v-for="wh in warehouses" :key="wh.id" :label="wh.name" :value="wh.id"></el-option>
        </el-select>
        <el-select v-model="productFilter" placeholder="产品" clearable filterable style="width: 150px">
          <el-option v-for="p in products" :key="p.id" :label="p.name" :value="p.id"></el-option>
        </el-select>
        <el-select v-model="statusFilter" placeholder="状态" clearable style="width: 150px">
          <el-option :value="0" label="全部"></el-option>
          <el-option :value="1" label="待审核"></el-option>
          <el-option :value="2" label="已审核"></el-option>
          <el-option :value="3" label="已完成"></el-option>
          <el-option :value="4" label="已拒绝"></el-option>
        </el-select>
        <el-button type="primary" @click="handleSearch">搜索</el-button>
        <el-button @click="handleReset">重置</el-button>
      </div>
      <el-table
        :data="transferOrders"
        style="width: 100%"
        stripe
        :fit="true"
        :cell-padding="12"
        :header-cell-padding="12"
      >
        <el-table-column prop="id" label="ID" width="60" align="center"></el-table-column>
        <el-table-column prop="transfer_no" label="调拨单号" min-width="160" show-overflow-tooltip></el-table-column>
        <el-table-column prop="from_warehouse" label="源仓库" min-width="120">
          <template #default="scope">
            {{ scope.row.from_warehouse?.name || '-' }}
          </template>
        </el-table-column>
        <el-table-column prop="to_warehouse" label="目标仓库" min-width="120">
          <template #default="scope">
            {{ scope.row.to_warehouse?.name || '-' }}
          </template>
        </el-table-column>
        <el-table-column prop="product" label="产品" min-width="140" show-overflow-tooltip>
          <template #default="scope">
            {{ scope.row.product?.name || '-' }}
          </template>
        </el-table-column>
        <el-table-column prop="quantity" label="调拨数量" width="100" align="right"></el-table-column>
        <el-table-column prop="status" label="状态" width="90" align="center">
          <template #default="scope">
            <el-tag :type="getStatusTagType(scope.row.status)" size="small">
              {{ getStatusText(scope.row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="remark" label="备注" min-width="130" show-overflow-tooltip></el-table-column>
        <el-table-column prop="created_at" label="创建时间" width="170"></el-table-column>
        <el-table-column label="操作" width="130" fixed="right" align="center">
          <template #default="scope">
            <el-dropdown trigger="click">
              <el-button link size="small">
                <el-icon><More /></el-icon> 更多
              </el-button>
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item @click="handleView(scope.row)">
                    <el-icon><Setting /></el-icon>
                    <span>查看</span>
                  </el-dropdown-item>
                  <template v-if="hasPermission('inventory_transfer:update')">
                    <el-dropdown-item @click="handleEdit(scope.row)" v-if="scope.row.status === 1">
                      <el-icon><Edit /></el-icon>
                      <span>编辑</span>
                    </el-dropdown-item>
                  </template>
                  <template v-if="hasPermission('inventory_transfer:audit')">
                    <el-dropdown-item @click="handleAudit(scope.row, 2)" v-if="scope.row.status === 1">
                      <el-icon><Menu /></el-icon>
                      <span>审核通过</span>
                    </el-dropdown-item>
                    <el-dropdown-item @click="handleAudit(scope.row, 4)" v-if="scope.row.status === 1">
                      <el-icon><Delete /></el-icon>
                      <span>拒绝</span>
                    </el-dropdown-item>
                  </template>
                  <template v-if="hasPermission('inventory_transfer:create')">
                    <el-dropdown-item @click="handleExecute(scope.row)" v-if="scope.row.status === 2">
                      <el-icon><Setting /></el-icon>
                      <span>执行调拨</span>
                    </el-dropdown-item>
                  </template>
                  <template v-if="hasPermission('inventory_transfer:delete')">
                    <el-dropdown-item divided @click="handleDelete(scope.row.id)" v-if="scope.row.status === 1 || scope.row.status === 4" danger>
                      <el-icon><Delete /></el-icon>
                      <span>删除</span>
                    </el-dropdown-item>
                  </template>
                </el-dropdown-menu>
              </template>
            </el-dropdown>
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

    <el-dialog
      :title="dialogTitle"
      v-model="dialogVisible"
      width="700px"
      :fullscreen="isMobile"
    >
      <template v-if="isView">
        <div class="detail-card">
          <div class="detail-header">
            <div class="transfer-no">调拨单号：{{ form.transferNo }}</div>
            <el-tag :type="getStatusTagType(form.status)">{{ getStatusText(form.status) }}</el-tag>
          </div>
          <div class="detail-body">
            <div class="detail-row">
              <span class="detail-label">基本信息</span>
            </div>
            <div class="detail-grid">
              <div class="detail-item">
                <span class="label">源仓库：</span>
                <span class="value">{{ getWarehouseName(form.fromWarehouseId) }}</span>
              </div>
              <div class="detail-item">
                <span class="label">目标仓库：</span>
                <span class="value">{{ getWarehouseName(form.toWarehouseId) }}</span>
              </div>
              <div class="detail-item">
                <span class="label">产品：</span>
                <span class="value">{{ getProductName(form.productId) }}</span>
              </div>
              <div class="detail-item">
                <span class="label">调拨数量：</span>
                <span class="value">{{ form.quantity }}</span>
              </div>
            </div>
            <div class="detail-row">
              <span class="detail-label">操作信息</span>
            </div>
            <div class="detail-grid">
              <div class="detail-item">
                <span class="label">创建人：</span>
                <span class="value">{{ form.createdByName || '-' }}</span>
              </div>
              <div class="detail-item">
                <span class="label">创建时间：</span>
                <span class="value">{{ form.createdAt || '-' }}</span>
              </div>
              <div class="detail-item">
                <span class="label">审核人：</span>
                <span class="value">{{ form.auditedByName || '-' }}</span>
              </div>
              <div class="detail-item">
                <span class="label">审核时间：</span>
                <span class="value">{{ form.auditedAt || '-' }}</span>
              </div>
              <div class="detail-item">
                <span class="label">执行人：</span>
                <span class="value">{{ form.executedByName || '-' }}</span>
              </div>
              <div class="detail-item">
                <span class="label">执行时间：</span>
                <span class="value">{{ form.executedAt || '-' }}</span>
              </div>
            </div>
            <div class="detail-row" v-if="form.remark">
              <span class="detail-label">备注</span>
              <div class="remark-content">{{ form.remark }}</div>
            </div>
          </div>
        </div>
      </template>
      <template v-else>
        <div class="dialog-description">
          <span class="dialog-description-icon">📦</span>
          <span class="dialog-description-text">库存调拨：将产品从一个仓库转移到另一个仓库。请选择源仓库、目标仓库、产品和调拨数量。调拨单需要审核通过后才能执行。</span>
        </div>
        <el-form :model="form" ref="formRef" :rules="rules" label-width="100px">
          <el-form-item label="源仓库" prop="fromWarehouseId">
            <el-select
              v-model="form.fromWarehouseId"
              placeholder="请选择源仓库"
              filterable
              style="width: 100%"
              @change="handleFromWarehouseChange"
            >
              <el-option v-for="wh in warehouses" :key="wh.id" :label="wh.name" :value="wh.id"></el-option>
            </el-select>
          </el-form-item>
          <el-form-item label="目标仓库" prop="toWarehouseId">
            <el-select
              v-model="form.toWarehouseId"
              placeholder="请选择目标仓库"
              filterable
              style="width: 100%"
            >
              <el-option v-for="wh in warehouses" :key="wh.id" :label="wh.name" :value="wh.id" :disabled="wh.id === form.fromWarehouseId"></el-option>
            </el-select>
          </el-form-item>
          <el-form-item label="产品" prop="productId">
            <el-select
              v-model="form.productId"
              placeholder="请选择产品"
              filterable
              style="width: 100%"
              @change="handleProductChange"
            >
              <el-option v-for="p in products" :key="p.id" :label="p.name" :value="p.id"></el-option>
            </el-select>
          </el-form-item>
          <el-form-item label="调拨数量" prop="quantity">
            <el-input-number
              v-model="form.quantity"
              :min="1"
              style="width: 100%"
              placeholder="请输入调拨数量"
            ></el-input-number>
            <div v-if="currentStock !== null" class="stock-info">
              当前库存：{{ currentStock }}
            </div>
          </el-form-item>
          <el-form-item label="备注">
            <el-input v-model="form.remark" type="textarea" :rows="3" placeholder="请输入备注信息"></el-input>
          </el-form-item>
        </el-form>
      </template>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="dialogVisible = false">关闭</el-button>
          <el-button type="primary" @click="handleSubmit" v-if="isCreate || isEdit">确认</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { More, Edit, Delete, Setting, Menu } from '@element-plus/icons-vue'
import { inventoryTransferApi } from '../../api/inventoryTransfer'
import { inventoryApi } from '../../api/inventory'
import { productApi } from '../../api/product'
import { warehouseApi } from '../../api/warehouse'
import { permission } from '../../utils/permission'

const hasPermission = (code) => permission.has(code)

const transferOrders = ref([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(10)
const searchQuery = ref('')
const fromWarehouseFilter = ref('')
const toWarehouseFilter = ref('')
const productFilter = ref('')
const statusFilter = ref('')

const dialogVisible = ref(false)
const dialogTitle = ref('新建调拨单')
const isCreate = ref(true)
const isEdit = ref(false)
const isView = ref(false)
const form = ref({
  id: null,
  fromWarehouseId: '',
  toWarehouseId: '',
  productId: '',
  quantity: 1,
  remark: ''
})
const formRef = ref(null)
const currentStock = ref(null)

const rules = {
  fromWarehouseId: [{ required: true, message: '请选择源仓库', trigger: 'change' }],
  toWarehouseId: [{ required: true, message: '请选择目标仓库', trigger: 'change' }],
  productId: [{ required: true, message: '请选择产品', trigger: 'change' }],
  quantity: [{ required: true, message: '请输入调拨数量', trigger: 'blur' }]
}

const warehouses = ref([])
const products = ref([])

const isMobile = computed(() => {
  if (typeof window !== 'undefined') {
    return window.innerWidth < 768
  }
  return false
})

const getStatusTagType = (status) => {
  const types = {
    1: 'info',
    2: 'warning',
    3: 'success',
    4: 'danger'
  }
  return types[status] || 'info'
}

const getStatusText = (status) => {
  const texts = {
    1: '待审核',
    2: '已审核',
    3: '已完成',
    4: '已拒绝'
  }
  return texts[status] || '未知'
}

const getWarehouseName = (warehouseId) => {
  const warehouse = warehouses.value.find(w => w.id === warehouseId)
  return warehouse ? warehouse.name : '-'
}

const getProductName = (productId) => {
  const product = products.value.find(p => p.id === productId)
  return product ? product.name : '-'
}

const loadTransferOrders = () => {
  const params = {
    page: currentPage.value,
    page_size: pageSize.value
  }
  if (searchQuery.value) {
    params.transfer_no = searchQuery.value
  }
  if (fromWarehouseFilter.value) {
    params.from_warehouse_id = fromWarehouseFilter.value
  }
  if (toWarehouseFilter.value) {
    params.to_warehouse_id = toWarehouseFilter.value
  }
  if (productFilter.value) {
    params.product_id = productFilter.value
  }
  if (statusFilter.value) {
    params.status = statusFilter.value
  }

  inventoryTransferApi.getInventoryTransferList(params)
  .then(response => {
    if (response.code === 0) {
      transferOrders.value = response.data.transfers || []
      total.value = response.data.total || 0
    } else {
      ElMessage.error('获取调拨单失败：' + response.message)
    }
  })
  .catch(error => {
    console.error('获取调拨单失败:', error)
    ElMessage.error('获取调拨单失败')
  })
}

const loadWarehouses = () => {
  warehouseApi.getActiveWarehouseList({})
  .then(response => {
    if (response.code === 0) {
      warehouses.value = response.data.warehouses || []
    }
  })
}

const loadProducts = () => {
  productApi.getActiveProductList({ page_size: 1000 })
  .then(response => {
    if (response.code === 0) {
      products.value = response.data.products || []
    }
  })
}

const handleCreate = () => {
  dialogTitle.value = '新建调拨单'
  isCreate.value = true
  isEdit.value = false
  isView.value = false
  currentStock.value = null
  form.value = {
    id: null,
    fromWarehouseId: '',
    toWarehouseId: '',
    productId: '',
    quantity: 1,
    remark: ''
  }
  dialogVisible.value = true
}

const handleEdit = (row) => {
  dialogTitle.value = '编辑调拨单'
  isCreate.value = false
  isEdit.value = true
  isView.value = false
  currentStock.value = null
  form.value = {
    id: row.id,
    fromWarehouseId: row.from_warehouse_id,
    toWarehouseId: row.to_warehouse_id,
    productId: row.product_id,
    quantity: row.quantity,
    remark: row.remark
  }
  dialogVisible.value = true
}

const handleView = (row) => {
  dialogTitle.value = '查看调拨单'
  isCreate.value = false
  isEdit.value = false
  isView.value = true
  form.value = {
    id: row.id,
    transferNo: row.transfer_no,
    fromWarehouseId: row.from_warehouse_id,
    toWarehouseId: row.to_warehouse_id,
    productId: row.product_id,
    quantity: row.quantity,
    status: row.status,
    remark: row.remark,
    createdAt: row.created_at,
    createdByName: row.created_by_name,
    auditedAt: row.audited_at,
    auditedByName: row.audited_by_name,
    executedByName: row.executed_by_name,
    executedAt: row.executed_at
  }
  dialogVisible.value = true
}

const handleFromWarehouseChange = () => {
  currentStock.value = null
  if (form.value.fromWarehouseId && form.value.productId) {
    loadCurrentStock()
  }
}

const handleProductChange = () => {
  currentStock.value = null
  if (form.value.fromWarehouseId && form.value.productId) {
    loadCurrentStock()
  }
}

const loadCurrentStock = () => {
  if (!form.value.fromWarehouseId || !form.value.productId) {
    currentStock.value = null
    return
  }
  inventoryApi.getCurrentStock({
    product_id: form.value.productId,
    warehouse_id: form.value.fromWarehouseId
  }).then(response => {
    if (response.code === 0) {
      currentStock.value = response.data.quantity || 0
    } else {
      currentStock.value = 0
    }
  }).catch(() => {
    currentStock.value = 0
  })
}

const handleSubmit = () => {
  formRef.value.validate((valid) => {
    if (valid) {
      if (isCreate.value) {
        if (currentStock.value !== null && form.value.quantity > currentStock.value) {
          ElMessage.error(`源仓库库存不足，当前库存：${currentStock.value}，申请调拨：${form.value.quantity}`)
          return
        }
        const data = {
          from_warehouse_id: form.value.fromWarehouseId,
          to_warehouse_id: form.value.toWarehouseId,
          product_id: form.value.productId,
          quantity: form.value.quantity,
          remark: form.value.remark
        }
        inventoryTransferApi.createInventoryTransfer(data)
        .then(response => {
          if (response.code === 0) {
            ElMessage.success('创建成功')
            dialogVisible.value = false
            loadTransferOrders()
          } else {
            ElMessage.error('创建失败：' + response.message)
          }
        })
      } else if (isEdit.value) {
        if (currentStock.value !== null && form.value.quantity > currentStock.value) {
          ElMessage.error(`源仓库库存不足，当前库存：${currentStock.value}，申请调拨：${form.value.quantity}`)
          return
        }
        const data = {
          id: form.value.id,
          from_warehouse_id: form.value.fromWarehouseId,
          to_warehouse_id: form.value.toWarehouseId,
          product_id: form.value.productId,
          quantity: form.value.quantity,
          remark: form.value.remark
        }
        inventoryTransferApi.updateInventoryTransfer(data)
        .then(response => {
          if (response.code === 0) {
            ElMessage.success('更新成功')
            dialogVisible.value = false
            loadTransferOrders()
          } else {
            ElMessage.error('更新失败：' + response.message)
          }
        })
      }
    }
  })
}

const handleAudit = (row, status) => {
  const message = status === 2 ? '确定要审核通过此调拨单吗？' : '确定要拒绝此调拨单吗？'
  ElMessageBox.confirm(message, '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消'
  }).then(() => {
    const data = {
      id: parseInt(row.id),
      status: status,
      remark: ''
    }
    console.log('审核请求数据:', JSON.stringify(data))
    inventoryTransferApi.auditInventoryTransfer(data)
    .then(response => {
      if (response.code === 0) {
        ElMessage.success('操作成功')
        loadTransferOrders()
      } else {
        ElMessage.error('操作失败：' + response.message)
      }
    })
    .catch(error => {
      console.error('审核请求失败:', error)
      ElMessage.error('审核请求失败')
    })
  })
}

const handleExecute = (row) => {
  ElMessageBox.confirm('确定要执行此调拨单吗？执行后将自动调整库存。', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消'
  }).then(() => {
    inventoryTransferApi.executeInventoryTransfer({ id: row.id })
    .then(response => {
      if (response.code === 0) {
        ElMessage.success('执行成功')
        loadTransferOrders()
      } else {
        ElMessage.error('执行失败：' + response.message)
      }
    })
  })
}

const handleDelete = (id) => {
  ElMessageBox.confirm('确定要删除此调拨单吗？', '警告', {
    type: 'warning'
  }).then(() => {
    inventoryTransferApi.deleteInventoryTransfer({ id })
    .then(response => {
      if (response.code === 0) {
        ElMessage.success('删除成功')
        loadTransferOrders()
      } else {
        ElMessage.error('删除失败：' + response.message)
      }
    })
  })
}

const handleSearch = () => {
  currentPage.value = 1
  loadTransferOrders()
}

const handleReset = () => {
  searchQuery.value = ''
  fromWarehouseFilter.value = ''
  toWarehouseFilter.value = ''
  productFilter.value = ''
  statusFilter.value = ''
  currentPage.value = 1
  loadTransferOrders()
}

const handleSizeChange = (size) => {
  pageSize.value = size
  loadTransferOrders()
}

const handleCurrentChange = (current) => {
  currentPage.value = current
  loadTransferOrders()
}

onMounted(() => {
  loadTransferOrders()
  loadWarehouses()
  loadProducts()
})
</script>

<style scoped>
.inventory-transfer {
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

.dialog-description {
  display: flex;
  align-items: flex-start;
  padding: 12px 16px;
  background-color: #fff7f0;
  border-left: 4px solid #ff6b35;
  margin-bottom: 20px;
  border-radius: 4px;
}

.dialog-description-icon {
  font-size: 16px;
  margin-right: 10px;
  flex-shrink: 0;
}

.dialog-description-text {
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

.pagination-container {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
  width: 100%;
  padding: 0 10px;
}

.stock-info {
  margin-top: 8px;
  color: #909399;
  font-size: 12px;
}

.detail-card {
  background: #fafafa;
  border-radius: 8px;
  padding: 20px;
}

.detail-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
  padding-bottom: 15px;
  border-bottom: 1px solid #e8e8e8;
}

.transfer-no {
  font-size: 18px;
  font-weight: bold;
  color: #303133;
}

.detail-body {
  padding: 10px 0;
}

.detail-row {
  margin-bottom: 15px;
}

.detail-label {
  display: inline-block;
  padding: 4px 12px;
  background: #e8f5e9;
  color: #2e7d32;
  border-radius: 4px;
  font-size: 14px;
  font-weight: 500;
  margin-bottom: 12px;
}

.detail-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 15px;
  margin-bottom: 20px;
}

.detail-item {
  display: flex;
  align-items: center;
  padding: 10px 15px;
  background: #fff;
  border-radius: 6px;
  box-shadow: 0 1px 2px rgba(0, 0, 0, 0.05);
}

.detail-item .label {
  color: #909399;
  font-size: 14px;
  min-width: 80px;
}

.detail-item .value {
  color: #303133;
  font-size: 14px;
  font-weight: 500;
}

.remark-content {
  padding: 12px 15px;
  background: #fff;
  border-radius: 6px;
  box-shadow: 0 1px 2px rgba(0, 0, 0, 0.05);
  color: #303133;
  font-size: 14px;
  line-height: 1.6;
}

:deep(.el-form-item) {
  margin-bottom: 20px;
}

:deep(.el-form-item__label) {
  text-align: right;
}

:deep(.el-form-item__content) {
  flex: 1;
  min-width: 0;
}

:deep(.el-input),
:deep(.el-select) {
  width: 100%;
}

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

  .search-container .el-input,
  .search-container .el-select {
    width: 100% !important;
  }

  .pagination-container {
    justify-content: center;
  }

  :deep(.el-form-item__label) {
    width: 80px;
  }
}
</style>
