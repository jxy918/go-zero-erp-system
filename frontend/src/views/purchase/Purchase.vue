<template>
  <div class="purchase-management">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>采购订单</span>
          <el-button type="primary" @click="handleAdd" v-has-permission="'btn_purchase_create'">新建订单</el-button>
        </div>
      </template>
      <div class="description-box">
        <span class="description-icon">ℹ️</span>
        <span class="description-text">采购管理：管理采购订单，包括新建订单、审核、入库等操作，采购订单审核通过后可执行入库。</span>
      </div>
      <div class="query-bar">
        <el-input
          v-model="searchQuery"
          placeholder="请输入订单号"
          style="width: 300px"
          prefix-icon="el-icon-search"
        ></el-input>
        <el-select v-model="statusFilter" placeholder="请选择状态" clearable style="width: 150px">
          <el-option :value="0" label="全部"></el-option>
          <el-option :value="1" label="待审核"></el-option>
          <el-option :value="2" label="已审核"></el-option>
          <el-option :value="3" label="已入库"></el-option>
          <el-option :value="4" label="已取消"></el-option>
        </el-select>
        <el-button type="primary" @click="handleSearch">搜索</el-button>
      </div>
      <el-table
        :data="orders"
        style="width: 100%"
        border
        stripe
        v-loading="loading"
      >
        <el-table-column prop="id" label="ID" width="80"></el-table-column>
        <el-table-column prop="order_no" label="订单号" min-width="150"></el-table-column>
        <el-table-column prop="supplier.name" label="供应商" min-width="120"></el-table-column>
        <el-table-column prop="warehouse.name" label="仓库" min-width="100"></el-table-column>
        <el-table-column prop="total_amount" label="订单金额" width="120">
          <template #default="scope">
            <span class="amount-text">¥{{ (scope.row.total_amount || 0).toFixed(2) }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="100">
          <template #default="scope">
            <el-tag :type="getStatusTagType(scope.row.status)">
              {{ getStatusText(scope.row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="remark" label="备注" min-width="120"></el-table-column>
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
                  <el-dropdown-item @click="handleApprove(scope.row)" v-if="scope.row.status === 1" v-has-permission="'btn_purchase_approve'">
                    <el-icon><Edit /></el-icon>
                    <span>审核</span>
                  </el-dropdown-item>
                  <el-dropdown-item @click="handleInbound(scope.row)" v-if="scope.row.status === 2" v-has-permission="'btn_purchase_inbound'">
                    <el-icon><Setting /></el-icon>
                    <span>入库</span>
                  </el-dropdown-item>
                  <el-dropdown-item @click="handleCancel(scope.row)" v-if="scope.row.status === 1 || scope.row.status === 2" v-has-permission="'btn_purchase_cancel'">
                    <el-icon><Delete /></el-icon>
                    <span>取消</span>
                  </el-dropdown-item>
                  <el-dropdown-item divided @click="handleDelete(scope.row.id)" v-if="scope.row.status === 1 || scope.row.status === 2" v-has-permission="'btn_purchase_delete'" danger>
                    <el-icon><Delete /></el-icon>
                    <span>删除</span>
                  </el-dropdown-item>
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
      width="900px"
      :fullscreen="isMobile"
    >
      <el-form :model="form" ref="formRef" :rules="formRules" label-width="100px">
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="供应商" prop="supplierId">
              <el-select
                v-model="form.supplierId"
                placeholder="请选择供应商"
                @change="handleSupplierChange"
              >
                <el-option v-for="sup in suppliers" :key="sup.id" :label="sup.name" :value="sup.id"></el-option>
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="入库仓库" prop="warehouseId">
              <el-select v-model="form.warehouseId" placeholder="请选择仓库">
                <el-option v-for="wh in filteredWarehouses" :key="wh.id" :label="wh.name" :value="wh.id"></el-option>
              </el-select>
              <div v-if="form.supplierId && filteredWarehouses.length === 0" class="empty-tip">
                该供应商暂无关联仓库
              </div>
            </el-form-item>
          </el-col>
        </el-row>
        <el-form-item label="订单明细" prop="items">
          <div class="order-items">
            <el-button type="primary" size="small" @click="addOrderItem">+ 添加商品</el-button>
            <el-table :data="form.items" style="width: 100%" border>
              <el-table-column prop="productId" label="商品" min-width="150">
                <template #default="scope">
                  <el-select
                    v-model="scope.row.productId"
                    placeholder="选择商品"
                    class="product-select"
                    @change="handleProductChange(scope.$index)"
                  >
                    <el-option v-for="prod in products" :key="prod.id" :label="prod.name" :value="prod.id"></el-option>
                  </el-select>
                </template>
              </el-table-column>
              <el-table-column prop="unitId" label="单位" width="120">
                <template #default="scope">
                  <el-select
                    v-model="scope.row.unitId"
                    placeholder="选择单位"
                    class="unit-select"
                    :disabled="!scope.row.productId"
                    @change="handleUnitChange(scope.$index)"
                  >
                    <el-option v-for="unit in getProductUnits(scope.row.productId)" :key="unit.id" :label="unit.unit_name" :value="unit.id"></el-option>
                  </el-select>
                </template>
              </el-table-column>
              <el-table-column label="当前库存" width="100">
                <template #default="scope">
                  <span class="inventory-count">{{ getProductInventory(scope.row.productId, form.warehouseId) }}</span>
                </template>
              </el-table-column>
              <el-table-column prop="quantity" label="采购数量" width="120">
                <template #default="scope">
                  <el-input-number v-model="scope.row.quantity" :min="1" :max="9999" class="input-number"></el-input-number>
                </template>
              </el-table-column>
              <el-table-column prop="unitPrice" label="单价" width="120">
                <template #default="scope">
                  <el-input-number v-model="scope.row.unitPrice" :min="0" :step="0.01" :precision="2" class="input-number"></el-input-number>
                </template>
              </el-table-column>
              <el-table-column prop="amount" label="金额" width="120">
                <template #default="scope">
                  <span class="amount-text">¥{{ ((scope.row.quantity || 0) * (scope.row.unitPrice || 0)).toFixed(2) }}</span>
                </template>
              </el-table-column>
              <el-table-column label="操作" width="80">
                <template #default="scope">
                  <el-button link @click="removeOrderItem(scope.$index)" danger>删除</el-button>
                </template>
              </el-table-column>
            </el-table>
            <div class="total-row">
              <span class="total-label">订单总计：</span>
              <span class="total-amount">¥{{ getTotalAmount().toFixed(2) }}</span>
            </div>
          </div>
        </el-form-item>
        <el-form-item label="备注">
          <el-input v-model="form.remark" type="textarea" :rows="3"></el-input>
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="dialogVisible = false">取消</el-button>
          <el-button type="primary" @click="handleSubmit">确定</el-button>
        </span>
      </template>
    </el-dialog>

    <el-dialog
      title="订单详情"
      v-model="viewVisible"
      width="800px"
      :fullscreen="isMobile"
    >
      <el-form :model="viewForm" label-width="100px">
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="订单号">
              <el-tag type="primary">{{ viewForm.order_no }}</el-tag>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="状态">
              <el-tag :type="getStatusTagType(viewForm.status)">
                {{ getStatusText(viewForm.status) }}
              </el-tag>
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="供应商">
              {{ viewForm.supplier?.name || '-' }}
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="入库仓库">
              {{ viewForm.warehouse?.name || '-' }}
            </el-form-item>
          </el-col>
        </el-row>
        <el-form-item label="订单金额">
          <span class="amount-large">¥{{ viewForm.total_amount?.toFixed(2) }}</span>
        </el-form-item>
        <el-form-item label="备注">
          {{ viewForm.remark || '-' }}
        </el-form-item>
        <el-form-item label="创建时间">
          {{ viewForm.created_at || '-' }}
        </el-form-item>
        <el-form-item label="订单明细">
          <el-table :data="viewForm.items" style="width: 100%" border>
            <el-table-column prop="product.name" label="商品名称" min-width="200"></el-table-column>
            <el-table-column prop="product.code" label="产品编码" width="120"></el-table-column>
            <el-table-column prop="unit_name" label="单位" width="80"></el-table-column>
            <el-table-column prop="ratio" label="换算比例" width="100">
              <template #default="scope">{{ scope.row.ratio || 1 }}:1</template>
            </el-table-column>
            <el-table-column prop="quantity" label="采购数量" width="100"></el-table-column>
            <el-table-column prop="base_qty" label="主单位数量" width="100"></el-table-column>
            <el-table-column prop="unit_price" label="单价" width="100">
              <template #default="scope">¥{{ (scope.row.unit_price || 0).toFixed(2) }}</template>
            </el-table-column>
            <el-table-column prop="amount" label="金额" width="100">
              <template #default="scope">¥{{ (scope.row.amount || 0).toFixed(2) }}</template>
            </el-table-column>
          </el-table>
        </el-form-item>
        <el-form-item label="状态日志">
          <el-table :data="orderLogs" style="width: 100%" border stripe>
            <el-table-column prop="operator_name" label="操作人" width="100"></el-table-column>
            <el-table-column prop="before_status_desc" label="原状态" width="100">
              <template #default="scope">
                <el-tag type="info">{{ scope.row.before_status_desc }}</el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="after_status_desc" label="新状态" width="100">
              <template #default="scope">
                <el-tag>{{ scope.row.after_status_desc }}</el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="remark" label="备注" min-width="150"></el-table-column>
            <el-table-column prop="created_at" label="操作时间" width="160"></el-table-column>
          </el-table>
          <el-empty v-if="orderLogs.length === 0" description="暂无状态变更记录" :image-size="60"></el-empty>
        </el-form-item>
      </el-form>
    </el-dialog>

    <el-dialog
      title="采购入库"
      v-model="inboundVisible"
      width="500px"
    >
      <el-form :model="inboundForm" :rules="inboundRules" ref="inboundFormRef" label-width="100px">
        <el-form-item label="仓库" prop="warehouse_id">
          <el-select v-model="inboundForm.warehouse_id" placeholder="请选择仓库">
            <el-option v-for="wh in warehouses" :key="wh.id" :label="wh.name" :value="wh.id"></el-option>
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="inboundVisible = false">取消</el-button>
          <el-button type="primary" @click="handleInboundSubmit">确认入库</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { More, Setting, Edit, Delete } from '@element-plus/icons-vue'
import { purchaseApi } from '../../api/purchase'
import { productApi } from '../../api/product'
import { inventoryApi } from '../../api/inventory'
import { orderLogApi } from '../../api/orderLog'
import { supplierApi } from '../../api/supplier'
import { warehouseApi } from '../../api/warehouse'

const orders = ref([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(10)
const searchQuery = ref('')
const statusFilter = ref(0)
const loading = ref(false)

const dialogVisible = ref(false)
const dialogTitle = ref('新建采购订单')
const form = ref({
  supplierId: '',
  warehouseId: '',
  items: [],
  remark: ''
})
const formRef = ref(null)

const validateItems = (rule, value, callback) => {
  if (!value || value.length === 0) {
    callback(new Error('请至少添加一个商品'))
    return
  }
  const validItems = value.filter(item => item.productId && item.quantity > 0 && item.unitPrice > 0)
  if (validItems.length === 0) {
    callback(new Error('请填写有效的商品信息'))
    return
  }
  callback()
}

const formRules = {
  supplierId: [{ required: true, message: '请选择供应商', trigger: 'change' }],
  warehouseId: [{ required: true, message: '请选择仓库', trigger: 'change' }],
  items: [{ validator: validateItems, trigger: 'change' }]
}

const viewVisible = ref(false)
const viewForm = ref({})
const orderLogs = ref([])

const inboundVisible = ref(false)
const inboundForm = ref({
  order_id: 0,
  warehouse_id: ''
})
const inboundRules = {
  warehouse_id: [{ required: true, message: '请选择仓库', trigger: 'change' }]
}
const inboundFormRef = ref(null)

const suppliers = ref([])
const products = ref([])
const warehouses = ref([])
const inventories = ref({})
const supplierWarehouses = ref({})
const productUnits = ref({})

const isMobile = computed(() => {
  if (typeof window !== 'undefined') {
    return window.innerWidth < 768
  }
  return false
})

const filteredWarehouses = computed(() => {
  if (!form.value.supplierId) {
    return warehouses.value
  }
  const supplierId = form.value.supplierId
  const linkedWarehouseIds = supplierWarehouses.value[supplierId] || []
  if (linkedWarehouseIds.length > 0) {
    return warehouses.value.filter(wh => linkedWarehouseIds.includes(wh.id))
  }
  return warehouses.value
})

const getStatusTagType = (status) => {
  const types = {
    1: 'warning',
    2: 'primary',
    3: 'success',
    4: 'danger'
  }
  return types[status] || 'info'
}

const getStatusText = (status) => {
  const texts = {
    1: '待审核',
    2: '已审核',
    3: '已入库',
    4: '已取消'
  }
  return texts[status] || '未知'
}

const getProductInventory = (productId, warehouseId) => {
  if (!productId) return '-'
  const key = `${productId}-${warehouseId}`
  return inventories.value[key] || 0
}

const getProductUnits = (productId) => {
  if (!productId) return []
  return productUnits.value[productId] || []
}

const loadProductUnits = (productId) => {
  if (!productId || productUnits.value[productId]) {
    return Promise.resolve()
  }
  return productApi.listProductUnit({ product_id: productId, page: 1, page_size: 100 })
    .then(response => {
      if (response.code === 0) {
        productUnits.value[productId] = response.data.units || []
      } else {
        productUnits.value[productId] = []
      }
    })
    .catch(() => {
      productUnits.value[productId] = []
    })
}

const getTotalAmount = () => {
  return form.value.items.reduce((sum, item) => {
    return sum + (item.quantity || 0) * (item.unitPrice || 0)
  }, 0)
}

const handleSupplierChange = () => {
  if (form.value.supplierId) {
    const linkedWarehouseIds = supplierWarehouses.value[form.value.supplierId] || []
    if (linkedWarehouseIds.length > 0) {
      form.value.warehouseId = linkedWarehouseIds[0]
    }
  } else {
    form.value.warehouseId = ''
  }
}

const handleProductChange = (index) => {
  const item = form.value.items[index]
  item.unitId = ''
  item.unitName = ''
  item.ratio = 1
  item.isMain = 0
  item.quantity = 1
  item.unitPrice = 0

  if (item.productId) {
    const product = products.value.find(p => p.id === item.productId)
    
    loadProductUnits(item.productId)
    .then(() => {
      const units = productUnits.value[item.productId] || []
      if (units.length > 0) {
        const mainUnit = units.find(u => u.is_main === 1) || units[0]
        item.unitId = mainUnit.id
        item.unitName = mainUnit.unit_name
        item.ratio = mainUnit.ratio
        item.isMain = mainUnit.is_main
        
        if (product && product.cost_price !== undefined && product.cost_price > 0) {
          item.unitPrice = product.cost_price * mainUnit.ratio
        }
      }
    })

    if (form.value.warehouseId) {
      const key = `${item.productId}-${form.value.warehouseId}`
      inventoryApi.getCurrentStock({
        product_id: item.productId,
        warehouse_id: form.value.warehouseId
      })
      .then(response => {
        if (response.code === 0) {
          inventories.value[key] = response.data.quantity || 0
        } else {
          inventories.value[key] = 0
        }
      })
      .catch(() => {
        inventories.value[key] = 0
      })
    }
  }
}

const handleUnitChange = (index) => {
  const item = form.value.items[index]
  if (item.productId && item.unitId) {
    const units = productUnits.value[item.productId] || []
    const unit = units.find(u => u.id === item.unitId)
    const product = products.value.find(p => p.id === item.productId)
    if (unit && product) {
      item.unitName = unit.unit_name
      item.ratio = unit.ratio
      item.isMain = unit.is_main
      const mainPrice = product.cost_price || 0
      item.unitPrice = mainPrice * unit.ratio
    }
  }
}

const loadOrders = () => {
  purchaseApi.getPurchaseOrderList({
    page: currentPage.value,
    page_size: pageSize.value,
    order_no: searchQuery.value,
    status: statusFilter.value
  })
  .then(response => {
    if (response.code === 0) {
      orders.value = response.data.orders || []
      total.value = response.data.total || 0
    } else {
      ElMessage.error('获取订单失败：' + response.message)
    }
  })
  .catch(error => {
    ElMessage.error('获取订单失败')
  })
}

const loadSuppliers = () => {
  supplierApi.getActiveSupplierList({})
  .then(response => {
    if (response.code === 0) {
      suppliers.value = response.data.suppliers || []
      suppliers.value.forEach(sup => {
        if (sup.warehouseIds) {
          supplierWarehouses.value[sup.id] = sup.warehouseIds
        }
      })
    }
  })
}

const loadProducts = () => {
  productApi.getActiveProductList({ page: 1, page_size: 100 })
  .then(response => {
    if (response.code === 0) {
      products.value = response.data.products || []
    }
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

const handleAdd = () => {
  dialogTitle.value = '新建采购订单'
  form.value = {
    supplierId: '',
    warehouseId: '',
    items: [{ productId: '', unitId: '', quantity: 1, unitPrice: 0, unitName: '', ratio: 1, isMain: 0 }],
    remark: ''
  }
  dialogVisible.value = true
}

const handleView = (order) => {
  purchaseApi.getPurchaseOrder(order.id)
  .then(response => {
    if (response.code === 0) {
      viewForm.value = response.data
      viewVisible.value = true
      loadOrderLogs(order.id, 1)
    }
  })
}

const loadOrderLogs = (orderId, orderType) => {
  orderLogApi.getOrderLogList({ order_id: orderId, order_type: orderType })
  .then(response => {
    if (response.code === 0) {
      orderLogs.value = response.data.logs || []
    }
  })
  .catch(() => {
    orderLogs.value = []
  })
}

const handleApprove = (order) => {
  ElMessageBox.confirm('确定要审核此订单吗？', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消'
  }).then(() => {
    purchaseApi.updatePurchaseOrder({ id: order.id, status: 2 })
    .then(response => {
      if (response.code === 0) {
        ElMessage.success('审核成功')
        loadOrders()
      } else {
        ElMessage.error('审核失败：' + response.message)
      }
    })
  })
}

const handleCancel = (order) => {
  ElMessageBox.confirm('确定要取消此订单吗？', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消'
  }).then(() => {
    purchaseApi.updatePurchaseOrder({ id: order.id, status: 4 })
    .then(response => {
      if (response.code === 0) {
        ElMessage.success('取消成功')
        loadOrders()
      } else {
        ElMessage.error('取消失败：' + response.message)
      }
    })
  })
}

const handleInbound = (order) => {
  inboundForm.value = {
    order_id: order.id,
    warehouse_id: order.warehouse_id || 0
  }
  inboundVisible.value = true
}

const handleInboundSubmit = () => {
  inboundFormRef.value.validate((valid) => {
    if (valid) {
      purchaseApi.purchaseInbound({
        order_id: inboundForm.value.order_id,
        warehouse_id: inboundForm.value.warehouse_id
      })
      .then(response => {
        if (response.code === 0) {
          ElMessage.success('入库成功')
          inboundVisible.value = false
          loadOrders()
        } else {
          ElMessage.error('入库失败：' + response.message)
        }
      })
    }
  })
}

const handleDelete = (id) => {
  ElMessageBox.confirm('确定要删除此订单吗？', '警告', {
    type: 'warning'
  }).then(() => {
    purchaseApi.deletePurchaseOrder({ id })
    .then(response => {
      if (response.code === 0) {
        ElMessage.success('删除成功')
        loadOrders()
      } else {
        ElMessage.error('删除失败：' + response.message)
      }
    })
  })
}

const handleSubmit = () => {
  formRef.value.validate((valid) => {
    if (valid) {
      const items = form.value.items.filter(item => item.productId && item.quantity > 0 && item.unitPrice > 0).map(item => ({
        product_id: item.productId,
        unit_id: item.unitId || 0,
        unit_name: item.unitName || '',
        ratio: item.ratio || 1,
        is_main: item.isMain || 0,
        quantity: item.quantity,
        unit_price: item.unitPrice
      }))
      purchaseApi.createPurchaseOrder({
        supplier_id: form.value.supplierId,
        warehouse_id: form.value.warehouseId,
        items,
        remark: form.value.remark
      })
      .then(response => {
        if (response.code === 0) {
          ElMessage.success('创建成功')
          dialogVisible.value = false
          loadOrders()
        } else {
          ElMessage.error('创建失败：' + response.message)
        }
      })
    }
  })
}

const addOrderItem = () => {
  form.value.items.push({ productId: '', unitId: '', quantity: 1, unitPrice: 0, unitName: '', ratio: 1, isMain: 0 })
}

const removeOrderItem = (index) => {
  form.value.items.splice(index, 1)
}

const handleSearch = () => {
  currentPage.value = 1
  loadOrders()
}

const handleSizeChange = (size) => {
  pageSize.value = size
  loadOrders()
}

const handleCurrentChange = (current) => {
  currentPage.value = current
  loadOrders()
}

onMounted(() => {
  loadOrders()
  loadSuppliers()
  loadProducts()
  loadWarehouses()
})
</script>

<style scoped>
.purchase-management {
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

.order-items {
  margin-top: 10px;
}

.order-items .total-row {
  display: flex;
  justify-content: flex-end;
  align-items: center;
  margin-top: 10px;
  padding-right: 20px;
}

.total-label {
  font-weight: 500;
  margin-right: 10px;
}

.total-amount {
  font-size: 18px;
  font-weight: 600;
  color: #e6a23c;
}

.empty-tip {
  font-size: 12px;
  color: #999;
  margin-top: 5px;
}

.inventory-count {
  color: #67c23a;
  font-weight: 500;
}

.amount-text {
  color: #409eff;
  font-weight: 500;
}

.amount-large {
  font-size: 18px;
  font-weight: 600;
  color: #e6a23c;
}

.product-select {
  width: 100%;
}

.unit-select {
  width: 100%;
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
:deep(.el-select),
:deep(.el-input-number) {
  width: 100%;
}

:deep(.el-input-number) {
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

  .search-container .el-input {
    width: 100% !important;
  }

  .pagination-container {
    justify-content: center;
  }

  :deep(.el-form-item) {
    margin-bottom: 16px;
  }

  :deep(.el-form-item__label) {
    width: 80px;
    margin-right: 15px;
  }

  .order-items .total-row {
    padding-right: 10px;
  }
}
</style>
