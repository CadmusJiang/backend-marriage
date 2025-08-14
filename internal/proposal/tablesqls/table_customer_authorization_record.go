package tablesqls

//CREATE TABLE `customer_authorization_record` (
//`id` int unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
//`name` varchar(60) NOT NULL COMMENT '客户姓名',
//`birth_year` int DEFAULT NULL COMMENT '出生年份',
//`gender` varchar(10) DEFAULT NULL COMMENT '性别',
//`height` int DEFAULT NULL COMMENT '身高(cm)',
//`city` varchar(50) DEFAULT NULL COMMENT '城市',
//`auth_store` varchar(100) DEFAULT NULL COMMENT '授权门店',
//`education` varchar(50) DEFAULT NULL COMMENT '学历',
//`profession` varchar(50) DEFAULT NULL COMMENT '职业',
//`income` varchar(20) DEFAULT NULL COMMENT '收入',
//`phone` varchar(20) DEFAULT NULL COMMENT '手机号',
//`wechat` varchar(50) DEFAULT NULL COMMENT '微信号',
//`drainage_account` varchar(50) DEFAULT NULL COMMENT '引流账户',
//`drainage_id` varchar(50) DEFAULT NULL COMMENT '引流ID',
//`drainage_channel` varchar(50) DEFAULT NULL COMMENT '引流渠道',
//`remark` text DEFAULT NULL COMMENT '备注',
//`is_authorized` tinyint(1) NOT NULL DEFAULT 0 COMMENT '是否已授权',
//`auth_photos` json DEFAULT NULL COMMENT '授权照片',
//`is_profile_complete` tinyint(1) NOT NULL DEFAULT 0 COMMENT '是否已完善资料',
//`is_assigned` tinyint(1) NOT NULL DEFAULT 0 COMMENT '是否已分配',
//`is_paid` tinyint(1) NOT NULL DEFAULT 0 COMMENT '是否已买单',
//`payment_amount` decimal(10,2) DEFAULT 0 COMMENT '支付金额',
//`refund_amount` decimal(10,2) DEFAULT 0 COMMENT '退款金额',
//`group` varchar(100) DEFAULT NULL COMMENT '归属组',
//`team` varchar(100) DEFAULT NULL COMMENT '归属团队',
//`account` varchar(100) DEFAULT NULL COMMENT '归属账户',
//`version` int NOT NULL DEFAULT 0 COMMENT '乐观锁版本号',
//`created_timestamp` bigint NOT NULL COMMENT '创建时间戳',
//`modified_timestamp` bigint NOT NULL COMMENT '修改时间戳',
//`created_user` varchar(60) NOT NULL COMMENT '创建人',
//`updated_user` varchar(60) NOT NULL COMMENT '更新人',
//PRIMARY KEY (`id`),
//KEY `idx_name` (`name`),
//KEY `idx_city` (`city`),
//KEY `idx_phone` (`phone`),
//KEY `idx_group` (`group`),
//KEY `idx_team` (`team`),
//KEY `idx_account` (`account`),
//KEY `idx_is_authorized` (`is_authorized`),
//KEY `idx_is_paid` (`is_paid`),
//KEY `idx_created_timestamp` (`created_timestamp`)
//) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='客户授权记录表';

func CreateCustomerAuthorizationRecordTableSql() (sql string) {
	sql = "CREATE TABLE `customer_authorization_record` ("
	sql += "`id` int unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',"
	sql += "`name` varchar(60) NOT NULL COMMENT '客户姓名',"
	sql += "`birth_year` int DEFAULT NULL COMMENT '出生年份',"
	sql += "`gender` varchar(10) DEFAULT NULL COMMENT '性别',"
	sql += "`height` int DEFAULT NULL COMMENT '身高(cm)',"
	sql += "`city` varchar(50) DEFAULT NULL COMMENT '城市',"
	sql += "`auth_store` varchar(100) DEFAULT NULL COMMENT '授权门店',"
	sql += "`education` varchar(50) DEFAULT NULL COMMENT '学历',"
	sql += "`profession` varchar(50) DEFAULT NULL COMMENT '职业',"
	sql += "`income` varchar(20) DEFAULT NULL COMMENT '收入',"
	sql += "`phone` varchar(20) DEFAULT NULL COMMENT '手机号',"
	sql += "`wechat` varchar(50) DEFAULT NULL COMMENT '微信号',"
	sql += "`drainage_account` varchar(50) DEFAULT NULL COMMENT '引流账户',"
	sql += "`drainage_id` varchar(50) DEFAULT NULL COMMENT '引流ID',"
	sql += "`drainage_channel` varchar(50) DEFAULT NULL COMMENT '引流渠道',"
	sql += "`remark` text DEFAULT NULL COMMENT '备注',"
	sql += "`is_authorized` tinyint(1) NOT NULL DEFAULT 0 COMMENT '是否已授权',"
	sql += "`auth_photos` json DEFAULT NULL COMMENT '授权照片',"
	sql += "`is_profile_complete` tinyint(1) NOT NULL DEFAULT 0 COMMENT '是否已完善资料',"
	sql += "`is_assigned` tinyint(1) NOT NULL DEFAULT 0 COMMENT '是否已分配',"
	sql += "`is_paid` tinyint(1) NOT NULL DEFAULT 0 COMMENT '是否已买单',"
	sql += "`payment_amount` decimal(10,2) DEFAULT 0 COMMENT '支付金额',"
	sql += "`refund_amount` decimal(10,2) DEFAULT 0 COMMENT '退款金额',"
	// 归属外键（行业惯例使用 *_id）
	sql += "`belong_group_id` int unsigned DEFAULT NULL COMMENT '归属组ID',"
	sql += "`belong_team_id` int unsigned DEFAULT NULL COMMENT '归属团队ID',"
	sql += "`belong_account_id` int unsigned DEFAULT NULL COMMENT '归属账户ID',"
	sql += "`version` int NOT NULL DEFAULT 0 COMMENT '乐观锁版本号',"
	sql += "`created_at` bigint NOT NULL COMMENT '创建时间戳',"
	sql += "`updated_at` bigint NOT NULL COMMENT '修改时间戳',"
	sql += "`created_user` varchar(60) NOT NULL COMMENT '创建人',"
	sql += "`updated_user` varchar(60) NOT NULL COMMENT '更新人',"
	sql += "PRIMARY KEY (`id`),"
	sql += "KEY `idx_name` (`name`),"
	sql += "KEY `idx_city` (`city`),"
	sql += "KEY `idx_phone` (`phone`),"
	sql += "KEY `idx_belong_group_id` (`belong_group_id`),"
	sql += "KEY `idx_belong_team_id` (`belong_team_id`),"
	sql += "KEY `idx_belong_account_id` (`belong_account_id`),"
	sql += "KEY `idx_is_authorized` (`is_authorized`),"
	sql += "KEY `idx_is_paid` (`is_paid`),"
	sql += "KEY `idx_created_timestamp` (`created_timestamp`)"
	sql += ") ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='客户授权记录表';"

	return
}

