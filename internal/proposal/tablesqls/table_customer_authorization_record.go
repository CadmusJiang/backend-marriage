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
//`authorization_status` enum('authorized','unauthorized') NOT NULL DEFAULT 'unauthorized' COMMENT '授权状态 authorized:已授权 unauthorized:未授权',
//`auth_photos` json DEFAULT NULL COMMENT '授权照片',
//`completion_status` enum('complete','incomplete') NOT NULL DEFAULT 'incomplete' COMMENT '完善状态 complete:已完善 incomplete:未完善',
//`assignment_status` enum('assigned','unassigned') NOT NULL DEFAULT 'unassigned' COMMENT '分配状态 assigned:已分配 unassigned:未分配',
//`payment_status` enum('paid','unpaid') NOT NULL DEFAULT 'unpaid' COMMENT '付费状态 paid:已付费 unpaid:未付费',
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
//KEY `idx_authorization_status` (`authorization_status`),
//KEY `idx_payment_status` (`payment_status`),
//KEY `idx_created_timestamp` (`created_timestamp`)
//) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='客户授权记录表';

func CreateCustomerAuthorizationRecordTableSql() (sql string) {
	sql = "CREATE TABLE `customer_authorization_record` ("
	sql += "`id` int unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',"
	sql += "`name` varchar(60) NOT NULL COMMENT '客户姓名',"
	sql += "`birth_year` int DEFAULT NULL COMMENT '出生年份',"
	sql += "`gender` varchar(10) DEFAULT NULL COMMENT '性别',"
	sql += "`height` int DEFAULT NULL COMMENT '身高(cm)',"
	sql += "`city_code` varchar(10) DEFAULT NULL COMMENT '城市编码',"
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
	sql += "`authorization_status` enum('authorized','unauthorized') NOT NULL DEFAULT 'unauthorized' COMMENT '授权状态 authorized:已授权 unauthorized:未授权',"
	sql += "`auth_photos` json DEFAULT NULL COMMENT '授权照片',"
	sql += "`completion_status` enum('complete','incomplete') NOT NULL DEFAULT 'incomplete' COMMENT '完善状态 complete:已完善 incomplete:未完善',"
	sql += "`assignment_status` enum('assigned','unassigned') NOT NULL DEFAULT 'unassigned' COMMENT '分配状态 assigned:已分配 unassigned:未分配',"
	sql += "`payment_status` enum('paid','unpaid') NOT NULL DEFAULT 'unpaid' COMMENT '付费状态 paid:已付费 unpaid:未付费',"
	sql += "`payment_amount` decimal(10,2) DEFAULT 0 COMMENT '支付金额',"
	sql += "`refund_amount` decimal(10,2) DEFAULT 0 COMMENT '退款金额',"
	sql += "`belong_group_id` int unsigned DEFAULT NULL COMMENT '归属组ID',"
	sql += "`belong_team_id` int unsigned DEFAULT NULL COMMENT '归属团队ID',"
	sql += "`belong_account_id` int unsigned DEFAULT NULL COMMENT '归属账户ID',"
	sql += "`version` int NOT NULL DEFAULT 0 COMMENT '乐观锁版本号',"
	sql += "`created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',"
	sql += "`updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',"
	sql += "`created_user` varchar(60) NOT NULL COMMENT '创建人',"
	sql += "`updated_user` varchar(60) NOT NULL COMMENT '更新人',"
	sql += "PRIMARY KEY (`id`),"
	sql += "KEY `idx_name` (`name`),"
	sql += "KEY `idx_city_code` (`city_code`),"
	sql += "KEY `idx_phone` (`phone`),"
	sql += "KEY `idx_belong_group_id` (`belong_group_id`),"
	sql += "KEY `idx_belong_team_id` (`belong_team_id`),"
	sql += "KEY `idx_belong_account_id` (`belong_account_id`),"
	sql += "KEY `idx_authorization_status` (`authorization_status`),"
	sql += "KEY `idx_payment_status` (`payment_status`),"
	sql += "KEY `idx_created_at` (`created_at`)"
	sql += ") ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='客户授权记录表';"

	return
}

