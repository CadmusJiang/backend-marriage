#### marriage_system.account 
账户表

| 序号 | 名称 | 描述 | 类型 | 键 | 为空 | 额外 | 默认值 |
| :--: | :--: | :--: | :--: | :--: | :--: | :--: | :--: |
| 1 | id | 主键 | int unsigned | PRI | NO | auto_increment |  |
| 2 | username | 用户名 | varchar(32) | UNI | NO |  |  |
| 3 | nickname | 姓名 | varchar(60) |  | NO |  |  |
| 4 | password | 密码 | varchar(100) |  | NO |  |  |
| 5 | phone | 手机号 | varchar(20) |  | YES |  |  |
| 6 | role_type | 角色类型 | varchar(20) |  | NO |  | employee |
| 7 | status | 状态 | varchar(20) |  | NO |  | enabled |
| 8 | belong_group_id | 所属组ID | int unsigned | MUL | YES |  |  |
| 9 | belong_team_id | 所属团队ID | int unsigned | MUL | YES |  |  |
| 10 | last_login_timestamp | 最后登录时间戳 | bigint |  | YES |  |  |
| 11 | created_timestamp | 创建时间戳 | bigint |  | NO |  |  |
| 12 | modified_timestamp | 修改时间戳 | bigint |  | NO |  |  |
| 13 | created_user | 创建人 | varchar(60) |  | NO |  |  |
| 14 | updated_user | 更新人 | varchar(60) |  | NO |  |  |
