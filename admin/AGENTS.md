# 后端管理系统 AGENTS.md

## 项目概述

基于 Go + go-zero + GORM + MySQL 的后端服务，为 RBAC 权限管理系统提供 REST API 接口。

## 常用命令

```bash
cd admin
go run admin.go -f etc/admin-api.yaml     # 启动服务（端口 8000）
goctl api go -api admin.api -dir .        # 从 API 定义重新生成代码
go mod tidy                              # 安装依赖
go build -o admin.exe                     # 编译可执行文件
```

## 项目结构

```
admin/
├── admin.go           # 应用入口（路由在此注册）
├── admin.api          # API 定义文件（修改后需运行 goctl 生成代码）
├── etc/               # YAML 配置文件
│   └── admin-api.yaml # 主配置文件（数据库连接、JWT 密钥等）
└── internal/
    ├── config/        # 配置结构体定义
    ├── handler/       # HTTP 请求处理器
    ├── logic/         # 业务逻辑层
    ├── middleware/    # 中间件（认证、CORS、数据权限）
    ├── model/         # GORM 数据模型
    │   ├── models.go  # 核心模型定义（User、Role、Permission、Menu）
    │   └── *_model.go # 各模块数据访问方法
    ├── svc/           # 服务上下文
    ├── types/         # 请求/响应类型定义
    └── util/          # 工具函数（JWT、IP 获取等）
```

## 开发规范

### 命名约定
- **Handler 命名**: `{操作}{实体}Handler`（例如：`CreateUserHandler`）
- **Logic 命名**: `{操作}{实体}Logic`（例如：`CreateUserLogic`）
- **Model 命名**: `{实体}Model` 接口（例如：`UserModel`）

### API 开发流程
1. 编辑 `admin.api` 定义 API
2. 运行 `goctl api go -api admin.api -dir . -style goZero` 生成代码
3. 在 `internal/logic/` 实现业务逻辑
4. 根据需要添加模型方法

## 注意事项（避免做法）
- **不要手动编辑 `routes.go`** - 该文件由 goctl 自动生成
- **不要在未更新 `admin.api` 的情况下添加路由**
- **不要在未通过 `go build` 验证的情况下提交代码**

## 关键文件说明

| 文件路径 | 用途 |
|----------|------|
| `admin.go` | 应用入口，路由注册 |
| `admin.api` | API 定义文件 |
| `etc/admin-api.yaml` | 配置文件（数据库、JWT、服务器端口） |
| `internal/handler/routes.go` | 路由注册（自动生成） |
| `internal/middleware/auth.go` | JWT 认证中间件 + 活动日志记录 |
| `internal/middleware/data_permission.go` | 数据权限控制工具函数 |
| `internal/util/jwt.go` | Token 生成与解析 |
| `internal/model/models.go` | 核心 GORM 模型定义 |
| `internal/svc/servicecontext.go` | 服务上下文初始化 |

## 认证机制
- 除 `/auth/login` 和 `/auth/refresh` 外，所有路由都需要 JWT Bearer token
- JWT 密钥在 `etc/admin-api.yaml` 中配置（生产环境需修改默认密钥）
- 前端通过 `Authorization: Bearer <token>` 头发送认证信息

## 数据权限
- **管理员**: 可以查看所有数据
- **普通用户**: 只能查看自己角色范围内的数据

## 数据库配置
```yaml
# etc/admin-api.yaml
DataSource: root:root@tcp(127.0.0.1:3306)/admin_system?charset=utf8mb4&parseTime=True&loc=Local
```

## 数据初始化
系统启动时自动创建表结构，可通过以下接口初始化基础数据：
```
POST /system/init-data
```