# RBAC Admin System + ERP Management System

## Project Introduction

This is an enterprise-level backend management system based on Role-Based Access Control (RBAC), integrating complete ERP management functions with a separated front-end and back-end architecture. The system provides comprehensive user authentication, role management, permission control, menu management, and ERP core functionality modules.

## System Preview

### Login Page
![Login Page](images/image1.png)

### Dashboard
![Dashboard](images/image2.png)

### ERP Statistics
![ERP Statistics](images/image3.png)

## Technology Stack

### Frontend
| Component | Technology | Version |
|-----------|------------|---------|
| Framework | Vue | 3.x |
| State Management | Pinia | - |
| Router | Vue Router | - |
| UI Library | Element Plus | - |
| HTTP Client | Axios | - |
| Build Tool | Vite | - |
| Chart Library | ECharts | - |

### Backend
| Component | Technology | Version |
|-----------|------------|---------|
| Framework | Go-zero | - |
| Language | Go | 1.20+ |
| Database | MySQL | 8.0+ |
| ORM | GORM | v1 |
| Authentication | JWT | - |

### Service Ports
| Service | Port |
|---------|------|
| Backend API | 8000 |
| Frontend Dev Server | 3000 |

## Core Features

### System Management Module
| Module | Features |
|--------|----------|
| User Authentication | Login, Logout, Token Refresh |
| User Management | Create, Edit, Delete, Role Assignment |
| Role Management | Create, Edit, Delete, Permission Assignment, Menu Assignment |
| Permission Management | Create, Edit, Delete, Tree Display |
| Menu Management | Create, Edit, Delete, Tree Structure |
| Activity Log | Operation Records, User Filtering |
| Data Permission | Role-based Data Access Control |

### ERP Management Module
| Module | Features |
|--------|----------|
| Product Management | Product List, Create, Edit, Delete, Category Management |
| Supplier Management | Supplier List, Create, Edit, Delete |
| Customer Management | Customer List, Create, Edit, Delete |
| Warehouse Management | Warehouse List, Create, Edit, Delete |
| Purchase Management | Purchase Orders, Approval, Inbound |
| Sales Management | Sales Orders, Approval, Outbound |
| Inventory Management | Inventory List, Adjust Request, Alert, Transfer, Check |
| Statistics Report | Purchase/Sales Trend, Inventory Alert, Hot Products |

## Project Structure

### Backend Structure
```
go-zero-erp/
├── admin/                    # Backend Application
│   ├── etc/                  # Configuration Files
│   │   └── admin-api.yaml    # Main Configuration
│   ├── internal/             # Internal Code
│   │   ├── config/           # Configuration Structs
│   │   ├── handler/          # HTTP Handlers
│   │   ├── logic/            # Business Logic
│   │   ├── middleware/       # Middleware (Auth, CORS, Data Permission)
│   │   ├── model/            # GORM Models
│   │   ├── svc/              # Service Context
│   │   ├── types/            # Request/Response Types
│   │   └── util/             # Utility Functions (JWT, IP)
│   ├── admin.api             # API Definition (goctl)
│   └── admin.go              # Application Entry
├── frontend/                 # Frontend Application
├── docs/                     # Project Documentation
├── images/                   # Screenshots
├── test/                     # Test Scripts
├── go.mod                    # Go Module File
└── go.sum                    # Dependency Checksum File
```

### Frontend Structure
```
frontend/
├── src/                      # Source Code
│   ├── api/                  # API Request Wrappers
│   ├── components/           # Common Components
│   ├── directives/           # Custom Directives (Permission)
│   ├── router/               # Router Configuration
│   ├── store/                # Pinia State Management
│   ├── utils/                # Utility Functions
│   ├── views/                # Page Components
│   │   ├── login/            # Login Page
│   │   ├── dashboard/        # Dashboard
│   │   ├── user/             # User Management
│   │   ├── role/             # Role Management
│   │   ├── permission/       # Permission Management
│   │   ├── menu/             # Menu Management
│   │   ├── activity/         # Activity Log
│   │   ├── product/          # Product Management
│   │   ├── supplier/         # Supplier Management
│   │   ├── customer/         # Customer Management
│   │   ├── warehouse/        # Warehouse Management
│   │   ├── purchase/         # Purchase Management
│   │   ├── sales/            # Sales Management
│   │   ├── inventory/        # Inventory Management
│   │   └── erp/              # ERP Statistics Report
│   ├── App.vue               # Root Component
│   └── main.js               # Entry File
├── index.html                # HTML Template
├── vite.config.js            # Vite Configuration (with Proxy)
└── package.json              # Dependency Management
```

## Installation & Running

### Environment Requirements
1. Go 1.20+
2. Node.js 16+
3. MySQL 8.0+

### Start Backend
```bash
cd admin
go mod tidy                    # Install dependencies
go run admin.go -f etc/admin-api.yaml  # Start server (port 8000)
```

### Start Frontend
```bash
cd frontend
npm install                    # Install dependencies
npm run dev                    # Start development server (port 3000)
npm run build                  # Build for production
```

### Database Configuration
Modify `admin/etc/admin-api.yaml`:
```yaml
DataSource: root:root@tcp(127.0.0.1:3306)/admin_system?charset=utf8mb4&parseTime=True&loc=Local
```

### Initialize Data
Table structures are automatically created on startup. Initialize basic data via API:
```
POST /system/init-data
```

## API Endpoints

### Authentication
| Method | Path | Description |
|--------|------|-------------|
| POST | /auth/login | User Login |
| POST | /auth/logout | User Logout |
| POST | /auth/refresh | Refresh Token |

