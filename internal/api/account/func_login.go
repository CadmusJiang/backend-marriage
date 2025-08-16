package account

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"crypto/rand"
	"encoding/base64"
	"strings"

	"github.com/xinliangnote/go-gin-api/configs"
	"github.com/xinliangnote/go-gin-api/internal/code"
	"github.com/xinliangnote/go-gin-api/internal/pkg/core"
	"github.com/xinliangnote/go-gin-api/internal/proposal"
	"github.com/xinliangnote/go-gin-api/internal/repository/redis"
	"go.uber.org/zap"
)

type loginRequest struct {
	Username string `json:"username" binding:"required"` // 用户名
	Password string `json:"password" binding:"required"` // 密码
}

type loginResponse struct {
	Status string `json:"status"`
	Data   struct {
		Token    string `json:"token"`
		Username string `json:"username"`
	} `json:"data"`
}

type logoutResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// Login 用户登录
// @Summary 用户登录
// @Description 用户登录接口
// @Tags CoreAuth
// @Accept application/json
// @Produce json
// @Param request body loginRequest true "登录请求"
// @Success 200 {object} loginResponse
// @Failure 400 {object} code.Failure
// @Router /api/v1/auth/login [post]
func (h *handler) Login() core.HandlerFunc {
	return func(c core.Context) {
		req := new(loginRequest)
		res := new(loginResponse)

		if err := c.ShouldBindJSON(req); err != nil {
			h.logger.Error("参数绑定失败", zap.Error(err))
			c.AbortWithError(core.Error(
				http.StatusBadRequest,
				code.ParamBindError,
				code.Text(code.ParamBindError)).WithError(err),
			)
			return
		}

		// 添加调试日志
		h.logger.Info("登录请求",
			zap.String("username", req.Username),
			zap.String("password", req.Password),
		)

		// 调用服务层进行登录验证
		accountInfo, err := h.accountService.Login(c, req.Username, req.Password)
		if err != nil {
			h.logger.Error("登录失败",
				zap.String("username", req.Username),
				zap.Error(err),
			)
			c.AbortWithError(core.Error(
				http.StatusUnauthorized,
				code.AdminLoginError,
				"用户名或密码错误").WithError(err),
			)
			return
		}

		// 添加成功日志
		h.logger.Info("登录成功",
			zap.String("username", req.Username),
			zap.Uint32("user_id", accountInfo.Id),
			zap.String("role_type", accountInfo.RoleType),
		)

		// 生成安全的session token
		token, err := generateSecureSessionToken()
		if err != nil {
			h.logger.Error("生成session token失败", zap.Error(err))
			c.AbortWithError(core.Error(
				http.StatusInternalServerError,
				code.ServerError,
				"登录失败").WithError(err),
			)
			return
		}

		// 将用户信息存储到Redis中
		sessionUserInfo := proposal.SessionUserInfo{
			UserID:    int32(accountInfo.Id),
			UserName:  accountInfo.Username,
			RoleType:  accountInfo.RoleType,
			Status:    accountInfo.Status,
			LoginTime: time.Now().Unix(), // 使用当前时间作为登录时间
		}

		// 序列化用户信息
		sessionData, err := json.Marshal(sessionUserInfo)
		if err != nil {
			h.logger.Error("序列化用户信息失败", zap.Error(err))
			c.AbortWithError(core.Error(
				http.StatusInternalServerError,
				code.ServerError,
				"登录失败").WithError(err),
			)
			return
		}

		// 存储到Redis，设置过期时间
		redisKey := configs.RedisKeyPrefixLoginUser + token
		err = h.cache.Set(redisKey, string(sessionData), configs.LoginSessionTTL, redis.WithTrace(c.Trace()))
		if err != nil {
			h.logger.Error("存储用户会话到Redis失败", zap.Error(err))
			c.AbortWithError(core.Error(
				http.StatusInternalServerError,
				code.ServerError,
				"登录失败").WithError(err),
			)
			return
		}

		h.logger.Info("用户会话已存储到Redis",
			zap.String("username", accountInfo.Username),
			zap.String("token", token),
			zap.String("redis_key", redisKey),
		)

		res.Status = "success"
		res.Data.Token = token
		res.Data.Username = accountInfo.Username

		c.Payload(res)
	}
}

// Logout 退出登录
// @Summary 退出登录
// @Description 用户退出登录
// @Tags CoreAuth
// @Accept application/x-www-form-urlencoded
// @Produce json
// @Success 200 {object} logoutResponse
// @Failure 400 {object} code.Failure
// @Router /api/v1/auth/logout [post]
func (h *handler) Logout() core.HandlerFunc {
	return func(c core.Context) {
		res := new(logoutResponse)

		// 从请求头获取token
		token := c.GetHeader(configs.HeaderLoginToken)
		if token == "" {
			// 尝试从Authorization Bearer头部获取
			authHeader := c.GetHeader(configs.HeaderSignToken)
			if authHeader != "" && len(authHeader) > 7 && authHeader[:7] == "Bearer " {
				token = authHeader[7:] // 去掉"Bearer "前缀
			}
		}

		// 验证token是否为空
		if token == "" {
			h.logger.Error("退出登录失败：缺少token")
			c.AbortWithError(core.Error(
				http.StatusBadRequest,
				code.ParamBindError,
				"缺少token参数").WithError(fmt.Errorf("token不能为空")),
			)
			return
		}

		// 验证token是否存在于Redis中
		redisKey := configs.RedisKeyPrefixLoginUser + token
		exists := h.cache.Exists(redisKey)
		if !exists {
			h.logger.Error("退出登录失败：token不存在或已过期",
				zap.String("token", token),
				zap.String("redis_key", redisKey),
			)
			c.AbortWithError(core.Error(
				http.StatusUnauthorized,
				code.AuthorizationError,
				"token无效或已过期").WithError(fmt.Errorf("token不存在或已过期")),
			)
			return
		}

		// 从Redis中删除token
		deleted := h.cache.Del(redisKey, redis.WithTrace(c.Trace()))
		if !deleted {
			h.logger.Error("删除Redis中的token失败",
				zap.String("token", token),
				zap.String("redis_key", redisKey),
			)
			c.AbortWithError(core.Error(
				http.StatusInternalServerError,
				code.ServerError,
				"退出登录失败").WithError(fmt.Errorf("删除token失败")),
			)
			return
		}

		h.logger.Info("用户会话已从Redis删除",
			zap.String("token", token),
			zap.String("redis_key", redisKey),
		)

		// 调用服务层退出登录
		logoutErr := h.accountService.Logout(c, token)
		if logoutErr != nil {
			h.logger.Error("退出登录失败", zap.Error(logoutErr))
		}

		res.Success = true
		res.Message = "退出登录成功"

		c.Payload(res)
	}
}

// generateSecureSessionToken 生成安全的session token
// 使用行业标准的32字节随机数 + Base64编码，确保token的唯一性和安全性
func generateSecureSessionToken() (string, error) {
	// 生成32字节的随机数
	randomBytes := make([]byte, 32)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", err
	}

	// 使用Base64编码，确保token只包含URL安全的字符
	// 去掉填充字符(=)以保持整洁
	token := base64.URLEncoding.EncodeToString(randomBytes)

	// 去掉Base64填充字符
	token = strings.TrimRight(token, "=")

	return token, nil
}
