INSERT IGNORE INTO permissions (name, code, `desc`, type, path, status) VALUES
('Create Inventory Adjust', 'inventory-adjust:create', 'Create inventory adjust request', 2, '/inventory/adjust-request/create', 1),
('List Inventory Adjust', 'inventory-adjust:list', 'List inventory adjust requests', 2, '/inventory/adjust-request/list', 1),
('Approve Inventory Adjust', 'inventory-adjust:approve', 'Approve inventory adjust request', 2, '/inventory/adjust-request/approve', 1),
('Reject Inventory Adjust', 'inventory-adjust:reject', 'Reject inventory adjust request', 2, '/inventory/adjust-request/reject', 1);
