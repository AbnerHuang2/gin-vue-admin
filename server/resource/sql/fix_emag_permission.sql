-- ================================================
-- 修复 Emag 品类管理 API 权限
-- ================================================

-- 先删除可能存在的旧记录（避免重复）
DELETE FROM `casbin_rule` WHERE `v1` LIKE '/emagCategory%';

-- 重新插入 casbin 规则给 admin 角色 (888)
INSERT INTO `casbin_rule` (`ptype`, `v0`, `v1`, `v2`)
VALUES 
('p', '888', '/emagCategory/createEmagCategory', 'POST'),
('p', '888', '/emagCategory/deleteEmagCategory', 'DELETE'),
('p', '888', '/emagCategory/deleteEmagCategoryByIds', 'DELETE'),
('p', '888', '/emagCategory/updateEmagCategory', 'PUT'),
('p', '888', '/emagCategory/findEmagCategory', 'GET'),
('p', '888', '/emagCategory/getEmagCategoryList', 'GET');

-- 验证是否插入成功
SELECT * FROM `casbin_rule` WHERE `v1` LIKE '/emagCategory%';

-- ================================================
-- 执行完成后，【必须重启后端服务】让 Casbin 重新加载规则！
-- ================================================

