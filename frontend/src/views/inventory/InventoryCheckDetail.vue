<template>
  <div class="inventory-check-detail">
    <el-card>
      <template #header>
        <div class="card-header">
          <div class="header-left">
            <el-button @click="handleBack" icon="el-icon-arrow-left">返回</el-button>
            <span class="header-title">{{ pageMode === 'view' ? '盘点单详情' : '盘点单录入' }}</span>
          </div>
          <div class="header-right" v-if="pageMode === 'edit' && (checkOrder.status === 1 || checkOrder.status === 2)">
            <el-button type="primary" @click="handleSave">保存</el-button>
            <el-button type="success" @click="handleSubmitCheck">提交盘点</el-button>
          </div>
        </div>
      </template>

      <el-form :model="checkOrder" label-width="100px" class="check-info">
        <el-row :gutter="20">
          <el-col :span="8">
            <el-form-item label="盘点单号">
              <el-tag type="primary">{{ checkOrder.check_no || '-' }}</el-tag>
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="仓库">
              <span>{{ checkOrder.warehouse?.name || '-' }}</span>
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="状态">
              <el-tag :type="getStatusTagType(checkOrder.status)">
                {{ getStatusText(checkOrder.status) }}
              </el-tag>
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="20">
          <el-col :span="8">
            <el-form-item label="产品总数">
              <span>{{ checkOrder.items?.length || 0 }}</span>
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="总差异">
              <span :class="getDiffClass(checkOrder.total_diff)">
                {{ (checkOrder.total_diff || 0) > 0 ? '+' : '' }}{{ checkOrder.total_diff || 0 }}</span>
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="创建时间">
              <span>{{ checkOrder.created_at || '-' }}</span>
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="20" v-if="pageMode === 'edit'">
          <el-col :span="24">
            <el-form-item label="备注">
              <el-input v-model="checkOrder.remark" placeholder="请输入备注信息" :disabled="checkOrder.status !== 1 && checkOrder.status !== 2"></el-input>
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="20" v-else>
          <el-col :span="24">
            <el-form-item label="备注">
              <span>{{ checkOrder.remark || '-' }}</span>
            </el-form-item>
          </el-col>
        </el-row>
      </el-form>

      <div class="detail-section">
        <div class="section-header">
          <span class="section-title">盘点明细</span>
          <div class="section-actions" v-if="pageMode === 'edit' && (checkOrder.status === 1 || checkOrder.status === 2)">
            <el-input
              v-model="batchActualStock"
              placeholder="批量录入实际库存"
              style="width: 150px; margin-right: 10px"
              @keyup.enter="handleBatchFill"
            >
              <template #append>
                <el-button @click="handleBatchFill">填充</el-button>
              </template>
            </el-input>
            <el-button type="primary" size="small" @click="handleBatchSame">与系统库存相同</el-button>
          </div>
        </div>
        <el-table
          :data="checkDetails"
          style="width: 100%"
          stripe
          border
          :fit="true"
        >
          <el-table-column prop="product" label="产品名称" min-width="150">
            <template #default="scope">
              {{ scope.row.product?.name || '-' }}
            </template>
          </el-table-column>
          <el-table-column prop="product" label="产品编码" width="120">
            <template #default="scope">
              {{ scope.row.product?.code || '-' }}
            </template>
          </el-table-column>
          <el-table-column prop="product" label="规格" width="100">
            <template #default="scope">
              {{ scope.row.product?.spec || '-' }}
            </template>
          </el-table-column>
          <el-table-column prop="product" label="单位" width="80">
            <template #default="scope">
              {{ scope.row.product?.main_unit || scope.row.product?.MainUnit || '-' }}
            </template>
          </el-table-column>
          <el-table-column prop="system_qty" label="系统库存" width="100" align="center"></el-table-column>
          <el-table-column prop="actual_qty" label="实际库存" width="150" align="center">
            <template #default="scope">
              <el-input-number
                v-if="pageMode === 'edit' && (checkOrder.status === 1 || checkOrder.status === 2)"
                v-model="scope.row.actual_qty"
                :min="0"
                controls-position="right"
                size="small"
                style="width: 100%"
                @change="handleStockChange(scope.$index)"
              ></el-input-number>
              <span v-else>{{ scope.row.actual_qty }}</span>
            </template>
          </el-table-column>
          <el-table-column prop="diff_qty" label="差异" width="100" align="center">
            <template #default="scope">
              <span :class="getDiffClass(scope.row.diff_qty)">
                {{ scope.row.diff_qty > 0 ? '+' : '' }}{{ scope.row.diff_qty || 0 }}
              </span>
            </template>
          </el-table-column>
        </el-table>
      </div>

      <div class="summary-section" v-if="checkDetails.length > 0">
        <div class="summary-item">
          <span class="summary-label">盘点产品数：</span>
          <span class="summary-value">{{ checkDetails.length }}</span>
        </div>
        <div class="summary-item">
          <span class="summary-label">盘盈数：</span>
          <span class="summary-value diff-positive">{{ getProfitCount() }}</span>
        </div>
        <div class="summary-item">
          <span class="summary-label">盘亏数：</span>
          <span class="summary-value diff-negative">{{ getLossCount() }}</span>
        </div>
        <div class="summary-item">
          <span class="summary-label">无差异数：</span>
          <span class="summary-value diff-zero">{{ getNoDiffCount() }}</span>
        </div>
      </div>

      <div class="action-section" v-if="pageMode === 'view' && checkOrder.status === 4 && hasDiff">
        <el-button type="warning" @click="handleGenerateAdjust" v-has-permission="'inventory_check:generate'">生成调整申请</el-button>
      </div>
    </el-card>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { inventoryCheckApi } from '../../api/inventoryCheck'

