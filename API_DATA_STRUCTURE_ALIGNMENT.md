# API数据结构与前端Mock代码对齐说明

## 概述

本文档说明后端API返回的数据结构如何与前端mock代码保持一致，特别是 `belongGroup` 和 `belongTeam` 字段作为内置对象的处理。

## 前端Mock数据结构

前端mock代码中的账户数据结构：

```typescript
{
  id: '0',
  username: 'admin',
  nickname: '系统管理员',
  phone: '13800138000',
  roleType: 'company_manager',
  belongGroup: {
    id: 0,
    username: 'admin_group',
    nickname: '系统管理组',
    createdTimestamp: 1705123200,
    modifiedTimestamp: 1705123200,
    currentCnt: 1
  },
  belongTeam: {
    id: 0,
    username: 'admin_team',
    nickname: '系统管理团队',
    createdTimestamp: 1705123200,
    modifiedTimestamp: 1705123200,
    currentCnt: 1
  },
  status: 'enabled',
  createdTimestamp: 1705123200,
  modifiedTimestamp: 1705123200,
  lastLoginTimestamp: 1705123200,
}
```

## 后端API数据结构

### 1. 账户列表API (`GET /api/v1/accounts`)

**请求参数：**
- `current`: 当前页码（默认1）
- `pageSize`: 每页数量（默认10）
- `includeGroup`: 是否包含组信息（默认true）
- `includeTeam`: 是否包含团队信息（默认true）
- 其他搜索参数：`username`, `nickname`, `roleType`, `status`, `phone`, `belongGroup`, `belongTeam`

**响应结构：**
```json
{
  "success": true,
  "data": [
    {
      "id": "0",
      "username": "admin",
      "nickname": "系统管理员",
      "phone": "13800138000",
      "roleType": "company_manager",
      "belongGroup": {
        "id": 0,
        "username": "admin_group",
        "nickname": "系统管理组",
        "createdTimestamp": 1705123200,
        "modifiedTimestamp": 1705123200,
        "currentCnt": 1
      },
      "belongTeam": {
        "id": 0,
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
  "total": 10,
  "pageSize": 10,
  "current": 1
}
```

### 2. 账户详情API (`GET /api/v1/accounts/{accountId}`)

**请求参数：**
- `accountId`: 账户ID（路径参数）
- `includeGroup`: 是否包含组信息（默认true）
- `includeTeam`: 是否包含团队信息（默认true）

**响应结构：**
```json
{
  "success": true,
  "data": {
    "id": "0",
    "username": "admin",
    "nickname": "系统管理员",
    "phone": "13800138000",
    "roleType": "company_manager",
    "belongGroup": {
      "id": 0,
      "username": "admin_group",
      "nickname": "系统管理组",
      "createdTimestamp": 1705123200,
      "modifiedTimestamp": 1705123200,
      "currentCnt": 1
    },
    "belongTeam": {
      "id": 0,
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
}
```

## 关键特性

### 1. 可选的组织信息

通过 `includeGroup` 和 `includeTeam` 参数控制是否返回组织信息：

```bash
# 包含组织信息（默认）
GET /api/v1/accounts?current=1&pageSize=10

# 不包含组织信息
GET /api/v1/accounts?current=1&pageSize=10&includeGroup=false&includeTeam=false
```

### 2. 内置对象结构

`belongGroup` 和 `belongTeam` 作为内置对象，包含完整的组织信息：

```json
{
  "id": 0,                    // 组织ID（数字）
  "username": "admin_group",   // 组织用户名
  "nickname": "系统管理组",     // 组织名称
  "createdTimestamp": 1705123200,   // 创建时间戳
  "modifiedTimestamp": 1705123200,  // 修改时间戳
  "currentCnt": 1             // 当前成员数量
}
```

### 3. 空值处理

当账户没有关联的组织时，`belongGroup` 或 `belongTeam` 字段为 `null`：

