package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// 数据库连接
	dsn := "marriage_system:19970901Zyt@tcp(gz-cdb-pepynap7.sql.tencentcdb.com:63623)/marriage_system?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("数据库连接失败:", err)
	}
	defer db.Close()

	// 测试连接
	err = db.Ping()
	if err != nil {
		log.Fatal("数据库连接测试失败:", err)
	}
	fmt.Println("✅ 数据库连接成功")

	// 1. company_manager不能有归属组和团队
	fmt.Println("\n1. 修复 company_manager 归属规则...")
	_, err = db.Exec("UPDATE account SET belong_group_id = NULL, belong_team_id = NULL WHERE role_type = 'company_manager'")
	if err != nil {
		log.Fatal("修复 company_manager 失败:", err)
	}
	fmt.Println("✅ company_manager 归属规则已修复")

	// 2. group_manager不能有归属团队，但必须有归属组
	fmt.Println("\n2. 修复 group_manager 归属规则...")
	_, err = db.Exec("UPDATE account SET belong_team_id = NULL WHERE role_type = 'group_manager'")
	if err != nil {
		log.Fatal("修复 group_manager 失败:", err)
	}
	fmt.Println("✅ group_manager 归属规则已修复")

	// 验证修复结果
	fmt.Println("\n3. 验证修复结果...")
	rows, err := db.Query(`
		SELECT 
			username,
			nickname,
			role_type,
			belong_group_id,
			belong_team_id,
			CASE 
				WHEN role_type = 'company_manager' AND belong_group_id IS NULL AND belong_team_id IS NULL THEN '✅ 正确'
				WHEN role_type = 'group_manager' AND belong_group_id IS NOT NULL AND belong_team_id IS NULL THEN '✅ 正确'
				WHEN role_type = 'team_manager' AND belong_group_id IS NOT NULL AND belong_team_id IS NOT NULL THEN '✅ 正确'
				WHEN role_type = 'employee' AND belong_group_id IS NOT NULL THEN '✅ 正确'
				ELSE '❌ 错误'
			END as rule_check
		FROM account 
		ORDER BY role_type, username
	`)
	if err != nil {
		log.Fatal("查询验证结果失败:", err)
	}
	defer rows.Close()

	fmt.Println("用户名 | 昵称 | 角色类型 | 归属组ID | 归属团队ID | 规则检查")
	fmt.Println("------|------|----------|----------|----------|----------")
	for rows.Next() {
		var username, nickname, roleType, ruleCheck string
		var belongGroupId, belongTeamId sql.NullInt32
		err := rows.Scan(&username, &nickname, &roleType, &belongGroupId, &belongTeamId, &ruleCheck)
		if err != nil {
			log.Fatal("读取验证结果失败:", err)
		}
		groupId := "NULL"
		teamId := "NULL"
		if belongGroupId.Valid {
			groupId = fmt.Sprintf("%d", belongGroupId.Int32)
		}
		if belongTeamId.Valid {
			teamId = fmt.Sprintf("%d", belongTeamId.Int32)
		}
		fmt.Printf("%s | %s | %s | %s | %s | %s\n",
			username, nickname, roleType, groupId, teamId, ruleCheck)
	}

	fmt.Println("\n✅ 账户归属规则修复完成！")
}
