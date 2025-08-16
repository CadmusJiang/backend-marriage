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

	// 直接尝试修复status字段
	log.Println("正在修复status字段...")

	// 方法1：尝试直接修改字段类型
	_, err = db.Exec("ALTER TABLE `org` MODIFY COLUMN `status` enum('enabled','disabled') NOT NULL DEFAULT 'enabled' COMMENT '状态 enabled:enabled disabled:disabled'")
	if err != nil {
		log.Printf("直接修改status字段失败: %v", err)
		log.Println("尝试使用替换方法...")

		// 方法2：添加新字段，更新数据，删除旧字段
		log.Println("正在添加新的status字段...")
		_, err = db.Exec("ALTER TABLE `org` ADD COLUMN `status_new` enum('enabled','disabled') NOT NULL DEFAULT 'enabled' COMMENT '状态 enabled:enabled disabled:disabled'")
		if err != nil {
			log.Printf("添加新status字段失败: %v", err)
		} else {
			log.Println("新status字段添加成功")

			// 设置所有记录为enabled
			log.Println("正在设置默认状态...")
			_, err = db.Exec("UPDATE `org` SET `status_new` = 'enabled'")
			if err != nil {
				log.Printf("设置默认状态失败: %v", err)
			} else {
				log.Println("默认状态设置成功")

				// 删除旧字段
				log.Println("正在删除旧status字段...")
				_, err = db.Exec("ALTER TABLE `org` DROP COLUMN `status`")
				if err != nil {
					log.Printf("删除旧status字段失败: %v", err)
				} else {
					log.Println("旧status字段删除成功")

					// 重命名新字段
					log.Println("正在重命名新status字段...")
					_, err = db.Exec("ALTER TABLE `org` CHANGE `status_new` `status` enum('enabled','disabled') NOT NULL DEFAULT 'enabled' COMMENT '状态 enabled:enabled disabled:disabled'")
					if err != nil {
						log.Printf("重命名status字段失败: %v", err)
					} else {
						log.Println("status字段重命名成功")
					}
				}
			}
		}
	} else {
		log.Println("status字段修复成功！")
	}

	log.Println("修复过程完成！")
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
