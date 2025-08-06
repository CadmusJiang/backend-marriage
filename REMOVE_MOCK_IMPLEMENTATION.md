# 移除Mock代码 - 使用真实数据库实现

## 概述

本文档说明如何将账户管理系统的Mock代码移除，改为使用真实数据库实现。

## 已完成的工作

### 1. 数据库模型生成

- ✅ 创建了 `internal/repository/mysql/account/gen_model.go` - 账户模型定义
- ✅ 创建了 `internal/repository/mysql/account_history/gen_model.go` - 历史记录模型定义
- ✅ 使用gormgen生成了数据库操作代码：
  - `internal/repository/mysql/account/gen_account.go`
  - `internal/repository/mysql/account_history/gen_account_history.go`

### 2. 服务层实现

- ✅ 创建了 `internal/services/account/service.go` - 服务接口定义
- ✅ 创建了 `internal/services/account/service_create.go` - 创建账户服务
- ✅ 创建了 `internal/services/account/service_list.go` - 列表和详情服务
- ✅ 创建了 `internal/services/account/service_update.go` - 更新和删除服务
- ✅ 创建了 `internal/services/account/service_auth.go` - 认证服务
- ✅ 创建了 `internal/services/account/service_history.go` - 历史记录服务

### 3. API层更新

- ✅ 更新了 `internal/api/account/handler.go` - 添加服务层依赖
- ✅ 更新了 `internal/api/account/func_list.go` - 使用真实数据库查询
- ✅ 更新了 `internal/api/account/func_login.go` - 使用真实数据库认证

### 4. 数据库初始化

- ✅ 创建了 `scripts/init_account_database.sql` - 数据库表结构和测试数据
- ✅ 创建了 `scripts/test_account_api_real.sh` - 真实数据库测试脚本

## 数据库表结构

### 账户表 (account)
```sql
CREATE TABLE `go_gin_api`.`account` (
  `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
  `account_id` varchar(32) NOT NULL COMMENT '账户ID',
  `username` varchar(32) NOT NULL COMMENT '用户名',
  `nickname` varchar(60) NOT NULL COMMENT '姓名',
  `password` varchar(100) NOT NULL COMMENT '密码',
  `phone` varchar(20) DEFAULT NULL COMMENT '手机号',
  `role_type` varchar(20) NOT NULL DEFAULT 'employee' COMMENT '角色类型',
  `status` varchar(20) NOT NULL DEFAULT 'enabled' COMMENT '状态',
  
  -- 所属组信息
  `belong_group_id` int DEFAULT NULL COMMENT '所属组ID',
  `belong_group_username` varchar(60) DEFAULT NULL COMMENT '所属组用户名',
  `belong_group_nickname` varchar(60) DEFAULT NULL COMMENT '所属组名称',
  `belong_group_created_timestamp` bigint DEFAULT NULL COMMENT '所属组创建时间戳',
  `belong_group_modified_timestamp` bigint DEFAULT NULL COMMENT '所属组修改时间戳',
  `belong_group_current_cnt` int DEFAULT 0 COMMENT '所属组当前成员数量',
  
  -- 所属团队信息
  `belong_team_id` int DEFAULT NULL COMMENT '所属团队ID',
  `belong_team_username` varchar(60) DEFAULT NULL COMMENT '所属团队用户名',
  `belong_team_nickname` varchar(60) DEFAULT NULL COMMENT '所属团队名称',
  `belong_team_created_timestamp` bigint DEFAULT NULL COMMENT '所属团队创建时间戳',
  `belong_team_modified_timestamp` bigint DEFAULT NULL COMMENT '所属团队修改时间戳',
  `belong_team_current_cnt` int DEFAULT 0 COMMENT '所属团队当前成员数量',
  
  `created_timestamp` bigint NOT NULL COMMENT '创建时间戳',
  `modified_timestamp` bigint NOT NULL COMMENT '修改时间戳',
  `last_login_timestamp` bigint DEFAULT NULL COMMENT '最后登录时间戳',
  
  `is_deleted` tinyint(1) NOT NULL DEFAULT '-1' COMMENT '是否删除',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `created_user` varchar(60) NOT NULL COMMENT '创建人',
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `updated_user` varchar(60) NOT NULL COMMENT '更新人',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_account_id` (`account_id`),
  UNIQUE KEY `uk_username` (`username`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='账户表';
```

