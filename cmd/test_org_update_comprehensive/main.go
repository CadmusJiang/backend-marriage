package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

func main() {
	log.Println("开始全面测试组织更新功能...")

	// 测试用例1：不包含status字段的更新（应该成功）
	log.Println("\n=== 测试用例1：不包含status字段的更新 ===")
	testUpdateWithoutStatus()

	// 测试用例2：包含有效status字段的更新（应该成功）
	log.Println("\n=== 测试用例2：包含有效status字段的更新 ===")
	testUpdateWithValidStatus()

	// 测试用例3：包含无效status字段的更新（应该失败）
	log.Println("\n=== 测试用例3：包含无效status字段的更新 ===")
	testUpdateWithInvalidStatus()

	log.Println("\n所有测试完成！")
}

func testUpdateWithoutStatus() {
	testData := map[string]interface{}{
		"name":    "测试组名称-无状态",
		"version": 0,
		// 故意不包含status字段
	}

	jsonData, _ := json.Marshal(testData)
	resp := sendUpdateRequest(jsonData)

	if resp.StatusCode == 200 {
		log.Println("✅ 测试用例1通过：不包含status字段的更新成功")
	} else {
		log.Printf("❌ 测试用例1失败：状态码 %d", resp.StatusCode)
	}
}

func testUpdateWithValidStatus() {
	testData := map[string]interface{}{
		"name":    "测试组名称-有效状态",
		"status":  "enabled",
		"version": 0,
	}

	jsonData, _ := json.Marshal(testData)
	resp := sendUpdateRequest(jsonData)

	if resp.StatusCode == 200 {
		log.Println("✅ 测试用例2通过：包含有效status字段的更新成功")
	} else {
		log.Printf("❌ 测试用例2失败：状态码 %d", resp.StatusCode)
	}
}

func testUpdateWithInvalidStatus() {
	testData := map[string]interface{}{
		"name":    "测试组名称-无效状态",
		"status":  "invalid_status",
		"version": 0,
	}

	jsonData, _ := json.Marshal(testData)
	resp := sendUpdateRequest(jsonData)

	if resp.StatusCode == 400 || resp.StatusCode == 500 {
		log.Println("✅ 测试用例3通过：包含无效status字段的更新被正确拒绝")
	} else {
		log.Printf("❌ 测试用例3失败：状态码 %d，应该被拒绝", resp.StatusCode)
	}
}

func sendUpdateRequest(jsonData []byte) *http.Response {
	url := "http://localhost:8080/api/v1/groups/10001"
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("创建请求失败: %v", err)
		return nil
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer test_token")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("发送请求失败: %v", err)
		return nil
	}

	// 读取响应内容
	var responseBody map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&responseBody); err != nil {
		log.Printf("解析响应失败: %v", err)
	} else {
		log.Printf("响应内容: %+v", responseBody)
	}

	return resp
}