```json
{
  "id": "2",
  "username": "group_manager",
  "nickname": "李明",
  "belongGroup": {
    "id": 1,
    "username": "nanjing_tianyuan",
    "nickname": "南京-天元大厦组",
    "createdTimestamp": 1705123200,
    "modifiedTimestamp": 1705123200,
    "currentCnt": 15
  },
  "belongTeam": null  // 没有关联的团队
}
```

## 实现细节

### 1. Go结构体定义

```go
type accountData struct {
    ID                 string   `json:"id"`                 // 账户ID
    Username           string   `json:"username"`           // 用户名
    Nickname           string   `json:"nickname"`           // 姓名
    Phone              string   `json:"phone"`              // 手机号
    RoleType           string   `json:"roleType"`           // 角色类型
    BelongGroup        *orgInfo `json:"belongGroup"`        // 所属组（指针类型，可为null）
    BelongTeam         *orgInfo `json:"belongTeam"`         // 所属团队（指针类型，可为null）
    Status             string   `json:"status"`             // 状态
    CreatedTimestamp   int64    `json:"createdTimestamp"`   // 创建时间戳
    ModifiedTimestamp  int64    `json:"modifiedTimestamp"`  // 修改时间戳
    LastLoginAt        int64    `json:"lastLoginAt"`        // 最后登录时间戳
}

type orgInfo struct {
    ID                int    `json:"id"`                // 组织ID
    Username          string `json:"username"`          // 组织用户名
    Nickname          string `json:"nickname"`          // 组织名称
    CreatedTimestamp  int64  `json:"createdTimestamp"`  // 创建时间戳
    ModifiedTimestamp int64  `json:"modifiedTimestamp"` // 修改时间戳
    CurrentCnt        int    `json:"currentCnt"`        // 当前成员数量
}
```

### 2. 条件返回逻辑

```go
// 根据includeGroup参数决定是否包含组信息
if req.IncludeGroup != "false" && acc.BelongGroupId > 0 {
    accountData.BelongGroup = &orgInfo{
        ID:                int(acc.BelongGroupId),
        Username:          acc.BelongGroupUsername,
        Nickname:          acc.BelongGroupNickname,
        CreatedTimestamp:  acc.BelongGroupCreatedTimestamp,
        ModifiedTimestamp: acc.BelongGroupModifiedTimestamp,
        CurrentCnt:        int(acc.BelongGroupCurrentCnt),
    }
}

// 根据includeTeam参数决定是否包含团队信息
if req.IncludeTeam != "false" && acc.BelongTeamId > 0 {
    accountData.BelongTeam = &orgInfo{
        ID:                int(acc.BelongTeamId),
        Username:          acc.BelongTeamUsername,
        Nickname:          acc.BelongTeamNickname,
        CreatedTimestamp:  acc.BelongTeamCreatedTimestamp,
        ModifiedTimestamp: acc.BelongTeamModifiedTimestamp,
        CurrentCnt:        int(acc.BelongTeamCurrentCnt),
    }
}
```

## 测试方法

### 1. 使用测试脚本
```bash
./scripts/test_account_api_structure.sh
```

### 2. 手动测试
```bash
# 测试账户列表
curl -X GET "http://localhost:9999/api/v1/accounts?current=1&pageSize=3" \
  -H "Authorization: Bearer <your_token>"

# 测试账户详情
curl -X GET "http://localhost:9999/api/v1/accounts/0" \
  -H "Authorization: Bearer <your_token>"

# 测试不包含组织信息
curl -X GET "http://localhost:9999/api/v1/accounts?includeGroup=false&includeTeam=false" \
  -H "Authorization: Bearer <your_token>"
```

## 注意事项

1. **数据类型一致性**：确保前端和后端使用相同的数据类型（如时间戳使用int64）
2. **空值处理**：使用指针类型允许返回null值
3. **参数验证**：正确处理查询参数的大小写和默认值
4. **性能考虑**：可选的组织信息可以减少不必要的数据传输

## 相关文件

- `internal/api/account/func_list.go`: 账户列表API实现
- `internal/api/account/func_detail.go`: 账户详情API实现
- `scripts/test_account_api_structure.sh`: 数据结构测试脚本 