package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

func main() {
	log.Println("开始测试更新后的API接口...")

	// 测试更新组接口
	log.Println("\n=== 测试更新组接口 ===")
	testUpdateGroup()

	// 测试更新团队接口
	log.Println("\n=== 测试更新团队接口 ===")
	testUpdateTeam()

	log.Println("\n所有测试完成！")
}

func testUpdateGroup() {
	testData := map[string]interface{}{
		"name":    "测试组名称-更新后",
		"version": 0,
		// 不包含status字段，测试我们的修复
	}

	jsonData, _ := json.Marshal(testData)
	resp := sendUpdateGroupRequest(jsonData)

	if resp != nil && resp.StatusCode == 200 {
		log.Println("✅ 更新组接口测试成功")
	} else {
		log.Printf("❌ 更新组接口测试失败，状态码: %d", resp.StatusCode)
	}
}

func testUpdateTeam() {
	testData := map[string]interface{}{
		"name":    "测试团队名称-更新后",
		"version": 0,
		// 不包含status字段，测试我们的修复
	}

	jsonData, _ := json.Marshal(testData)
	resp := sendUpdateTeamRequest(jsonData)

	if resp != nil && resp.StatusCode == 200 {
		log.Println("✅ 更新团队接口测试成功")
	} else {
		log.Printf("❌ 更新团队接口测试失败，状态码: %d", resp.StatusCode)
	}
}

func sendUpdateGroupRequest(jsonData []byte) *http.Response {
	url := "http://localhost:8080/api/v1/groups/10001"
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("创建更新组请求失败: %v", err)
		return nil
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer test_token")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("发送更新组请求失败: %v", err)
		return nil
	}

	// 读取响应内容
	var responseBody map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&responseBody); err != nil {
		log.Printf("解析更新组响应失败: %v", err)
	} else {
		log.Printf("更新组响应内容: %+v", responseBody)

		// 检查是否包含data字段
		if data, exists := responseBody["data"]; exists {
			log.Printf("✅ 更新组接口返回了更新后的实体数据: %+v", data)
		} else {
			log.Printf("❌ 更新组接口没有返回data字段")
		}
	}

	return resp
}

func sendUpdateTeamRequest(jsonData []byte) *http.Response {
	url := "http://localhost:8080/api/v1/teams/100001"
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("创建更新团队请求失败: %v", err)
		return nil
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer test_token")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("发送更新团队请求失败: %v", err)
		return nil
	}

	// 读取响应内容
	var responseBody map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&responseBody); err != nil {
		log.Printf("解析更新团队响应失败: %v", err)
	} else {
		log.Printf("更新团队响应内容: %+v", responseBody)

		// 检查是否包含data字段
		if data, exists := responseBody["data"]; exists {
			log.Printf("✅ 更新团队接口返回了更新后的实体数据: %+v", data)
		} else {
			log.Printf("❌ 更新团队接口没有返回data字段")
		}
	}

	return resp
}
