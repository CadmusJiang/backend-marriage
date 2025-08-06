-- 更新表结构脚本
-- 将ID字段改为bigint unsigned，其他数字字段改为bigint

USE marriage_system;

-- 1. 更新account表
ALTER TABLE `account` 
MODIFY COLUMN `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
MODIFY COLUMN `belong_group_id` bigint unsigned DEFAULT NULL COMMENT '所属组ID',
MODIFY COLUMN `belong_team_id` bigint unsigned DEFAULT NULL COMMENT '所属团队ID';

-- 2. 更新account_history表
ALTER TABLE `account_history` 
MODIFY COLUMN `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
MODIFY COLUMN `account_id` bigint unsigned NOT NULL COMMENT '账户ID';

-- 3. 创建account_org_relation表（如果不存在）
CREATE TABLE IF NOT EXISTS `account_org_relation` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
  `account_id` bigint unsigned NOT NULL COMMENT '账户ID',
  `org_id` bigint unsigned NOT NULL COMMENT '组织ID',
  `relation_type` varchar(20) NOT NULL DEFAULT 'member' COMMENT '关系类型: member, manager, owner',
  `status` tinyint(1) NOT NULL DEFAULT 1 COMMENT '状态: 1-有效, 0-无效',
  `created_timestamp` bigint NOT NULL COMMENT '创建时间戳',
  `modified_timestamp` bigint NOT NULL COMMENT '修改时间戳',
  `created_user` varchar(60) NOT NULL COMMENT '创建人',
  `updated_user` varchar(60) NOT NULL COMMENT '更新人',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_account_org` (`account_id`, `org_id`),
  KEY `idx_account_id` (`account_id`),
  KEY `idx_org_id` (`org_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='账户组织关系表';

-- 4. 更新org表结构
ALTER TABLE `org` 
MODIFY COLUMN `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
CHANGE COLUMN `org_name` `organization_name` varchar(60) NOT NULL COMMENT '组织名称',
CHANGE COLUMN `org_type` `organization_type` tinyint(1) NOT NULL COMMENT '组织类型: 1-公司, 2-部门, 3-团队',
CHANGE COLUMN `org_path` `organization_path` varchar(255) NOT NULL COMMENT '组织路径',
CHANGE COLUMN `org_level` `organization_level` tinyint(1) NOT NULL DEFAULT 1 COMMENT '组织层级',
CHANGE COLUMN `current_cnt` `current_count` int NOT NULL DEFAULT 0 COMMENT '当前成员数量',
CHANGE COLUMN `max_cnt` `max_count` int DEFAULT NULL COMMENT '最大成员数量',
CHANGE COLUMN `ext_data` `extension_data` json DEFAULT NULL COMMENT '扩展数据';

-- 5. 添加外键约束（可选）
-- ALTER TABLE `account` 
-- ADD CONSTRAINT `fk_account_group` FOREIGN KEY (`belong_group_id`) REFERENCES `org` (`id`) ON DELETE SET NULL,
-- ADD CONSTRAINT `fk_account_team` FOREIGN KEY (`belong_team_id`) REFERENCES `org` (`id`) ON DELETE SET NULL;

-- ALTER TABLE `account_history` 
-- ADD CONSTRAINT `fk_history_account` FOREIGN KEY (`account_id`) REFERENCES `account` (`id`) ON DELETE CASCADE;

-- ALTER TABLE `account_org_relation` 
-- ADD CONSTRAINT `fk_relation_account` FOREIGN KEY (`account_id`) REFERENCES `account` (`id`) ON DELETE CASCADE,
-- ADD CONSTRAINT `fk_relation_org` FOREIGN KEY (`org_id`) REFERENCES `org` (`id`) ON DELETE CASCADE; 