# AGENTS.md

## 项目概览

RBAC + ERP 管理系统，前后端分离架构。

| 组件 | 技术栈 | 端口 |
|------|--------|------|
| 后端 | Go + go-zero + GORM + MySQL | 8001 |
| 前端 | Vue 3 + Element Plus + Vite + Pinia | 3000 |
| 认证 | JWT | - |

**默认账号**: admin / admin123

---

## 关键命令

### 后端
```bash
cd admin
go mod tidy                                # 安装依赖
go run admin.go -f etc/admin-api.yaml      # 启动 (端口 8001)
go build -o admin.exe                       # 编译
```

### 监控端点（DevServer）
服务启动后自动启用，端口 **6060**（配置在 `admin/etc/admin-api.yaml`）：

| 端点 | 功能 | 说明 |
|------|------|------|
| `http://localhost:6060/healthz` | 健康检查 | 返回 "OK" 表示服务正常 |
| `http://localhost:6060/metrics` | Prometheus 指标 | CPU、内存、GC、HTTP 请求统计 |
| `http://localhost:6060/debug/pprof/` | pprof 性能分析 | CPU/Memory/Goroutine/Block 分析 |

**pprof 使用示例**:
```bash
# CPU 分析（30秒）
go tool pprof http://localhost:6060/debug/pprof/profile?seconds=30

# 内存分析
go tool pprof http://localhost:6060/debug/pprof/heap

# Goroutine 分析
go tool pprof http://localhost:6060/debug/pprof/goroutine

# 生成火焰图（需安装 graphviz）
go tool pprof -http=:8080 http://localhost:6060/debug/pprof/heap
```

**⚠️ 生产环境注意**: DevServer 默认仅绑定 `127.0.0.1`，如需远程访问需修改配置或通过反向代理。

### 监控配置开关
所有监控功能均可通过配置文件独立开关：

```yaml
# DevServer 开关（控制 /metrics, /debug/pprof/, /healthz）
DevServer:
  Enabled: true          # false 完全关闭 DevServer
  EnableMetrics: true    # false 仅关闭 Prometheus 指标
  EnablePprof: false     # false 关闭 pprof 性能分析端点

# 分布式追踪开关（需要外部 Jaeger 依赖）
# 注释掉整个 Telemetry 块或设置 Sampler: 0.0 即可关闭
# Telemetry:
#   Name: admin-api
#   Batcher: otlpgrpc
#   Endpoint: http://localhost:4317
#   Sampler: 1.0

# 自定义业务指标开关（控制登录计数、订单统计等）
Metrics:
  Enabled: true          # false 关闭自定义业务指标
```

**关闭外部依赖的配置**（推荐默认状态）:
- ✅ DevServer: 开启（无外部依赖，提供基础监控）
- ❌ Telemetry: 关闭（需要 Jaeger 容器）
- ✅ Metrics: 开启（无外部依赖，业务指标收集）

### 自定义业务指标
服务启动后自动启用，通过 `http://localhost:6060/metrics` 暴露 Prometheus 格式指标：

| 指标名称 | 类型 | 标签 | 说明 |
|----------|------|------|------|
| `erp_auth_login_total` | Counter | `status` (success/failure) | 登录尝试次数 |
| `erp_api_request_total` | Counter | `method`, `path`, `status` | API 请求总数 |
| `erp_api_request_duration_ms` | Histogram | `method`, `path` | API 请求耗时（ms） |
| `erp_db_query_duration_ms` | Histogram | `table`, `operation` | 数据库查询耗时（ms） |
| `erp_order_create_total` | Counter | `type`, `status` | 订单创建次数（purchase/sales） |
| `erp_inventory_adjust_total` | Counter | `type` (inbound/outbound/adjust) | 库存调整次数 |
| `erp_system_active_users` | Gauge | `status` (enabled/disabled) | 活跃用户数 |

**指标集成示例**:
```go
// 在 Logic 层记录业务指标
import "myproject/admin/internal/metric"

// 登录成功
metric.LoginCounter.Inc("success")

// 记录 API 耗时（中间件自动收集）
metric.ObserveApiRequest(r.Method, r.URL.Path, duration)

// 记录订单创建
metric.OrderCreateCounter.Inc("purchase", "success")
```

**Grafana 查询示例**:
```promql
# 登录成功率
rate(erp_auth_login_total{status="success"}[5m]) / rate(erp_auth_login_total[5m])

# API P99 延迟
histogram_quantile(0.99, rate(erp_api_request_duration_ms_bucket[5m]))

# 采购订单创建速率
rate(erp_order_create_total{type="purchase"}[1m])
```

### 分布式追踪（OpenTelemetry）
服务启动后自动启用，追踪数据上报到 Jaeger（需先启动 Jaeger 容器）：

**启动 Jaeger**:
```bash
cd tools
docker-compose -f docker-compose-jaeger.yaml up -d
# 或运行 start-jaeger.bat
```

**访问 Jaeger UI**: http://localhost:16686

**追踪配置** (`admin/etc/admin-api.yaml`):
```yaml
Telemetry:
  Name: admin-api
  Batcher: otlpgrpc              # otlpgrpc|otlphttp|zipkin|file
  Endpoint: http://localhost:4317 # Jaeger OTLP gRPC 端点
  Sampler: 1.0                   # 1.0=100%, 0.1=10%
```

**调试模式**: 使用 `Batcher: file` 输出到 `traces.json` 文件验证追踪是否生效。

### 前端
```bash
cd frontend
npm install          # 安装依赖
npm run dev          # 开发服务器 (端口 3000)
npm run build        # 生产构建
```

