package interceptor

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/xinliangnote/go-gin-api/configs"
	"github.com/xinliangnote/go-gin-api/internal/code"
	"github.com/xinliangnote/go-gin-api/internal/pkg/core"
	"github.com/xinliangnote/go-gin-api/internal/pkg/token"
	"github.com/xinliangnote/go-gin-api/internal/proposal"
	"github.com/xinliangnote/go-gin-api/internal/repository/redis"
)

func (i *interceptor) CheckLogin(ctx core.Context) (sessionUserInfo proposal.SessionUserInfo, err core.BusinessError) {
	// 首先尝试从Token头部获取
	tokenStr := ctx.GetHeader(configs.HeaderLoginToken)

	// 如果Token头部为空，尝试从Authorization Bearer头部获取
	if tokenStr == "" {
		authHeader := ctx.GetHeader(configs.HeaderSignToken)
		if authHeader != "" && len(authHeader) > 7 && authHeader[:7] == "Bearer " {
			tokenStr = authHeader[7:] // 去掉"Bearer "前缀
		}
	}

	// 添加调试信息
	fmt.Printf("DEBUG: Token header = %s\n", ctx.GetHeader(configs.HeaderLoginToken))
	fmt.Printf("DEBUG: Authorization header = %s\n", ctx.GetHeader(configs.HeaderSignToken))
	fmt.Printf("DEBUG: Final extracted token = %s\n", tokenStr)

	if tokenStr == "" {
		fmt.Printf("DEBUG: Token is empty, returning error\n")
		err = core.Error(
			http.StatusUnauthorized,
			code.AuthorizationError,
			"Header 中缺少 Token 参数").WithError(errors.New("Header 中缺少 Token 参数"))

		return
	}

	// 解密token获取用户ID
	tokenData, decryptErr := token.DecryptTokenFromString(tokenStr)
	if decryptErr != nil {
		fmt.Printf("DEBUG: Token解密失败: %v\n", decryptErr)
		err = core.Error(
			http.StatusUnauthorized,
			code.AuthorizationError,
			"Token无效").WithError(decryptErr)

		return
	}

	// 直接通过用户ID查找session
	userSessionKey := fmt.Sprintf("%suser:%d", configs.RedisKeyPrefixLoginUser, tokenData.UserID)
	fmt.Printf("DEBUG: 查找用户session key: %s\n", userSessionKey)

	exists := i.cache.Exists(userSessionKey)
	fmt.Printf("DEBUG: 用户session是否存在: %v\n", exists)

	if !exists {
		fmt.Printf("DEBUG: 用户session不存在\n")
		err = core.Error(
			http.StatusUnauthorized,
			code.AuthorizationError,
			"请先登录").WithError(errors.New("请先登录"))

		return
	}

	cacheData, cacheErr := i.cache.Get(userSessionKey, redis.WithTrace(ctx.Trace()))
	if cacheErr != nil {
		fmt.Printf("DEBUG: Redis get error: %v\n", cacheErr)
		err = core.Error(
			http.StatusUnauthorized,
			code.AuthorizationError,
			"获取用户session失败").WithError(cacheErr)

		return
	}

	fmt.Printf("DEBUG: Retrieved cache data: %s\n", cacheData)

	jsonErr := json.Unmarshal([]byte(cacheData), &sessionUserInfo)
	if jsonErr != nil {
		fmt.Printf("DEBUG: JSON unmarshal error: %v\n", jsonErr)
		err = core.Error(
			http.StatusUnauthorized,
			code.AuthorizationError,
			"解析用户session数据失败").WithError(jsonErr)

		return
	}

	// 检查session是否过期（7天）
	if time.Now().Unix()-sessionUserInfo.LoginTime > 7*24*3600 {
		fmt.Printf("DEBUG: Session已过期\n")
		err = core.Error(
			http.StatusUnauthorized,
			code.AuthorizationError,
			"Session已过期").WithError(errors.New("Session已过期"))

		return
	}

	fmt.Printf("DEBUG: Successfully parsed session user info: %+v\n", sessionUserInfo)
	fmt.Printf("DEBUG: CheckLogin completed successfully\n")

	return
}

// 注意：加密解密逻辑已移至 internal/pkg/token/token.go 工具类中
