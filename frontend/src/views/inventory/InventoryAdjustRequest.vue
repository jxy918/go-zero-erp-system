<template>
  <div class="inventory-adjust-request">
    <el-card>
      <template #header>
        <div class="card-header">
          <div class="header-left">
            <el-button @click="goBack" v-if="showBackButton">返回</el-button>
            <span>库存调整申请</span>
          </div>
          <el-button type="primary" @click="showCreateDialog">新建申请</el-button>
        </div>
      </template>

      <div class="description-box">
        <span class="description-icon">ℹ️</span>
        <span class="description-text">库存调整申请：用于处理日常库存差异调整。类型说明：【盘盈】实际库存大于系统库存（输入正数）；【盘亏】实际库存小于系统库存（输入负数）；【其他】特殊情况的库存调整（可正可负）。</span>
      </div>

      <div class="query-bar">
        <el-select v-model="queryForm.status" placeholder="请选择状态" clearable style="width: 150px">
          <el-option label="全部" :value="0"></el-option>
          <el-option label="待审核" :value="1"></el-option>
          <el-option label="已审核" :value="2"></el-option>
          <el-option label="已拒绝" :value="3"></el-option>
        </el-select>
        <el-button type="primary" @click="handleQuery">查询</el-button>
      </div>

      <el-table :data="tableData" border stripe v-loading="loading">
        <el-table-column prop="id" label="ID" width="60"></el-table-column>
        <el-table-column prop="request_no" label="单据号" width="180"></el-table-column>
        <el-table-column prop="product.name" label="产品名称" min-width="120">
          <template #default="{ row }">
            {{ row.product?.name || '-' }}
          </template>
        </el-table-column>
        <el-table-column prop="warehouse.name" label="仓库" width="100">
          <template #default="{ row }">
            {{ row.warehouse?.name || '-' }}
          </template>
        </el-table-column>
        <el-table-column prop="before_qty" label="调整前库存" width="100"></el-table-column>
        <el-table-column prop="quantity" label="调整数量" width="100">
          <template #default="{ row }">
            <span :class="row.quantity > 0 ? 'text-success' : 'text-danger'">
              {{ row.quantity > 0 ? '+' : '' }}{{ row.quantity }}
            </span>
          </template>
        </el-table-column>
        <el-table-column prop="after_qty" label="调整后库存" width="100"></el-table-column>
        <el-table-column prop="type_desc" label="调整类型" width="80"></el-table-column>
        <el-table-column prop="reason" label="调整原因" min-width="150" show-overflow-tooltip></el-table-column>
        <el-table-column prop="status_desc" label="状态" width="80">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)">{{ row.status_desc }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="applicant.nickname" label="申请人" width="100">
          <template #default="{ row }">
            {{ row.applicant?.nickname || '-' }}
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="申请时间" width="160"></el-table-column>
        <el-table-column label="操作" width="150" fixed="right">
          <template #default="{ row }">
            <template v-if="row.status === 1">
              <el-button type="success" size="small" @click="handleApprove(row)">审核</el-button>
              <el-button type="danger" size="small" @click="handleReject(row)">拒绝</el-button>
            </template>
            <template v-else-if="row.status === 2">
              <el-tag type="success">已审核</el-tag>
            </template>
            <template v-else>
              <el-tag type="danger">已拒绝</el-tag>
            </template>
          </template>
        </el-table-column>
      </el-table>

      <el-pagination
        v-model:current-page="pagination.page"
        v-model:page-size="pagination.pageSize"
        :total="pagination.total"
        :page-sizes="[10, 20, 50, 100]"
        layout="total, sizes, prev, pager, next, jumper"
        @size-change="handleQuery"
        @current-change="handleQuery"
        style="margin-top: 20px"
      />
    </el-card>

    <el-dialog v-model="createDialogVisible" title="新建库存调整申请" width="500px">
      <div class="dialog-description">
        <span class="dialog-description-icon">🔄</span>
        <span class="dialog-description-text">库存调整申请用于处理日常库存差异。调整申请需要审核通过后，系统才会更新库存数量。</span>
      </div>
      <el-form :model="createForm" :rules="createRules" ref="createFormRef" label-width="100px">
        <el-form-item label="产品" prop="productId">
          <el-select
            v-model="createForm.productId"
            placeholder="请选择产品"
            filterable
            @change="handleProductChange"
          >
            <el-option
              v-for="product in productList"
              :key="product.id"
              :label="product.name"
              :value="product.id"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="仓库" prop="warehouseId">
          <el-select v-model="createForm.warehouseId" placeholder="请选择仓库" @change="handleWarehouseChange">
            <el-option
              v-for="warehouse in warehouseList"
              :key="warehouse.id"
              :label="warehouse.name"
              :value="warehouse.id"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="当前库存">
          <span>{{ currentStock }}</span>
        </el-form-item>
        <el-form-item label="调整类型" prop="type">
          <el-select v-model="createForm.type" placeholder="请选择调整类型" @change="handleTypeChange">
            <el-option label="盘盈" :value="1"></el-option>
            <el-option label="盘亏" :value="2"></el-option>
            <el-option label="其他" :value="4"></el-option>
          </el-select>
          <div class="type-desc">
            <p><strong>盘盈：</strong>实际库存大于系统库存，增加库存数量（输入正数）</p>
            <p><strong>盘亏：</strong>实际库存小于系统库存，减少库存数量（输入负数）</p>
            <p><strong>其他：</strong>特殊情况的库存调整，可增可减（正数增加，负数减少）</p>
          </div>
        </el-form-item>
        <el-form-item label="调整数量" prop="quantity">
          <el-input-number 
            v-model="createForm.quantity" 
            :min="quantityMin" 
            :max="quantityMax"
            :disabled="!createForm.type"
          ></el-input-number>
          <span v-if="createForm.type === 1" class="quantity-tip">盘盈只能输入正数</span>
          <span v-else-if="createForm.type === 2" class="quantity-tip">盘亏只能输入负数</span>
          <span v-else-if="createForm.type === 4" class="quantity-tip">其他调整可输入正数（增加）或负数（减少）</span>
        </el-form-item>
        <el-form-item label="预估调整后">
          <span>{{ previewAfterQty }}</span>
        </el-form-item>
        <el-form-item label="调整原因" prop="reason">
          <el-input v-model="createForm.reason" type="textarea" :rows="3" placeholder="请输入调整原因"></el-input>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="createDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleCreate" :loading="submitLoading">提交申请</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="auditDialogVisible" :title="auditAction === 'approve' ? '审核通过' : '审核拒绝'" width="500px">
      <el-form :model="auditForm" label-width="100px">
        <el-form-item label="单据号">
          <span>{{ currentRow?.request_no }}</span>
        </el-form-item>
        <el-form-item label="产品">
          <span>{{ currentRow?.product?.name }}</span>
        </el-form-item>
        <el-form-item label="调整数量">
          <span :class="currentRow?.quantity > 0 ? 'text-success' : 'text-danger'">
            {{ currentRow?.quantity > 0 ? '+' : '' }}{{ currentRow?.quantity }}
          </span>
        </el-form-item>
        <el-form-item :label="auditAction === 'approve' ? '审核备注' : '拒绝原因'" prop="note">
          <el-input v-model="auditForm.note" type="textarea" :rows="3" :placeholder="auditAction === 'approve' ? '请输入审核备注' : '请输入拒绝原因'"></el-input>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="auditDialogVisible = false">取消</el-button>
        <el-button :type="auditAction === 'approve' ? 'success' : 'danger'" @click="handleAuditConfirm" :loading="submitLoading">
          {{ auditAction === 'approve' ? '确认通过' : '确认拒绝' }}
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { inventoryApi } from '@/api/inventory'
import { productApi } from '@/api/product'
import { warehouseApi } from '@/api/warehouse'

