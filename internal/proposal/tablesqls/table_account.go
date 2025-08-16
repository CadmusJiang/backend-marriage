package tablesqls

//CREATE TABLE `account` (
//`id` int unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
//`username` varchar(32) NOT NULL COMMENT '用户名',
//`name` varchar(60) NOT NULL COMMENT '姓名',
//`password` varchar(255) NOT NULL COMMENT '密码',
//`phone` varchar(20) NOT NULL COMMENT '手机号',
//`role_type` varchar(20) NOT NULL COMMENT '角色类型',
//`status` enum('enabled','disabled') NOT NULL DEFAULT 'enabled' COMMENT '状态 enabled:enabled disabled:disabled',
//`belong_group_id` int unsigned NOT NULL DEFAULT 0 COMMENT '所属组ID',
//`belong_team_id` int unsigned NOT NULL DEFAULT 0 COMMENT '所属团队ID',
//`last_login_at` timestamp NULL DEFAULT NULL COMMENT '最后登录时间',
//`created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
//`updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
//`created_user` varchar(60) NOT NULL COMMENT '创建人',
//`updated_user` varchar(60) NOT NULL COMMENT '更新人',
//PRIMARY KEY (`id`),
//UNIQUE KEY `uk_username` (`username`),
//UNIQUE KEY `uk_phone` (`phone`),
//KEY `idx_belong_group_id` (`belong_group_id`),
//KEY `idx_belong_team_id` (`belong_team_id`),
//KEY `idx_role_type` (`role_type`)
//) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='账户表';

func CreateAccountTableSql() (sql string) {
	sql = "CREATE TABLE `account` ("
	sql += "`id` int unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',"
	sql += "`username` varchar(32) NOT NULL COMMENT '用户名',"
	sql += "`name` varchar(60) NOT NULL COMMENT '姓名',"
	sql += "`password` varchar(255) NOT NULL COMMENT '密码',"
	sql += "`phone` varchar(20) NOT NULL COMMENT '手机号',"
	sql += "`role_type` varchar(20) NOT NULL COMMENT '角色类型',"
	sql += "`status` enum('enabled','disabled') NOT NULL DEFAULT 'enabled' COMMENT '状态 enabled:enabled disabled:disabled',"
	sql += "`belong_group_id` int unsigned NOT NULL DEFAULT 0 COMMENT '所属组ID',"
	sql += "`belong_team_id` int unsigned NOT NULL DEFAULT 0 COMMENT '所属团队ID',"
	sql += "`last_login_at` timestamp NULL DEFAULT NULL COMMENT '最后登录时间',"
	sql += "`version` int NOT NULL DEFAULT 0 COMMENT '乐观锁版本号',"
	sql += "`created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',"
	sql += "`updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',"
	sql += "`created_user` varchar(60) NOT NULL COMMENT '创建人',"
	sql += "`updated_user` varchar(60) NOT NULL COMMENT '更新人',"
	sql += "PRIMARY KEY (`id`),"
	sql += "UNIQUE KEY `uk_username` (`username`),"
	sql += "UNIQUE KEY `uk_phone` (`phone`),"
	sql += "KEY `idx_belong_group_id` (`belong_group_id`),"
	sql += "KEY `idx_belong_team_id` (`belong_team_id`),"
	sql += "KEY `idx_role_type` (`role_type`)"
	sql += ") ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='账户表';"

	return
}

func CreateAccountTableDataSql() (sql string) {
	// 预置账户数据，确保与 org 和 account_org_relation 的示例数据相匹配
	// 密码: 123456 (bcrypt哈希)
	sql = "INSERT INTO `account` (`username`, `name`, `password`, `phone`, `role_type`, `status`, `belong_group_id`, `belong_team_id`, `last_login_at`, `created_user`, `updated_user`) VALUES"
	sql += "('admin', '系统管理员', '$2a$10$YVK93ajcju7xXxBkmkBOwOiD24IsZrWSxXFKbmvRToRP/T9g2TPMa', '13800138001', 'company_manager', 'enabled', 1, 4, NULL, '系统管理员', '系统管理员'),"
	sql += "('employee001', '张三', '$2a$10$YVK93ajcju7xXxBkmkBOwOiD24IsZrWSxXFKbmvRToRP/T9g2TPMa', '13800138002', 'employee', 'enabled', 2, 5, NULL, '系统管理员', '系统管理员'),"
	sql += "('employee002', '李四', '$2a$10$YVK93ajcju7xXxBkmkBOwOiD24IsZrWSxXFKbmvRToRP/T9g2TPMa', '13800138003', 'employee', 'enabled', 3, 6, NULL, '系统管理员', '系统管理员'),"
	sql += "('employee003', '王五', '$2a$10$YVK93ajcju7xXxBkmkBOwOiD24IsZrWSxXFKbmvRToRP/T9g2TPMa', '13800138004', 'employee', 'enabled', 3, 7, NULL, '系统管理员', '系统管理员'),"
	sql += "('group_manager001', '陈静', '$2a$10$YVK93ajcju7xXxBkmkBOwOiD24IsZrWSxXFKbmvRToRP/T9g2TPMa', '13800138005', 'group_manager', 'enabled', 2, 0, NULL, '系统管理员', '系统管理员'),"
	sql += "('team_manager001', '赵六', '$2a$10$YVK93ajcju7xXxBkmkBOwOiD24IsZrWSxXFKbmvRToRP/T9g2TPMa', '13800138006', 'team_manager', 'enabled', 3, 6, NULL, '系统管理员', '系统管理员'),"
	sql += "('employee004', '孙七', '$2a$10$YVK93ajcju7xXxBkmkBOwOiD24IsZrWSxXFKbmvRToRP/T9g2TPMa', '13800138007', 'employee', 'enabled', 2, 5, NULL, '系统管理员', '系统管理员'),"
	sql += "('employee005', '周八', '$2a$10$YVK93ajcju7xXxBkmkBOwOiD24IsZrWSxXFKbmvRToRP/T9g2TPMa', '13800138008', 'employee', 'enabled', 3, 6, NULL, '系统管理员', '系统管理员'),"
	sql += "('employee006', '吴九', '$2a$10$YVK93ajcju7xXxBkmkBOwOiD24IsZrWSxXFKbmvRToRP/T9g2TPMa', '13800138009', 'employee', 'enabled', 2, 5, NULL, '系统管理员', '系统管理员'),"
	sql += "('employee007', '郑十', '$2a$10$YVK93ajcju7xXxBkmkBOwOiD24IsZrWSxXFKbmvRToRP/T9g2TPMa', '13800138010', 'employee', 'enabled', 3, 6, NULL, '系统管理员', '系统管理员'),"
	sql += "('employee008', '王十一', '$2a$10$YVK93ajcju7xXxBkmkBOwOiD24IsZrWSxXFKbmvRToRP/T9g2TPMa', '13800138011', 'employee', 'enabled', 2, 5, NULL, '系统管理员', '系统管理员');"

	return
}
