# 后端服务

## 概述

基于 Go + go-zero + GORM + MySQL 的后端服务，为 RBAC 权限管理系统和 ERP 管理系统提供 REST API 接口。

## 技术栈

- **框架**: go-zero
- **语言**: Go 1.20+
- **数据库**: MySQL 8.0+
- **ORM**: GORM v1
- **认证**: JWT

## 项目结构

```
admin/
├── admin.go              # 应用入口
├── admin.api             # API 定义文件
├── etc/
│   └── admin-api.yaml    # 配置文件
└── internal/
    ├── config/           # 配置结构体
    ├── handler/          # HTTP 处理器
    ├── logic/            # 业务逻辑
    ├── middleware/       # 中间件
    ├── model/            # 数据模型
    ├── svc/              # 服务上下文
    ├── types/            # 类型定义
    └── util/             # 工具函数
```

## 快速开始

### 环境要求

- Go 1.20+
- MySQL 8.0+

### 安装依赖

```bash
go mod tidy
```

### 配置数据库

修改 `etc/admin-api.yaml`：

```yaml
DataSource: root:root@tcp(127.0.0.1:3306)/admin_system?charset=utf8mb4&parseTime=True&loc=Local
```

### 启动服务

```bash
go run admin.go -f etc/admin-api.yaml
```

服务将在 `http://localhost:8000` 启动。

## API 接口

### 认证接口
- `POST /auth/login` - 用户登录
- `POST /auth/logout` - 用户登出
- `POST /auth/refresh` - 刷新 Token

### 用户管理
- `GET /user/list` - 获取用户列表
- `GET /user/get/:id` - 获取用户详情
- `POST /user/create` - 创建用户
- `POST /user/update` - 更新用户
- `POST /user/delete` - 删除用户
- `POST /user/assign-roles` - 分配角色

### 角色管理
- `GET /role/list` - 获取角色列表
- `GET /role/get/:id` - 获取角色详情
- `POST /role/create` - 创建角色
- `PUT /role/update` - 更新角色
- `DELETE /role/delete` - 删除角色
- `POST /role/assign-permissions` - 分配权限
- `POST /role/assign-menus` - 分配菜单

### 权限管理
- `GET /permission/list` - 获取权限列表
- `GET /permission/get/:id` - 获取权限详情
- `POST /permission/create` - 创建权限
- `PUT /permission/update` - 更新权限
- `DELETE /permission/delete` - 删除权限

### 菜单管理
- `GET /menu/tree` - 获取菜单树
- `GET /menu/list` - 获取菜单列表
- `GET /menu/get/:id` - 获取菜单详情
- `POST /menu/create` - 创建菜单
- `POST /menu/update` - 更新菜单
- `POST /menu/delete` - 删除菜单
- `POST /menu/assign-permissions` - 分配菜单权限

### 活动日志
- `GET /activity/list` - 获取活动日志

### 产品管理
- `GET /product/list` - 获取产品列表
- `GET /product/get/:id` - 获取产品详情
- `POST /product/create` - 创建产品
- `POST /product/update` - 更新产品
- `POST /product/delete` - 删除产品
- `GET /product/category/list` - 获取产品分类列表
- `POST /product/category/create` - 创建产品分类

### 供应商管理
- `GET /supplier/list` - 获取供应商列表
- `GET /supplier/get/:id` - 获取供应商详情
- `POST /supplier/create` - 创建供应商
- `POST /supplier/update` - 更新供应商
- `POST /supplier/delete` - 删除供应商

### 客户管理
- `GET /customer/list` - 获取客户列表
- `GET /customer/get/:id` - 获取客户详情
- `POST /customer/create` - 创建客户
- `POST /customer/update` - 更新客户
- `POST /customer/delete` - 删除客户

### 仓库管理
- `GET /warehouse/list` - 获取仓库列表
- `GET /warehouse/get/:id` - 获取仓库详情
- `POST /warehouse/create` - 创建仓库
- `POST /warehouse/update` - 更新仓库
- `POST /warehouse/delete` - 删除仓库

### 采购管理
- `GET /purchase/list` - 获取采购订单列表
- `GET /purchase/get/:id` - 获取采购订单详情
- `POST /purchase/create` - 创建采购订单
- `POST /purchase/update` - 更新采购订单
- `POST /purchase/delete` - 删除采购订单
- `POST /purchase/approve` - 审核采购订单
- `POST /purchase/inbound` - 采购入库

### 销售管理
- `GET /sales/list` - 获取销售订单列表
- `GET /sales/get/:id` - 获取销售订单详情
- `POST /sales/create` - 创建销售订单
- `POST /sales/update` - 更新销售订单
- `POST /sales/delete` - 删除销售订单
- `POST /sales/approve` - 审核销售订单
- `POST /sales/outbound` - 销售出库

### 库存管理
- `GET /inventory/list` - 获取库存列表
- `GET /inventory/history` - 获取库存历史记录
- `GET /inventory/current-stock` - 获取当前库存
- `POST /inventory/adjust-request/create` - 创建库存调整申请
- `GET /inventory/adjust-request/list` - 获取库存调整申请列表
- `POST /inventory/adjust-request/approve` - 审核库存调整申请
- `POST /inventory/adjust-request/reject` - 拒绝库存调整申请

### ERP 统计报表
- `GET /erp/statistics/overview` - 获取概览统计
- `GET /erp/statistics/trend` - 获取采购/销售趋势
- `GET /erp/statistics/inventory-alert` - 获取库存预警
- `GET /erp/statistics/top-products` - 获取热销商品排行
- `GET /erp/statistics/order-status` - 获取订单状态分布
- `GET /erp/statistics/business` - 获取供应商/客户统计

### 系统接口
- `POST /system/init-data` - 初始化数据

## 开发流程

1. 编辑 `admin.api` 定义 API
2. 运行 `goctl api go -api admin.api -dir . -style goZero` 生成代码
3. 在 `internal/logic/` 实现业务逻辑

## 默认账号

- **用户名**: admin
- **密码**: admin123

## 安全注意

1. JWT 密钥默认使用 `your-secret-key`，生产环境请修改
2. CORS 当前允许所有来源，生产环境请限制
3. 建议生产环境启用 HTTPS
