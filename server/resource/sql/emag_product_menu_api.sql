-- ================================================
-- Emag 产品管理 - 菜单和数据库表初始化
-- 在已有 Emag 一级菜单的基础上执行此 SQL
-- ================================================

-- 1. 获取 Emag 一级菜单的 ID
SET @emag_menu_id = (SELECT id FROM `sys_base_menus` WHERE `name` = 'emag' LIMIT 1);

-- 2. 添加二级菜单 "产品管理"
INSERT INTO `sys_base_menus` (`created_at`, `updated_at`, `menu_level`, `parent_id`, `path`, `name`, `hidden`, `component`, `sort`, `keep_alive`, `default_menu`, `title`, `icon`, `close_tab`)
VALUES (NOW(), NOW(), 1, @emag_menu_id, 'product', 'emagProduct', 0, 'view/emag/product/index.vue', 4, 0, 0, '产品管理', 'goods', 0);

-- 3. 给 admin 角色(888)分配菜单权限
INSERT INTO `sys_authority_menus` (`sys_base_menu_id`, `sys_authority_authority_id`)
SELECT id, 888 FROM `sys_base_menus` WHERE `name` = 'emagProduct';

-- 4. 创建产品表
CREATE TABLE IF NOT EXISTS `emag_product` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `product_id` VARCHAR(64) COMMENT 'id',
  `category_id` INT COMMENT 'categoryId',
  `ext_id` VARCHAR(64) COMMENT 'extId ｜',
  `pn` VARCHAR(64) COMMENT 'extPartNumber',
  `pnk` VARCHAR(64) COMMENT 'docProductPartNumberKey ｜',
  `ean` VARCHAR(64) COMMENT 'eans数组下第一个',
  `title` VARCHAR(256) COMMENT 'docProductName ｜',
  `status` VARCHAR(32) COMMENT 'statusName ｜',
  `sale_price` DECIMAL(10,2) COMMENT 'extSalePrice ｜',
  `after_tax_price` DECIMAL(10,2) COMMENT '+ extSalePrice * 1 + vat ｜',
  `currency` VARCHAR(32) COMMENT 'currency ｜',
  `country` VARCHAR(20) COMMENT '+根据currency设置吧。目前只有ro，bg, hu三个国家 ｜',
  `sale_price_cn` DECIMAL(10,2) COMMENT '+ 根据price + currency转换 ｜',
  `cost_price_cn` DECIMAL(10,2) COMMENT '- 后续自己去数据库补充 ｜',
  `vat` VARCHAR(64) COMMENT ' vatValues对象的key ｜',
  `stock` INT COMMENT 'extStock ｜',
  `url` VARCHAR(256) COMMENT 'links对象的details属性对应的value ｜',
  `buy_button_rank` TINYINT COMMENT 'productPerformanceBuyButtonRank ｜',
  `buy_button_cnt` TINYINT COMMENT 'productPerformanceMultiofferOffersCount ｜',
  `created_at` TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_product_id` (`product_id`),
  KEY `idx_ext_id` (`ext_id`),
  KEY `idx_country` (`country`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='emag_product';

-- ================================================
-- 执行完成后，重启后端服务，刷新前端页面即可看到菜单
-- ================================================
