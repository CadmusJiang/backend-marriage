# Marriage System API 文档

## 概述

Marriage System 提供了完整的组织管理和账户管理功能，包括以下核心模块：
- **org**: 组织信息管理
- **account**: 账户管理
- **account_history**: 账户历史记录
- **account_org_relation**: 账户组织关系

## 访问信息

- **服务地址**: `http://localhost:9999`
- **Swagger文档**: `http://localhost:9999/swagger/index.html`
- **API基础路径**: `/api/v1`

## 认证

所有API（除了登录接口）都需要在请求头中包含认证token：
```
Authorization: Bearer <token>
```

## 1. 认证相关API

### 1.1 用户登录
```
POST http://localhost:9999/api/v1/auth/login
```

**请求体：**
```json
{
  "username": "admin",
  "password": "123456"
}
```

**响应：**
```json
{
  "success": true,
  "message": "登录成功",
  "token": "token_20231201120000_admin",
  "userInfo": {
    "id": "acc_001",
    "username": "admin",
    "nickname": "系统管理员",
    "roleType": "company_manager",
    "status": "enabled",
    "phone": "13800138000",
    "belongGroup": "系统管理组",
    "belongTeam": "系统管理团队",
    "createdTime": 1705123200,
    "lastLoginTime": 1705123200,
    "lastModifiedTime": 1705123200
  }
}
```

### 1.2 退出登录
```
POST http://localhost:9999/api/v1/auth/logout
```

**响应：**
```json
{
  "success": true,
  "message": "退出登录成功"
}
```

## 2. 账户管理API

### 2.1 获取账户列表（批量获取）
```
GET http://localhost:9999/api/v1/accounts?current=1&pageSize=10
```

**查询参数：**
- `current`: 当前页码 (默认: 1)
- `pageSize`: 每页数量 (默认: 10)
- `username`: 用户名搜索
- `nickname`: 姓名搜索
- `roleType`: 角色类型筛选
- `status`: 状态筛选
- `phone`: 手机号搜索
- `belongGroup`: 所属组筛选
- `belongTeam`: 所属团队筛选

**响应：**
```json
{
  "data": [
    {
      "id": "acc_001",
      "username": "admin",
      "nickname": "系统管理员",
      "phone": "13800138000",
      "roleType": "company_manager",
      "belongGroup": {
        "id": 1,
        "username": "admin_group",
        "nickname": "系统管理组",
        "createdTimestamp": 1705123200,
        "modifiedTimestamp": 1705123200,
        "currentCnt": 1
      },
      "belongTeam": {
        "id": 1,
        "username": "admin_team",
        "nickname": "系统管理团队",
        "createdTimestamp": 1705123200,
        "modifiedTimestamp": 1705123200,
        "currentCnt": 1
      },
      "status": "enabled",
      "createdTimestamp": 1705123200,
      "modifiedTimestamp": 1705123200,
      "lastLoginTimestamp": 1705123200
    }
  ],
  "total": 5,
  "success": true,
  "pageSize": 10,
  "current": 1
}
```

### 2.2 获取单个账户详情
```
GET http://localhost:9999/api/v1/accounts/acc_001
```

**响应：**
```json
{
  "data": {
    "id": "acc_001",
    "username": "admin",
    "nickname": "系统管理员",
    "phone": "13800138000",
    "roleType": "company_manager",
    "belongGroup": {
      "id": 1,
      "username": "admin_group",
      "nickname": "系统管理组",
      "createdTimestamp": 1705123200,
      "modifiedTimestamp": 1705123200,
      "currentCnt": 1
    },
    "belongTeam": {
      "id": 1,
      "username": "admin_team",
      "nickname": "系统管理团队",
      "createdTimestamp": 1705123200,
      "modifiedTimestamp": 1705123200,
      "currentCnt": 1
    },
    "status": "enabled",
    "createdTimestamp": 1705123200,
    "modifiedTimestamp": 1705123200,
    "lastLoginTimestamp": 1705123200
  },
  "success": true
}
```

### 2.3 创建账户
```
POST http://localhost:9999/api/v1/accounts
```

