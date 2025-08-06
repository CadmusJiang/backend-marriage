package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// 腾讯云数据库连接信息
	dsn := "marriage_system:19970901Zyt@tcp(gz-cdb-pepynap7.sql.tencentcdb.com:63623)/marriage_system?charset=utf8mb4&parseTime=True&loc=Local"

	fmt.Println("正在连接腾讯云数据库...")

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

	// 清空现有数据
	fmt.Println("清空现有数据...")
	_, err = db.Exec("DELETE FROM org")
	if err != nil {
		log.Fatal("清空数据失败:", err)
	}

	// 重置自增ID
	_, err = db.Exec("ALTER TABLE org AUTO_INCREMENT = 1")
	if err != nil {
		log.Fatal("重置自增ID失败:", err)
	}

	fmt.Println("✅ 数据已清空")

	// 插入正确的组织结构
	fmt.Println("插入正确的组织结构...")

	now := time.Now().Unix()

	// 准备插入语句 - 使用正确的字段名
	insertSQL := `INSERT INTO org (
		org_name, org_type, org_path, org_level, 
		current_cnt, max_cnt, status, ext_data, 
		created_at, modified_at, created_user, updated_user
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	// 插入数据
	data := []struct {
		name, path, extData                   string
		orgType, orgLevel, currentCnt, maxCnt int
	}{
		// 组 (org_type = 1)
		{"系统管理组", "/1", `{"prefix": "weilan_"}`, 1, 1, 1, 10},
		{"南京-天元大厦组", "/2", `{"prefix": "weilan_"}`, 1, 1, 15, 50},
		{"南京-南京南站组", "/3", `{"prefix": "weilan_"}`, 1, 1, 12, 50},

		// 团队 (org_type = 2) - 每个团队只属于一个组
		{"系统管理团队", "/1/1", `{"team_leader": "系统管理员", "project": "系统维护"}`, 2, 2, 1, 10},
		{"营销团队A", "/2/1", `{"team_leader": "张伟", "project": "产品推广"}`, 2, 2, 8, 20},
		{"营销团队B", "/3/1", `{"team_leader": "刘强", "project": "市场拓展"}`, 2, 2, 10, 20},
		{"营销团队C", "/2/2", `{"team_leader": "王芳", "project": "品牌建设"}`, 2, 2, 6, 20},
	}

	for _, item := range data {
		_, err := db.Exec(insertSQL,
			item.name, item.orgType, item.path, item.orgLevel,
			item.currentCnt, item.maxCnt, 1, item.extData,
			now, now, "系统管理员", "系统管理员")

		if err != nil {
			log.Fatalf("插入数据失败 [%s]: %v", item.name, err)
		}

		fmt.Printf("✅ 插入: %s\n", item.name)
	}

	fmt.Println("\n✅ 所有数据插入完成！")

	// 验证数据
	fmt.Println("\n验证数据:")
	rows, err := db.Query("SELECT id, org_name, org_type, org_path FROM org ORDER BY org_type, id")
	if err != nil {
		log.Fatal("查询数据失败:", err)
	}
	defer rows.Close()

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

	fmt.Println("\n🎉 数据修复完成！每个team现在只属于一个组。")
}
