-- Marriage System 数据库初始化脚本
-- 数据库名称: 从配置文件读取

-- 注意: 此脚本需要在运行时替换数据库名称
-- 请使用以下命令运行:
-- sed 's/__DB_NAME__/marriage_system/g' scripts/init_marriage_system_database.sql | mysql -h HOST -P PORT -u USER -pPASS

-- 创建数据库
CREATE DATABASE IF NOT EXISTS `__DB_NAME__` 
DEFAULT CHARACTER SET utf8mb4 
COLLATE utf8mb4_unicode_ci;

USE `__DB_NAME__`;

-- 删除旧表（如果存在）
DROP TABLE IF EXISTS `account_org_relation`;
DROP TABLE IF EXISTS `account_history`;
DROP TABLE IF EXISTS `account`;
DROP TABLE IF EXISTS `org_info`;

-- 创建组织信息表
CREATE TABLE `org_info` (
  `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
  `org_name` varchar(60) NOT NULL COMMENT '组织名称',
  `org_type` tinyint(1) NOT NULL COMMENT '组织类型: 1-group, 2-team',
  `org_path` varchar(255) NOT NULL COMMENT '组织路径',
  `org_level` tinyint(1) NOT NULL DEFAULT '1' COMMENT '组织层级: 1-组, 2-团队',
  `current_cnt` int NOT NULL DEFAULT '0' COMMENT '当前成员数量',
  `max_cnt` int NOT NULL DEFAULT '0' COMMENT '最大成员数量',
  `status` tinyint(1) NOT NULL DEFAULT '1' COMMENT '状态: 1-启用, 0-禁用',
  `ext_data` json DEFAULT NULL COMMENT '扩展数据 (JSON格式)',
  `created_timestamp` bigint NOT NULL COMMENT '创建时间戳',
  `modified_timestamp` bigint NOT NULL COMMENT '修改时间戳',
  `created_user` varchar(60) NOT NULL COMMENT '创建人',
  `updated_user` varchar(60) NOT NULL COMMENT '更新人',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_org_name` (`org_name`),
  UNIQUE KEY `uk_org_path` (`org_path`),
  KEY `idx_org_type` (`org_type`),
  KEY `idx_org_level` (`org_level`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='组织信息表';

-- 创建账户表（不包含组织关联字段）
CREATE TABLE `account` (
  `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
  `account_id` varchar(32) NOT NULL COMMENT '账户ID',
  `username` varchar(32) NOT NULL COMMENT '用户名',
  `nickname` varchar(60) NOT NULL COMMENT '姓名',
  `password` varchar(100) NOT NULL COMMENT '密码',
  `phone` varchar(20) DEFAULT NULL COMMENT '手机号',
  `role_type` varchar(20) NOT NULL DEFAULT 'employee' COMMENT '角色类型: company_manager, group_manager, team_manager, employee',
  `status` varchar(20) NOT NULL DEFAULT 'enabled' COMMENT '状态: enabled, disabled',
  `last_login_timestamp` bigint DEFAULT NULL COMMENT '最后登录时间戳',
  `created_timestamp` bigint NOT NULL COMMENT '创建时间戳',
  `modified_timestamp` bigint NOT NULL COMMENT '修改时间戳',
  `created_user` varchar(60) NOT NULL COMMENT '创建人',
  `updated_user` varchar(60) NOT NULL COMMENT '更新人',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_account_id` (`account_id`),
  UNIQUE KEY `uk_username` (`username`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='账户表';

-- 创建账户组织关联表
CREATE TABLE `account_org_relation` (
  `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
  `account_id` varchar(32) NOT NULL COMMENT '账户ID',
  `org_id` varchar(32) NOT NULL COMMENT '组织ID',
  `relation_type` varchar(20) NOT NULL COMMENT '关联类型: group, team',
  `created_timestamp` bigint NOT NULL COMMENT '创建时间戳',
  `modified_timestamp` bigint NOT NULL COMMENT '修改时间戳',
  `created_user` varchar(60) NOT NULL COMMENT '创建人',
  `updated_user` varchar(60) NOT NULL COMMENT '更新人',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_account_org` (`account_id`, `org_id`, `relation_type`),
  KEY `idx_account_id` (`account_id`),
  KEY `idx_org_id` (`org_id`),
  KEY `idx_relation_type` (`relation_type`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='账户组织关联表';

-- 创建账户历史记录表
CREATE TABLE `account_history` (
  `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
  `history_id` varchar(32) NOT NULL COMMENT '历史记录ID',
  `account_id` varchar(32) NOT NULL COMMENT '账户ID',
  `operate_type` varchar(20) NOT NULL COMMENT '操作类型: created, modified, deleted',
  `operate_timestamp` bigint NOT NULL COMMENT '操作时间戳',
  `content` json DEFAULT NULL COMMENT '操作内容 (JSON格式)',
  `operator_username` varchar(32) NOT NULL COMMENT '操作人用户名',
  `operator_nickname` varchar(60) NOT NULL COMMENT '操作人姓名',
  `operator_role_type` varchar(20) NOT NULL COMMENT '操作人角色',
  `created_timestamp` bigint NOT NULL COMMENT '创建时间戳',
  `modified_timestamp` bigint NOT NULL COMMENT '修改时间戳',
  `created_user` varchar(60) NOT NULL COMMENT '创建人',
  `updated_user` varchar(60) NOT NULL COMMENT '更新人',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_history_id` (`history_id`),
  KEY `idx_account_id` (`account_id`),
  KEY `idx_operate_type` (`operate_type`),
  KEY `idx_operate_timestamp` (`operate_timestamp`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='账户历史记录表';

-- 插入组织数据
INSERT INTO `org_info` (`org_name`, `org_type`, `org_path`, `org_level`, `current_cnt`, `max_cnt`, `status`, `created_timestamp`, `modified_timestamp`, `created_user`, `updated_user`) VALUES
('系统管理组', 1, '/system/admin', 1, 1, 10, 1, 1705123200, 1705123200, 'system', 'system'),
('南京-天元大厦组', 1, '/nanjing/tianyuan', 1, 15, 50, 1, 1705123200, 1705123200, 'system', 'system'),
('北京办公室组', 1, '/beijing/office', 1, 12, 40, 1, 1705123200, 1705123200, 'system', 'system'),
('上海分公司组', 1, '/shanghai/branch', 1, 18, 60, 1, 1705123200, 1705123200, 'system', 'system'),
('广州办公室组', 1, '/guangzhou/office', 1, 8, 30, 1, 1705123200, 1705123200, 'system', 'system'),
('系统管理团队', 2, '/system/admin/team', 2, 1, 5, 1, 1705123200, 1705123200, 'system', 'system'),
('营销团队A', 2, '/nanjing/tianyuan/marketing_a', 2, 8, 20, 1, 1705123200, 1705123200, 'system', 'system'),
('销售团队B', 2, '/beijing/office/sales_b', 2, 10, 25, 1, 1705123200, 1705123200, 'system', 'system'),
('技术团队C', 2, '/shanghai/branch/tech_c', 2, 6, 15, 1, 1705123200, 1705123200, 'system', 'system'),
('客服团队D', 2, '/guangzhou/office/service_d', 2, 4, 12, 1, 1705123200, 1705123200, 'system', 'system');

-- 插入账户数据
INSERT INTO `account` (`account_id`, `username`, `nickname`, `password`, `phone`, `role_type`, `status`, `last_login_timestamp`, `created_timestamp`, `modified_timestamp`, `created_user`, `updated_user`) VALUES
('acc_001', 'admin', '系统管理员', 'e10adc3949ba59abbe56e057f20f883e', '13800138000', 'company_manager', 'enabled', 1705123200, 1705123200, 1705123200, 'system', 'system'),
('acc_002', 'company_manager', '张伟', 'e10adc3949ba59abbe56e057f20f883e', '13800138001', 'company_manager', 'enabled', 1705123200, 1705123200, 1705123200, 'system', 'system'),
('acc_003', 'group_manager', '李娜', 'e10adc3949ba59abbe56e057f20f883e', '13800138002', 'group_manager', 'enabled', 1705123200, 1705123200, 1705123200, 'system', 'system'),
('acc_004', 'team_manager', '王强', 'e10adc3949ba59abbe56e057f20f883e', '13800138003', 'team_manager', 'enabled', 1705123200, 1705123200, 1705123200, 'system', 'system'),
('acc_005', 'employee', '赵敏', 'e10adc3949ba59abbe56e057f20f883e', '13800138004', 'employee', 'enabled', 1705123200, 1705123200, 1705123200, 'system', 'system');

-- 插入账户组织关联数据
INSERT INTO `account_org_relation` (`account_id`, `org_id`, `relation_type`, `created_timestamp`, `modified_timestamp`, `created_user`, `updated_user`) VALUES
-- admin 关联系统管理组和系统管理团队
('acc_001', '1', 'group', 1705123200, 1705123200, 'system', 'system'),
('acc_001', '6', 'team', 1705123200, 1705123200, 'system', 'system'),
-- company_manager 关联南京-天元大厦组和营销团队A
('acc_002', '2', 'group', 1705123200, 1705123200, 'system', 'system'),
('acc_002', '7', 'team', 1705123200, 1705123200, 'system', 'system'),
-- group_manager 关联北京办公室组和销售团队B
('acc_003', '3', 'group', 1705123200, 1705123200, 'system', 'system'),
('acc_003', '8', 'team', 1705123200, 1705123200, 'system', 'system'),
-- team_manager 关联上海分公司组和技术团队C
('acc_004', '4', 'group', 1705123200, 1705123200, 'system', 'system'),
('acc_004', '9', 'team', 1705123200, 1705123200, 'system', 'system'),
-- employee 关联广州办公室组和客服团队D
('acc_005', '5', 'group', 1705123200, 1705123200, 'system', 'system'),
('acc_005', '10', 'team', 1705123200, 1705123200, 'system', 'system');

-- 插入账户历史记录数据
INSERT INTO `account_history` (`history_id`, `account_id`, `operate_type`, `operate_timestamp`, `content`, `operator_username`, `operator_nickname`, `operator_role_type`, `created_timestamp`, `modified_timestamp`, `created_user`, `updated_user`) VALUES
('hist_001', 'acc_001', 'created', 1705123200, '{"action": "create", "details": "系统管理员账户创建"}', 'system', '系统', 'company_manager', 1705123200, 1705123200, 'system', 'system'),
('hist_002', 'acc_002', 'created', 1705123200, '{"action": "create", "details": "公司管理员账户创建"}', 'system', '系统', 'company_manager', 1705123200, 1705123200, 'system', 'system'),
('hist_003', 'acc_003', 'created', 1705123200, '{"action": "create", "details": "组管理员账户创建"}', 'system', '系统', 'company_manager', 1705123200, 1705123200, 'system', 'system'),
('hist_004', 'acc_004', 'created', 1705123200, '{"action": "create", "details": "团队管理员账户创建"}', 'system', '系统', 'company_manager', 1705123200, 1705123200, 'system', 'system'),
('hist_005', 'acc_005', 'created', 1705123200, '{"action": "create", "details": "员工账户创建"}', 'system', '系统', 'company_manager', 1705123200, 1705123200, 'system', 'system');

-- 显示统计信息
SELECT 'Marriage System 数据库初始化完成' as message;
SELECT COUNT(*) as org_count FROM org_info;
SELECT COUNT(*) as account_count FROM account;
SELECT COUNT(*) as relation_count FROM account_org_relation;
SELECT COUNT(*) as history_count FROM account_history; 