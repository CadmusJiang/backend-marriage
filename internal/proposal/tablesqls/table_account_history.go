package tablesqls

//CREATE TABLE `account_history` (
//`id` int unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
//`account_id` int unsigned NOT NULL COMMENT '账户ID',
//`operate_type` varchar(20) NOT NULL COMMENT '操作类型',
//`operate_timestamp` bigint NOT NULL COMMENT '操作时间戳',
//`content` json DEFAULT NULL COMMENT '操作内容',
//`operator_username` varchar(32) NOT NULL COMMENT '操作人用户名',
//`operator_nickname` varchar(60) NOT NULL COMMENT '操作人姓名',
//`operator_role_type` varchar(20) NOT NULL COMMENT '操作人角色',
//`created_timestamp` bigint NOT NULL COMMENT '创建时间戳',
//`modified_timestamp` bigint NOT NULL COMMENT '修改时间戳',
//`created_user` varchar(60) NOT NULL COMMENT '创建人',
//`updated_user` varchar(60) NOT NULL COMMENT '更新人',
//PRIMARY KEY (`id`),
//KEY `idx_account_id` (`account_id`)
//) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='账户历史记录表';

func CreateAccountHistoryTableSql() (sql string) {
	sql = "CREATE TABLE `account_history` ("
	sql += "`id` int unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',"
	sql += "`account_id` int unsigned NOT NULL COMMENT '账户ID',"
	sql += "`operate_type` varchar(20) NOT NULL COMMENT '操作类型',"
	sql += "`operate_timestamp` bigint NOT NULL COMMENT '操作时间戳',"
	sql += "`content` json DEFAULT NULL COMMENT '操作内容',"
	sql += "`operator_username` varchar(32) NOT NULL COMMENT '操作人用户名',"
	sql += "`operator_nickname` varchar(60) NOT NULL COMMENT '操作人姓名',"
	sql += "`operator_role_type` varchar(20) NOT NULL COMMENT '操作人角色',"
	sql += "`created_timestamp` bigint NOT NULL COMMENT '创建时间戳',"
	sql += "`modified_timestamp` bigint NOT NULL COMMENT '修改时间戳',"
	sql += "`created_user` varchar(60) NOT NULL COMMENT '创建人',"
	sql += "`updated_user` varchar(60) NOT NULL COMMENT '更新人',"
	sql += "PRIMARY KEY (`id`),"
	sql += "KEY `idx_account_id` (`account_id`)"
	sql += ") ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='账户历史记录表';"

	return
}

func CreateAccountHistoryTableDataSql() (sql string) {
	sql = "INSERT INTO `account_history` (`account_id`, `operate_type`, `operate_timestamp`, `content`, `operator_username`, `operator_nickname`, `operator_role_type`, `created_timestamp`, `modified_timestamp`, `created_user`, `updated_user`) VALUES"
	sql += "(2, 'modified', 1705923000, '{\"roleType\": {\"old\": \"员工\", \"new\": \"小队管理员\"}, \"belongTeam\": {\"old\": \"无\", \"new\": \"营销团队A\"}, \"status\": {\"old\": \"enabled\", \"new\": \"disabled\"}}', 'zhangwei', '张伟', 'company_manager', 1705923000, 1705923000, '张伟', '张伟'),"
	sql += "(2, 'modified', 1705754700, '{\"belongGroup\": {\"old\": \"南京-天元大厦组\", \"new\": \"南京-南京南站组\"}, \"belongTeam\": {\"old\": \"营销团队A\", \"new\": \"营销团队C\"}}', 'liming', '李明', 'team_manager', 1705754700, 1705754700, '李明', '李明'),"
	sql += "(2, 'modified', 1705565700, '{\"status\": {\"old\": \"enabled\", \"new\": \"disabled\"}}', 'wangfang', '王芳', 'team_manager', 1705565700, 1705565700, '王芳', '王芳'),"
	sql += "(2, 'modified', 1705303200, '{\"phone\": {\"old\": \"13800138000\", \"new\": \"13900139000\"}, \"belongTeam\": {\"old\": \"营销团队A\", \"new\": \"营销团队B\"}}', 'liuqiang', '刘强', 'team_manager', 1705303200, 1705303200, '刘强', '刘强'),"
	sql += "(2, 'modified', 1705032000, '{\"nickname\": {\"old\": \"张三\", \"new\": \"张明\"}, \"roleType\": {\"old\": \"员工\", \"new\": \"小队管理员\"}}', 'chenjing', '陈静', 'team_manager', 1705032000, 1705032000, '陈静', '陈静'),"
	sql += "(2, 'created', 1704877800, '{\"username\": {\"old\": \"\", \"new\": \"employee001\"}, \"nickname\": {\"old\": \"\", \"new\": \"张三\"}, \"roleType\": {\"old\": \"\", \"new\": \"员工\"}, \"phone\": {\"old\": \"\", \"new\": \"13800138000\"}, \"belongGroup\": {\"old\": \"\", \"new\": \"南京-天元大厦组\"}, \"belongTeam\": {\"old\": \"\", \"new\": \"营销团队A\"}}', 'admin', '系统管理员', 'company_manager', 1704877800, 1704877800, '系统管理员', '系统管理员'),"
	sql += "(2, 'modified', 1704709200, '{\"belongGroup\": {\"old\": \"南京-天元大厦组\", \"new\": \"南京-夫子庙组\"}, \"belongTeam\": {\"old\": \"营销团队A\", \"new\": \"营销团队D\"}}', 'zhaomin', '赵敏', 'team_manager', 1704709200, 1704709200, '赵敏', '赵敏'),"
	sql += "(2, 'modified', 1704457500, '{\"status\": {\"old\": \"disabled\", \"new\": \"enabled\"}, \"reason\": {\"old\": \"\", \"new\": \"问题已解决，恢复账户\"}}', 'liming', '李明', 'team_manager', 1704457500, 1704457500, '李明', '李明'),"
	sql += "(2, 'modified', 1704265800, '{\"password\": {\"old\": \"****\", \"new\": \"******\"}}', 'zhangwei', '张伟', 'company_manager', 1704265800, 1704265800, '张伟', '张伟'),"
	sql += "(2, 'modified', 1704081600, '{\"password\": {\"old\": \"****\", \"new\": \"******\"}}', 'wangfang', '王芳', 'team_manager', 1704081600, 1704081600, '王芳', '王芳');"

	return
}
