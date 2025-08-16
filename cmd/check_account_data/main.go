package main

import (
	"fmt"
	"log"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Account struct {
	Id       uint32 `gorm:"primaryKey;column:id"`
	Username string `gorm:"column:username"`
	Name     string `gorm:"column:name"`
	Password string `gorm:"column:password"`
	Phone    string `gorm:"column:phone"`
	RoleType string `gorm:"column:role_type"`
	Status   string `gorm:"column:status"`

	// 所属组信息
	BelongGroupId uint32 `gorm:"column:belong_group_id"`
	BelongTeamId  uint32 `gorm:"column:belong_team_id"`

	LastLoginAt *time.Time `gorm:"column:last_login_at"`
	CreatedAt   time.Time  `gorm:"column:created_at"`
	UpdatedAt   time.Time  `gorm:"column:updated_at"`
	CreatedUser string     `gorm:"column:created_user"`
	UpdatedUser string     `gorm:"column:updated_user"`
}

func main() {
	// 数据库连接
	dsn := "marriage_system:19970901Zyt@tcp(gz-cdb-pepynap7.sql.tencentcdb.com:63623)/marriage_system?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// 更新role_type: company_manager -> company_manager
	fmt.Println("正在更新role_type...")
	result := db.Exec("UPDATE account SET role_type = ? WHERE role_type = ?", "company_manager", "company_manager")
	if result.Error != nil {
		log.Fatal("更新role_type失败:", result.Error)
	}
	fmt.Printf("更新完成，影响行数: %d\n", result.RowsAffected)

	// 查询所有账户数据
	var accounts []Account
	db.Table("account").Find(&accounts)

	fmt.Println("\nDatabase account data:")
	fmt.Println("==================================================================================")
	for _, acc := range accounts {
		fmt.Printf("ID: %d, Username: %s, Name: %s, RoleType: %s, Status: %s, BelongGroupId: %d, BelongTeamId: %d\n",
			acc.Id, acc.Username, acc.Name, acc.RoleType, acc.Status, acc.BelongGroupId, acc.BelongTeamId)
	}

	// 特别显示admin用户的信息
	fmt.Println("\n==================================================================================")
	fmt.Println("Admin user details:")
	var adminAccount Account
	db.Table("account").Where("username = ?", "admin").First(&adminAccount)
	fmt.Printf("ID: %d\n", adminAccount.Id)
	fmt.Printf("Username: %s\n", adminAccount.Username)
	fmt.Printf("Name: %s\n", adminAccount.Name)
	fmt.Printf("Phone: %s\n", adminAccount.Phone)
	fmt.Printf("RoleType: %s\n", adminAccount.RoleType)
	fmt.Printf("Status: %s\n", adminAccount.Status)
	fmt.Printf("BelongGroupId: %d\n", adminAccount.BelongGroupId)
	fmt.Printf("BelongTeamId: %d\n", adminAccount.BelongTeamId)
	fmt.Printf("CreatedAt: %s\n", adminAccount.CreatedAt.Format("2006-01-02 15:04:05"))
	fmt.Printf("UpdatedAt: %s\n", adminAccount.UpdatedAt.Format("2006-01-02 15:04:05"))

	// 统计各角色类型数量
	fmt.Println("\n==================================================================================")
	fmt.Println("角色类型统计:")
	var roleCounts []struct {
		RoleType string
		Count    int64
	}
	db.Table("account").Select("role_type, count(*) as count").Group("role_type").Find(&roleCounts)

	for _, rc := range roleCounts {
		fmt.Printf("%s: %d个\n", rc.RoleType, rc.Count)
	}
}
