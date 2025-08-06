-- 数据库初始化脚本
-- 删除旧表并创建新表结构

-- 删除旧表（如果存在）
DROP TABLE IF EXISTS `go_gin_api`.`account_history`;
DROP TABLE IF EXISTS `go_gin_api`.`account`;

-- 创建新的账户表
CREATE TABLE `go_gin_api`.`account` (
  `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
  `account_id` varchar(32) NOT NULL COMMENT '账户ID',
  `username` varchar(32) NOT NULL COMMENT '用户名',
  `nickname` varchar(60) NOT NULL COMMENT '姓名',
  `password` varchar(100) NOT NULL COMMENT '密码',
  `phone` varchar(20) DEFAULT NULL COMMENT '手机号',
  `role_type` varchar(20) NOT NULL DEFAULT 'employee' COMMENT '角色类型: company_manager, group_manager, team_manager, employee',
  `status` varchar(20) NOT NULL DEFAULT 'enabled' COMMENT '状态: enabled, disabled',
  
  -- 所属组信息
  `belong_group_id` int DEFAULT NULL COMMENT '所属组ID',
  `belong_group_username` varchar(60) DEFAULT NULL COMMENT '所属组用户名',
  `belong_group_nickname` varchar(60) DEFAULT NULL COMMENT '所属组名称',
  `belong_group_created_timestamp` bigint DEFAULT NULL COMMENT '所属组创建时间戳',
  `belong_group_modified_timestamp` bigint DEFAULT NULL COMMENT '所属组修改时间戳',
  `belong_group_current_cnt` int DEFAULT 0 COMMENT '所属组当前成员数量',
  
  -- 所属团队信息
  `belong_team_id` int DEFAULT NULL COMMENT '所属团队ID',
  `belong_team_username` varchar(60) DEFAULT NULL COMMENT '所属团队用户名',
  `belong_team_nickname` varchar(60) DEFAULT NULL COMMENT '所属团队名称',
  `belong_team_created_timestamp` bigint DEFAULT NULL COMMENT '所属团队创建时间戳',
  `belong_team_modified_timestamp` bigint DEFAULT NULL COMMENT '所属团队修改时间戳',
  `belong_team_current_cnt` int DEFAULT 0 COMMENT '所属团队当前成员数量',
  
  `created_timestamp` bigint NOT NULL COMMENT '创建时间戳',
  `modified_timestamp` bigint NOT NULL COMMENT '修改时间戳',
  `last_login_timestamp` bigint DEFAULT NULL COMMENT '最后登录时间戳',
  
  `is_deleted` tinyint(1) NOT NULL DEFAULT '-1' COMMENT '是否删除 1:是  -1:否',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `created_user` varchar(60) NOT NULL COMMENT '创建人',
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `updated_user` varchar(60) NOT NULL COMMENT '更新人',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_account_id` (`account_id`),
  UNIQUE KEY `uk_username` (`username`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='账户表';

-- 创建新的账户历史记录表
CREATE TABLE `go_gin_api`.`account_history` (
  `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
  `history_id` varchar(32) NOT NULL COMMENT '历史记录ID',
  `account_id` varchar(32) NOT NULL COMMENT '账户ID',
  `operate_type` varchar(20) NOT NULL COMMENT '操作类型: created, modified, deleted',
  `operate_timestamp` bigint NOT NULL COMMENT '操作时间戳',
  `content` json DEFAULT NULL COMMENT '操作内容',
  
  -- 操作人信息
  `operator_username` varchar(60) NOT NULL COMMENT '操作人用户名',
  `operator_nickname` varchar(60) NOT NULL COMMENT '操作人姓名',
  `operator_role_type` varchar(20) NOT NULL COMMENT '操作人角色类型',
  
  `is_deleted` tinyint(1) NOT NULL DEFAULT '-1' COMMENT '是否删除 1:是  -1:否',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `created_user` varchar(60) NOT NULL COMMENT '创建人',
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `updated_user` varchar(60) NOT NULL COMMENT '更新人',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_history_id` (`history_id`),
  KEY `idx_account_id` (`account_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='账户历史记录表';

-- 插入Mock数据
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

-- 插入账户历史记录数据
INSERT INTO `go_gin_api`.`account_history` (
  `history_id`, `account_id`, `operate_type`, `operate_timestamp`, `content`,
  `operator_username`, `operator_nickname`, `operator_role_type`,
  `created_user`, `updated_user`
) VALUES 
-- 账户6的历史记录
('1', '6', 'modified', 1705923000, '{"roleType": {"old": "员工", "new": "小队管理员"}, "belongTeam": {"old": "无", "new": "营销团队A"}, "status": {"old": "enabled", "new": "disabled"}}', 'zhangwei', '张伟', 'company_manager', '系统管理员', '系统管理员'),

('2', '6', 'modified', 1705754700, '{"belongGroup": {"old": "南京-天元大厦组", "new": "南京-南京南站组"}, "belongTeam": {"old": "营销团队A", "new": "营销团队C"}}', 'liming', '李明', 'team_manager', '系统管理员', '系统管理员'),

('3', '6', 'modified', 1705565700, '{"status": {"old": "enabled", "new": "disabled"}}', 'wangfang', '王芳', 'team_manager', '系统管理员', '系统管理员'),

('4', '6', 'modified', 1705303200, '{"phone": {"old": "13800138000", "new": "13900139000"}, "belongTeam": {"old": "营销团队A", "new": "营销团队B"}}', 'liuqiang', '刘强', 'team_manager', '系统管理员', '系统管理员'),

('5', '6', 'modified', 1705032000, '{"nickname": {"old": "张三", "new": "张明"}, "roleType": {"old": "员工", "new": "小队管理员"}}', 'chenjing', '陈静', 'team_manager', '系统管理员', '系统管理员'),

('6', '6', 'created', 1704877800, '{"username": {"old": "", "new": "employee001"}, "nickname": {"old": "", "new": "张三"}, "roleType": {"old": "", "new": "员工"}, "phone": {"old": "", "new": "13800138000"}, "belongGroup": {"old": "", "new": "南京-天元大厦组"}, "belongTeam": {"old": "", "new": "营销团队A"}}', 'admin', '系统管理员', 'company_manager', '系统管理员', '系统管理员'),

('7', '6', 'modified', 1704709200, '{"belongGroup": {"old": "南京-天元大厦组", "new": "南京-夫子庙组"}, "belongTeam": {"old": "营销团队A", "new": "营销团队D"}}', 'zhaomin', '赵敏', 'team_manager', '系统管理员', '系统管理员'),

('8', '6', 'modified', 1704457500, '{"status": {"old": "disabled", "new": "enabled"}, "reason": {"old": "", "new": "问题已解决，恢复账户"}}', 'liming', '李明', 'team_manager', '系统管理员', '系统管理员'),

('9', '6', 'modified', 1704265800, '{"password": {"old": "****", "new": "******"}}', 'zhangwei', '张伟', 'company_manager', '系统管理员', '系统管理员'),

('10', '6', 'modified', 1704081600, '{"password": {"old": "****", "new": "******"}}', 'wangfang', '王芳', 'team_manager', '系统管理员', '系统管理员'),

-- 账户7的历史记录
('11', '7', 'modified', 1705923000, '{"roleType": {"old": "员工", "new": "小队管理员"}}', 'zhangwei', '张伟', 'company_manager', '系统管理员', '系统管理员'),

('12', '7', 'created', 1704877800, '{"username": {"old": "", "new": "employee002"}, "nickname": {"old": "", "new": "孙丽"}, "roleType": {"old": "", "new": "员工"}, "phone": {"old": "", "new": "13800138007"}, "belongGroup": {"old": "", "new": "南京-南京南站组"}, "belongTeam": {"old": "", "new": "营销团队B"}}', 'admin', '系统管理员', 'company_manager', '系统管理员', '系统管理员'),

-- 账户8的历史记录
('13', '8', 'modified', 1705754700, '{"belongTeam": {"old": "营销团队A", "new": "营销团队C"}}', 'liming', '李明', 'team_manager', '系统管理员', '系统管理员'),

('14', '8', 'created', 1704877800, '{"username": {"old": "", "new": "employee003"}, "nickname": {"old": "", "new": "周杰"}, "roleType": {"old": "", "new": "员工"}, "phone": {"old": "", "new": "13800138008"}, "belongGroup": {"old": "", "new": "南京-天元大厦组"}, "belongTeam": {"old": "", "new": "营销团队A"}}', 'admin', '系统管理员', 'company_manager', '系统管理员', '系统管理员'),

-- 账户9的历史记录
('15', '9', 'modified', 1705565700, '{"status": {"old": "enabled", "new": "disabled"}}', 'wangfang', '王芳', 'team_manager', '系统管理员', '系统管理员'),

('16', '9', 'created', 1704877800, '{"username": {"old": "", "new": "employee004"}, "nickname": {"old": "", "new": "吴婷"}, "roleType": {"old": "", "new": "员工"}, "phone": {"old": "", "new": "13800138009"}, "belongGroup": {"old": "", "new": "南京-南京南站组"}, "belongTeam": {"old": "", "new": "营销团队B"}}', 'admin', '系统管理员', 'company_manager', '系统管理员', '系统管理员'),

-- 账户10的历史记录
('17', '10', 'created', 1704877800, '{"username": {"old": "", "new": "employee005"}, "nickname": {"old": "", "new": "郑华"}, "roleType": {"old": "", "new": "员工"}, "phone": {"old": "", "new": "13800138010"}, "belongGroup": {"old": "", "new": "南京-天元大厦组"}, "belongTeam": {"old": "", "new": "营销团队A"}}', 'admin', '系统管理员', 'company_manager', '系统管理员', '系统管理员'); 