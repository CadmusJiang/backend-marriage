package mysql

import (
	"fmt"
	"strings"
	"time"

	"github.com/xinliangnote/go-gin-api/configs"
	"github.com/xinliangnote/go-gin-api/pkg/errors"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

// Predicate is a string that acts as a condition in the where clause
type Predicate string

var (
	EqualPredicate              = Predicate("=")
	NotEqualPredicate           = Predicate("<>")
	GreaterThanPredicate        = Predicate(">")
	GreaterThanOrEqualPredicate = Predicate(">=")
	SmallerThanPredicate        = Predicate("<")
	SmallerThanOrEqualPredicate = Predicate("<=")
	LikePredicate               = Predicate("LIKE")
)

var _ Repo = (*dbRepo)(nil)

type Repo interface {
	i()
	GetDbR() *gorm.DB
	GetDbW() *gorm.DB
	DbRClose() error
	DbWClose() error
}

type dbRepo struct {
	DbR *gorm.DB
	DbW *gorm.DB
}

func New() (Repo, error) {
	cfg := configs.Get().MySQL

	// 连接读数据库
	dbr, err := dbConnect(cfg.Read.User, cfg.Read.Pass, cfg.Read.Addr, cfg.Read.Name, "read")
	if err != nil {
		return nil, fmt.Errorf("读数据库连接失败: %v", err)
	}

	// 连接写数据库
	dbw, err := dbConnect(cfg.Write.User, cfg.Write.Pass, cfg.Write.Addr, cfg.Write.Name, "write")
	if err != nil {
		return nil, fmt.Errorf("写数据库连接失败: %v", err)
	}

	return &dbRepo{
		DbR: dbr,
		DbW: dbw,
	}, nil
}

func (d *dbRepo) i() {}

func (d *dbRepo) GetDbR() *gorm.DB {
	return d.DbR
}

func (d *dbRepo) GetDbW() *gorm.DB {
	return d.DbW
}

func (d *dbRepo) DbRClose() error {
	sqlDB, err := d.DbR.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

func (d *dbRepo) DbWClose() error {
	sqlDB, err := d.DbW.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

func dbConnect(user, pass, addr, dbName, connType string) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=%t&loc=%s&allowNativePasswords=true",
		user,
		pass,
		addr,
		dbName,
		true,
		"Local")

	// 记录连接信息（不包含密码）
	connInfo := fmt.Sprintf("%s@%s/%s", user, addr, dbName)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		Logger: logger.Default.LogMode(logger.Error), // 只记录错误，不记录所有SQL
	})

	if err != nil {
		// 提供详细的错误诊断信息
		var errorMsg string
		if strings.Contains(err.Error(), "Access denied") {
			errorMsg = fmt.Sprintf("数据库访问被拒绝，请检查用户名和密码是否正确。连接信息: %s", connInfo)
		} else if strings.Contains(err.Error(), "Unknown database") {
			errorMsg = fmt.Sprintf("数据库不存在，请检查数据库名称是否正确。连接信息: %s", connInfo)
		} else if strings.Contains(err.Error(), "connection refused") {
			errorMsg = fmt.Sprintf("无法连接到数据库服务器，请检查服务器地址和端口是否正确。连接信息: %s", connInfo)
		} else if strings.Contains(err.Error(), "timeout") {
			errorMsg = fmt.Sprintf("数据库连接超时，请检查网络连接和服务器状态。连接信息: %s", connInfo)
		} else {
			errorMsg = fmt.Sprintf("数据库连接失败: %v。连接信息: %s", err, connInfo)
		}
		return nil, errors.Wrap(err, errorMsg)
	}

	db.Set("gorm:table_options", "CHARSET=utf8mb4")

	cfg := configs.Get().MySQL.Base

	sqlDB, err := db.DB()
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("获取数据库连接池失败: %s", connInfo))
	}

	// 设置连接池 用于设置最大打开的连接数，默认值为0表示不限制.设置最大的连接数，可以避免并发太高导致连接mysql出现too many connections的错误。
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConn)

	// 设置最大连接数 用于设置闲置的连接数.设置闲置的连接数则当开启的一个连接使用完成后可以放在池里等候下一次使用。
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConn)

	// 设置最大连接超时
	sqlDB.SetConnMaxLifetime(time.Minute * cfg.ConnMaxLifeTime)

	// 使用插件
	db.Use(&TracePlugin{})

	return db, nil
}
