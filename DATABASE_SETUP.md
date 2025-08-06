# 账户管理系统数据库设置

## 概述

本文档说明如何设置账户管理系统的数据库，包括表结构创建和数据初始化。

## 数据库表结构

### 1. 账户表 (account)

| 字段名 | 类型 | 说明 | 默认值 |
|--------|------|------|--------|
| id | int unsigned | 主键 | AUTO_INCREMENT |
| account_id | varchar(32) | 账户ID | - |
| username | varchar(32) | 用户名 | - |
| nickname | varchar(60) | 姓名 | - |
| password | varchar(100) | 密码(MD5) | - |
| phone | varchar(20) | 手机号 | NULL |
| role_type | varchar(20) | 角色类型 | user |
| status | varchar(20) | 状态 | active |
| belong_group | varchar(60) | 所属组 | NULL |
| belong_team | varchar(60) | 所属团队 | NULL |
| last_login_time | timestamp | 最后登录时间 | NULL |
| is_deleted | tinyint(1) | 是否删除 | -1 |
| created_at | timestamp | 创建时间 | CURRENT_TIMESTAMP |
| created_user | varchar(60) | 创建人 | - |
| updated_at | timestamp | 更新时间 | CURRENT_TIMESTAMP |
| updated_user | varchar(60) | 更新人 | - |

### 2. 账户历史记录表 (account_history)

| 字段名 | 类型 | 说明 | 默认值 |
|--------|------|------|--------|
| id | int unsigned | 主键 | AUTO_INCREMENT |
| history_id | varchar(32) | 历史记录ID | - |
| account_id | varchar(32) | 账户ID | - |
| operate_type | varchar(20) | 操作类型 | - |
| operate_timestamp | timestamp | 操作时间戳 | CURRENT_TIMESTAMP |
| content | json | 操作内容 | NULL |
| operator | varchar(60) | 操作人 | - |
| operator_role_type | varchar(20) | 操作人角色 | - |
| is_deleted | tinyint(1) | 是否删除 | -1 |
| created_at | timestamp | 创建时间 | CURRENT_TIMESTAMP |
| created_user | varchar(60) | 创建人 | - |
| updated_at | timestamp | 更新时间 | CURRENT_TIMESTAMP |
| updated_user | varchar(60) | 更新人 | - |

## 初始化步骤

### 1. 创建数据库

```sql
CREATE DATABASE IF NOT EXISTS `go_gin_api` 
DEFAULT CHARACTER SET utf8mb4 
COLLATE utf8mb4_unicode_ci;
```

### 2. 执行初始化脚本

```bash
# 方法1：使用MySQL命令行
mysql -u your_username -p go_gin_api < scripts/init_database.sql

# 方法2：直接在MySQL中执行
mysql -u your_username -p
USE go_gin_api;
SOURCE scripts/init_database.sql;
```

### 3. 验证数据

```sql
-- 查看账户表数据
SELECT account_id, username, nickname, role_type, status, belong_group, belong_team 
FROM go_gin_api.account 
WHERE is_deleted = -1;

-- 查看历史记录数据
SELECT history_id, account_id, operate_type, operate_timestamp, operator, operator_role_type 
FROM go_gin_api.account_history 
WHERE is_deleted = -1 
ORDER BY operate_timestamp DESC;
```

## Mock数据说明

### 账户数据

系统预置了10个测试账户：

1. **张三** (employee001) - 小队管理员 - 南京-天元大厦组/营销团队A
2. **李四** (employee002) - 团队负责人 - 南京-天元大厦组/营销团队A
3. **王五** (employee003) - 员工 - 南京-夫子庙组/营销团队C (inactive)
4. **赵六** (employee004) - 员工 - 南京-天元大厦组/营销团队D
5. **钱七** (employee005) - 小队管理员 - 南京-夫子庙组/营销团队B
6. **孙八** (employee006) - 员工 - 南京-南京南站组/营销团队A
7. **周九** (employee007) - 员工 - 南京-夫子庙组/营销团队B
8. **吴十** (employee008) - 员工 - 南京-天元大厦组/营销团队C
9. **郑十一** (employee009) - 员工 - 南京-南京南站组/营销团队D
10. **王十二** (employee010) - 员工 - 南京-夫子庙组/营销团队A

### 历史记录数据

系统预置了18条历史记录，主要包含：

- **账户创建记录** - 记录账户的初始创建
- **角色变更记录** - 记录角色从员工到小队管理员、团队负责人的变更
- **团队调整记录** - 记录在不同营销团队间的调动
- **状态变更记录** - 记录账户启用/禁用的状态变更
- **信息修改记录** - 记录姓名、手机号等信息的修改

### 操作人员

- **系统管理员** (super_admin) - 负责账户创建
- **张伟** (company_manager) - 公司管理员
- **李明** (team_manager) - 团队经理
- **王芳** (team_manager) - 团队经理
- **刘强** (team_manager) - 团队经理
- **陈静** (team_manager) - 团队经理
- **赵敏** (team_manager) - 团队经理

## 角色类型说明

- **员工** - 基础员工角色
- **小队管理员** - 管理小团队的负责人
- **团队负责人** - 管理整个团队的负责人
- **company_manager** - 公司级管理员
- **team_manager** - 团队级管理员
- **super_admin** - 超级管理员

## 状态说明

- **active** - 账户启用
- **inactive** - 账户禁用

## 操作类型说明

- **created** - 账户创建
- **modified** - 账户信息修改

## 注意事项

1. 所有密码都使用MD5加密存储，测试密码为：`123456`
2. 历史记录使用JSON格式存储操作内容，便于前端解析
3. 时间戳使用标准的MySQL timestamp格式
4. 所有记录都包含审计字段（创建人、更新人、创建时间、更新时间）
5. 使用软删除机制（is_deleted字段）

## 测试账号

可以使用以下账号进行测试：

- **用户名**: employee001, employee002, employee003
- **密码**: 123456
- **角色**: 小队管理员, 团队负责人, 员工 