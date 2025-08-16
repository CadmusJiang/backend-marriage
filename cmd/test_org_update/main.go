package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

func main() {
	// 测试更新组织状态
	log.Println("开始测试组织更新...")

	// 测试数据
	testData := map[string]interface{}{
		"name":    "测试组名称",
		"status":  "enabled",
		"version": 0,
	}

	// 序列化测试数据
	jsonData, err := json.Marshal(testData)
	if err != nil {
		log.Fatalf("序列化测试数据失败: %v", err)
	}

	// 发送PUT请求更新组织
	url := "http://localhost:8080/api/v1/groups/10001"
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatalf("创建请求失败: %v", err)
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer test_token") // 这里需要真实的token

	// 发送请求
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("发送请求失败: %v", err)
	}
	defer resp.Body.Close()

	log.Printf("响应状态码: %d", resp.StatusCode)

	// 读取响应内容
	var responseBody map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&responseBody); err != nil {
		log.Printf("解析响应失败: %v", err)
	} else {
		log.Printf("响应内容: %+v", responseBody)
	}

	if resp.StatusCode == 200 {
		log.Println("✅ 组织更新测试成功！")
	} else {
		log.Printf("❌ 组织更新测试失败，状态码: %d", resp.StatusCode)
	}
}
