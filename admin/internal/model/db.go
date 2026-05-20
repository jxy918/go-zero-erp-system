package model

import (
	"fmt"
	"log"
	"strings"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// InitDB 初始化数据库连接
func InitDB(dataSource string) error {
	var err error

	// 从数据源字符串中提取数据库名称
	dbName := extractDBName(dataSource)

	// 先连接到 MySQL 服务器（不指定数据库）
	rootDataSource := strings.Replace(dataSource, "/"+dbName+"?", "/?", 1)

	// 配置GORM - 使用Silent模式减少日志输出
	config := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	}

	// 连接到 MySQL 服务器
	rootDB, err := gorm.Open(mysql.Open(rootDataSource), config)
	if err != nil {
		return fmt.Errorf("failed to connect to MySQL server: %v", err)
	}

	// 创建数据库（如果不存在）
	if err := createDatabase(rootDB, dbName); err != nil {
		return fmt.Errorf("failed to create database: %v", err)
	}

	// 关闭临时连接
	sqlDB, err := rootDB.DB()
	if err != nil {
		return fmt.Errorf("failed to get database connection: %v", err)
	}
	sqlDB.Close()

	// 连接到目标数据库
	// 确保使用utf8mb4字符集，解决中文乱码问题
	// 构建新的数据源字符串，确保字符集参数正确
	// 格式: user:pass@tcp(host:port)/dbname?charset=utf8mb4&parseTime=True&loc=Local
	// 从原始数据源字符串中提取基本信息
	baseDSN := dataSource
	// 移除所有查询参数
	if strings.Contains(baseDSN, "?") {
		baseDSN = strings.Split(baseDSN, "?")[0]
	}
	// 构建新的DSN，明确指定所有必要的参数，确保所有字符集相关参数都设置为utf8mb4
	newDSN := baseDSN + "?charset=utf8mb4&parseTime=True&loc=Local&collation=utf8mb4_unicode_ci"

	// 配置MySQL驱动，明确指定字符集
	mysqlConfig := mysql.Config{
		DSN:                       newDSN,
		DefaultStringSize:         256,
		DisableDatetimePrecision:  true,
		DontSupportRenameIndex:    true,
		DontSupportRenameColumn:   true,
		SkipInitializeWithVersion: false,
	}

	// 连接到数据库
	DB, err = gorm.Open(mysql.New(mysqlConfig), config)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %v", err)
	}

	// 获取底层的sql.DB对象
	sqlDB, err = DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get database connection: %v", err)
	}

	// 测试连接
	if err := sqlDB.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %v", err)
	}

	// 执行SET NAMES语句，确保客户端和连接的字符集都是utf8mb4
	if err := DB.Exec("SET NAMES utf8mb4 COLLATE utf8mb4_unicode_ci").Error; err != nil {
		return fmt.Errorf("failed to set names: %v", err)
	}

	log.Println("✅ 数据库连接成功")

	// 执行数据库迁移 - 同步模型结构到数据库
	if err := migrateDatabase(); err != nil {
		log.Printf("⚠️ 数据库迁移警告: %v", err)
	}

	// 设置连接池参数
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// 快速检查表是否存在
	if err := quickCheckTables(); err != nil {
		log.Printf("⚠️ 表检查失败: %v", err)
		log.Println("💡 尝试手动创建缺失的表...")
		// 创建 inventory_changes 表
		if err := createInventoryChangesTable(); err != nil {
			log.Printf("⚠️ 创建 inventory_changes 表失败: %v", err)
		} else {
			log.Println("✅ inventory_changes 表创建成功")
		}
		// 创建 inventory_adjust_requests 表
		if err := createInventoryAdjustRequestsTable(); err != nil {
			log.Printf("⚠️ 创建 inventory_adjust_requests 表失败: %v", err)
		} else {
			log.Println("✅ inventory_adjust_requests 表创建成功")
		}
		// 创建 inventory_checks 表
		if err := createInventoryChecksTable(); err != nil {
			log.Printf("⚠️ 创建 inventory_checks 表失败: %v", err)
		} else {
			log.Println("✅ inventory_checks 表创建成功")
		}
		// 创建 inventory_check_items 表
		if err := createInventoryCheckItemsTable(); err != nil {
			log.Printf("⚠️ 创建 inventory_check_items 表失败: %v", err)
		} else {
			log.Println("✅ inventory_check_items 表创建成功")
		}
		// 创建 inventory_transfers 表
		if err := createInventoryTransfersTable(); err != nil {
			log.Printf("⚠️ 创建 inventory_transfers 表失败: %v", err)
		} else {
			log.Println("✅ inventory_transfers 表创建成功")
		}
	}

	log.Println("✅ 服务启动完成")
	log.Println("💡 如需初始化数据库数据，请调用 POST /system/init-data 接口")
	return nil
}

