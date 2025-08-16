# 添加组织表字段迁移脚本

这个脚本用于为 `org` 表添加 `address` 和 `extra` 字段。

## 新增字段

### 1. address 字段
- 类型: `varchar(255)`
- 默认值: `NULL`
- 说明: 地址信息

### 2. extra 字段
- 类型: `json`
- 默认值: `NULL`
- 说明: 扩展信息，支持存储JSON格式的额外数据

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
cd cmd/add_org_fields
go run main.go
```

### 3. 或者编译后运行

```bash
cd cmd/add_org_fields
go build -o add_org_fields
./add_org_fields
```

## 注意事项

1. **备份数据库**: 在运行迁移脚本之前，请务必备份您的数据库
2. **测试环境**: 建议先在测试环境中运行迁移脚本
3. **停机时间**: 添加字段通常很快，但建议在低峰期执行
4. **权限要求**: 确保数据库用户有 ALTER TABLE 权限

## 验证结果

迁移完成后，您可以运行以下SQL来验证结果：

```sql
-- 检查表结构
DESCRIBE org;

-- 检查新字段是否存在
SELECT COLUMN_NAME, DATA_TYPE, IS_NULLABLE, COLUMN_DEFAULT, COLUMN_COMMENT
FROM INFORMATION_SCHEMA.COLUMNS 
WHERE TABLE_SCHEMA = 'marriage_system' 
AND TABLE_NAME = 'org' 
AND COLUMN_NAME IN ('address', 'extra');
```

## 支持

如果在迁移过程中遇到问题，请检查：

1. 数据库连接是否正常
2. 用户权限是否足够
3. 数据库版本是否支持JSON类型
4. 是否有外键约束阻止了表修改
