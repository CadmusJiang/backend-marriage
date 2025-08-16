package tablesqls

//CREATE TABLE `org_history` (
//`id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
//`org_id` int unsigned NOT NULL COMMENT '组织ID',
//`org_type` enum('group','team') NOT NULL COMMENT '组织类型 group:group team:team',
//`operate_type` varchar(20) NOT NULL COMMENT '操作类型: created, updated',
//`operated_at` timestamp NOT NULL COMMENT '操作时间',
//`content` json DEFAULT NULL COMMENT '操作内容',
//`operator` varchar(60) NOT NULL COMMENT '操作人',
//`operator_role_type` varchar(20) NOT NULL COMMENT '操作人角色',
//`created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
//`updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
//`created_user` varchar(60) NOT NULL COMMENT '创建人',
//`updated_user` varchar(60) NOT NULL COMMENT '更新人',
//PRIMARY KEY (`id`),
//KEY `idx_org_id` (`org_id`),
//KEY `idx_operated_at` (`operated_at`)
//) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='组织历史表';

func CreateOrgHistoryTableSql() (sql string) {
	sql = "CREATE TABLE `org_history` ("
	sql += "`id` int unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',"
	sql += "`org_id` int unsigned NOT NULL COMMENT '组织ID',"
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
	sql += "KEY `idx_org_id` (`org_id`),"
	sql += "KEY `idx_operated_at` (`operated_at`)"
	sql += ") ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='组织历史表';"
	return
}

func CreateOrgHistoryTableDataSql() (sql string) {
	// 默认不插入任何历史数据
	sql = ""
	return
}
