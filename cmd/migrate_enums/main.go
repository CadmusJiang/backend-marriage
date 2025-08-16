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
	dbPass := getEnv("DB_PASSWORD", "")
	dbName := getEnv("DB_NAME", "marriage_system")

	// 构建数据库连接字符串
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbUser, dbPass, dbHost, dbPort, dbName)

	// 连接数据库
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// 测试连接
	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	log.Println("Successfully connected to database")

	// 开始迁移
	if err := migrateEnums(db); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	log.Println("Enum migration completed successfully")
}

func migrateEnums(db *sql.DB) error {
	// 1. 迁移 account.status
	log.Println("Migrating account.status...")
	if err := migrateAccountStatus(db); err != nil {
		return fmt.Errorf("failed to migrate account.status: %v", err)
	}

	// 2. 迁移 org.org_type
	log.Println("Migrating org.org_type...")
	if err := migrateOrgType(db); err != nil {
		return fmt.Errorf("failed to migrate org.org_type: %v", err)
	}

	// 3. 迁移 org.status
	log.Println("Migrating org.status...")
	if err := migrateOrgStatus(db); err != nil {
		return fmt.Errorf("failed to migrate org.status: %v", err)
	}

	// 4. 迁移 account_org_relation.relation_type
	log.Println("Migrating account_org_relation.relation_type...")
	if err := migrateRelationType(db); err != nil {
		return fmt.Errorf("failed to migrate account_org_relation.relation_type: %v", err)
	}

	// 5. 迁移 account_org_relation.status
	log.Println("Migrating account_org_relation.status...")
	if err := migrateRelationStatus(db); err != nil {
		return fmt.Errorf("failed to migrate account_org_relation.status: %v", err)
	}

	// 6. 迁移 customer_authorization_record 布尔字段
	log.Println("Migrating customer_authorization_record boolean fields...")
	if err := migrateCustomerBooleanFields(db); err != nil {
		return fmt.Errorf("failed to migrate customer_authorization_record boolean fields: %v", err)
	}

	// 7. 迁移 outbox_events.status
	log.Println("Migrating outbox_events.status...")
	if err := migrateOutboxStatus(db); err != nil {
		return fmt.Errorf("failed to migrate outbox_events.status: %v", err)
	}

	return nil
}

func migrateAccountStatus(db *sql.DB) error {
	// 先添加新列
	_, err := db.Exec("ALTER TABLE `account` ADD COLUMN `status_new` enum('enabled','disabled') NOT NULL DEFAULT 'enabled' COMMENT '状态 enabled:enabled disabled:disabled'")
	if err != nil {
		return err
	}

	// 更新数据
	_, err = db.Exec("UPDATE `account` SET `status_new` = CASE WHEN `status` = 1 THEN 'enabled' WHEN `status` = 2 THEN 'disabled' ELSE 'enabled' END")
	if err != nil {
		return err
	}

	// 删除旧列，重命名新列
	_, err = db.Exec("ALTER TABLE `account` DROP COLUMN `status`")
	if err != nil {
		return err
	}

	_, err = db.Exec("ALTER TABLE `account` CHANGE `status_new` `status` enum('enabled','disabled') NOT NULL DEFAULT 'enabled' COMMENT '状态 enabled:enabled disabled:disabled'")
	return err
}

func migrateOrgType(db *sql.DB) error {
	// 先添加新列
	_, err := db.Exec("ALTER TABLE `org` ADD COLUMN `org_type_new` enum('group','team') NOT NULL COMMENT '组织类型 group:group team:team'")
	if err != nil {
		return err
	}

	// 更新数据
	_, err = db.Exec("UPDATE `org` SET `org_type_new` = CASE WHEN `org_type` = 1 THEN 'group' WHEN `org_type` = 2 THEN 'team' ELSE 'group' END")
	if err != nil {
		return err
	}

	// 删除旧列，重命名新列
	_, err = db.Exec("ALTER TABLE `org` DROP COLUMN `org_type`")
	if err != nil {
		return err
	}

	_, err = db.Exec("ALTER TABLE `org` CHANGE `org_type_new` `org_type` enum('group','team') NOT NULL COMMENT '组织类型 group:group team:team'")
	return err
}

func migrateOrgStatus(db *sql.DB) error {
	// 先添加新列
	_, err := db.Exec("ALTER TABLE `org` ADD COLUMN `status_new` enum('enabled','disabled') NOT NULL DEFAULT 'enabled' COMMENT '状态 enabled:enabled disabled:disabled'")
	if err != nil {
		return err
	}

	// 更新数据
	_, err = db.Exec("UPDATE `org` SET `status_new` = CASE WHEN `status` = 1 THEN 'enabled' WHEN `status` = 2 THEN 'disabled' ELSE 'enabled' END")
	if err != nil {
		return err
	}

	// 删除旧列，重命名新列
	_, err = db.Exec("ALTER TABLE `org` DROP COLUMN `status`")
	if err != nil {
		return err
	}

	_, err = db.Exec("ALTER TABLE `org` CHANGE `status_new` `status` enum('enabled','disabled') NOT NULL DEFAULT 'enabled' COMMENT '状态 enabled:enabled disabled:disabled'")
	return err
}

