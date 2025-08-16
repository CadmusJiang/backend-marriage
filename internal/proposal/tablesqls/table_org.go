package tablesqls

//CREATE TABLE `org` (
//`id` int unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
//`org_type` enum('group','team') NOT NULL COMMENT '组织类型 group:group team:team',
//`parent_id` int unsigned NOT NULL DEFAULT 0 COMMENT '父级ID，根为0',
//`path` varchar(255) NOT NULL COMMENT '层级路径，例如 /1/5/ 表示根->1->5',
//`username` varchar(32) NOT NULL COMMENT '唯一标识',
//`name` varchar(60) NOT NULL COMMENT '显示名称',
//`status` enum('enabled','disabled') NOT NULL DEFAULT 'enabled' COMMENT '状态 enabled:enabled disabled:disabled',
//`created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
//`updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
//`created_user` varchar(60) NOT NULL COMMENT '创建人',
//`updated_user` varchar(60) NOT NULL COMMENT '更新人',
//PRIMARY KEY (`id`),
//UNIQUE KEY `uk_username` (`username`),
//KEY `idx_parent_id` (`parent_id`),
//KEY `idx_org_type` (`org_type`)
//) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='组织表';

func CreateOrgTableSql() (sql string) {
	sql = "CREATE TABLE `org` ("
	sql += "`id` int unsigned NOT NULL COMMENT '主键（构造ID，不使用自增）',"
	sql += "`org_type` enum('company','group','team') NOT NULL COMMENT '组织类型 company:company group:group team:team',"
	sql += "`parent_id` int unsigned NOT NULL DEFAULT 0 COMMENT '父级ID，根为0',"
	sql += "`path` varchar(255) NOT NULL COMMENT '层级路径，例如 /1001/10001/100001/ 表示公司->组->团队',"
	sql += "`username` varchar(32) NOT NULL COMMENT '唯一标识',"
	sql += "`name` varchar(60) NOT NULL COMMENT '显示名称',"
	sql += "`address` varchar(255) DEFAULT NULL COMMENT '地址',"
	sql += "`extra` json DEFAULT NULL COMMENT '扩展信息',"
	sql += "`status` enum('enabled','disabled') NOT NULL DEFAULT 'enabled' COMMENT '状态 enabled:enabled disabled:disabled',"
	sql += "`version` int NOT NULL DEFAULT 0 COMMENT '乐观锁版本号',"
	sql += "`created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',"
	sql += "`updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',"
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
	// 三级组织架构：公司 -> 组 -> 团队
	// 使用构造ID设计：公司(1000-9999), 组(10000-99999), 团队(100000-999999)
	sql = "INSERT INTO `org` (`id`, `org_type`, `parent_id`, `path`, `username`, `name`, `address`, `extra`, `status`, `version`, `created_at`, `updated_at`, `created_user`, `updated_user`) VALUES"

	// 1级：公司 (ID: 1000-9999)
	sql += "(1001, 'company', 0, '/1001/', 'company_main', '主公司', '上海市浦东新区', '{\"prefix\": \"company\", \"maxMemberCount\": 1000, \"maxTeamCount\": 100}', 'enabled', 0, '2024-01-13 10:00:00', '2024-01-13 10:00:00', '系统管理员', '系统管理员'),"

	// 2级：组（属于公司）(ID: 10000-99999)
	sql += "(10001, 'group', 1001, '/1001/10001/', 'admin_group', '系统管理组', '上海市浦东新区', '{\"prefix\": \"admin\", \"maxMemberCount\": 50, \"maxTeamCount\": 10}', 'enabled', 0, '2024-01-13 10:00:00', '2024-01-13 10:00:00', '系统管理员', '系统管理员'),"
	sql += "(10002, 'group', 1001, '/1001/10002/', 'nanjing_tianyuan', '南京-天元大厦组', '南京市建邺区天元大厦', '{\"prefix\": \"nanjing\", \"maxMemberCount\": 100, \"maxTeamCount\": 20}', 'enabled', 0, '2024-01-13 10:00:00', '2024-01-13 10:00:00', '系统管理员', '系统管理员'),"
	sql += "(10003, 'group', 1001, '/1001/10003/', 'beijing_office', '北京办公室组', '北京市朝阳区', '{\"prefix\": \"beijing\", \"maxMemberCount\": 80, \"maxTeamCount\": 15}', 'enabled', 0, '2024-01-13 10:00:00', '2024-01-13 10:00:00', '系统管理员', '系统管理员'),"

	// 3级：团队（属于组）(ID: 100000-999999)
	sql += "(100001, 'team', 10001, '/1001/10001/100001/', 'admin_team', '系统管理团队', '上海市浦东新区', '{\"prefix\": \"admin_team\", \"maxMemberCount\": 20, \"maxTeamCount\": 0}', 'enabled', 0, '2024-01-13 10:00:00', '2024-01-13 10:00:00', '系统管理员', '系统管理员'),"
	sql += "(100002, 'team', 10002, '/1001/10002/100002/', 'marketing_team_a', '营销团队A', '南京市建邺区天元大厦', '{\"prefix\": \"marketing\", \"maxMemberCount\": 30, \"maxTeamCount\": 0}', 'enabled', 0, '2024-01-13 10:00:00', '2024-01-13 10:00:00', '系统管理员', '系统管理员'),"
	sql += "(100003, 'team', 10003, '/1001/10003/100003/', 'sales_team_b', '销售团队B', '北京市朝阳区', '{\"prefix\": \"sales\", \"maxMemberCount\": 25, \"maxTeamCount\": 0}', 'enabled', 0, '2024-01-01 10:00:00', '2024-01-13 10:00:00', '系统管理员', '系统管理员'),"
	sql += "(100004, 'team', 10003, '/1001/10003/100004/', 'tech_team_c', '技术团队C', '北京市朝阳区', '{\"prefix\": \"tech\", \"maxMemberCount\": 35, \"maxTeamCount\": 0}', 'enabled', 0, '2024-01-13 10:00:00', '2024-01-13 10:00:00', '系统管理员', '系统管理员');"

	return
}
