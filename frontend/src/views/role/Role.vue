<template>
  <div class="role-management">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>角色管理</span>
          <el-button v-if="hasRoleCreate" type="primary" @click="handleAddRole">添加角色</el-button>
        </div>
      </template>
      <div class="description-box">
        <span class="description-icon">ℹ️</span>
        <span class="description-text">角色管理：管理系统角色，支持创建、编辑、删除角色，可分配权限和菜单给角色，控制用户的系统访问范围。</span>
      </div>
      <div class="search-container">
        <el-input
          v-model="searchQuery"
          placeholder="请输入角色名称"
          style="width: 300px"
          prefix-icon="el-icon-search"
        ></el-input>
        <el-button type="primary" @click="handleSearch">搜索</el-button>
      </div>
      <el-table 
        :data="roles" 
        style="width: 100%"
        stripe
        :header-cell-style="{ backgroundColor: 'var(--card-light)', fontWeight: 'bold' }"
        :cell-style="{ transition: 'background-color 0.3s ease' }"
        :row-style="{ transition: 'background-color 0.3s ease' }"
        @row-mouse-enter="rowHover = true"
        @row-mouse-leave="rowHover = false"
      >
        <el-table-column prop="id" label="ID" width="80" fixed></el-table-column>
        <el-table-column prop="name" label="角色名称"></el-table-column>
        <el-table-column prop="code" label="角色编码"></el-table-column>
        <el-table-column prop="desc" label="角色描述"></el-table-column>
        <el-table-column prop="status" label="状态" width="100">
          <template #default="scope">
            <el-tag :type="scope.row.status === 1 ? 'success' : 'danger'">
              {{ scope.row.status === 1 ? '启用' : '禁用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="permissions" label="权限数量" width="120">
          <template #default="scope">
            {{ scope.row.permissions ? scope.row.permissions.length : 0 }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="scope">
            <!-- 测试：始终显示下拉菜单（暂时移除权限检查） -->
            <el-dropdown trigger="click">
              <el-button link size="small">
                <el-icon><More /></el-icon> 更多
              </el-button>
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item v-if="hasRoleUpdate" @click="handleEditRole(scope.row)">
                    <el-icon><Edit /></el-icon>
                    <span>编辑</span>
                  </el-dropdown-item>
                  <el-dropdown-item v-if="hasRoleDelete" @click="handleDeleteRole(scope.row.id)" danger>
                    <el-icon><Delete /></el-icon>
                    <span>删除</span>
                  </el-dropdown-item>
                  <el-dropdown-item v-if="hasRoleAssign" @click="handleAssignPermissions(scope.row)">
                    <el-icon><Setting /></el-icon>
                    <span>分配权限</span>
                  </el-dropdown-item>
                  <el-dropdown-item v-if="hasRoleAssignMenus" @click="handleAssignMenus(scope.row)">
                    <el-icon><Menu /></el-icon>
                    <span>分配菜单</span>
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

    <!-- 添加/编辑角色对话框 -->
    <el-dialog
      :title="dialogTitle"
      v-model="dialogVisible"
      width="500px"
    >
      <el-form :model="roleForm" :rules="rules" ref="roleFormRef">
        <el-form-item label="角色名称" prop="name">
          <el-input v-model="roleForm.name" placeholder="请输入角色名称"></el-input>
        </el-form-item>
        <el-form-item label="角色编码" prop="code">
          <el-input v-model="roleForm.code" placeholder="请输入角色编码"></el-input>
        </el-form-item>
        <el-form-item label="角色描述" prop="desc">
          <el-input
            v-model="roleForm.desc"
            type="textarea"
            placeholder="请输入角色描述"
          ></el-input>
        </el-form-item>
        <el-form-item label="状态">
          <el-switch v-model="roleForm.status"></el-switch>
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="dialogVisible = false">取消</el-button>
          <el-button type="primary" @click="handleSubmit">确定</el-button>
        </span>
      </template>
    </el-dialog>

    <!-- 分配权限对话框 -->
    <el-dialog
      title="分配权限"
      v-model="assignPermissionsVisible"
      width="600px"
    >
      <el-form :model="assignPermissionsForm">
        <el-form-item label="角色">
          <el-tag>{{ assignPermissionsForm.name }}</el-tag>
        </el-form-item>
        <el-form-item label="权限">
          <el-tree
            ref="permissionTreeRef"
            :key="permissionTreeKey"
            :data="permissions"
            show-checkbox
            node-key="id"
            :default-checked-keys="assignPermissionsForm.permissionIds"
            @check="handlePermissionCheck"
          >
            <template #default="{ node, data }">
              <span class="tree-node">
                <span>{{ data.name }}</span>
                <span v-if="data.code" class="code">({{ data.code }})</span>
              </span>
            </template>
          </el-tree>
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="assignPermissionsVisible = false">取消</el-button>
          <el-button type="primary" @click="handleAssignPermissionsSubmit">确定</el-button>
        </span>
      </template>
    </el-dialog>

    <!-- 分配菜单对话框 -->
    <el-dialog
      title="分配菜单"
      v-model="assignMenusVisible"
      width="600px"
    >
      <el-form :model="assignMenusForm">
        <el-form-item label="角色">
          <el-tag>{{ assignMenusForm.name }}</el-tag>
        </el-form-item>
        <el-form-item label="菜单">
          <el-tree
            ref="menuTreeRef"
            :data="menus"
            show-checkbox
            node-key="id"
            :default-checked-keys="assignMenusForm.menuIds"
            @check-change="handleCheckChange"
          >
            <template #default="{ node, data }">
              <span class="tree-node">
                <span>{{ data.name }}</span>
                <span v-if="data.code" class="code">({{ data.code }})</span>
              </span>
            </template>
          </el-tree>
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="assignMenusVisible = false">取消</el-button>
          <el-button type="primary" @click="handleAssignMenusSubmit">确定</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted, nextTick, watch, computed } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { More, Edit, Delete, Setting, Menu } from '@element-plus/icons-vue'
import { roleApi, request, menuApi } from '../../api'
import { permission, PermissionCode } from '../../utils/permission'
import { useUserStore } from '../../store'

const userStore = useUserStore()

// 直接获取权限码列表
const permissionCodes = computed(() => {
  console.log('[Role.vue] permissionCodes computed, store.permissions:', userStore.permissions)
  const codes = userStore.permissions.map(p => p.code || p.Code)
  console.log('[Role.vue] permissionCodes result:', codes)
  console.log('[Role.vue] has btn_role_assign_menus:', codes.includes('btn_role_assign_menus'))
  return codes
})

const hasAnyRolePermission = computed(() => {
  console.log('[Role.vue] hasAnyRolePermission computed')
  console.log('[Role.vue] current permissionCodes:', permissionCodes.value)
  
  const codes = permissionCodes.value
  const hasUpdate = codes.includes(PermissionCode.ROLE_UPDATE)
  const hasDelete = codes.includes(PermissionCode.ROLE_DELETE)
  const hasAssign = codes.includes(PermissionCode.ROLE_ASSIGN_PERMISSIONS)
  const hasAssignMenus = codes.includes(PermissionCode.ROLE_ASSIGN_MENUS)
  const hasCreate = codes.includes(PermissionCode.ROLE_CREATE)
  
  const result = hasUpdate || hasDelete || hasAssign || hasAssignMenus || hasCreate
  console.log('[Role.vue] hasCreate:', hasCreate, 'hasUpdate:', hasUpdate, 'hasDelete:', hasDelete, 'hasAssign:', hasAssign, 'hasAssignMenus:', hasAssignMenus)
  console.log('[Role.vue] hasAnyRolePermission result:', result)
  
  return result
})

// 单独检查各权限（管理员拥有所有权限）
const hasRoleCreate = computed(() => userStore.isAdmin || permissionCodes.value.includes(PermissionCode.ROLE_CREATE))
const hasRoleUpdate = computed(() => userStore.isAdmin || permissionCodes.value.includes(PermissionCode.ROLE_UPDATE))
const hasRoleDelete = computed(() => userStore.isAdmin || permissionCodes.value.includes(PermissionCode.ROLE_DELETE))
const hasRoleAssign = computed(() => userStore.isAdmin || permissionCodes.value.includes(PermissionCode.ROLE_ASSIGN_PERMISSIONS))
const hasRoleAssignMenus = computed(() => userStore.isAdmin || permissionCodes.value.includes(PermissionCode.ROLE_ASSIGN_MENUS))

const roles = ref([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(10)
const searchQuery = ref('')

const dialogVisible = ref(false)
const dialogTitle = ref('添加角色')
const roleForm = ref({
  id: 0,
  name: '',
  code: '',
  desc: '',
  status: true
})
const rules = {
  name: [
    { required: true, message: '请输入角色名称', trigger: 'blur' },
    { max: 20, message: '角色名称不超过20个字符', trigger: 'blur' }
  ],
  code: [
    { required: true, message: '请输入角色编码', trigger: 'blur' },
    { pattern: /^[A-Za-z0-9_]{1,20}$/, message: '角色编码支持大小写字母、数字和下划线，不超过20个字符', trigger: 'blur' }
  ],
  desc: [
    { max: 200, message: '角色描述不超过200个字符', trigger: 'blur' }
  ]
}
const roleFormRef = ref(null)

const assignPermissionsVisible = ref(false)
const assignPermissionsForm = ref({
  roleId: 0,
  name: '',
  permissionIds: []
})
const permissions = ref([])
const permissionTreeRef = ref(null)
const permissionTreeKey = ref(0)

const assignMenusVisible = ref(false)
const assignMenusForm = ref({
  roleId: 0,
  name: '',
  menuIds: []
})
const menus = ref([])
const menuTreeRef = ref(null)

// 监听分配权限对话框的显示状态
watch(assignPermissionsVisible, async (newValue) => {
  if (newValue) {
    await nextTick()
    if (permissionTreeRef.value) {
      permissionTreeRef.value.setCheckedKeys([])
      await nextTick()
      permissionTreeRef.value.setCheckedKeys(assignPermissionsForm.value.permissionIds)
    }
  }
})

// 监听分配菜单对话框的显示状态
watch(assignMenusVisible, async (newValue) => {
  if (newValue) {
    // 对话框打开后，等待 DOM 更新，然后设置选中的节点
    await nextTick()
    if (menuTreeRef.value) {
      menuTreeRef.value.setCheckedKeys(assignMenusForm.value.menuIds)
    }
  }
})

const loadRoles = () => {
  // 从 API 获取角色列表
  roleApi.getRoleList({
    page: currentPage.value,
    page_size: pageSize.value,
    name: searchQuery.value
  })
  .then(response => {
    if (response.code === 0) {
      roles.value = response.data.roles
      total.value = response.data.total
      // 打印角色信息，检查是否包含 Menus 字段
      console.log('角色列表:', response.data.roles)
      // 检查第一个角色的菜单信息
      if (response.data.roles && response.data.roles.length > 0) {
        console.log('第一个角色的菜单信息:', response.data.roles[0].Menus || response.data.roles[0].menus)
      }
    } else {
      ElMessage.error('获取角色列表失败：' + response.message)
    }
  })
  .catch(error => {
    console.error('获取角色列表失败:', error)
    ElMessage.error('获取角色列表失败')
  })
}

const loadPermissions = () => {
  loadPermissionsAsync().catch(() => {})
}

const loadPermissionsAsync = async () => {
  return new Promise((resolve, reject) => {
    request({
      url: '/permission/list',
      method: 'get',
      params: {
        page: 1,
        page_size: 500
      }
    })
    .then(response => {
      if (response.code === 0) {
        const permList = response.data.permissions
        
        // 分离菜单和按钮权限
        const menus = permList.filter(p => p.type === 1)
        const buttons = permList.filter(p => p.type === 2)
        
        // 构建菜单映射
        const menuMap = {}
        menus.forEach(menu => {
          menuMap[menu.id] = {
            ...menu,
            children: []
          }
        })
        
        // 将按钮权限添加到对应的父菜单下
        // 给按钮权限添加前缀 'p_' 避免与菜单ID冲突
        buttons.forEach(btn => {
          const parentId = btn.parent_id
          if (parentId && menuMap[parentId]) {
            menuMap[parentId].children.push({
              ...btn,
              id: 'p_' + btn.id, // 添加前缀
              children: []
            })
          }
        })
        
        // 构建树形结构（菜单树）
        const result = []
        Object.values(menuMap).forEach(node => {
          const parentId = node.parent_id
          if (parentId && menuMap[parentId]) {
            menuMap[parentId].children.push(node)
          } else {
            result.push(node)
          }
        })
        
        permissions.value = result
        resolve(result)
      } else {
        ElMessage.error('获取权限列表失败：' + response.message)
        reject(response.message)
      }
    })
    .catch(error => {
      console.error('获取权限列表失败:', error)
      ElMessage.error('获取权限列表失败')
      reject(error)
    })
  })
}

const loadMenus = () => {
  // 从 API 获取菜单列表
  menuApi.getMenuList()
  .then(response => {
    if (response.code === 0) {
      // 构建菜单树
      const buildTree = (items, parentId = 0) => {
        return items
          .filter(item => item.parent_id === parentId)
          .map(item => ({
            ...item,
            children: buildTree(items, item.id)
          }))
      }
      // 后端返回的菜单数据在menus字段中
      const menuData = response.data.menus || response.data.Menus || []
      menus.value = buildTree(menuData)
    } else {
      ElMessage.error('获取菜单列表失败：' + response.message)
    }
  })
  .catch(error => {
    console.error('获取菜单列表失败:', error)
    ElMessage.error('获取菜单列表失败')
  })
}

const handleAddRole = () => {
  if (!permission.check(PermissionCode.ROLE_CREATE)) {
    return
  }
  dialogTitle.value = '添加角色'
  roleForm.value = {
    id: 0,
    name: '',
    code: '',
    desc: '',
    status: true
  }
  dialogVisible.value = true
}

const handleEditRole = (role) => {
  if (!permission.check(PermissionCode.ROLE_UPDATE)) {
    return
  }
  dialogTitle.value = '编辑角色'
  // 将数字状态转换为布尔值，1表示启用（true），0表示禁用（false）
  const roleWithBooleanStatus = {
    ...role,
    status: role.status === 1
  }
  roleForm.value = roleWithBooleanStatus
  dialogVisible.value = true
}

const handleDeleteRole = (id) => {
  if (!permission.check(PermissionCode.ROLE_DELETE)) {
    return
  }
  console.log('准备删除角色，ID:', id, '类型:', typeof id)
  ElMessageBox.confirm('确定要删除这个角色吗？', '警告', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(() => {
    console.log('用户确认删除，发送删除请求')
    // 调用 API 删除角色
    roleApi.deleteRole(id)
    .then(response => {
      console.log('删除角色响应:', response)
      if (response.code === 0) {
        ElMessage.success('删除成功')
        loadRoles()
      } else {
        ElMessage.error('删除失败：' + response.message)
      }
    })
    .catch(error => {
      console.error('删除角色失败:', error)
      console.error('错误详情:', error.response)
      ElMessage.error('删除角色失败')
    })
  })
}

const handleAssignPermissions = async (role) => {
  if (!permission.check(PermissionCode.ROLE_ASSIGN_PERMISSIONS)) {
    return
  }
  
  // 从后端重新获取角色的最新权限数据，避免缓存问题
  try {
    const response = await roleApi.getRole(role.id)
    if (response.code === 0) {
      role = response.data
    }
  } catch (error) {
    console.error('获取角色信息失败:', error)
  }
  
  // 只获取按钮权限ID，不包含父菜单
  // 给按钮权限ID添加前缀 'p_' 与权限树中的节点ID一致
  const buttonPermissions = role.permissions ? role.permissions.filter(perm => perm.type === 2) : []
  const buttonPermissionIds = buttonPermissions.map(perm => 'p_' + perm.id)
  
  assignPermissionsForm.value = {
    roleId: role.id,
    name: role.name,
    permissionIds: buttonPermissionIds
  }
  
  // 确保权限树数据已加载
  if (permissions.value.length === 0) {
    await loadPermissionsAsync()
  }
  
  // 更新权限树的key，强制重建组件，清除之前的勾选状态
  permissionTreeKey.value++
  
  // 打开对话框
  assignPermissionsVisible.value = true
}

const handleAssignMenus = async (role) => {
  if (!permission.check(PermissionCode.ROLE_ASSIGN_MENUS)) {
    return
  }
  
  // 确保菜单列表已加载
  if (menus.value.length === 0) {
    await new Promise((resolve, reject) => {
      menuApi.getMenuTree()
      .then(response => {
        if (response.code === 0) {
          // 后端已经返回了树形结构，直接使用
          const menuData = response.data.menus || response.data.Menus || []
          menus.value = menuData
          resolve()
        } else {
          ElMessage.error('获取菜单列表失败：' + response.message)
          reject()
        }
      })
      .catch(error => {
        console.error('获取菜单列表失败:', error)
        ElMessage.error('获取菜单列表失败')
        reject()
      })
    })
  }
  
  console.log('角色信息:', role)
  console.log('角色的菜单信息:', role.menus || role.Menus)
  assignMenusForm.value = {
    roleId: role.id,
    name: role.name,
    menuIds: role.menus ? role.menus.map(menu => menu.id) : (role.Menus ? role.Menus.map(menu => menu.id) : [])
  }
  console.log('menuIds:', assignMenusForm.value.menuIds)
  assignMenusVisible.value = true
}

const handleSubmit = () => {
  roleFormRef.value.validate((valid) => {
    if (valid) {
      // 将布尔值转换为数字，true表示启用（1），false表示禁用（0）
      const formData = {
        ...roleForm.value,
        status: roleForm.value.status ? 1 : 0
      }
      // 调用 API 创建或更新角色
      if (roleForm.value.id) {
        // 更新角色 - 使用 POST 方法
        roleApi.updateRole(formData)
        .then(response => {
          if (response.code === 0) {
            ElMessage.success('更新成功')
            dialogVisible.value = false
            loadRoles()
          } else {
            ElMessage.error('操作失败：' + response.message)
          }
        })
        .catch(error => {
          console.error('操作失败:', error)
          ElMessage.error('操作失败')
        })
      } else {
        // 创建角色 - 使用 POST 方法
        roleApi.createRole(formData)
        .then(response => {
          if (response.code === 0) {
            ElMessage.success('创建成功')
            dialogVisible.value = false
            loadRoles()
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

const handlePermissionCheck = (data, checkInfo) => {
	// 处理权限树选择变化
	// checkInfo.checkedKeys 包含所有被选中的节点的 key（包括父子联动后的结果）
	if (permissionTreeRef.value) {
		assignPermissionsForm.value.permissionIds = [...checkInfo.checkedKeys]
	}
}

const handleAssignPermissionsSubmit = () => {
  // 获取所有被选中的节点（包括父子联动的）
  const allCheckedIds = permissionTreeRef.value ? permissionTreeRef.value.getCheckedKeys() : []
  
  // 只保留按钮权限的ID（以 'p_' 前缀开头），并移除前缀
  const buttonPermissionIds = []
  const findButtonPermissions = (nodes) => {
    nodes.forEach(node => {
      if (node.type === 2 && allCheckedIds.includes(node.id)) {
        // 移除 'p_' 前缀，恢复为原始ID
        const originalId = String(node.id).replace('p_', '')
        buttonPermissionIds.push(Number(originalId))
      }
      if (node.children && node.children.length > 0) {
        findButtonPermissions(node.children)
      }
    })
  }
  findButtonPermissions(permissions.value)
  
  roleApi.assignPermissions(assignPermissionsForm.value.roleId, buttonPermissionIds)
  .then(response => {
    if (response.code === 0) {
      ElMessage.success('分配权限成功')
      assignPermissionsVisible.value = false
      loadRoles()
    } else {
      ElMessage.error('分配权限失败：' + response.message)
    }
  })
  .catch(error => {
    console.error('分配权限失败:', error)
    ElMessage.error('分配权限失败')
  })
}

const handleAssignMenusSubmit = () => {
  // 从 el-tree 获取当前选中的菜单ID（包括全选和半选的父节点）
  const checkedMenuIds = menuTreeRef.value ? menuTreeRef.value.getCheckedKeys() : []
  const halfCheckedMenuIds = menuTreeRef.value ? menuTreeRef.value.getHalfCheckedKeys() : []
  
  // 合并全选和半选的菜单ID（确保父级菜单也被写入）
  const allMenuIds = [...new Set([...checkedMenuIds, ...halfCheckedMenuIds])]
  
  // 调用 API 分配菜单
  roleApi.assignMenus(assignMenusForm.value.roleId, allMenuIds)
  .then(response => {
    if (response.code === 0) {
      ElMessage.success('分配菜单成功')
      assignMenusVisible.value = false
      loadRoles()
    } else {
      ElMessage.error('分配菜单失败：' + response.message)
    }
  })
  .catch(error => {
    console.error('分配菜单失败:', error)
    ElMessage.error('分配菜单失败')
  })
}

const handleSearch = () => {
  // 模拟搜索
  loadRoles()
}

const handleSizeChange = (size) => {
  pageSize.value = size
  loadRoles()
}

const handleCurrentChange = (current) => {
  currentPage.value = current
  loadRoles()
}

onMounted(() => {
  loadRoles()
})
</script>

<style scoped>
.role-management {
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
  display: flex;
  justify-content: flex-end;
}

.tree-node {
  display: flex;
  align-items: center;
}

.code {
  margin-left: 10px;
  font-size: 12px;
  color: #909399;
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
  
  .el-tree {
    font-size: 12px;
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
