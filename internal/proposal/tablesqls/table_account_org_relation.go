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
	sql = "INSERT INTO `account_org_relation` (`id`, `account_id`, `org_id`, `relation_type`, `status`, `created_at`, `updated_at`, `created_user`, `updated_user`) VALUES"

	// 注意：account_id现在对应实际的账户ID（1-11），确保与account表的新ID一致

	// company_manager (系统管理员, ID:1) - 与公司有manage关系
	sql += "(1, 1, 1001, 'manage', 'active', '2024-01-13 10:00:00', '2024-01-13 10:00:00', '系统管理员', '系统管理员'),"

	// group_manager (陈静, ID:5) - 与组有manage关系，与公司有belong关系
	sql += "(2, 5, 1001, 'belong', 'active', '2024-01-13 10:00:00', '2024-01-13 10:00:00', '系统管理员', '系统管理员'),"
	sql += "(3, 5, 10003, 'manage', 'active', '2024-01-13 10:00:00', '2024-01-13 10:00:00', '系统管理员', '系统管理员'),"

	// team_manager (赵六, ID:6) - 与组有belong关系，与团队有manage关系，与公司有belong关系
	sql += "(4, 6, 1001, 'belong', 'active', '2024-01-13 10:00:00', '2024-01-13 10:00:00', '系统管理员', '系统管理员'),"
	sql += "(5, 6, 10003, 'belong', 'active', '2024-01-13 10:00:00', '2024-01-13 10:00:00', '系统管理员', '系统管理员'),"
	sql += "(6, 6, 100004, 'manage', 'active', '2024-01-13 10:00:00', '2024-01-13 10:00:00', '系统管理员', '系统管理员'),"

	// employee (张三, ID:2) - 与组和团队有belong关系，与公司有belong关系
	sql += "(7, 2, 1001, 'belong', 'active', '2024-01-13 10:00:00', '2024-01-13 10:00:00', '系统管理员', '系统管理员'),"
	sql += "(8, 2, 10002, 'belong', 'active', '2024-01-13 10:00:00', '2024-01-13 10:00:00', '系统管理员', '系统管理员'),"
	sql += "(9, 2, 100002, 'belong', 'active', '2024-01-13 10:00:00', '2024-01-13 10:00:00', '系统管理员', '系统管理员'),"

	// employee (李四, ID:3) - 只有组关系，无团队关系，与公司有belong关系
	sql += "(10, 3, 1001, 'belong', 'active', '2024-01-13 10:00:00', '2024-01-13 10:00:00', '系统管理员', '系统管理员'),"
	sql += "(11, 3, 10002, 'belong', 'active', '2024-01-13 10:00:00', '2024-01-13 10:00:00', '系统管理员', '系统管理员'),"

	// employee (王五, ID:4) - 与组和团队有belong关系，与公司有belong关系
	sql += "(12, 4, 1001, 'belong', 'active', '2024-01-13 10:00:00', '2024-01-13 10:00:00', '系统管理员', '系统管理员'),"
	sql += "(13, 4, 10002, 'belong', 'active', '2024-01-13 10:00:00', '2024-01-13 10:00:00', '系统管理员', '系统管理员'),"
	sql += "(14, 4, 100002, 'belong', 'active', '2024-01-13 10:00:00', '2024-01-13 10:00:00', '系统管理员', '系统管理员'),"

	// employee (孙七, ID:7) - 与组和团队有belong关系，与公司有belong关系
	sql += "(15, 7, 1001, 'belong', 'active', '2024-01-13 10:00:00', '2024-01-13 10:00:00', '系统管理员', '系统管理员'),"
	sql += "(16, 7, 10002, 'belong', 'active', '2024-01-13 10:00:00', '2024-01-13 10:00:00', '系统管理员', '系统管理员'),"
	sql += "(17, 7, 100002, 'belong', 'active', '2024-01-13 10:00:00', '2024-01-13 10:00:00', '系统管理员', '系统管理员'),"

	// employee (周八, ID:8) - 与组和团队有belong关系，与公司有belong关系
	sql += "(18, 8, 1001, 'belong', 'active', '2024-01-13 10:00:00', '2024-01-13 10:00:00', '系统管理员', '系统管理员'),"
	sql += "(19, 8, 10003, 'belong', 'active', '2024-01-13 10:00:00', '2024-01-13 10:00:00', '系统管理员', '系统管理员'),"
	sql += "(20, 8, 100003, 'belong', 'active', '2024-01-13 10:00:00', '2024-01-13 10:00:00', '系统管理员', '系统管理员'),"

	// employee (吴九, ID:9) - 与组和团队有belong关系，与公司有belong关系
	sql += "(21, 9, 1001, 'belong', 'active', '2024-01-13 10:00:00', '2024-01-13 10:00:00', '系统管理员', '系统管理员'),"
	sql += "(22, 9, 10003, 'belong', 'active', '2024-01-13 10:00:00', '2024-01-13 10:00:00', '系统管理员', '系统管理员'),"
	sql += "(23, 9, 100004, 'belong', 'active', '2024-01-13 10:00:00', '2024-01-13 10:00:00', '系统管理员', '系统管理员'),"

	// employee (郑十, ID:10) - 与组和团队有belong关系，与公司有belong关系
	sql += "(24, 10, 1001, 'belong', 'active', '2024-01-13 10:00:00', '2024-01-13 10:00:00', '系统管理员', '系统管理员'),"
	sql += "(25, 10, 10003, 'belong', 'active', '2024-01-13 10:00:00', '2024-01-13 10:00:00', '系统管理员', '系统管理员'),"
	sql += "(26, 10, 100003, 'belong', 'active', '2024-01-13 10:00:00', '2024-01-13 10:00:00', '系统管理员', '系统管理员'),"

	// employee (王十一, ID:11) - 与组和团队有belong关系，与公司有belong关系
	sql += "(27, 11, 1001, 'belong', 'active', '2024-01-13 10:00:00', '2024-01-13 10:00:00', '系统管理员', '系统管理员'),"
	sql += "(28, 11, 10002, 'belong', 'active', '2024-01-13 10:00:00', '2024-01-13 10:00:00', '系统管理员', '系统管理员'),"
	sql += "(29, 11, 100002, 'belong', 'active', '2024-01-13 10:00:00', '2024-01-13 10:00:00', '系统管理员', '系统管理员');"

	return
}
