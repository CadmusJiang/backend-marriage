# 账户管理系统更新总结

## 概述

根据前端mock数据结构，对后端账户管理系统进行了全面更新，确保前后端数据结构一致，并实现了自动化的数据库初始化。

## 主要更新内容

### 1. 数据结构统一

#### 前端数据结构
```typescript
{
  id: number,           // 账户ID
  username: string,     // 用户名
  nickname: string,     // 姓名
  phone: string,        // 手机号
  roleType: number,     // 角色类型: 0-4
  belongGroup: string,  // 所属组
  belongTeam: string,   // 所属团队
  status: string,       // 状态: active/inactive
  createdAt: string,    // 创建时间
  updatedAt: string     // 更新时间
}
```

#### 后端数据结构
```go
type accountData struct {
    ID          int    `json:"id"`          // 账户ID
    Username    string `json:"username"`    // 用户名
    Nickname    string `json:"nickname"`    // 姓名
    Phone       string `json:"phone"`       // 手机号
    RoleType    int    `json:"roleType"`    // 角色类型: 0-4
    BelongGroup string `json:"belongGroup"` // 所属组
    BelongTeam  string `json:"belongTeam"`  // 所属团队
    Status      string `json:"status"`      // 状态: active/inactive
    CreatedAt   string `json:"createdAt"`   // 创建时间
    UpdatedAt   string `json:"updatedAt"`   // 更新时间
}
```

### 2. 角色类型映射

| 数值 | 角色类型 | 说明 |
|------|----------|------|
| 0 | 超级管理员 | 系统最高权限 |
| 1 | 公司管理员 | 公司级管理权限 |
| 2 | 组管理员 | 组级管理权限 |
| 3 | 小队管理员 | 小队管理权限 |
| 4 | 员工 | 普通员工权限 |

### 3. 状态类型

| 状态值 | 说明 |
|--------|------|
| active | 账户启用 |
| inactive | 账户禁用 |

### 4. 数据库表结构更新

#### account表字段
- `role_type`: 改为int类型，存储0-4的数值
- `status`: 保持varchar类型，存储'active'或'inactive'
- 移除了不必要的时间戳字段，使用标准的created_at和updated_at

#### account_history表
- 保持JSON格式存储操作内容
- 支持复杂的历史记录数据结构

### 5. 自动化初始化

#### 新增文件
- `internal/proposal/tablesqls/table_account.go` - 账户表SQL生成
- `internal/proposal/tablesqls/table_account_history.go` - 历史记录表SQL生成

#### 更新文件
- `internal/render/install/execute.go` - 添加account和account_history到初始化列表

#### 初始化流程
项目启动时会自动：
1. 创建account表
2. 创建account_history表
3. 插入预设的mock数据

### 6. Mock数据更新

#### 账户数据
- 9个预设账户，涵盖所有角色类型
- 包含超级管理员、公司管理员、组管理员、小队管理员、员工
- 分布在不同组和团队中

#### 历史记录数据
- 18条历史记录
- 包含账户创建、角色变更、团队调整、状态变更等操作
- 使用JSON格式存储详细的操作内容

### 7. API接口更新

#### 获取账户列表
- 支持按角色类型筛选（数值比较）
- 支持按状态筛选
- 支持按组和团队筛选
- 支持分页

#### 获取账户详情
- 支持通过ID或username查找
- 返回完整的账户信息

#### 创建账户
- 支持设置角色类型（数值）
- 支持设置状态
- 支持设置所属组和团队

#### 更新账户
- 支持更新基本信息
- 支持更新状态
- 支持更新所属组和团队

### 8. 测试账号

| 用户名 | 密码 | 角色类型 | 说明 |
|--------|------|----------|------|
| admin | 123456 | 0 | 超级管理员 |
| company_manager | 123456 | 1 | 公司管理员 |
| group_manager | 123456 | 2 | 组管理员 |
| team_manager | 123456 | 3 | 小队管理员 |
| employee1 | 123456 | 4 | 员工 |

## 技术实现

### 1. 类型转换处理
```go
// 角色类型字符串转数值
roleTypeInt, err := strconv.Atoi(req.RoleType)
if err == nil {
    for _, acc := range filteredAccounts {
        if acc.RoleType == roleTypeInt {
            filtered = append(filtered, acc)
        }
    }
}
```

### 2. 数据库初始化
```go
installTableList := map[string]map[string]string{
    "account": {
        "table_sql":      tablesqls.CreateAccountTableSql(),
        "table_data_sql": tablesqls.CreateAccountTableDataSql(),
    },
    "account_history": {
        "table_sql":      tablesqls.CreateAccountHistoryTableSql(),
        "table_data_sql": tablesqls.CreateAccountHistoryTableDataSql(),
    },
}
```

### 3. 数据结构验证
- 前后端数据结构完全一致
- 角色类型使用数值存储和传输
- 状态使用字符串存储和传输
- 时间格式统一使用ISO 8601格式

## 部署说明

### 1. 首次部署
```bash
# 启动项目，会自动创建表和插入数据
go run main.go
```

### 2. 验证数据
```sql
-- 查看账户数据
SELECT * FROM account WHERE is_deleted = -1;

-- 查看历史记录
SELECT * FROM account_history WHERE is_deleted = -1;
```

### 3. API测试
```bash
# 测试登录
curl -X POST http://localhost:8000/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username": "admin", "password": "123456"}'

# 测试获取账户列表
curl "http://localhost:8000/api/v1/accounts?current=1&pageSize=10"

# 测试获取账户详情
curl "http://localhost:8000/api/v1/accounts/0"
```

## 注意事项

1. **角色类型**: 前端传递数值，后端存储数值，确保类型一致性
2. **状态值**: 使用'active'和'inactive'字符串
3. **时间格式**: 统一使用ISO 8601格式
4. **密码**: 所有测试账号密码都是'123456'的MD5值
5. **软删除**: 使用is_deleted字段实现软删除机制

## 后续优化建议

1. 添加角色类型和状态的枚举定义
2. 实现真实的用户认证和授权
3. 添加数据验证和错误处理
4. 实现真实的数据库操作
5. 添加单元测试和集成测试 