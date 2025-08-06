# Marriage System 婚姻管理系统

## 概述

Marriage System 是一个基于Go-Gin框架的婚姻管理系统，包含账户管理和组织信息管理功能。系统使用MySQL数据库存储数据，支持完整的CRUD操作和历史记录追踪。

## 功能特性

### 1. 账户管理
- ✅ 用户注册和登录
- ✅ 账户信息管理（创建、查询、更新、删除）
- ✅ 角色权限管理（公司经理、组经理、团队经理、普通员工）
- ✅ 账户状态管理（启用、禁用）
- ✅ 密码加密存储（MD5）
- ✅ 操作历史记录追踪

### 2. 组织信息管理
- ✅ 组织架构管理（组、团队）
- ✅ 层级关系管理（父子组织）
- ✅ 成员数量统计
- ✅ 组织信息CRUD操作
- ✅ 组织搜索和筛选

### 3. 数据库设计
- ✅ 使用marriage_system数据库
- ✅ 完整的表结构设计
- ✅ 索引优化
- ✅ 软删除机制
- ✅ 审计字段记录

## 系统架构

```
marriage_system/
├── internal/
│   ├── api/                    # API层
│   │   ├── account/           # 账户API
│   │   └── org_info/          # 组织信息API
│   ├── services/              # 服务层
│   │   ├── account/           # 账户服务
│   │   └── org_info/          # 组织信息服务
│   └── repository/mysql/      # 数据访问层
│       ├── account/           # 账户数据模型
│       ├── account_history/   # 历史记录数据模型
│       └── org_info/          # 组织信息数据模型
├── scripts/                   # 脚本文件
│   ├── install_marriage_system.sh    # 安装脚本
│   ├── stop_marriage_system.sh       # 停止脚本
│   ├── test_account_api_real.sh      # 账户API测试
│   ├── test_org_info_api.sh          # 组织信息API测试
│   └── init_marriage_system_database.sql  # 数据库初始化
└── main.go                   # 应用入口
```

## 数据库设计

### 1. 组织信息表 (org)
| 字段名 | 类型 | 说明 |
|--------|------|------|
| id | int unsigned | 主键 |
| org_id | varchar(32) | 组织ID |
| org_name | varchar(60) | 组织名称 |
| org_type | varchar(20) | 组织类型 (group/team) |
| org_level | int | 组织层级 (1-组, 2-团队) |
| parent_org_id | varchar(32) | 父组织ID |
| org_description | varchar(200) | 组织描述 |
| current_cnt | int | 当前成员数量 |
| max_cnt | int | 最大成员数量 |
| created_timestamp | bigint | 创建时间戳 |
| modified_timestamp | bigint | 修改时间戳 |

### 2. 账户表 (account)
| 字段名 | 类型 | 说明 |
|--------|------|------|
| id | int unsigned | 主键 |
| account_id | varchar(32) | 账户ID |
| username | varchar(32) | 用户名 |
| nickname | varchar(60) | 姓名 |
| password | varchar(100) | 密码(MD5) |
| phone | varchar(20) | 手机号 |
| role_type | varchar(20) | 角色类型 |
| status | varchar(20) | 状态 |
| belong_group_* | varchar/int | 所属组信息 |
| belong_team_* | varchar/int | 所属团队信息 |

### 3. 历史记录表 (account_history)
| 字段名 | 类型 | 说明 |
|--------|------|------|
| id | int unsigned | 主键 |
| history_id | varchar(32) | 历史记录ID |
| account_id | varchar(32) | 账户ID |
| operate_type | varchar(20) | 操作类型 |
| operate_timestamp | bigint | 操作时间戳 |
| content | json | 操作内容 |
| operator | varchar(60) | 操作人 |
| operator_role_type | varchar(20) | 操作人角色 |

## 快速开始

### 1. 环境要求
- Go 1.16+
- MySQL 5.7+
- Git

### 2. 安装步骤

```bash
# 1. 克隆项目
git clone <repository-url>
cd backend-marriage

# 2. 运行安装脚本
./scripts/install_marriage_system.sh
```

安装脚本会自动：
- 检查MySQL服务状态
- 创建marriage_system数据库
- 导入表结构和测试数据
- 安装Go依赖
- 生成数据库模型代码
- 编译应用
- 启动服务

### 3. 测试API

```bash
# 测试账户API
./scripts/test_account_api_real.sh

# 测试组织信息API
./scripts/test_org_info_api.sh
```

### 4. 停止服务

```bash
./scripts/stop_marriage_system.sh
```

## API接口

### 账户管理API

#### 认证相关
- `POST /api/v1/auth/login` - 用户登录
- `POST /api/v1/auth/logout` - 退出登录

#### 账户管理
- `GET /api/v1/accounts` - 获取账户列表
- `POST /api/v1/accounts` - 创建账户
- `GET /api/v1/accounts/{accountId}` - 获取账户详情
- `PUT /api/v1/accounts/{accountId}` - 更新账户
- `DELETE /api/v1/accounts/{accountId}` - 删除账户

