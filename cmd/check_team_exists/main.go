package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

type Config struct {
	MySQL struct {
		Read struct {
			Addr string `toml:"addr"`
			Name string `toml:"name"`
			Pass string `toml:"pass"`
			User string `toml:"user"`
		} `toml:"read"`
	} `toml:"mysql"`
}

func main() {
	// 直接使用配置信息
	dsn := "marriage_system:19970901Zyt@tcp(gz-cdb-pepynap7.sql.tencentcdb.com:63623)/marriage_system?charset=utf8mb4&parseTime=True&loc=Local"

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("连接数据库失败:", err)
	}
	defer db.Close()

	// 测试连接
	if err := db.Ping(); err != nil {
		log.Fatal("数据库连接测试失败:", err)
	}

	fmt.Println("数据库连接成功！")

	// 检查团队ID 100001 是否存在
	teamID := 100001
	var team struct {
		ID      uint32 `db:"id"`
		OrgType string `db:"org_type"`
		Status  string `db:"status"`
		Name    string `db:"name"`
	}

	query := `SELECT id, org_type, status, name 
			  FROM org 
			  WHERE id = ? AND org_type = 'team'`

	err = db.QueryRow(query, teamID).Scan(
		&team.ID,
		&team.OrgType,
		&team.Status,
		&team.Name,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Printf("❌ 团队ID %d 不存在\n", teamID)
		} else {
			fmt.Printf("❌ 查询团队失败: %v\n", err)
		}
	} else {
		fmt.Printf("✅ 找到团队:\n")
		fmt.Printf("   ID: %d\n", team.ID)
		fmt.Printf("   类型: %s\n", team.OrgType)
		fmt.Printf("   状态: %s\n", team.Status)
		fmt.Printf("   名称: %s\n", team.Name)
	}

	// 检查所有团队
	fmt.Println("\n=== 所有团队列表 ===")
	rows, err := db.Query(`SELECT id, org_type, status, name 
						   FROM org 
						   WHERE org_type = 'team' 
						   ORDER BY id`)
	if err != nil {
		log.Fatal("查询所有团队失败:", err)
	}
	defer rows.Close()

	fmt.Printf("%-8s %-10s %-8s %-20s\n", "ID", "类型", "状态", "名称")
	fmt.Println("------------------------------------------------")

	for rows.Next() {
		var t struct {
			ID      uint32 `db:"id"`
			OrgType string `db:"org_type"`
			Status  string `db:"status"`
			Name    string `db:"name"`
		}

		err := rows.Scan(&t.ID, &t.OrgType, &t.Status, &t.Name)
		if err != nil {
			fmt.Printf("读取团队数据失败: %v\n", err)
			continue
		}

		fmt.Printf("%-8d %-10s %-8s %-20s\n",
			t.ID, t.OrgType, t.Status, t.Name)
	}

	if err = rows.Err(); err != nil {
		log.Fatal("遍历团队数据失败:", err)
	}
}
","trace_id":"x222222","method":"PUT","path":"/api/v1/teams/100001/members/7","remote_addr":"::1","trace_id":"x222222","method":"PUT","path":"/api/v1/teams/100001/members/7","query":"","user_agent":"PostmanRuntime/7.45.0"}
{"level":"info","time":"2025-08-16 20:36:53","caller":"interceptor/trace_middleware.go:68","msg":"HTTP请求结束","domain":"marriage_system[fat]","trace_id":"x222222","method":"PUT","path":"/api/v1/teams/100001/members/7","remote_addr":"::1","trace_id":"x222222","status_code":404,"response_size_kb":-0.0009765625}
{"level":"info","time":"2025-08-16 20:36:57","caller":"interceptor/trace_mid