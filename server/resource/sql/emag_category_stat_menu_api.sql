-- ================================================
-- Emag 品类指标 - 菜单和 API 权限初始化
-- 在已有 Emag 一级菜单的基础上执行此 SQL
-- ================================================

-- 1. 获取 Emag 一级菜单的 ID
SET @emag_menu_id = (SELECT id FROM `sys_base_menus` WHERE `name` = 'emag' LIMIT 1);

-- 2. 添加二级菜单 "品类指标"
INSERT INTO `sys_base_menus` (`created_at`, `updated_at`, `menu_level`, `parent_id`, `path`, `name`, `hidden`, `component`, `sort`, `keep_alive`, `default_menu`, `title`, `icon`, `close_tab`)
VALUES (NOW(), NOW(), 1, @emag_menu_id, 'categoryStat', 'emagCategoryStat', 0, 'view/emag/categoryStat/index.vue', 2, 0, 0, '品类指标', 'data-analysis', 0);

-- 3. 添加 API 权限
INSERT INTO `sys_apis` (`created_at`, `updated_at`, `path`, `description`, `api_group`, `method`)
VALUES
(NOW(), NOW(), '/emagCategoryStat/getSnapshotDateList', '获取快照日期列表', 'Emag品类指标', 'GET'),
(NOW(), NOW(), '/emagCategoryStat/getTop20', '获取品类指标Top20', 'Emag品类指标', 'GET'),
(NOW(), NOW(), '/emagCategoryStat/getGrowthRank', '获取品类指标同比增长排名', 'Emag品类指标', 'GET');

-- 4. 给 admin 角色(888)分配菜单权限
INSERT INTO `sys_authority_menus` (`sys_base_menu_id`, `sys_authority_authority_id`)
SELECT id, 888 FROM `sys_base_menus` WHERE `name` = 'emagCategoryStat';

-- 5. 给 admin 角色(888)分配 API 权限（casbin_rule 表）
INSERT INTO `casbin_rule` (`ptype`, `v0`, `v1`, `v2`)
VALUES
('p', '888', '/emagCategoryStat/getSnapshotDateList', 'GET'),
('p', '888', '/emagCategoryStat/getTop20', 'GET'),
('p', '888', '/emagCategoryStat/getGrowthRank', 'GET');

-- ================================================
-- 执行完成后，重启后端服务，刷新前端页面即可看到菜单
-- ================================================

