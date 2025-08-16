package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// 数据库连接信息
	dsn := "marriage_system:19970901Zyt@tcp(gz-cdb-pepynap7.sql.tencentcdb.com:63623)/marriage_system?charset=utf8mb4&parseTime=True&loc=Local"

	// 连接数据库
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("连接数据库失败: %v", err)
	}
	defer db.Close()

	// 测试连接
	if err := db.Ping(); err != nil {
		log.Fatalf("数据库连接测试失败: %v", err)
	}

	log.Println("数据库连接成功，开始检查表结构...")

	// 检查org表结构 - 使用SHOW COLUMNS命令
	log.Println("=== 检查 org 表结构 (SHOW COLUMNS) ===")
	rows, err := db.Query("SHOW COLUMNS FROM org")
	if err != nil {
		log.Fatalf("查询表结构失败: %v", err)
	}
	defer rows.Close()

	// 获取列信息
	columns, err := rows.Columns()
	if err != nil {
		log.Fatalf("获取列信息失败: %v", err)
	}
	log.Printf("SHOW COLUMNS 返回的列数: %d", len(columns))
	for i, col := range columns {
		log.Printf("列 %d: %s", i, col)
	}

	fmt.Printf("%-20s %-20s %-10s %-10s %-10s %-10s\n", "Field", "Type", "Null", "Key", "Default", "Extra")
	fmt.Println(string(make([]byte, 80, 80)))

	for rows.Next() {
		// 根据实际列数创建变量
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))
		for i := range values {
			valuePtrs[i] = &values[i]
		}

		if err := rows.Scan(valuePtrs...); err != nil {
			log.Printf("扫描行数据失败: %v", err)
			continue
		}

		// 打印前6列（如果存在）
		field := ""
		typ := ""
		null := ""
		key := ""
		defaultVal := ""
		extra := ""

		if len(values) > 0 && values[0] != nil {
			field = fmt.Sprintf("%v", values[0])
		}
		if len(values) > 1 && values[1] != nil {
			typ = fmt.Sprintf("%v", values[1])
		}
		if len(values) > 2 && values[2] != nil {
			null = fmt.Sprintf("%v", values[2])
		}
		if len(values) > 3 && values[3] != nil {
			key = fmt.Sprintf("%v", values[3])
		}
		if len(values) > 4 && values[4] != nil {
			defaultVal = fmt.Sprintf("%v", values[4])
		}
		if len(values) > 5 && values[5] != nil {
			extra = fmt.Sprintf("%v", values[5])
		}

		fmt.Printf("%-20s %-20s %-10s %-10s %-10s %-10s\n", field, typ, null, key, defaultVal, extra)
	}

	// 检查org表中的实际数据
	log.Println("\n=== 检查 org 表中的数据 ===")
	dataRows, err := db.Query("SELECT id, org_type, status FROM org LIMIT 5")
	if err != nil {
		log.Fatalf("查询表数据失败: %v", err)
	}
	defer dataRows.Close()

	fmt.Printf("%-10s %-15s %-15s\n", "ID", "Org Type", "Status")
	fmt.Println(string(make([]byte, 50, 50)))

	for dataRows.Next() {
		var id uint32
		var orgType, status string
		if err := dataRows.Scan(&id, &orgType, &status); err != nil {
			log.Printf("扫描数据行失败: %v", err)
			continue
		}
		fmt.Printf("%-10d %-15s %-15s\n", id, orgType, status)
	}

	// 检查information_schema中的列信息
	log.Println("\n=== 检查 information_schema 中的列信息 ===")
	infoRows, err := db.Query(`
		SELECT COLUMN_NAME, DATA_TYPE, IS_NULLABLE, COLUMN_DEFAULT, COLUMN_TYPE, COLUMN_COMMENT
		FROM INFORMATION_SCHEMA.COLUMNS 
		WHERE TABLE_SCHEMA = 'marriage_system' AND TABLE_NAME = 'org' AND COLUMN_NAME = 'status'
	`)
	if err != nil {
		log.Fatalf("查询information_schema失败: %v", err)
	}
	defer infoRows.Close()

	fmt.Printf("%-20s %-20s %-15s %-20s %-30s %-30s\n", "Column", "Data Type", "Nullable", "Default", "Column Type", "Comment")
	fmt.Println(string(make([]byte, 150, 150)))

	for infoRows.Next() {
		var columnName, dataType, isNullable, columnDefault, columnType, columnComment sql.NullString
		if err := infoRows.Scan(&columnName, &dataType, &isNullable, &columnDefault, &columnType, &columnComment); err != nil {
			log.Printf("扫描information_schema行失败: %v", err)
			continue
		}
		fmt.Printf("%-20s %-20s %-15s %-20s %-30s %-30s\n",
			columnName.String, dataType.String, isNullable.String, columnDefault.String, columnType.String, columnComment.String)
	}

	// 尝试更新一条记录来测试
	log.Println("\n=== 测试更新操作 ===")
	result, err := db.Exec("UPDATE org SET status = ? WHERE id = ?", "enabled", 10001)
	if err != nil {
		log.Printf("更新测试失败: %v", err)
	} else {
		rowsAffected, _ := result.RowsAffected()
		log.Printf("更新测试成功，影响行数: %d", rowsAffected)
	}

	log.Println("检查完成！")
}
