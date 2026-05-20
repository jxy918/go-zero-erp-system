# 前端应用

## 概述

基于 Vue 3 + Element Plus + Pinia + Vue Router 的单页应用（SPA），使用 Vite 作为构建工具。

## 技术栈

- **框架**: Vue 3
- **状态管理**: Pinia
- **路由**: Vue Router
- **UI组件**: Element Plus
- **构建工具**: Vite

## 项目结构

```
frontend/src/
├── api/                # API 请求封装
├── components/         # 公共组件
├── router/             # 路由配置
├── store/              # Pinia 状态管理
├── views/              # 页面组件
│   ├── login/          # 登录页面
│   ├── dashboard/      # 控制台
│   ├── user/           # 用户管理
│   ├── role/           # 角色管理
│   ├── permission/     # 权限管理
│   ├── menu/           # 菜单管理
│   └── activity/       # 活动日志
├── directives/         # 自定义指令
├── utils/              # 工具函数
├── App.vue             # 根组件
└── main.js             # 入口文件
```

## 快速开始

### 环境要求

- Node.js 16+

### 安装依赖

```bash
npm install
```

### 开发模式

```bash
npm run dev
```

开发服务器将在 `http://localhost:3000` 启动。

### 生产构建

```bash
npm run build
```

构建产物将输出到 `dist/` 目录。

## 核心功能

| 模块 | 功能 |
|------|------|
| 用户认证 | 登录、登出、Token刷新 |
| 用户管理 | 创建、编辑、删除、角色分配 |
| 角色管理 | 创建、编辑、删除、权限分配 |
| 权限管理 | 创建、编辑、删除、树形展示 |
| 菜单管理 | 创建、编辑、删除、树形结构 |
| 活动日志 | 操作记录查询 |

## 权限控制

### 指令方式
```vue
<button v-has-permission="'btn_user_create'">创建用户</button>
```

### 编程方式
```javascript
import { permission } from './utils/permission'
if (permission.check('btn_user_create')) {
  // 有权限
}
```

## 认证流程

1. 用户登录成功后，Token 和用户信息存储到 localStorage
2. 每个请求通过 Axios 拦截器添加 `Authorization: Bearer <token>` 头
3. 遇到 401 错误时，自动跳转到登录页面

## 代理配置

Vite 配置将 `/api` 前缀转发到后端：

```javascript
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

## 默认账号

- **用户名**: admin
- **密码**: admin123