-- 更新库存管理相关权限为中文
UPDATE permissions SET name='创建库存调整申请', `desc`='创建库存调整申请', parent_id=13 WHERE code='inventory-adjust:create';
UPDATE permissions SET name='查看库存调整申请', `desc`='查看库存调整申请', parent_id=13 WHERE code='inventory-adjust:list';
UPDATE permissions SET name='审核库存调整申请', `desc`='审核库存调整申请', parent_id=13 WHERE code='inventory-adjust:approve';
UPDATE permissions SET name='拒绝库存调整申请', `desc`='拒绝库存调整申请', parent_id=13 WHERE code='inventory-adjust:reject';

-- 查看更新结果
SELECT id, name, code, `desc`, parent_id FROM permissions WHERE code LIKE 'inventory%';
