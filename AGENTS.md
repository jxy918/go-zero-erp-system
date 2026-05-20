# AGENTS.md

## 项目概览

RBAC 后台管理系统，前后端分离架构。

| 组件 | 技术栈 | 端口 |
|------|--------|------|
| 后端 | Go + go-zero + GORM + MySQL | 8000 |
| 前端 | Vue 3 + Element Plus + Vite + Pinia | 3000 |
| 认证 | JWT | - |

**默认账号**: admin / admin123

---

## 关键命令

### 后端
```bash
cd admin
go mod tidy                                # 安装依赖
go run admin.go -f etc/admin-api.yaml      # 启动 (端口 8000)
go build -o admin.exe                       # 编译
```

### 前端
```bash
cd frontend
npm install          # 安装依赖
npm run dev          # 开发服务器 (端口 3000)
npm run build        # 生产构建
```

### 数据库
```bash
mysql -h 127.0.0.1 -P 3306 -u root -proot admin_system
```
- Host: `127.0.0.1:3306`
- Database: `admin_system`
- User: `root`
- Password: `root`
- 配置: `admin/etc/admin-api.yaml`

---

## 统一响应格式规范

所有后端接口必须遵循以下统一响应格式：

```json
{
  "code": 0,
  "data": {},
  "message": "success"
}
```

### 字段说明

| 字段 | 类型 | 说明 |
|------|------|------|
| `code` | int | 状态码，0 = 成功，非 0 = 失败 |
| `data` | any | 业务数据，成功时返回实际数据，失败时可为 null |
| `message` | string | 提示信息，成功时为 "success"，失败时为错误描述 |

### 状态码规范

| 状态码 | 含义 |
|--------|------|
| 0 | 成功 |
| 400 | 请求参数错误 |
| 401 | 未认证（Token无效或过期） |
| 403 | 无权限 |
| 500 | 服务器内部错误 |

---

## 字段命名规范

### 后端 API 定义

所有后端 API 定义（*.api 文件）中，JSON 字段必须使用**下划线格式（snake_case）**，严禁使用驼峰格式。

**正确示例**:
```go
type ProductInfo {
    ID          uint   `json:"id"`
    Name        string `json:"name"`
    ProductID   uint   `json:"product_id"`
    ProductName string `json:"product_name"`
    MinStock    int    `json:"min_stock"`
    MaxStock    int    `json:"max_stock"`
    CreatedAt   string `json:"created_at"`
}
```

**错误示例**（禁止使用）:
```go
type ProductInfo {
    ProductID   uint   `json:"productId"`      // ❌ 使用了驼峰格式
    ProductName string `json:"productName"`    // ❌ 使用了驼峰格式
    MinStock    int    `json:"minStock"`       // ❌ 使用了驼峰格式
}
```

### 前端数据引用

前端所有对后端返回数据的引用，必须使用**下划线格式（snake_case）**，与后端 API 定义保持一致。

**正确示例**:
```javascript
// Vue 模板中
{{ product.product_name }}

// JavaScript 中
const productId = data.product_id
const minStock = data.min_stock
```

**错误示例**（禁止使用）:
```javascript
// ❌ 使用了驼峰格式
{{ product.productName }}
const productId = data.productId
```

### 统一响应格式规范

### 示例

**成功响应**:
```json
{
  "code": 0,
  "data": { "id": 1, "name": "产品名称" },
  "message": "success"
}
```

**失败响应**:
```json
{
  "code": 400,
  "data": null,
  "message": "参数校验失败"
}
```

---

## 项目结构

