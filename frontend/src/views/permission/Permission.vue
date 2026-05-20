<template>
  <div class="permission-management">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>权限管理</span>
          <el-button v-if="hasPermissionCreate" type="primary" @click="handleAddPermission">添加权限</el-button>
        </div>
      </template>
      <div class="description-box">
        <span class="description-icon">ℹ️</span>
        <span class="description-text">权限管理：管理系统权限点，包括菜单权限和按钮权限，通过角色分配实现细粒度的访问控制。</span>
      </div>
      <div class="search-container">
        <el-input
          v-model="searchQuery"
          placeholder="请输入权限或菜单名称"
          style="width: 300px"
          prefix-icon="el-icon-search"
        ></el-input>
        <el-button type="primary" @click="handleSearch">搜索</el-button>
      </div>
      
      <el-tree
        :data="treeData"
        :props="treeProps"
        node-key="id"
        :default-expand-all="true"
        class="permission-tree"
      >
        <template #default="{ node, data }">
          <div class="tree-node-wrapper">
            <span class="tree-node">
              <span :class="['icon-wrapper', data.type === 'menu' ? 'menu-icon-wrapper' : 'permission-icon-wrapper']">
                <el-icon class="icon" :class="data.type === 'menu' ? 'menu-icon' : 'permission-icon'">
                  <Folder v-if="data.type === 'menu'" />
                  <Key v-else />
                </el-icon>
                <span :class="['icon-label', data.type === 'menu' ? 'menu-label' : 'permission-label']">
                  {{ data.type === 'menu' ? '菜单' : '权限' }}
                </span>
              </span>
              <span class="node-name">{{ data.name }}</span>
              <span v-if="data.code" class="code">({{ data.code }})</span>
              <span :class="['status-tag', data.status === 1 ? 'status-enabled' : 'status-disabled']">
                {{ data.status === 1 ? '启用' : '禁用' }}
              </span>
            </span>
            <div class="tree-actions" v-if="data.type === 'permission'">
              <el-button 
                v-if="hasPermissionUpdate" 
                type="primary" 
                size="small"
                @click.stop="handleEditPermission(data)">
                <template #icon><Edit /></template>
                <span>编辑</span>
              </el-button>
              <el-button 
                v-if="hasPermissionDelete" 
                type="danger" 
                size="small"
                @click.stop="handleDeletePermission(data.id)">
                <template #icon><Delete /></template>
                <span>删除</span>
              </el-button>
            </div>
          </div>
        </template>
      </el-tree>
      
      <div class="pagination-container">
        <span class="total-count">共 {{ total }} 条权限</span>
      </div>
    </el-card>

    <el-dialog
      :title="dialogTitle"
      v-model="dialogVisible"
      width="500px"
    >
      <el-form :model="permissionForm" :rules="rules" ref="permissionFormRef">
        <el-form-item label="权限名称" prop="name">
          <el-input v-model="permissionForm.name" placeholder="请输入权限名称"></el-input>
        </el-form-item>
        <el-form-item label="权限编码" prop="code">
          <el-input v-model="permissionForm.code" placeholder="请输入权限编码"></el-input>
        </el-form-item>
        <el-form-item label="权限描述" prop="desc">
          <el-input
            v-model="permissionForm.desc"
            type="textarea"
            placeholder="请输入权限描述"
          ></el-input>
        </el-form-item>
        <el-form-item label="API路径" prop="path">
          <el-input v-model="permissionForm.path" placeholder="请输入API路径，如 /user/list"></el-input>
        </el-form-item>
        <el-form-item label="所属菜单" prop="menu_id">
          <el-tree-select
            v-model="permissionForm.menu_id"
            :data="menuTreeData"
            :props="treeProps"
            node-key="id"
            placeholder="请选择所属菜单"
            check-strictly
            style="width: 100%"
          />
        </el-form-item>
        <el-form-item label="排序" prop="sort">
          <el-input-number v-model="permissionForm.sort" :min="0"></el-input-number>
        </el-form-item>
        <el-form-item label="状态">
          <el-switch v-model="permissionForm.status"></el-switch>
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
import { ref, computed, onMounted, reactive } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Folder, Key, Edit, Delete } from '@element-plus/icons-vue'
import { permissionApi } from '../../api/permission'
import { permission, PermissionCode } from '../../utils/permission'
import { useUserStore } from '../../store'

const userStore = useUserStore()

const hasPermissionCreate = computed(() => {
  const codes = userStore.permissions.map(p => p.code || p.Code)
  return userStore.isAdmin || codes.includes(PermissionCode.PERMISSION_CREATE)
})

