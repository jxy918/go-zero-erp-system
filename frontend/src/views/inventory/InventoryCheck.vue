<template>
  <div class="inventory-check">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>库存盘点</span>
          <el-button type="primary" @click="handleCreate" v-has-permission="'inventory_check:create'">新建盘点单</el-button>
        </div>
      </template>
      <div class="description-box">
        <span class="description-icon">ℹ️</span>
        <span class="description-text">库存盘点：定期对仓库库存进行实际盘点，对比系统库存与实际库存差异。流程：新建盘点单 → 录入实际库存 → 提交审核 → 系统自动调整库存。</span>
      </div>
      <div class="query-bar">
        <el-input
          v-model="searchQuery"
          placeholder="请输入盘点单号"
          style="width: 200px"
          prefix-icon="el-icon-search"
          clearable
        ></el-input>
        <el-select v-model="warehouseFilter" placeholder="请选择仓库" clearable style="width: 150px">
          <el-option v-for="wh in warehouses" :key="wh.id" :label="wh.name" :value="wh.id"></el-option>
        </el-select>
        <el-select v-model="statusFilter" placeholder="请选择状态" clearable style="width: 150px">
          <el-option :value="1" label="待盘点"></el-option>
          <el-option :value="2" label="盘点中"></el-option>
          <el-option :value="3" label="已完成"></el-option>
          <el-option :value="4" label="已提交"></el-option>
        </el-select>
        <el-button type="primary" @click="handleSearch">搜索</el-button>
        <el-button @click="handleReset">重置</el-button>
      </div>
      <el-table :data="checkOrders" border stripe v-loading="loading">
        <el-table-column prop="id" label="ID" width="60"></el-table-column>
        <el-table-column prop="check_no" label="盘点单号" width="180"></el-table-column>
        <el-table-column prop="warehouse" label="仓库" width="100">
          <template #default="scope">
            {{ scope.row.warehouse?.name || '-' }}
          </template>
        </el-table-column>
        <el-table-column label="产品总数" width="100">
          <template #default="scope">
            {{ scope.row.item_count || scope.row.items?.length || 0 }}
          </template>
        </el-table-column>
        <el-table-column prop="total_diff" label="总差异" width="100">
          <template #default="scope">
            <span :class="getDiffClass(scope.row.total_diff)">
              {{ scope.row.total_diff > 0 ? '+' : '' }}{{ scope.row.total_diff || 0 }}
            </span>
          </template>
        </el-table-column>
        <el-table-column prop="remark" label="备注" min-width="150" show-overflow-tooltip></el-table-column>
        <el-table-column prop="status" label="状态" width="80">
          <template #default="scope">
            <el-tag :type="getStatusTagType(scope.row.status)">{{ getStatusText(scope.row.status) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="创建时间" width="160"></el-table-column>
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
                    <span>查看详情</span>
                  </el-dropdown-item>
                  <el-dropdown-item @click="handleRecord(scope.row)" v-if="scope.row.status === 1 || scope.row.status === 2">
                    <el-icon><Edit /></el-icon>
                    <span>录入</span>
                  </el-dropdown-item>
                  <el-dropdown-item @click="handleSubmit(scope.row)" v-if="scope.row.status === 1 || scope.row.status === 2" v-has-permission="'inventory_check:submit'">
                    <el-icon><Menu /></el-icon>
                    <span>提交</span>
                  </el-dropdown-item>
                  <el-dropdown-item divided @click="handleDelete(scope.row.id)" v-if="scope.row.status === 1" v-has-permission="'inventory_check:delete'" danger>
                    <el-icon><Delete /></el-icon>
                    <span>删除</span>
                  </el-dropdown-item>
                </el-dropdown-menu>
              </template>
            </el-dropdown>
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
      title="新建盘点单"
      v-model="createVisible"
      width="600px"
      :fullscreen="isMobile"
    >
      <div class="dialog-description">
        <span class="dialog-description-icon">📋</span>
        <span class="dialog-description-text">库存盘点：选择仓库后，系统会自动获取该仓库下的所有产品。盘点单创建后，需要录入实际库存数量并提交审核。</span>
      </div>
      <el-form :model="createForm" ref="createFormRef" :rules="createRules" label-width="100px">
        <el-form-item label="仓库" prop="warehouseId">
          <el-select
            v-model="createForm.warehouseId"
            placeholder="请选择仓库"
            filterable
            style="width: 100%"
          >
            <el-option v-for="wh in warehouses" :key="wh.id" :label="wh.name" :value="wh.id"></el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="备注">
          <el-input v-model="createForm.remark" type="textarea" :rows="3" placeholder="请输入备注信息"></el-input>
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="createVisible = false">取消</el-button>
          <el-button type="primary" @click="handleCreateSubmit">确认创建</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { More, Setting, Edit, Delete, Menu } from '@element-plus/icons-vue'
import { inventoryCheckApi } from '../../api/inventoryCheck'
import { inventoryApi } from '../../api/inventory'
import { warehouseApi } from '../../api/warehouse'

const router = useRouter()

const loading = ref(false)
const checkOrders = ref([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(10)
const searchQuery = ref('')
const warehouseFilter = ref('')
const statusFilter = ref('')

const createVisible = ref(false)
const createForm = ref({
  warehouseId: '',
  remark: ''
})
const createFormRef = ref(null)

const createRules = {
  warehouseId: [{ required: true, message: '请选择仓库', trigger: 'change' }]
}

const warehouses = ref([])

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
    4: 'primary'
  }
  return types[status] || 'info'
}

const getStatusText = (status) => {
  const texts = {
    1: '待盘点',
    2: '盘点中',
    3: '已完成',
    4: '已提交'
  }
  return texts[status] || '未知'
}

const getDiffClass = (diff) => {
  if (diff > 0) return 'diff-positive'
  if (diff < 0) return 'diff-negative'
  return 'diff-zero'
}

const loadCheckOrders = async () => {
  loading.value = true
  try {
    const params = {
      page: currentPage.value,
      page_size: pageSize.value
    }
    if (searchQuery.value) {
      params.check_no = searchQuery.value
    }
    if (warehouseFilter.value) {
      params.warehouse_id = warehouseFilter.value
    }
    if (statusFilter.value) {
      params.status = statusFilter.value
    }

    const response = await inventoryCheckApi.getInventoryCheckList(params)
    if (response.code === 0) {
      checkOrders.value = response.data.checks || []
      total.value = response.data.total || 0
    } else {
      ElMessage.error('获取盘点单失败：' + response.message)
    }
  } catch (error) {
    console.error('获取盘点单失败:', error)
    ElMessage.error('获取盘点单失败')
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

const handleCreate = () => {
  createForm.value = {
    warehouseId: '',
    remark: ''
  }
  createVisible.value = true
}

const handleCreateSubmit = () => {
  createFormRef.value.validate((valid) => {
    if (valid) {
      inventoryCheckApi.createInventoryCheck({
        warehouse_id: createForm.value.warehouseId,
        items: [],
        remark: createForm.value.remark
      })
      .then(response => {
        if (response.code === 0) {
          ElMessage.success('创建成功')
          createVisible.value = false
          loadCheckOrders()
          if (response.data && response.data.id) {
            router.push(`/inventory/check/${response.data.id}?mode=edit`)
          }
        } else {
          ElMessage.error('创建失败：' + response.message)
        }
      })
    }
  })
}

const handleView = (row) => {
  router.push(`/inventory/check/${row.id}?mode=view`)
}

const handleRecord = (row) => {
  router.push(`/inventory/check/${row.id}?mode=edit`)
}

const handleSubmit = (row) => {
  ElMessageBox.confirm('确定要提交此盘点单吗？提交后将无法继续编辑。', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消'
  }).then(() => {
    inventoryCheckApi.submitInventoryCheck({ id: row.id })
    .then(response => {
      if (response.code === 0) {
        ElMessage.success('提交成功')
        loadCheckOrders()
      } else {
        ElMessage.error('提交失败：' + response.message)
      }
    })
  })
}

const handleDelete = (id) => {
  ElMessageBox.confirm('确定要删除此盘点单吗？', '警告', {
    type: 'warning'
  }).then(() => {
    inventoryCheckApi.deleteInventoryCheck({ id })
    .then(response => {
      if (response.code === 0) {
        ElMessage.success('删除成功')
        loadCheckOrders()
      } else {
        ElMessage.error('删除失败：' + response.message)
      }
    })
  })
}

const handleSearch = () => {
  currentPage.value = 1
  loadCheckOrders()
}

const handleReset = () => {
  searchQuery.value = ''
  warehouseFilter.value = ''
  statusFilter.value = ''
  currentPage.value = 1
  loadCheckOrders()
}

const handleSizeChange = (size) => {
  pageSize.value = size
  loadCheckOrders()
}

const handleCurrentChange = (current) => {
  currentPage.value = current
  loadCheckOrders()
}

onMounted(() => {
  loadCheckOrders()
  loadWarehouses()
})
</script>

<style scoped>
.inventory-check {
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

.diff-positive {
  color: #67c23a;
}

.diff-negative {
  color: #f56c6c;
}

.diff-zero {
  color: #909399;
}
</style>