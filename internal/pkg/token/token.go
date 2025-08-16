package token

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"
)

// TokenData token数据结构
type TokenData struct {
	UserID    int32  `json:"user_id"`
	Timestamp int64  `json:"timestamp"`
	Random    string `json:"random"`
}

// GenerateSecureSessionToken 生成安全的加密session token
// 包含用户ID、时间戳、随机数，使用AES加密确保安全性
func GenerateSecureSessionToken(userID int32) (string, error) {
	// 构造token数据
	tokenData := TokenData{
		UserID:    userID,
		Timestamp: time.Now().Unix(),
		Random:    generateRandomString(16), // 16位随机字符串
	}

	// 序列化数据
	jsonData, err := json.Marshal(tokenData)
	if err != nil {
		return "", fmt.Errorf("序列化token数据失败: %v", err)
	}

	// 使用AES加密
	encryptedToken, err := encryptToken(jsonData)
	if err != nil {
		return "", fmt.Errorf("加密token失败: %v", err)
	}

	// Base64编码，确保URL安全，不去掉填充字符
	token := base64.URLEncoding.EncodeToString(encryptedToken)
	return token, nil
}

// DecryptTokenFromString 从字符串解密token数据
func DecryptTokenFromString(encryptedToken string) (*TokenData, error) {
	// Base64解码
	encryptedData, err := base64.URLEncoding.DecodeString(encryptedToken)
	if err != nil {
		return nil, fmt.Errorf("Base64解码失败: %v", err)
	}

	// 解密
	decryptedData, err := decryptToken(encryptedData)
	if err != nil {
		return nil, fmt.Errorf("解密失败: %v", err)
	}

	// 解析数据
	var tokenData TokenData
	if err := json.Unmarshal(decryptedData, &tokenData); err != nil {
		return nil, fmt.Errorf("解析token数据失败: %v", err)
	}

	return &tokenData, nil
}

// ValidateToken 验证token是否有效
func ValidateToken(token string) (*TokenData, error) {
	tokenData, err := DecryptTokenFromString(token)
	if err != nil {
		return nil, err
	}

	// 检查token是否过期（7天）
	if time.Now().Unix()-tokenData.Timestamp > 7*24*3600 {
		return nil, fmt.Errorf("token已过期")
	}

	return tokenData, nil
}

// generateRandomString 生成指定长度的随机字符串
func generateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	randomBytes := make([]byte, length)
	rand.Read(randomBytes)
	for i := range b {
		b[i] = charset[int(randomBytes[i])%len(charset)]
	}
	return string(b)
}

// encryptToken 使用AES加密token数据
func encryptToken(data []byte) ([]byte, error) {
	// 从配置获取密钥
	key := getEncryptionKey()
	if len(key) != 32 {
		return nil, fmt.Errorf("密钥长度必须为32字节，当前长度: %d", len(key))
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// 生成随机IV
	iv := make([]byte, aes.BlockSize)
	if _, err := rand.Read(iv); err != nil {
		return nil, err
	}

	// 加密
	ciphertext := make([]byte, len(data)+aes.BlockSize)
	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], data)
	copy(ciphertext[:aes.BlockSize], iv)

	return ciphertext, nil
}

// decryptToken 解密token数据
func decryptToken(encryptedData []byte) ([]byte, error) {
	// 从配置获取密钥
	key := getEncryptionKey()
	if len(key) != 32 {
		return nil, fmt.Errorf("密钥长度必须为32字节，当前长度: %d", len(key))
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	if len(encryptedData) < aes.BlockSize {
		return nil, fmt.Errorf("加密数据长度不足")
	}

	// 提取IV
	iv := encryptedData[:aes.BlockSize]
	ciphertext := encryptedData[aes.BlockSize:]

	// 解密
	plaintext := make([]byte, len(ciphertext))
	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(plaintext, ciphertext)

	return plaintext, nil
}

// getEncryptionKey 获取加密密钥
func getEncryptionKey() []byte {
	// 这里应该从配置文件读取，暂时使用硬编码
	// TODO: 从配置文件读取密钥
	return []byte("marriage-system-secretkey-aeskey")
}
