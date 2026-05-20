# RBAC 后台管理系统 - 数据字典

## 1. 概述

本文档描述 RBAC 后台管理系统的数据库表结构和字段定义，包含所有数据表的详细说明。

## 2. 数据库表清单

| 序号 | 表名 | 说明 | 状态 |
| :---: | :--- | :--- | :---: |
| 1 | `users` | 用户表 | ✅ |
| 2 | `roles` | 角色表 | ✅ |
| 3 | `permissions` | 权限表 | ✅ |
| 4 | `menus` | 菜单表 | ✅ |
| 5 | `user_roles` | 用户角色关联表 | ✅ |
| 6 | `role_permissions` | 角色权限关联表 | ✅ |
| 7 | `role_menus` | 角色菜单关联表 | ✅ |
| 8 | `menu_permissions` | 菜单权限关联表 | ✅ |
| 9 | `activities` | 活动日志表 | ✅ |

## 3. 表结构详情

### 3.1 用户表 (users)

| 字段名 | 类型 | 约束 | 说明 | 默认值 |
| :--- | :--- | :--- | :--- | :--- |
| `id` | bigint unsigned | PRIMARY KEY, AUTO_INCREMENT | 用户ID | - |
| `username` | varchar(50) | NOT NULL, UNIQUE | 用户名（登录账号） | - |
| `password` | varchar(100) | NOT NULL | 密码（BCrypt加密） | - |
| `nickname` | varchar(50) | NULL | 昵称（显示名称） | '' |
| `email` | varchar(100) | NULL | 邮箱地址 | '' |
| `phone` | varchar(20) | NULL | 手机号码 | '' |
| `status` | tinyint | NULL | 状态：1=启用，0=禁用 | 1 |
| `created_at` | datetime | NULL | 创建时间 | CURRENT_TIMESTAMP |
| `updated_at` | datetime | NULL | 更新时间 | CURRENT_TIMESTAMP ON UPDATE |
| `deleted_at` | datetime | NULL, INDEX | 删除时间（软删除） | NULL |

**索引**:
- PRIMARY KEY: `id`
- UNIQUE KEY: `uk_username` (`username`)
- INDEX: `idx_deleted_at` (`deleted_at`)

### 3.2 角色表 (roles)

| 字段名 | 类型 | 约束 | 说明 | 默认值 |
| :--- | :--- | :--- | :--- | :--- |
| `id` | bigint unsigned | PRIMARY KEY, AUTO_INCREMENT | 角色ID | - |
| `name` | varchar(50) | NOT NULL, UNIQUE | 角色名称 | - |
| `code` | varchar(50) | NOT NULL, UNIQUE | 角色编码（唯一标识） | - |
| `desc` | varchar(200) | NULL | 角色描述 | '' |
| `status` | tinyint | NULL, INDEX | 状态：1=启用，0=禁用 | 1 |
| `created_at` | datetime | NULL | 创建时间 | CURRENT_TIMESTAMP |
| `updated_at` | datetime | NULL | 更新时间 | CURRENT_TIMESTAMP ON UPDATE |
| `deleted_at` | datetime | NULL, INDEX | 删除时间（软删除） | NULL |

**索引**:
- PRIMARY KEY: `id`
- UNIQUE KEY: `uk_name` (`name`)
- UNIQUE KEY: `uk_code` (`code`)
- INDEX: `idx_status` (`status`)
- INDEX: `idx_deleted_at` (`deleted_at`)

### 3.3 权限表 (permissions)

| 字段名 | 类型 | 约束 | 说明 | 默认值 |
| :--- | :--- | :--- | :--- | :--- |
| `id` | bigint unsigned | PRIMARY KEY, AUTO_INCREMENT | 权限ID | - |
| `name` | varchar(50) | NOT NULL, UNIQUE | 权限名称 | - |
| `code` | varchar(50) | NOT NULL, UNIQUE | 权限编码（唯一标识） | - |
| `desc` | varchar(200) | NULL | 权限描述 | '' |
| `type` | tinyint | NULL, INDEX | 类型：1=菜单权限，2=按钮权限 | 1 |
| `parent_id` | bigint unsigned | NULL, INDEX | 父权限ID | 0 |
| `path` | varchar(200) | NULL | 路由路径 | '' |
| `component` | varchar(200) | NULL | 组件路径 | '' |
| `icon` | varchar(50) | NULL | 图标 | '' |
| `sort` | int | NULL | 排序号 | 0 |
| `status` | tinyint | NULL, INDEX | 状态：1=启用，0=禁用 | 1 |
| `created_at` | datetime | NULL | 创建时间 | CURRENT_TIMESTAMP |
| `updated_at` | datetime | NULL | 更新时间 | CURRENT_TIMESTAMP ON UPDATE |
| `deleted_at` | datetime | NULL, INDEX | 删除时间（软删除） | NULL |

**索引**:
- PRIMARY KEY: `id`
- UNIQUE KEY: `uk_name` (`name`)
- UNIQUE KEY: `uk_code` (`code`)
- INDEX: `idx_type` (`type`)
- INDEX: `idx_parent_id` (`parent_id`)
- INDEX: `idx_status` (`status`)
- INDEX: `idx_deleted_at` (`deleted_at`)