const route = useRoute()
const router = useRouter()
const loading = ref(false)
const submitLoading = ref(false)
const tableData = ref([])
const productList = ref([])
const warehouseList = ref([])

const showBackButton = computed(() => {
  // 只有当路由参数中包含 back=true 时才显示返回按钮
  return route.query.back === 'true'
})

const goBack = () => {
  router.push('/inventory')
}

const queryForm = reactive({
  status: 0
})

const pagination = reactive({
  page: 1,
  pageSize: 10,
  total: 0
})

const createDialogVisible = ref(false)
const auditDialogVisible = ref(false)
const auditAction = ref('')
const currentRow = ref(null)
const createFormRef = ref(null)

const createForm = reactive({
  productId: null,
  warehouseId: null,
  type: 1,
  quantity: 0,
  reason: ''
})

const auditForm = reactive({
  note: ''
})

const currentStock = ref(0)

const handleTypeChange = () => {
  createForm.quantity = 0
}

const loadCurrentStock = async () => {
  if (!createForm.productId || !createForm.warehouseId) {
    currentStock.value = 0
    return
  }
  try {
    const res = await inventoryApi.getCurrentStock({
      product_id: createForm.productId,
      warehouse_id: createForm.warehouseId
    })
    if (res.code === 0) {
      currentStock.value = res.data.quantity || 0
    } else {
      currentStock.value = 0
    }
  } catch (error) {
    currentStock.value = 0
  }
}

