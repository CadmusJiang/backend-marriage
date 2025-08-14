package tablesqls

//CREATE TABLE `org_history` (
//`id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
//`org_id` int unsigned NOT NULL COMMENT '组织ID',
//`org_type` tinyint(1) NOT NULL COMMENT '组织类型 1:group 2:team',
//`operate_type` varchar(20) NOT NULL COMMENT '操作类型: created, modified',
//`occurred_at` bigint NOT NULL COMMENT '操作发生时间戳',
//`content` json DEFAULT NULL COMMENT '操作内容',
//`operator` varchar(60) NOT NULL COMMENT '操作人',
//`operator_role_type` varchar(20) NOT NULL COMMENT '操作人角色',
//`created_at` bigint NOT NULL COMMENT '创建时间戳',
//`updated_at` bigint NOT NULL COMMENT '修改时间戳',
//`created_user` varchar(60) NOT NULL COMMENT '创建人',
//`updated_user` varchar(60) NOT NULL COMMENT '更新人',
//PRIMARY KEY (`id`),
//KEY `idx_org_id` (`org_id`),
//KEY `idx_operate_timestamp` (`operate_timestamp`)
//) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='组织历史表';

func CreateOrgHistoryTableSql() (sql string) {
	sql = "CREATE TABLE `org_history` ("
	sql += "`id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',"
	sql += "`org_id` int unsigned NOT NULL COMMENT '组织ID',"
	sql += "`org_type` tinyint(1) NOT NULL COMMENT '组织类型 1:group 2:team',"
	sql += "`operate_type` varchar(20) NOT NULL COMMENT '操作类型: created, modified',"
	sql += "`occurred_at` bigint NOT NULL COMMENT '操作发生时间戳',"
	sql += "`content` json DEFAULT NULL COMMENT '操作内容',"
	sql += "`operator` varchar(60) NOT NULL COMMENT '操作人',"
	sql += "`operator_role_type` varchar(20) NOT NULL COMMENT '操作人角色',"
	sql += "`created_at` bigint NOT NULL COMMENT '创建时间戳',"
	sql += "`updated_at` bigint NOT NULL COMMENT '修改时间戳',"
	sql += "`created_user` varchar(60) NOT NULL COMMENT '创建人',"
	sql += "`updated_user` varchar(60) NOT NULL COMMENT '更新人',"
	sql += "PRIMARY KEY (`id`),"
	sql += "KEY `idx_org_id` (`org_id`),"
	sql += "KEY `idx_operate_timestamp` (`operate_timestamp`)"
	sql += ") ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='组织历史表';"
	return
}

func CreateOrgHistoryTableDataSql() (sql string) {
	// 默认不插入任何历史数据
	sql = ""
	return
}
