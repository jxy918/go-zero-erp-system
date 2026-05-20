<template>
  <div class="user-management">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>用户管理</span>
          <el-button type="primary" @click="handleAddUser" v-has-permission="'btn_user_create'">添加用户</el-button>
        </div>
      </template>
      <div class="description-box">
        <span class="description-icon">ℹ️</span>
        <span class="description-text">用户管理：管理系统用户账号，支持创建、编辑、删除用户，可分配角色和设置用户状态。</span>
      </div>
      <div class="search-container">
        <el-input
          v-model="searchQuery"
          placeholder="请输入用户名"
          style="width: 300px"
          prefix-icon="el-icon-search"
        ></el-input>
        <el-button type="primary" @click="handleSearch">搜索</el-button>
      </div>
      <el-table 
        :data="users" 
        style="width: 100%"
        stripe
        :header-cell-style="{ backgroundColor: 'var(--card-light)', fontWeight: 'bold' }"
        :cell-style="{ transition: 'background-color 0.3s ease' }"
        :row-style="{ transition: 'background-color 0.3s ease' }"
        @row-mouse-enter="rowHover = true"
        @row-mouse-leave="rowHover = false"
      >
        <el-table-column prop="id" label="ID" width="80" fixed></el-table-column>
        <el-table-column prop="username" label="用户名"></el-table-column>
        <el-table-column prop="nickname" label="昵称"></el-table-column>
        <el-table-column prop="email" label="邮箱"></el-table-column>
        <el-table-column prop="phone" label="电话"></el-table-column>
        <el-table-column prop="status" label="状态" width="100">
          <template #default="scope">
            <el-tag :type="scope.row.status === 1 ? 'success' : 'danger'">
              {{ scope.row.status === 1 ? '启用' : '禁用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="roles" label="角色" width="200">
          <template #default="scope">
            <el-tag v-for="role in (scope.row.roles || [])" :key="role.id" size="small" style="margin-right: 5px">
              {{ role.name }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="scope">
            <el-dropdown trigger="click" v-if="hasAnyUserPermission">
              <el-button link size="small">
                <el-icon><More /></el-icon> 更多
              </el-button>
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item @click="handleEditUser(scope.row)" v-if="permission.has(PermissionCode.USER_UPDATE)">
                    <el-icon><Edit /></el-icon>
                    <span>编辑</span>
                  </el-dropdown-item>
                  <el-dropdown-item @click="handleDeleteUser(scope.row.id)" danger v-if="permission.has(PermissionCode.USER_DELETE)">
                    <el-icon><Delete /></el-icon>
                    <span>删除</span>
                  </el-dropdown-item>
                  <el-dropdown-item @click="handleAssignRoles(scope.row)" v-if="permission.has(PermissionCode.USER_ASSIGN_ROLES)">
                    <el-icon><Setting /></el-icon>
                    <span>分配角色</span>
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

    <!-- 添加/编辑用户对话框 -->
    <el-dialog
      :title="dialogTitle"
      v-model="dialogVisible"
      width="500px"
    >
      <el-form :model="userForm" :rules="rules" ref="userFormRef">
        <el-form-item label="用户名" prop="username">
          <el-input v-model="userForm.username" placeholder="请输入用户名"></el-input>
        </el-form-item>
        <el-form-item label="密码" prop="password" v-if="!userForm.id">
          <el-input
            v-model="userForm.password"
            type="password"
            placeholder="请输入密码"
            show-password
          ></el-input>
        </el-form-item>
        <el-form-item label="昵称" prop="nickname">
          <el-input v-model="userForm.nickname" placeholder="请输入昵称"></el-input>
        </el-form-item>
        <el-form-item label="邮箱" prop="email">
          <el-input v-model="userForm.email" placeholder="请输入邮箱"></el-input>
        </el-form-item>
        <el-form-item label="电话" prop="phone">
          <el-input v-model="userForm.phone" placeholder="请输入电话"></el-input>
        </el-form-item>
        <el-form-item label="状态">
          <el-switch v-model="userForm.status"></el-switch>
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="dialogVisible = false">取消</el-button>
          <el-button type="primary" @click="handleSubmit">确定</el-button>
        </span>
      </template>
    </el-dialog>

    <!-- 分配角色对话框 -->
    <el-dialog
      title="分配角色"
      v-model="assignRolesVisible"
      width="500px"
    >
      <el-form :model="assignRolesForm">
        <el-form-item label="用户">
          <el-tag>{{ assignRolesForm.username }}</el-tag>
        </el-form-item>
        <el-form-item label="角色">
          <el-checkbox-group v-model="assignRolesForm.roleIds">
            <el-checkbox v-for="role in roles" :key="role.id" :label="role.id">
              {{ role.name }}
            </el-checkbox>
          </el-checkbox-group>
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="assignRolesVisible = false">取消</el-button>
          <el-button type="primary" @click="handleAssignRolesSubmit">确定</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { More, Edit, Delete, Setting } from '@element-plus/icons-vue'
import request from '../../api/request'
import { permission, PermissionCode } from '../../utils/permission'

const hasAnyUserPermission = computed(() => {
  const result = permission.has([
    PermissionCode.USER_UPDATE,
    PermissionCode.USER_DELETE,
    PermissionCode.USER_ASSIGN_ROLES
  ])
  console.log('hasAnyUserPermission computed:', result, 'PermissionCode values:', PermissionCode)
  return result
})

const users = ref([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(10)
const searchQuery = ref('')

const dialogVisible = ref(false)
const dialogTitle = ref('添加用户')
const userForm = ref({
  id: 0,
  username: '',
  password: '',
  nickname: '',
  email: '',
  phone: '',
  status: 1
})
const rules = {
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' },
    { pattern: /^[A-Za-z0-9]{1,20}$/, message: '用户名只支持英文大小写和数字，不超过20个字符', trigger: 'blur' }
  ],
  password: [{ required: true, message: '请输入密码', trigger: 'blur' }],
  nickname: [
    { required: true, message: '请输入昵称', trigger: 'blur' },
    { max: 20, message: '昵称不超过20个字符', trigger: 'blur' }
  ],
  email: [
    { required: true, message: '请输入邮箱', trigger: 'blur' },
    { type: 'email', message: '请输入正确的邮箱格式', trigger: 'blur' }
  ],
  phone: [
    { pattern: /^1[3-9]\d{9}$/, message: '请输入正确的手机号码格式', trigger: 'blur' }
  ]
}
const userFormRef = ref(null)

const assignRolesVisible = ref(false)
const assignRolesForm = ref({
  userId: 0,
  username: '',
  roleIds: []
})
const roles = ref([])

const loadUsers = () => {
  // 从 API 获取用户列表
  request.get('/user/list', {
    params: {
      page: currentPage.value,
      page_size: pageSize.value,
      username: searchQuery.value
    }
  })
  .then(response => {
    console.log('用户列表响应:', response)
    if (response.code === 0) {
      users.value = response.data.users || []
      total.value = response.data.total || 0
    } else {
      ElMessage.error('获取用户列表失败：' + response.message)
    }
  })
  .catch(error => {
    console.error('获取用户列表失败:', error)
    ElMessage.error('获取用户列表失败')
  })
}

const loadRoles = () => {
  // 从 API 获取角色列表
  request.get('/role/list', {
    params: {
      page: 1,
      page_size: 100
    }
  })
  .then(response => {
    if (response.code === 0) {
      roles.value = response.data.roles
    } else {
      ElMessage.error('获取角色列表失败：' + response.message)
    }
  })
  .catch(error => {
    console.error('获取角色列表失败:', error)
    ElMessage.error('获取角色列表失败')
  })
}

const handleAddUser = () => {
  if (!permission.check(PermissionCode.USER_CREATE)) {
    return
  }
  dialogTitle.value = '添加用户'
  userForm.value = {
    id: 0,
    username: '',
    password: '',
    nickname: '',
    email: '',
    phone: '',
    status: 1
  }
  dialogVisible.value = true
}

const handleEditUser = (user) => {
  if (!permission.check(PermissionCode.USER_UPDATE)) {
    return
  }
  dialogTitle.value = '编辑用户'
  // 将数字状态转换为布尔值，1表示启用（true），0表示禁用（false）
  const userWithBooleanStatus = {
    ...user,
    status: user.status === 1
  }
  userForm.value = userWithBooleanStatus
  dialogVisible.value = true
}

const handleDeleteUser = (id) => {
  if (!permission.check(PermissionCode.USER_DELETE)) {
    return
  }
  console.log('准备删除用户，ID:', id)
  ElMessageBox.confirm('确定要删除这个用户吗？', '警告', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(() => {
    console.log('用户确认删除，发送删除请求')
    const url = `/user/delete?id=${id}`
    console.log('删除请求URL:', url)
    // 调用 API 删除用户
    request.post(url)
    .then(response => {
      console.log('删除用户响应:', response)
      if (response.code === 0) {
        ElMessage.success('删除成功')
        loadUsers()
      } else {
        ElMessage.error('删除失败：' + response.message)
      }
    })
    .catch(error => {
      console.error('删除用户失败:', error)
      console.error('错误详情:', error.response)
      ElMessage.error('删除用户失败')
    })
  })
}

const handleAssignRoles = async (user) => {
  if (!permission.check(PermissionCode.USER_ASSIGN_ROLES)) {
    return
  }
  
  // 确保角色列表已加载（按需加载）
  if (roles.value.length === 0) {
    await new Promise((resolve, reject) => {
      request.get('/role/list', {
        params: {
          page: 1,
          page_size: 100
        }
      })
      .then(response => {
        if (response.code === 0) {
          roles.value = response.data.roles
          resolve()
        } else {
          ElMessage.error('获取角色列表失败：' + response.message)
          reject()
        }
      })
      .catch(error => {
        console.error('获取角色列表失败:', error)
        ElMessage.error('获取角色列表失败')
        reject()
      })
    })
  }
  
  assignRolesForm.value = {
    userId: user.id,
    username: user.username,
    roleIds: user.roles ? user.roles.map(role => role.id) : []
  }
  assignRolesVisible.value = true
}

const handleSubmit = () => {
  userFormRef.value.validate((valid) => {
    if (valid) {
      // 将布尔值转换为数字，true表示启用（1），false表示禁用（0）
      const formData = {
        ...userForm.value,
        status: userForm.value.status ? 1 : 0
      }
      // 调用 API 创建或更新用户
      if (userForm.value.id) {
        // 更新用户 - 使用 POST 方法
        request.post('/user/update', formData)
        .then(response => {
          if (response.code === 0) {
            ElMessage.success('更新成功')
            dialogVisible.value = false
            loadUsers()
          } else {
            ElMessage.error('操作失败：' + response.message)
          }
        })
        .catch(error => {
          console.error('操作失败:', error)
          ElMessage.error('操作失败')
        })
      } else {
        // 创建用户 - 使用 POST 方法
        request.post('/user/create', formData)
        .then(response => {
          if (response.code === 0) {
            ElMessage.success('创建成功')
            dialogVisible.value = false
            loadUsers()
          } else {
            ElMessage.error('操作失败：' + response.message)
          }
        })
        .catch(error => {
          console.error('操作失败:', error)
          ElMessage.error('操作失败')
        })
      }
    }
  })
}

const handleAssignRolesSubmit = () => {
  // 调用 API 分配角色
  request.post('/user/assign-roles', {
    user_id: assignRolesForm.value.userId,
    role_ids: assignRolesForm.value.roleIds
  })
  .then(response => {
    if (response.code === 0) {
      ElMessage.success('分配角色成功')
      assignRolesVisible.value = false
      loadUsers()
    } else {
      ElMessage.error('分配角色失败：' + response.message)
    }
  })
  .catch(error => {
    console.error('分配角色失败:', error)
    ElMessage.error('分配角色失败')
  })
}

const handleSearch = () => {
  // 模拟搜索
  loadUsers()
}

const handleSizeChange = (size) => {
  pageSize.value = size
  loadUsers()
}

const handleCurrentChange = (current) => {
  currentPage.value = current
  loadUsers()
}

onMounted(() => {
  loadUsers()
})
</script>

<style scoped>
.user-management {
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

.search-container {
  margin: 20px 0;
  display: flex;
  gap: 10px;
}

.pagination-container {
  margin-top: 20px;
  margin-bottom: 20px;
  display: flex;
  justify-content: flex-end;
  width: 100%;
  padding: 0 10px;
}

/* 表单样式 - 确保输入框对齐 */
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

/* 表格样式 */
:deep(.el-table__header th) {
  border-radius: var(--border-radius) var(--border-radius) 0 0 !important;
}

:deep(.el-table__row:hover) {
  background-color: rgba(76, 175, 80, 0.05) !important;
}

body.dark-mode :deep(.el-table__row:hover) {
  background-color: rgba(76, 175, 80, 0.15) !important;
}

body.dark-mode :deep(.el-table__header th) {
  background-color: var(--card-dark) !important;
  color: var(--text-dark) !important;
}

/* 响应式设计 */
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
  
  .el-table {
    font-size: 12px;
  }
  
  .el-table th,
  .el-table td {
    padding: 8px 4px;
  }
  
  .el-button {
    font-size: 12px;
    padding: 4px 8px;
  }
  
  .el-dialog {
    width: 90% !important;
    margin: 20px auto !important;
  }
  
  /* 响应式表单样式 */
  :deep(.el-form-item) {
    margin-bottom: 16px;
  }
  
  :deep(.el-form-item__label) {
    width: 80px;
    margin-right: 15px;
  }
}
</style>