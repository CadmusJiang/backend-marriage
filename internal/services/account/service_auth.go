package account

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/xinliangnote/go-gin-api/internal/pkg/core"
	"github.com/xinliangnote/go-gin-api/internal/repository/mysql"
	"github.com/xinliangnote/go-gin-api/internal/repository/mysql/account"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

// Login 用户登录
func (s *service) Login(ctx core.Context, username, password string) (accountInfo *account.Account, err error) {
	// 添加调试日志
	ctx.Logger().Info("开始登录验证",
		zap.String("username", username),
		zap.String("password_length", fmt.Sprintf("%d", len(password))),
	)

	// 检查数据库连接
	if s.db == nil {
		ctx.Logger().Error("数据库连接未初始化")
		return nil, fmt.Errorf("数据库连接未初始化")
	}

	dbR := s.db.GetDbR()
	if dbR == nil {
		ctx.Logger().Error("数据库读连接未初始化")
		return nil, fmt.Errorf("数据库读连接未初始化")
	}

	// 查询用户是否存在
	accountQueryBuilder := account.NewQueryBuilder()
	accountQueryBuilder.WhereUsername(mysql.EqualPredicate, username)

	ctx.Logger().Info("查询用户信息", zap.String("username", username))

	// 确保数据库查询传递trace信息
	dbR = dbR.WithContext(ctx.RequestContext())
	accountInfo, err = accountQueryBuilder.QueryOne(dbR)
	if err != nil {
		ctx.Logger().Error("查询账户失败",
			zap.String("username", username),
			zap.Error(err),
		)
		return nil, fmt.Errorf("查询账户失败: %v", err)
	}
	if accountInfo == nil {
		ctx.Logger().Warn("用户不存在", zap.String("username", username))
		return nil, fmt.Errorf("用户不存在")
	}

	ctx.Logger().Info("用户查询成功",
		zap.String("username", username),
		zap.Uint32("user_id", accountInfo.Id),
		zap.String("status", accountInfo.Status),
	)

	// 验证密码
	// 支持bcrypt和MD5两种密码格式
	passwordMatch := false

	// 首先尝试bcrypt验证
	if err := bcrypt.CompareHashAndPassword([]byte(accountInfo.Password), []byte(password)); err == nil {
		passwordMatch = true
	} else {
		// 如果bcrypt验证失败，尝试MD5验证
		hashedPassword := s.generateMD5(password)
		passwordMatch = accountInfo.Password == hashedPassword
	}

	ctx.Logger().Info("密码验证",
		zap.String("stored_password_hash", accountInfo.Password),
		zap.Bool("password_match", passwordMatch),
	)

	if !passwordMatch {
		ctx.Logger().Warn("密码错误",
			zap.String("username", username),
			zap.String("stored_password_hash", accountInfo.Password),
		)
		return nil, fmt.Errorf("密码错误")
	}

	// 检查账户状态
	if accountInfo.Status != "enabled" { // "enabled"表示启用，"disabled"表示禁用
		ctx.Logger().Warn("账户已被禁用",
			zap.String("username", username),
			zap.String("status", accountInfo.Status),
		)
		return nil, fmt.Errorf("账户已被禁用")
	}

	// 更新最后登录时间
	updateFields := map[string]interface{}{
		"LastLoginAt": time.Now(),
		"UpdatedUser": username,
	}

	// 确保数据库更新传递trace信息
	dbW := s.db.GetDbW().WithContext(ctx.RequestContext())
	err = accountQueryBuilder.Updates(dbW, updateFields)
	if err != nil {
		// 记录日志但不影响登录流程
		ctx.Logger().Error("更新登录时间失败", zap.Error(err))
	} else {
		ctx.Logger().Info("登录时间更新成功", zap.String("username", username))
	}

	ctx.Logger().Info("登录验证完成",
		zap.String("username", username),
		zap.Uint32("user_id", accountInfo.Id),
	)

	return accountInfo, nil
}

// Logout 退出登录
func (s *service) Logout(ctx core.Context, token string) (err error) {
	// 这里可以实现token黑名单等逻辑
	// 暂时只是简单返回成功
	return nil
}

// generateMD5 生成MD5哈希
func (s *service) generateMD5(str string) string {
	m := md5.New()
	m.Write([]byte(str))
	return hex.EncodeToString(m.Sum(nil))
}