**请求体：**
```json
{
  "username": "newuser",
  "nickname": "新用户",
  "password": "123456",
  "phone": "13800138006",
  "belongGroup": "技术组",
  "roleType": "employee"
}
```

**响应：**
```json
{
  "data": {
    "id": "acc_006",
    "username": "newuser",
    "nickname": "新用户",
    "roleType": "employee",
    "status": "enabled",
    "phone": "13800138006",
    "belongGroup": "技术组",
    "belongTeam": "",
    "createdTime": 1705123200,
    "lastLoginTime": 0,
    "lastModifiedTime": 1705123200
  },
  "success": true,
  "message": "账户创建成功"
}
```

### 2.4 更新账户
```
PUT http://localhost:9999/api/v1/accounts/acc_001
```

**请求体：**
```json
{
  "nickname": "新姓名",
  "phone": "13800138007",
  "belongGroup": "运营组",
  "belongTeam": "市场团队",
  "status": "enabled"
}
```

**响应：**
```json
{
  "success": true,
  "message": "账户更新成功"
}
```

### 2.5 获取账户历史记录
```
GET http://localhost:9999/api/v1/account-histories?current=1&pageSize=10
```

**响应：**
```json
{
  "data": [
    {
      "id": "hist_001",
      "accountId": "acc_001",
      "operateType": "create",
      "operateTimestamp": 1705123200,
      "content": {
        "username": "admin",
        "nickname": "系统管理员",
        "roleType": "company_manager"
      },
      "operator": "system",
      "operatorRoleType": "system"
    }
  ],
  "total": 5,
  "success": true,
  "pageSize": 10,
  "current": 1
}
```

## 3. 组织信息API

### 3.1 获取groups列表（专门获取groups信息）
```
GET http://localhost:9999/api/v1/groups
```

**查询参数：**
- `current`: 当前页码 (默认: 1)
- `pageSize`: 每页数量 (默认: 10)
- `orgName`: 组织名称搜索
- `orgLevel`: 组织层级筛选
- `parentOrgId`: 父组织ID筛选

**响应：**
```json
{
  "data": [
    {
      "id": 1,
      "orgName": "系统管理组",
      "orgType": 1,
      "orgPath": "/system/admin",
      "orgLevel": 1,
      "currentCnt": 1,
      "maxCnt": 10,
      "status": 1,
      "createdTimestamp": 1705123200,
      "modifiedTimestamp": 1705123200
    },
    {
      "id": 2,
      "orgName": "南京-天元大厦组",
      "orgType": 1,
      "orgPath": "/nanjing/tianyuan",
      "orgLevel": 1,
      "currentCnt": 15,
      "maxCnt": 50,
      "status": 1,
      "createdTimestamp": 1705123200,
      "modifiedTimestamp": 1705123200
    }
  ],
  "total": 5,
  "success": true,
  "pageSize": 10,
  "current": 1
}
```

### 3.2 获取组织列表（批量获取）
```
GET http://localhost:9999/api/v1/org-infos?current=1&pageSize=10
```

**查询参数：**
- `current`: 当前页码 (默认: 1)
- `pageSize`: 每页数量 (默认: 10)
- `orgName`: 组织名称搜索
- `orgType`: 组织类型筛选 (group/team)
- `orgLevel`: 组织层级筛选
- `parentOrgId`: 父组织ID筛选

**响应：**
```json
{
  "data": [
    {
      "id": "org_001",
      "orgName": "系统管理组",
      "orgType": "group",
      "orgLevel": 1,
      "parentOrgId": "",
      "orgDescription": "系统管理相关组织",
      "currentCnt": 1,
      "maxCnt": 10,
      "createdTimestamp": 1705123200,
      "modifiedTimestamp": 1705123200
    },
    {
      "id": "org_002",
      "orgName": "南京-天元大厦组",
      "orgType": "group",
      "orgLevel": 1,
      "parentOrgId": "",
      "orgDescription": "南京天元大厦办公区域",
      "currentCnt": 15,
      "maxCnt": 50,
      "createdTimestamp": 1705123200,
      "modifiedTimestamp": 1705123200
    }
  ],
  "total": 10,
  "success": true,
  "pageSize": 10,
  "current": 1
}
```