```
go-zero-erp/
├── admin/                          # 后端 (go-zero)
│   ├── admin.go                    # 入口点 (⚠️ 路由在此注册，不再使用 routes.go)
│   ├── admin.api                   # API 定义 (goctl)
│   ├── etc/admin-api.yaml          # 配置文件
│   └── internal/
│       ├── config/                 # 配置结构体
│       ├── handler/                # HTTP 处理 (routes.go 已废弃)
│       ├── logic/                  # 业务逻辑
│       ├── middleware/             # 中间件 (认证、CORS)
│       ├── model/                  # GORM 模型
│       │   ├── models.go           # 核心模型 (User, Role, Permission, Menu)
│       │   └── *_model.go          # 各模块数据访问
│       ├── svc/                    # 服务上下文
│       ├── types/                  # 类型定义
│       └── util/                   # 工具 (JWT)
├── frontend/                        # 前端 (Vue 3)
│   ├── src/
│   │   ├── api/                    # Axios API 模块
│   │   ├── components/             # 公共组件
│   │   ├── router/                 # Vue Router
│   │   ├── store/                  # Pinia 状态
│   │   ├── views/                  # 页面
│   │   ├── directives/             # 自定义指令 (权限指令)
│   │   └── utils/                  # 工具函数
│   └── vite.config.js              # 代理: /api → localhost:8000
├── docs/                           # 项目文档
│   ├── requirements/               # 需求文档
│   ├── api/                        # API 文档
│   ├── bugs/                       # Bug 日志
│   └── changelog.md                 # 版本变更记录
├── go.mod
└── main.go                          # ⚠️ 独立 Gin 测试程序，非主入口
```

---

## 项目开发流程

### 流程总览

```
需求分析 → 菜单配置 → 权限配置 → 接口设计 → 文档生成
       ↓
后端开发 ←───────────→ 前端开发
       ↓
测试验证 → 部署上线 → Bug修复 → 文档更新
```

### 开发阶段

#### 阶段 1: 需求分析

**任务**:
- 明确功能目标和业务价值
- 分析功能边界和依赖关系
- 设计数据模型和业务流程
- 定义用户角色和权限需求

**输出**:
- 需求文档 (`docs/requirements/{功能名称}.md`)
- 数据模型设计文档
- API 接口设计文档

#### 阶段 2: 菜单配置（关键步骤）

**操作位置**: 后台 → 系统管理 → 菜单管理

| 字段 | 必填 | 说明 | 示例 |
|------|------|------|------|
| 菜单名称 | 是 | 显示名称 | 用户管理 |
| 菜单编码 | 是 | 唯一标识 | menu_user |
| 上级菜单 | 是 | 父菜单ID | 系统管理 |
| 路由路径 | 是 | 前端路由 | /user |
| 组件路径 | 是 | 组件文件 | User |
| 菜单图标 | 否 | Element图标 | User |
| 排序号 | 是 | 显示顺序 | 2 |
| 状态 | 是 | 启用/禁用 | 启用 |

**命名规范**:
- 菜单编码: `menu_{模块名}`
- 路由路径: `/模块名`
- 组件路径: 首字母大写驼峰

#### 阶段 3: 权限配置（关键步骤）

**操作位置**: 后台 → 系统管理 → 权限管理

| 字段 | 必填 | 说明 | 示例 |
|------|------|------|------|
| 权限名称 | 是 | 显示名称 | 查看用户列表 |
| 权限编码 | 是 | 唯一标识 | user:list |
| 权限类型 | 是 | 菜单/按钮 | 按钮权限 |
| 所属菜单 | 是 | 关联菜单 | 用户管理 |
| 状态 | 是 | 启用/禁用 | 启用 |

**权限编码规范**:
```
{模块}:{操作}

模块名: user, role, permission, menu, activity
操作: list, create, update, delete, assign

示例:
- user:list      # 查看用户列表
- user:create    # 创建用户
- user:update    # 更新用户
- user:delete    # 删除用户
```

#### 阶段 4: 接口设计与文档

**文档位置**: `docs/api/{接口名}.md`

#### 阶段 5: 后端开发

1. API 定义 (`admin/admin.api`)
2. 代码生成: `goctl api go -api admin.api -dir . -style goZero`
3. 逻辑实现 (`internal/logic/`)
4. 模型定义 (`internal/model/`)
5. 单元测试 (`internal/logic/*_test.go`)

#### API 接口参数修改规范

**⚠️ 重要规则：** 所有 API 接口的参数变更（新增、修改、删除字段）**必须首先**在 `admin/admin.api` 文件中进行，然后再执行代码生成命令。

**参数变更流程：**
1. **修改 `admin/admin.api`** - 在 API 定义文件中添加、修改或删除字段
2. **运行代码生成** - 执行 `goctl api go -api admin.api -dir . -style goZero`
3. **检查生成结果** - 确认 `internal/types/types.go` 已正确更新
4. **更新逻辑层** - 根据新增字段更新 `internal/logic/` 中的业务逻辑
5. **更新数据库模型** - 如果需要，更新 `internal/model/models.go`
6. **数据库迁移** - 必要时执行数据库字段添加或修改

