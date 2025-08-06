-- 插入详细的Mock账户历史记录数据
INSERT INTO `go_gin_api`.`account_history` (
  `history_id`, `account_id`, `operate_type`, `operate_timestamp`, `content`,
  `operator_username`, `operator_nickname`, `operator_role_type`,
  `created_user`, `updated_user`
) VALUES 
-- 账户6的历史记录
('1', '6', 'modified', 1705923000, '{"roleType": {"old": "员工", "new": "小队管理员"}, "belongTeam": {"old": "无", "new": "营销团队A"}, "status": {"old": "enabled", "new": "disabled"}}', 'zhangwei', '张伟', 'company_manager', '系统管理员', '系统管理员'),

('2', '6', 'modified', 1705754700, '{"belongGroup": {"old": "南京-天元大厦组", "new": "南京-南京南站组"}, "belongTeam": {"old": "营销团队A", "new": "营销团队C"}}', 'liming', '李明', 'team_manager', '系统管理员', '系统管理员'),

('3', '6', 'modified', 1705565700, '{"status": {"old": "enabled", "new": "disabled"}}', 'wangfang', '王芳', 'team_manager', '系统管理员', '系统管理员'),

('4', '6', 'modified', 1705303200, '{"phone": {"old": "13800138000", "new": "13900139000"}, "belongTeam": {"old": "营销团队A", "new": "营销团队B"}}', 'liuqiang', '刘强', 'team_manager', '系统管理员', '系统管理员'),

('5', '6', 'modified', 1705032000, '{"nickname": {"old": "张三", "new": "张明"}, "roleType": {"old": "员工", "new": "小队管理员"}}', 'chenjing', '陈静', 'team_manager', '系统管理员', '系统管理员'),

('6', '6', 'created', 1704877800, '{"username": {"old": "", "new": "employee001"}, "nickname": {"old": "", "new": "张三"}, "roleType": {"old": "", "new": "员工"}, "phone": {"old": "", "new": "13800138000"}, "belongGroup": {"old": "", "new": "南京-天元大厦组"}, "belongTeam": {"old": "", "new": "营销团队A"}}', 'admin', '系统管理员', 'company_manager', '系统管理员', '系统管理员'),

('7', '6', 'modified', 1704709200, '{"belongGroup": {"old": "南京-天元大厦组", "new": "南京-夫子庙组"}, "belongTeam": {"old": "营销团队A", "new": "营销团队D"}}', 'zhaomin', '赵敏', 'team_manager', '系统管理员', '系统管理员'),

('8', '6', 'modified', 1704457500, '{"status": {"old": "disabled", "new": "enabled"}, "reason": {"old": "", "new": "问题已解决，恢复账户"}}', 'liming', '李明', 'team_manager', '系统管理员', '系统管理员'),

('9', '6', 'modified', 1704265800, '{"password": {"old": "****", "new": "******"}}', 'zhangwei', '张伟', 'company_manager', '系统管理员', '系统管理员'),

('10', '6', 'modified', 1704081600, '{"password": {"old": "****", "new": "******"}}', 'wangfang', '王芳', 'team_manager', '系统管理员', '系统管理员'),

-- 账户7的历史记录
('11', '7', 'modified', 1705923000, '{"roleType": {"old": "员工", "new": "小队管理员"}}', 'zhangwei', '张伟', 'company_manager', '系统管理员', '系统管理员'),

('12', '7', 'created', 1704877800, '{"username": {"old": "", "new": "employee002"}, "nickname": {"old": "", "new": "孙丽"}, "roleType": {"old": "", "new": "员工"}, "phone": {"old": "", "new": "13800138007"}, "belongGroup": {"old": "", "new": "南京-南京南站组"}, "belongTeam": {"old": "", "new": "营销团队B"}}', 'admin', '系统管理员', 'company_manager', '系统管理员', '系统管理员'),

-- 账户8的历史记录
('13', '8', 'modified', 1705754700, '{"belongTeam": {"old": "营销团队A", "new": "营销团队C"}}', 'liming', '李明', 'team_manager', '系统管理员', '系统管理员'),

('14', '8', 'created', 1704877800, '{"username": {"old": "", "new": "employee003"}, "nickname": {"old": "", "new": "周杰"}, "roleType": {"old": "", "new": "员工"}, "phone": {"old": "", "new": "13800138008"}, "belongGroup": {"old": "", "new": "南京-天元大厦组"}, "belongTeam": {"old": "", "new": "营销团队A"}}', 'admin', '系统管理员', 'company_manager', '系统管理员', '系统管理员'),

-- 账户9的历史记录
('15', '9', 'modified', 1705565700, '{"status": {"old": "enabled", "new": "disabled"}}', 'wangfang', '王芳', 'team_manager', '系统管理员', '系统管理员'),

('16', '9', 'created', 1704877800, '{"username": {"old": "", "new": "employee004"}, "nickname": {"old": "", "new": "吴婷"}, "roleType": {"old": "", "new": "员工"}, "phone": {"old": "", "new": "13800138009"}, "belongGroup": {"old": "", "new": "南京-南京南站组"}, "belongTeam": {"old": "", "new": "营销团队B"}}', 'admin', '系统管理员', 'company_manager', '系统管理员', '系统管理员'),

-- 账户10的历史记录
('17', '10', 'created', 1704877800, '{"username": {"old": "", "new": "employee005"}, "nickname": {"old": "", "new": "郑华"}, "roleType": {"old": "", "new": "员工"}, "phone": {"old": "", "new": "13800138010"}, "belongGroup": {"old": "", "new": "南京-天元大厦组"}, "belongTeam": {"old": "", "new": "营销团队A"}}', 'admin', '系统管理员', 'company_manager', '系统管理员', '系统管理员'); 