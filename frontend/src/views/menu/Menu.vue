<template>
  <div class="menu-management">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>菜单管理</span>
          <el-button type="primary" @click="handleAddMenu" v-has-permission="'btn_menu_create'">添加菜单</el-button>
        </div>
      </template>
      <div class="description-box">
        <span class="description-icon">ℹ️</span>
        <span class="description-text">菜单管理：管理系统导航菜单结构，支持多级菜单配置，控制前端页面的展示和路由访问。</span>
      </div>
      <el-tree
        :data="menuTree"
        node-key="id"
        default-expand-all
      >
        <template #default="{ node, data }">
          <span class="tree-node">
            <span class="tree-node-content">
              <span class="icon menu-icon">
                <el-icon><Folder /></el-icon>
              </span>
              <span>{{ data.name }}</span>
              <span v-if="data.code" class="code">({{ data.code }})</span>
              <span :class="['status-tag', data.status === 1 ? 'status-enabled' : 'status-disabled']">
                {{ data.status === 1 ? '启用' : '禁用' }}
              </span>
            </span>
            <span class="tree-actions">
              <el-button v-if="hasMenuUpdate" type="primary" size="small" @click="handleEditMenu(data)">
                <template #icon><Edit /></template>
                <span>编辑</span>
              </el-button>
              <el-button v-if="hasMenuDelete" type="danger" size="small" @click="handleDeleteMenu(data.id)">
                <template #icon><Delete /></template>
                <span>删除</span>
              </el-button>
            </span>
          </span>
        </template>
      </el-tree>
    </el-card>

    <!-- 添加/编辑菜单对话框 -->
    <el-dialog
      :title="dialogTitle"
      v-model="dialogVisible"
      width="500px"
    >
      <el-form :model="menuForm" :rules="rules" ref="menuFormRef">
        <el-form-item label="菜单名称" prop="name">
          <el-input v-model="menuForm.name" placeholder="请输入菜单名称"></el-input>
        </el-form-item>
        <el-form-item label="菜单编码" prop="code">
          <el-input v-model="menuForm.code" placeholder="请输入菜单编码"></el-input>
        </el-form-item>
        <el-form-item label="菜单描述" prop="desc">
          <el-input
            v-model="menuForm.desc"
            type="textarea"
            placeholder="请输入菜单描述"
          ></el-input>
        </el-form-item>
        <el-form-item label="父菜单" prop="parent_id">
          <el-select v-model="menuForm.parent_id" placeholder="请选择父菜单">
            <el-option label="根菜单" :value="0"></el-option>
            <el-option
              v-for="menu in menuOptions"
              :key="menu.id"
              :label="menu.name"
              :value="menu.id"
            ></el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="菜单路径" prop="path">
          <el-input v-model="menuForm.path" placeholder="请输入菜单路径"></el-input>
        </el-form-item>
        <el-form-item label="组件路径" prop="component">
          <el-input v-model="menuForm.component" placeholder="请输入组件路径"></el-input>
        </el-form-item>
        <el-form-item label="菜单图标" prop="icon">
          <el-input v-model="menuForm.icon" placeholder="请输入菜单图标"></el-input>
        </el-form-item>
        <el-form-item label="排序" prop="sort">
          <el-input-number v-model="menuForm.sort" :min="0"></el-input-number>
        </el-form-item>
        <el-form-item label="状态">
          <el-switch v-model="menuForm.status"></el-switch>
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
import { useMenuStore } from '../../store'
import { useUserStore } from '../../store'
import { permission, PermissionCode } from '../../utils/permission'

const menuStore = useMenuStore()
const userStore = useUserStore()

// 直接从store获取权限列表，而不是通过permission.has()
const hasMenuUpdate = computed(() => {
  if (userStore.isAdmin) return true
  const codes = userStore.permissions.map(p => p.code || p.Code)
  return codes.includes(PermissionCode.MENU_UPDATE)
})

