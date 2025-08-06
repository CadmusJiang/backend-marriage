#### go_gin_api.account_history 
账户历史记录表

| 序号 | 名称 | 描述 | 类型 | 键 | 为空 | 额外 | 默认值 |
| :--: | :--: | :--: | :--: | :--: | :--: | :--: | :--: |
| 1 | id | 主键 | int unsigned | PRI | NO | auto_increment |  |
| 2 | history_id | 历史记录ID | varchar(32) | UNI | NO |  |  |
| 3 | account_id | 账户ID | varchar(32) |  | NO |  |  |
| 4 | operate_type | 操作类型 | varchar(20) |  | NO |  |  |
| 5 | operate_timestamp | 操作时间戳 | timestamp |  | NO |  | CURRENT_TIMESTAMP |
| 6 | content | 操作内容 | json |  | YES |  | NULL |
| 7 | operator | 操作人 | varchar(60) |  | NO |  |  |
| 8 | operator_role_type | 操作人角色 | varchar(20) |  | NO |  |  |
| 9 | is_deleted | 是否删除 1:是  -1:否 | tinyint(1) |  | NO |  | -1 |
| 10 | created_at | 创建时间 | timestamp |  | NO | DEFAULT_GENERATED | CURRENT_TIMESTAMP |
| 11 | created_user | 创建人 | varchar(60) |  | NO |  |  |
| 12 | updated_at | 更新时间 | timestamp |  | NO | DEFAULT_GENERATED on update CURRENT_TIMESTAMP | CURRENT_TIMESTAMP |
| 13 | updated_user | 更新人 | varchar(60) |  | NO |  |  | 