<template>
  <div class="customer-management">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>客户管理</span>
          <el-button type="primary" @click="handleAdd">添加客户</el-button>
        </div>
      </template>
      <div class="query-bar">
        <el-input
          v-model="searchQuery"
          placeholder="请输入客户名称或编码"
          style="width: 300px"
          prefix-icon="el-icon-search"
        ></el-input>
        <el-button type="primary" @click="handleSearch">搜索</el-button>
      </div>
      <el-table 
        :data="customers" 
        style="width: 100%"
        border
        stripe
        v-loading="loading"
      >
        <el-table-column prop="id" label="ID" width="80" fixed></el-table-column>
        <el-table-column prop="name" label="客户名称"></el-table-column>
        <el-table-column prop="code" label="客户编码"></el-table-column>
        <el-table-column prop="contact" label="联系人"></el-table-column>
        <el-table-column prop="phone" label="联系电话"></el-table-column>
        <el-table-column prop="address" label="地址"></el-table-column>
        <el-table-column prop="status" label="状态" width="100">
          <template #default="scope">
            <el-tag :type="scope.row.status === 1 ? 'success' : 'danger'">
              {{ scope.row.status === 1 ? '启用' : '禁用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="130" fixed="right" align="center">
          <template #default="scope">
            <div class="action-buttons">
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
        <el-form-item label="客户名称" prop="name">
          <el-input v-model="form.name" placeholder="请输入客户名称"></el-input>
        </el-form-item>
        <el-form-item label="客户编码" prop="code">
          <el-input v-model="form.code" placeholder="请输入客户编码"></el-input>
        </el-form-item>
        <el-form-item label="联系人" prop="contact">
          <el-input v-model="form.contact" placeholder="请输入联系人"></el-input>
        </el-form-item>
        <el-form-item label="联系电话" prop="phone">
          <el-input v-model="form.phone" placeholder="请输入联系电话"></el-input>
        </el-form-item>
        <el-form-item label="邮箱" prop="email">
          <el-input v-model="form.email" placeholder="请输入邮箱"></el-input>
        </el-form-item>
        <el-form-item label="地址">
          <el-input v-model="form.address" type="textarea" placeholder="请输入地址"></el-input>
        </el-form-item>
        <el-form-item label="描述">
          <el-input v-model="form.desc" type="textarea" placeholder="请输入描述"></el-input>
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
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { customerApi } from '../../api/customer'

const customers = ref([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(10)
const searchQuery = ref('')
const loading = ref(false)

const dialogVisible = ref(false)
const dialogTitle = ref('添加客户')
const form = ref({
  id: 0,
  name: '',
  code: '',
  contact: '',
  phone: '',
  email: '',
  address: '',
  desc: '',
  status: true
})
const rules = {
  name: [
    { required: true, message: '请输入客户名称', trigger: 'blur' },
    { max: 100, message: '客户名称不超过100个字符', trigger: 'blur' }
  ],
  code: [
    { required: true, message: '请输入客户编码', trigger: 'blur' },
    { max: 50, message: '客户编码不超过50个字符', trigger: 'blur' },
    { pattern: /^[A-Za-z0-9_]{1,50}$/, message: '客户编码只支持字母、数字和下划线', trigger: 'blur' }
  ],
  contact: [
    { max: 50, message: '联系人不超过50个字符', trigger: 'blur' }
  ],
  phone: [
    { pattern: /^1[3-9]\d{9}$/, message: '请输入正确的手机号码', trigger: 'blur' }
  ],
  email: [
    { type: 'email', message: '请输入正确的邮箱格式', trigger: 'blur' }
  ]
}
const formRef = ref(null)

const loadCustomers = () => {
  customerApi.getCustomerList({
    page: currentPage.value,
    page_size: pageSize.value,
    name: searchQuery.value,
    code: searchQuery.value
  })
  .then(response => {
    if (response.code === 0) {
      customers.value = response.data.customers || []
      total.value = response.data.total || 0
    } else {
      ElMessage.error('获取客户列表失败：' + response.message)
    }
  })
  .catch(error => {
    console.error('获取客户列表失败:', error)
    ElMessage.error('获取客户列表失败')
  })
}

const handleAdd = () => {
  dialogTitle.value = '添加客户'
  form.value = {
    id: 0,
    name: '',
    code: '',
    contact: '',
    phone: '',
    email: '',
    address: '',
    desc: '',
    status: true
  }
  dialogVisible.value = true
}

const handleEdit = (customer) => {
  dialogTitle.value = '编辑客户'
  form.value = {
    id: customer.id,
    name: customer.name,
    code: customer.code,
    contact: customer.contact,
    phone: customer.phone,
    email: customer.email,
    address: customer.address,
    desc: customer.desc,
    status: customer.status === 1
  }
  dialogVisible.value = true
}

const handleDelete = (id) => {
  ElMessageBox.confirm('确定要删除这个客户吗？', '警告', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(() => {
    customerApi.deleteCustomer({ id })
    .then(response => {
      if (response.code === 0) {
        ElMessage.success('删除成功')
        loadCustomers()
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
        ...form.value,
        status: form.value.status ? 1 : 0
      }
      if (form.value.id) {
        customerApi.updateCustomer(formData)
        .then(response => {
          if (response.code === 0) {
            ElMessage.success('更新成功')
            dialogVisible.value = false
            loadCustomers()
          } else {
            ElMessage.error('操作失败：' + response.message)
          }
        })
      } else {
        customerApi.createCustomer(formData)
        .then(response => {
          if (response.code === 0) {
            ElMessage.success('创建成功')
            dialogVisible.value = false
            loadCustomers()
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
  loadCustomers()
}

const handleSizeChange = (size) => {
  pageSize.value = size
  loadCustomers()
}

const handleCurrentChange = (current) => {
  currentPage.value = current
  loadCustomers()
}

onMounted(() => {
  loadCustomers()
})
</script>

<style scoped>
.customer-management {
  padding: 0;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
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

.search-container {
  margin: 20px 0;
  display: flex;
  gap: 10px;
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