// quickCheckTables 快速检查表是否存在
func quickCheckTables() error {
	tables := []string{"users", "roles", "permissions", "menus", "user_roles", "role_permissions", "menu_permissions", "role_menus", "activities",
		"categories", "products", "suppliers", "customers", "warehouses",
		"purchase_orders", "purchase_order_items", "sales_orders", "sales_order_items",
		"inventory_records", "inventory_changes", "inventory_adjust_requests",
		"inventory_checks", "inventory_check_items", "inventory_transfers"}

	for _, table := range tables {
		var count int64
		if err := DB.Raw("SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = DATABASE() AND table_name = ?", table).Scan(&count).Error; err != nil {
			return err
		}
		if count == 0 {
			return fmt.Errorf("表 %s 不存在", table)
		}
	}

	return nil
}

// extractDBName 从数据源字符串中提取数据库名称
func extractDBName(dataSource string) string {
	// 从数据源字符串中提取数据库名称
	// 格式: user:pass@tcp(host:port)/dbname?charset=utf8mb4&parseTime=True&loc=Local
	parts := strings.Split(dataSource, "/")
	if len(parts) < 2 {
		return ""
	}

	dbPart := parts[len(parts)-1]
	dbName := strings.Split(dbPart, "?")[0]
	return dbName
}

// createDatabase 创建数据库（如果不存在）
func createDatabase(db *gorm.DB, dbName string) error {
	// 执行创建数据库的 SQL 语句
	createDBQuery := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s DEFAULT CHARSET utf8mb4 COLLATE utf8mb4_unicode_ci", dbName)
	return db.Exec(createDBQuery).Error
}

// createInventoryChangesTable 创建 inventory_changes 表
func createInventoryChangesTable() error {
	sql := `
		CREATE TABLE IF NOT EXISTS inventory_changes (
			id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
			product_id BIGINT UNSIGNED NOT NULL,
			warehouse_id BIGINT UNSIGNED NOT NULL,
			before_quantity INT NOT NULL,
			after_quantity INT NOT NULL,
			quantity INT NOT NULL,
			type INT NOT NULL,
			order_id BIGINT UNSIGNED,
			order_type VARCHAR(50),
			remark VARCHAR(500),
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			INDEX idx_product_id (product_id),
			INDEX idx_warehouse_id (warehouse_id),
			INDEX idx_type (type),
			INDEX idx_order_id (order_id)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
	`
	return DB.Exec(sql).Error
}

// createInventoryAdjustRequestsTable 创建库存调整申请表
func createInventoryAdjustRequestsTable() error {
	sql := `
		CREATE TABLE IF NOT EXISTS inventory_adjust_requests (
			id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
			request_no VARCHAR(50) NOT NULL UNIQUE,
			product_id BIGINT UNSIGNED NOT NULL,
			warehouse_id BIGINT UNSIGNED NOT NULL,
			before_qty INT NOT NULL,
			quantity INT NOT NULL,
			after_qty INT NOT NULL,
			type INT NOT NULL,
			reason VARCHAR(500),
			status INT DEFAULT 1,
			applicant_id BIGINT UNSIGNED NOT NULL,
			approver_id BIGINT UNSIGNED,
			approve_time DATETIME,
			approve_note VARCHAR(500),
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			deleted_at DATETIME,
			INDEX idx_product_id (product_id),
			INDEX idx_warehouse_id (warehouse_id),
			INDEX idx_status (status),
			INDEX idx_applicant_id (applicant_id),
			INDEX idx_request_no (request_no)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
	`
	return DB.Exec(sql).Error
}