const router = useRouter()
const route = useRoute()

const pageMode = ref('view')
const checkOrder = ref({
  id: 0,
  checkNo: '',
  warehouseId: '',
  warehouse: null,
  status: 1,
  totalDiff: 0,
  remark: '',
  createdAt: '',
  items: []
})
const checkDetails = ref([])
const batchActualStock = ref('')

const isMobile = computed(() => {
  if (typeof window !== 'undefined') {
    return window.innerWidth < 768
  }
  return false
})

const hasDiff = computed(() => {
  return checkDetails.value.some(item => item.diff_qty !== 0)
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

const getProfitCount = () => {
  return checkDetails.value.filter(item => item.diff_qty > 0).length
}

const getLossCount = () => {
  return checkDetails.value.filter(item => item.diff_qty < 0).length
}

const getNoDiffCount = () => {
  return checkDetails.value.filter(item => item.diff_qty === 0).length
}

const loadCheckOrder = () => {
  const id = route.params.id
  if (!id) {
    ElMessage.error('无效的盘点单ID')
    router.back()
    return
  }

  pageMode.value = route.query.mode || 'view'

  inventoryCheckApi.getInventoryCheck(id)
  .then(response => {
    if (response.code === 0) {
      console.log('盘点单详情数据:', response.data)
      const checkData = response.data.check || response.data
      if (checkData.items && checkData.items.length > 0) {
        console.log('第一个商品数据:', checkData.items[0])
        console.log('第一个商品的product:', checkData.items[0].product)
      }
      checkOrder.value = checkData
      checkDetails.value = checkData.items || []
      checkDetails.value.forEach(item => {
        if (item.actual_qty === undefined || item.actual_qty === null) {
          item.actual_qty = item.system_qty
        }
        if (item.diff_qty === undefined) {
          item.diff_qty = (item.actual_qty || 0) - (item.system_qty || 0)
        }
      })
    } else {
      ElMessage.error('获取盘点单详情失败：' + response.message)
    }
  })
  .catch(error => {
    console.error('获取盘点单详情失败:', error)
    ElMessage.error('获取盘点单详情失败')
  })
}

const handleStockChange = (index) => {
  const item = checkDetails.value[index]
  item.diff_qty = (item.actual_qty || 0) - (item.system_qty || 0)
}

const handleBatchFill = () => {
  const value = parseInt(batchActualStock.value)
  if (isNaN(value) || value < 0) {
    ElMessage.warning('请输入有效的数字')
    return
  }
  checkDetails.value.forEach(item => {
    item.actualQty = value
    item.diffQty = (item.actualQty || 0) - (item.systemQty || 0)
  })
  ElMessage.success('已批量填充实际库存')
  batchActualStock.value = ''
}

const handleBatchSame = () => {
  checkDetails.value.forEach(item => {
    item.actual_qty = item.system_qty
    item.diff_qty = 0
  })
  ElMessage.success('已设置所有实际库存与系统库存相同')
}

const handleSave = () => {
  const items = checkDetails.value.map(item => ({
    id: item.id,
    actual_qty: item.actual_qty
  }))

  inventoryCheckApi.updateInventoryCheck({
    id: checkOrder.value.id,
    remark: checkOrder.value.remark,
    items: items
  })
  .then(response => {
    if (response.code === 0) {
      ElMessage.success('保存成功')
      loadCheckOrder()
    } else {
      ElMessage.error('保存失败：' + response.message)
    }
  })
}

const handleSubmitCheck = () => {
  const hasEmptyStock = checkDetails.value.some(item => item.actual_qty === undefined || item.actual_qty === null)
  if (hasEmptyStock) {
    ElMessage.warning('请填写所有产品的实际库存')
    return
  }

  ElMessageBox.confirm('确定要提交盘点单吗？提交后将无法继续编辑。', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消'
  }).then(() => {
      const items = checkDetails.value.map(item => ({
        id: item.id,
        actual_qty: item.actual_qty
      }))

      inventoryCheckApi.updateInventoryCheck({
        id: checkOrder.value.id,
        remark: checkOrder.value.remark,
        items: items
      })
    .then(response => {
      if (response.code === 0) {
        return inventoryCheckApi.submitInventoryCheck({ 
          id: checkOrder.value.id,
          remark: checkOrder.value.remark || ''
        })
      } else {
        throw new Error(response.message || '保存失败')
      }
    })
    .then(response => {
      if (response.code === 0) {
        ElMessage.success('提交成功')
        loadCheckOrder()
      } else {
        ElMessage.error('提交失败：' + response.message)
      }
    })
    .catch(error => {
      ElMessage.error(error.message || '操作失败')
    })
  })
}

const handleGenerateAdjust = () => {
  ElMessageBox.confirm('确定要生成调整申请吗？将根据盘点差异生成库存调整单。', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消'
  }).then(() => {
    inventoryCheckApi.generateInventoryAdjust({ id: checkOrder.value.id })
    .then(response => {
      if (response.code === 0) {
        ElMessage.success('生成调整申请成功')
        router.push('/inventory/adjust-request')
      } else {
        ElMessage.error('生成调整申请失败：' + response.message)
      }
    })
  })
}

