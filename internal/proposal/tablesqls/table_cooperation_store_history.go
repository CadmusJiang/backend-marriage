package tablesqls

//CREATE TABLE `cooperation_store_history` (
//`id` int unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
//`store_id` int unsigned NOT NULL COMMENT '合作门店ID',
//`operate_type` varchar(20) NOT NULL COMMENT '操作类型: created, modified, deleted',
//`occurred_at` bigint NOT NULL COMMENT '操作发生时间戳',
//`content` json DEFAULT NULL COMMENT '操作内容',
//`operator_username` varchar(60) NOT NULL COMMENT '操作人用户名',
//`operator_role_type` varchar(20) NOT NULL COMMENT '操作人角色类型',
//`created_at` bigint NOT NULL COMMENT '创建时间戳',
//`updated_at` bigint NOT NULL COMMENT '修改时间戳',
//`created_user` varchar(60) NOT NULL COMMENT '创建人',
//`updated_user` varchar(60) NOT NULL COMMENT '更新人',
//PRIMARY KEY (`id`),
//KEY `idx_store_id` (`store_id`),
//KEY `idx_occurred_at` (`occurred_at`)
//) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='合作门店历史记录表';

func CreateCooperationStoreHistoryTableSql() (sql string) {
	sql = "CREATE TABLE `cooperation_store_history` ("
	sql += "`id` int unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',"
	sql += "`store_id` int unsigned NOT NULL COMMENT '合作门店ID',"
	sql += "`operate_type` varchar(20) NOT NULL COMMENT '操作类型: created, modified, deleted',"
	sql += "`occurred_at` bigint NOT NULL COMMENT '操作发生时间戳',"
	sql += "`content` json DEFAULT NULL COMMENT '操作内容',"
	sql += "`operator_username` varchar(60) NOT NULL COMMENT '操作人用户名',"
	sql += "`operator_role_type` varchar(20) NOT NULL COMMENT '操作人角色类型',"
	sql += "`created_at` bigint NOT NULL COMMENT '创建时间戳',"
	sql += "`updated_at` bigint NOT NULL COMMENT '修改时间戳',"
	sql += "`created_user` varchar(60) NOT NULL COMMENT '创建人',"
	sql += "`updated_user` varchar(60) NOT NULL COMMENT '更新人',"
	sql += "PRIMARY KEY (`id`),"
	sql += "KEY `idx_store_id` (`store_id`),"
	sql += "KEY `idx_occurred_at` (`occurred_at`)"
	sql += ") ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='合作门店历史记录表';"
	return
}
