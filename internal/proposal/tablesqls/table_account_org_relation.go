package tablesqls

//CREATE TABLE `account_org_relation` (
//`id` int unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
//`account_id` int unsigned NOT NULL COMMENT '账户ID',
//`org_id` int unsigned NOT NULL COMMENT '组织ID',
//`relation_type` enum('belong','manage') NOT NULL COMMENT '关联类型 belong:belong manage:manage',
//`status` enum('active','inactive') NOT NULL DEFAULT 'active' COMMENT '状态 active:active inactive:inactive',
//`version` int NOT NULL DEFAULT 0 COMMENT '乐观锁版本号',
//`created_at` bigint NOT NULL COMMENT '创建时间戳',
//`updated_at` bigint NOT NULL COMMENT '修改时间戳',
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
	sql += "`relation_type` enum('belong','manage') NOT NULL COMMENT '关联类型 belong:belong manage:manage',"
	sql += "`status` enum('active','inactive') NOT NULL DEFAULT 'active' COMMENT '状态 active:active inactive:inactive',"
	sql += "`version` int NOT NULL DEFAULT 0 COMMENT '乐观锁版本号',"
	sql += "`created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',"
	sql += "`updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',"
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
	sql = "INSERT INTO `account_org_relation` (`account_id`, `org_id`, `relation_type`, `status`, `created_at`, `updated_at`, `created_user`, `updated_user`) VALUES"
	sql += "(1, 1, 'belong', 'active', '2024-01-13 10:00:00', '2024-01-13 10:00:00', '系统管理员', '系统管理员'),"
	sql += "(1, 4, 'belong', 'active', '2024-01-13 10:00:00', '2024-01-13 10:00:00', '系统管理员', '系统管理员'),"
	sql += "(2, 2, 'belong', 'active', '2024-01-13 10:00:00', '2024-01-13 10:00:00', '系统管理员', '系统管理员'),"
	sql += "(2, 5, 'belong', 'active', '2024-01-13 10:00:00', '2024-01-13 10:00:00', '系统管理员', '系统管理员'),"
	sql += "(3, 2, 'belong', 'active', '2024-01-13 10:00:00', '2024-01-13 10:00:00', '系统管理员', '系统管理员'),"
	sql += "(4, 2, 'belong', 'active', '2024-01-13 10:00:00', '2024-01-13 10:00:00', '系统管理员', '系统管理员'),"
	sql += "(4, 5, 'belong', 'active', '2024-01-13 10:00:00', '2024-01-13 10:00:00', '系统管理员', '系统管理员'),"
	sql += "(5, 3, 'belong', 'active', '2024-01-13 10:00:00', '2024-01-13 10:00:00', '系统管理员', '系统管理员'),"
	sql += "(5, 6, 'belong', 'active', '2024-01-13 10:00:00', '2024-01-13 10:00:00', '系统管理员', '系统管理员'),"
	sql += "(6, 3, 'belong', 'active', '2024-01-13 10:00:00', '2024-01-13 10:00:00', '系统管理员', '系统管理员'),"
	sql += "(6, 7, 'belong', 'active', '2024-01-13 10:00:00', '2024-01-13 10:00:00', '系统管理员', '系统管理员'),"
	sql += "(7, 2, 'belong', 'active', '2024-01-13 10:00:00', '2024-01-13 10:00:00', '系统管理员', '系统管理员'),"
	sql += "(7, 5, 'belong', 'active', '2024-01-13 10:00:00', '2024-01-13 10:00:00', '系统管理员', '系统管理员'),"
	sql += "(8, 3, 'belong', 'active', '2024-01-13 10:00:00', '2024-01-13 10:00:00', '系统管理员', '系统管理员'),"
	sql += "(8, 6, 'belong', 'active', '2024-01-13 10:00:00', '2024-01-13 10:00:00', '系统管理员', '系统管理员'),"
	sql += "(9, 3, 'belong', 'active', '2024-01-13 10:00:00', '2024-01-13 10:00:00', '系统管理员', '系统管理员'),"
	sql += "(9, 7, 'belong', 'active', '2024-01-13 10:00:00', '2024-01-13 10:00:00', '系统管理员', '系统管理员'),"
	sql += "(10, 3, 'belong', 'active', '2024-01-13 10:00:00', '2024-01-13 10:00:00', '系统管理员', '系统管理员'),"
	sql += "(10, 6, 'belong', 'active', '2024-01-13 10:00:00', '2024-01-13 10:00:00', '系统管理员', '系统管理员'),"
	sql += "(11, 2, 'belong', 'active', '2024-01-13 10:00:00', '2024-01-13 10:00:00', '系统管理员', '系统管理员'),"
	sql += "(11, 5, 'belong', 'active', '2024-01-13 10:00:00', '2024-01-13 10:00:00', '系统管理员', '系统管理员');"

	return
}
