-- 添加 executed_by 字段
ALTER TABLE inventory_transfers ADD COLUMN executed_by BIGINT UNSIGNED NULL AFTER audited_by;
