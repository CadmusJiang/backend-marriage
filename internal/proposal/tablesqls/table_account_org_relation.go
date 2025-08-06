package tablesqls

//CREATE TABLE `account_org_relation` (
//`id` int unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
//`account_id` int unsigned NOT NULL COMMENT '账户ID',
//`org_id` int unsigned NOT NULL COMMENT '组织ID',
//`relation_type` tinyint(1) NOT NULL COMMENT '关联类型 1:belong 2:manage',
//`status` tinyint(1) NOT NULL DEFAULT 1 COMMENT '状态 1:active 2:inactive',
//`created_timestamp` bigint NOT NULL COMMENT '创建时间戳',
//`modified_timestamp` bigint NOT NULL COMMENT '修改时间戳',
//`created_user` varchar(60) NOT NULL COMMENT '创建人',
//`updated_user` varchar(60) NOT NULL COMMENT '更新人',
//PRIMARY KEY (`id`),
//UNIQUE KEY `uk_account_org` (`account_id`, `org_id`),
//KEY `idx_account_id` (`account_id`),
//KEY `idx_org_id` (`org_id`)
//) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='账户组织关联表';

func CreateAccountOrgRelationTableSql() (sql string) {
	sql = "CREATE TABLE `account_org_relation` ("
	sql += "`id` int unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',"
	sql += "`account_id` int unsigned NOT NULL COMMENT '账户ID',"
	sql += "`org_id` int unsigned NOT NULL COMMENT '组织ID',"
	sql += "`relation_type` tinyint(1) NOT NULL COMMENT '关联类型 1:belong 2:manage',"
	sql += "`status` tinyint(1) NOT NULL DEFAULT 1 COMMENT '状态 1:active 2:inactive',"
	sql += "`created_timestamp` bigint NOT NULL COMMENT '创建时间戳',"
	sql += "`modified_timestamp` bigint NOT NULL COMMENT '修改时间戳',"
	sql += "`created_user` varchar(60) NOT NULL COMMENT '创建人',"
	sql += "`updated_user` varchar(60) NOT NULL COMMENT '更新人',"
	sql += "PRIMARY KEY (`id`),"
	sql += "UNIQUE KEY `uk_account_org` (`account_id`, `org_id`),"
	sql += "KEY `idx_account_id` (`account_id`),"
	sql += "KEY `idx_org_id` (`org_id`)"
	sql += ") ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='账户组织关联表';"

	return
}

func CreateAccountOrgRelationTableDataSql() (sql string) {
	sql = "INSERT INTO `account_org_relation` (`account_id`, `org_id`, `relation_type`, `status`, `created_timestamp`, `modified_timestamp`, `created_user`, `updated_user`) VALUES"
	sql += "(1, 1, 1, 1, 1705123200, 1705123200, '系统管理员', '系统管理员'),"
	sql += "(1, 4, 1, 1, 1705123200, 1705123200, '系统管理员', '系统管理员'),"
	sql += "(2, 2, 1, 1, 1705123200, 1705123200, '系统管理员', '系统管理员'),"
	sql += "(2, 5, 1, 1, 1705123200, 1705123200, '系统管理员', '系统管理员'),"
	sql += "(3, 2, 1, 1, 1705123200, 1705123200, '系统管理员', '系统管理员'),"
	sql += "(4, 2, 1, 1, 1705123200, 1705123200, '系统管理员', '系统管理员'),"
	sql += "(4, 5, 1, 1, 1705123200, 1705123200, '系统管理员', '系统管理员'),"
	sql += "(5, 3, 1, 1, 1705123200, 1705123200, '系统管理员', '系统管理员'),"
	sql += "(5, 6, 1, 1, 1705123200, 1705123200, '系统管理员', '系统管理员'),"
	sql += "(6, 2, 1, 1, 1705123200, 1705123200, '系统管理员', '系统管理员'),"
	sql += "(6, 7, 1, 1, 1705123200, 1705123200, '系统管理员', '系统管理员'),"
	sql += "(7, 2, 1, 1, 1705123200, 1705123200, '系统管理员', '系统管理员'),"
	sql += "(7, 5, 1, 1, 1705123200, 1705123200, '系统管理员', '系统管理员'),"
	sql += "(8, 3, 1, 1, 1705123200, 1705123200, '系统管理员', '系统管理员'),"
	sql += "(8, 6, 1, 1, 1705123200, 1705123200, '系统管理员', '系统管理员'),"
	sql += "(9, 2, 1, 1, 1705123200, 1705123200, '系统管理员', '系统管理员'),"
	sql += "(9, 7, 1, 1, 1705123200, 1705123200, '系统管理员', '系统管理员'),"
	sql += "(10, 3, 1, 1, 1705123200, 1705123200, '系统管理员', '系统管理员'),"
	sql += "(10, 6, 1, 1, 1705123200, 1705123200, '系统管理员', '系统管理员'),"
	sql += "(11, 2, 1, 1, 1705123200, 1705123200, '系统管理员', '系统管理员'),"
	sql += "(11, 5, 1, 1, 1705123200, 1705123200, '系统管理员', '系统管理员');"

	return
}
