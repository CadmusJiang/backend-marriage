package interceptor

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/xinliangnote/go-gin-api/configs"
	"github.com/xinliangnote/go-gin-api/internal/code"
	"github.com/xinliangnote/go-gin-api/internal/pkg/core"
	"github.com/xinliangnote/go-gin-api/internal/proposal"
	"github.com/xinliangnote/go-gin-api/internal/repository/redis"
)

func (i *interceptor) CheckLogin(ctx core.Context) (sessionUserInfo proposal.SessionUserInfo, err core.BusinessError) {
	// 首先尝试从Token头部获取
	token := ctx.GetHeader(configs.HeaderLoginToken)

	// 如果Token头部为空，尝试从Authorization Bearer头部获取
	if token == "" {
		authHeader := ctx.GetHeader(configs.HeaderSignToken)
		if authHeader != "" && len(authHeader) > 7 && authHeader[:7] == "Bearer " {
			token = authHeader[7:] // 去掉"Bearer "前缀
		}
	}

	// 添加调试信息
	fmt.Printf("DEBUG: Token header = %s\n", ctx.GetHeader(configs.HeaderLoginToken))
	fmt.Printf("DEBUG: Authorization header = %s\n", ctx.GetHeader(configs.HeaderSignToken))
	fmt.Printf("DEBUG: Final extracted token = %s\n", token)

	if token == "" {
		fmt.Printf("DEBUG: Token is empty, returning error\n")
		err = core.Error(
			http.StatusUnauthorized,
			code.AuthorizationError,
			code.Text(code.AuthorizationError)).WithError(errors.New("Header 中缺少 Token 参数"))

		return
	}

	redisKey := configs.RedisKeyPrefixLoginUser + token
	fmt.Printf("DEBUG: RedisKeyPrefixLoginUser = '%s'\n", configs.RedisKeyPrefixLoginUser)
	fmt.Printf("DEBUG: Token = '%s'\n", token)
	fmt.Printf("DEBUG: Complete Redis key = '%s'\n", redisKey)
	fmt.Printf("DEBUG: Redis key length = %d\n", len(redisKey))

	exists := i.cache.Exists(redisKey)
	fmt.Printf("DEBUG: Redis key exists: %v\n", exists)

	// 尝试获取所有匹配的key来调试
	fmt.Printf("DEBUG: Trying to list all keys with pattern: %s*\n", configs.RedisKeyPrefixLoginUser)

	if !exists {
		fmt.Printf("DEBUG: Redis key does not exist, returning error\n")
		err = core.Error(
			http.StatusUnauthorized,
			code.AuthorizationError,
			code.Text(code.AuthorizationError)).WithError(errors.New("请先登录"))

		return
	}

	cacheData, cacheErr := i.cache.Get(redisKey, redis.WithTrace(ctx.Trace()))
	if cacheErr != nil {
		fmt.Printf("DEBUG: Redis get error: %v\n", cacheErr)
		err = core.Error(
			http.StatusUnauthorized,
			code.AuthorizationError,
			code.Text(code.AuthorizationError)).WithError(cacheErr)

		return
	}

	fmt.Printf("DEBUG: Retrieved cache data: %s\n", cacheData)

	jsonErr := json.Unmarshal([]byte(cacheData), &sessionUserInfo)
	if jsonErr != nil {
		fmt.Printf("DEBUG: JSON unmarshal error: %v\n", jsonErr)
		err = core.Error(
			http.StatusUnauthorized,
			code.AuthorizationError,
			code.Text(code.AuthorizationError)).WithError(jsonErr)

		return
	}

	fmt.Printf("DEBUG: Successfully parsed session user info: %+v\n", sessionUserInfo)
	fmt.Printf("DEBUG: CheckLogin completed successfully\n")

	return
}
