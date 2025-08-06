-- 插入Mock账户数据
INSERT INTO `go_gin_api`.`account` (
  `account_id`, `username`, `nickname`, `password`, `phone`, 
  `role_type`, `status`, `belong_group`, `belong_team`, 
  `last_login_time`, `created_user`, `updated_user`
) VALUES 
('acc_001', 'admin', '管理员', 'e10adc3949ba59abbe56e057f20f883e', '13800138001', 'admin', 'active', '技术组', '开发团队', NOW() - INTERVAL 2 HOUR, 'system', 'system'),
('acc_002', 'user001', '张三', 'e10adc3949ba59abbe56e057f20f883e', '13800138002', 'user', 'active', '运营组', '市场团队', NOW() - INTERVAL 1 HOUR, 'admin', 'admin'),
('acc_003', 'user002', '李四', 'e10adc3949ba59abbe56e057f20f883e', '13800138003', 'user', 'inactive', '销售组', '销售团队', NOW() - INTERVAL 1 DAY, 'admin', 'admin'),
('acc_004', 'manager001', '王五', 'e10adc3949ba59abbe56e057f20f883e', '13800138004', 'manager', 'active', '管理组', '管理团队', NOW() - INTERVAL 30 MINUTE, 'admin', 'admin'),
('acc_005', 'user003', '赵六', 'e10adc3949ba59abbe56e057f20f883e', '13800138005', 'user', 'active', '技术组', '测试团队', NOW() - INTERVAL 3 HOUR, 'admin', 'admin'); 