**禁止操作：**
- ❌ 禁止直接修改 `internal/types/types.go` 而不同步 `admin/admin.api`
- ❌ 禁止在未更新 `admin.api` 的情况下运行 goctl 命令
- ❌ 禁止修改数据库字段而不更新模型定义
- ❌ 禁止修改后端代码后不重新启动服务

**⚠️ 重要规则：** 所有后端代码修改（包括逻辑层、模型层、配置文件等）**必须**重新编译并重启后端服务才能生效。

**重启命令：**
```bash
cd admin
go run admin.go -f etc/admin-api.yaml
```

**字段映射一致性检查清单：**
- [ ] `admin/admin.api` 中的 Request/Response 定义
- [ ] `internal/types/types.go` 中的结构体定义
- [ ] `internal/model/models.go` 中的数据库模型
- [ ] `internal/logic/` 中的字段赋值逻辑
- [ ] 数据库表结构（如需要）

**示例场景：**

```go
// 正确做法：先修改 admin.api
CreateProductRequest {
    Name       string  `json:"name"`
    Code       string  `json:"code"`
    Spec       string  `json:"spec"`       // ✅ 新增字段
    CostPrice  float64 `json:"cost_price"` // ✅ 新增字段
}

// 然后运行代码生成
goctl api go -api admin.api -dir . -style goZero

// 最后更新逻辑层和模型
```

**同步检查时机：**
- 在执行 goctl 命令之前，必须确认 `admin.api` 已更新
- 在提交代码之前，必须确认 `admin.api` 与 `types.go` 字段一致
- 在部署之前，必须运行完整的编译检查

#### 阶段 6: 前端开发

1. 创建组件 (`src/views/`)
2. 配置路由 (`src/router/`)
3. API 调用 (`src/api/`)
4. 权限控制

#### 阶段 7: 测试验证

- [ ] 菜单显示正确
- [ ] 页面访问正常
- [ ] 权限控制生效
- [ ] CRUD 操作正常
- [ ] 数据正确持久化
- [ ] API 返回正确

#### 阶段 8: 部署上线

- [ ] 代码审查通过
- [ ] 构建验证通过
- [ ] 测试用例全部通过
- [ ] 变更日志已更新

---

## 后端开发规范

### 职责范围

- 使用 go-zero 框架实现 RESTful API
- 在 `internal/logic/` 目录编写业务逻辑
- 使用 GORM 创建数据库模型
- 开发中间件（认证、CORS、日志）
- 实现基于 JWT 的认证机制
- 使用 GORM AutoMigrate 管理数据库迁移

### 添加新 API 流程

1. 编辑 `admin/admin.api` 定义 API
2. 运行 `goctl api go -api admin.api -dir . -style goZero` 生成代码
3. 在 `internal/logic/` 实现业务逻辑
4. 根据需要添加模型方法

### 代码生成

```bash
# 生成 go-zero API 代码
goctl api go -api admin.api -dir . -style goZero
```

### ⚠️ 生成文件修改规范

**重要规则：** 所有由 `goctl` 命令生成的文件（包括 `handler/`、`logic/`、`types/` 目录下的文件）**禁止直接修改**。如需修改，必须按照以下流程操作：

1. **修改 API 定义** - 在 `admin/api/` 目录下的对应 `.api` 文件中修改接口定义
2. **重新生成代码** - 执行 `goctl api go -api admin.api -dir . -style goZero` 命令重新生成
3. **检查生成结果** - 确认生成的文件已正确更新
4. **更新逻辑层** - 仅在 `internal/logic/` 目录中添加业务逻辑（不修改自动生成的函数签名）

**禁止操作：**
- ❌ 禁止直接修改 `internal/handler/` 目录下的文件
- ❌ 禁止直接修改 `internal/types/` 目录下的文件
- ❌ 禁止修改 `internal/logic/` 中自动生成的函数签名

