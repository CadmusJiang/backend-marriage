package tablesqls

//CREATE TABLE `cooperation_store` (
//`id` int unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
//`store_name` varchar(100) NOT NULL COMMENT '门店名称',
//`cooperation_city` varchar(50) NOT NULL COMMENT '合作城市',
//`cooperation_type` json DEFAULT NULL COMMENT '合作类型',
//`store_short_name` varchar(50) DEFAULT NULL COMMENT '门店简称',
//`company_name` varchar(100) DEFAULT NULL COMMENT '公司名称',
//`cooperation_method` json DEFAULT NULL COMMENT '合作方式',
//`cooperation_status` varchar(20) NOT NULL DEFAULT 'pending' COMMENT '合作状态: active, inactive, pending',
//`business_license` varchar(500) DEFAULT NULL COMMENT '营业执照',
//`store_photos` json DEFAULT NULL COMMENT '门店照片',
//`actual_business_address` varchar(200) DEFAULT NULL COMMENT '实际经营地址',
//`contract_photos` json DEFAULT NULL COMMENT '合同照片',
//`version` int NOT NULL DEFAULT 0 COMMENT '乐观锁版本号',
//`created_timestamp` bigint NOT NULL COMMENT '创建时间戳',
//`modified_timestamp` bigint NOT NULL COMMENT '修改时间戳',
//`created_user` varchar(60) NOT NULL COMMENT '创建人',
//`updated_user` varchar(60) NOT NULL COMMENT '更新人',
//PRIMARY KEY (`id`),
//KEY `idx_store_name` (`store_name`),
//KEY `idx_cooperation_city` (`cooperation_city`),
//KEY `idx_cooperation_status` (`cooperation_status`)
//) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='合作门店表';

func CreateCooperationStoreTableSql() (sql string) {
	sql = "CREATE TABLE `cooperation_store` ("
	sql += "`id` int unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',"
	sql += "`store_name` varchar(100) NOT NULL COMMENT '门店名称',"
	sql += "`cooperation_city` varchar(50) NOT NULL COMMENT '合作城市',"
	sql += "`cooperation_type` json DEFAULT NULL COMMENT '合作类型',"
	sql += "`store_short_name` varchar(50) DEFAULT NULL COMMENT '门店简称',"
	sql += "`company_name` varchar(100) DEFAULT NULL COMMENT '公司名称',"
	sql += "`cooperation_method` json DEFAULT NULL COMMENT '合作方式',"
	sql += "`cooperation_status` varchar(20) NOT NULL DEFAULT 'pending' COMMENT '合作状态: active, inactive, pending',"
	sql += "`business_license` varchar(500) DEFAULT NULL COMMENT '营业执照',"
	sql += "`store_photos` json DEFAULT NULL COMMENT '门店照片',"
	sql += "`actual_business_address` varchar(200) DEFAULT NULL COMMENT '实际经营地址',"
	sql += "`contract_photos` json DEFAULT NULL COMMENT '合同照片',"
	sql += "`version` int NOT NULL DEFAULT 0 COMMENT '乐观锁版本号',"
	sql += "`created_at` bigint NOT NULL COMMENT '创建时间戳',"
	sql += "`updated_at` bigint NOT NULL COMMENT '修改时间戳',"
	sql += "`created_user` varchar(60) NOT NULL COMMENT '创建人',"
	sql += "`updated_user` varchar(60) NOT NULL COMMENT '更新人',"
	// 归属组（用于范围过滤，行业惯例）
	sql += "`belong_group_id` int unsigned DEFAULT NULL COMMENT '归属组ID',"
	sql += "PRIMARY KEY (`id`),"
	sql += "KEY `idx_store_name` (`store_name`),"
	sql += "KEY `idx_cooperation_city` (`cooperation_city`),"
	sql += "KEY `idx_cooperation_status` (`cooperation_status`),"
	sql += "KEY `idx_belong_group_id` (`belong_group_id`)"
	sql += ") ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='合作门店表';"

	return
}

