-- 账户管理系统数据库初始化脚本

-- 删除旧表（如果存在）
DROP TABLE IF EXISTS `go_gin_api`.`account_history`;
DROP TABLE IF EXISTS `go_gin_api`.`account`;

-- 创建账户表
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

-- 创建账户历史记录表
CREATE TABLE `go_gin_api`.`account_history` (
  `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
  `history_id` varchar(32) NOT NULL COMMENT '历史记录ID',
  `account_id` varchar(32) NOT NULL COMMENT '账户ID',
  `operate_type` varchar(20) NOT NULL COMMENT '操作类型: created, modified, deleted',
  `operate_timestamp` bigint NOT NULL COMMENT '操作时间戳',
  `content` json DEFAULT NULL COMMENT '操作内容 (JSON格式)',
  `operator` varchar(60) NOT NULL COMMENT '操作人',
  `operator_role_type` varchar(20) NOT NULL COMMENT '操作人角色',
  
  `is_deleted` tinyint(1) NOT NULL DEFAULT '-1' COMMENT '是否删除 1:是  -1:否',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `created_user` varchar(60) NOT NULL COMMENT '创建人',
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `updated_user` varchar(60) NOT NULL COMMENT '更新人',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_history_id` (`history_id`),
  KEY `idx_account_id` (`account_id`),
  KEY `idx_operate_timestamp` (`operate_timestamp`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='账户历史记录表';

-- 插入测试数据
-- 注意：密码都是 '123456' 的MD5值：e10adc3949ba59abbe56e057f20f883e

INSERT INTO `go_gin_api`.`account` (
  `account_id`, `username`, `nickname`, `password`, `phone`, `role_type`, `status`,
  `belong_group_id`, `belong_group_username`, `belong_group_nickname`, 
  `belong_group_created_timestamp`, `belong_group_modified_timestamp`, `belong_group_current_cnt`,
  `belong_team_id`, `belong_team_username`, `belong_team_nickname`,
  `belong_team_created_timestamp`, `belong_team_modified_timestamp`, `belong_team_current_cnt`,
  `created_timestamp`, `modified_timestamp`, `last_login_timestamp`,
  `is_deleted`, `created_user`, `updated_user`
) VALUES
-- 管理员账户
('acc_001', 'admin', '系统管理员', 'e10adc3949ba59abbe56e057f20f883e', '13800138000', 'company_manager', 'enabled',
 1, 'admin_group', '系统管理组', 1705123200, 1705123200, 1,
 1, 'admin_team', '系统管理团队', 1705123200, 1705123200, 1,
 1705123200, 1705123200, 1705123200, -1, 'system', 'system'),

-- 公司经理
('acc_002', 'company_manager', '张伟', 'e10adc3949ba59abbe56e057f20f883e', '13800138001', 'company_manager', 'enabled',
 2, 'nanjing_tianyuan', '南京-天元大厦组', 1705123200, 1705123200, 15,
 2, 'marketing_team_a', '营销团队A', 1705123200, 1705123200, 8,
 1705123200, 1705123200, 1705123200, -1, 'system', 'system'),

-- 组经理
('acc_003', 'group_manager', '李娜', 'e10adc3949ba59abbe56e057f20f883e', '13800138002', 'group_manager', 'enabled',
 3, 'beijing_office', '北京办公室组', 1705123200, 1705123200, 12,
 3, 'sales_team_b', '销售团队B', 1705123200, 1705123200, 6,
 1705123200, 1705123200, 1705123200, -1, 'system', 'system'),

-- 团队经理
('acc_004', 'team_manager', '王强', 'e10adc3949ba59abbe56e057f20f883e', '13800138003', 'team_manager', 'enabled',
 4, 'shanghai_branch', '上海分公司组', 1705123200, 1705123200, 20,
 4, 'tech_team_c', '技术团队C', 1705123200, 1705123200, 10,
 1705123200, 1705123200, 1705123200, -1, 'system', 'system'),

-- 普通员工
('acc_005', 'employee', '赵敏', 'e10adc3949ba59abbe56e057f20f883e', '13800138004', 'employee', 'enabled',
 5, 'guangzhou_office', '广州办公室组', 1705123200, 1705123200, 8,
 5, 'support_team_d', '客服团队D', 1705123200, 1705123200, 4,
 1705123200, 1705123200, 1705123200, -1, 'system', 'system');

-- 插入历史记录数据
INSERT INTO `go_gin_api`.`account_history` (
  `history_id`, `account_id`, `operate_type`, `operate_timestamp`, `content`, `operator`, `operator_role_type`,
  `is_deleted`, `created_user`, `updated_user`
) VALUES
('hist_001', 'acc_001', 'created', 1705123200, 
 '{"username":{"old":"","new":"admin"},"nickname":{"old":"","new":"系统管理员"},"roleType":{"old":"","new":"company_manager"}}',
 'system', 'admin', -1, 'system', 'system'),

('hist_002', 'acc_002', 'created', 1705123200,
 '{"username":{"old":"","new":"company_manager"},"nickname":{"old":"","new":"张伟"},"roleType":{"old":"","new":"company_manager"}}',
 'system', 'admin', -1, 'system', 'system'),

('hist_003', 'acc_003', 'created', 1705123200,
 '{"username":{"old":"","new":"group_manager"},"nickname":{"old":"","new":"李娜"},"roleType":{"old":"","new":"group_manager"}}',
 'system', 'admin', -1, 'system', 'system'),

('hist_004', 'acc_004', 'created', 1705123200,
 '{"username":{"old":"","new":"team_manager"},"nickname":{"old":"","new":"王强"},"roleType":{"old":"","new":"team_manager"}}',
 'system', 'admin', -1, 'system', 'system'),

('hist_005', 'acc_005', 'created', 1705123200,
 '{"username":{"old":"","new":"employee"},"nickname":{"old":"","new":"赵敏"},"roleType":{"old":"","new":"employee"}}',
 'system', 'admin', -1, 'system', 'system');

-- 显示创建结果
SELECT '账户表创建完成' as message;
SELECT COUNT(*) as account_count FROM `go_gin_api`.`account`;
SELECT COUNT(*) as history_count FROM `go_gin_api`.`account_history`; 