func CreateCustomerAuthorizationRecordTableDataSql() (sql string) {
	// 预置客户授权记录数据，确保与 account 和 org 的示例数据相匹配
	// 对应：组2=南京-天元大厦组，组3=北京办公室组；队5=营销团队A(组2)，队6=销售团队B(组3)，队7=技术团队C(组3)
	// 账户：2=张伟, 5=刘强, 4=王芳, 6=赵敏, 3=李明
	sql = "INSERT INTO `customer_authorization_record` (`name`, `birth_year`, `gender`, `height`, `city_code`, `auth_store`, `education`, `profession`, `income`, `phone`, `wechat`, `drainage_account`, `drainage_id`, `drainage_channel`, `remark`, `authorization_status`, `auth_photos`, `completion_status`, `assignment_status`, `payment_status`, `payment_amount`, `refund_amount`, `belong_group_id`, `belong_team_id`, `belong_account_id`, `created_at`, `updated_at`, `created_user`, `updated_user`) VALUES"
	sql += "('用户1', 1985, 'male', 175, '110000', '朝阳门店***', '本科', '工程师', '50w', '13800001234', 'wx_12345****', 'drainage_001', 'D12345', '小红书', '备注信息1', 'authorized', '[\"https://picsum.photos/300/200?random=0\",\"https://picsum.photos/300/200?random=1\",\"https://picsum.photos/300/200?random=2\"]', 'complete', 'assigned', 'paid', 25000.00, 0.00, 2, 5, 2, '2024-01-13 10:00:00', '2024-01-13 10:00:00', '系统管理员', '系统管理员'),"
	sql += "('用户2', 1990, 'female', 165, '310000', '浦东门店***', '硕士', '设计师', '80w', '13900005678', 'wx_67890****', 'drainage_002', 'D67890', '小红书', '', 'authorized', '[\"https://picsum.photos/300/200?random=3\",\"https://picsum.photos/300/200?random=4\",\"https://picsum.photos/300/200?random=5\"]', 'complete', 'assigned', 'unpaid', 0.00, 0.00, 3, 6, 5, '2024-01-14 10:00:00', '2024-01-14 10:00:00', '系统管理员', '系统管理员'),"
	sql += "('用户3', 1988, 'male', 180, '440100', '天河门店***', '大专', '销售', '30w', '13700009012', 'wx_34567****', 'drainage_003', 'D34567', '小红书', '备注信息3', 'unauthorized', '[]', 'incomplete', 'unassigned', 'unpaid', 0.00, 0.00, 2, NULL, 4, '2024-01-15 10:00:00', '2024-01-15 10:00:00', '系统管理员', '系统管理员'),"
	sql += "('用户4', 1992, 'female', 160, '440300', '南山门店***', '本科', '教师', '40w', '13600003456', 'wx_78901****', 'drainage_004', 'D78901', '小红书', '', 'authorized', '[\"https://picsum.photos/300/200?random=6\",\"https://picsum.photos/300/200?random=7\",\"https://picsum.photos/300/200?random=8\"]', 'complete', 'unassigned', 'unpaid', 0.00, 0.00, 3, 7, 6, '2024-01-16 10:00:00', '2024-01-16 10:00:00', '系统管理员', '系统管理员'),"
	sql += "('用户5', 1987, 'male', 178, '330100', '西湖门店***', '博士', '医生', '120w', '13500007890', 'wx_23456****', 'drainage_005', 'D23456', '小红书', '备注信息5', 'authorized', '[\"https://picsum.photos/300/200?random=9\",\"https://picsum.photos/300/200?random=10\",\"https://picsum.photos/300/200?random=11\"]', 'complete', 'assigned', 'paid', 35000.00, 2000.00, 3, 6, 3, '2024-01-17 10:00:00', '2024-01-17 10:00:00', '系统管理员', '系统管理员');"

	return
}
