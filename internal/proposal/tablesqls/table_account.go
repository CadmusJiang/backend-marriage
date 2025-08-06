package tablesqls

//CREATE TABLE `account` (
//`id` int unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
//`username` varchar(32) NOT NULL COMMENT '用户名',
//`nickname` varchar(60) NOT NULL COMMENT '姓名',
//`password` varchar(100) NOT NULL COMMENT '密码',
//`phone` varchar(20) DEFAULT NULL COMMENT '手机号',
//`role_type` varchar(20) NOT NULL DEFAULT 'employee' COMMENT '角色类型',
//`status` varchar(20) NOT NULL DEFAULT 'enabled' COMMENT '状态',
//`belong_group_id` int unsigned DEFAULT NULL COMMENT '所属组ID',
//`belong_team_id` int unsigned DEFAULT NULL COMMENT '所属团队ID',
//`last_login_timestamp` bigint DEFAULT NULL COMMENT '最后登录时间戳',
//`created_timestamp` bigint NOT NULL COMMENT '创建时间戳',
//`modified_timestamp` bigint NOT NULL COMMENT '修改时间戳',
//`created_user` varchar(60) NOT NULL COMMENT '创建人',
//`updated_user` varchar(60) NOT NULL COMMENT '更新人',
//PRIMARY KEY (`id`),
//UNIQUE KEY `uk_username` (`username`),
//KEY `idx_belong_group_id` (`belong_group_id`),
//KEY `idx_belong_team_id` (`belong_team_id`)
//) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='账户表';

func CreateAccountTableSql() (sql string) {
	sql = "CREATE TABLE `account` ("
	sql += "`id` int unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',"
	sql += "`username` varchar(32) NOT NULL COMMENT '用户名',"
	sql += "`nickname` varchar(60) NOT NULL COMMENT '姓名',"
	sql += "`password` varchar(100) NOT NULL COMMENT '密码',"
	sql += "`phone` varchar(20) DEFAULT NULL COMMENT '手机号',"
	sql += "`role_type` varchar(20) NOT NULL DEFAULT 'employee' COMMENT '角色类型',"
	sql += "`status` varchar(20) NOT NULL DEFAULT 'enabled' COMMENT '状态',"
	sql += "`belong_group_id` int unsigned DEFAULT NULL COMMENT '所属组ID',"
	sql += "`belong_team_id` int unsigned DEFAULT NULL COMMENT '所属团队ID',"
	sql += "`last_login_timestamp` bigint DEFAULT NULL COMMENT '最后登录时间戳',"
	sql += "`created_timestamp` bigint NOT NULL COMMENT '创建时间戳',"
	sql += "`modified_timestamp` bigint NOT NULL COMMENT '修改时间戳',"
	sql += "`created_user` varchar(60) NOT NULL COMMENT '创建人',"
	sql += "`updated_user` varchar(60) NOT NULL COMMENT '更新人',"
	sql += "PRIMARY KEY (`id`),"
	sql += "UNIQUE KEY `uk_username` (`username`),"
	sql += "KEY `idx_belong_group_id` (`belong_group_id`),"
	sql += "KEY `idx_belong_team_id` (`belong_team_id`)"
	sql += ") ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='账户表';"

	return
}

func CreateAccountTableDataSql() (sql string) {
	sql = "INSERT INTO `account` (`username`, `nickname`, `password`, `phone`, `role_type`, `status`, `belong_group_id`, `belong_team_id`, `last_login_timestamp`, `created_timestamp`, `modified_timestamp`, `created_user`, `updated_user`) VALUES"
	sql += "('admin', '系统管理员', 'e10adc3949ba59abbe56e057f20f883e', '13800138000', 'company_manager', 'enabled', 1, 4, 1705123200, 1705123200, 1705123200, '系统管理员', '系统管理员'),"
	sql += "('company_manager', '张伟', 'e10adc3949ba59abbe56e057f20f883e', '13800138001', 'company_manager', 'enabled', 2, 5, 1705123200, 1705123200, 1705123200, '系统管理员', '张伟'),"
	sql += "('group_manager', '李明', 'e10adc3949ba59abbe56e057f20f883e', '13800138002', 'group_manager', 'enabled', 2, NULL, 1705123200, 1705123200, 1705123200, '系统管理员', '李明'),"
	sql += "('team_manager', '王芳', 'e10adc3949ba59abbe56e057f20f883e', '13800138003', 'team_manager', 'enabled', 2, 5, 1705123200, 1705123200, 1705123200, '系统管理员', '刘强'),"
	sql += "('group_manager2', '刘强', 'e10adc3949ba59abbe56e057f20f883e', '13800138004', 'group_manager', 'enabled', 3, 6, 1705123200, 1705123200, 1705123200, '系统管理员', '王芳'),"
	sql += "('team_manager2', '赵敏', 'e10adc3949ba59abbe56e057f20f883e', '13800138005', 'team_manager', 'disabled', 2, 7, 1705123200, 1705123200, 1705123200, '系统管理员', '系统管理员'),"
	sql += "('employee001', '陈静', 'e10adc3949ba59abbe56e057f20f883e', '13800138006', 'employee', 'enabled', 2, 5, 1705123200, 1705123200, 1705123200, '系统管理员', '系统管理员'),"
	sql += "('employee002', '孙丽', 'e10adc3949ba59abbe56e057f20f883e', '13800138007', 'employee', 'enabled', 3, 6, 1705123200, 1705123200, 1705123200, '系统管理员', '系统管理员'),"
	sql += "('employee003', '周杰', 'e10adc3949ba59abbe56e057f20f883e', '13800138008', 'employee', 'enabled', 2, 7, 1705123200, 1705123200, 1705123200, '系统管理员', '系统管理员'),"
	sql += "('employee004', '吴婷', 'e10adc3949ba59abbe56e057f20f883e', '13800138009', 'employee', 'disabled', 3, 6, 1705123200, 1705123200, 1705123200, '系统管理员', '系统管理员'),"
	sql += "('employee005', '郑华', 'e10adc3949ba59abbe56e057f20f883e', '13800138010', 'employee', 'enabled', 2, 5, 1705123200, 1705123200, 1705123200, '系统管理员', '系统管理员');"

	return
}
