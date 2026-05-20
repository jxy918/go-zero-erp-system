<template>
  <div class="category-management">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>分类管理</span>
          <el-button type="primary" @click="handleAdd">添加分类</el-button>
        </div>
      </template>
      <div class="search-container">
        <el-input
          v-model="searchQuery"
          placeholder="请输入分类名称"
          style="width: 300px"
          prefix-icon="el-icon-search"
        ></el-input>
        <el-button type="primary" @click="handleSearch">搜索</el-button>
      </div>
      <el-tree
        :data="categoryTree"
        node-key="id"
        default-expand-all
        class="category-tree"
      >
        <template #default="{ node, data }">
          <span class="tree-node">
            <span class="tree-node-content">
              <span class="icon category-icon">
                <el-icon><Folder /></el-icon>
              </span>
              <span>{{ data.name }}</span>
              <span v-if="data.code" class="code">({{ data.code }})</span>
              <span :class="['status-tag', data.status === 1 ? 'status-enabled' : 'status-disabled']">
                {{ data.status === 1 ? '启用' : '禁用' }}
              </span>
            </span>
            <span class="tree-actions">
              <el-button type="primary" size="small" @click="handleEdit(data)">
                <template #icon><Edit /></template>
                <span>编辑</span>
              </el-button>
              <el-button type="danger" size="small" @click="handleDelete(data.id)">
                <template #icon><Delete /></template>
                <span>删除</span>
              </el-button>
            </span>
          </span>
        </template>
      </el-tree>
    </el-card>

    <el-dialog
      :title="dialogTitle"
      v-model="dialogVisible"
      width="500px"
    >
      <el-form :model="form" :rules="rules" ref="formRef">
        <el-form-item label="分类名称" prop="name">
          <el-input v-model="form.name" placeholder="请输入分类名称"></el-input>
        </el-form-item>
        <el-form-item label="分类编码" prop="code">
          <el-input v-model="form.code" placeholder="请输入分类编码"></el-input>
        </el-form-item>
        <el-form-item label="上级分类">
          <el-select v-model="form.parentId" placeholder="请选择上级分类">
            <el-option :key="0" :label="'- 顶级分类 -'" :value="0"></el-option>
            <el-option v-for="cat in parentCategories" :key="cat.id" :label="cat.name" :value="cat.id"></el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="排序号">
          <el-input-number v-model="form.sort" :min="0" placeholder="请输入排序号"></el-input-number>
        </el-form-item>
        <el-form-item label="描述">
          <el-input v-model="form.desc" type="textarea" placeholder="请输入分类描述"></el-input>
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
import { ElMessage, ElMessageBox } from 'element-plus'
import { Folder, Edit, Delete } from '@element-plus/icons-vue'
import { productApi } from '../../api/product'

const categories = ref([])
const searchQuery = ref('')

const dialogVisible = ref(false)
const dialogTitle = ref('添加分类')
const form = ref({
  id: 0,
  name: '',
  code: '',
  parentId: '',
  sort: 0,
  desc: '',
  status: true
})
const rules = {
  name: [{ required: true, message: '请输入分类名称', trigger: 'blur' }],
  code: [{ required: true, message: '请输入分类编码', trigger: 'blur' }]
}
const formRef = ref(null)

const parentCategories = ref([])

const buildTree = (items, parentId = 0) => {
  return items
    .filter(item => (item.parentId || item.parent_id || 0) === parentId)
    .map(item => ({
      ...item,
      parentId: item.parentId || item.parent_id || 0,
      children: buildTree(items, item.id)
    }))
    .sort((a, b) => (a.sort || 0) - (b.sort || 0))
}

const categoryTree = computed(() => {
  return buildTree(categories.value, 0)
})

const loadCategories = () => {
  productApi.getCategoryList({
    page: 1,
    page_size: 100,
    name: searchQuery.value
  })
  .then(response => {
    if (response.code === 0) {
      categories.value = response.data.categories || []
    } else {
      ElMessage.error('获取分类列表失败：' + response.message)
    }
  })
  .catch(error => {
    console.error('获取分类列表失败:', error)
    ElMessage.error('获取分类列表失败')
  })
}

