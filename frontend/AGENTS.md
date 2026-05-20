# 前端管理系统 AGENTS.md

## 项目概述

基于 Vue 3 + Element Plus + Pinia + Vue Router 的单页应用（SPA），使用 Vite 作为构建工具。

## 常用命令

```bash
cd frontend
npm install              # 安装依赖
npm run dev             # 启动开发服务器（端口 3000）
npm run build           # 构建生产版本到 dist/ 目录
```

## 项目结构

```
frontend/src/
├── api/                # Axios API 模块（按业务模块划分）
│   ├── request.js      # 统一请求封装（拦截器配置）
│   ├── user.js         # 用户管理 API
│   ├── role.js         # 角色管理 API
│   ├── permission.js   # 权限管理 API
│   ├── menu.js         # 菜单管理 API
│   └── activity.js     # 活动日志 API
├── components/         # 公共组件
│   └── Layout.vue      # 布局组件（侧边栏、顶部导航）
├── router/             # Vue Router 配置
│   └── index.js        # 路由定义与守卫
├── store/              # Pinia 状态管理
│   ├── user.js         # 用户状态（token、权限、用户信息）
│   ├── menu.js         # 菜单状态（菜单树）
│   └── activity.js     # 活动日志状态
├── views/              # 页面组件
│   ├── login/          # 登录页面
│   ├── dashboard/      # 控制台首页
│   ├── user/           # 用户管理
│   ├── role/           # 角色管理
│   ├── permission/     # 权限管理
│   ├── menu/           # 菜单管理
│   └── activity/       # 活动日志
├── directives/         # 自定义指令
│   └── permission.js   # 权限指令（v-has-permission）
├── utils/              # 工具函数
│   └── permission.js   # 权限检查工具
├── App.vue             # 根组件
└── main.js             # 入口文件（Axios 拦截器、Element Plus 配置）
```

## 开发规范

### API 调用
- 所有 API 请求使用 `/api` 前缀
- Vite 代理会自动将 `/api` 转发到后端 `localhost:8000`

### Token 存储
- `localStorage.getItem('token')` - 访问令牌
- `localStorage.getItem('user')` - 用户信息（JSON 字符串）

### 主题色
- Element Plus 主题色: `#2c7d34`（绿色）

### 权限控制
- **指令方式**: `<button v-has-permission="'btn_user_create'">创建</button>`
- **编程方式**: `permission.check('btn_user_create')`

## 注意事项（避免做法）
- **不要硬编码后端 URL** - 使用 Vite 代理配置
- **不要将 token 存储在 cookies 中**（使用 localStorage）
- **不要在请求拦截器中使用 Pinia store**（直接从 localStorage 读取）

## 关键文件说明

| 文件路径 | 用途 |
|----------|------|
| `main.js` | 入口文件，配置 Axios 拦截器、Element Plus |
| `vite.config.js` | Vite 配置（代理 `/api` → `localhost:8000`） |
| `router/index.js` | 路由配置与导航守卫 |
| `store/user.js` | 用户状态管理（token、权限、用户信息） |
| `store/menu.js` | 菜单状态管理（菜单树） |
| `api/request.js` | 统一请求封装（请求/响应拦截器） |
| `directives/permission.js` | 权限指令 `v-has-permission` |
| `utils/permission.js` | 权限检查工具函数 |

## 认证流程

1. 用户通过 `/auth/login` 登录
2. Token 和用户信息存储到 localStorage
3. 每个请求通过 Axios 拦截器添加 `Authorization: Bearer <token>` 头
4. 遇到 401 错误时，用户重定向到登录页面

## 响应格式

后端统一响应格式：
```json
{ "code": 0, "data": {}, "message": "success" }
```

- `code: 0` = 成功，非 0 = 失败
- Axios 拦截器在 `src/api/request.js` 中处理响应

## 权限编码规范

权限 code 格式：`{模块}_{操作}`

```
示例:
- btn_user_create    # 创建用户
- btn_user_update    # 更新用户
- btn_role_assign    # 分配角色权限
- btn_menu_delete    # 删除菜单
```