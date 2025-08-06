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
	BelongGroupId                int32  `gorm:"column:belong_group_id"`
	BelongGroupUsername          string `gorm:"column:belong_group_username"`
	BelongGroupNickname          string `gorm:"column:belong_group_nickname"`
	BelongGroupCreatedTimestamp  int64  `gorm:"column:belong_group_created_timestamp"`
	BelongGroupModifiedTimestamp int64  `gorm:"column:belong_group_modified_timestamp"`
	BelongGroupCurrentCnt        int32  `gorm:"column:belong_group_current_cnt"`

	// 所属团队信息
	BelongTeamId                int32  `gorm:"column:belong_team_id"`
	BelongTeamUsername          string `gorm:"column:belong_team_username"`
	BelongTeamNickname          string `gorm:"column:belong_team_nickname"`
	BelongTeamCreatedTimestamp  int64  `gorm:"column:belong_team_created_timestamp"`
	BelongTeamModifiedTimestamp int64  `gorm:"column:belong_team_modified_timestamp"`
	BelongTeamCurrentCnt        int32  `gorm:"column:belong_team_current_cnt"`

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

	// 设置表名
	db.Table("account")

	// 更新账户的组织信息
	updates := []struct {
		username      string
		groupID       int32
		groupUsername string
		groupNickname string
		teamID        int32
		teamUsername  string
		teamNickname  string
	}{
		{"admin", 1, "admin_group", "系统管理组", 1, "admin_team", "系统管理团队"},
		{"company_manager", 2, "nanjing_tianyuan", "南京-天元大厦组", 2, "marketing_team_a", "营销团队A"},
		{"group_manager", 3, "beijing_office", "北京办公室组", 3, "sales_team_b", "销售团队B"},
		{"team_manager", 4, "shanghai_branch", "上海分公司组", 4, "tech_team_c", "技术团队C"},
		{"group_manager2", 5, "guangzhou_office", "广州办公室组", 5, "service_team_d", "客服团队D"},
		{"employee001", 2, "nanjing_tianyuan", "南京-天元大厦组", 2, "marketing_team_a", "营销团队A"},
		{"employee002", 3, "beijing_office", "北京办公室组", 3, "sales_team_b", "销售团队B"},
		{"employee003", 4, "shanghai_branch", "上海分公司组", 4, "tech_team_c", "技术团队C"},
	}

	for _, update := range updates {
		result := db.Table("account").Where("username = ?", update.username).Updates(map[string]interface{}{
			"belong_group_id": update.groupID,
			"belong_team_id":  update.teamID,
		})

		if result.Error != nil {
			log.Printf("Failed to update %s: %v", update.username, result.Error)
		} else {
			fmt.Printf("Updated %s: %d rows affected\n", update.username, result.RowsAffected)
		}
	}

	// 验证更新结果
	var accounts []Account
	db.Table("account").Find(&accounts)

	fmt.Println("\nUpdated accounts:")
	for _, acc := range accounts {
		fmt.Printf("ID: %d, Username: %s, Group: %s, Team: %s\n",
			acc.Id, acc.Username, acc.BelongGroupNickname, acc.BelongTeamNickname)
	}
}
