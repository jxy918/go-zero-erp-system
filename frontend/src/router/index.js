import { createRouter, createWebHistory } from 'vue-router'
import Login from '../views/login/Login.vue'
import Layout from '../components/Layout.vue'
import { useUserStore, useMenuStore } from '../store'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/login',
      name: 'Login',
      component: Login
    },
    {
      path: '/',
      name: 'Layout',
      component: Layout,
      children: [
        {
          path: '',
          name: 'Dashboard',
          component: () => import('../views/dashboard/Dashboard.vue')
        },
        {
          path: 'user',
          name: 'User',
          component: () => import('../views/user/User.vue')
        },
        {
          path: 'role',
          name: 'Role',
          component: () => import('../views/role/Role.vue')
        },
        {
          path: 'permission',
          name: 'Permission',
          component: () => import('../views/permission/Permission.vue')
        },
        {
          path: 'menu',
          name: 'Menu',
          component: () => import('../views/menu/Menu.vue')
        },
        {
          path: 'activity',
          name: 'Activity',
          component: () => import('../views/activity/Activity.vue')
        },
        {
          path: 'product',
          name: 'Product',
          component: () => import('../views/product/Product.vue')
        },
        {
          path: 'product/category',
          name: 'Category',
          component: () => import('../views/product/Category.vue')
        },
        {
          path: 'supplier',
          name: 'Supplier',
          component: () => import('../views/supplier/Supplier.vue')
        },
        {
          path: 'customer',
          name: 'Customer',
          component: () => import('../views/customer/Customer.vue')
        },
        {
          path: 'warehouse',
          name: 'Warehouse',
          component: () => import('../views/warehouse/Warehouse.vue')
        },
        {
          path: 'purchase',
          name: 'Purchase',
          component: () => import('../views/purchase/Purchase.vue')
        },
        {
          path: 'sales',
          name: 'Sales',
          component: () => import('../views/sales/Sales.vue')
        },
        {
          path: 'inventory',
          name: 'Inventory',
          component: () => import('../views/inventory/Inventory.vue')
        },
        {
          path: 'inventory/adjust-request',
          name: 'InventoryAdjustRequest',
          component: () => import('../views/inventory/InventoryAdjustRequest.vue')
        },
        {
          path: 'inventory/check',
          name: 'InventoryCheck',
          component: () => import('../views/inventory/InventoryCheck.vue')
        },
        {
          path: 'inventory/check/:id',
          name: 'InventoryCheckDetail',
          component: () => import('../views/inventory/InventoryCheckDetail.vue')
        },
        {
          path: 'inventory/transfer',
          name: 'InventoryTransfer',
          component: () => import('../views/inventory/InventoryTransfer.vue')
        },
        {
          path: 'inventory/alert',
          name: 'InventoryAlert',
          component: () => import('../views/inventory/InventoryAlert.vue')
        },
        {
          path: 'product/unit',
          name: 'ProductUnit',
          component: () => import('../views/product/ProductUnit.vue')
        },
        {
          path: 'erp',
          name: 'Erp',
          component: () => import('../views/erp/Erp.vue')
        }
      ]
    }
  ]
})

// 路由守卫
router.beforeEach(async (to, from, next) => {
  const userStore = useUserStore()
  const menuStore = useMenuStore()
  
  // 从localStorage加载用户信息
  userStore.loadUserInfo()
  
  if (to.path === '/login') {
    next()
  } else {
    if (userStore.getIsAuthenticated) {
      const userInfo = userStore.getUserInfo
      if (userInfo) {
        // 检查用户状态
        if (userInfo.status === 0) {
          userStore.logout()
          next('/login')
          return
        }
        
        // 加载菜单数据
        if (menuStore.getMenuTree.length === 0) {
          await menuStore.loadMenuTree()
        }
        
        // 检查用户权限
        const hasPermission = checkPermission(to.path, userInfo, menuStore.getMenuTree)
        if (hasPermission) {
          next()
        } else {
          next('/')
        }
      } else {
        next('/login')
      }
    } else {
      next('/login')
    }
  }
})

// 检查权限
function checkPermission(path, user, menuTree) {
  // 控制台页面所有人都可以访问
  if (path === '/') {
    return true
  }
  
  // 检查用户是否为管理员（这里假设管理员角色的code为'admin'）
  const isAdmin = user.roles && user.roles.some(role => role.code === 'admin')
  if (isAdmin) {
    return true
  }
  
  // 检查用户是否有对应路径的菜单权限
  if (menuTree.length > 0) {
    // 遍历菜单树，检查是否有匹配的路径
    const hasMenuPermission = checkMenuPermission(path, menuTree)
    if (hasMenuPermission) {
      return true
    }
  }
  
  // 检查用户是否有对应路径的权限
  if (user.roles) {
    for (const role of user.roles) {
      if (role.permissions) {
        for (const permission of role.permissions) {
          // 检查权限路径是否匹配
          if (permission.path === path) {
            return true
          }
        }
      }
    }
  }
  
  return false
}

// 检查菜单权限
function checkMenuPermission(path, menuTree) {
  for (const menu of menuTree) {
    if (menu.path === path) {
      return true
    }
    if (menu.children && menu.children.length > 0) {
      const hasChildPermission = checkMenuPermission(path, menu.children)
      if (hasChildPermission) {
        return true
      }
    }
  }
  return false
}

export default router