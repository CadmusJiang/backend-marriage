package tablesqls

// CreateOutboxTableSql returns DDL for the outbox table used for reliable messaging
func CreateOutboxTableSql() (sql string) {
	sql = "CREATE TABLE `outbox_events` ("
	sql += "`id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',"
	sql += "`topic` varchar(128) NOT NULL COMMENT '投递目标/主题',"
	sql += "`payload` json NOT NULL COMMENT '事件载荷(JSON)',"
	sql += "`status` enum('unpublished','published') NOT NULL DEFAULT 'unpublished' COMMENT '发布状态 unpublished:未发布 published:已发布',"
	sql += "`retry_count` int NOT NULL DEFAULT 0 COMMENT '重试次数',"
	sql += "`created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',"
	sql += "`published_at` timestamp NULL DEFAULT NULL COMMENT '发布时间',"
	sql += "PRIMARY KEY (`id`),"
	sql += "KEY `idx_status_created` (`status`,`created_at`)"
	sql += ") ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='Outbox 事件表';"
	return
}
