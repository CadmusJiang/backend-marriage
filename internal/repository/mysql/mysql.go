package mysql

import (
	"context"
	"fmt"
	"time"

	"github.com/xinliangnote/go-gin-api/configs"
	"github.com/xinliangnote/go-gin-api/pkg/trace"

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

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		NowFunc: func() time.Time {
			return time.Now()
		},
	})

	if err != nil {
		return nil, fmt.Errorf("连接%s数据库失败: %v", connType, err)
	}

	// 添加Trace支持
	db = db.Session(&gorm.Session{
		Logger: &traceLogger{},
	})

	return db, nil
}

// traceLogger 实现GORM的Logger接口，用于记录SQL执行信息
type traceLogger struct{}

func (l *traceLogger) LogMode(level logger.LogLevel) logger.Interface {
	return l
}

func (l *traceLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	// 暂时简化实现，后续完善Trace功能
}

func (l *traceLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	l.Info(ctx, msg, data...)
}

func (l *traceLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	l.Info(ctx, msg, data...)
}

func (l *traceLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	// 暂时简化实现，后续完善Trace功能
}

// getTraceFromContext 从context中获取Trace信息
func getTraceFromContext(ctx context.Context) *trace.Trace {
	// 这里需要根据实际的context类型来获取Trace
	// 暂时返回nil，后续会完善
	return nil
}
