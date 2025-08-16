package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

func main() {
	log.Println("开始测试团队成员管理功能...")

	// 测试添加团队成员
	log.Println("\n=== 测试添加团队成员 ===")
	testAddTeamMember()

	// 测试更新成员角色
	log.Println("\n=== 测试更新成员角色 ===")
	testUpdateMemberRole()

	// 测试移除团队成员
	log.Println("\n=== 测试移除团队成员 ===")
	testRemoveTeamMember()

	log.Println("\n所有测试完成！")
}

func testAddTeamMember() {
	testData := map[string]interface{}{
		"accountId": "1001", // 使用账户ID
	}

	jsonData, _ := json.Marshal(testData)
	resp := sendAddMemberRequest(jsonData)

	if resp != nil && resp.StatusCode == 200 {
		log.Println("✅ 添加团队成员测试成功")
	} else {
		log.Printf("❌ 添加团队成员测试失败，状态码: %d", resp.StatusCode)
	}
}

func testUpdateMemberRole() {
	testData := map[string]interface{}{
		"roleType": "team_manager", // 更新为团队管理员
	}

	jsonData, _ := json.Marshal(testData)
	resp := sendUpdateMemberRoleRequest(jsonData)

	if resp != nil && resp.StatusCode == 200 {
		log.Println("✅ 更新成员角色测试成功")
	} else {
		log.Printf("❌ 更新成员角色测试失败，状态码: %d", resp.StatusCode)
	}
}

func testRemoveTeamMember() {
	// 移除成员不需要请求体，直接发送DELETE请求
	resp := sendRemoveMemberRequest()

	if resp != nil && resp.StatusCode == 200 {
		log.Println("✅ 移除团队成员测试成功")
	} else {
		log.Printf("❌ 移除团队成员测试失败，状态码: %d", resp.StatusCode)
	}
}

func sendAddMemberRequest(jsonData []byte) *http.Response {
	url := "http://localhost:8080/api/v1/teams/100001/members"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("创建添加成员请求失败: %v", err)
		return nil
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer test_token")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("发送添加成员请求失败: %v", err)
		return nil
	}

	// 读取响应内容
	var responseBody map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&responseBody); err != nil {
		log.Printf("解析添加成员响应失败: %v", err)
	} else {
		log.Printf("添加成员响应内容: %+v", responseBody)
	}

	return resp
}

func sendUpdateMemberRoleRequest(jsonData []byte) *http.Response {
	url := "http://localhost:8080/api/v1/teams/100001/members/1001/role"
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("创建更新成员角色请求失败: %v", err)
		return nil
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer test_token")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("发送更新成员角色请求失败: %v", err)
		return nil
	}

	// 读取响应内容
	var responseBody map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&responseBody); err != nil {
		log.Printf("解析更新成员角色响应失败: %v", err)
	} else {
		log.Printf("更新成员角色响应内容: %+v", responseBody)
	}

	return resp
}

func sendRemoveMemberRequest() *http.Response {
	url := "http://localhost:8080/api/v1/teams/100001/members/1001"
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		log.Printf("创建移除成员请求失败: %v", err)
		return nil
	}

	req.Header.Set("Authorization", "Bearer test_token")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("发送移除成员请求失败: %v", err)
		return nil
	}

	// 读取响应内容
	var responseBody map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&responseBody); err != nil {
		log.Printf("解析移除成员响应失败: %v", err)
	} else {
		log.Printf("移除成员响应内容: %+v", responseBody)
	}

	return resp
}