func CreateCustomerAuthorizationRecordTableDataSql() (sql string) {
	sql = "INSERT INTO `customer_authorization_record` (`name`, `birth_year`, `gender`, `height`, `city`, `auth_store`, `education`, `profession`, `income`, `phone`, `wechat`, `drainage_account`, `drainage_id`, `drainage_channel`, `remark`, `is_authorized`, `auth_photos`, `is_profile_complete`, `is_assigned`, `is_paid`, `payment_amount`, `refund_amount`, `belong_group_id`, `belong_team_id`, `belong_account_id`, `created_at`, `updated_at`, `created_user`, `updated_user`) VALUES"
	// 对应：组2=南京-天元大厦组，组3=北京办公室组；队5=营销团队A(组2)，队6=销售团队B(组3)，队7=技术团队C(组3)
	// 账户：2=张伟, 5=刘强, 4=王芳, 6=赵敏, 3=李明
	sql += "('用户1', 1985, 'male', 175, '北京', '朝阳门店***', '本科', '工程师', '50w', '13800001234', 'wx_12345****', 'drainage_001', 'D12345', '小红书', '备注信息1', 1, '[\"https://picsum.photos/300/200?random=0\",\"https://picsum.photos/300/200?random=1\",\"https://picsum.photos/300/200?random=2\"]', 1, 1, 1, 25000.00, 0.00, 2, 5, 2, 1705123200, 1705123200, '系统管理员', '系统管理员'),"
	sql += "('用户2', 1990, 'female', 165, '上海', '浦东门店***', '硕士', '设计师', '80w', '13900005678', 'wx_67890****', 'drainage_002', 'D67890', '小红书', '', 1, '[\"https://picsum.photos/300/200?random=3\",\"https://picsum.photos/300/200?random=4\",\"https://picsum.photos/300/200?random=5\"]', 1, 1, 0, 0.00, 0.00, 3, 6, 5, 1705209600, 1705209600, '系统管理员', '系统管理员'),"
	sql += "('用户3', 1988, 'male', 180, '广州', '天河门店***', '大专', '销售', '30w', '13700009012', 'wx_34567****', 'drainage_003', 'D34567', '小红书', '备注信息3', 0, '[]', 0, 0, 0, 0.00, 0.00, 2, NULL, 4, 1705296000, 1705296000, '系统管理员', '系统管理员'),"
	sql += "('用户4', 1992, 'female', 160, '深圳', '南山门店***', '本科', '教师', '40w', '13600003456', 'wx_78901****', 'drainage_004', 'D78901', '小红书', '', 1, '[\"https://picsum.photos/300/200?random=6\",\"https://picsum.photos/300/200?random=7\",\"https://picsum.photos/300/200?random=8\"]', 1, 0, 0, 0.00, 0.00, 3, 7, 6, 1705382400, 1705382400, '系统管理员', '系统管理员'),"
	sql += "('用户5', 1987, 'male', 178, '杭州', '西湖门店***', '博士', '医生', '120w', '13500007890', 'wx_23456****', 'drainage_005', 'D23456', '小红书', '备注信息5', 1, '[\"https://picsum.photos/300/200?random=9\",\"https://picsum.photos/300/200?random=10\",\"httpsum.photos/300/200?random=11\"]', 1, 1, 1, 35000.00, 2000.00, 3, 6, 3, 1705468800, 1705468800, '系统管理员', '系统管理员');"

	return
}
