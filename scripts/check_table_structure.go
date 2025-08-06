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

	// 检查表结构
	fmt.Println("\n检查org表结构:")
	rows, err := db.Query("DESCRIBE org")
	if err != nil {
		log.Fatal("查询表结构失败:", err)
	}
	defer rows.Close()

	fmt.Println("字段名 | 类型 | 空 | 键 | 默认值 | 额外")
	fmt.Println("-------|------|----|----|--------|------")

	for rows.Next() {
		var field, typ, null, key, defaultVal, extra sql.NullString

		err := rows.Scan(&field, &typ, &null, &key, &defaultVal, &extra)
		if err != nil {
			log.Fatal("读取表结构失败:", err)
		}

		fmt.Printf("%s | %s | %s | %s | %s | %s\n",
			field.String, typ.String, null.String, key.String,
			defaultVal.String, extra.String)
	}
}