### 3.4 菜单表 (menus)

| 字段名 | 类型 | 约束 | 说明 | 默认值 |
| :--- | :--- | :--- | :--- | :--- |
| `id` | bigint unsigned | PRIMARY KEY, AUTO_INCREMENT | 菜单ID | - |
| `name` | varchar(50) | NOT NULL, UNIQUE | 菜单名称 | - |
| `code` | varchar(50) | NOT NULL, UNIQUE | 菜单编码（唯一标识） | - |
| `desc` | varchar(200) | NULL | 菜单描述 | '' |
| `parent_id` | bigint unsigned | NULL, INDEX | 父菜单ID（0表示顶级菜单） | 0 |
| `path` | varchar(200) | NOT NULL | 路由路径 | - |
| `component` | varchar(200) | NOT NULL | 组件路径 | - |
| `icon` | varchar(50) | NULL | 菜单图标（Element Plus图标名） | '' |
| `sort` | int | NULL | 排序号（越小越靠前） | 0 |
| `status` | tinyint | NULL, INDEX | 状态：1=启用，0=禁用 | 1 |
| `created_at` | datetime | NULL | 创建时间 | CURRENT_TIMESTAMP |
| `updated_at` | datetime | NULL | 更新时间 | CURRENT_TIMESTAMP ON UPDATE |
| `deleted_at` | datetime | NULL, INDEX | 删除时间（软删除） | NULL |

**索引**:
- PRIMARY KEY: `id`
- UNIQUE KEY: `uk_name` (`name`)
- UNIQUE KEY: `uk_code` (`code`)
- INDEX: `idx_parent_id` (`parent_id`)
- INDEX: `idx_status` (`status`)
- INDEX: `idx_deleted_at` (`deleted_at`)

### 3.5 用户角色关联表 (user_roles)

| 字段名 | 类型 | 约束 | 说明 | 默认值 |
| :--- | :--- | :--- | :--- | :--- |
| `user_id` | bigint unsigned | PRIMARY KEY, NOT NULL | 用户ID | - |
| `role_id` | bigint unsigned | PRIMARY KEY, NOT NULL | 角色ID | - |
| `created_at` | datetime | NULL | 创建时间 | CURRENT_TIMESTAMP |
| `updated_at` | datetime | NULL | 更新时间 | CURRENT_TIMESTAMP ON UPDATE |
| `deleted_at` | datetime | NULL, INDEX | 删除时间（软删除） | NULL |

**索引**:
- PRIMARY KEY: (`user_id`, `role_id`)
- INDEX: `idx_user_id` (`user_id`)
- INDEX: `idx_role_id` (`role_id`)
- INDEX: `idx_deleted_at` (`deleted_at`)

### 3.6 角色权限关联表 (role_permissions)

| 字段名 | 类型 | 约束 | 说明 | 默认值 |
| :--- | :--- | :--- | :--- | :--- |
| `role_id` | bigint unsigned | PRIMARY KEY, NOT NULL | 角色ID | - |
| `permission_id` | bigint unsigned | PRIMARY KEY, NOT NULL | 权限ID | - |
| `created_at` | datetime | NULL | 创建时间 | CURRENT_TIMESTAMP |
| `updated_at` | datetime | NULL | 更新时间 | CURRENT_TIMESTAMP ON UPDATE |
| `deleted_at` | datetime | NULL, INDEX | 删除时间（软删除） | NULL |

**索引**:
- PRIMARY KEY: (`role_id`, `permission_id`)
- INDEX: `idx_role_id` (`role_id`)
- INDEX: `idx_permission_id` (`permission_id`)
- INDEX: `idx_deleted_at` (`deleted_at`)

### 3.7 角色菜单关联表 (role_menus)

| 字段名 | 类型 | 约束 | 说明 | 默认值 |
| :--- | :--- | :--- | :--- | :--- |
| `role_id` | bigint unsigned | PRIMARY KEY, NOT NULL | 角色ID | - |
| `menu_id` | bigint unsigned | PRIMARY KEY, NOT NULL | 菜单ID | - |
| `created_at` | datetime | NULL | 创建时间 | CURRENT_TIMESTAMP |
| `updated_at` | datetime | NULL | 更新时间 | CURRENT_TIMESTAMP ON UPDATE |
| `deleted_at` | datetime | NULL, INDEX | 删除时间（软删除） | NULL |

**索引**:
- PRIMARY KEY: (`role_id`, `menu_id`)
- INDEX: `idx_role_id` (`role_id`)
- INDEX: `idx_menu_id` (`menu_id`)
- INDEX: `idx_deleted_at` (`deleted_at`)

### 3.8 菜单权限关联表 (menu_permissions)