func migrateRelationType(db *sql.DB) error {
	// 先添加新列
	_, err := db.Exec("ALTER TABLE `account_org_relation` ADD COLUMN `relation_type_new` enum('belong','manage') NOT NULL COMMENT '关联类型 belong:belong manage:manage'")
	if err != nil {
		return err
	}

	// 更新数据
	_, err = db.Exec("UPDATE `account_org_relation` SET `relation_type_new` = CASE WHEN `relation_type` = 1 THEN 'belong' WHEN `relation_type` = 2 THEN 'manage' ELSE 'belong' END")
	if err != nil {
		return err
	}

	// 删除旧列，重命名新列
	_, err = db.Exec("ALTER TABLE `account_org_relation` DROP COLUMN `relation_type`")
	if err != nil {
		return err
	}

	_, err = db.Exec("ALTER TABLE `account_org_relation` CHANGE `relation_type_new` `relation_type` enum('belong','manage') NOT NULL COMMENT '关联类型 belong:belong manage:manage'")
	return err
}

func migrateRelationStatus(db *sql.DB) error {
	// 先添加新列
	_, err := db.Exec("ALTER TABLE `account_org_relation` ADD COLUMN `status_new` enum('active','inactive') NOT NULL DEFAULT 'active' COMMENT '状态 active:active inactive:inactive'")
	if err != nil {
		return err
	}

	// 更新数据
	_, err = db.Exec("UPDATE `account_org_relation` SET `status_new` = CASE WHEN `status` = 1 THEN 'active' WHEN `status` = 2 THEN 'inactive' ELSE 'active' END")
	if err != nil {
		return err
	}

	// 删除旧列，重命名新列
	_, err = db.Exec("ALTER TABLE `account_org_relation` DROP COLUMN `status`")
	if err != nil {
		return err
	}

	_, err = db.Exec("ALTER TABLE `account_org_relation` CHANGE `status_new` `status` enum('active','inactive') NOT NULL DEFAULT 'active' COMMENT '状态 active:active inactive:inactive'")
	return err
}

func migrateCustomerBooleanFields(db *sql.DB) error {
	// 迁移 authorization_status
	_, err := db.Exec("ALTER TABLE `customer_authorization_record` MODIFY COLUMN `authorization_status` enum('authorized','unauthorized') NOT NULL DEFAULT 'unauthorized' COMMENT '授权状态'")
	if err != nil {
		return err
	}

	// 迁移 completion_status
	_, err = db.Exec("ALTER TABLE `customer_authorization_record` MODIFY COLUMN `completion_status` enum('complete','incomplete') NOT NULL DEFAULT 'incomplete' COMMENT '完善状态'")
	if err != nil {
		return err
	}

	// 迁移 assignment_status
	_, err = db.Exec("ALTER TABLE `customer_authorization_record` MODIFY COLUMN `assignment_status` enum('assigned','unassigned') NOT NULL DEFAULT 'unassigned' COMMENT '分配状态'")
	if err != nil {
		return err
	}

	// 迁移 payment_status
	_, err = db.Exec("ALTER TABLE `customer_authorization_record` MODIFY COLUMN `payment_status` enum('paid','unpaid') NOT NULL DEFAULT 'unpaid' COMMENT '付费状态'")
	if err != nil {
		return err
	}

	// 更新数据：将 1 改为对应的状态值，0 改为对应的状态值
	_, err = db.Exec("UPDATE `customer_authorization_record` SET `authorization_status` = CASE WHEN `authorization_status` = 1 THEN 'authorized' ELSE 'unauthorized' END")
	if err != nil {
		return err
	}

	_, err = db.Exec("UPDATE `customer_authorization_record` SET `completion_status` = CASE WHEN `completion_status` = 1 THEN 'complete' ELSE 'incomplete' END")
	if err != nil {
		return err
	}

	_, err = db.Exec("UPDATE `customer_authorization_record` SET `assignment_status` = CASE WHEN `assignment_status` = 1 THEN 'assigned' ELSE 'unassigned' END")
	if err != nil {
		return err
	}

	_, err = db.Exec("UPDATE `customer_authorization_record` SET `payment_status` = CASE WHEN `payment_status` = 1 THEN 'paid' ELSE 'unpaid' END")
	return err
}

func migrateOutboxStatus(db *sql.DB) error {
	// 先添加新列
	_, err := db.Exec("ALTER TABLE `outbox_events` ADD COLUMN `status_new` enum('unpublished','published') NOT NULL DEFAULT 'unpublished' COMMENT '发布状态 unpublished:未发布 published:已发布'")
	if err != nil {
		return err
	}

	// 更新数据
	_, err = db.Exec("UPDATE `outbox_events` SET `status_new` = CASE WHEN `status` = 0 THEN 'unpublished' WHEN `status` = 1 THEN 'published' ELSE 'unpublished' END")
	if err != nil {
		return err
	}

	// 删除旧列，重命名新列
	_, err = db.Exec("ALTER TABLE `outbox_events` DROP COLUMN `status`")
	if err != nil {
		return err
	}

	_, err = db.Exec("ALTER TABLE `outbox_events` CHANGE `status_new` `status` enum('unpublished','published') NOT NULL DEFAULT 'unpublished' COMMENT '发布状态 unpublished:未发布 published:已发布'")
	return err
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
