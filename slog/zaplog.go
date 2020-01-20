package slog

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	//"github.com/shawnfeng/sutil/stime"
	"github.com/natefinch/lumberjack"
	"io"
	"os"
	"sync"
	"sync/atomic"
	"time"
)

// log 级别
const (
	LV_TRACE int = 0
	LV_DEBUG int = 1
	LV_INFO  int = 2
	LV_WARN  int = 3
	LV_ERROR int = 4
	LV_FATAL int = 5
	LV_PANIC int = 6
)

var (
	// log count
	cnTrace int64
	cnDebug int64
	cnInfo  int64
	cnWarn  int64
	cnError int64
	cnFatal int64
	cnPanic int64
	// log count stat stamp
	cnStamp   int64
	slogMutex sync.Mutex
	sLog  *zap.SugaredLogger
)

func TimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006/01/02 15:04:05.000000"))
}

func CapitalLevelEncoder(l zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(l.CapitalString())
}

func Sync() {
	sLog.Sync()
}

func Init(logdir string, level string) {
	InitV2(logdir, level, 10240000, 0, 0)
}

func InitV2(logDir, level string, maxSize int, maxAge, maxBackups int) {
	logLevel := zap.InfoLevel
	if level == "TRACE" {
		logLevel = zap.DebugLevel
	} else if level == "DEBUG" {
		logLevel = zap.DebugLevel
	} else if level == "INFO" {
		logLevel = zap.InfoLevel
	} else if level == "WARN" {
		logLevel = zap.WarnLevel
	} else if level == "ERROR" {
		logLevel = zap.ErrorLevel
	} else if level == "FATAL" {
		logLevel = zap.FatalLevel
	} else if level == "PANIC" {
		logLevel = zap.PanicLevel
	} else {
		logLevel = zap.InfoLevel
	}

	logInfoFile := ""
	if logDir != "" {
		logInfoFile = logDir + "/ser.log"
	}

	var out io.Writer
	if len(logDir) > 0 {
		lumberjackLogger := NewLogger(logInfoFile, maxSize, maxAge, maxBackups, true, false)
		go func() {
			for {
				now := time.Now().Unix()
				duration := 3600 - now%3600
				select {
				case <-time.After(time.Second * time.Duration(duration)):
					lumberjackLogger.Rotate()
				}
			}
		}()
		out = lumberjackLogger
	} else {
		out = os.Stdout
	}
	w := zapcore.AddSync(out)

	enconf := zap.NewProductionEncoderConfig()
	enconf.EncodeTime = TimeEncoder
	enconf.CallerKey = "caller"
	enconf.EncodeCaller = zapcore.FullCallerEncoder
	enconf.EncodeLevel = CapitalLevelEncoder
	core := zapcore.NewCore(
		//zapcore.NewJSONEncoder(enconf),
		zapcore.NewConsoleEncoder(enconf),
		w,
		logLevel,
	)
	logger := zap.New(core)
	sLog = logger.Sugar()
}

func init() {
	Init("", "TRACE")
	atomic.StoreInt64(&cnStamp, time.Now().Unix())
}

//type SLogger struct {
//}
//
//func GetLogger() *SLogger {
//	return &SLogger{}
//}
//
//func (m *SLogger) Printf(format string, items ...interface{}) {
//	Errorf(format, items...)
//}

func NewLogger(filename string, maxSize, maxAge, maxBackups int, localTime, compress bool) *lumberjack.Logger {
	return &lumberjack.Logger{
		Filename:   filename,
		MaxSize:    maxSize,
		MaxAge:     maxAge,
		MaxBackups: maxBackups,
		LocalTime:  localTime,
		Compress:   compress,
	}
}

//todo 调用写入日志

func formatFromContext(ctx context.Context, includeHead bool, format string) string {
	if cs := extractContextAsString(ctx, includeHead); cs != "" {
		return fmt.Sprintf("%s%s", cs, format)
	}
	return format
}

func vFromContext(ctx context.Context, includeHead bool, v ...interface{}) []interface{} {
	if cs := extractContextAsString(ctx, includeHead); len(cs) > 0 {
		return append([]interface{}{cs}, v...)
	}
	return v
}

func Tracef(ctx context.Context, format string, v ...interface{}) {
	format = formatFromContext(ctx, false, format)
	sLog.Debugf(format, v...)
	atomic.AddInt64(&cnTrace, 1)
}

func Traceln(ctx context.Context, v ...interface{}) {
	v = vFromContext(ctx, false, v...)
	sLog.Debug(v...)
	atomic.AddInt64(&cnTrace, 1)
}

func Debugf(ctx context.Context, format string, v ...interface{}) {
	format = formatFromContext(ctx, false, format)
	sLog.Debugf(format, v...)
	atomic.AddInt64(&cnDebug, 1)
}

func Debugln(ctx context.Context, v ...interface{}) {
	v = vFromContext(ctx, false, v...)
	sLog.Debug(v...)
	atomic.AddInt64(&cnDebug, 1)
}

func Infof(ctx context.Context, format string, v ...interface{}) {
	format = formatFromContext(ctx, false, format)
	sLog.Infof(format, v...)
	atomic.AddInt64(&cnInfo, 1)
}

func Infoln(ctx context.Context, v ...interface{}) {
	v = vFromContext(ctx, false, v...)
	sLog.Info(v...)
	atomic.AddInt64(&cnInfo, 1)
}

func Warnf(ctx context.Context, format string, v ...interface{}) {
	format = formatFromContext(ctx, false, format)
	sLog.Warnf(format, v...)
	atomic.AddInt64(&cnWarn, 1)
}

func Warnln(ctx context.Context, v ...interface{}) {
	v = vFromContext(ctx, false, v...)
	sLog.Warn(v...)
	atomic.AddInt64(&cnWarn, 1)
}

func Errorf(ctx context.Context, format string, v ...interface{}) {
	format = formatFromContext(ctx, true, format)
	sLog.Errorf(format, v...)
	atomic.AddInt64(&cnError, 1)
}

func Errorln(ctx context.Context, v ...interface{}) {
	v = vFromContext(ctx, true, v...)
	sLog.Error(v...)
	atomic.AddInt64(&cnError, 1)
}

func Fatalf(ctx context.Context, format string, v ...interface{}) {
	format = formatFromContext(ctx, true, format)
	sLog.Fatalf(format, v...)
	atomic.AddInt64(&cnFatal, 1)
}

func Fatalln(ctx context.Context, v ...interface{}) {
	v = vFromContext(ctx, true, v...)
	sLog.Fatal(v...)
	atomic.AddInt64(&cnFatal, 1)
}

func Panicf(ctx context.Context, format string, v ...interface{}) {
	format = formatFromContext(ctx, true, format)
	sLog.Panicf(format, v...)
	atomic.AddInt64(&cnPanic, 1)
}

func Panicln(ctx context.Context, v ...interface{}) {
	v = vFromContext(ctx, true, v...)
	sLog.Panic(v...)
	atomic.AddInt64(&cnPanic, 1)
}