// createInventoryChecksTable 创建盘点单主表
func createInventoryChecksTable() error {
	sql := `
		CREATE TABLE IF NOT EXISTS inventory_checks (
			id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
			check_no VARCHAR(50) NOT NULL UNIQUE,
			warehouse_id BIGINT UNSIGNED NOT NULL,
			status INT DEFAULT 1,
			total_diff INT DEFAULT 0,
			remark VARCHAR(500),
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			deleted_at DATETIME,
			INDEX idx_warehouse_id (warehouse_id),
			INDEX idx_status (status),
			INDEX idx_deleted_at (deleted_at)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
	`
	return DB.Exec(sql).Error
}

// createInventoryCheckItemsTable 创建盘点单明细表
func createInventoryCheckItemsTable() error {
	sql := `
		CREATE TABLE IF NOT EXISTS inventory_check_items (
			id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
			check_id BIGINT UNSIGNED NOT NULL,
			product_id BIGINT UNSIGNED NOT NULL,
			system_qty INT NOT NULL,
			actual_qty INT NOT NULL,
			diff_qty INT NOT NULL,
			status INT DEFAULT 1,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			deleted_at DATETIME,
			INDEX idx_check_id (check_id),
			INDEX idx_product_id (product_id),
			INDEX idx_status (status),
			INDEX idx_deleted_at (deleted_at)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
	`
	return DB.Exec(sql).Error
}

// createInventoryTransfersTable 创建调拨单表
func createInventoryTransfersTable() error {
	sql := `
		CREATE TABLE IF NOT EXISTS inventory_transfers (
			id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
			transfer_no VARCHAR(50) NOT NULL UNIQUE,
			from_warehouse_id BIGINT UNSIGNED NOT NULL,
			to_warehouse_id BIGINT UNSIGNED NOT NULL,
			product_id BIGINT UNSIGNED NOT NULL,
			quantity INT NOT NULL,
			status TINYINT NOT NULL DEFAULT 1,
			remark VARCHAR(500),
			created_by BIGINT UNSIGNED,
			audited_by BIGINT UNSIGNED,
			audited_at DATETIME,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			deleted_at DATETIME,
			INDEX idx_transfer_no (transfer_no),
			INDEX idx_from_warehouse_id (from_warehouse_id),
			INDEX idx_to_warehouse_id (to_warehouse_id),
			INDEX idx_product_id (product_id),
			INDEX idx_status (status),
			INDEX idx_created_at (created_at)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
	`
	return DB.Exec(sql).Error
}