const hasPermissionUpdate = computed(() => {
  const codes = userStore.permissions.map(p => p.code || p.Code)
  return userStore.isAdmin || codes.includes(PermissionCode.PERMISSION_UPDATE)
})

const hasPermissionDelete = computed(() => {
  const codes = userStore.permissions.map(p => p.code || p.Code)
  return userStore.isAdmin || codes.includes(PermissionCode.PERMISSION_DELETE)
})

const permissions = ref([])
const total = ref(0)
const searchQuery = ref('')

const dialogVisible = ref(false)
const dialogTitle = ref('添加权限')
const permissionForm = reactive({
  id: 0,
  name: '',
  code: '',
  desc: '',
  path: '',
  menu_id: 0,
  sort: 0,
  status: 1
})
const rules = {
  name: [{ required: true, message: '请输入权限名称', trigger: 'blur' }],
  code: [{ required: true, message: '请输入权限编码', trigger: 'blur' }],
  menu_id: [{ required: true, message: '请选择所属菜单', trigger: 'change' }]
}
const permissionFormRef = ref(null)

const treeProps = {
  children: 'children',
  label: 'name'
}

const menuTreeData = computed(() => {
  const result = []
  const menuMap = {}
  
  permissions.value.forEach(perm => {
    if (perm.type === 1) {
      menuMap[perm.id] = {
        id: perm.id,
        name: perm.name,
        code: perm.code,
        parent_id: perm.ParentID || perm.parent_id || 0,
        type: 1,
        children: []
      }
    }
  })
  
  Object.values(menuMap).forEach(menuNode => {
    const parentId = menuNode.parent_id
    if (parentId > 0 && menuMap[parentId]) {
      menuMap[parentId].children.push(menuNode)
    } else {
      result.push(menuNode)
    }
  })
  
  return result
})

const treeData = computed(() => {
  const result = []
  const menuMap = {}
  
  permissions.value.forEach(perm => {
    if (perm.type === 1) {
      menuMap[perm.id] = {
        id: perm.id,
        name: perm.name,
        code: perm.code,
        type: 'menu',
        parent_id: perm.parent_id || 0,
        sort: perm.sort,
        status: perm.status,
        children: []
      }
    }
  })
  
  permissions.value.forEach(perm => {
    if (perm.type === 2) {
      const parentId = perm.parent_id || 0
      const permNode = {
        id: perm.id,
        name: perm.name,
        code: perm.code,
        permissionType: '按钮权限',
        parent_id: parentId,
        sort: perm.sort,
        status: perm.status,
        children: [],
        ...perm,
        type: 'permission'
      }
      
      if (menuMap[parentId]) {
        menuMap[parentId].children.push(permNode)
      } else {
        result.push(permNode)
      }
    }
  })
  
  Object.values(menuMap).forEach(menuNode => {
    const parentId = menuNode.parent_id
    if (parentId > 0 && menuMap[parentId]) {
      menuMap[parentId].children.push(menuNode)
    } else {
      result.push(menuNode)
    }
  })
  
  result.sort((a, b) => {
    return (a.sort || 0) - (b.sort || 0)
  })
  
  const sortChildren = (nodes) => {
    nodes.forEach(node => {
      if (node.children && node.children.length > 0) {
        node.children.sort((a, b) => {
          return (a.sort || 0) - (b.sort || 0)
        })
        sortChildren(node.children)
      }
    })
  }
  sortChildren(result)
  
  total.value = permissions.value.length
  
  return result
})

const loadPermissions = () => {
  permissionApi.getPermissionList({
    page: 1,
    page_size: 500,
    name: searchQuery.value
  })
  .then(response => {
    if (response.code === 0) {
      permissions.value = response.data.permissions
    } else {
      ElMessage.error('获取权限列表失败：' + response.message)
    }
  })
  .catch(error => {
    console.error('获取权限列表失败:', error)
    ElMessage.error('获取权限列表失败')
  })
}

const handleAddPermission = () => {
  if (!permission.check(PermissionCode.PERMISSION_CREATE)) {
    return
  }
  dialogTitle.value = '添加权限'
  permissionForm.id = 0
  permissionForm.name = ''
  permissionForm.code = ''
  permissionForm.desc = ''
  permissionForm.path = ''
  permissionForm.menu_id = null
  permissionForm.sort = 0
  permissionForm.status = 1
  dialogVisible.value = true
}