**允许操作：**
- ✅ 在 `internal/logic/` 中实现业务逻辑（在函数体内添加代码）
- ✅ 在 `internal/util/` 中添加工具函数
- ✅ 在 `internal/model/` 中添加数据模型
- ✅ 在 `internal/middleware/` 中添加中间件

### 重要约定

- 使用 GORM v1，标签格式 `gorm:"..."`
- 软删除: `DeletedAt gorm.DeletedAt`
- 多对多关系: `gorm:"many2many:table_name;"`
- 遵循 go-zero 项目结构和命名规范
- Handler 命名: `{操作}{实体}Handler`（如 `CreateUserHandler`）
- Logic 命名: `{操作}{实体}Logic`（如 `CreateUserLogic`）
- Model 命名: `{实体}Model` 接口（如 `UserModel`）

### 文件命名规范

**⚠️ 后端所有目录的 Go 文件必须统一使用小驼峰命名法**（lowerCamelCase），禁止使用下划线命名法**：

| 文件类型 | 命名规则 | 示例 |
|----------|----------|------|
| Handler 文件 | `{操作}{实体}Handler.go` | `createUserHandler.go`, `listMenuHandler.go` |
| Logic 文件 | `{操作}{实体}Logic.go` | `createUserLogic.go`, `listMenuLogic.go` |
| Model 文件 | `{实体}Model.go` | `userModel.go`, `menuModel.go`, `inventoryChangeModel.go` |
| 配置文件 | `{模块}Config.go` | `config.go` |
| 中间件 | `{功能}Middleware.go` | `authMiddleware.go` |
| 工具函数 | `{功能}.go` | `jwt.go`, `dataPermission.go`, `response.go` |
| 测试文件 | `{功能}Test.go` | `purchaseOrderTest.go` |
| 其他文件 | 小驼峰命名 | `initData.go`, `migration.go` |

**操作前缀说明**:
- `create` - 创建
- `update` - 更新
- `delete` - 删除
- `get` - 获取单个
- `list` - 获取列表
- `assign` - 分配（如角色分配权限）
- `login` - 登录
- `logout` - 登出
- `refresh` - 刷新

**正确示例**:
- ✅ `createUserHandler.go`
- ✅ `getMenuListHandler.go`
- ✅ `updateRoleLogic.go`
- ✅ `userModel.go`
- ✅ `dataPermission.go`
- ✅ `initData.go`

**错误示例**:
- ❌ `CreateUserHandler.go`（首字母大写错误）
- ❌ `createuserhandler.go`（全小写，不符合驼峰）
- ❌ `CreateuserHandler.go`（不规范的大小写）
- ❌ `user_model.go`（使用下划线，错误）
- ❌ `data_permission.go`（使用下划线，错误）
- ❌ `init_data.go`（使用下划线，错误）
- ❌ `purchase_order_test.go`（使用下划线，错误）

**重要规则**：
- 所有 Go 目录下的文件都必须使用小驼峰命名
- 禁止在文件名中使用下划线分隔单词
- 多个单词组合时，除第一个单词外，其他单词首字母大写

### 认证机制

- JWT token 在 `internal/util/jwt.go` 中生成和验证
- 认证中间件在 `internal/middleware/auth.go`
- 前端通过 `Authorization: Bearer <token>` 头发送 token
- 所有路由（除 `/auth/login` 和 `/auth/refresh`）都需要 JWT 认证
- Activity 日志通过 `go svc.ActivityModel.Create()` 异步写入


## 后端统一返回格式

```json
{ "code": 0, "data": {}, "message": "success" }
```

---

## 前端开发规范

### 职责范围

- 开发 Vue 3 组件和页面
- 使用 Vue Router 实现路由
- 使用 Pinia 管理应用状态
- 与后端 API 集成
- 使用 Element Plus 实现响应式 UI
- 处理认证和授权

### 项目结构

```
frontend/
├── src/
│   ├── views/                  # 页面组件
│   ├── components/             # 可复用组件
│   ├── router/                 # 路由配置
│   ├── store/                  # Pinia stores
│   ├── api/                    # API 客户端
│   ├── directives/             # 自定义指令 (权限指令)
│   └── utils/                 # 工具函数
└── vite.config.js              # Vite 配置
```

### 常用页面路由