func CreateCooperationStoreTableDataSql() (sql string) {
	sql = "INSERT INTO `cooperation_store` (`store_name`, `cooperation_city`, `cooperation_type`, `store_short_name`, `company_name`, `cooperation_method`, `cooperation_status`, `business_license`, `store_photos`, `actual_business_address`, `contract_photos`, `created_at`, `updated_at`, `created_user`, `updated_user`, `belong_group_id`) VALUES"
	sql += "('上海婚恋门店1', '上海', '[\"自然流\", \"投流\"]', '上海店1', '某某婚恋服务有限公司', '[\"CPA(by条数)\", \"CPS(by结算)\"]', 'active', 'https://picsum.photos/200/200?random=1', '[\"https://picsum.photos/300/200?random=1_1\", \"https://picsum.photos/300/200?random=1_2\", \"https://picsum.photos/300/200?random=1_3\"]', '上海市浦东新区某某路123号1层', '[\"https://picsum.photos/400/300?random=1_contract1\", \"https://picsum.photos/400/300?random=1_contract2\"]', 1705123200, 1705123200, '系统管理员', '系统管理员', 2),"
	sql += "('北京婚恋门店2', '北京', '[\"自然流\"]', '北京店2', '幸福婚恋集团', '[\"CPA(by条数)\"]', 'active', 'https://picsum.photos/200/200?random=2', '[\"https://picsum.photos/300/200?random=2_1\", \"https://picsum.photos/300/200?random=2_2\"]', '北京市朝阳区某某路456号2层', '[\"https://picsum.photos/400/300?random=2_contract1\"]', 1705209600, 1705209600, '系统管理员', '系统管理员', 2),"
	sql += "('深圳婚恋门店3', '深圳', '[\"投流\"]', '深圳店3', '缘来婚恋连锁', '[\"CPS(by结算)\"]', 'inactive', '', '[]', '深圳市南山区某某路789号3层', '[]', 1705296000, 1705296000, '系统管理员', '系统管理员', 3),"
	sql += "('杭州婚恋门店4', '杭州', '[\"自然流\", \"投流\"]', '杭州店4', '真爱婚恋服务', '[\"CPA(by条数)\", \"CPS(by结算)\"]', 'pending', 'https://picsum.photos/200/200?random=4', '[\"https://picsum.photos/300/200?random=4_1\"]', '杭州市西湖区某某路321号4层', '[\"https://picsum.photos/400/300?random=4_contract1\", \"https://picsum.photos/400/300?random=4_contract2\"]', 1705382400, 1705382400, '系统管理员', '系统管理员', 3),"
	sql += "('广州婚恋门店5', '广州', '[\"自然流\"]', '广州店5', '美满婚恋机构', '[\"CPA(by条数)\"]', 'active', 'https://picsum.photos/200/200?random=5', '[\"https://picsum.photos/300/200?random=5_1\", \"https://picsum.photos/300/200?random=5_2\", \"https://picsum.photos/300/200?random=5_3\"]', '广州市天河区某某路654号5层', '[\"https://picsum.photos/400/300?random=5_contract1\"]', 1705468800, 1705468800, '系统管理员', '系统管理员', 2),"
	sql += "('南京婚恋门店6', '南京', '[\"投流\"]', '南京店6', '红娘婚恋连锁', '[\"CPS(by结算)\"]', 'active', '', '[\"https://picsum.photos/300/200?random=6_1\"]', '南京市鼓楼区某某路987号6层', '[]', 1705555200, 1705555200, '系统管理员', '系统管理员', 3),"
	sql += "('成都婚恋门店7', '成都', '[\"自然流\", \"投流\"]', '成都店7', '缘分婚恋服务', '[\"CPA(by条数)\", \"CPS(by结算)\"]', 'inactive', 'https://picsum.photos/200/200?random=7', '[\"https://picsum.photos/300/200?random=7_1\", \"https://picsum.photos/300/200?random=7_2\"]', '成都市锦江区某某路147号7层', '[\"https://picsum.photos/400/300?random=7_contract1\"]', 1705641600, 1705641600, '系统管理员', '系统管理员', 2),"
	sql += "('武汉婚恋门店8', '武汉', '[\"自然流\"]', '武汉店8', '佳缘婚恋机构', '[\"CPA(by条数)\"]', 'pending', '', '[]', '武汉市江汉区某某路258号8层', '[]', 1705728000, 1705728000, '系统管理员', '系统管理员', 2),"
	sql += "('西安婚恋门店9', '西安', '[\"投流\"]', '西安店9', '百合婚恋集团', '[\"CPS(by结算)\"]', 'active', 'https://picsum.photos/200/200?random=9', '[\"https://picsum.photos/300/200?random=9_1\", \"https://picsum.photos/300/200?random=9_2\"]', '西安市雁塔区某某路369号9层', '[\"https://picsum.photos/400/300?random=9_contract1\", \"https://picsum.photos/400/300?random=9_contract2\"]', 1705814400, 1705814400, '系统管理员', '系统管理员', 3),"
	sql += "('重庆婚恋门店10', '重庆', '[\"自然流\", \"投流\"]', '重庆店10', '珍爱婚恋服务', '[\"CPA(by条数)\", \"CPS(by结算)\"]', 'active', 'https://picsum.photos/200/200?random=10', '[\"https://picsum.photos/300/200?random=10_1\", \"https://picsum.photos/300/200?random=10_2\", \"https://picsum.photos/300/200?random=10_3\"]', '重庆市渝中区某某路741号10层', '[\"https://picsum.photos/400/300?random=10_contract1\"]', 1705900800, 1705900800, '系统管理员', '系统管理员', 3);"

	return
}