### 3.3 获取单个组织详情
```
GET http://localhost:9999/api/v1/org-infos/org_001
```

**响应：**
```json
{
  "success": true,
  "data": {
    "id": "org_001",
    "orgName": "系统管理组",
    "orgType": "group",
    "orgLevel": 1,
    "parentOrgId": "",
    "orgDescription": "系统管理相关组织",
    "currentCnt": 1,
    "maxCnt": 10,
    "createdTimestamp": 1705123200,
    "modifiedTimestamp": 1705123200
  }
}
```

### 3.4 创建组织
```
POST http://localhost:9999/api/v1/org-infos
```

**请求体：**
```json
{
  "orgName": "新组织",
  "orgType": "team",
  "orgLevel": 2,
  "parentOrgId": "org_001",
  "orgDescription": "新创建的团队",
  "maxCnt": 20
}
```

**响应：**
```json
{
  "success": true,
  "message": "组织信息创建成功",
  "id": 11
}
```

### 3.5 更新组织
```
PUT http://localhost:9999/api/v1/org-infos/org_001
```

**请求体：**
```json
{
  "orgName": "更新后的组织名称",
  "orgDescription": "更新后的组织描述",
  "maxCnt": 15
}
```

**响应：**
```json
{
  "success": true,
  "message": "组织信息更新成功"
}
```

### 3.6 删除组织
```
DELETE http://localhost:9999/api/v1/org-infos/org_001
```

**响应：**
```json
{
  "success": true,
  "message": "组织信息删除成功"
}
```

### 3.7 获取子组织
```
GET http://localhost:9999/api/v1/org-infos/org_001/children
```

**响应：**
```json
{
  "success": true,
  "data": [
    {
      "id": "org_006",
      "orgName": "系统管理团队",
      "orgType": "team",
      "orgLevel": 2,
      "parentOrgId": "org_001",
      "orgDescription": "系统管理相关团队",
      "currentCnt": 1,
      "maxCnt": 5,
      "createdTimestamp": 1705123200,
      "modifiedTimestamp": 1705123200
    }
  ]
}
```

### 3.8 获取父组织
```
GET http://localhost:9999/api/v1/org-infos/org_006/parent
```

**响应：**
```json
{
  "success": true,
  "data": {
    "id": "org_001",
    "orgName": "系统管理组",
    "orgType": "group",
    "orgLevel": 1,
    "parentOrgId": "",
    "orgDescription": "系统管理相关组织",
    "currentCnt": 1,
    "maxCnt": 10,
    "createdTimestamp": 1705123200,
    "modifiedTimestamp": 1705123200
  }
}
```

## 4. 数据库检查API

### 4.1 检查数据库状态
```
GET http://localhost:9999/api/v1/check-db
```

**响应：**
```json
{
  "account": {
    "exists": true,
    "columns": ["id", "account_id", "username", "nickname", "password", "phone", "role_type", "status", "belong_group_id", "belong_team_id", "created_timestamp", "modified_timestamp", "last_login_timestamp", "created_user", "updated_user"],
    "count": 5
  },
  "account_history": {
    "exists": true,
    "columns": ["id", "history_id", "account_id", "operate_type", "operate_timestamp", "content", "operator", "operator_role_type"],
    "count": 5
  },
  "org": {
    "exists": true,
    "columns": ["id", "org_id", "org_name", "org_type", "org_level", "parent_org_id", "org_description", "current_cnt", "max_cnt", "created_timestamp", "modified_timestamp", "is_deleted", "created_at", "created_user", "updated_at", "updated_user"],
    "count": 10
  }
}
```

## 5. 日志查看API

### 5.1 获取日志列表
```
GET http://localhost:9999/api/v1/logs
```

**查询参数：**
- `limit`: 获取条数 (默认: 50)