| 路径 | 页面 | 说明 |
|------|------|------|
| `/` | 控制台 | Dashboard 首页 |
| `/login` | 登录页 | 用户登录 |
| `/user` | 用户管理 | 用户列表、创建、编辑 |
| `/role` | 角色管理 | 角色列表、权限分配 |
| `/permission` | 权限管理 | 权限列表、创建、编辑 |
| `/menu` | 菜单管理 | 菜单树、创建、编辑 |
| `/activity` | 活动日志 | 操作日志查询 |

### 重要约定

#### API 代理配置
```js
// vite.config.js
server: {
  proxy: {
    '/api': {
      target: 'http://localhost:8000',
      changeOrigin: true,
      rewrite: (path) => path.replace(/^\/api/, '')
    }
  }
}
```
- `/api` 前缀会被重写移除
- 例如: `/api/user/list` → 后端收到 `/user/list`

#### Token 存储
- `localStorage.getItem('token')` - 访问令牌
- `localStorage.getItem('user')` - 用户信息（JSON 字符串）

#### 主题色
- Element Plus 主题色: `#2c7d34` (绿色)

### 认证流程

1. 用户通过 `/auth/login` 登录
2. Token 存储在 localStorage
3. 每个请求通过 Axios 拦截器添加 `Authorization: Bearer <token>` 头
4. 遇到 401 错误时，用户重定向到登录页面

### 状态管理 (Pinia)

主要 Store：
- `store/user.js` - 用户状态、Token、权限
- `store/menu.js` - 菜单树
- `store/activity.js` - 活动日志

### 权限控制

- 自定义指令 `v-has-permission`: 控制按钮/元素的显示权限
- 路由守卫: 验证用户登录状态和菜单权限

---

## API 路由

| 方法 | 路径 | 认证 | 说明 |
|------|------|------|------|
| POST | /auth/login | ❌ | 登录 |
| POST | /auth/refresh | ❌ | 刷新 Token |
| POST | /auth/logout | ✅ | 登出 |
| GET | /user/list | ✅ | 用户列表 |
| POST | /user/create | ✅ | 创建用户 |
| POST | /user/update | ✅ | 更新用户 |
| POST | /user/delete | ✅ | 删除用户 |
| GET | /user/get/:id | ✅ | 获取用户 |
| POST | /user/assign-roles | ✅ | 分配角色 |
| GET | /role/list | ✅ | 角色列表 |
| POST | /role/create | ✅ | 创建角色 |
| PUT | /role/update | ✅ | 更新角色 |
| DELETE | /role/delete | ✅ | 删除角色 |
| GET | /role/get/:id | ✅ | 获取角色 |
| POST | /role/assign-permissions | ✅ | 分配权限 |
| POST | /role/assign-menus | ✅ | 分配菜单 |
| GET | /permission/list | ✅ | 权限列表 |
| POST | /permission/create | ✅ | 创建权限 |
| PUT | /permission/update | ✅ | 更新权限 |
| DELETE | /permission/delete | ✅ | 删除权限 |
| GET | /permission/get/:id | ✅ | 获取权限 |
| GET | /menu/tree | ✅ | 菜单树 |
| GET | /menu/list | ✅ | 菜单列表 |
| POST | /menu/create | ✅ | 创建菜单 |
| POST | /menu/update | ✅ | 更新菜单 |
| POST | /menu/delete | ✅ | 删除菜单 |
| GET | /menu/get/:id | ✅ | 获取菜单 |
| POST | /menu/assign-permissions | ✅ | 分配菜单权限 |
| GET | /activity/list | ✅ | 活动日志 |

---

## 统一响应格式

```json
{ "code": 0, "data": {}, "message": "success" }
```

- `code: 0` = 成功，非 0 = 失败
- 前端 Axios 拦截器在 `frontend/src/main.js`

---

## 数据库管理规范

### 数据库配置

```yaml
# admin/etc/admin-api.yaml
DataSource: root:root@tcp(127.0.0.1:3306)/admin_system?charset=utf8mb4&parseTime=True&loc=Local
```

### 数据库表结构

1. **users** - 用户信息
2. **roles** - 角色定义
3. **permissions** - 权限定义
4. **menus** - 菜单定义
5. **user_roles** - 用户-角色关系
6. **role_permissions** - 角色-权限关系
7. **role_menus** - 角色-菜单关系
8. **menu_permissions** - 菜单-权限关系
9. **activities** - 活动日志

