-- 创建 org 表
USE marriage_system;

-- 删除旧表（如果存在）
DROP TABLE IF EXISTS `org`;

-- 创建组织信息表
CREATE TABLE `org` (
  `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
  `org_name` varchar(60) NOT NULL COMMENT '组织名称',
  `org_type` tinyint(1) NOT NULL COMMENT '组织类型: 1-group, 2-team',
  `org_path` varchar(255) NOT NULL COMMENT '组织路径',
  `org_level` tinyint(1) NOT NULL DEFAULT '1' COMMENT '组织层级: 1-组, 2-团队',
  `current_cnt` int NOT NULL DEFAULT '0' COMMENT '当前成员数量',
  `max_cnt` int NOT NULL DEFAULT '0' COMMENT '最大成员数量',
  `status` tinyint(1) NOT NULL DEFAULT '1' COMMENT '状态: 1-启用, 0-禁用',
  `ext_data` json DEFAULT NULL COMMENT '扩展数据 (JSON格式)',
  `created_at` bigint NOT NULL COMMENT '创建时间戳',
  `updated_at` bigint NOT NULL COMMENT '修改时间戳',
  `created_user` varchar(60) NOT NULL COMMENT '创建人',
  `updated_user` varchar(60) NOT NULL COMMENT '更新人',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_org_name` (`org_name`),
  UNIQUE KEY `uk_org_path` (`org_path`),
  KEY `idx_org_type` (`org_type`),
  KEY `idx_org_level` (`org_level`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='组织信息表';

-- 插入一些测试数据
INSERT INTO `org` (
  `org_name`, 
  `org_type`, 
  `org_path`, 
  `org_level`, 
  `current_cnt`, 
  `max_cnt`, 
  `status`, 
  `ext_data`, 
  `created_at`, 
  `updated_at`, 
  `created_user`, 
  `updated_user`
) VALUES 
('系统管理组', 1, '/system_admin', 1, 1, 10, 1, '{}', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 'admin', 'admin'),
('南京-天元大厦组', 1, '/nanjing_tianyuan', 1, 15, 50, 1, '{}', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 'admin', 'admin'),
('南京-南京南站组', 1, '/nanjing_nanjing', 1, 12, 30, 1, '{}', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 'admin', 'admin'),
('营销团队A', 2, '/marketing_team_a', 2, 8, 20, 1, '{}', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 'admin', 'admin'),
('营销团队B', 2, '/marketing_team_b', 2, 10, 25, 1, '{}', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 'admin', 'admin'),
('营销团队C', 2, '/marketing_team_c', 2, 6, 15, 1, '{}', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), 'admin', 'admin');

-- 验证表创建成功
SELECT COUNT(*) as org_count FROM `org`;
SELECT 'org 表创建成功' as result; 