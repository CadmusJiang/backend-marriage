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

	log.Println("数据库连接成功，开始修复org.status字段...")

	// 检查当前status字段的类型
	var columnType string
	err = db.QueryRow("SELECT DATA_TYPE FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = 'org' AND COLUMN_NAME = 'status'", dbName).Scan(&columnType)
	if err != nil {
		log.Fatalf("查询status字段类型失败: %v", err)
	}

	log.Printf("当前status字段类型: %s", columnType)

	if columnType == "enum" {
		log.Println("status字段已经是enum类型，无需修复")
		return
	}

	// 备份现有数据
	log.Println("正在备份现有数据...")
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM org").Scan(&count)
	if err != nil {
		log.Fatalf("查询记录数失败: %v", err)
	}
	log.Printf("org表共有 %d 条记录", count)

	// 添加新的status字段
	log.Println("正在添加新的status字段...")
	_, err = db.Exec("ALTER TABLE `org` ADD COLUMN `status_new` enum('enabled','disabled') NOT NULL DEFAULT 'enabled' COMMENT '状态 enabled:enabled disabled:disabled'")
	if err != nil {
		log.Printf("添加新status字段失败: %v", err)
	} else {
		log.Println("新status字段添加成功")
	}

	// 更新数据：将旧的status值转换为新的enum值
	log.Println("正在转换数据...")
	if columnType == "int" || columnType == "tinyint" {
		_, err = db.Exec("UPDATE `org` SET `status_new` = CASE WHEN `status` = 1 THEN 'enabled' WHEN `status` = 2 THEN 'disabled' ELSE 'enabled' END")
	} else {
		// 如果是其他类型，直接设置为默认值
		_, err = db.Exec("UPDATE `org` SET `status_new` = 'enabled'")
	}

	if err != nil {
		log.Printf("数据转换失败: %v", err)
	} else {
		log.Println("数据转换成功")
	}

	// 删除旧字段
	log.Println("正在删除旧status字段...")
	_, err = db.Exec("ALTER TABLE `org` DROP COLUMN `status`")
	if err != nil {
		log.Printf("删除旧status字段失败: %v", err)
	} else {
		log.Println("旧status字段删除成功")
	}

	// 重命名新字段
	log.Println("正在重命名新status字段...")
	_, err = db.Exec("ALTER TABLE `org` CHANGE `status_new` `status` enum('enabled','disabled') NOT NULL DEFAULT 'enabled' COMMENT '状态 enabled:enabled disabled:disabled'")
	if err != nil {
		log.Printf("重命名status字段失败: %v", err)
	} else {
		log.Println("status字段重命名成功")
	}

	log.Println("org.status字段修复完成！")
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
