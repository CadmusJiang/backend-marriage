# 端口配置说明

## 默认端口

Marriage System 项目的默认端口是 **9999**，不是 8000。

## 端口配置位置

端口配置在 `configs/constants.go` 文件中：

```go
// ProjectPort 项目端口
ProjectPort = ":9999"
```

## 正确的访问地址

### 服务地址
- **本地访问**: `http://localhost:9999`
- **API基础路径**: `http://localhost:9999/api/v1`

### Swagger文档
- **Swagger UI**: `http://localhost:9999/swagger/index.html`

### 主要API接口

#### 认证相关
- `POST http://localhost:9999/api/v1/auth/login` - 用户登录
- `POST http://localhost:9999/api/v1/auth/logout` - 退出登录

#### 账户管理
- `GET http://localhost:9999/api/v1/accounts` - 获取账户列表（批量获取）
- `GET http://localhost:9999/api/v1/accounts/{accountId}` - 获取单个账户详情
- `POST http://localhost:9999/api/v1/accounts` - 创建账户
- `PUT http://localhost:9999/api/v1/accounts/{accountId}` - 更新账户
- `GET http://localhost:9999/api/v1/account-histories` - 获取历史记录

#### 组织管理
- `GET http://localhost:9999/api/v1/org-infos` - 获取组织列表（批量获取）
- `GET http://localhost:9999/api/v1/org-infos/{orgId}` - 获取单个组织详情
- `POST http://localhost:9999/api/v1/org-infos` - 创建组织
- `PUT http://localhost:9999/api/v1/org-infos/{orgId}` - 更新组织
- `DELETE http://localhost:9999/api/v1/org-infos/{orgId}` - 删除组织
- `GET http://localhost:9999/api/v1/org-infos/{orgId}/children` - 获取子组织
- `GET http://localhost:9999/api/v1/org-infos/{orgId}/parent` - 获取父组织

#### 系统检查
- `GET http://localhost:9999/api/v1/check-db` - 检查数据库状态

## 快速测试

### 1. 检查服务状态
```bash
curl http://localhost:9999/api/v1/check-db
```

### 2. 用户登录
```bash
curl -X POST http://localhost:9999/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"123456"}'
```

### 3. 获取账户列表
```bash
curl -X GET "http://localhost:9999/api/v1/accounts?current=1&pageSize=5" \
  -H "Authorization: Bearer <your_token>"
```

### 4. 获取组织列表
```bash
curl -X GET "http://localhost:9999/api/v1/org-infos?current=1&pageSize=5" \
  -H "Authorization: Bearer <your_token>"
```

## 启动服务

### 开发环境
```bash
go run main.go
```

### 生产环境
```bash
./marriage_system
```

服务启动后会在控制台显示：
```
应用将在后台运行，端口: 9999
访问地址: http://localhost:9999
```

## 注意事项

1. **端口冲突**: 如果9999端口被占用，可以修改 `configs/constants.go` 中的 `ProjectPort` 值
2. **防火墙**: 确保防火墙允许9999端口的访问
3. **Docker**: 如果使用Docker，确保端口映射正确：`-p 9999:9999`

## 常见问题

### Q: 为什么访问8000端口没有响应？
A: 项目默认端口是9999，不是8000。请使用正确的端口号访问。

### Q: 如何修改端口？
A: 修改 `configs/constants.go` 文件中的 `ProjectPort` 值，然后重新编译运行。

### Q: 服务启动失败怎么办？
A: 检查端口是否被占用，可以使用 `lsof -i :9999` 命令检查。 