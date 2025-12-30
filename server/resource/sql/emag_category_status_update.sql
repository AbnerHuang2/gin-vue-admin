-- =====================================================
-- emag_category 表结构更新
-- 添加状态管理相关字段，用于追踪 API 调用失败的分类
-- 执行时间：2025-12-30
-- =====================================================

-- 一次性添加所有字段
ALTER TABLE emag_category
    ADD COLUMN status VARCHAR(20) DEFAULT 'normal' COMMENT '状态: normal/bad_request',
ADD COLUMN fail_count INT DEFAULT 0 COMMENT '连续失败次数',
ADD COLUMN last_fail_reason VARCHAR(255) DEFAULT NULL COMMENT '最后失败原因',
ADD COLUMN last_fail_at DATETIME DEFAULT NULL COMMENT '最后失败时间',
ADD INDEX idx_emag_category_status (status);