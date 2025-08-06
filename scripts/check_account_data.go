package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	dsn := "marriage_system:19970901Zyt@tcp(gz-cdb-pepynap7.sql.tencentcdb.com:63623)/marriage_system?charset=utf8mb4&parseTime=True&loc=Local"
	fmt.Println("正在连接腾讯云数据库...")
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("连接数据库失败:", err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal("Ping数据库失败:", err)
	}
	fmt.Println("✅ 数据库连接成功！")

	// 检查account表是否存在
	var tableExists int
	err = db.QueryRow("SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = 'marriage_system' AND table_name = 'account'").Scan(&tableExists)
	if err != nil {
		log.Fatal("查询表失败:", err)
	}

	if tableExists > 0 {
		fmt.Println("✅ account表存在")
		var count int
		err = db.QueryRow("SELECT COUNT(*) FROM account").Scan(&count)
		if err != nil {
			log.Fatal("查询数据失败:", err)
		}
		fmt.Printf("📊 account表中有 %d 条记录\n", count)

		// 查询所有账户数据
		rows, err := db.Query("SELECT id, username, nickname, password, phone, role_type, status FROM account ORDER BY id")
		if err != nil {
			log.Fatal("查询数据失败:", err)
		}
		defer rows.Close()

		fmt.Println("\n现有账户数据:")
		fmt.Println("ID | 用户名 | 姓名 | 密码 | 手机号 | 角色类型 | 状态")
		fmt.Println("---|--------|------|------|--------|----------|------")
		for rows.Next() {
			var id int
			var username, nickname, password, phone, roleType, status string
			err := rows.Scan(&id, &username, &nickname, &password, &phone, &roleType, &status)
			if err != nil {
				log.Fatal("读取数据失败:", err)
			}
			fmt.Printf("%d | %s | %s | %s | %s | %s | %s\n", id, username, nickname, password, phone, roleType, status)
		}

		// 特别检查company_manager账户
		var companyManagerExists int
		err = db.QueryRow("SELECT COUNT(*) FROM account WHERE username = 'company_manager'").Scan(&companyManagerExists)
		if err != nil {
			log.Fatal("查询company_manager失败:", err)
		}

		if companyManagerExists > 0 {
			fmt.Println("\n✅ company_manager账户存在")
			var password string
			err = db.QueryRow("SELECT password FROM account WHERE username = 'company_manager'").Scan(&password)
			if err != nil {
				log.Fatal("查询密码失败:", err)
			}
			fmt.Printf("company_manager的密码是: %s\n", password)
		} else {
			fmt.Println("\n❌ company_manager账户不存在")
		}
	} else {
		fmt.Println("❌ account表不存在")
	}
}
