-- ================================================
-- Emag 订单管理 - 菜单和数据库表初始化
-- 在已有 Emag 一级菜单的基础上执行此 SQL
-- ================================================

-- 1. 获取 Emag 一级菜单的 ID
SET @emag_menu_id = (SELECT id FROM `sys_base_menus` WHERE `name` = 'emag' LIMIT 1);

-- 2. 添加二级菜单 "订单管理"
INSERT INTO `sys_base_menus` (`created_at`, `updated_at`, `menu_level`, `parent_id`, `path`, `name`, `hidden`, `component`, `sort`, `keep_alive`, `default_menu`, `title`, `icon`, `close_tab`)
VALUES (NOW(), NOW(), 1, @emag_menu_id, 'order', 'emagOrder', 0, 'view/emag/order/index.vue', 3, 0, 0, '订单管理', 'shopping-cart', 0);

-- 3. 给 admin 角色(888)分配菜单权限
INSERT INTO `sys_authority_menus` (`sys_base_menu_id`, `sys_authority_authority_id`)
SELECT id, 888 FROM `sys_base_menus` WHERE `name` = 'emagOrder';

-- 4. 创建订单表
CREATE TABLE IF NOT EXISTS `emag_order` (
  `id` BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY,
  `order_id` VARCHAR(32) NOT NULL COMMENT '订单ID，对应mktId',
  `order_date_local` VARCHAR(32) COMMENT '订单日期，对应date',
  `country` VARCHAR(20) COMMENT '国家，对应customerBillingCountry',
  `currency` VARCHAR(20) COMMENT '货币，对应currency',
  `price` DECIMAL(10, 2) NOT NULL DEFAULT 0 COMMENT '金额',
  `price_cny` DECIMAL(10, 2) NOT NULL DEFAULT 0 COMMENT '金额(人民币)',
  `status` VARCHAR(32) COMMENT '订单状态',
  `customer_name` VARCHAR(128) COMMENT '客户名称',
  `created_at` TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  UNIQUE KEY `uk_order_id` (`order_id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='Emag订单表';

-- 5. 创建订单产品表
CREATE TABLE IF NOT EXISTS `emag_order_product` (
  `id` BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY,
  `order_id` VARCHAR(32) NOT NULL COMMENT '订单ID，对应mktId',
  `product_id` VARCHAR(32) COMMENT '产品ID，对应partNumberKey',
  `product_url` VARCHAR(512) COMMENT '产品链接',
  `product_name` VARCHAR(512) COMMENT '产品名称',
  `sale_price` DECIMAL(10, 2) COMMENT '销售价格',
  `quantity` FLOAT COMMENT '数量',
  `vat_rate` FLOAT COMMENT '增值税率',
  `created_at` TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  KEY `idx_order_id` (`order_id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='Emag订单产品表';

-- ================================================
-- 执行完成后，重启后端服务，刷新前端页面即可看到菜单
-- ================================================

