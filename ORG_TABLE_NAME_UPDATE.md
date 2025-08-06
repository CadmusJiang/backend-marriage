# 组织表名更新说明

## 更新概述

将组织信息表名从 `org_info` 更改为 `org`，同时保留内部字段 `org_name` 不变。

## 更改内容

### 1. 数据库表结构
- **表名**: `org_info` → `org`
- **字段**: 保持 `org_name` 字段不变
- **数据**: 保留所有现有数据
- **索引**: 保留所有现有索引

### 2. 代码文件更新

#### 数据库相关文件
- `internal/proposal/tablesqls/table_org_info.go` - 更新表创建SQL
- `scripts/init_marriage_system_database.sql` - 更新初始化脚本
- `scripts/update_table_structure.sql` - 更新表结构更新脚本
- `internal/api/helper/func_check_db.go` - 更新数据库检查表名

#### 模型文件
- `internal/repository/mysql/org_info/gen_model.go` - 添加表名映射方法

#### 脚本文件
- `scripts/install_marriage_system.sh` - 更新安装脚本输出
- `scripts/install_dev.sh` - 更新开发环境安装脚本
- `scripts/test_org_info_api.sh` - 更新测试脚本
- `internal/render/install/execute.go` - 更新安装执行器

#### 文档文件
- `MARRIAGE_SYSTEM_README.md` - 更新文档说明

### 3. 新增文件
- `scripts/migrate_org_info_to_org.sql` - 数据库迁移脚本
- `scripts/update_org_table_name.sh` - 迁移执行脚本
- `ORG_TABLE_NAME_UPDATE.md` - 本说明文档

## 迁移步骤

### 方法一：使用迁移脚本（推荐）

```bash
# 1. 备份数据库
mysqldump -h HOST -P PORT -u USER -p marriage_system > backup_$(date +%Y%m%d_%H%M%S).sql

# 2. 执行迁移
./scripts/update_org_table_name.sh -h localhost -P 3306 -u root -p123456

# 3. 重新生成模型代码
go run cmd/gormgen/main.go -structs OrgInfo -input ./internal/repository/mysql/org_info/

# 4. 重新编译应用
go build -o bin/marriage-system main.go
```

### 方法二：手动执行

```sql
-- 1. 备份数据库
mysqldump -h HOST -P PORT -u USER -p marriage_system > backup.sql

-- 2. 执行迁移
USE marriage_system;
RENAME TABLE `org_info` TO `org`;

-- 3. 验证结果
SELECT COUNT(*) FROM `org`;
```

## 注意事项

1. **备份**: 执行迁移前请务必备份数据库
2. **兼容性**: 此更改向后兼容，不会影响现有数据
3. **代码生成**: 迁移后需要重新生成数据库模型代码
4. **应用重启**: 需要重新编译并重启应用

## 验证步骤

1. 检查数据库表名是否正确：
   ```sql
   SHOW TABLES LIKE 'org';
   ```

2. 检查数据是否完整：
   ```sql
   SELECT COUNT(*) FROM org;
   ```

3. 测试API功能：
   ```bash
   ./scripts/test_org_info_api.sh
   ```

## 回滚方案

如果需要回滚，可以执行：

```sql
USE marriage_system;
RENAME TABLE `org` TO `org_info`;
```

然后恢复相关的代码文件。

## 影响范围

- ✅ 数据库表名
- ✅ 初始化脚本
- ✅ 测试脚本
- ✅ 文档说明
- ✅ 模型代码（通过TableName方法映射）
- ❌ 包名和目录结构（保持不变）
- ❌ API接口路径（保持不变）
- ❌ 业务逻辑（保持不变）

## 总结

此次更新将组织信息表名从 `org_info` 改为 `org`，同时保留了所有现有数据和功能。更新过程安全且可回滚，不会影响系统的正常运行。 