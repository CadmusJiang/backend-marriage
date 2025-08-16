package tablesqls

//CREATE TABLE `cooperation_store_history` (
//`id` int unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
//`store_id` int unsigned NOT NULL COMMENT '合作门店ID',
//`operate_type` varchar(20) NOT NULL COMMENT '操作类型: created, updated, deleted',
//`operated_at` timestamp NOT NULL COMMENT '操作时间',
//`content` json DEFAULT NULL COMMENT '操作内容',
//`operator_username` varchar(60) NOT NULL COMMENT '操作人用户名',
//`operator_role_type` varchar(20) NOT NULL COMMENT '操作人角色类型',
//`created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
//`updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
//`created_user` varchar(60) NOT NULL COMMENT '创建人',
//`updated_user` varchar(60) NOT NULL COMMENT '更新人',
//PRIMARY KEY (`id`),
//KEY `idx_store_id` (`store_id`),
//KEY `idx_operated_at` (`operated_at`)
//) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='合作门店历史记录表';

func CreateCooperationStoreHistoryTableSql() (sql string) {
	sql = "CREATE TABLE `cooperation_store_history` ("
	sql += "`id` int unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',"
	sql += "`store_id` int unsigned NOT NULL COMMENT '合作门店ID',"
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
	sql += "KEY `idx_store_id` (`store_id`),"
	sql += "KEY `idx_operated_at` (`operated_at`)"
	sql += ") ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='合作门店历史记录表';"
	return
}

func CreateCooperationStoreHistoryTableDataSql() (sql string) {
	sql = "INSERT INTO `cooperation_store_history` (`store_id`, `operate_type`, `operated_at`, `content`, `operator_username`, `operator_name`, `operator_role_type`, `created_at`, `updated_at`, `created_user`, `updated_user`) VALUES"
	sql += "(1, 'created', '2024-01-13 10:00:00', '{\"store_name\":\"上海婚恋门店1\",\"cooperation_city_code\":\"310000\"}', 'admin', '系统管理员', 'company_manager', '2024-01-13 10:00:00', '2024-01-13 10:00:00', '系统管理员', '系统管理员'),"
	sql += "(2, 'created', '2024-01-14 10:00:00', '{\"store_name\":\"北京婚恋门店2\",\"cooperation_city_code\":\"110000\"}', 'admin', '系统管理员', 'company_manager', '2024-01-14 10:00:00', '2024-01-14 10:00:00', '系统管理员', '系统管理员'),"
	sql += "(3, 'created', '2024-01-15 10:00:00', '{\"store_name\":\"深圳婚恋门店3\",\"cooperation_city_code\":\"440300\"}', 'admin', '系统管理员', 'company_manager', '2024-01-15 10:00:00', '2024-01-15 10:00:00', '系统管理员', '系统管理员');"
	return
}
