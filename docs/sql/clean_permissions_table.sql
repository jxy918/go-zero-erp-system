-- 清理 permissions 表中无用的字段
-- 注意：先备份数据再执行此脚本

-- 1. 删除 type 字段
ALTER TABLE permissions DROP COLUMN IF EXISTS `type`;

-- 2. 删除 path 字段
ALTER TABLE permissions DROP COLUMN IF EXISTS `path`;

-- 3. 删除 component 字段
ALTER TABLE permissions DROP COLUMN IF EXISTS `component`;

-- 4. 删除 icon 字段
ALTER TABLE permissions DROP COLUMN IF EXISTS `icon`;

-- 验证结果
DESCRIBE permissions;