// migrateDatabase 执行数据库迁移
func migrateDatabase() error {
	log.Println("🔄 开始执行数据库迁移...")

	// 0. 使用 GORM AutoMigrate 自动同步所有核心模型结构
	log.Println("📦 执行 AutoMigrate 同步核心表结构...")
	if err := DB.AutoMigrate(
		&User{}, &Role{}, &Permission{}, &Menu{},
		&UserRole{}, &RolePermission{}, &RoleMenu{}, &MenuPermission{},
		&Activity{},
		&Category{}, &Product{}, &ProductUnit{},
		&Supplier{}, &Customer{}, &Warehouse{},
		&PurchaseOrder{}, &PurchaseOrderItem{},
		&SalesOrder{}, &SalesOrderItem{},
		&InventoryRecord{}, &WarehouseInventory{}, &InventoryChange{},
		&InventoryAdjustRequest{},
	); err != nil {
		log.Printf("⚠️ AutoMigrate 失败: %v", err)
	} else {
		log.Println("✅ AutoMigrate 核心表同步完成")
	}

	// 1. 删除 products.stock 字段（如果存在）
	if err := dropColumnIfExists("products", "stock"); err != nil {
		log.Printf("⚠️ 删除 products.stock 字段: %v", err)
	} else {
		log.Println("✅ 删除 products.stock 字段成功")
	}

	// 1. 为 products 表添加库存预警字段（如果不存在）
	if err := addColumnIfNotExists("products", "min_stock", "INT DEFAULT 0"); err != nil {
		log.Printf("⚠️ 添加 min_stock 字段失败: %v", err)
	} else {
		log.Println("✅ 添加 min_stock 字段成功")
	}

	if err := addColumnIfNotExists("products", "max_stock", "INT DEFAULT 99999"); err != nil {
		log.Printf("⚠️ 添加 max_stock 字段失败: %v", err)
	} else {
		log.Println("✅ 添加 max_stock 字段成功")
	}

	if err := addColumnIfNotExists("products", "safety_stock", "INT DEFAULT 0"); err != nil {
		log.Printf("⚠️ 添加 safety_stock 字段失败: %v", err)
		log.Println("💡 safety_stock 字段缺失不影响基本功能，可以手动添加")
	} else {
		log.Println("✅ 添加 safety_stock 字段成功")
	}

	// 2. 创建 product_units 表（如果不存在）
	if err := createProductUnitsTable(); err != nil {
		log.Printf("⚠️ 创建 product_units 表失败: %v", err)
	} else {
		log.Println("✅ 创建 product_units 表成功")
	}

	// 3. 创建 warehouse_inventories 表（如果不存在）
	if err := createWarehouseInventoriesTable(); err != nil {
		log.Printf("⚠️ 创建 warehouse_inventories 表失败: %v", err)
	} else {
		log.Println("✅ 创建 warehouse_inventories 表成功")
	}

	// 4. 为采购订单明细表添加单位相关字段
	if err := addColumnIfNotExists("purchase_order_items", "unit_id", "INT DEFAULT 0"); err != nil {
		log.Printf("⚠️ 添加 purchase_order_items.unit_id 字段失败: %v", err)
	} else {
		log.Println("✅ 添加 purchase_order_items.unit_id 字段成功")
	}

	if err := addColumnIfNotExists("purchase_order_items", "unit_name", "VARCHAR(50)"); err != nil {
		log.Printf("⚠️ 添加 purchase_order_items.unit_name 字段失败: %v", err)
	} else {
		log.Println("✅ 添加 purchase_order_items.unit_name 字段成功")
	}

	if err := addColumnIfNotExists("purchase_order_items", "ratio", "INT DEFAULT 1"); err != nil {
		log.Printf("⚠️ 添加 purchase_order_items.ratio 字段失败: %v", err)
	} else {
		log.Println("✅ 添加 purchase_order_items.ratio 字段成功")
	}

	if err := addColumnIfNotExists("purchase_order_items", "base_qty", "INT DEFAULT 0"); err != nil {
		log.Printf("⚠️ 添加 purchase_order_items.base_qty 字段失败: %v", err)
	} else {
		log.Println("✅ 添加 purchase_order_items.base_qty 字段成功")
	}

	// 4. 为销售订单明细表添加单位相关字段
	if err := addColumnIfNotExists("sales_order_items", "unit_id", "INT DEFAULT 0"); err != nil {
		log.Printf("⚠️ 添加 sales_order_items.unit_id 字段失败: %v", err)
	} else {
		log.Println("✅ 添加 sales_order_items.unit_id 字段成功")
	}

	if err := addColumnIfNotExists("sales_order_items", "unit_name", "VARCHAR(50)"); err != nil {
		log.Printf("⚠️ 添加 sales_order_items.unit_name 字段失败: %v", err)
	} else {
		log.Println("✅ 添加 sales_order_items.unit_name 字段成功")
	}

	if err := addColumnIfNotExists("sales_order_items", "ratio", "INT DEFAULT 1"); err != nil {
		log.Printf("⚠️ 添加 sales_order_items.ratio 字段失败: %v", err)
	} else {
		log.Println("✅ 添加 sales_order_items.ratio 字段成功")
	}

	if err := addColumnIfNotExists("sales_order_items", "base_qty", "INT DEFAULT 0"); err != nil {
		log.Printf("⚠️ 添加 sales_order_items.base_qty 字段失败: %v", err)
	} else {
		log.Println("✅ 添加 sales_order_items.base_qty 字段成功")
	}

	log.Println("✅ 数据库迁移完成")
	return nil
}

