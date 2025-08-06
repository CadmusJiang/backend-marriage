# 账户管理系统 API

## 概述

账户管理系统提供了完整的账户管理功能，包括账户的增删改查、登录登出、历史记录等功能。

## API 列表

### 1. 用户登录
- **URL**: `POST /api/v1/auth/login`
- **描述**: 用户登录接口
- **请求体**:
  ```json
  {
    "username": "admin",
    "password": "123456"
  }
  ```
- **响应**:
  ```json
  {
    "success": true,
    "message": "登录成功",
    "token": "token_20231201120000_admin",
    "userInfo": {
      "id": "acc_001",
      "username": "admin",
      "nickname": "管理员",
      "roleType": "admin",
      "status": "active",
      "phone": "13800138001",
      "belongGroup": "技术组",
      "belongTeam": "开发团队",
      "createdTime": 1701234567,
      "lastLoginTime": 1701234567,
      "lastModifiedTime": 1701234567
    }
  }
  ```

### 2. 退出登录
- **URL**: `POST /api/v1/auth/logout`
- **描述**: 用户退出登录
- **响应**:
  ```json
  {
    "success": true,
    "message": "退出登录成功"
  }
  ```

### 3. 获取账户列表
- **URL**: `GET /api/v1/accounts`
- **描述**: 分页获取账户列表，支持搜索和筛选
- **查询参数**:
  - `current`: 当前页码 (默认: 1)
  - `pageSize`: 每页数量 (默认: 10)
  - `username`: 用户名搜索
  - `nickname`: 姓名搜索
  - `roleType`: 角色类型筛选
  - `status`: 状态筛选
  - `phone`: 手机号搜索
  - `belongGroup`: 所属组筛选
  - `belongTeam`: 所属团队筛选
- **响应**:
  ```json
  {
    "data": [
      {
        "id": "acc_001",
        "username": "admin",
        "nickname": "管理员",
        "roleType": "admin",
        "status": "active",
        "phone": "13800138001",
        "belongGroup": "技术组",
        "belongTeam": "开发团队",
        "createdTime": 1701234567,
        "lastLoginTime": 1701234567,
        "lastModifiedTime": 1701234567
      }
    ],
    "total": 5,
    "success": true,
    "pageSize": 10,
    "current": 1
  }
  ```

### 4. 创建账户
- **URL**: `POST /api/v1/accounts`
- **描述**: 创建新账户
- **请求体**:
  ```json
  {
    "username": "newuser",
    "nickname": "新用户",
    "password": "123456",
    "phone": "13800138006",
    "belongGroup": "技术组",
    "roleType": "user"
  }
  ```
- **响应**:
  ```json
  {
    "data": {
      "id": "acc_006",
      "username": "newuser",
      "nickname": "新用户",
      "roleType": "user",
      "status": "active",
      "phone": "13800138006",
      "belongGroup": "技术组",
      "belongTeam": "",
      "createdTime": 1701234567,
      "lastLoginTime": 0,
      "lastModifiedTime": 1701234567
    },
    "success": true,
    "message": "账户创建成功"
  }
  ```

### 5. 获取账户详情
- **URL**: `GET /api/v1/accounts/{accountId}`
- **描述**: 根据账户ID获取详细信息
- **路径参数**:
  - `accountId`: 账户ID
- **响应**:
  ```json
  {
    "data": {
      "id": "acc_001",
      "username": "admin",
      "nickname": "管理员",
      "roleType": "admin",
      "status": "active",
      "phone": "13800138001",
      "belongGroup": "技术组",
      "belongTeam": "开发团队",
      "createdTime": 1701234567,
      "lastLoginTime": 1701234567,
      "lastModifiedTime": 1701234567
    },
    "success": true
  }
  ```

### 6. 更新账户
- **URL**: `PUT /api/v1/accounts/{accountId}`
- **描述**: 更新账户信息
- **路径参数**:
  - `accountId`: 账户ID
- **请求体**:
  ```json
  {
    "nickname": "新姓名",
    "phone": "13800138007",
    "belongGroup": "运营组",
    "belongTeam": "市场团队",
    "status": "active"
  }
  ```
- **响应**:
  ```json
  {
    "success": true,
    "message": "账户更新成功"
  }
  ```

### 7. 获取账户历史记录
- **URL**: `GET /api/v1/account-histories`
- **描述**: 分页获取账户操作历史记录
- **查询参数**:
  - `accountId`: 账户ID
  - `accountUsername`: 账户用户名
  - `current`: 当前页码 (默认: 1)
  - `pageSize`: 每页数量 (默认: 10)
- **响应**:
  ```json
  {
    "data": [
      {
        "id": "hist_001",
        "accountId": "acc_001",
        "operateType": "created",
        "operateTimestamp": 1701234567,
        "content": {
          "username": {
            "old": "",
            "new": "admin"
          },
          "nickname": {
            "old": "",
            "new": "管理员"
          }
        },
        "operator": "system",
        "operatorRoleType": "admin"
      }
    ],
    "total": 4,
    "success": true,
    "pageSize": 10,
    "current": 1
  }
  ```

## 数据库表结构

### 账户表 (account)
- `id`: 主键
- `account_id`: 账户ID (唯一)
- `username`: 用户名 (唯一)
- `nickname`: 姓名
- `password`: 密码 (MD5加密)
- `phone`: 手机号
- `role_type`: 角色类型 (admin/user/manager)
- `status`: 状态 (active/inactive)
- `belong_group`: 所属组
- `belong_team`: 所属团队
- `last_login_time`: 最后登录时间
- `created_at`: 创建时间
- `updated_at`: 更新时间

### 账户历史记录表 (account_history)
- `id`: 主键
- `history_id`: 历史记录ID (唯一)
- `account_id`: 账户ID
- `operate_type`: 操作类型 (created/modified)
- `operate_timestamp`: 操作时间戳
- `content`: 操作内容 (JSON格式)
- `operator`: 操作人
- `operator_role_type`: 操作人角色

## Mock数据

系统已预置了以下Mock数据：

### 账户数据
- admin (管理员) - 技术组/开发团队
- user001 (张三) - 运营组/市场团队
- user002 (李四) - 销售组/销售团队 (inactive)
- manager001 (王五) - 管理组/管理团队
- user003 (赵六) - 技术组/测试团队

### 测试账号
- 用户名: admin, user001
- 密码: 123456

## 注意事项

1. 所有需要认证的API都需要在请求头中包含有效的token
2. 密码使用MD5加密存储
3. 历史记录会自动记录账户的创建和修改操作
4. 分页查询支持多种筛选条件
5. 所有时间戳使用Unix时间戳格式 