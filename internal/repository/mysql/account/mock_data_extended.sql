-- 插入详细的Mock账户数据
INSERT INTO `go_gin_api`.`account` (
  `account_id`, `username`, `nickname`, `password`, `phone`, 
  `role_type`, `status`, 
  `belong_group_id`, `belong_group_username`, `belong_group_nickname`, `belong_group_created_timestamp`, `belong_group_modified_timestamp`, `belong_group_current_cnt`,
  `belong_team_id`, `belong_team_username`, `belong_team_nickname`, `belong_team_created_timestamp`, `belong_team_modified_timestamp`, `belong_team_current_cnt`,
  `created_timestamp`, `modified_timestamp`, `last_login_timestamp`, 
  `created_user`, `updated_user`
) VALUES 
-- 系统管理员
('0', 'admin', '系统管理员', 'e10adc3949ba59abbe56e057f20f883e', '13800138000', 'company_manager', 'enabled', 
 0, 'admin_group', '系统管理组', 1705123200, 1705123200, 1,
 0, 'admin_team', '系统管理团队', 1705123200, 1705123200, 1,
 1705123200, 1705123200, 1705123200, '系统管理员', '系统管理员'),

-- 公司管理员
('1', 'company_manager', '张伟', 'e10adc3949ba59abbe56e057f20f883e', '13800138001', 'company_manager', 'enabled', 
 1, 'nanjing_tianyuan', '南京-天元大厦组', 1705123200, 1705123200, 15,
 1, 'marketing_team_a', '营销团队A', 1705123200, 1705123200, 8,
 1705123200, 1705123200, 1705123200, '系统管理员', '系统管理员'),

-- 组管理员
('2', 'group_manager', '李明', 'e10adc3949ba59abbe56e057f20f883e', '13800138002', 'group_manager', 'enabled', 
 1, 'nanjing_tianyuan', '南京-天元大厦组', 1705123200, 1705123200, 15,
 NULL, NULL, NULL, NULL, NULL, 0,
 1705123200, 1705123200, 1705123200, '系统管理员', '系统管理员'),

-- 小队管理员
('3', 'team_manager', '王芳', 'e10adc3949ba59abbe56e057f20f883e', '13800138003', 'team_manager', 'enabled', 
 1, 'nanjing_tianyuan', '南京-天元大厦组', 1705123200, 1705123200, 15,
 1, 'marketing_team_a', '营销团队A', 1705123200, 1705123200, 8,
 1705123200, 1705123200, 1705123200, '系统管理员', '系统管理员'),

-- 组管理员2
('4', 'group_manager2', '刘强', 'e10adc3949ba59abbe56e057f20f883e', '13800138004', 'group_manager', 'enabled', 
 2, 'nanjing_nanjing', '南京-南京南站组', 1705123200, 1705123200, 12,
 2, 'marketing_team_b', '营销团队B', 1705123200, 1705123200, 10,
 1705123200, 1705123200, 1705123200, '系统管理员', '系统管理员'),

-- 小队管理员2
('5', 'team_manager2', '赵敏', 'e10adc3949ba59abbe56e057f20f883e', '13800138005', 'team_manager', 'disabled', 
 1, 'nanjing_tianyuan', '南京-天元大厦组', 1705123200, 1705123200, 15,
 3, 'marketing_team_c', '营销团队C', 1705123200, 1705123200, 6,
 1705123200, 1705123200, 1705123200, '系统管理员', '系统管理员'),

-- 员工1
('6', 'employee001', '陈静', 'e10adc3949ba59abbe56e057f20f883e', '13800138006', 'employee', 'enabled', 
 1, 'nanjing_tianyuan', '南京-天元大厦组', 1705123200, 1705123200, 15,
 1, 'marketing_team_a', '营销团队A', 1705123200, 1705123200, 8,
 1705123200, 1705123200, 1705123200, '系统管理员', '系统管理员'),

-- 员工2
('7', 'employee002', '孙丽', 'e10adc3949ba59abbe56e057f20f883e', '13800138007', 'employee', 'enabled', 
 2, 'nanjing_nanjing', '南京-南京南站组', 1705123200, 1705123200, 12,
 2, 'marketing_team_b', '营销团队B', 1705123200, 1705123200, 10,
 1705123200, 1705123200, 1705123200, '系统管理员', '系统管理员'),

-- 员工3
('8', 'employee003', '周杰', 'e10adc3949ba59abbe56e057f20f883e', '13800138008', 'employee', 'enabled', 
 1, 'nanjing_tianyuan', '南京-天元大厦组', 1705123200, 1705123200, 15,
 3, 'marketing_team_c', '营销团队C', 1705123200, 1705123200, 6,
 1705123200, 1705123200, 1705123200, '系统管理员', '系统管理员'),

-- 员工4
('9', 'employee004', '吴婷', 'e10adc3949ba59abbe56e057f20f883e', '13800138009', 'employee', 'disabled', 
 2, 'nanjing_nanjing', '南京-南京南站组', 1705123200, 1705123200, 12,
 2, 'marketing_team_b', '营销团队B', 1705123200, 1705123200, 10,
 1705123200, 1705123200, 1705123200, '系统管理员', '系统管理员'),

-- 员工5
('10', 'employee005', '郑华', 'e10adc3949ba59abbe56e057f20f883e', '13800138010', 'employee', 'enabled', 
 1, 'nanjing_tianyuan', '南京-天元大厦组', 1705123200, 1705123200, 15,
 1, 'marketing_team_a', '营销团队A', 1705123200, 1705123200, 8,
 1705123200, 1705123200, 1705123200, '系统管理员', '系统管理员'); 