#### 历史记录
- `GET /api/v1/account-histories` - 获取账户历史记录

### 组织信息API

#### 组织管理
- `GET /api/v1/org-infos` - 获取组织信息列表
- `POST /api/v1/org-infos` - 创建组织信息
- `GET /api/v1/org-infos/{orgId}` - 获取组织信息详情
- `PUT /api/v1/org-infos/{orgId}` - 更新组织信息
- `DELETE /api/v1/org-infos/{orgId}` - 删除组织信息

#### 层级管理
- `GET /api/v1/org-infos/{orgId}/children` - 获取子组织
- `GET /api/v1/org-infos/{orgId}/parent` - 获取父组织

## 测试数据

### 测试账户
| 用户名 | 姓名 | 角色 | 密码 | 状态 |
|--------|------|------|------|------|
| admin | 系统管理员 | company_manager | 123456 | enabled |
| company_manager | 张伟 | company_manager | 123456 | enabled |
| group_manager | 李娜 | group_manager | 123456 | enabled |
| team_manager | 王强 | team_manager | 123456 | enabled |
| employee | 赵敏 | employee | 123456 | enabled |

### 测试组织
| 组织ID | 组织名称 | 类型 | 层级 | 父组织 |
|--------|----------|------|------|--------|
| org_001 | 系统管理组 | group | 1 | - |
| org_002 | 南京-天元大厦组 | group | 1 | - |
| org_003 | 北京办公室组 | group | 1 | - |
| org_004 | 上海分公司组 | group | 1 | - |
| org_005 | 广州办公室组 | group | 1 | - |
| org_006 | 系统管理团队 | team | 2 | org_001 |
| org_007 | 营销团队A | team | 2 | org_002 |
| org_008 | 销售团队B | team | 2 | org_003 |
| org_009 | 技术团队C | team | 2 | org_004 |
| org_010 | 客服团队D | team | 2 | org_005 |

## 开发指南

### 1. 添加新的数据模型

```bash
# 1. 创建模型文件
touch internal/repository/mysql/new_model/gen_model.go

# 2. 定义模型结构
# 参考 internal/repository/mysql/account/gen_model.go

# 3. 生成数据库操作代码
go run cmd/gormgen/main.go -structs NewModel -input ./internal/repository/mysql/new_model/
```

### 2. 添加新的服务

```bash
# 1. 创建服务目录
mkdir -p internal/services/new_service

# 2. 创建服务接口
touch internal/services/new_service/service.go

# 3. 创建服务实现
touch internal/services/new_service/service_*.go
```

### 3. 添加新的API

```bash
# 1. 创建API目录
mkdir -p internal/api/new_api

# 2. 创建处理器
touch internal/api/new_api/handler.go

# 3. 创建API实现
touch internal/api/new_api/func_*.go
```

## 部署说明

### 1. 生产环境部署

```bash
# 1. 编译应用
go build -o marriage_system main.go

# 2. 创建systemd服务文件
sudo tee /etc/systemd/system/marriage-system.service << EOF
[Unit]
Description=Marriage System
After=network.target

[Service]
Type=simple
User=marriage
WorkingDirectory=/opt/marriage-system
ExecStart=/opt/marriage-system/marriage_system
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
EOF

# 3. 启动服务
sudo systemctl enable marriage-system
sudo systemctl start marriage-system
```

### 2. Docker部署

```dockerfile
FROM golang:1.16-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o marriage_system main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/marriage_system .
EXPOSE 8000
CMD ["./marriage_system"]
```

## 监控和日志

### 1. 日志文件
- 应用日志: `logs/marriage_system.log`
- 错误日志: `logs/error.log`

### 2. 监控指标
- 应用状态: `http://localhost:8000/health`
- 数据库连接: 检查MySQL连接状态
- API响应时间: 通过日志分析

## 故障排除

### 1. 常见问题

**问题**: 数据库连接失败
**解决**: 检查MySQL服务状态和连接配置

**问题**: 应用启动失败
**解决**: 检查端口占用和配置文件

**问题**: API返回500错误
**解决**: 查看应用日志定位具体错误

### 2. 日志分析

```bash
# 查看实时日志
tail -f logs/marriage_system.log

# 查看错误日志
grep "ERROR" logs/marriage_system.log

# 查看API访问日志
grep "GET\|POST\|PUT\|DELETE" logs/marriage_system.log
```

## 贡献指南

1. Fork 项目
2. 创建功能分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 打开 Pull Request

## 许可证

本项目采用 MIT 许可证 - 查看 [LICENSE](LICENSE) 文件了解详情。

## 联系方式

如有问题或建议，请通过以下方式联系：
- 提交 Issue
- 发送邮件
- 项目讨论区 