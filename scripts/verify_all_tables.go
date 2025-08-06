package main

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	fmt.Println("=== 验证所有表结构 ===")

	// 连接数据库
	dsn := "marriage_system:19970901Zyt@tcp(gz-cdb-pepynap7.sql.tencentcdb.com:63623)/marriage_system?charset=utf8mb4&parseTime=true&loc=Local"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		fmt.Printf("❌ 连接数据库失败: %v\n", err)
		return
	}
	defer db.Close()

	// 查询所有表
	rows, err := db.Query("SHOW TABLES")
	if err != nil {
		fmt.Printf("❌ 查询表失败: %v\n", err)
		return
	}
	defer rows.Close()

	var tables []string
	for rows.Next() {
		var tableName string
		if err := rows.Scan(&tableName); err != nil {
			fmt.Printf("❌ 读取表名失败: %v\n", err)
			continue
		}
		tables = append(tables, tableName)
	}

	fmt.Printf("📋 数据库中共有 %d 个表:\n", len(tables))
	for _, tableName := range tables {
		fmt.Printf("  - %s\n", tableName)
	}

	// 检查每个表的结构
	fmt.Println("\n🔍 检查各表结构:")
	for _, tableName := range tables {
		fmt.Printf("\n📋 %s表结构:\n", tableName)

		rows2, err := db.Query(fmt.Sprintf("DESCRIBE %s", tableName))
		if err != nil {
			fmt.Printf("❌ 查询%s表结构失败: %v\n", tableName, err)
			continue
		}

		var hasIsDeleted, hasCreatedAt, hasUpdatedAt bool
		for rows2.Next() {
			var field, typ, null, key, defaultVal, extra sql.NullString
			err := rows2.Scan(&field, &typ, &null, &key, &defaultVal, &extra)
			if err != nil {
				continue
			}

			fieldName := field.String
			fmt.Printf("  - %s: %s\n", fieldName, typ.String)

			if fieldName == "is_deleted" {
				hasIsDeleted = true
			}
			if fieldName == "created_at" {
				hasCreatedAt = true
			}
			if fieldName == "updated_at" {
				hasUpdatedAt = true
			}
		}
		rows2.Close()

		// 检查是否还有需要移除的字段
		if hasIsDeleted || hasCreatedAt || hasUpdatedAt {
			fmt.Printf("  ⚠️  仍有需要移除的字段:")
			if hasIsDeleted {
				fmt.Printf(" is_deleted")
			}
			if hasCreatedAt {
				fmt.Printf(" created_at")
			}
			if hasUpdatedAt {
				fmt.Printf(" updated_at")
			}
			fmt.Printf("\n")
		} else {
			fmt.Printf("  ✅ 表结构正确，无需要移除的字段\n")
		}

		// 检查数据量
		var count int
		err = db.QueryRow(fmt.Sprintf("SELECT COUNT(*) FROM %s", tableName)).Scan(&count)
		if err != nil {
			fmt.Printf("❌ 查询%s表数据失败: %v\n", tableName, err)
			continue
		}
		fmt.Printf("  📊 数据量: %d 条\n", count)
	}

	fmt.Println("\n=== 验证完成 ===")
}