const handleEditPermission = (data) => {
  if (!permission.check(PermissionCode.PERMISSION_UPDATE)) {
    return
  }
  dialogTitle.value = '编辑权限'
  permissionForm.id = data.id || 0
  permissionForm.name = data.name || ''
  permissionForm.code = data.code || ''
  permissionForm.desc = data.desc || ''
  permissionForm.path = data.path || ''
  const menuId = data.menu_id || data.MenuID
  permissionForm.menu_id = menuId && menuId > 0 ? menuId : null
  permissionForm.sort = data.sort || 0
  permissionForm.status = data.status === 1 || data.status === true
  dialogVisible.value = true
}

const handleDeletePermission = (id) => {
  if (!permission.check(PermissionCode.PERMISSION_DELETE)) {
    return
  }
  ElMessageBox.confirm('确定要删除这个权限吗？', '警告', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(() => {
    permissionApi.deletePermission(id)
    .then(response => {
      if (response.code === 0) {
        ElMessage.success('删除成功')
        loadPermissions()
      } else {
        ElMessage.error('删除失败：' + response.message)
      }
    })
    .catch(error => {
      console.error('删除权限失败:', error)
      ElMessage.error('删除权限失败')
    })
  })
}

const handleSubmit = () => {
  permissionFormRef.value.validate((valid) => {
    if (valid) {
      const formData = permissionForm.id ? {
        id: permissionForm.id,
        name: permissionForm.name,
        code: permissionForm.code,
        desc: permissionForm.desc,
        path: permissionForm.path || '',
        sort: parseInt(permissionForm.sort) || 0,
        status: permissionForm.status ? 1 : 0,
        menu_id: parseInt(permissionForm.menu_id) || 0
      } : {
        name: permissionForm.name,
        code: permissionForm.code,
        desc: permissionForm.desc,
        path: permissionForm.path || '',
        sort: parseInt(permissionForm.sort) || 0,
        status: permissionForm.status ? 1 : 0,
        menu_id: parseInt(permissionForm.menu_id) || 0
      }
      if (permissionForm.id) {
        permissionApi.updatePermission(formData)
        .then(response => {
          if (response.code === 0) {
            ElMessage.success('更新成功')
            dialogVisible.value = false
            loadPermissions()
          } else {
            ElMessage.error('操作失败：' + response.message)
          }
        })
        .catch(error => {
          console.error('操作失败:', error)
          ElMessage.error('操作失败')
        })
      } else {
        permissionApi.createPermission(formData)
        .then(response => {
          if (response.code === 0) {
            ElMessage.success('创建成功')
            dialogVisible.value = false
            loadPermissions()
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

const handleSearch = () => {
  loadPermissions()
}

onMounted(() => {
  loadPermissions()
})
</script>

<style scoped>
.permission-management {
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

.permission-tree {
  margin-top: 20px;
}

.tree-node-wrapper {
  display: flex;
  align-items: center;
  justify-content: space-between;
  width: 100%;
}

.tree-node {
  display: flex;
  align-items: center;
  gap: 12px;
  flex: 1;
}

.icon-wrapper {
  display: flex;
  align-items: center;
  gap: 4px;
  padding: 4px 8px;
  border-radius: 6px;
  min-width: 60px;
  justify-content: center;
}

.menu-icon-wrapper {
  background: linear-gradient(135deg, #e6f7ff 0%, #bae7ff 100%);
}

.permission-icon-wrapper {
  background: linear-gradient(135deg, #f6ffed 0%, #b7eb8f 100%);
}

.icon {
  font-size: 14px;
}

.menu-icon {
  color: #1890ff;
}

.permission-icon {
  color: #52c41a;
}

.icon-label {
  font-size: 11px;
  font-weight: 500;
}

.menu-label {
  color: #1890ff;
}

.permission-label {
  color: #52c41a;
}

.node-name {
  font-weight: 500;
  color: #303133;
}

.tree-node .code {
  font-size: 12px;
  color: #909399;
  font-family: monospace;
}

.tree-node .status-tag {
  font-size: 11px;
  padding: 2px 8px;
  border-radius: 10px;
  margin-left: 8px;
}

.tree-node .status-enabled {
  background-color: #f6ffed;
  color: #52c41a;
  border: 1px solid #b7eb8f;
}

.tree-node .status-disabled {
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

.pagination-container {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
}

.total-count {
  font-size: 14px;
  color: #666;
}

:deep(.el-tree-node__content) {
  display: flex;
  align-items: center;
  width: 100%;
}

:deep(.el-tree-node__children) {
  padding-left: 24px;
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
    flex-wrap: wrap;
  }
  
  .tree-actions {
    margin-left: 0;
    margin-top: 8px;
    width: 100%;
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