### 核心模型关系

| 模型 | 说明 |
|------|------|
| User ↔ Role | 多对多 (user_roles) |
| Role ↔ Permission | 多对多 (role_permissions) |
| Role ↔ Menu | 多对多 (role_menus) |
| Menu ↔ Permission | 多对多 (menu_permissions) |

### 用户状态

- `Status = 1`: 启用
- `Status = 0`: 禁用

### 模型约定

- 使用 GORM v1，标签格式 `gorm:"..."`
- 软删除: `DeletedAt gorm.DeletedAt`
- 多对多: `gorm:"many2many:table_name;"`
- 表名: 结构体名称的复数形式（如 `User` → `users`）

### 迁移

项目使用 GORM AutoMigrate，启动时自动创建或更新表：

```go
db.AutoMigrate(&User{}, &Role{}, &Permission{}, &Menu{}, &Activity{})
```

### 常见任务

- **重置数据库**: `DROP DATABASE admin_system;` 后端会自动重建
- **添加新模型**: 在 `internal/model/` 创建 → 在 `db.go` 添加 AutoMigrate → 在 `init_data.go` 添加初始数据

---

## API 管理规范

### API 定义文件结构

为了提高代码可维护性和协作效率，API 定义已按业务模块拆分为多个文件：

```
admin/
├── api/                    # API 模块目录
│   ├── auth.api            # 认证模块（登录/登出/刷新Token）
│   ├── user.api            # 用户管理
│   ├── role.api            # 角色管理
│   ├── permission.api      # 权限管理
│   ├── menu.api            # 菜单管理
│   ├── product.api         # 产品管理（含分类、单位）
│   ├── supplier.api        # 供应商管理
│   ├── customer.api        # 客户管理
│   ├── warehouse.api       # 仓库管理
│   ├── purchase.api        # 采购订单
│   ├── sales.api           # 销售订单
│   ├── inventory.api       # 库存管理（记录/调整/盘点/调拨）
│   ├── erp.api             # ERP 统计报表
│   └── activity.api        # 活动日志
└── admin.api               # 主入口文件
```

### 主入口文件格式

`admin/admin.api` 作为主入口，通过 `import` 导入所有模块：

```go
syntax = "v1"

import (
    "api/auth.api"
    "api/user.api"
    "api/role.api"
    "api/permission.api"
    "api/menu.api"
    "api/product.api"
    "api/supplier.api"
    "api/customer.api"
    "api/warehouse.api"
    "api/purchase.api"
    "api/sales.api"
    "api/inventory.api"
    "api/erp.api"
    "api/activity.api"
)

type (
    EmptyResponse {}
)
```

### 添加新接口规范

**⚠️ 重要规则**：新增接口必须按照所属业务模块添加到对应的 `.api` 文件中：

| 模块 | 文件 | 说明 |
|------|------|------|
| 认证 | `api/auth.api` | 登录、登出、刷新Token |
| 用户管理 | `api/user.api` | 用户CRUD、角色分配 |
| 角色管理 | `api/role.api` | 角色CRUD、权限分配、菜单分配 |
| 权限管理 | `api/permission.api` | 权限CRUD |
| 菜单管理 | `api/menu.api` | 菜单CRUD、权限分配 |
| 产品管理 | `api/product.api` | 产品CRUD、分类、单位管理 |
| 供应商管理 | `api/supplier.api` | 供应商CRUD |
| 客户管理 | `api/customer.api` | 客户CRUD |
| 仓库管理 | `api/warehouse.api` | 仓库CRUD |
| 采购订单 | `api/purchase.api` | 采购订单CRUD、状态变更 |
| 销售订单 | `api/sales.api` | 销售订单CRUD、状态变更 |
| 库存管理 | `api/inventory.api` | 库存记录、调整、盘点、调拨 |
| ERP报表 | `api/erp.api` | 统计报表、仪表盘数据 |
| 活动日志 | `api/activity.api` | 操作日志查询 |

### 代码生成

```bash
goctl api go -api admin.api -dir . -style goZero
```

