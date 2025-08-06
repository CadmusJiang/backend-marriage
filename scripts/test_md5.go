package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
)

func main() {
	password := "123456"
	m := md5.New()
	m.Write([]byte(password))
	hashedPassword := hex.EncodeToString(m.Sum(nil))

	fmt.Printf("密码: %s\n", password)
	fmt.Printf("MD5值: %s\n", hashedPassword)

	// 验证是否与数据库中的值匹配
	expectedHash := "e10adc3949ba59abbe56e057f20f883e"
	if hashedPassword == expectedHash {
		fmt.Println("✅ MD5值匹配！")
	} else {
		fmt.Println("❌ MD5值不匹配！")
		fmt.Printf("期望: %s\n", expectedHash)
		fmt.Printf("实际: %s\n", hashedPassword)
	}
}
