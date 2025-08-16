# 枚举字段迁移脚本

这个脚本用于将数据库中的整数枚举字段迁移到字符串枚举字段，并使用MySQL的ENUM约束。

## 迁移的字段

### 1. account.status
- 从: `tinyint(1)` (1:enabled, 2:disabled)
- 到: `enum('enabled','disabled')`

### 2. org.org_type
- 从: `tinyint(1)` (1:group, 2:team)
- 到: `enum('group','team')`

### 3. org.status
- 从: `tinyint(1)` (1:enabled, 2:disabled)
- 到: `enum('enabled','disabled')`

### 4. account_org_relation.relation_type
- 从: `tinyint(1)` (1:belong, 2:manage)
- 到: `enum('belong','manage')`

### 5. account_org_relation.status
- 从: `tinyint(1)` (1:active, 2:inactive)
- 到: `enum('active','inactive')`

### 6. customer_authorization_record 状态字段
- `authorization_status`: 从 `tinyint(1)` 到 `enum('authorized','unauthorized')`
- `completion_status`: 从 `tinyint(1)` 到 `enum('complete','incomplete')`
- `assignment_status`: 从 `tinyint(1)` 到 `enum('assigned','unassigned')`
- `payment_status`: 从 `tinyint(1)` 到 `enum('paid','unpaid')`

### 7. outbox_events.status
- 从: `tinyint` (0:未发布, 1:已发布)
- 到: `enum('unpublished','published')`

## 使用方法

### 1. 设置环境变量

```bash
export DB_HOST=localhost
export DB_PORT=3306
export DB_USER=root
export DB_PASSWORD=your_password
export DB_NAME=marriage_system
```

### 2. 运行迁移脚本

```bash
cd cmd/migrate_enums
go run main.go
```

### 3. 或者编译后运行

```bash
cd cmd/migrate_enums
go build -o migrate_enums
./migrate_enums
```

## 注意事项

1. **备份数据库**: 在运行迁移脚本之前，请务必备份您的数据库
2. **测试环境**: 建议先在测试环境中运行迁移脚本
3. **停机时间**: 迁移过程中可能需要短暂的停机时间
4. **回滚计划**: 如果迁移失败，请准备好回滚计划

## 迁移过程

迁移脚本会按以下步骤进行：

1. 为每个字段添加新的ENUM列
2. 将旧数据从整数转换为对应的字符串值
3. 删除旧的整数列
4. 将新列重命名为原来的列名

## 数据映射

### 状态字段
- `1` → `'enabled'`
- `2` → `'disabled'`

### 组织类型
- `1` → `'group'`
- `2` → `'team'`

### 关联类型
- `1` → `'belong'`
- `2` → `'manage'`

### 关联状态
- `1` → `'active'`
- `2` → `'inactive'`

### 客资状态字段
- `1` → `'authorized'` / `'complete'` / `'assigned'` / `'paid'`
- `0` → `'unauthorized'` / `'incomplete'` / `'unassigned'` / `'unpaid'`

### 发布状态
- `0` → `'unpublished'`
- `1` → `'published'`

## 错误处理

如果迁移过程中出现错误，脚本会：

1. 记录详细的错误信息
2. 停止迁移过程
3. 保持数据库在一致状态

## 验证迁移结果

迁移完成后，您可以运行以下SQL来验证结果：

```sql
-- 检查account表
SELECT status, COUNT(*) FROM account GROUP BY status;

-- 检查org表
SELECT org_type, status, COUNT(*) FROM org GROUP BY org_type, status;

-- 检查account_org_relation表
SELECT relation_type, status, COUNT(*) FROM account_org_relation GROUP BY relation_type, status;

-- 检查customer_authorization_record表
SELECT authorization_status, completion_status, assignment_status, payment_status, COUNT(*) 
FROM customer_authorization_record 
GROUP BY authorization_status, completion_status, assignment_status, payment_status;

-- 检查outbox_events表
SELECT status, COUNT(*) FROM outbox_events GROUP BY status;
```

## 支持

如果在迁移过程中遇到问题，请检查：

1. 数据库连接是否正常
2. 用户权限是否足够
3. 数据库版本是否支持ENUM类型
4. 是否有外键约束阻止了列修改
