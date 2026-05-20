# Frontend Application

## Overview

Vue 3 SPA with Element Plus + Pinia + Vue Router. Built with Vite.

## Technology Stack

- **Framework**: Vue 3
- **State Management**: Pinia
- **Router**: Vue Router
- **UI Library**: Element Plus
- **Build Tool**: Vite

## Project Structure

```
frontend/src/
├── api/                # API request wrappers
├── components/         # Shared components
├── router/             # Router configuration
├── store/              # Pinia stores
├── views/              # Page components
│   ├── login/          # Login page
│   ├── dashboard/      # Dashboard
│   ├── user/           # User management
│   ├── role/           # Role management
│   ├── permission/     # Permission management
│   ├── menu/           # Menu management
│   └── activity/       # Activity log
├── directives/         # Custom directives
├── utils/              # Utility functions
├── App.vue             # Root component
└── main.js             # Entry point
```

## Quick Start

### Requirements

- Node.js 16+

### Install Dependencies

```bash
npm install
```

### Development Mode

```bash
npm run dev
```

Dev server will start at `http://localhost:3000`.

### Production Build

```bash
npm run build
```

Build output will be in `dist/` directory.

## Core Features

| Module | Features |
|--------|----------|
| Authentication | Login, Logout, Token Refresh |
| User Management | Create, Edit, Delete, Role Assignment |
| Role Management | Create, Edit, Delete, Permission Assignment |
| Permission Management | Create, Edit, Delete, Tree Display |
| Menu Management | Create, Edit, Delete, Tree Structure |
| Activity Log | Operation Records |

## Permission Control

### Directive Usage
```vue
<button v-has-permission="'btn_user_create'">Create User</button>
```

### Programmatic Usage
```javascript
import { permission } from './utils/permission'
if (permission.check('btn_user_create')) {
  // Has permission
}
```

## Authentication Flow

1. After login, token and user info stored in localStorage
2. Each request adds `Authorization: Bearer <token>` via Axios interceptor
3. On 401 error, redirect to login page automatically

## Proxy Configuration

Vite proxies `/api` requests to backend:

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

## Default Account

- **Username**: admin
- **Password**: admin123