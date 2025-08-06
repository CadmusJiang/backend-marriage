-- 插入Mock账户历史记录数据
INSERT INTO `go_gin_api`.`account_history` (
  `history_id`, `account_id`, `operate_type`, `operate_timestamp`, 
  `content`, `operator`, `operator_role_type`, `created_user`, `updated_user`
) VALUES 
('hist_001', 'acc_001', 'created', NOW() - INTERVAL 1 MONTH, 
 '{"username": {"old": "", "new": "admin"}, "nickname": {"old": "", "new": "管理员"}}', 
 'system', 'admin', 'system', 'system'),
('hist_002', 'acc_001', 'modified', NOW() - INTERVAL 2 HOUR, 
 '{"phone": {"old": "13800138000", "new": "13800138001"}}', 
 'admin', 'admin', 'admin', 'admin'),
('hist_003', 'acc_002', 'created', NOW() - INTERVAL 2 MONTH, 
 '{"username": {"old": "", "new": "user001"}, "nickname": {"old": "", "new": "张三"}}', 
 'admin', 'admin', 'admin', 'admin'),
('hist_004', 'acc_002', 'modified', NOW() - INTERVAL 1 HOUR, 
 '{"belongGroup": {"old": "技术组", "new": "运营组"}}', 
 'manager001', 'manager', 'manager001', 'manager001'),
('hist_005', 'acc_003', 'created', NOW() - INTERVAL 3 MONTH, 
 '{"username": {"old": "", "new": "user002"}, "nickname": {"old": "", "new": "李四"}}', 
 'admin', 'admin', 'admin', 'admin'),
('hist_006', 'acc_003', 'modified', NOW() - INTERVAL 1 DAY, 
 '{"status": {"old": "active", "new": "inactive"}}', 
 'manager001', 'manager', 'manager001', 'manager001'); 