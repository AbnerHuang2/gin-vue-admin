-- ================================================
-- Emag 品类指标 - 建表语句
-- 使用 InnoDB 引擎，id 自增
-- ================================================

CREATE TABLE IF NOT EXISTS `emag_category_stat` (
    `id` BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY COMMENT '主键ID',
    `category_id` VARCHAR(64) NOT NULL COMMENT '品类ID',
    `total` INT DEFAULT NULL COMMENT '总数',
    `supper_hot_total` INT DEFAULT NULL COMMENT '超热销总数',
    `oem_total` INT DEFAULT NULL COMMENT 'OEM总数',
    `oem_supper_hot_total` INT DEFAULT NULL COMMENT 'OEM超热销总数',
    `supper_hot_rate` FLOAT DEFAULT NULL COMMENT '超热销率',
    `oem_supper_hot_rate` FLOAT DEFAULT NULL COMMENT 'OEM超热销率',
    `snapshot_date` DATE NOT NULL COMMENT '快照日期',
    `tags` VARCHAR(30) DEFAULT NULL COMMENT '标签',
    `create_at` DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `update_at` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    INDEX `idx_snapshot_date` (`snapshot_date`),
    INDEX `idx_category_id` (`category_id`),
    INDEX `idx_supper_hot_rate` (`supper_hot_rate`),
    INDEX `idx_oem_supper_hot_rate` (`oem_supper_hot_rate`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='Emag品类指标统计表';

