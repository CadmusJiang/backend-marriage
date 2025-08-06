package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// 腾讯云数据库连接信息
	dsn := "marriage_system:19970901Zyt@tcp(gz-cdb-pepynap7.sql.tencentcdb.com:63623)/marriage_system?charset=utf8mb4&parseTime=True&loc=Local"

	fmt.Println("正在连接腾讯云数据库...")
	fmt.Println("DSN:", dsn)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("连接数据库失败:", err)
	}
	defer db.Close()

	// 测试连接
	err = db.Ping()
	if err != nil {
		log.Fatal("Ping数据库失败:", err)
	}

	fmt.Println("✅ 数据库连接成功！")

	// 检查org表是否存在
	var tableExists int
	err = db.QueryRow("SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = 'marriage_system' AND table_name = 'org'").Scan(&tableExists)
	if err != nil {
		log.Fatal("查询表失败:", err)
	}

	if tableExists > 0 {
		fmt.Println("✅ org表存在")

		// 检查数据
		var count int
		err = db.QueryRow("SELECT COUNT(*) FROM org").Scan(&count)
		if err != nil {
			log.Fatal("查询数据失败:", err)
		}

		fmt.Printf("📊 org表中有 %d 条记录\n", count)

		// 显示现有数据
		rows, err := db.Query("SELECT id, org_name, org_type, org_path FROM org ORDER BY id")
		if err != nil {
			log.Fatal("查询数据失败:", err)
		}
		defer rows.Close()

		fmt.Println("\n现有数据:")
		fmt.Println("ID | 名称 | 类型 | 路径")
		fmt.Println("---|------|------|------")

		for rows.Next() {
			var id int
			var name, path string
			var orgType int

			err := rows.Scan(&id, &name, &orgType, &path)
			if err != nil {
				log.Fatal("读取数据失败:", err)
			}

			typeName := "组"
			if orgType == 2 {
				typeName = "团队"
			}

			fmt.Printf("%d | %s | %s | %s\n", id, name, typeName, path)
		}
	} else {
		fmt.Println("❌ org表不存在")
	}

	// 检查account表是否存在
	fmt.Println("\n检查account表...")
	err = db.QueryRow("SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = 'marriage_system' AND table_name = 'account'").Scan(&tableExists)
	if err != nil {
		log.Fatal("查询account表失败:", err)
	}

	if tableExists > 0 {
		fmt.Println("✅ account表存在")

		// 检查表结构
		fmt.Println("\naccount表结构:")
		rows, err := db.Query("DESCRIBE account")
		if err != nil {
			log.Fatal("查询表结构失败:", err)
		}
		defer rows.Close()

		fmt.Println("字段名 | 类型 | 是否为空 | 键 | 默认值")
		fmt.Println("------|------|----------|----|--------")

		for rows.Next() {
			var field, fieldType, null, key, defaultValue, extra sql.NullString
			err := rows.Scan(&field, &fieldType, &null, &key, &defaultValue, &extra)
			if err != nil {
				log.Fatal("读取表结构失败:", err)
			}
			fmt.Printf("%s | %s | %s | %s | %s\n",
				field.String, fieldType.String, null.String, key.String, defaultValue.String)
		}

		// 检查account数据
		var count int
		err = db.QueryRow("SELECT COUNT(*) FROM account").Scan(&count)
		if err != nil {
			log.Fatal("查询account数据失败:", err)
		}

		fmt.Printf("\n📊 account表中有 %d 条记录\n", count)

		// 查询company_manager用户
		fmt.Println("\n查询company_manager用户...")
		var username, password, status string
		err = db.QueryRow("SELECT username, password, status FROM account WHERE username = 'company_manager'").Scan(&username, &password, &status)
		if err != nil {
			if err == sql.ErrNoRows {
				fmt.Println("❌ company_manager用户不存在")
			} else {
				log.Fatal("查询用户失败:", err)
			}
		} else {
			fmt.Printf("✅ 找到用户: username=%s, status=%s, password=%s\n", username, status, password)
		}

		// 显示所有用户
		fmt.Println("\n所有用户:")
		fmt.Println("用户名 | 昵称 | 状态 | 归属组ID | 归属团队ID")
		fmt.Println("------|------|------|----------|----------")
		rows, err = db.Query("SELECT username, nickname, status, belong_group_id, belong_team_id FROM account ORDER BY username")
		if err != nil {
			log.Fatal("查询用户列表失败:", err)
		}
		defer rows.Close()

		for rows.Next() {
			var username, nickname, status string
			var belongGroupId, belongTeamId sql.NullInt32
			err := rows.Scan(&username, &nickname, &status, &belongGroupId, &belongTeamId)
			if err != nil {
				log.Fatal("读取用户数据失败:", err)
			}
			groupId := 0
			teamId := 0
			if belongGroupId.Valid {
				groupId = int(belongGroupId.Int32)
			}
			if belongTeamId.Valid {
				teamId = int(belongTeamId.Int32)
			}
			fmt.Printf("%s | %s | %s | %d | %d\n",
				username, nickname, status, groupId, teamId)
		}
	} else {
		fmt.Println("❌ account表不存在")
	}
}