### 数据库
```bash
mysql -h 127.0.0.1 -P 3306 -u root -proot admin_erp_system
```
- Host: `127.0.0.1:3306`
- Database: `admin_erp_system`
- User: `root`
- Password: `root`

---

## ⚠️ 关键注意事项

### 1. 代码生成流程
- **API 定义文件**: `admin/admin.api`（主入口，import 所有 `admin/api/*.api` 模块）
- **代码生成命令**: `goctl api go -api admin.api -dir . -style goZero`
- **禁止直接修改**: `internal/handler/`、`internal/types/` 下的生成文件
- **只允许修改**: `internal/logic/` 中的业务逻辑函数体（不修改函数签名）

### 2. 路由注册
- 路由通过 `handler.RegisterHandlers(server, ctx)` 注册（`internal/handler/routes.go` 由 goctl 自动生成）
- 中间件在 `admin.go` 中**手动**按顺序注册：ResponseMiddleware → CorsMiddleware → AuthMiddleware → PermissionMiddleware

### 3. 跳过认证的路径
`/auth/login`、`/auth/refresh`、`/system/init-data`

### 4. 端口与数据库
- 后端端口: **8001**（配置在 `admin/etc/admin-api.yaml`）
- 数据库名: **admin_erp_system**
- 前端代理目标: `localhost:8001`（配置在 `frontend/vite.config.js`）

### 5. go.mod 模块名
模块名为 `myproject`，import 路径示例: `myproject/admin/internal/model`

---

## 统一响应格式

```json
{ "code": 0, "data": {}, "message": "success" }
```

- `code: 0` = 成功，非 0 = 失败
- 通过 `internal/middleware/responseMiddleware.go` 中间件统一包装
- 工具函数: `internal/util/response.go` 提供 `SuccessResponse` / `ErrorResponse` / `ParseRequest`

---

## 字段命名规范

**所有 API JSON 字段必须使用 snake_case**（下划线格式）。

```go
// ✅ 正确
ProductID   uint   `json:"product_id"`
CreatedAt   string `json:"created_at"`

// ❌ 错误
ProductID   uint   `json:"productId"`
```

前端引用后端数据同样使用 snake_case。

---

## 后端开发规范

### 文件命名
所有 Go 文件使用 **lowerCamelCase**（小驼峰），禁止下划线：

| 类型 | 示例 |
|------|------|
| Handler 文件 | `createUserHandler.go` |
| Logic 文件 | `createUserLogic.go` |
| Model 文件 | `userModel.go` |
| 工具文件 | `jwt.go`, `response.go` |

### API 模块划分
API 定义按业务模块拆分在 `admin/api/` 目录：

| 模块 | 文件 | 说明 |
|------|------|------|
| 认证 | `auth.api` | 登录/登出/刷新Token |
| 系统管理 | `user.api`, `role.api`, `permission.api`, `menu.api`, `activity.api` | RBAC 核心 |
| ERP 基础 | `product.api`, `supplier.api`, `customer.api`, `warehouse.api` | 基础数据 |
| ERP 业务 | `purchase.api`, `sales.api`, `inventory.api` | 订单/库存 |
| ERP 报表 | `erp.api` | 统计报表 |

### 认证机制
- JWT 工具: `internal/util/jwt.go`
- 认证中间件: `internal/middleware/auth.go`
- 权限中间件: `internal/middleware/` (路径级权限检查)
- 请求头: `Authorization: Bearer <token>`
- Activity 日志通过 `svc.ActivityModel.Create()` 异步写入

### 数据模型
- 核心模型定义: `internal/model/models.go`
- 各模块数据访问: `internal/model/*Model.go`
- 数据库初始化/迁移: `internal/model/db.go`（自动建库建表 + 手动 SQL 迁移）
- 数据初始化: `internal/model/initData.go`（通过 `POST /system/init-data` 触发）
- 软删除: `DeletedAt gorm.DeletedAt`
- 多对多: `gorm:"many2many:table_name;"`

### ServiceContext
`internal/svc/servicecontext.go` 初始化所有 Model，新增 Model 需在此注册。

---

## 前端开发规范

### API 代理
```js
// vite.config.js - /api 前缀会被重写移除
'/api': { target: 'http://localhost:8001', rewrite: (path) => path.replace(/^\/api/, '') }
```
例: 前端请求 `/api/user/list` → 后端收到 `/user/list`

### Token 存储
- `localStorage.getItem('token')` - 访问令牌
- `localStorage.getItem('user')` - 用户信息（JSON 字符串）

### 主题色
Element Plus 主题色: `#2c7d34`（绿色）

### 权限控制
- 指令: `v-has-permission="'btn_user_create'"`
- 编程: `permission.check('btn_user_create')`

---

## 数据库表清单

| 表名 | 说明 |
|------|------|
| users / roles / permissions / menus | RBAC 核心表 |
| user_roles / role_permissions / role_menus / menu_permissions | 关联表 |
| activities | 操作日志 |
| categories / products / product_units | 产品及分类/单位 |
| suppliers / customers / warehouses | 供应商/客户/仓库 |
| purchase_orders / purchase_order_items | 采购订单 |
| sales_orders / sales_order_items | 销售订单 |
| inventory_records / inventory_changes / warehouse_inventories | 库存记录/变更/缓存 |
| inventory_adjust_requests | 库存调整申请 |
| inventory_checks / inventory_check_items | 库存盘点 |
| inventory_transfers | 库存调拨 |

---

## 安全注意

- JWT 密钥默认 `your-secret-key`，生产需修改
- CORS 当前允许所有来源 (`*`)，生产需限制
- 密码使用 bcrypt 加密存储
