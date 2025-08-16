package tablesqls

//CREATE TABLE `customer_authorization_record_history` (
//`id` int unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
//`history_id` varchar(64) NOT NULL COMMENT '历史记录ID',
//`customer_authorization_record_id` varchar(64) NOT NULL COMMENT '客户授权记录ID',
//`operate_type` varchar(20) NOT NULL COMMENT '操作类型: created, updated, deleted',
//`operated_at` timestamp NOT NULL COMMENT '操作时间',
//`content` json DEFAULT NULL COMMENT '操作内容',
//`operator_username` varchar(32) NOT NULL COMMENT '操作人用户名',
//`operator_name` varchar(60) NOT NULL COMMENT '操作人姓名',
//`operator_role_type` varchar(20) NOT NULL COMMENT '操作人角色类型',
//`is_deleted` enum('true','false') NOT NULL DEFAULT 'false' COMMENT '是否删除 true:是 false:否',
//`created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
//`updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
//`created_user` varchar(60) NOT NULL COMMENT '创建人',
//`updated_user` varchar(60) NOT NULL COMMENT '更新人',
//PRIMARY KEY (`id`),
//UNIQUE KEY `uk_history_id` (`history_id`),
//KEY `idx_customer_authorization_record_id` (`customer_authorization_record_id`),
//KEY `idx_operated_at` (`operated_at`)
//) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='客户授权记录历史记录表';

func CreateCustomerAuthorizationRecordHistoryTableSql() (sql string) {
	sql = "CREATE TABLE `customer_authorization_record_history` ("
	sql += "`id` int unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',"
	sql += "`customer_authorization_record_id` int unsigned NOT NULL COMMENT '客户授权记录ID',"
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
	sql += "KEY `idx_customer_authorization_record_id` (`customer_authorization_record_id`),"
	sql += "KEY `idx_operated_at` (`operated_at`)"
	sql += ") ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='客户授权记录历史记录表';"

	return
}

func CreateCustomerAuthorizationRecordHistoryTableDataSql() (sql string) {
	// 预置客户授权记录历史数据
	sql = "INSERT INTO `customer_authorization_record_history` (`customer_authorization_record_id`, `operate_type`, `operated_at`, `content`, `operator_username`, `operator_name`, `operator_role_type`, `created_user`, `updated_user`) VALUES"
	sql += "(1, 'created', '2024-01-13 10:00:00', '{\"customerName\": \"张三\", \"phone\": \"13800138000\", \"authorizationType\": \"full\"}', 'admin', '系统管理员', 'company_manager', '系统管理员', '系统管理员'),"
	sql += "(2, 'created', '2024-01-13 10:00:00', '{\"customerName\": \"李四\", \"phone\": \"13800138001\", \"authorizationType\": \"partial\"}', 'employee001', '张三', 'employee', '张三', '张三');"

	return
}
