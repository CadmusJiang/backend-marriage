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