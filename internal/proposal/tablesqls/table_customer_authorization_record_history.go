package tablesqls

//CREATE TABLE `customer_authorization_record_history` (
//`id` int unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
//`history_id` varchar(32) NOT NULL COMMENT '历史记录ID',
//`customer_authorization_record_id` varchar(32) NOT NULL COMMENT '客户授权记录ID',
//`operate_type` varchar(20) NOT NULL COMMENT '操作类型: created, modified, deleted',
//`occurred_at` bigint NOT NULL COMMENT '操作发生时间戳',
//`content` json DEFAULT NULL COMMENT '操作内容',
//
//-- 操作人信息
//`operator_username` varchar(60) NOT NULL COMMENT '操作人用户名',
//`operator_nickname` varchar(60) NOT NULL COMMENT '操作人姓名',
//`operator_role_type` varchar(20) NOT NULL COMMENT '操作人角色类型',
//
//`deleted_at` timestamp NULL DEFAULT NULL COMMENT '软删除时间',
//`created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
//`created_user` varchar(60) NOT NULL COMMENT '创建人',
//`updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
//`updated_user` varchar(60) NOT NULL COMMENT '更新人',
//PRIMARY KEY (`id`),
//UNIQUE KEY `uk_history_id` (`history_id`),
//KEY `idx_customer_authorization_record_id` (`customer_authorization_record_id`)
//) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='客户授权记录历史记录表';

func CreateCustomerAuthorizationRecordHistoryTableSql() (sql string) {
	sql = "CREATE TABLE `customer_authorization_record_history` ("
	sql += "`id` int unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',"
	sql += "`history_id` varchar(32) NOT NULL COMMENT '历史记录ID',"
	sql += "`customer_authorization_record_id` varchar(32) NOT NULL COMMENT '客户授权记录ID',"
	sql += "`operate_type` varchar(20) NOT NULL COMMENT '操作类型: created, modified, deleted',"
	sql += "`occurred_at` bigint NOT NULL COMMENT '操作发生时间戳',"
	sql += "`content` json DEFAULT NULL COMMENT '操作内容',"
	sql += "`operator_username` varchar(60) NOT NULL COMMENT '操作人用户名',"
	sql += "`operator_nickname` varchar(60) NOT NULL COMMENT '操作人姓名',"
	sql += "`operator_role_type` varchar(20) NOT NULL COMMENT '操作人角色类型',"
	sql += "`deleted_at` timestamp NULL DEFAULT NULL COMMENT '软删除时间',"
	sql += "`created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',"
	sql += "`created_user` varchar(60) NOT NULL COMMENT '创建人',"
	sql += "`updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',"
	sql += "`updated_user` varchar(60) NOT NULL COMMENT '更新人',"
	sql += "PRIMARY KEY (`id`),"
	sql += "UNIQUE KEY `uk_history_id` (`history_id`),"
	sql += "KEY `idx_customer_authorization_record_id` (`customer_authorization_record_id`),"
	sql += "KEY `idx_occurred_at` (`occurred_at`)"
	sql += ") ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='客户授权记录历史记录表';"

	return
}

