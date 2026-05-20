-- ==============================================
-- RBAC Admin System - Database Schema and Initial Data
-- ==============================================

-- ----------------------------
-- 1. Users Table (users)
-- ----------------------------
DROP TABLE IF EXISTS `users`;
CREATE TABLE `users` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT 'User ID',
  `username` varchar(50) NOT NULL COMMENT 'Username',
  `password` varchar(100) NOT NULL COMMENT 'Password (BCrypt)',
  `nickname` varchar(50) DEFAULT '' COMMENT 'Nickname',
  `email` varchar(100) DEFAULT '' COMMENT 'Email',
  `phone` varchar(20) DEFAULT '' COMMENT 'Phone',
  `status` tinyint DEFAULT '1' COMMENT 'Status (1:active, 0:inactive)',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT 'Created At',
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'Updated At',
  `deleted_at` datetime DEFAULT NULL COMMENT 'Deleted At (soft delete)',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_username` (`username`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='Users Table';

-- ----------------------------
-- 2. Roles Table (roles)
-- ----------------------------
DROP TABLE IF EXISTS `roles`;
CREATE TABLE `roles` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT 'Role ID',
  `name` varchar(50) NOT NULL COMMENT 'Role Name',
  `code` varchar(50) NOT NULL COMMENT 'Role Code',
  `desc` varchar(200) DEFAULT '' COMMENT 'Description',
  `status` tinyint DEFAULT '1' COMMENT 'Status (1:active, 0:inactive)',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT 'Created At',
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'Updated At',
  `deleted_at` datetime DEFAULT NULL COMMENT 'Deleted At (soft delete)',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_name` (`name`),
  UNIQUE KEY `uk_code` (`code`),
  KEY `idx_status` (`status`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='Roles Table';

-- ----------------------------
-- 3. Permissions Table (permissions)
-- ----------------------------
DROP TABLE IF EXISTS `permissions`;
CREATE TABLE `permissions` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT 'Permission ID',
  `name` varchar(50) NOT NULL COMMENT 'Permission Name',
  `code` varchar(50) NOT NULL COMMENT 'Permission Code',
  `desc` varchar(200) DEFAULT '' COMMENT 'Description',
  `type` tinyint DEFAULT '1' COMMENT 'Type (1:menu, 2:button)',
  `parent_id` bigint unsigned DEFAULT '0' COMMENT 'Parent ID',
  `path` varchar(200) DEFAULT '' COMMENT 'Path',
  `component` varchar(200) DEFAULT '' COMMENT 'Component',
  `icon` varchar(50) DEFAULT '' COMMENT 'Icon',
  `sort` int DEFAULT '0' COMMENT 'Sort Order',
  `status` tinyint DEFAULT '1' COMMENT 'Status (1:active, 0:inactive)',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT 'Created At',
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'Updated At',
  `deleted_at` datetime DEFAULT NULL COMMENT 'Deleted At (soft delete)',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_name` (`name`),
  UNIQUE KEY `uk_code` (`code`),
  KEY `idx_type` (`type`),
  KEY `idx_parent_id` (`parent_id`),
  KEY `idx_status` (`status`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='Permissions Table';

-- ----------------------------
-- 4. Menus Table (menus)
-- ----------------------------
DROP TABLE IF EXISTS `menus`;
CREATE TABLE `menus` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT 'Menu ID',
  `name` varchar(50) NOT NULL COMMENT 'Menu Name',
  `code` varchar(50) NOT NULL COMMENT 'Menu Code',
  `desc` varchar(200) DEFAULT '' COMMENT 'Description',
  `parent_id` bigint unsigned DEFAULT '0' COMMENT 'Parent Menu ID',
  `path` varchar(200) NOT NULL COMMENT 'Route Path',
  `component` varchar(200) NOT NULL COMMENT 'Component Path',
  `icon` varchar(50) DEFAULT '' COMMENT 'Icon',
  `sort` int DEFAULT '0' COMMENT 'Sort Order',
  `status` tinyint DEFAULT '1' COMMENT 'Status (1:active, 0:inactive)',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT 'Created At',
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'Updated At',
  `deleted_at` datetime DEFAULT NULL COMMENT 'Deleted At (soft delete)',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_name` (`name`),
  UNIQUE KEY `uk_code` (`code`),
  KEY `idx_parent_id` (`parent_id`),
  KEY `idx_status` (`status`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='Menus Table';

-- ----------------------------
-- 5. User Roles Table (user_roles)
-- ----------------------------
DROP TABLE IF EXISTS `user_roles`;
CREATE TABLE `user_roles` (
  `user_id` bigint unsigned NOT NULL COMMENT 'User ID',
  `role_id` bigint unsigned NOT NULL COMMENT 'Role ID',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT 'Created At',
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'Updated At',
  `deleted_at` datetime DEFAULT NULL COMMENT 'Deleted At (soft delete)',
  PRIMARY KEY (`user_id`,`role_id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_role_id` (`role_id`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='User Roles Table';

-- ----------------------------
-- 6. Role Permissions Table (role_permissions)
-- ----------------------------
DROP TABLE IF EXISTS `role_permissions`;
CREATE TABLE `role_permissions` (
  `role_id` bigint unsigned NOT NULL COMMENT 'Role ID',
  `permission_id` bigint unsigned NOT NULL COMMENT 'Permission ID',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT 'Created At',
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'Updated At',
  `deleted_at` datetime DEFAULT NULL COMMENT 'Deleted At (soft delete)',
  PRIMARY KEY (`role_id`,`permission_id`),
  KEY `idx_role_id` (`role_id`),
  KEY `idx_permission_id` (`permission_id`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='Role Permissions Table';

-- ----------------------------
-- 7. Role Menus Table (role_menus)
-- ----------------------------
DROP TABLE IF EXISTS `role_menus`;
CREATE TABLE `role_menus` (
  `role_id` bigint unsigned NOT NULL COMMENT 'Role ID',
  `menu_id` bigint unsigned NOT NULL COMMENT 'Menu ID',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT 'Created At',
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'Updated At',
  `deleted_at` datetime DEFAULT NULL COMMENT 'Deleted At (soft delete)',
  PRIMARY KEY (`role_id`,`menu_id`),
  KEY `idx_role_id` (`role_id`),
  KEY `idx_menu_id` (`menu_id`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='Role Menus Table';

-- ----------------------------
-- 8. Menu Permissions Table (menu_permissions)
-- ----------------------------
DROP TABLE IF EXISTS `menu_permissions`;
CREATE TABLE `menu_permissions` (
  `menu_id` bigint unsigned NOT NULL COMMENT 'Menu ID',
  `permission_id` bigint unsigned NOT NULL COMMENT 'Permission ID',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT 'Created At',
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'Updated At',
  `deleted_at` datetime DEFAULT NULL COMMENT 'Deleted At (soft delete)',
  PRIMARY KEY (`menu_id`,`permission_id`),
  KEY `idx_menu_id` (`menu_id`),
  KEY `idx_permission_id` (`permission_id`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='Menu Permissions Table';

-- ----------------------------
-- 9. Activities Table (activities)
-- ----------------------------
DROP TABLE IF EXISTS `activities`;
CREATE TABLE `activities` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT 'Activity ID',
  `user_id` bigint unsigned NOT NULL COMMENT 'User ID',
  `username` varchar(50) NOT NULL COMMENT 'Username',
  `action` varchar(255) NOT NULL COMMENT 'Action',
  `url` varchar(255) DEFAULT '' COMMENT 'URL',
  `ip` varchar(50) DEFAULT '' COMMENT 'IP Address',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT 'Created At',
  PRIMARY KEY (`id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_username` (`username`),
  KEY `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='Activities Table';

-- ==============================================
-- Initial Data
-- ==============================================

-- ----------------------------
-- Insert Roles
-- ----------------------------
INSERT INTO `roles` (`id`, `name`, `code`, `desc`, `status`) VALUES 
(1, 'Admin', 'admin', 'System administrator with all permissions', 1),
(2, 'User', 'user', 'Normal user with basic permissions', 1),
(3, 'Auditor', 'auditor', 'Auditor for content review', 1);

-- ----------------------------
-- Insert Users (password: admin123, BCrypt encrypted)
-- ----------------------------
INSERT INTO `users` (`id`, `username`, `password`, `nickname`, `email`, `status`) VALUES 
(1, 'admin', '$2a$10$N9qo8uLOickgx2ZMRZoMye.IjzqAKL9xL5jvMFVdNJHvGCgTq/VEq', 'Admin', 'admin@example.com', 1),
(2, 'test1', '$2a$10$N9qo8uLOickgx2ZMRZoMye.IjzqAKL9xL5jvMFVdNJHvGCgTq/VEq', 'Test User', 'test1@example.com', 1),
(3, 'auditor', '$2a$10$N9qo8uLOickgx2ZMRZoMye.IjzqAKL9xL5jvMFVdNJHvGCgTq/VEq', 'Auditor', 'auditor@example.com', 1);

-- ----------------------------
-- Insert User Roles
-- ----------------------------
INSERT INTO `user_roles` (`user_id`, `role_id`) VALUES 
(1, 1),
(2, 2),
(3, 3);

-- ----------------------------
-- Insert Menus
-- ----------------------------
-- Top Level Menus
INSERT INTO `menus` (`id`, `name`, `code`, `desc`, `parent_id`, `path`, `component`, `icon`, `sort`, `status`) VALUES 
(1, 'Dashboard', 'menu_dashboard', 'Dashboard', 0, '/', 'Dashboard', 'Dashboard', 0, 1),
(2, 'System', 'menu_system', 'System Management', 0, '/system', '', 'Settings', 1, 1),
(3, 'Content', 'menu_content', 'Content Management', 0, '/content', '', 'Document', 2, 1),
(4, 'Statistics', 'menu_statistics', 'Data Statistics', 0, '/statistics', '', 'BarChart', 3, 1);

-- System Management Submenus
INSERT INTO `menus` (`id`, `name`, `code`, `desc`, `parent_id`, `path`, `component`, `icon`, `sort`, `status`) VALUES 
(5, 'Users', 'menu_user', 'User Management', 2, '/user', 'User', 'User', 1, 1),
(6, 'Roles', 'menu_role', 'Role Management', 2, '/role', 'Role', 'UserFilled', 2, 1),
(7, 'Permissions', 'menu_permission', 'Permission Management', 2, '/permission', 'Permission', 'Lock', 3, 1),
(8, 'Menus', 'menu_menu', 'Menu Management', 2, '/menu', 'Menu', 'Menu', 4, 1),
(9, 'Activity Log', 'menu_activity', 'Activity Log', 2, '/activity', 'Activity', 'Clock', 5, 1);

-- Content Management Submenus
INSERT INTO `menus` (`id`, `name`, `code`, `desc`, `parent_id`, `path`, `component`, `icon`, `sort`, `status`) VALUES 
(10, 'Articles', 'menu_article', 'Article Management', 3, '/article', 'Article', 'FileText', 1, 1),
(11, 'Categories', 'menu_category', 'Category Management', 3, '/category', 'Category', 'Folder', 2, 1),
(12, 'Tags', 'menu_tag', 'Tag Management', 3, '/tag', 'Tag', 'Tags', 3, 1);

-- Statistics Submenus
INSERT INTO `menus` (`id`, `name`, `code`, `desc`, `parent_id`, `path`, `component`, `icon`, `sort`, `status`) VALUES 
(13, 'User Stats', 'menu_user_stat', 'User Statistics', 4, '/statistics/user', 'UserStat', 'User', 1, 1),
(14, 'Behavior', 'menu_behavior', 'Behavior Analysis', 4, '/statistics/behavior', 'Behavior', 'TrendingUp', 2, 1),
(15, 'Reports', 'menu_report', 'Data Reports', 4, '/statistics/report', 'Report', 'FileSpreadsheet', 3, 1);

-- ----------------------------
-- Insert Permissions
-- ----------------------------

-- Dashboard Permissions
INSERT INTO `permissions` (`id`, `name`, `code`, `desc`, `type`, `parent_id`, `path`, `component`, `icon`, `sort`, `status`) VALUES 
(1, 'View Dashboard', 'dashboard:view', 'View dashboard', 2, 0, '', '', '', 1, 1);

-- User Management Permissions
INSERT INTO `permissions` (`id`, `name`, `code`, `desc`, `type`, `parent_id`, `path`, `component`, `icon`, `sort`, `status`) VALUES 
(2, 'List Users', 'user:list', 'List users', 2, 0, '', '', '', 1, 1),
(3, 'Create User', 'user:create', 'Create user', 2, 0, '', '', '', 2, 1),
(4, 'Update User', 'user:update', 'Update user', 2, 0, '', '', '', 3, 1),
(5, 'Delete User', 'user:delete', 'Delete user', 2, 0, '', '', '', 4, 1),
(6, 'Assign Roles', 'user:assign-roles', 'Assign roles to user', 2, 0, '', '', '', 5, 1),
(7, 'Reset Password', 'user:reset-password', 'Reset user password', 2, 0, '', '', '', 6, 1);

-- Role Management Permissions
INSERT INTO `permissions` (`id`, `name`, `code`, `desc`, `type`, `parent_id`, `path`, `component`, `icon`, `sort`, `status`) VALUES 
(8, 'List Roles', 'role:list', 'List roles', 2, 0, '', '', '', 1, 1),
(9, 'Create Role', 'role:create', 'Create role', 2, 0, '', '', '', 2, 1),
(10, 'Update Role', 'role:update', 'Update role', 2, 0, '', '', '', 3, 1),
(11, 'Delete Role', 'role:delete', 'Delete role', 2, 0, '', '', '', 4, 1),
(12, 'Assign Permissions', 'role:assign-permissions', 'Assign permissions to role', 2, 0, '', '', '', 5, 1),
(13, 'Assign Menus', 'role:assign-menus', 'Assign menus to role', 2, 0, '', '', '', 6, 1);

-- Permission Management Permissions
INSERT INTO `permissions` (`id`, `name`, `code`, `desc`, `type`, `parent_id`, `path`, `component`, `icon`, `sort`, `status`) VALUES 
(14, 'List Permissions', 'permission:list', 'List permissions', 2, 0, '', '', '', 1, 1),
(15, 'Create Permission', 'permission:create', 'Create permission', 2, 0, '', '', '', 2, 1),
(16, 'Update Permission', 'permission:update', 'Update permission', 2, 0, '', '', '', 3, 1),
(17, 'Delete Permission', 'permission:delete', 'Delete permission', 2, 0, '', '', '', 4, 1);

-- Menu Management Permissions
INSERT INTO `permissions` (`id`, `name`, `code`, `desc`, `type`, `parent_id`, `path`, `component`, `icon`, `sort`, `status`) VALUES 
(18, 'List Menus', 'menu:list', 'List menus', 2, 0, '', '', '', 1, 1),
(19, 'Create Menu', 'menu:create', 'Create menu', 2, 0, '', '', '', 2, 1),
(20, 'Update Menu', 'menu:update', 'Update menu', 2, 0, '', '', '', 3, 1),
(21, 'Delete Menu', 'menu:delete', 'Delete menu', 2, 0, '', '', '', 4, 1),
(22, 'Assign Menu Permissions', 'menu:assign-permissions', 'Assign permissions to menu', 2, 0, '', '', '', 5, 1);

-- Activity Log Permissions
INSERT INTO `permissions` (`id`, `name`, `code`, `desc`, `type`, `parent_id`, `path`, `component`, `icon`, `sort`, `status`) VALUES 
(23, 'List Activities', 'activity:list', 'List activities', 2, 0, '', '', '', 1, 1),
(24, 'Export Activities', 'activity:export', 'Export activities', 2, 0, '', '', '', 2, 1);

-- Article Management Permissions
INSERT INTO `permissions` (`id`, `name`, `code`, `desc`, `type`, `parent_id`, `path`, `component`, `icon`, `sort`, `status`) VALUES 
(25, 'List Articles', 'article:list', 'List articles', 2, 0, '', '', '', 1, 1),
(26, 'Create Article', 'article:create', 'Create article', 2, 0, '', '', '', 2, 1),
(27, 'Update Article', 'article:update', 'Update article', 2, 0, '', '', '', 3, 1),
(28, 'Delete Article', 'article:delete', 'Delete article', 2, 0, '', '', '', 4, 1),
(29, 'Audit Article', 'article:audit', 'Audit article', 2, 0, '', '', '', 5, 1),
(30, 'Publish Article', 'article:publish', 'Publish article', 2, 0, '', '', '', 6, 1);

-- Category Management Permissions
INSERT INTO `permissions` (`id`, `name`, `code`, `desc`, `type`, `parent_id`, `path`, `component`, `icon`, `sort`, `status`) VALUES 
(31, 'List Categories', 'category:list', 'List categories', 2, 0, '', '', '', 1, 1),
(32, 'Create Category', 'category:create', 'Create category', 2, 0, '', '', '', 2, 1),
(33, 'Update Category', 'category:update', 'Update category', 2, 0, '', '', '', 3, 1),
(34, 'Delete Category', 'category:delete', 'Delete category', 2, 0, '', '', '', 4, 1);

-- Tag Management Permissions
INSERT INTO `permissions` (`id`, `name`, `code`, `desc`, `type`, `parent_id`, `path`, `component`, `icon`, `sort`, `status`) VALUES 
(35, 'List Tags', 'tag:list', 'List tags', 2, 0, '', '', '', 1, 1),
(36, 'Create Tag', 'tag:create', 'Create tag', 2, 0, '', '', '', 2, 1),
(37, 'Update Tag', 'tag:update', 'Update tag', 2, 0, '', '', '', 3, 1),
(38, 'Delete Tag', 'tag:delete', 'Delete tag', 2, 0, '', '', '', 4, 1);

-- User Statistics Permissions
INSERT INTO `permissions` (`id`, `name`, `code`, `desc`, `type`, `parent_id`, `path`, `component`, `icon`, `sort`, `status`) VALUES 
(39, 'View User Stats', 'user_stat:view', 'View user statistics', 2, 0, '', '', '', 1, 1),
(40, 'Export User Stats', 'user_stat:export', 'Export user statistics', 2, 0, '', '', '', 2, 1);

-- Behavior Analysis Permissions
INSERT INTO `permissions` (`id`, `name`, `code`, `desc`, `type`, `parent_id`, `path`, `component`, `icon`, `sort`, `status`) VALUES 
(41, 'View Behavior', 'behavior:view', 'View behavior analysis', 2, 0, '', '', '', 1, 1),
(42, 'Export Behavior', 'behavior:export', 'Export behavior analysis', 2, 0, '', '', '', 2, 1);

-- Report Permissions
INSERT INTO `permissions` (`id`, `name`, `code`, `desc`, `type`, `parent_id`, `path`, `component`, `icon`, `sort`, `status`) VALUES 
(43, 'View Reports', 'report:view', 'View reports', 2, 0, '', '', '', 1, 1),
(44, 'Export Reports', 'report:export', 'Export reports', 2, 0, '', '', '', 2, 1);

-- ----------------------------
-- Assign All Permissions to Admin Role
-- ----------------------------
INSERT INTO `role_permissions` (`role_id`, `permission_id`) VALUES 
(1, 1),
(1, 2), (1, 3), (1, 4), (1, 5), (1, 6), (1, 7),
(1, 8), (1, 9), (1, 10), (1, 11), (1, 12), (1, 13),
(1, 14), (1, 15), (1, 16), (1, 17),
(1, 18), (1, 19), (1, 20), (1, 21), (1, 22),
(1, 23), (1, 24),
(1, 25), (1, 26), (1, 27), (1, 28), (1, 29), (1, 30),
(1, 31), (1, 32), (1, 33), (1, 34),
(1, 35), (1, 36), (1, 37), (1, 38),
(1, 39), (1, 40),
(1, 41), (1, 42),
(1, 43), (1, 44);

-- ----------------------------
-- Assign All Menus to Admin Role
-- ----------------------------
INSERT INTO `role_menus` (`role_id`, `menu_id`) VALUES 
(1, 1), (1, 2), (1, 3), (1, 4),
(1, 5), (1, 6), (1, 7), (1, 8), (1, 9),
(1, 10), (1, 11), (1, 12),
(1, 13), (1, 14), (1, 15);

-- ----------------------------
-- Assign Permissions to User Role
-- ----------------------------
INSERT INTO `role_permissions` (`role_id`, `permission_id`) VALUES 
(2, 1),
(2, 2),
(2, 8),
(2, 14),
(2, 18),
(2, 23),
(2, 25), (2, 26), (2, 27),
(2, 31),
(2, 35),
(2, 39),
(2, 41),
(2, 43);

-- ----------------------------
-- Assign Menus to User Role
-- ----------------------------
INSERT INTO `role_menus` (`role_id`, `menu_id`) VALUES 
(2, 1), (2, 2), (2, 3), (2, 4),
(2, 5), (2, 6), (2, 7), (2, 8), (2, 9),
(2, 10), (2, 11), (2, 12),
(2, 13), (2, 14), (2, 15);

-- ----------------------------
-- Assign Permissions to Auditor Role
-- ----------------------------
INSERT INTO `role_permissions` (`role_id`, `permission_id`) VALUES 
(3, 1),
(3, 25), (3, 29), (3, 30),
(3, 23);

-- ----------------------------
-- Assign Menus to Auditor Role
-- ----------------------------
INSERT INTO `role_menus` (`role_id`, `menu_id`) VALUES 
(3, 1), (3, 3), (3, 10);

-- ----------------------------
-- Assign Menu Permissions
-- ----------------------------
INSERT INTO `menu_permissions` (`menu_id`, `permission_id`) VALUES 
(1, 1),
(5, 2), (5, 3), (5, 4), (5, 5), (5, 6), (5, 7),
(6, 8), (6, 9), (6, 10), (6, 11), (6, 12), (6, 13),
(7, 14), (7, 15), (7, 16), (7, 17),
(8, 18), (8, 19), (8, 20), (8, 21), (8, 22),
(9, 23), (9, 24),
(10, 25), (10, 26), (10, 27), (10, 28), (10, 29), (10, 30),
(11, 31), (11, 32), (11, 33), (11, 34),
(12, 35), (12, 36), (12, 37), (12, 38),
(13, 39), (13, 40),
(14, 41), (14, 42),
(15, 43), (15, 44);