// addColumnIfNotExists 如果字段不存在则添加字段
func addColumnIfNotExists(tableName, columnName, columnType string) error {
	var count int64
	checkSQL := `SELECT COUNT(*) FROM information_schema.columns WHERE table_schema = DATABASE() AND table_name = ? AND column_name = ?`
	if err := DB.Raw(checkSQL, tableName, columnName).Scan(&count).Error; err != nil {
		log.Printf("检查字段 %s.%s 失败: %v", tableName, columnName, err)
		return err
	}

	if count > 0 {
		log.Printf("字段 %s.%s 已存在，跳过添加", tableName, columnName)
		return nil
	}

	alterSQL := fmt.Sprintf("ALTER TABLE `%s` ADD COLUMN `%s` %s", tableName, columnName, columnType)
	log.Printf("执行SQL: %s", alterSQL)
	return DB.Exec(alterSQL).Error
}

// dropColumnIfExists 删除字段（如果存在）
func dropColumnIfExists(tableName, columnName string) error {
	var count int64
	checkSQL := `SELECT COUNT(*) FROM information_schema.columns WHERE table_schema = DATABASE() AND table_name = ? AND column_name = ?`
	if err := DB.Raw(checkSQL, tableName, columnName).Scan(&count).Error; err != nil {
		return err
	}

	if count == 0 {
		log.Printf("字段 %s.%s 不存在，跳过删除", tableName, columnName)
		return nil
	}

	dropSQL := fmt.Sprintf("ALTER TABLE `%s` DROP COLUMN `%s`", tableName, columnName)
	log.Printf("执行SQL: %s", dropSQL)
	return DB.Exec(dropSQL).Error
}

// createWarehouseInventoriesTable 创建仓库库存缓存表
func createWarehouseInventoriesTable() error {
	var count int64
	if err := DB.Raw("SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = DATABASE() AND table_name = 'warehouse_inventories'").Scan(&count).Error; err != nil {
		return err
	}

	if count > 0 {
		return nil
	}

	sql := `
		CREATE TABLE IF NOT EXISTS warehouse_inventories (
			id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
			product_id BIGINT UNSIGNED NOT NULL,
			warehouse_id BIGINT UNSIGNED NOT NULL,
			quantity INT DEFAULT 0,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			UNIQUE KEY uk_product_warehouse (product_id, warehouse_id),
			INDEX idx_product_id (product_id),
			INDEX idx_warehouse_id (warehouse_id)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
	`
	return DB.Exec(sql).Error
}

// createProductUnitsTable 创建产品计量单位表
func createProductUnitsTable() error {
	var count int64
	if err := DB.Raw("SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = DATABASE() AND table_name = 'product_units'").Scan(&count).Error; err != nil {
		return err
	}

	if count > 0 {
		return nil
	}

	sql := `
		CREATE TABLE IF NOT EXISTS product_units (
			id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
			product_id BIGINT UNSIGNED NOT NULL,
			unit_name VARCHAR(20) NOT NULL,
			ratio DECIMAL(10,4) DEFAULT 1,
			is_main TINYINT DEFAULT 0,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			deleted_at DATETIME,
			INDEX idx_product_id (product_id),
			INDEX idx_deleted_at (deleted_at),
			UNIQUE KEY uk_product_unit (product_id, unit_name)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
	`
	return DB.Exec(sql).Error
}
