USE admin_system;

-- 添加 code 字段
ALTER TABLE categories ADD COLUMN code VARCHAR(50) AFTER name;

-- 添加 desc 字段
ALTER TABLE categories ADD COLUMN `desc` VARCHAR(255) AFTER deleted_at;

-- 更新现有数据的 code 字段
UPDATE categories SET code = CONCAT('cat_', id) WHERE code IS NULL;

-- 添加唯一索引
ALTER TABLE categories MODIFY COLUMN code VARCHAR(50) NOT NULL UNIQUE;