### ⚠️ 重要注意事项

**routes.go 已启用**：路由不在 `admin.go` 中注册，采用命令生成 `routes.go`

### API 定义格式

```go
type LoginRequest struct {
    Username string `json:"username"`
    Password string `json:"password"`
}

type LoginResponse struct {
    Token  string   `json:"token"`
    User   UserInfo `json:"user"`
}

service admin-api {
    @handler LoginHandler
    post /auth/login (LoginRequest) returns (LoginResponse)
}
```

---

## Bug 修复规范

### Bug 日志结构

**文档位置**: `docs/bugs/{日期}_{bug编号}.md`

### Bug 编号规范

```
BUG-{三位数字}

示例:
- BUG-001
- BUG-002
```

### 严重程度定义

| 级别 | 定义 | 处理优先级 |
|------|------|------------|
| 高 | 影响核心功能，导致系统不可用 | 立即处理 |
| 中 | 影响部分功能，用户可通过其他方式操作 | 24小时内处理 |
| 低 | 界面显示问题或不影响功能的小问题 | 下个迭代处理 |

### Bug 修复检查清单

- [ ] Bug 日志已记录
- [ ] 修复方案已确认
- [ ] 验证测试通过
- [ ] 相关文档已更新
- [ ] Bug 状态已关闭

---

## 文档管理规范

### 文档结构

```
go-zero-erp/
├── test/                      # Python 测试脚本（数据修复、数据库操作等）
│   └── {功能名}.py
├── docs/
│   ├── requirements/           # 需求文档
│   │   └── {功能名}_v{版本}.md
│   ├── api/                    # API 文档
│   │   └── {接口名}.md
│   ├── sql/                    # SQL 脚本
│   │   └── {功能名}_{操作}.sql
│   ├── testing/                # 测试用例文档
│   │   └── {功能名}_test.md
│   ├── bugs/                   # Bug 日志
│   │   └── {日期}_{编号}.md
│   └── changelog.md            # 版本变更记录
```

### 文件存放规范

| 文件类型 | 存放目录 | 命名规范 | 示例 |
|----------|----------|----------|------|
| Python 脚本 | `test/` | `{功能名}.py` | `update_menus.py` |
| SQL 脚本 | `docs/sql/` | `{功能名}_{操作}.sql` | `add_permissions.sql` |
| 测试用例 | `docs/testing/` | `{功能名}_test.md` | `erp_statistics_test.md` |
| 需求文档 | `docs/requirements/` | `{功能名}_v{版本}.md` | `user_management_v1.0.md` |
| API 文档 | `docs/api/` | `{接口名}.md` | `user_list.md` |
| Bug 日志 | `docs/bugs/` | `{日期}_{编号}.md` | `20240115_BUG-001.md` |

### 文档命名规范

| 文档类型 | 命名格式 | 示例 |
|----------|----------|------|
| 需求文档 | `{功能名}_v{版本}.md` | user_management_v1.0.md |
| API 文档 | `{接口名}.md` | user_list.md |
| Bug 日志 | `{日期}_{编号}.md` | 20240115_BUG-001.md |

### 文档更新规则

1. **及时性**: 代码变更后 24 小时内更新相关文档
2. **完整性**: 文档必须包含所有必要信息
3. **一致性**: 文档与代码保持同步
4. **可追溯性**: 文档之间相互引用

---

## 安全注意

- JWT 密钥默认 `your-secret-key`，生产需修改
- CORS 当前允许所有来源 (`*`)，生产需限制
- 密码使用 bcrypt 加密存储

---

## 参考文档

- `admin/admin.api` - API 定义
- `admin/internal/middleware/auth.go` - 认证逻辑
- `admin/internal/model/models.go` - 数据模型
- `PROJECT_ANALYSIS.md` - 详细项目分析
- `docs/backend-admin-system-prd.md` - 产品需求
- `.trae/skills/backend-developer/SKILL.md` - 后端开发规范
- `.trae/skills/frontend-developer/SKILL.md` - 前端开发规范
- `.trae/skills/project-workflow/SKILL.md` - 项目工作流规范
- `.trae/skills/database-manager/SKILL.md` - 数据库管理规范
- `.trae/skills/api-manager/SKILL.md` - API 管理规范
