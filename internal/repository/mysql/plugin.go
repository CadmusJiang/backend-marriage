package mysql

import (
	"fmt"
	"strings"
	"time"

	"github.com/xinliangnote/go-gin-api/internal/pkg/core"
	"github.com/xinliangnote/go-gin-api/pkg/timeutil"
	"github.com/xinliangnote/go-gin-api/pkg/trace"

	"gorm.io/gorm"
	"gorm.io/gorm/utils"
)

const (
	callBackBeforeName = "core:before"
	callBackAfterName  = "core:after"
	startTime          = "_start_time"
)

type TracePlugin struct{}

func (op *TracePlugin) Name() string {
	return "tracePlugin"
}

func (op *TracePlugin) Initialize(db *gorm.DB) (err error) {
	// 开始前
	_ = db.Callback().Create().Before("gorm:before_create").Register(callBackBeforeName, before)
	_ = db.Callback().Query().Before("gorm:query").Register(callBackBeforeName, before)
	_ = db.Callback().Delete().Before("gorm:before_delete").Register(callBackBeforeName, before)
	_ = db.Callback().Update().Before("gorm:setup_reflect_value").Register(callBackBeforeName, before)
	_ = db.Callback().Row().Before("gorm:row").Register(callBackBeforeName, before)
	_ = db.Callback().Raw().Before("gorm:raw").Register(callBackBeforeName, before)

	// 结束后
	_ = db.Callback().Create().After("gorm:after_create").Register(callBackAfterName, after)
	_ = db.Callback().Query().After("gorm:after_query").Register(callBackAfterName, after)
	_ = db.Callback().Delete().After("gorm:after_delete").Register(callBackAfterName, after)
	_ = db.Callback().Update().After("gorm:after_update").Register(callBackAfterName, after)
	_ = db.Callback().Row().After("gorm:row").Register(callBackAfterName, after)
	_ = db.Callback().Raw().After("gorm:raw").Register(callBackAfterName, after)
	return
}

var _ gorm.Plugin = &TracePlugin{}

func before(db *gorm.DB) {
	db.InstanceSet(startTime, time.Now())
	return
}

func after(db *gorm.DB) {
	_ctx := db.Statement.Context
	ctx, ok := _ctx.(core.StdContext)
	if !ok {
		// 如果没有正确的context，直接返回
		return
	}

	_ts, isExist := db.InstanceGet(startTime)
	if !isExist {
		return
	}

	ts, ok := _ts.(time.Time)
	if !ok {
		return
	}

	sql := db.Dialector.Explain(db.Statement.SQL.String(), db.Statement.Vars...)

	// 过滤掉外表的扫描SQL，避免日志过多
	if shouldFilterSQL(sql) {
		return
	}

	sqlInfo := new(trace.SQL)
	sqlInfo.Timestamp = timeutil.CSTLayoutString()
	sqlInfo.SQL = sql
	sqlInfo.Stack = utils.FileWithLineNum()

	// 根据SQL类型设置正确的行数信息
	sqlLower := strings.ToLower(sql)
	if strings.HasPrefix(sqlLower, "select") {
		// SELECT查询：使用RowsAffected（通常为0）或者尝试获取结果行数
		if db.Statement.Dest != nil {
			// 尝试从结果中获取行数
			switch v := db.Statement.Dest.(type) {
			case *[]interface{}:
				sqlInfo.Rows = int64(len(*v))
			case *[]map[string]interface{}:
				sqlInfo.Rows = int64(len(*v))
			default:
				// 对于单个结果，如果有结果则为1，否则为0
				if db.Statement.RowsAffected > 0 {
					sqlInfo.Rows = db.Statement.RowsAffected
				} else {
					// 检查是否有结果
					if db.Statement.RowsAffected == 0 && db.Statement.DB.Error == nil {
						sqlInfo.Rows = 1 // 假设查询成功返回了结果
					} else {
						sqlInfo.Rows = 0
					}
				}
			}
		} else {
			sqlInfo.Rows = 0
		}
	} else {
		// UPDATE, INSERT, DELETE：使用RowsAffected
		sqlInfo.Rows = db.Statement.RowsAffected
	}

	sqlInfo.CostSeconds = time.Since(ts).Seconds()

	if ctx.Trace != nil {
		ctx.Trace.AppendSQL(sqlInfo)
	} else {
		// 调试信息：如果没有trace，说明可能有问题
		fmt.Printf("DEBUG: No trace found for SQL: %s\n", sql)
	}

	return
}

// shouldFilterSQL 判断是否应该过滤掉这个SQL
func shouldFilterSQL(sql string) bool {
	// 转换为小写进行匹配
	sqlLower := strings.ToLower(sql)

	// 过滤条件列表 - 只过滤outbox相关的定期扫描SQL
	filterPatterns := []string{
		// 过滤outbox相关的定期扫描SQL，但不过滤正常的业务查询
		"from outbox_events where status = 0 and created_at <",
		// 可以根据需要添加更多过滤条件
	}

	for _, pattern := range filterPatterns {
		if strings.Contains(sqlLower, pattern) {
			return true
		}
	}

	return false
}