const previewAfterQty = computed(() => {
  let quantity = createForm.quantity
  if (createForm.type === 1) { // 盘盈
    quantity = Math.abs(createForm.quantity)
  } else if (createForm.type === 2) { // 盘亏
    if (quantity > 0) {
      quantity = -quantity
    }
  }
  return currentStock.value + quantity
})

const quantityMin = computed(() => {
  if (createForm.type === 1) { // 盘盈
    return 1
  } else if (createForm.type === 2) { // 盘亏
    return -99999
  } else if (createForm.type === 4) { // 其他
    return -99999
  }
  return -99999
})

const quantityMax = computed(() => {
  if (createForm.type === 1) { // 盘盈
    return 99999
  } else if (createForm.type === 2) { // 盘亏
    return -1
  } else if (createForm.type === 4) { // 其他
    return 99999
  }
  return 99999
})

const validateQuantity = (rule, value, callback) => {
  if (value === 0) {
    callback(new Error('调整数量不能为0'))
  } else {
    callback()
  }
}

const createRules = {
  productId: [{ required: true, message: '请选择产品', trigger: 'change' }],
  warehouseId: [{ required: true, message: '请选择仓库', trigger: 'change' }],
  type: [{ required: true, message: '请选择调整类型', trigger: 'change' }],
  quantity: [
    { required: true, message: '请输入调整数量', trigger: 'blur' },
    { validator: validateQuantity, trigger: 'blur' }
  ],
  reason: [{ required: true, message: '请输入调整原因', trigger: 'blur' }]
}

const getStatusType = (status) => {
  const types = { 1: 'warning', 2: 'success', 3: 'danger' }
  return types[status] || 'info'
}

const handleQuery = async () => {
  loading.value = true
  try {
    const params = {
      page: pagination.page,
      page_size: pagination.pageSize,
      status: queryForm.status || undefined
    }
    const res = await inventoryApi.getInventoryAdjustRequestList(params)
    if (res.code === 0) {
      tableData.value = res.data.requests || []
      pagination.total = res.data.total || 0
    }
  } catch (error) {
    console.error('查询失败:', error)
  } finally {
    loading.value = false
  }
}

const showCreateDialog = async () => {
  createDialogVisible.value = true
  createForm.productId = null
  createForm.warehouseId = null
  createForm.type = 1
  createForm.quantity = 0
  createForm.reason = ''
  await loadProducts()
  await loadWarehouses()
}

const loadProducts = async () => {
  try {
    const res = await productApi.getActiveProductList({ page: 1, page_size: 1000 })
    if (res.code === 0) {
      productList.value = res.data.products || []
    }
  } catch (error) {
    console.error('加载产品失败:', error)
  }
}

