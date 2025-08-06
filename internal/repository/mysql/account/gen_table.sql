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