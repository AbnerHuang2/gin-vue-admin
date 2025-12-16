-- ================================================
-- Emag 品类管理 - 菜单和 API 权限初始化
-- 在已初始化的系统中执行此 SQL
-- ================================================

-- 1. 添加一级菜单 "Emag"
INSERT INTO `sys_base_menus` (`created_at`, `updated_at`, `menu_level`, `parent_id`, `path`, `name`, `hidden`, `component`, `sort`, `keep_alive`, `default_menu`, `title`, `icon`, `close_tab`)
VALUES (NOW(), NOW(), 0, 0, 'emag', 'emag', 0, 'view/routerHolder.vue', 2, 0, 0, 'Emag', 'goods', 0);

-- 获取刚插入的菜单 ID（如果你的 MySQL 版本支持，可以用 LAST_INSERT_ID()）
SET @emag_menu_id = LAST_INSERT_ID();

-- 2. 添加二级菜单 "品类管理"
INSERT INTO `sys_base_menus` (`created_at`, `updated_at`, `menu_level`, `parent_id`, `path`, `name`, `hidden`, `component`, `sort`, `keep_alive`, `default_menu`, `title`, `icon`, `close_tab`)
VALUES (NOW(), NOW(), 1, @emag_menu_id, 'category', 'emagCategory', 0, 'view/emag/category/index.vue', 1, 0, 0, '品类管理', 'menu', 0);

-- 3. 添加 API 权限
INSERT INTO `sys_apis` (`created_at`, `updated_at`, `path`, `description`, `api_group`, `method`)
VALUES 
(NOW(), NOW(), '/emagCategory/createEmagCategory', '创建品类', 'Emag品类管理', 'POST'),
(NOW(), NOW(), '/emagCategory/deleteEmagCategory', '删除品类', 'Emag品类管理', 'DELETE'),
(NOW(), NOW(), '/emagCategory/deleteEmagCategoryByIds', '批量删除品类', 'Emag品类管理', 'DELETE'),
(NOW(), NOW(), '/emagCategory/updateEmagCategory', '更新品类', 'Emag品类管理', 'PUT'),
(NOW(), NOW(), '/emagCategory/findEmagCategory', '获取单个品类', 'Emag品类管理', 'GET'),
(NOW(), NOW(), '/emagCategory/getEmagCategoryList', '获取品类列表', 'Emag品类管理', 'GET');

-- 4. 给 admin 角色(888)分配菜单权限
-- 先获取刚才插入的两个菜单 ID
INSERT INTO `sys_authority_menus` (`sys_base_menu_id`, `sys_authority_authority_id`)
SELECT id, 888 FROM `sys_base_menus` WHERE `name` IN ('emag', 'emagCategory');

-- 5. 给 admin 角色(888)分配 API 权限（casbin_rule 表）
INSERT INTO `casbin_rule` (`ptype`, `v0`, `v1`, `v2`)
VALUES 
('p', '888', '/emagCategory/createEmagCategory', 'POST'),
('p', '888', '/emagCategory/deleteEmagCategory', 'DELETE'),
('p', '888', '/emagCategory/deleteEmagCategoryByIds', 'DELETE'),
('p', '888', '/emagCategory/updateEmagCategory', 'PUT'),
('p', '888', '/emagCategory/findEmagCategory', 'GET'),
('p', '888', '/emagCategory/getEmagCategoryList', 'GET');

-- ================================================
-- 执行完成后，重启后端服务，刷新前端页面即可看到菜单
-- ================================================