const handleBack = () => {
  router.push('/inventory/check')
}

onMounted(() => {
  loadCheckOrder()
})
</script>

<style scoped>
.inventory-check-detail {
  padding: 0;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 15px;
}

.header-title {
  font-size: 18px;
  font-weight: 500;
}

.header-right {
  display: flex;
  gap: 10px;
}

.check-info {
  margin-bottom: 20px;
  padding-bottom: 20px;
  border-bottom: 1px solid #ebeef5;
}

.detail-section {
  margin-top: 20px;
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 15px;
}

.section-title {
  font-size: 16px;
  font-weight: 500;
}

.section-actions {
  display: flex;
  align-items: center;
}

.summary-section {
  display: flex;
  gap: 30px;
  margin-top: 20px;
  padding: 15px;
  background: #f5f7fa;
  border-radius: 4px;
}

.summary-item {
  display: flex;
  align-items: center;
}

.summary-label {
  color: #606266;
  margin-right: 5px;
}

.summary-value {
  font-weight: 600;
  font-size: 16px;
}

.action-section {
  margin-top: 20px;
  padding-top: 20px;
  border-top: 1px solid #ebeef5;
  text-align: center;
}

.diff-positive {
  color: #67c23a;
  font-weight: 500;
}

.diff-negative {
  color: #f56c6c;
  font-weight: 500;
}

.diff-zero {
  color: #909399;
}

:deep(.el-form-item) {
  margin-bottom: 0;
}

:deep(.el-form-item__label) {
  text-align: right;
}

:deep(.el-input-number) {
  width: 100%;
}

@media screen and (max-width: 768px) {
  .card-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 15px;
  }

  .header-right {
    width: 100%;
    justify-content: flex-end;
  }

  .section-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 10px;
  }

  .section-actions {
    width: 100%;
    flex-wrap: wrap;
  }

  .summary-section {
    flex-wrap: wrap;
    gap: 15px;
  }
}
</style>