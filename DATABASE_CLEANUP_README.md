# 数据库清理脚本使用说明

## 概述

本项目提供了两个脚本来帮助进行数据库清理和管理：

1. **`run.sh`** - 主运行脚本，包含基本的数据库清理功能
2. **`cleanup_db.sh`** - 专门的数据库强制清理脚本，提供更高级的清理选项

## 脚本功能

### run.sh 增强功能

原有的 `run.sh` 脚本已经增强，现在包含：

- ✅ 自动检测和停止运行中的服务
- ✅ 强制数据库连接断开
- ✅ 数据库缓存清理
- ✅ 数据库状态重置
- ✅ 临时表清理
- ✅ 表优化
- ✅ 自动使用 `-rebuild-db -force-reseed` 参数启动

### cleanup_db.sh 新功能

新创建的 `cleanup_db.sh` 脚本提供：

- 🔥 **强制断开连接** - 强制断开所有数据库连接
- 🧹 **缓存清理** - 清理各种数据库缓存
- 🔄 **状态重置** - 重置数据库全局状态
- 🗑️ **临时表清理** - 清理所有临时表
- ⚡ **表优化** - 优化所有数据库表
- 🔢 **自增ID重置** - 重置所有表的AUTO_INCREMENT值
- 📊 **状态监控** - 显示数据库当前状态

## 使用方法

### 1. 基本使用（run.sh）

```bash
# 启动服务（自动包含数据库清理）
./run.sh
```

### 2. 高级数据库清理（cleanup_db.sh）

```bash
# 显示帮助信息
./cleanup_db.sh --help

# 执行所有清理操作
./cleanup_db.sh

# 强制执行所有操作（跳过确认）
./cleanup_db.sh -f

# 只断开数据库连接
./cleanup_db.sh -c

# 只清理缓存
./cleanup_db.sh -k

# 只重置状态
./cleanup_db.sh -s

# 只清理临时表
./cleanup_db.sh -t

# 只优化表
./cleanup_db.sh -o

# 只重置自增ID
./cleanup_db.sh -i

# 组合操作
./cleanup_db.sh -c -k -s
```

## 命令行选项

| 选项 | 长选项 | 描述 |
|------|--------|------|
| `-h` | `--help` | 显示帮助信息 |
| `-f` | `--force` | 强制执行，跳过确认提示 |
| `-c` | `--connections` | 强制断开所有数据库连接 |
| `-k` | `--cache` | 清理数据库缓存 |
| `-s` | `--status` | 重置数据库状态 |
| `-t` | `--temp` | 清理临时表 |
| `-o` | `--optimize` | 优化数据库表 |
| `-i` | `--increment` | 重置所有表的自增ID |
| `-a` | `--all` | 执行所有清理操作（默认） |

## 前置要求

### 必需软件

1. **MySQL客户端**
   ```bash
   # macOS
   brew install mysql-client
   
   # Ubuntu/Debian
   sudo apt-get install mysql-client
   
   # CentOS/RHEL
   sudo yum install mysql-client
   ```

2. **Bash shell** (大多数Linux/macOS系统默认包含)

### 数据库权限

确保数据库用户具有以下权限：
- `PROCESS` - 查看和终止连接
- `RELOAD` - 执行FLUSH命令
- `SUPER` - 设置全局变量
- `DROP` - 删除临时表
- `OPTIMIZE` - 优化表

## 安全注意事项

⚠️ **警告**: 这些脚本会执行强制性的数据库操作，请谨慎使用：

1. **生产环境**: 不建议在生产环境中使用
2. **数据备份**: 执行前请确保有数据备份
3. **权限控制**: 确保只有授权用户能执行这些脚本
4. **连接影响**: 会强制断开所有数据库连接

## 故障排除

### 常见问题

1. **MySQL客户端未安装**
   ```bash
   # 安装MySQL客户端
   brew install mysql-client  # macOS
   sudo apt-get install mysql-client  # Ubuntu
   ```

2. **权限不足**
   ```bash
   # 检查脚本权限
   ls -la run.sh cleanup_db.sh
   
   # 添加执行权限
   chmod +x run.sh cleanup_db.sh
   ```

3. **数据库连接失败**
   - 检查网络连接
   - 验证数据库配置
   - 确认防火墙设置

4. **清理操作失败**
   - 某些操作可能需要更高权限
   - 部分失败是正常的，脚本会继续执行

### 日志和调试

脚本会显示详细的执行过程，包括：
- 每个步骤的执行状态
- 成功/失败的操作
- 数据库连接信息
- 错误信息（如果有）

## 示例输出

### run.sh 执行示例

```
================================
    marriage_system 统一管理脚本
================================
检测到服务正在运行
正在停止进程 PID: 12345
进程 12345 已停止
端口已释放
已删除可执行文件
Go缓存已清理
开始强制数据库清理...
1. 强制断开所有数据库连接...
  执行: KILL 12346;
  执行: KILL 12347;
  所有连接已断开
2. 清理数据库缓存...
  数据库缓存已清理
3. 重置数据库状态...
  数据库状态已重置
4. 清理临时表...
  临时表已清理
5. 优化数据库表...
  表优化完成
强制数据库清理完成！
```

### cleanup_db.sh 执行示例

```
================================
    开始数据库强制清理
================================
目标数据库: gz-cdb-pepynap7.sql.tencentcdb.com:63623/marriage_system
执行操作: connections cache status temp optimize
强制模式: false

测试数据库连接...
数据库连接成功

警告: 即将执行数据库清理操作
这些操作可能会影响数据库性能
是否继续? (y/N): y

1. 强制断开所有数据库连接...
找到以下活动连接:
12346  marriage_system  192.168.1.100:12345  marriage_system  Sleep  120  NULL  NULL
正在断开连接...
  执行: KILL 12346;
  所有连接已断开

2. 清理数据库缓存...
  执行: FLUSH PRIVILEGES
  执行: FLUSH HOSTS
  执行: FLUSH LOGS
  执行: FLUSH STATUS
  执行: FLUSH TABLES
  执行: FLUSH TABLES WITH READ LOCK
  执行: UNLOCK TABLES
  数据库缓存已清理

[... 更多输出 ...]
```

## 联系和支持

如果在使用过程中遇到问题，请：

1. 检查脚本输出和错误信息
2. 确认前置要求是否满足
3. 验证数据库配置是否正确
4. 查看项目文档和README

---

**注意**: 这些脚本主要用于开发和测试环境，生产环境使用前请充分测试。
