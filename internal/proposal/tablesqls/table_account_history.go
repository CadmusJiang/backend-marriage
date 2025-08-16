package tablesqls

//CREATE TABLE `account_history` (
//`id` int unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
//`account_id` int unsigned NOT NULL COMMENT '账户ID',
//`operate_type` varchar(20) NOT NULL COMMENT '操作类型: created, updated',
//`operated_at` timestamp NOT NULL COMMENT '操作时间',
//`content` text NOT NULL COMMENT '操作内容 (JSON格式)',
//`operator_username` varchar(32) NOT NULL COMMENT '操作人用户名',
//`operator_name` varchar(60) NOT NULL COMMENT '操作人姓名',
//`operator_role_type` varchar(20) NOT NULL COMMENT '操作人角色类型',
//`created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
//`updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
//`created_user` varchar(60) NOT NULL COMMENT '创建人',
//`updated_user` varchar(60) NOT NULL COMMENT '更新人',
//PRIMARY KEY (`id`),
//KEY `idx_account_id` (`account_id`),
//KEY `idx_operated_at` (`operated_at`)
//) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='账户历史记录表';

func CreateAccountHistoryTableSql() (sql string) {
	sql = "CREATE TABLE `account_history` ("
	sql += "`id` int unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',"
	sql += "`account_id` int unsigned NOT NULL COMMENT '账户ID',"
	sql += "`operate_type` varchar(20) NOT NULL COMMENT '操作类型: created, updated',"
	sql += "`operated_at` timestamp NOT NULL COMMENT '操作时间',"
	sql += "`content` json DEFAULT NULL COMMENT '操作内容',"
	sql += "`operator_username` varchar(32) NOT NULL COMMENT '操作人用户名',"
	sql += "`operator_name` varchar(60) NOT NULL COMMENT '操作人姓名',"
	sql += "`operator_role_type` varchar(20) NOT NULL COMMENT '操作人角色类型',"
	sql += "`created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',"
	sql += "`updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',"
	sql += "`created_user` varchar(60) NOT NULL COMMENT '创建人',"
	sql += "`updated_user` varchar(60) NOT NULL COMMENT '更新人',"
	sql += "PRIMARY KEY (`id`),"
	sql += "KEY `idx_account_id` (`account_id`),"
	sql += "KEY `idx_operated_at` (`operated_at`)"
	sql += ") ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='账户历史记录表';"

	return
}

func CreateAccountHistoryTableDataSql() (sql string) {
	// 预置账户历史记录数据
	sql = "INSERT INTO `account_history` (`id`, `account_id`, `operate_type`, `operated_at`, `content`, `operator_username`, `operator_name`, `operator_role_type`, `created_user`, `updated_user`) VALUES"
	sql += "(1, 2, 'updated', '2024-01-12 10:00:00', '{\"name\": {\"old\": \"张三\", \"new\": \"张明\"}, \"roleType\": {\"old\": \"员工\", \"new\": \"小队管理员\"}}', 'chenjing', '陈静', 'team_manager', '陈静', '陈静'),"
	sql += "(2, 2, 'created', '2024-01-10 10:00:00', '{\"username\": {\"old\": \"\", \"new\": \"employee001\"}, \"name\": {\"old\": \"\", \"new\": \"张三\"}, \"roleType\": {\"old\": \"\", \"new\": \"员工\"}, \"phone\": {\"old\": \"\", \"new\": \"13800138000\"}, \"belongGroup\": {\"old\": \"\", \"new\": \"南京-天元大厦组\"}, \"belongTeam\": {\"old\": \"\", \"new\": \"营销团队A\"}}', 'admin', '系统管理员', 'company_manager', '系统管理员', '系统管理员');"

	return
}
