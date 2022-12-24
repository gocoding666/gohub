package logger

import (
	"context"
	"errors"
	"go.uber.org/zap"
	"gohub/pkg/helpers"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

// GormLogger
// @Description: 操作对象，实现gormlogger.Interface
type GormLogger struct {
	ZapLogger     *zap.Logger
	SlowThreshold time.Duration
}

// NewGormLogger 外部调用。实例化一个GormLogger对象，示例：
// DB,err:=gorm.Open(dbConfig,&gorm.Config{Logger:logger.NewGormLogger()})
func NewGormLogger() GormLogger {
	return GormLogger{
		ZapLogger:     Logger,                 //使用全局的 logger.Logger对象
		SlowThreshold: 200 * time.Millisecond, //慢查询阈值，单位为千分之一秒
	}
}

// LogMode 实现gormlogger.Interface的logMode方法
func (l GormLogger) LogMode(level gormlogger.LogLevel) gormlogger.Interface {
	return GormLogger{
		ZapLogger:     l.ZapLogger,
		SlowThreshold: l.SlowThreshold,
	}
}
func (l GormLogger) Info(ctx context.Context, str string, args ...interface{}) {
	l.logger().Sugar().Debugf(str, args...)
}
func (l GormLogger) Warn(ctx context.Context, str string, args ...interface{}) {
	l.logger().Sugar().Warnf(str, args...)
}
func (l GormLogger) Error(ctx context.Context, str string, args ...interface{}) {
	l.logger().Sugar().Errorf(str, args...)
}
func (l GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	//获取运行时间
	elapsed := time.Since(begin)
	//获取SQL请求和返回条数
	sql, rows := fc()
	//通用字段
	logFields := []zap.Field{
		zap.String("sql", sql),
		zap.String("time", helpers.MicrosecondsStr(elapsed)),
		zap.Int64("rows", rows),
	}
	//Gorm错误
	if err != nil {
		//记录未找到的错误使用 warning等级
		if errors.Is(err, gorm.ErrRecordNotFound) {
			l.logger().Warn("Database ErrReordNotFound", logFields...)
		} else {
			//其他错误使用error等级
			logFields = append(logFields, zap.Error(err))
			l.logger().Error("Database Error", logFields...)
		}
	}
	//慢查询日志
	if l.SlowThreshold != 0 && elapsed > l.SlowThreshold {
		l.logger().Warn("Database Slow log", logFields...)
	}
	//记录所有SQL请求
	l.logger().Debug("Database Query", logFields...)
}

// logger内用的辅助方法，确保Zap内置信息Caller的准确性（如paginator/paginator.go:148）
func (l GormLogger) logger() *zap.Logger {
	//跳过gorm内置的调用
	var (
		gormPackage    = filepath.Join("gorm.io", "gorm")
		zapgormPackage = filepath.Join("moul.io", "zapgorm2")
	)
	//减去一次封装，以及一次在logger初始化里添加zap.AddCallerSkip(1)
	clone := l.ZapLogger.WithOptions(zap.AddCallerSkip(-2))

	for i := 2; i < 15; i++ {
		_, file, _, ok := runtime.Caller(i)
		switch {
		case !ok:
		case strings.HasSuffix(file, "_test.go"):
		case strings.Contains(file, gormPackage):
		case strings.Contains(file, zapgormPackage):
		default:
			//返回一个附带跳过行号的新的zap logger
			return clone.WithOptions(zap.AddCallerSkip(i))
		}
	}
	return l.ZapLogger
}
