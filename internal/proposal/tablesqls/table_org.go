package tablesqls

//CREATE TABLE `org` (
//`id` int unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
//`org_type` tinyint(1) NOT NULL COMMENT '组织类型 1:group 2:team',
//`parent_id` int unsigned NOT NULL DEFAULT 0 COMMENT '父级ID，根为0',
//`path` varchar(255) NOT NULL COMMENT '层级路径，例如 /1/5/ 表示根->1->5',
//`username` varchar(32) NOT NULL COMMENT '唯一标识',
//`nickname` varchar(60) NOT NULL COMMENT '显示名称',
//`status` tinyint(1) NOT NULL DEFAULT 1 COMMENT '状态 1:enabled 2:disabled',
//`created_at` bigint NOT NULL COMMENT '创建时间戳',
//`updated_at` bigint NOT NULL COMMENT '修改时间戳',
//`created_user` varchar(60) NOT NULL COMMENT '创建人',
//`updated_user` varchar(60) NOT NULL COMMENT '更新人',
//PRIMARY KEY (`id`),
//UNIQUE KEY `uk_username` (`username`),
//KEY `idx_parent_id` (`parent_id`),
//KEY `idx_org_type` (`org_type`)
//) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='组织表';

func CreateOrgTableSql() (sql string) {
	sql = "CREATE TABLE `org` ("
	sql += "`id` int unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',"
	sql += "`org_type` tinyint(1) NOT NULL COMMENT '组织类型 1:group 2:team',"
	sql += "`parent_id` int unsigned NOT NULL DEFAULT 0 COMMENT '父级ID，根为0',"
	sql += "`path` varchar(255) NOT NULL COMMENT '层级路径，例如 /1/5/ 表示根->1->5',"
	sql += "`username` varchar(32) NOT NULL COMMENT '唯一标识',"
	sql += "`nickname` varchar(60) NOT NULL COMMENT '显示名称',"
	sql += "`status` tinyint(1) NOT NULL DEFAULT 1 COMMENT '状态 1:enabled 2:disabled',"
	sql += "`version` int NOT NULL DEFAULT 0 COMMENT '乐观锁版本号',"
	sql += "`created_at` bigint NOT NULL COMMENT '创建时间戳',"
	sql += "`updated_at` bigint NOT NULL COMMENT '修改时间戳',"
	sql += "`created_user` varchar(60) NOT NULL COMMENT '创建人',"
	sql += "`updated_user` varchar(60) NOT NULL COMMENT '更新人',"
	sql += "PRIMARY KEY (`id`),"
	sql += "UNIQUE KEY `uk_username` (`username`),"
	sql += "KEY `idx_parent_id` (`parent_id`),"
	sql += "KEY `idx_org_type` (`org_type`)"
	sql += ") ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='组织表';"

	return
}

func CreateOrgTableDataSql() (sql string) {
	// 预置三个组与四个团队，确保与 account 和 account_org_relation 的示例数据相匹配
	sql = "INSERT INTO `org` (`id`, `org_type`, `parent_id`, `path`, `username`, `nickname`, `status`, `version`, `created_at`, `updated_at`, `created_user`, `updated_user`) VALUES"
	// groups
	sql += "(1, 1, 0, '/1/', 'admin_group', '系统管理组', 1, 0, 1705123200, 1705123200, '系统管理员', '系统管理员'),"
	sql += "(2, 1, 0, '/2/', 'nanjing_tianyuan', '南京-天元大厦组', 1, 0, 1705123200, 1705123200, '系统管理员', '系统管理员'),"
	sql += "(3, 1, 0, '/3/', 'beijing_office', '北京办公室组', 1, 0, 1705123200, 1705123200, '系统管理员', '系统管理员'),"
	// teams
	sql += "(4, 2, 1, '/1/4/', 'admin_team', '系统管理团队', 1, 0, 1705123200, 1705123200, '系统管理员', '系统管理员'),"
	sql += "(5, 2, 2, '/2/5/', 'marketing_team_a', '营销团队A', 1, 0, 1705123200, 1705123200, '系统管理员', '系统管理员'),"
	sql += "(6, 2, 3, '/3/6/', 'sales_team_b', '销售团队B', 1, 0, 1705123200, 1705123200, '系统管理员', '系统管理员'),"
	sql += "(7, 2, 3, '/3/7/', 'tech_team_c', '技术团队C', 1, 0, 1705123200, 1705123200, '系统管理员', '系统管理员');"

	return
}