**响应：**
```json
{
  "data": [
    {
      "level": "info",
      "time": "2023-12-01T12:00:00Z",
      "path": "/api/v1/accounts",
      "http_code": 200,
      "method": "GET",
      "msg": "trace-log",
      "trace_id": "trace_20231201120000_001",
      "content": "完整的日志内容",
      "cost_seconds": 0.123,
      "success": true
    }
  ],
  "total": 50,
  "success": true,
  "last_time": "2023-12-01T12:00:00Z"
}
```

### 5.2 获取实时日志
```
GET http://localhost:9999/api/v1/logs/realtime
```

**请求头：**
- `since`: 起始时间 (可选)

**响应：**
```json
{
  "data": [
    {
      "level": "info",
      "time": "2023-12-01T12:01:00Z",
      "path": "/api/v1/groups",
      "http_code": 200,
      "method": "GET",
      "msg": "trace-log",
      "trace_id": "trace_20231201120100_002",
      "content": "完整的日志内容",
      "cost_seconds": 0.045,
      "success": true
    }
  ],
  "total": 10,
  "success": true,
  "last_time": "2023-12-01T12:01:00Z"
}
```

### 5.3 日志查看页面
```
GET http://localhost:9999/logs
```

**功能特性：**
- 实时显示API请求日志
- 支持按方法、状态、路径过滤
- 自动刷新功能
- 统计信息展示
- 响应时间监控
- 成功/失败请求统计

## 6. 测试数据

### 6.1 测试账户
| 用户名 | 姓名 | 角色 | 密码 | 状态 |
|--------|------|------|------|------|
| admin | 系统管理员 | company_manager | 123456 | enabled |
| company_manager | 张伟 | company_manager | 123456 | enabled |
| group_manager | 李娜 | group_manager | 123456 | enabled |
| team_manager | 王强 | team_manager | 123456 | enabled |
| employee | 赵敏 | employee | 123456 | enabled |

### 6.2 测试组织
| 组织ID | 组织名称 | 类型 | 层级 | 父组织 |
|--------|----------|------|------|--------|
| org_001 | 系统管理组 | group | 1 | - |
| org_002 | 南京-天元大厦组 | group | 1 | - |
| org_003 | 北京办公室组 | group | 1 | - |
| org_004 | 上海分公司组 | group | 1 | - |
| org_005 | 广州办公室组 | group | 1 | - |
| org_006 | 系统管理团队 | team | 2 | org_001 |
| org_007 | 营销团队A | team | 2 | org_002 |
| org_008 | 销售团队B | team | 2 | org_003 |
| org_009 | 技术团队C | team | 2 | org_004 |
| org_010 | 客服团队D | team | 2 | org_005 |

## 7. 快速测试命令

### 7.1 使用curl测试登录
```bash
curl -X POST http://localhost:9999/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"123456"}'
```

### 7.2 使用curl测试获取账户列表
```bash
curl -X GET "http://localhost:9999/api/v1/accounts?current=1&pageSize=5" \
  -H "Authorization: Bearer <your_token>"
```

### 7.3 使用curl测试获取组织列表
```bash
curl -X GET "http://localhost:9999/api/v1/org-infos?current=1&pageSize=5" \
  -H "Authorization: Bearer <your_token>"
```

## 8. 错误码说明

| 错误码 | 说明 |
|--------|------|
| 200 | 成功 |
| 400 | 请求参数错误 |
| 401 | 未授权 |
| 403 | 禁止访问 |
| 404 | 资源不存在 |
| 500 | 服务器内部错误 |

## 9. 注意事项

1. **认证**: 除了登录接口，所有API都需要在请求头中包含有效的token
2. **分页**: 列表接口都支持分页，默认每页10条记录
3. **搜索**: 支持按名称、类型等字段进行搜索和筛选
4. **数据格式**: 所有时间戳都是Unix时间戳格式
5. **状态码**: 成功响应状态码为200，错误响应会包含详细的错误信息

## 10. Swagger文档

启动服务后，可以通过以下地址访问Swagger文档：
```
http://localhost:9999/swagger/index.html
```

Swagger文档提供了完整的API接口说明，包括：
- 请求参数说明
- 响应格式说明
- 在线测试功能
- 错误码说明 