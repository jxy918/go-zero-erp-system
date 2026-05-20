-- 库存调拨功能数据库迁移脚本
-- 创建日期: 2026-05-07

-- 1. 创建调拨单表
CREATE TABLE IF NOT EXISTS `inventory_transfers` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `transfer_no` VARCHAR(50) NOT NULL COMMENT '调拨单号（格式：TF-YYYYMMDD-0001）',
    `from_warehouse_id` BIGINT UNSIGNED NOT NULL COMMENT '源仓库ID',
    `to_warehouse_id` BIGINT UNSIGNED NOT NULL COMMENT '目标仓库ID',
    `product_id` BIGINT UNSIGNED NOT NULL COMMENT '产品ID',
    `quantity` INT NOT NULL COMMENT '调拨数量',
    `status` TINYINT NOT NULL DEFAULT 1 COMMENT '状态（1:待审核, 2:已审核, 3:已完成, 4:已拒绝）',
    `remark` VARCHAR(500) DEFAULT NULL COMMENT '备注',
    `created_by` BIGINT UNSIGNED DEFAULT NULL COMMENT '创建人ID',
    `audited_by` BIGINT UNSIGNED DEFAULT NULL COMMENT '审核人ID',
    `audited_at` DATETIME DEFAULT NULL COMMENT '审核时间',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `deleted_at` DATETIME DEFAULT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_transfer_no` (`transfer_no`),
    KEY `idx_from_warehouse` (`from_warehouse_id`),
    KEY `idx_to_warehouse` (`to_warehouse_id`),
    KEY `idx_product` (`product_id`),
    KEY `idx_status` (`status`),
    KEY `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='库存调拨单表';

-- 2. 插入调拨菜单数据（如果不存在）
INSERT IGNORE INTO `menus` (`name`, `code`, `desc`, `parent_id`, `path`, `component`, `icon`, `sort`, `status`, `created_at`, `updated_at`)
SELECT '库存调拨', 'menu_inventory_transfer', '库存调拨模块', id, '/inventory/transfer', 'inventory/InventoryTransfer', 'Transfer', 5, 1, NOW(), NOW()
FROM `menus` WHERE `code` = 'menu_inventory' OR `code` = 'menu_erp' LIMIT 1;

-- 3. 插入调拨权限数据（如果不存在）
INSERT IGNORE INTO `permissions` (`name`, `code`, `desc`, `type`, `path`, `sort`, `status`, `created_at`, `updated_at`)
VALUES
('查看调拨单', 'inventory_transfer:list', '查看库存调拨单列表', 2, '/inventory/transfer/list', 1, 1, NOW(), NOW()),
('创建调拨单', 'inventory_transfer:create', '创建库存调拨单', 2, '/inventory/transfer/create', 2, 1, NOW(), NOW()),
('更新调拨单', 'inventory_transfer:update', '更新库存调拨单', 2, '/inventory/transfer/update', 3, 1, NOW(), NOW()),
('删除调拨单', 'inventory_transfer:delete', '删除库存调拨单', 2, '/inventory/transfer/delete', 4, 1, NOW(), NOW()),
('审核调拨单', 'inventory_transfer:audit', '审核库存调拨单', 2, '/inventory/transfer/audit', 5, 1, NOW(), NOW());

-- 4. 关联菜单与权限
INSERT IGNORE INTO `menu_permissions` (`menu_id`, `permission_id`, `created_at`, `updated_at`)
SELECT m.id, p.id, NOW(), NOW()
FROM `menus` m
CROSS JOIN `permissions` p
WHERE m.code = 'menu_inventory_transfer'
AND p.code IN (
    'inventory_transfer:list',
    'inventory_transfer:create',
    'inventory_transfer:update',
    'inventory_transfer:delete',
    'inventory_transfer:audit'
);

-- 5. 查询验证
SELECT '数据库迁移完成！' AS message;
