package tablesqls

//CREATE TABLE `org` (
//`id` int unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
//`org_name` varchar(60) NOT NULL COMMENT '组织名称',
//`org_type` tinyint(1) NOT NULL COMMENT '组织类型 1:group 2:team',
//`org_path` varchar(255) NOT NULL COMMENT '组织路径，如: /1/5/12',
//`org_level` tinyint(1) NOT NULL DEFAULT 1 COMMENT '组织层级',
//`current_cnt` int NOT NULL DEFAULT 0 COMMENT '当前成员数量',
//`max_cnt` int DEFAULT NULL COMMENT '最大成员数量',
//`status` tinyint(1) NOT NULL DEFAULT 1 COMMENT '状态 1:enabled 2:disabled',
//`ext_data` json DEFAULT NULL COMMENT '扩展数据(JSON)',
//`created_timestamp` bigint NOT NULL COMMENT '创建时间戳',
//`modified_timestamp` bigint NOT NULL COMMENT '修改时间戳',
//`created_user` varchar(60) NOT NULL COMMENT '创建人',
//`updated_user` varchar(60) NOT NULL COMMENT '更新人',
//PRIMARY KEY (`id`),
//UNIQUE KEY `uk_org_path` (`org_path`),
//KEY `idx_org_type` (`org_type`),
//KEY `idx_org_level` (`org_level`)
//) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='组织信息表';

func CreateOrgInfoTableSql() (sql string) {
	sql = "CREATE TABLE `org` ("
	sql += "`id` int unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',"
	sql += "`org_name` varchar(60) NOT NULL COMMENT '组织名称',"
	sql += "`org_type` tinyint(1) NOT NULL COMMENT '组织类型 1:group 2:team',"
	sql += "`org_path` varchar(255) NOT NULL COMMENT '组织路径，如: /1/5/12',"
	sql += "`org_level` tinyint(1) NOT NULL DEFAULT 1 COMMENT '组织层级',"
	sql += "`current_cnt` int NOT NULL DEFAULT 0 COMMENT '当前成员数量',"
	sql += "`max_cnt` int DEFAULT NULL COMMENT '最大成员数量',"
	sql += "`status` tinyint(1) NOT NULL DEFAULT 1 COMMENT '状态 1:enabled 2:disabled',"
	sql += "`ext_data` json DEFAULT NULL COMMENT '扩展数据(JSON)',"
	sql += "`created_timestamp` bigint NOT NULL COMMENT '创建时间戳',"
	sql += "`modified_timestamp` bigint NOT NULL COMMENT '修改时间戳',"
	sql += "`created_user` varchar(60) NOT NULL COMMENT '创建人',"
	sql += "`updated_user` varchar(60) NOT NULL COMMENT '更新人',"
	sql += "PRIMARY KEY (`id`),"
	sql += "UNIQUE KEY `uk_org_path` (`org_path`),"
	sql += "KEY `idx_org_type` (`org_type`),"
	sql += "KEY `idx_org_level` (`org_level`)"
	sql += ") ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='组织信息表';"

	return
}

func CreateOrgInfoTableDataSql() (sql string) {
	sql = "INSERT INTO `org` (`org_name`, `org_type`, `org_path`, `org_level`, `current_cnt`, `max_cnt`, `status`, `ext_data`, `created_timestamp`, `modified_timestamp`, `created_user`, `updated_user`) VALUES"
	// 组 (org_type = 1) - 顶级组
	sql += "('系统管理组', 1, '/1', 1, 1, 10, 1, '{\"prefix\": \"weilan_\"}', 1705123200, 1705123200, '系统管理员', '系统管理员'),"
	sql += "('南京-天元大厦组', 1, '/2', 1, 15, 50, 1, '{\"prefix\": \"weilan_\"}', 1705123200, 1705123200, '系统管理员', '系统管理员'),"
	sql += "('南京-南京南站组', 1, '/3', 1, 12, 50, 1, '{\"prefix\": \"weilan_\"}', 1705123200, 1705123200, '系统管理员', '系统管理员'),"
	// 团队 (org_type = 2) - 每个团队只属于一个组
	sql += "('系统管理团队', 2, '/1/1', 2, 1, 10, 1, '{\"team_leader\": \"系统管理员\", \"project\": \"系统维护\"}', 1705123200, 1705123200, '系统管理员', '系统管理员'),"
	sql += "('营销团队A', 2, '/2/1', 2, 8, 20, 1, '{\"team_leader\": \"张伟\", \"project\": \"产品推广\"}', 1705123200, 1705123200, '系统管理员', '系统管理员'),"
	sql += "('营销团队B', 2, '/3/1', 2, 10, 20, 1, '{\"team_leader\": \"刘强\", \"project\": \"市场拓展\"}', 1705123200, 1705123200, '系统管理员', '系统管理员'),"
	sql += "('营销团队C', 2, '/2/2', 2, 6, 20, 1, '{\"team_leader\": \"王芳\", \"project\": \"品牌建设\"}', 1705123200, 1705123200, '系统管理员', '系统管理员');"

	return
}
