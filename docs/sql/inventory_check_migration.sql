-- ==============================================
-- 库存盘点相关表结构
-- inventory_checks: 盘点单主表
-- inventory_check_items: 盘点单明细表
-- ==============================================

-- ----------------------------
-- 1. 盘点单主表 (inventory_checks)
-- ----------------------------
DROP TABLE IF EXISTS `inventory_checks`;
CREATE TABLE `inventory_checks` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '盘点单ID',
  `check_no` varchar(50) NOT NULL COMMENT '盘点单号（格式：CK-YYYYMMDD-XXXX）',
  `warehouse_id` bigint unsigned NOT NULL COMMENT '仓库ID',
  `status` int DEFAULT '1' COMMENT '状态（1:进行中, 2:已完成）',
  `total_diff` int DEFAULT '0' COMMENT '总差异数量',
  `remark` varchar(500) DEFAULT '' COMMENT '备注',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` datetime DEFAULT NULL COMMENT '删除时间（软删除）',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_check_no` (`check_no`),
  KEY `idx_warehouse_id` (`warehouse_id`),
  KEY `idx_status` (`status`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='盘点单主表';

-- ----------------------------
-- 2. 盘点单明细表 (inventory_check_items)
-- ----------------------------
DROP TABLE IF EXISTS `inventory_check_items`;
CREATE TABLE `inventory_check_items` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '盘点明细ID',
  `check_id` bigint unsigned NOT NULL COMMENT '盘点单ID',
  `product_id` bigint unsigned NOT NULL COMMENT '产品ID',
  `system_qty` int NOT NULL COMMENT '系统库存数量',
  `actual_qty` int NOT NULL COMMENT '实际盘点数量',
  `diff_qty` int NOT NULL COMMENT '差异数量（实际-系统）',
  `status` int DEFAULT '1' COMMENT '状态（1:待处理, 2:已处理）',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` datetime DEFAULT NULL COMMENT '删除时间（软删除）',
  PRIMARY KEY (`id`),
  KEY `idx_check_id` (`check_id`),
  KEY `idx_product_id` (`product_id`),
  KEY `idx_status` (`status`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='盘点单明细表';

-- ----------------------------
-- 3. 仓库库存表 (warehouse_inventories) - 依赖表
-- ----------------------------
DROP TABLE IF EXISTS `warehouse_inventories`;
CREATE TABLE `warehouse_inventories` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '库存记录ID',
  `product_id` bigint unsigned NOT NULL COMMENT '产品ID',
  `warehouse_id` bigint unsigned NOT NULL COMMENT '仓库ID',
  `quantity` int NOT NULL DEFAULT '0' COMMENT '库存数量',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_product_warehouse` (`product_id`, `warehouse_id`),
  KEY `idx_product_id` (`product_id`),
  KEY `idx_warehouse_id` (`warehouse_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='仓库库存表';

-- ==============================================
-- 外键约束
-- ==============================================
ALTER TABLE `inventory_checks`
  ADD CONSTRAINT `fk_inventory_checks_warehouse` FOREIGN KEY (`warehouse_id`) REFERENCES `warehouses` (`id`) ON DELETE RESTRICT ON UPDATE CASCADE;

ALTER TABLE `inventory_check_items`
  ADD CONSTRAINT `fk_inventory_check_items_check` FOREIGN KEY (`check_id`) REFERENCES `inventory_checks` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  ADD CONSTRAINT `fk_inventory_check_items_product` FOREIGN KEY (`product_id`) REFERENCES `products` (`id`) ON DELETE RESTRICT ON UPDATE CASCADE;

ALTER TABLE `warehouse_inventories`
  ADD CONSTRAINT `fk_warehouse_inventories_product` FOREIGN KEY (`product_id`) REFERENCES `products` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  ADD CONSTRAINT `fk_warehouse_inventories_warehouse` FOREIGN KEY (`warehouse_id`) REFERENCES `warehouses` (`id`) ON DELETE CASCADE ON UPDATE CASCADE;

-- ==============================================
-- 初始化测试数据
-- ==============================================

-- 1. 创建测试仓库
INSERT INTO `warehouses` (`name`, `code`, `address`, `status`) VALUES
('主仓库', 'WH-MAIN', '北京市朝阳区某某路1号', 1),
('备货仓库', 'WH-BACKUP', '上海市浦东新区某某路2号', 1);

-- 2. 创建测试产品
INSERT INTO `products` (`name`, `code`, `category_id`, `spec`, `unit`, `price`, `cost_price`, `stock`, `status`) VALUES
('产品A', 'P-001', 1, '规格1', '件', 100.00, 60.00, 0, 1),
('产品B', 'P-002', 1, '规格2', '箱', 500.00, 300.00, 0, 1),
('产品C', 'P-003', 2, '规格3', '个', 50.00, 30.00, 0, 1);

-- 3. 添加仓库库存记录
INSERT INTO `warehouse_inventories` (`product_id`, `warehouse_id`, `quantity`) VALUES
(1, 1, 100),  -- 产品A在主仓库有100件
(2, 1, 50),   -- 产品B在主仓库有50箱
(3, 1, 200),  -- 产品C在主仓库有200个
(1, 2, 20),   -- 产品A在备货仓库有20件
(2, 2, 10);   -- 产品B在备货仓库有10箱

-- 4. 更新产品库存总量
UPDATE `products` SET `stock` = (SELECT SUM(`quantity`) FROM `warehouse_inventories` WHERE `product_id` = `products`.`id`) WHERE `id` IN (1, 2, 3);
