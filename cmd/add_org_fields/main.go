package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// 从环境变量获取数据库连接信息
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "3306")
	dbUser := getEnv("DB_USER", "root")
	dbPassword := getEnv("DB_PASSWORD", "")
	dbName := getEnv("DB_NAME", "marriage_system")

	// 构建数据库连接字符串
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbUser, dbPassword, dbHost, dbPort, dbName)

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

	log.Println("数据库连接成功，开始添加字段...")

	// 添加 address 字段
	log.Println("正在添加 address 字段...")
	_, err = db.Exec("ALTER TABLE `org` ADD COLUMN `address` varchar(255) DEFAULT NULL COMMENT '地址'")
	if err != nil {
		log.Printf("添加 address 字段失败: %v", err)
	} else {
		log.Println("address 字段添加成功")
	}

	// 添加 extra 字段
	log.Println("正在添加 extra 字段...")
	_, err = db.Exec("ALTER TABLE `org` ADD COLUMN `extra` json DEFAULT NULL COMMENT '扩展信息'")
	if err != nil {
		log.Printf("添加 extra 字段失败: %v", err)
	} else {
		log.Println("extra 字段添加成功")
	}

	log.Println("字段添加完成！")
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
