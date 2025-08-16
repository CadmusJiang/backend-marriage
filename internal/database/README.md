# Database Manager Package

## 概述

`internal/database` 包负责管理数据库的初始化、表结构管理和数据操作。这个包将原本分散在 router 包中的数据库管理逻辑集中管理，提高了代码的可维护性和职责分离。

## 功能特性

### 1. 数据库连接管理
- 测试数据库读写连接
- 连接状态验证

### 2. 表结构管理
- 自动检查并创建缺失的数据表
- 支持强制重建所有表
- 表依赖关系管理

### 3. 数据管理
- Mock 数据的插入和管理
- 数据重置功能

## 主要组件

### Manager 结构体
```go
type Manager struct {
    db     mysql.Repo
    logger *zap.Logger
}
```

### 核心方法

#### TestConnection()
测试数据库连接是否正常，包括读写连接。

#### EnsureTables()
检查并创建核心业务相关的数据表，如果表已存在则跳过。

#### RebuildTables()
强制删除所有现有表并重新创建，适用于开发环境或数据重置场景。

#### ReinsertMockData()
重新插入测试数据，按依赖顺序执行。

## 使用示例

```go
// 创建数据库管理器
dbManager := database.New(dbRepo, logger)

// 测试连接
if err := dbManager.TestConnection(); err != nil {
    return err
}

// 确保表存在
if err := dbManager.EnsureTables(); err != nil {
    return err
}

// 重建表（可选）
if rebuildDatabase {
    if err := dbManager.RebuildTables(); err != nil {
        return err
    }
}

// 插入测试数据（可选）
if forceReseed {
    if err := dbManager.ReinsertMockData(); err != nil {
        return err
    }
}
```

## 支持的数据表

- `org` - 组织表
- `org_history` - 组织历史表
- `account` - 账户表
- `account_history` - 账户历史表
- `account_org_relation` - 账户组织关系表
- `customer_authorization_record` - 客户授权记录表
- `customer_authorization_record_history` - 客户授权记录历史表
- `cooperation_store` - 合作门店表
- `cooperation_store_history` - 合作门店历史表
- `outbox_events` - 事件输出表

## 设计原则

1. **单一职责**: 专门负责数据库管理，不涉及业务逻辑
2. **依赖注入**: 通过构造函数注入依赖，便于测试和配置
3. **错误处理**: 提供详细的错误信息和日志记录
4. **事务安全**: 确保数据操作的原子性
5. **可配置性**: 支持不同的初始化策略（检查、重建、重置）

## 迁移说明

从 router 包迁移到 database 包的好处：

1. **职责分离**: router 包专注于路由配置，database 包专注于数据库管理
2. **代码复用**: 数据库管理逻辑可以在其他地方复用
3. **测试友好**: 独立的包更容易进行单元测试
4. **维护性**: 数据库相关的修改不会影响路由配置
5. **扩展性**: 可以轻松添加新的数据库管理功能