| 字段名 | 类型 | 约束 | 说明 | 默认值 |
| :--- | :--- | :--- | :--- | :--- |
| `menu_id` | bigint unsigned | PRIMARY KEY, NOT NULL | 菜单ID | - |
| `permission_id` | bigint unsigned | PRIMARY KEY, NOT NULL | 权限ID | - |
| `created_at` | datetime | NULL | 创建时间 | CURRENT_TIMESTAMP |
| `updated_at` | datetime | NULL | 更新时间 | CURRENT_TIMESTAMP ON UPDATE |
| `deleted_at` | datetime | NULL, INDEX | 删除时间（软删除） | NULL |

**索引**:
- PRIMARY KEY: (`menu_id`, `permission_id`)
- INDEX: `idx_menu_id` (`menu_id`)
- INDEX: `idx_permission_id` (`permission_id`)
- INDEX: `idx_deleted_at` (`deleted_at`)

### 3.9 活动日志表 (activities)

| 字段名 | 类型 | 约束 | 说明 | 默认值 |
| :--- | :--- | :--- | :--- | :--- |
| `id` | bigint unsigned | PRIMARY KEY, AUTO_INCREMENT | 日志ID | - |
| `user_id` | bigint unsigned | NOT NULL, INDEX | 用户ID | - |
| `username` | varchar(50) | NOT NULL, INDEX | 用户名 | - |
| `action` | varchar(255) | NOT NULL | 操作内容 | - |
| `url` | varchar(255) | NULL | 请求URL | '' |
| `ip` | varchar(50) | NULL | 客户端IP地址 | '' |
| `created_at` | datetime | NULL, INDEX | 创建时间 | CURRENT_TIMESTAMP |

**索引**:
- PRIMARY KEY: `id`
- INDEX: `idx_user_id` (`user_id`)
- INDEX: `idx_username` (`username`)
- INDEX: `idx_created_at` (`created_at`)

## 4. 数据模型关系图

```
┌───────────┐     N:N      ┌───────────┐
│   User    │◄────────────►│   Role    │
│  (用户)   │   user_roles │  (角色)   │
└───────────┘              └─────┬─────┘
                                 │
           ┌─────────────────────┼─────────────────────┐
           │                     │                     │
           ▼                     ▼                     ▼
┌───────────────┐     N:N    ┌───────────────┐     N:N    ┌───────────────┐
│ Permission    │◄───────────│ Role_Permission│───────────►│    Menu       │
│   (权限)      │             │  (角色权限关联) │             │   (菜单)      │
└───────────────┘             └───────────────┘             └───────────────┘
           │                                                     │
           └─────────────────────┬─────────────────────────────┘
                                 ▼
                       ┌───────────────┐
                       │Menu_Permission│
                       │ (菜单权限关联) │
                       └───────────────┘
```

## 5. 枚举值说明

### 5.1 用户状态 (users.status)

| 值 | 说明 |
| :---: | :--- |
| 1 | 启用 |
| 0 | 禁用 |

### 5.2 角色状态 (roles.status)

| 值 | 说明 |
| :---: | :--- |
| 1 | 启用 |
| 0 | 禁用 |

### 5.3 权限类型 (permissions.type)

| 值 | 说明 |
| :---: | :--- |
| 1 | 菜单权限 |
| 2 | 按钮权限 |

### 5.4 权限状态 (permissions.status)

| 值 | 说明 |
| :---: | :--- |
| 1 | 启用 |
| 0 | 禁用 |

### 5.5 菜单状态 (menus.status)

| 值 | 说明 |
| :---: | :--- |
| 1 | 启用 |
| 0 | 禁用 |

## 6. 初始数据说明

### 6.1 默认账号

| 用户名 | 密码 | 角色 | 说明 |
| :--- | :--- | :--- | :--- |
| admin | admin123 | 管理员 | 系统管理员，拥有所有权限 |
| test1 | admin123 | 普通用户 | 测试用户，拥有基础权限 |

### 6.2 默认角色

| 角色名称 | 角色编码 | 描述 |
| :--- | :--- | :--- |
| 管理员 | admin | 系统管理员，拥有所有权限 |
| 普通用户 | user | 普通用户，拥有基础权限 |

### 6.3 权限编码规范

```
{模块}:{操作}

模块名: user, role, permission, menu, activity
操作: list, create, update, delete, assign

示例:
- user:list              # 查看用户列表
- user:create            # 创建用户
- user:update            # 编辑用户
- user:delete            # 删除用户
- user:assign-roles      # 分配角色
- role:assign-permissions # 分配权限
- role:assign-menus      # 分配菜单
```

## 7. 使用说明

### 7.1 创建数据库

```sql
CREATE DATABASE IF NOT EXISTS admin_system 
CHARACTER SET utf8mb4 
COLLATE utf8mb4_unicode_ci;
```

### 7.2 执行建表脚本

```bash
mysql -h 127.0.0.1 -P 3306 -u root -proot admin_system < docs/sql/schema.sql
```

### 7.3 注意事项

1. 密码使用 BCrypt 加密存储，初始密码为 `admin123`
2. 所有关联表使用复合主键
3. 支持软删除（通过 `deleted_at` 字段）
4. 使用 UTC 时间戳