func CreateCustomerAuthorizationRecordHistoryTableDataSql() (sql string) {
	sql = "INSERT INTO `customer_authorization_record_history` (`history_id`, `customer_authorization_record_id`, `operate_type`, `occurred_at`, `content`, `operator_username`, `operator_nickname`, `operator_role_type`, `created_user`, `updated_user`) VALUES"
	sql += "('hist_car_001', '1', 'created', 1705123200, '{\"name\": {\"old\": \"\", \"new\": \"用户1\"}, \"birthYear\": {\"old\": \"\", \"new\": 1985}, \"gender\": {\"old\": \"\", \"new\": \"male\"}, \"height\": {\"old\": \"\", \"new\": 175}, \"city\": {\"old\": \"\", \"new\": \"北京\"}, \"authStore\": {\"old\": \"\", \"new\": \"朝阳门店***\"}, \"education\": {\"old\": \"\", \"new\": \"本科\"}, \"profession\": {\"old\": \"\", \"new\": \"工程师\"}, \"income\": {\"old\": \"\", \"new\": \"50w\"}, \"phone\": {\"old\": \"\", \"new\": \"138****1234\"}, \"wechat\": {\"old\": \"\", \"new\": \"wx_12345****\"}, \"drainageAccount\": {\"old\": \"\", \"new\": \"drainage_001\"}, \"drainageId\": {\"old\": \"\", \"new\": \"D12345\"}, \"drainageChannel\": {\"old\": \"\", \"new\": \"小红书\"}, \"remark\": {\"old\": \"\", \"new\": \"备注信息1\"}, \"isAuthorized\": {\"old\": \"\", \"new\": true}, \"authPhotos\": {\"old\": \"\", \"new\": [\"https://picsum.photos/300/200?random=0\", \"https://picsum.photos/300/200?random=1\", \"https://picsum.photos/300/200?random=2\"]}, \"isProfileComplete\": {\"old\": \"\", \"new\": true}, \"isAssigned\": {\"old\": \"\", \"new\": true}, \"isPaid\": {\"old\": \"\", \"new\": true}, \"paymentAmount\": {\"old\": \"\", \"new\": 25000.00}, \"refundAmount\": {\"old\": \"\", \"new\": 0.00}, \"group\": {\"old\": \"\", \"new\": \"南京-天元大厦组\"}, \"team\": {\"old\": \"\", \"new\": \"营销团队A\"}, \"account\": {\"old\": \"\", \"new\": \"张伟\"}}', 'admin', '系统管理员', 'company_manager', '系统管理员', '系统管理员'),"
	sql += "('hist_car_002', '1', 'modified', 1705209600, '{\"isAuthorized\": {\"old\": false, \"new\": true}, \"authPhotos\": {\"old\": [], \"new\": [\"https://picsum.photos/300/200?random=0\", \"https://picsum.photos/300/200?random=1\", \"https://picsum.photos/300/200?random=2\"]}}', 'zhangwei', '张伟', 'company_manager', '系统管理员', '系统管理员'),"
	sql += "('hist_car_003', '1', 'modified', 1705296000, '{\"isAssigned\": {\"old\": false, \"new\": true}, \"team\": {\"old\": \"\", \"new\": \"营销团队A\"}}', 'liming', '李明', 'team_manager', '系统管理员', '系统管理员'),"
	sql += "('hist_car_004', '1', 'modified', 1705382400, '{\"isPaid\": {\"old\": false, \"new\": true}, \"paymentAmount\": {\"old\": 0.00, \"new\": 25000.00}}', 'wangfang', '王芳', 'team_manager', '系统管理员', '系统管理员'),"
	sql += "('hist_car_005', '2', 'created', 1705209600, '{\"name\": {\"old\": \"\", \"new\": \"用户2\"}, \"birthYear\": {\"old\": \"\", \"new\": 1990}, \"gender\": {\"old\": \"\", \"new\": \"female\"}, \"height\": {\"old\": \"\", \"new\": 165}, \"city\": {\"old\": \"\", \"new\": \"上海\"}, \"authStore\": {\"old\": \"\", \"new\": \"浦东门店***\"}, \"education\": {\"old\": \"\", \"new\": \"硕士\"}, \"profession\": {\"old\": \"\", \"new\": \"设计师\"}, \"income\": {\"old\": \"\", \"new\": \"80w\"}, \"phone\": {\"old\": \"\", \"new\": \"139****5678\"}, \"wechat\": {\"old\": \"\", \"new\": \"wx_67890****\"}, \"drainageAccount\": {\"old\": \"\", \"new\": \"drainage_002\"}, \"drainageId\": {\"old\": \"\", \"new\": \"D67890\"}, \"drainageChannel\": {\"old\": \"\", \"new\": \"小红书\"}, \"remark\": {\"old\": \"\", \"new\": \"\"}, \"isAuthorized\": {\"old\": \"\", \"new\": true}, \"authPhotos\": {\"old\": \"\", \"new\": [\"https://picsum.photos/300/200?random=3\", \"https://picsum.photos/300/200?random=4\", \"https://picsum.photos/300/200?random=5\"]}, \"isProfileComplete\": {\"old\": \"\", \"new\": true}, \"isAssigned\": {\"old\": \"\", \"new\": true}, \"isPaid\": {\"old\": \"\", \"new\": false}, \"paymentAmount\": {\"old\": \"\", \"new\": 0.00}, \"refundAmount\": {\"old\": \"\", \"new\": 0.00}, \"group\": {\"old\": \"\", \"new\": \"南京-南京南站组\"}, \"team\": {\"old\": \"\", \"new\": \"营销团队B\"}, \"account\": {\"old\": \"\", \"new\": \"刘强\"}}', 'admin', '系统管理员', 'company_manager', '系统管理员', '系统管理员'),"
	sql += "('hist_car_006', '2', 'modified', 1705478400, '{\"isAssigned\": {\"old\": false, \"new\": true}, \"team\": {\"old\": \"\", \"new\": \"营销团队B\"}}', 'liuqiang', '刘强', 'team_manager', '系统管理员', '系统管理员'),"
	sql += "('hist_car_007', '3', 'created', 1705296000, '{\"name\": {\"old\": \"\", \"new\": \"用户3\"}, \"birthYear\": {\"old\": \"\", \"new\": 1988}, \"gender\": {\"old\": \"\", \"new\": \"male\"}, \"height\": {\"old\": \"\", \"new\": 180}, \"city\": {\"old\": \"\", \"new\": \"广州\"}, \"authStore\": {\"old\": \"\", \"new\": \"天河门店***\"}, \"education\": {\"old\": \"\", \"new\": \"大专\"}, \"profession\": {\"old\": \"\", \"new\": \"销售\"}, \"income\": {\"old\": \"\", \"new\": \"30w\"}, \"phone\": {\"old\": \"\", \"new\": \"137****9012\"}, \"wechat\": {\"old\": \"\", \"new\": \"wx_34567****\"}, \"drainageAccount\": {\"old\": \"\", \"new\": \"drainage_003\"}, \"drainageId\": {\"old\": \"\", \"new\": \"D34567\"}, \"drainageChannel\": {\"old\": \"\", \"new\": \"小红书\"}, \"remark\": {\"old\": \"\", \"new\": \"备注信息3\"}, \"isAuthorized\": {\"old\": \"\", \"new\": false}, \"authPhotos\": {\"old\": \"\", \"new\": []}, \"isProfileComplete\": {\"old\": \"\", \"new\": false}, \"isAssigned\": {\"old\": \"\", \"new\": false}, \"isPaid\": {\"old\": \"\", \"new\": false}, \"paymentAmount\": {\"old\": \"\", \"new\": 0.00}, \"refundAmount\": {\"old\": \"\", \"new\": 0.00}, \"group\": {\"old\": \"\", \"new\": \"南京-天元大厦组\"}, \"team\": {\"old\": \"\", \"new\": null}, \"account\": {\"old\": \"\", \"new\": \"王芳\"}}', 'admin', '系统管理员', 'company_manager', '系统管理员', '系统管理员'),"
	sql += "('hist_car_008', '3', 'modified', 1705564800, '{\"isAuthorized\": {\"old\": false, \"new\": true}, \"authPhotos\": {\"old\": [], \"new\": [\"https://picsum.photos/300/200?random=6\", \"https://picsum.photos/300/200?random=7\"]}}', 'wangfang', '王芳', 'team_manager', '系统管理员', '系统管理员'),"
	sql += "('hist_car_009', '3', 'modified', 1705651200, '{\"isProfileComplete\": {\"old\": false, \"new\": true}, \"education\": {\"old\": \"大专\", \"new\": \"本科\"}, \"profession\": {\"old\": \"销售\", \"new\": \"销售经理\"}}', 'zhangwei', '张伟', 'company_manager', '系统管理员', '系统管理员'),"
	sql += "('hist_car_010', '3', 'modified', 1705737600, '{\"isAssigned\": {\"old\": false, \"new\": true}, \"team\": {\"old\": null, \"new\": \"营销团队C\"}}', 'liming', '李明', 'team_manager', '系统管理员', '系统管理员');"

	return
}
