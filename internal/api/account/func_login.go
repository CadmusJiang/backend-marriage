package account

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/xinliangnote/go-gin-api/configs"
	"github.com/xinliangnote/go-gin-api/internal/code"
	"github.com/xinliangnote/go-gin-api/internal/pkg/core"
	"github.com/xinliangnote/go-gin-api/internal/pkg/token"
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
		token, err := token.GenerateSecureSessionToken(int32(accountInfo.Id))
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

		// 清理该用户之前的旧session（防止同一用户多个session）
		userSessionKey := fmt.Sprintf("%suser:%d", configs.RedisKeyPrefixLoginUser, accountInfo.Id)
		oldSessionData, oldSessionErr := h.cache.Get(userSessionKey, redis.WithTrace(c.Trace()))
		if oldSessionErr == nil && oldSessionData != "" {
			// 解析旧session数据，获取旧token
			var oldSessionUserInfo proposal.SessionUserInfo
			if json.Unmarshal([]byte(oldSessionData), &oldSessionUserInfo) == nil {
				// 这里可以记录日志，但不删除旧数据（因为会被覆盖）
				h.logger.Info("用户重新登录，将覆盖旧session",
					zap.String("username", accountInfo.Username),
					zap.Uint32("user_id", accountInfo.Id),
				)
			}
		}

		// 直接存储用户session数据到用户ID对应的key
		err = h.cache.Set(userSessionKey, string(sessionData), configs.LoginSessionTTL, redis.WithTrace(c.Trace()))
		if err != nil {
			h.logger.Error("存储用户会话到Redis失败", zap.Error(err))
			c.AbortWithError(core.Error(
				http.StatusInternalServerError,
				code.ServerError,
				"登录失败").WithError(err),
			)
			return
		}

		// 注意：这里不再存储token映射，因为token本身包含用户ID信息
		// CheckLogin拦截器可以直接解密token获取用户ID，然后查找用户session

		h.logger.Info("用户会话已存储到Redis",
			zap.String("username", accountInfo.Username),
			zap.String("token", token),
			zap.String("user_session_key", userSessionKey),
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
		tokenStr := c.GetHeader(configs.HeaderLoginToken)
		if tokenStr == "" {
			// 尝试从Authorization Bearer头部获取
			authHeader := c.GetHeader(configs.HeaderSignToken)
			if authHeader != "" && len(authHeader) > 7 && authHeader[:7] == "Bearer " {
				tokenStr = authHeader[7:] // 去掉"Bearer "前缀
			}
		}

		// 验证token是否为空
		if tokenStr == "" {
			h.logger.Error("退出登录失败：缺少token")
			c.AbortWithError(core.Error(
				http.StatusBadRequest,
				code.ParamBindError,
				"缺少token参数").WithError(fmt.Errorf("token不能为空")),
			)
			return
		}

		// 解密token获取用户ID
		tokenData, decryptErr := token.DecryptTokenFromString(tokenStr)
		if decryptErr != nil {
			h.logger.Error("退出登录失败：token解密失败", zap.Error(decryptErr))
			c.AbortWithError(core.Error(
				http.StatusUnauthorized,
				code.AuthorizationError,
				"token无效").WithError(decryptErr),
			)
			return
		}

		// 直接通过用户ID查找并删除session
		userSessionKey := fmt.Sprintf("%suser:%d", configs.RedisKeyPrefixLoginUser, tokenData.UserID)
		exists := h.cache.Exists(userSessionKey)
		if !exists {
			h.logger.Error("退出登录失败：用户session不存在",
				zap.Int32("user_id", tokenData.UserID),
				zap.String("user_session_key", userSessionKey),
			)
			c.AbortWithError(core.Error(
				http.StatusUnauthorized,
				code.AuthorizationError,
				"用户session不存在").WithError(fmt.Errorf("用户session不存在")),
			)
			return
		}

		// 从Redis中删除用户session
		deleted := h.cache.Del(userSessionKey, redis.WithTrace(c.Trace()))
		if !deleted {
			h.logger.Error("删除Redis中的用户session失败",
				zap.Int32("user_id", tokenData.UserID),
				zap.String("user_session_key", userSessionKey),
			)
			c.AbortWithError(core.Error(
				http.StatusInternalServerError,
				code.ServerError,
				"退出登录失败").WithError(fmt.Errorf("删除用户session失败")),
			)
			return
		}

		h.logger.Info("用户会话已从Redis删除",
			zap.String("token", tokenStr),
			zap.Int32("user_id", tokenData.UserID),
			zap.String("user_session_key", userSessionKey),
		)

		// 调用服务层退出登录
		logoutErr := h.accountService.Logout(c, tokenStr)
		if logoutErr != nil {
			h.logger.Error("退出登录失败", zap.Error(logoutErr))
		}

		res.Success = true
		res.Message = "退出登录成功"

		c.Payload(res)
	}
}

// 注意：加密解密逻辑已移至 internal/pkg/token/token.go 工具类中
