-- 检查并添加 safety_stock 字段
USE admin_system;

-- 检查 safety_stock 字段是否存在
SELECT COUNT(*) INTO @exists 
FROM information_schema.columns 
WHERE table_schema = 'admin_system' 
AND table_name = 'products' 
AND column_name = 'safety_stock';

-- 如果不存在则添加字段
SET @sql = IF(@exists = 0, 
    'ALTER TABLE products ADD COLUMN safety_stock INT DEFAULT 0', 
    'SELECT "safety_stock column already exists" AS result');

PREPARE stmt FROM @sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- 显示 products 表结构
DESCRIBE products;