const hasMenuDelete = computed(() => {
  if (userStore.isAdmin) return true
  const codes = userStore.permissions.map(p => p.code || p.Code)
  return codes.includes(PermissionCode.MENU_DELETE)
})

const hasAnyMenuPermission = computed(() => {
  return hasMenuUpdate.value || hasMenuDelete.value
})

const dialogVisible = ref(false)
const dialogTitle = ref('添加菜单')
const menuForm = ref({
  id: 0,
  name: '',
  code: '',
  desc: '',
  parent_id: 0,
  path: '',
  component: '',
  icon: '',
  sort: 0,
  status: true
})
const rules = {
  name: [{ required: true, message: '请输入菜单名称', trigger: 'blur' }],
  code: [{ required: true, message: '请输入菜单编码', trigger: 'blur' }],
  path: [{ required: true, message: '请输入菜单路径', trigger: 'blur' }],
  component: [{ required: true, message: '请输入组件路径', trigger: 'blur' }]
}
const menuFormRef = ref(null)

const menuTree = computed(() => {
  return menuStore.getMenuTree
})

const menuOptions = computed(() => {
  return menuStore.getMenuList
})

const handleAddMenu = () => {
  if (!permission.check(PermissionCode.MENU_CREATE)) {
    return
  }
  dialogTitle.value = '添加菜单'
  menuForm.value = {
    id: 0,
    name: '',
    code: '',
    desc: '',
    parent_id: 0,
    path: '',
    component: '',
    icon: '',
    sort: 0,
    status: true
  }
  dialogVisible.value = true
}

const handleEditMenu = (menu) => {
  if (!permission.check(PermissionCode.MENU_UPDATE)) {
    return
  }
  dialogTitle.value = '编辑菜单'
  // 将数字状态转换为布尔值，1表示启用（true），0表示禁用（false）
  const menuWithBooleanStatus = {
    ...menu,
    status: menu.status === 1
  }
  menuForm.value = menuWithBooleanStatus
  dialogVisible.value = true
}

const handleDeleteMenu = (id) => {
  if (!permission.check(PermissionCode.MENU_DELETE)) {
    return
  }
  ElMessageBox.confirm('确定要删除这个菜单吗？', '警告', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async () => {
    const success = await menuStore.deleteMenu(id)
    if (success) {
      ElMessage.success('删除成功')
    }
  })
}

const handleSubmit = async () => {
  menuFormRef.value.validate(async (valid) => {
    if (valid) {
      // 准备表单数据，true表示启用（1），false表示禁用（0）
      const formData = {
        name: menuForm.value.name,
        code: menuForm.value.code,
        desc: menuForm.value.desc,
        parent_id: parseInt(menuForm.value.parent_id) || 0,
        path: menuForm.value.path,
        component: menuForm.value.component,
        icon: menuForm.value.icon,
        sort: parseInt(menuForm.value.sort) || 0,
        status: menuForm.value.status ? 1 : 0
      }
      
      console.log('提交的表单数据:', formData)
      
      let success
      if (menuForm.value.id) {
        formData.id = menuForm.value.id
        // 更新菜单
        success = await menuStore.updateMenu(formData)
      } else {
        // 创建菜单
        success = await menuStore.createMenu(formData)
      }
      
      if (success) {
        ElMessage.success(menuForm.value.id ? '更新成功' : '创建成功')
        dialogVisible.value = false
      }
    }
  })
}

onMounted(async () => {
  await menuStore.loadMenuTree()
})
</script>

<style scoped>
.menu-management {
  padding: 0;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
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

.tree-node-content .menu-icon {
  background-color: #e6f7ff;
}

.tree-node-content .menu-icon :deep(.el-icon) {
  color: #1890ff;
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

/* 响应式设计 */
@media screen and (max-width: 768px) {
  .card-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 10px;
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
  
  /* 响应式表单样式 */
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