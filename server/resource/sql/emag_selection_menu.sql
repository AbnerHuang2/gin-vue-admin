-- ================================================
-- Emag 选品中心 & 选品分析 - 菜单初始化
-- 在已有 Emag 一级菜单的基础上执行此 SQL
-- ================================================

-- 1. 获取 Emag 一级菜单的 ID
SET @emag_menu_id = (SELECT id FROM `sys_base_menus` WHERE `name` = 'emag' LIMIT 1);

-- 2. 添加二级菜单 "选品中心"
INSERT INTO `sys_base_menus` (`created_at`, `updated_at`, `menu_level`, `parent_id`, `path`, `name`, `hidden`, `component`, `sort`, `keep_alive`, `default_menu`, `title`, `icon`, `close_tab`)
VALUES (NOW(), NOW(), 1, @emag_menu_id, 'discovery', 'emagDiscovery', 0, 'view/emag/discovery/index.vue', 5, 0, 0, '选品中心', 'search', 0);

-- 3. 添加二级菜单 "选品分析"
INSERT INTO `sys_base_menus` (`created_at`, `updated_at`, `menu_level`, `parent_id`, `path`, `name`, `hidden`, `component`, `sort`, `keep_alive`, `default_menu`, `title`, `icon`, `close_tab`)
VALUES (NOW(), NOW(), 1, @emag_menu_id, 'selection', 'emagSelection', 0, 'view/emag/selection/index.vue', 6, 0, 0, '选品分析', 'data-analysis', 0);

-- 4. 给 admin 角色(888)分配菜单权限
INSERT INTO `sys_authority_menus` (`sys_base_menu_id`, `sys_authority_authority_id`)
SELECT id, 888 FROM `sys_base_menus` WHERE `name` IN ('emagDiscovery', 'emagSelection');

-- ================================================
-- 执行完成后，重启后端服务，刷新前端页面即可看到菜单
-- ================================================
