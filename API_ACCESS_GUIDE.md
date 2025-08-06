# 账户管理系统API访问指南

## 访问端口

**主端口：9999**

## 完整的API访问地址

### 1. 认证相关接口

#### 用户登录
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

#### 退出登录
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

### 2. 账户管理接口

#### 获取账户列表
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

#### 创建账户
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
  "roleType": "user"
}
```

#### 获取账户详情
```
GET http://localhost:9999/api/v1/accounts/{accountId}
```

#### 更新账户
```
PUT http://localhost:9999/api/v1/accounts/{accountId}
```

**请求体：**
```json
{
  "nickname": "新姓名",
  "phone": "13800138007",
  "belongGroup": "运营组",
  "belongTeam": "市场团队",
  "status": "active"
}
```

#### 获取账户历史记录
```
GET http://localhost:9999/api/v1/account-histories?current=1&pageSize=10
```

**查询参数：**
- `accountId`: 账户ID
- `accountUsername`: 账户用户名
- `current`: 当前页码 (默认: 1)
- `pageSize`: 每页数量 (默认: 10)

## 测试方法

### 1. 使用curl测试

#### 登录测试
```bash
curl -X POST http://localhost:8000/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "admin",
    "password": "123456"
  }'
```

#### 获取账户列表测试
```bash
curl -X GET "http://localhost:8000/api/v1/accounts?current=1&pageSize=5" \
  -H "Authorization: Bearer YOUR_TOKEN"
```

### 2. 使用测试脚本

```bash
./scripts/test_account_api.sh
```

### 3. 使用Postman或其他API测试工具

1. 导入API集合
2. 设置基础URL: `http://localhost:8000`
3. 先调用登录接口获取token
4. 在其他请求的Header中添加: `Authorization: Bearer YOUR_TOKEN`

## 认证说明

1. **登录接口** - 无需认证，直接调用
2. **其他接口** - 需要在请求头中包含有效的token
3. **Token格式** - `Authorization: Bearer YOUR_TOKEN`

## 测试账号

- **用户名**: `admin`, `user001`
- **密码**: `123456`

## 注意事项

1. 所有需要认证的API都需要在请求头中包含有效的token
2. 当前使用Mock数据，无需数据库
3. 分页查询支持多种筛选条件
4. 所有时间戳使用Unix时间戳格式
5. 响应格式统一，包含success、message、data等字段

## 错误处理

### 常见错误码
- `10103`: 参数绑定错误
- `20206`: 登录错误
- `20213`: 账户不存在

### 错误响应格式
```json
{
  "code": 10103,
  "message": "参数绑定错误"
}
``` 