const handleAdd = () => {
  dialogTitle.value = '添加分类'
  form.value = {
    id: 0,
    name: '',
    code: '',
    parentId: '',
    sort: 0,
    desc: '',
    status: true
  }
  dialogVisible.value = true
  loadParentCategories()
}

const handleEdit = (category) => {
  dialogTitle.value = '编辑分类'
  form.value = {
    id: category.id,
    name: category.name,
    code: category.code,
    parentId: category.parentId || category.parent_id || 0,
    sort: category.sort,
    desc: category.desc,
    status: category.status === 1
  }
  dialogVisible.value = true
  loadParentCategories()
}

const loadParentCategories = () => {
  productApi.getCategoryList({})
  .then(response => {
    if (response.code === 0) {
      parentCategories.value = (response.data.categories || []).filter(c => c.id !== form.value.id)
    }
  })
}

const handleDelete = (id) => {
  ElMessageBox.confirm('确定要删除这个分类吗？', '警告', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(() => {
    productApi.deleteCategory({ id })
    .then(response => {
      if (response.code === 0) {
        ElMessage.success('删除成功')
        loadCategories()
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
      const formData = {
        name: form.value.name,
        code: form.value.code,
        parent_id: parseInt(form.value.parentId) || 0,
        sort: parseInt(form.value.sort) || 0,
        status: form.value.status ? 1 : 0,
        desc: form.value.desc || ''
      }
      if (form.value.id) {
        productApi.updateCategory({
          ...formData,
          id: form.value.id
        })
        .then(response => {
          if (response.code === 0) {
            ElMessage.success('更新成功')
            dialogVisible.value = false
            loadCategories()
          } else {
            ElMessage.error('操作失败：' + response.message)
          }
        })
      } else {
        productApi.createCategory(formData)
        .then(response => {
          if (response.code === 0) {
            ElMessage.success('创建成功')
            dialogVisible.value = false
            loadCategories()
          } else {
            ElMessage.error('操作失败：' + response.message)
          }
        })
      }
    }
  })
}

const handleSearch = () => {
  loadCategories()
}

onMounted(() => {
  loadCategories()
})
</script>

<style scoped>
.category-management {
  padding: 0;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.search-container {
  margin: 20px 0;
  display: flex;
  gap: 10px;
}

.category-tree {
  margin-top: 20px;
}

.tree-node {
  display: flex;
  justify-content: space-between;
  align-items: center;
  width: 100%;
}

.tree-node-content {
  display: flex;
  align-items: center;
  gap: 8px;
  flex: 1;
}

.tree-node-content .icon {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 24px;
  height: 24px;
  border-radius: 4px;
}

.tree-node-content .category-icon {
  background-color: #fff7e6;
}

.tree-node-content .category-icon :deep(.el-icon) {
  color: #fa8c16;
}

.tree-node-content .code {
  font-size: 12px;
  color: #909399;
}

.tree-node-content .status-tag {
  font-size: 11px;
  padding: 2px 8px;
  border-radius: 10px;
  margin-left: 8px;
}

.tree-node-content .status-enabled {
  background-color: #f6ffed;
  color: #52c41a;
  border: 1px solid #b7eb8f;
}

.tree-node-content .status-disabled {
  background-color: #fff2f0;
  color: #ff4d4f;
  border: 1px solid #ffccc7;
}

.tree-actions {
  display: flex;
  gap: 5px;
}

:deep(.tree-actions .el-button) {
  min-width: 60px;
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
  
  .tree-node {
    flex-direction: column;
    align-items: flex-start;
    gap: 10px;
  }
  
  .tree-actions {
    width: 100%;
    justify-content: flex-start;
  }
  
  .el-button {
    font-size: 12px;
    padding: 4px 8px;
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
  
  .el-tree {
    font-size: 12px;
  }
}
</style>