# Groups API 设计文档

## 概述

本文档描述了新实现的 `/api/v1/groups` API，该API专门用于获取组织中的groups信息。

## API设计理念

### 1. 资源导向设计
- **URL**: `/api/v1/groups` - 直接表达获取groups的意图
- **方法**: GET - 符合RESTful规范
- **语义**: 清晰明确，客户端可以直接理解这是获取groups的接口

### 2. 数据结构设计

#### 使用自增ID的优势
- **性能**: 自增ID查询性能更好
- **存储**: 节省存储空间
- **索引**: 主键索引效率高
- **关联**: 外键关联更简单

#### 组织结构设计
采用 **parent_org_id + org_path 混合方案**：

**parent_org_id (推荐)**
- ✅ 查询父子关系简单直接
- ✅ 支持灵活的树形结构
- ✅ 便于移动节点位置
- ✅ 数据库查询性能好
- ❌ 查询完整路径需要递归

**org_path (辅助)**
- ✅ 查询完整路径简单
- ✅ 便于按层级筛选
- ✅ 支持路径模式匹配
- ❌ 移动节点时需要更新所有子节点路径

## API接口详情

### 请求格式
```
GET /api/v1/groups?current=1&pageSize=10&orgName=系统&orgLevel=1&parentOrgId=0
```

### 查询参数
| 参数名 | 类型 | 必填 | 默认值 | 说明 |
|--------|------|------|--------|------|
| current | int | 否 | 1 | 当前页码 |
| pageSize | int | 否 | 10 | 每页数量 |
| orgName | string | 否 | - | 组织名称搜索 |
| orgLevel | int | 否 | - | 组织层级筛选 |
| parentOrgId | int | 否 | - | 父组织ID筛选 |

### 响应格式
```json
{
  "data": [
    {
      "id": 1,
      "orgName": "系统管理组",
      "orgType": 1,
      "orgPath": "/system/admin",
      "orgLevel": 1,
      "currentCnt": 1,
      "maxCnt": 10,
      "status": 1,
      "createdTimestamp": 1705123200,
      "modifiedTimestamp": 1705123200
    }
  ],
  "total": 5,
  "success": true,
  "pageSize": 10,
  "current": 1
}
```

### 字段说明
| 字段名 | 类型 | 说明 |
|--------|------|------|
| id | uint64 | 组织ID (自增ID) |
| orgName | string | 组织名称 |
| orgType | int32 | 组织类型: 1-group, 2-team |
| orgPath | string | 组织路径 |
| orgLevel | int32 | 组织层级 |
| currentCnt | int32 | 当前成员数量 |
| maxCnt | int32 | 最大成员数量 |
| status | int32 | 状态: 1-启用, 0-禁用 |
| createdTimestamp | int64 | 创建时间戳 |
| modifiedTimestamp | int64 | 修改时间戳 |

## 与现有API的关系

### 1. `/api/v1/groups` vs `/api/v1/org-infos`
- **groups**: 专门获取groups信息，固定orgType=1
- **org-infos**: 获取所有组织信息，支持多种类型

### 2. 设计优势
- **语义清晰**: groups API直接表达意图
- **性能优化**: 固定类型查询，索引效率更高
- **扩展性好**: 未来可添加 `/api/v1/teams` 等
- **缓存友好**: 不同资源类型可独立缓存

## 行业最佳实践

### 1. RESTful设计
- 资源导向的URL设计
- 使用HTTP方法表达操作意图
- 统一的响应格式

### 2. 分页设计
- 支持页码和页大小
- 返回总数便于前端分页
- 默认值设置合理

### 3. 搜索和筛选
- 支持多字段搜索
- 参数可选，不影响基础功能
- 查询参数语义明确

## 测试用例

### 1. 基础功能测试
```bash
# 获取所有groups
GET /api/v1/groups

# 分页获取
GET /api/v1/groups?current=1&pageSize=5

# 按名称搜索
GET /api/v1/groups?orgName=系统

# 按层级筛选
GET /api/v1/groups?orgLevel=1
```

### 2. 边界测试
```bash
# 空结果
GET /api/v1/groups?orgName=不存在的组织

# 大页码
GET /api/v1/groups?current=999&pageSize=100

# 无效参数
GET /api/v1/groups?orgLevel=invalid
```

## 性能考虑

### 1. 数据库索引
- `org_type` 索引：快速筛选groups
- `org_name` 索引：支持名称搜索
- `org_level` 索引：支持层级筛选
- `parent_org_id` 索引：支持父子关系查询

### 2. 缓存策略
- 可对groups列表进行缓存
- 按查询参数组合缓存
- 设置合理的缓存过期时间

### 3. 查询优化
- 使用固定条件 `org_type = 1` 提高查询效率
- 分页查询避免大量数据返回
- 只返回必要字段

## 未来扩展

### 1. 可能的扩展
- `/api/v1/teams` - 专门获取teams
- `/api/v1/departments` - 专门获取departments
- `/api/v1/groups/{id}/members` - 获取group成员

### 2. 功能增强
- 支持排序参数
- 支持更多搜索条件
- 支持批量操作

## 总结

新的 `/api/v1/groups` API 遵循了RESTful设计原则，使用自增ID和合理的组织结构设计，提供了清晰、高效、可扩展的groups管理功能。这种设计既符合行业惯例，又满足了实际业务需求。 