### 历史记录表 (account_history)
```sql
CREATE TABLE `go_gin_api`.`account_history` (
  `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
  `history_id` varchar(32) NOT NULL COMMENT '历史记录ID',
  `account_id` varchar(32) NOT NULL COMMENT '账户ID',
  `operate_type` varchar(20) NOT NULL COMMENT '操作类型',
  `operate_timestamp` bigint NOT NULL COMMENT '操作时间戳',
  `content` json DEFAULT NULL COMMENT '操作内容 (JSON格式)',
  `operator` varchar(60) NOT NULL COMMENT '操作人',
  `operator_role_type` varchar(20) NOT NULL COMMENT '操作人角色',
  
  `is_deleted` tinyint(1) NOT NULL DEFAULT '-1' COMMENT '是否删除',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `created_user` varchar(60) NOT NULL COMMENT '创建人',
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `updated_user` varchar(60) NOT NULL COMMENT '更新人',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_history_id` (`history_id`),
  KEY `idx_account_id` (`account_id`),
  KEY `idx_operate_timestamp` (`operate_timestamp`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='账户历史记录表';
```

## 部署步骤

### 1. 初始化数据库

```bash
# 连接到MySQL数据库
mysql -u root -p

# 执行初始化脚本
source scripts/init_account_database.sql
```

### 2. 启动应用

```bash
# 启动Go应用
go run main.go
```

### 3. 测试API

```bash
# 运行测试脚本
chmod +x scripts/test_account_api_real.sh
./scripts/test_account_api_real.sh
```

## 测试数据

系统预置了以下测试账户：

| 用户名 | 姓名 | 角色 | 密码 | 状态 |
|--------|------|------|------|------|
| admin | 系统管理员 | company_manager | 123456 | enabled |
| company_manager | 张伟 | company_manager | 123456 | enabled |
| group_manager | 李娜 | group_manager | 123456 | enabled |
| team_manager | 王强 | team_manager | 123456 | enabled |
| employee | 赵敏 | employee | 123456 | enabled |

## API接口

### 认证相关
- `POST /api/v1/auth/login` - 用户登录
- `POST /api/v1/auth/logout` - 退出登录

### 账户管理
- `GET /api/v1/accounts` - 获取账户列表（支持分页、搜索、筛选）
- `POST /api/v1/accounts` - 创建账户
- `GET /api/v1/accounts/{accountId}` - 获取账户详情
- `PUT /api/v1/accounts/{accountId}` - 更新账户

### 历史记录
- `GET /api/v1/account-histories` - 获取账户历史记录

## 主要改进

### 1. 数据持久化
- 所有数据现在存储在MySQL数据库中
- 支持数据的增删改查操作
- 实现了软删除机制

### 2. 历史记录追踪
- 自动记录账户的创建、修改、删除操作
- 使用JSON格式存储操作内容
- 支持按账户ID查询历史记录

### 3. 密码安全
- 使用MD5加密存储密码
- 支持密码验证和更新

### 4. 搜索和筛选
- 支持按用户名、姓名、角色类型、状态搜索
- 支持分页查询
- 支持按所属组筛选

### 5. 组织架构支持
- 支持账户所属组和团队信息
- 记录组织的创建时间、修改时间、成员数量

## 注意事项

1. **密码加密**: 所有密码都使用MD5加密存储，测试密码为 `123456`
2. **软删除**: 删除操作使用软删除，不会真正删除数据
3. **历史记录**: 所有操作都会自动记录到历史记录表中
4. **数据一致性**: 使用事务确保数据一致性
5. **错误处理**: 完善的错误处理和日志记录

## 后续优化建议

1. **JWT认证**: 实现真正的JWT token认证
2. **密码加密**: 使用更安全的密码加密算法（如bcrypt）
3. **缓存优化**: 添加Redis缓存提高查询性能
4. **权限控制**: 实现基于角色的权限控制
5. **数据验证**: 添加更严格的数据验证规则
6. **API文档**: 完善Swagger API文档
7. **单元测试**: 添加完整的单元测试
8. **性能监控**: 添加性能监控和日志分析 