const loadWarehouses = async () => {
  try {
    const res = await warehouseApi.getActiveWarehouseList({ page: 1, page_size: 100 })
    if (res.code === 0) {
      warehouseList.value = res.data.warehouses || []
    }
  } catch (error) {
    console.error('加载仓库失败:', error)
  }
}

const handleProductChange = () => {
  createForm.warehouseId = null
  currentStock.value = 0
}

const handleWarehouseChange = () => {
  loadCurrentStock()
}

const handleCreate = async () => {
  const valid = await createFormRef.value.validate().catch(() => false)
  if (!valid) return

  submitLoading.value = true
  try {
    const data = {
      product_id: createForm.productId,
      warehouse_id: createForm.warehouseId,
      'type': createForm.type,
      quantity: createForm.quantity,
      reason: createForm.reason
    }
    const res = await inventoryApi.createInventoryAdjustRequest(data)
    if (res.code === 0) {
      ElMessage.success('申请提交成功')
      createDialogVisible.value = false
      handleQuery()
    } else {
      ElMessage.error(res.message || '提交失败')
    }
  } catch (error) {
    ElMessage.error('提交失败')
  } finally {
    submitLoading.value = false
  }
}

const handleApprove = (row) => {
  currentRow.value = row
  auditAction.value = 'approve'
  auditForm.note = ''
  auditDialogVisible.value = true
}

const handleReject = (row) => {
  ElMessageBox.confirm('确定要拒绝此申请吗?', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async () => {
    currentRow.value = row
    auditAction.value = 'reject'
    auditForm.note = ''
    auditDialogVisible.value = true
  }).catch(() => {})
}

const handleAuditConfirm = async () => {
  if (!auditForm.note) {
    ElMessage.warning('请输入审核备注/拒绝原因')
    return
  }

  submitLoading.value = true
  try {
    const data = auditAction.value === 'approve'
      ? {
          id: currentRow.value.id,
          status: 2,
          remark: '',
          note: auditForm.note
        }
      : {
          id: currentRow.value.id,
          remark: '',
          note: auditForm.note
        }
    const res = auditAction.value === 'approve'
      ? await inventoryApi.approveInventoryAdjustRequest(data)
      : await inventoryApi.rejectInventoryAdjustRequest(data)

    if (res.code === 0) {
      ElMessage.success(auditAction.value === 'approve' ? '审核通过' : '已拒绝')
      auditDialogVisible.value = false
      handleQuery()
    } else {
      ElMessage.error(res.message || '操作失败')
    }
  } catch (error) {
    ElMessage.error('操作失败')
  } finally {
    submitLoading.value = false
  }
}

onMounted(async () => {
  await handleQuery()
  await loadProducts()
  await loadWarehouses()

  // 检查路由参数
  const { productId, warehouseId } = route.query
  if (productId && warehouseId) {
    // 自动打开新建申请弹窗并预填参数
    createDialogVisible.value = true
    createForm.productId = Number(productId)
    createForm.warehouseId = Number(warehouseId)
    createForm.type = 1
    createForm.quantity = 0
    createForm.reason = ''
    // 加载当前库存
    await loadCurrentStock()
  }
})
</script>

<style scoped>
.inventory-adjust-request {
  padding: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 10px;
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

.text-success {
  color: #67c23a;
}

.text-danger {
  color: #f56c6c;
}

.type-desc {
  margin-top: 10px;
  padding: 12px 16px;
  background-color: #fff7f0;
  border-left: 4px solid #ff6b35;
  border-radius: 4px;
  font-size: 13px;
  line-height: 1.6;
}

.type-desc p {
  margin: 4px 0;
  color: #d93026;
}

.type-desc strong {
  color: #d93026;
  font-weight: 600;
}

.quantity-tip {
  margin-left: 10px;
  font-size: 12px;
  color: #8c8c8c;
}
</style>