### User Management
| Method | Path | Description |
|--------|------|-------------|
| GET | /user/list | Get User List |
| GET | /user/get/:id | Get User Detail |
| POST | /user/create | Create User |
| POST | /user/update | Update User |
| POST | /user/delete | Delete User |
| POST | /user/assign-roles | Assign Roles |

### Role Management
| Method | Path | Description |
|--------|------|-------------|
| GET | /role/list | Get Role List |
| GET | /role/get/:id | Get Role Detail |
| POST | /role/create | Create Role |
| PUT | /role/update | Update Role |
| DELETE | /role/delete | Delete Role |
| POST | /role/assign-permissions | Assign Permissions |
| POST | /role/assign-menus | Assign Menus |

### Permission Management
| Method | Path | Description |
|--------|------|-------------|
| GET | /permission/list | Get Permission List |
| GET | /permission/get/:id | Get Permission Detail |
| POST | /permission/create | Create Permission |
| PUT | /permission/update | Update Permission |
| DELETE | /permission/delete | Delete Permission |

### Menu Management
| Method | Path | Description |
|--------|------|-------------|
| GET | /menu/tree | Get Menu Tree |
| GET | /menu/list | Get Menu List |
| POST | /menu/create | Create Menu |
| POST | /menu/update | Update Menu |
| POST | /menu/delete | Delete Menu |
| GET | /menu/get/:id | Get Menu Detail |
| POST | /menu/assign-permissions | Assign Menu Permissions |

### Activity Log
| Method | Path | Description |
|--------|------|-------------|
| GET | /activity/list | Get Activity Log |

### Product Management
| Method | Path | Description |
|--------|------|-------------|
| GET | /product/list | Get Product List |
| GET | /product/get/:id | Get Product Detail |
| POST | /product/create | Create Product |
| POST | /product/update | Update Product |
| POST | /product/delete | Delete Product |
| GET | /product/category/list | Get Category List |
| POST | /product/category/create | Create Category |

### Supplier Management
| Method | Path | Description |
|--------|------|-------------|
| GET | /supplier/list | Get Supplier List |
| GET | /supplier/get/:id | Get Supplier Detail |
| POST | /supplier/create | Create Supplier |
| POST | /supplier/update | Update Supplier |
| POST | /supplier/delete | Delete Supplier |

### Customer Management
| Method | Path | Description |
|--------|------|-------------|
| GET | /customer/list | Get Customer List |
| GET | /customer/get/:id | Get Customer Detail |
| POST | /customer/create | Create Customer |
| POST | /customer/update | Update Customer |
| POST | /customer/delete | Delete Customer |

### Warehouse Management
| Method | Path | Description |
|--------|------|-------------|
| GET | /warehouse/list | Get Warehouse List |
| GET | /warehouse/get/:id | Get Warehouse Detail |
| POST | /warehouse/create | Create Warehouse |
| POST | /warehouse/update | Update Warehouse |
| POST | /warehouse/delete | Delete Warehouse |

### Purchase Management
| Method | Path | Description |
|--------|------|-------------|
| GET | /purchase/list | Get Purchase Order List |
| GET | /purchase/get/:id | Get Purchase Order Detail |
| POST | /purchase/create | Create Purchase Order |
| POST | /purchase/update | Update Purchase Order |
| POST | /purchase/delete | Delete Purchase Order |
| POST | /purchase/approve | Approve Purchase Order |
| POST | /purchase/inbound | Purchase Inbound |

### Sales Management
| Method | Path | Description |
|--------|------|-------------|
| GET | /sales/list | Get Sales Order List |
| GET | /sales/get/:id | Get Sales Order Detail |
| POST | /sales/create | Create Sales Order |
| POST | /sales/update | Update Sales Order |
| POST | /sales/delete | Delete Sales Order |
| POST | /sales/approve | Approve Sales Order |
| POST | /sales/outbound | Sales Outbound |

### Inventory Management
| Method | Path | Description |
|--------|------|-------------|
| GET | /inventory/list | Get Inventory List |
| GET | /inventory/history | Get Inventory History |
| GET | /inventory/current-stock | Get Current Stock |
| POST | /inventory/adjust-request/create | Create Inventory Adjust Request |
| GET | /inventory/adjust-request/list | Get Adjust Request List |
| POST | /inventory/adjust-request/approve | Approve Adjust Request |
| POST | /inventory/adjust-request/reject | Reject Adjust Request |
| POST | /inventory/transfer/create | Create Inventory Transfer |
| POST | /inventory/check/create | Create Inventory Check |

### ERP Statistics
| Method | Path | Description |
|--------|------|-------------|
| GET | /erp/statistics/overview | Get Overview Statistics |
| GET | /erp/statistics/trend | Get Purchase/Sales Trend |
| GET | /erp/statistics/inventory-alert | Get Inventory Alert |
| GET | /erp/statistics/top-products | Get Top Products |
| GET | /erp/statistics/order-status | Get Order Status |
| GET | /erp/statistics/business | Get Business Statistics |

### System
| Method | Path | Description |
|--------|------|-------------|
| POST | /system/init-data | Initialize Data |

## Default Account

| Username | Password | Role |
|----------|----------|------|
| admin | admin123 | Super Administrator |

## Permission Description

### Permission Code Format
```
{module}:{action}

Examples:
- btn_user_create    # Create User
- btn_user_update    # Update User
- btn_role_assign    # Assign Role Permissions
```

### Data Permission
- **Admin**: Can view all data
- **Normal User**: Can only view data within their role scope

## Security Notes

1. JWT secret key defaults to `your-secret-key`, please change it in production
2. CORS currently allows all origins, please restrict in production
3. Passwords are stored with bcrypt encryption
4. HTTPS is recommended for production environments

## License

MIT License