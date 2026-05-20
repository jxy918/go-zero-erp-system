# Backend Service

## Overview

Go backend service using go-zero + GORM + MySQL. Provides REST API for RBAC admin system and ERP management system.

## Technology Stack

- **Framework**: go-zero
- **Language**: Go 1.20+
- **Database**: MySQL 8.0+
- **ORM**: GORM v1
- **Authentication**: JWT

## Project Structure

```
admin/
├── admin.go              # Entry point
├── admin.api             # API definition
├── etc/
│   └── admin-api.yaml    # Configuration
└── internal/
    ├── config/           # Config structs
    ├── handler/          # HTTP handlers
    ├── logic/            # Business logic
    ├── middleware/       # Middleware
    ├── model/            # Data models
    ├── svc/              # Service context
    ├── types/            # Type definitions
    └── util/             # Utility functions
```

## Quick Start

### Requirements

- Go 1.20+
- MySQL 8.0+

### Install Dependencies

```bash
go mod tidy
```

### Configure Database

Modify `etc/admin-api.yaml`:

```yaml
DataSource: root:root@tcp(127.0.0.1:3306)/admin_system?charset=utf8mb4&parseTime=True&loc=Local
```

### Start Server

```bash
go run admin.go -f etc/admin-api.yaml
```

Server will start at `http://localhost:8000`.

## API Endpoints

### Authentication
- `POST /auth/login` - User login
- `POST /auth/logout` - User logout
- `POST /auth/refresh` - Refresh token

### User Management
- `GET /user/list` - Get user list
- `GET /user/get/:id` - Get user detail
- `POST /user/create` - Create user
- `POST /user/update` - Update user
- `POST /user/delete` - Delete user
- `POST /user/assign-roles` - Assign roles

### Role Management
- `GET /role/list` - Get role list
- `GET /role/get/:id` - Get role detail
- `POST /role/create` - Create role
- `PUT /role/update` - Update role
- `DELETE /role/delete` - Delete role
- `POST /role/assign-permissions` - Assign permissions
- `POST /role/assign-menus` - Assign menus

### Permission Management
- `GET /permission/list` - Get permission list
- `GET /permission/get/:id` - Get permission detail
- `POST /permission/create` - Create permission
- `PUT /permission/update` - Update permission
- `DELETE /permission/delete` - Delete permission

### Menu Management
- `GET /menu/tree` - Get menu tree
- `GET /menu/list` - Get menu list
- `GET /menu/get/:id` - Get menu detail
- `POST /menu/create` - Create menu
- `POST /menu/update` - Update menu
- `POST /menu/delete` - Delete menu
- `POST /menu/assign-permissions` - Assign menu permissions

### Activity Log
- `GET /activity/list` - Get activity log

### Product Management
- `GET /product/list` - Get product list
- `GET /product/get/:id` - Get product detail
- `POST /product/create` - Create product
- `POST /product/update` - Update product
- `POST /product/delete` - Delete product
- `GET /product/category/list` - Get category list
- `POST /product/category/create` - Create category

### Supplier Management
- `GET /supplier/list` - Get supplier list
- `GET /supplier/get/:id` - Get supplier detail
- `POST /supplier/create` - Create supplier
- `POST /supplier/update` - Update supplier
- `POST /supplier/delete` - Delete supplier

### Customer Management
- `GET /customer/list` - Get customer list
- `GET /customer/get/:id` - Get customer detail
- `POST /customer/create` - Create customer
- `POST /customer/update` - Update customer
- `POST /customer/delete` - Delete customer

### Warehouse Management
- `GET /warehouse/list` - Get warehouse list
- `GET /warehouse/get/:id` - Get warehouse detail
- `POST /warehouse/create` - Create warehouse
- `POST /warehouse/update` - Update warehouse
- `POST /warehouse/delete` - Delete warehouse

### Purchase Management
- `GET /purchase/list` - Get purchase order list
- `GET /purchase/get/:id` - Get purchase order detail
- `POST /purchase/create` - Create purchase order
- `POST /purchase/update` - Update purchase order
- `POST /purchase/delete` - Delete purchase order
- `POST /purchase/approve` - Approve purchase order
- `POST /purchase/inbound` - Purchase inbound

### Sales Management
- `GET /sales/list` - Get sales order list
- `GET /sales/get/:id` - Get sales order detail
- `POST /sales/create` - Create sales order
- `POST /sales/update` - Update sales order
- `POST /sales/delete` - Delete sales order
- `POST /sales/approve` - Approve sales order
- `POST /sales/outbound` - Sales outbound

### Inventory Management
- `GET /inventory/list` - Get inventory list
- `GET /inventory/history` - Get inventory history
- `GET /inventory/current-stock` - Get current stock
- `POST /inventory/adjust-request/create` - Create inventory adjust request
- `GET /inventory/adjust-request/list` - Get adjust request list
- `POST /inventory/adjust-request/approve` - Approve adjust request
- `POST /inventory/adjust-request/reject` - Reject adjust request

### ERP Statistics
- `GET /erp/statistics/overview` - Get overview statistics
- `GET /erp/statistics/trend` - Get purchase/sales trend
- `GET /erp/statistics/inventory-alert` - Get inventory alert
- `GET /erp/statistics/top-products` - Get top products
- `GET /erp/statistics/order-status` - Get order status
- `GET /erp/statistics/business` - Get business statistics

### System
- `POST /system/init-data` - Initialize data

## Development Workflow

1. Edit `admin.api` to define API
2. Run `goctl api go -api admin.api -dir . -style goZero` to generate code
3. Implement business logic in `internal/logic/`

## Default Account

- **Username**: admin
- **Password**: admin123

## Security Notes

1. JWT secret defaults to `your-secret-key`, change in production
2. CORS allows all origins, restrict in production
3. HTTPS is recommended for production
