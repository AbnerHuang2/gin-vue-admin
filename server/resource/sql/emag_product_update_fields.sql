-- ================================================
-- Emag 产品表字段更新脚本
-- 用于已创建表的字段添加和索引创建
-- ================================================

-- 1. 检查并添加 ext_id 字段
ALTER TABLE `emag_product` 
ADD COLUMN IF NOT EXISTS `ext_id` VARCHAR(64) COMMENT 'extId ｜' AFTER `category_id`;

-- 2. 检查并添加 country 字段
ALTER TABLE `emag_product` 
ADD COLUMN IF NOT EXISTS `country` VARCHAR(20) COMMENT '+根据currency设置吧。目前只有ro，bg, hu三个国家 ｜' AFTER `currency`;

-- 3. 创建索引（如果不存在）
ALTER TABLE `emag_product` 
ADD INDEX IF NOT EXISTS `idx_ext_id` (`ext_id`);

ALTER TABLE `emag_product` 
ADD INDEX IF NOT EXISTS `idx_country` (`country`);

-- ================================================
-- 执行完成后，重启后端服务即可
-- ================================================
