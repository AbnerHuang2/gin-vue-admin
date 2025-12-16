-- Emag 品类管理表
-- 如果系统已初始化，可手动执行此 SQL 创建表

CREATE TABLE IF NOT EXISTS `emag_category` (
    `id` INT AUTO_INCREMENT PRIMARY KEY COMMENT '主键ID',
    `category_id` VARCHAR(64) NOT NULL COMMENT '品类ID',
    `category_name` VARCHAR(128) NOT NULL COMMENT '品类名称',
    `subcategory_name` VARCHAR(128) DEFAULT NULL COMMENT '子品类名称',
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    INDEX `idx_category_id` (`category_id`),
    INDEX `idx_category_name` (`category_name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='Emag品类管理表';

