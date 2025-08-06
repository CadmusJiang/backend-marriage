package main

import (
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Account struct {
	Id       int32  `gorm:"primaryKey;column:id"`
	Username string `gorm:"column:username"`
	Nickname string `gorm:"column:nickname"`
	Password string `gorm:"column:password"`
	Phone    string `gorm:"column:phone"`
	RoleType string `gorm:"column:role_type"`
	Status   string `gorm:"column:status"`

	// 所属组信息
	BelongGroupId int32 `gorm:"column:belong_group_id"`
	BelongTeamId  int32 `gorm:"column:belong_team_id"`

	LastLoginTimestamp int64  `gorm:"column:last_login_timestamp"`
	CreatedTimestamp   int64  `gorm:"column:created_timestamp"`
	ModifiedTimestamp  int64  `gorm:"column:modified_timestamp"`
	CreatedUser        string `gorm:"column:created_user"`
	UpdatedUser        string `gorm:"column:updated_user"`
}

func main() {
	// 数据库连接
	dsn := "marriage_system:19970901Zyt@tcp(gz-cdb-pepynap7.sql.tencentcdb.com:63623)/marriage_system?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// 查询所有账户数据
	var accounts []Account
	db.Table("account").Find(&accounts)

	fmt.Println("Database account data:")
	for _, acc := range accounts {
		fmt.Printf("ID: %d, Username: %s, BelongGroupId: %d, BelongTeamId: %d\n",
			acc.Id, acc.Username, acc.BelongGroupId, acc.BelongTeamId)
	}
}
