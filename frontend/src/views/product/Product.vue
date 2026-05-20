<template>
  <div class="product-management">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>产品管理</span>
          <el-button type="primary" @click="handleAdd">添加产品</el-button>
        </div>
      </template>
      <div class="description-box">
        <span class="description-icon">ℹ️</span>
        <span class="description-text">产品管理：管理系统产品信息，包括产品名称、编码、规格、分类、单价等，支持产品的增删改查操作。</span>
      </div>
      <div class="query-bar">
        <el-input
          v-model="searchQuery"
          placeholder="请输入产品名称或编码"
          style="width: 300px"
          prefix-icon="el-icon-search"
        ></el-input>
        <el-button type="primary" @click="handleSearch">搜索</el-button>
      </div>
      <el-table 
        :data="products" 
        style="width: 100%"
        border
        stripe
        v-loading="loading"
      >
        <el-table-column prop="id" label="ID" width="80" fixed></el-table-column>
        <el-table-column prop="name" label="产品名称"></el-table-column>
        <el-table-column prop="code" label="产品编码"></el-table-column>
        <el-table-column prop="spec" label="规格型号"></el-table-column>
        <el-table-column prop="category.name" label="分类"></el-table-column>
        <el-table-column prop="main_unit" label="主单位" width="80"></el-table-column>
        <el-table-column prop="price" label="单价">
          <template #default="scope">
            ¥{{ scope.row.price.toFixed(2) }}
          </template>
        </el-table-column>
        <el-table-column prop="cost_price" label="成本价">
          <template #default="scope">
            ¥{{ (scope.row.cost_price || 0).toFixed(2) }}
          </template>
        </el-table-column>
        <el-table-column prop="stock" label="库存"></el-table-column>
        <el-table-column prop="min_stock" label="最低库存">
          <template #default="scope">
            {{ scope.row.min_stock || 0 }}
          </template>
        </el-table-column>
        <el-table-column prop="safety_stock" label="安全库存">
          <template #default="scope">
            {{ scope.row.safety_stock || 0 }}
          </template>
        </el-table-column>
        <el-table-column prop="max_stock" label="最高库存">
          <template #default="scope">
            {{ scope.row.max_stock || 99999 }}
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="100">
          <template #default="scope">
            <el-tag :type="scope.row.status === 1 ? 'success' : 'danger'">
              {{ scope.row.status === 1 ? '启用' : '禁用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="200" fixed="right" align="center">
          <template #default="scope">
            <div class="action-buttons">
              <el-button size="small" type="success" @click="handleManageUnits(scope.row)">管理单位</el-button>
              <el-button size="small" @click="handleEdit(scope.row)">编辑</el-button>
              <el-button size="small" type="danger" @click="handleDelete(scope.row.id)">删除</el-button>
            </div>
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
      width="500px"
    >
      <el-form :model="form" :rules="rules" ref="formRef">
        <el-form-item label="产品名称" prop="name">
          <el-input v-model="form.name" placeholder="请输入产品名称"></el-input>
        </el-form-item>
        <el-form-item label="产品编码" prop="code">
          <el-input v-model="form.code" placeholder="请输入产品编码"></el-input>
        </el-form-item>
        <el-form-item label="规格型号">
          <el-input v-model="form.spec" placeholder="请输入规格型号"></el-input>
        </el-form-item>
        <el-form-item label="产品分类">
          <el-cascader
            v-model="form.categoryPath"
            :options="categoryTreeOptions"
            :props="{ label: 'name', value: 'id', children: 'children' }"
            placeholder="请选择分类"
            clearable
          ></el-cascader>
        </el-form-item>
        <el-form-item label="单价" prop="price">
          <el-input-number v-model="form.price" :min="0" :step="0.01" placeholder="请输入单价"></el-input-number>
          <span class="field-tip">💡 <span style="color: #d93026;">销售单价：产品对外销售时的价格</span></span>
        </el-form-item>
        <el-form-item label="成本价">
          <el-input-number v-model="form.costPrice" :min="0" :step="0.01" placeholder="请输入成本价"></el-input-number>
          <span class="field-tip">💡 <span style="color: #d93026;">成本价：采购产品时的成本价格</span></span>
        </el-form-item>

        <el-form-item label="最低库存">
          <el-input-number v-model="form.minStock" :min="0" placeholder="请输入最低库存"></el-input-number>
          <span class="field-tip">💡 <span style="color: #d93026;">最低库存：库存低于此值时触发库存预警</span></span>
        </el-form-item>
        <el-form-item label="安全库存">
          <el-input-number v-model="form.safetyStock" :min="0" placeholder="请输入安全库存"></el-input-number>
          <span class="field-tip">💡 <span style="color: #d93026;">安全库存：建议维持的最低安全库存数量</span></span>
        </el-form-item>
        <el-form-item label="最高库存">
          <el-input-number v-model="form.maxStock" :min="0" placeholder="请输入最高库存"></el-input-number>
          <span class="field-tip">💡 <span style="color: #d93026;">最高库存：库存上限，建议不超过此值</span></span>
        </el-form-item>
        <el-form-item label="描述">
          <el-input v-model="form.desc" type="textarea" placeholder="请输入产品描述"></el-input>
        </el-form-item>
        <el-form-item label="状态">
          <el-switch v-model="form.status"></el-switch>
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="dialogVisible = false">取消</el-button>
          <el-button type="primary" @click="handleSubmit">确定</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { productApi } from '../../api/product'

const router = useRouter()
const route = useRoute()
const products = ref([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(10)
const searchQuery = ref('')
const loading = ref(false)

const dialogVisible = ref(false)
const dialogTitle = ref('添加产品')
const form = ref({
  id: 0,
  name: '',
  code: '',
  spec: '',
  categoryPath: [],
  categoryId: '',
  price: 0,
  costPrice: 0,
  minStock: 0,
  safetyStock: 0,
  maxStock: 99999,
  desc: '',
  status: true
})
const rules = {
  name: [
    { required: true, message: '请输入产品名称', trigger: 'blur' },
    { max: 100, message: '产品名称不超过100个字符', trigger: 'blur' }
  ],
  code: [
    { required: true, message: '请输入产品编码', trigger: 'blur' },
    { max: 50, message: '产品编码不超过50个字符', trigger: 'blur' },
    { pattern: /^[A-Za-z0-9_]{1,50}$/, message: '产品编码只支持字母、数字和下划线', trigger: 'blur' }
  ],
  price: [
    { required: true, message: '请输入单价', trigger: 'blur' },
    { validator: (rule, value, callback) => {
      if (value <= 0) {
        callback(new Error('单价必须大于0'))
      } else {
        callback()
      }
    }, trigger: 'blur' }
  ],
  costPrice: [
    { validator: (rule, value, callback) => {
      if (value < 0) {
        callback(new Error('成本价不能为负数'))
      } else {
        callback()
      }
    }, trigger: 'blur' }
  ],

}
const formRef = ref(null)

const categories = ref([])

const buildTreeOptions = (items, parentId = 0) => {
  return items
    .filter(item => (item.parentId || item.parent_id || 0) === parentId)
    .map(item => ({
      id: item.id,
      name: item.name,
      children: buildTreeOptions(items, item.id)
    }))
    .sort((a, b) => (a.sort || 0) - (b.sort || 0))
}

const categoryTreeOptions = computed(() => {
  return buildTreeOptions(categories.value, 0)
})

const loadProducts = () => {
  productApi.getProductList({
    page: currentPage.value,
    page_size: pageSize.value,
    name: searchQuery.value,
    code: searchQuery.value
  })
  .then(response => {
    if (response.code === 0) {
      products.value = response.data.products || []
      total.value = response.data.total || 0
    } else {
      ElMessage.error('获取产品列表失败：' + response.message)
    }
  })
  .catch(error => {
    console.error('获取产品列表失败:', error)
    ElMessage.error('获取产品列表失败')
  })
}

const loadCategories = () => {
  return new Promise((resolve) => {
    productApi.getCategoryList({})
    .then(response => {
      if (response.code === 0) {
        categories.value = response.data.categories || []
      }
      resolve()
    })
    .catch(error => {
      console.error('获取分类失败:', error)
      resolve()
    })
  })
}

const handleAdd = () => {
  dialogTitle.value = '添加产品'
  form.value = {
    id: 0,
    name: '',
    code: '',
    spec: '',
    categoryPath: [],
    categoryId: '',
    price: 0,
    costPrice: 0,
    minStock: 0,
    safetyStock: 0,
    maxStock: 99999,
    desc: '',
    status: true
  }
  dialogVisible.value = true
}

const handleEdit = async (product) => {
  if (categories.value.length === 0) {
    await loadCategories()
  }
  
  dialogTitle.value = '编辑产品'
  const categoryId = product.categoryId || product.CategoryID || (product.category_id || 0)
  form.value = {
    id: product.id,
    name: product.name,
    code: product.code,
    spec: product.spec || '',
    categoryPath: buildCategoryPath(categoryId),
    categoryId: categoryId,
    price: product.price || 0,
    costPrice: product.costPrice || product.cost_price || 0,
    minStock: product.minStock || product.min_stock || 0,
    safetyStock: product.safetyStock || product.safety_stock || 0,
    maxStock: product.maxStock || product.max_stock || 99999,
    desc: product.desc || '',
    status: product.status === 1
  }
  dialogVisible.value = true
}

const buildCategoryPath = (categoryId) => {
  if (!categoryId || categoryId === 0) return []
  const path = []
  let currentId = categoryId
  const findParent = (id) => {
    const cat = categories.value.find(c => c.id === id)
    if (cat) {
      path.unshift(id)
      const parentId = cat.parentId || cat.parent_id || 0
      if (parentId && parentId !== 0) {
        findParent(parentId)
      }
    }
  }
  findParent(categoryId)
  return path
}

const handleManageUnits = (product) => {
  // 跳转到产品计量单位管理页面，并携带产品ID
  router.push({
    path: '/product/unit',
    query: { productId: product.id, productName: product.name }
  })
}

const handleDelete = (id) => {
  ElMessageBox.confirm('确定要删除这个产品吗？', '警告', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(() => {
    productApi.deleteProduct({ id })
    .then(response => {
      if (response.code === 0) {
        ElMessage.success('删除成功')
        loadProducts()
      } else {
        ElMessage.error('删除失败：' + response.message)
      }
    })
    .catch(error => {
      ElMessage.error('删除失败')
    })
  })
}

const handleSubmit = () => {
  formRef.value.validate((valid) => {
    if (valid) {
      const categoryPath = form.value.categoryPath || []
      const categoryId = categoryPath.length > 0 
        ? categoryPath[categoryPath.length - 1] 
        : 0
      
      const formData = {
        name: form.value.name,
        code: form.value.code,
        spec: form.value.spec,
        price: form.value.price,
        cost_price: form.value.costPrice,
        min_stock: form.value.minStock,
        safety_stock: form.value.safetyStock,
        max_stock: form.value.maxStock,
        desc: form.value.desc,
        status: form.value.status ? 1 : 0,
        category_id: categoryId
      }
      if (form.value.id) {
        formData.id = form.value.id
      }
      
      if (form.value.id) {
        productApi.updateProduct(formData)
        .then(response => {
          if (response.code === 0) {
            ElMessage.success('更新成功')
            dialogVisible.value = false
            loadProducts()
          } else {
            ElMessage.error('操作失败：' + response.message)
          }
        })
      } else {
        productApi.createProduct(formData)
        .then(response => {
          if (response.code === 0) {
            ElMessage.success('创建成功')
            dialogVisible.value = false
            loadProducts()
          } else {
            ElMessage.error('操作失败：' + response.message)
          }
        })
      }
    }
  })
}

const handleSearch = () => {
  currentPage.value = 1
  loadProducts()
}

const handleSizeChange = (size) => {
  pageSize.value = size
  loadProducts()
}

const handleCurrentChange = (current) => {
  currentPage.value = current
  loadProducts()
}

onMounted(() => {
  // 如果路由中带有搜索参数，则自动搜索
  if (route.query.search) {
    searchQuery.value = route.query.search
  }
  loadProducts()
  loadCategories()
})
</script>

<style scoped>
.product-management {
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

.field-tip {
  display: block;
  font-size: 12px;
  margin-top: 6px;
  padding-left: 2px;
}

.query-bar {
  margin-bottom: 20px;
  display: flex;
  align-items: center;
  gap: 16px;
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

.pagination-container {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
  width: 100%;
  padding: 0 10px;
}

:deep(.el-form-item) {
  margin-bottom: 20px;
}

:deep(.el-form-item__label) {
  width: 100px;
  text-align: right;
  margin-right: 20px;
}

:deep(.el-form-item__content) {
  flex: 1;
  min-width: 0;
}

:deep(.el-input),
:deep(.el-select),
:deep(.el-input-number),
:deep(.el-cascader) {
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
  
  .el-dialog {
    width: 90% !important;
    margin: 20px auto !important;
  }
  
  :deep(.el-form-item) {
    margin-bottom: 16px;
  }
  
  :deep(.el-form-item__label) {
    width: 80px;
    margin-right: 15px;
  }
}
</style>