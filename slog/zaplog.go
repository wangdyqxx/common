package slog

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"runtime"

	//"github.com/shawnfeng/sutil/stime"
	"github.com/natefinch/lumberjack"
	"io"
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
	sLog      *zap.SugaredLogger
)

func init() {
	Init("", "TRACE")
	atomic.StoreInt64(&cnStamp, time.Now().Unix())
}

//初始化
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
	var out io.Writer
	if logDir == "" {
		out = os.Stdout
	} else {
		logInfoFile = logDir + "/ser.log"
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

func TimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006/01/02 15:04:05.000000"))
}

func CapitalLevelEncoder(l zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(l.CapitalString())
}

func Sync() {
	sLog.Sync()
}

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

func formatFromContext(ctx context.Context, includeHead bool, format string, isCaller bool) string {
	if cs := extractContextAsString(ctx, includeHead); cs != "" {
		if isCaller {
			_, file, line, ok := runtime.Caller(2)
			if ok {
				cs = fmt.Sprintf("%s %s:%d", cs, file, line)
			}
		}
		return fmt.Sprintf("%s %s", cs, format)
	}
	return format
}

func vFromContext(ctx context.Context, includeHead bool, v ...interface{}) []interface{} {
	if cs := extractContextAsString(ctx, includeHead); len(cs) > 0 {
		return append([]interface{}{cs}, v...)
	}
	return v
}

//todo 可调用的写入日志函数

func Tracef(ctx context.Context, format string, v ...interface{}) {
	format = formatFromContext(ctx, false, format, false)
	sLog.Debugf(format, v...)
	atomic.AddInt64(&cnTrace, 1)
}

func Traceln(ctx context.Context, v ...interface{}) {
	v = vFromContext(ctx, false, v...)
	sLog.Debug(v...)
	atomic.AddInt64(&cnTrace, 1)
}

func Debugf(ctx context.Context, format string, v ...interface{}) {
	format = formatFromContext(ctx, false, format, false)
	sLog.Debugf(format, v...)
	atomic.AddInt64(&cnDebug, 1)
}

func Debugln(ctx context.Context, v ...interface{}) {
	v = vFromContext(ctx, false, v...)
	sLog.Debug(v...)
	atomic.AddInt64(&cnDebug, 1)
}

func Infof(ctx context.Context, format string, v ...interface{}) {
	format = formatFromContext(ctx, false, format, false)
	sLog.Infof(format, v...)
	atomic.AddInt64(&cnInfo, 1)
}

func Info(format string, v ...interface{}) {
	format = formatFromContext(context.TODO(),false, format, false)
	sLog.Infof(format, v...)
	atomic.AddInt64(&cnInfo, 1)
}

func Infoln(ctx context.Context, v ...interface{}) {
	v = vFromContext(ctx, false, v...)
	sLog.Info(v...)
	atomic.AddInt64(&cnInfo, 1)
}

func Warnf(ctx context.Context, format string, v ...interface{}) {
	format = formatFromContext(ctx, false, format, false)
	sLog.Warnf(format, v...)
	atomic.AddInt64(&cnWarn, 1)
}

func Warnln(ctx context.Context, v ...interface{}) {
	v = vFromContext(ctx, false, v...)
	sLog.Warn(v...)
	atomic.AddInt64(&cnWarn, 1)
}

func Errorf(ctx context.Context, format string, v ...interface{}) {
	format = formatFromContext(ctx, true, format, true)
	sLog.Errorf(format, v...)
	atomic.AddInt64(&cnError, 1)
}

func Fatalf(ctx context.Context, format string, v ...interface{}) {
	format = formatFromContext(ctx, true, format, true)
	sLog.Fatalf(format, v...)
	atomic.AddInt64(&cnFatal, 1)
}

func Panicf(ctx context.Context, format string, v ...interface{}) {
	format = formatFromContext(ctx, true, format, true)
	sLog.Panicf(format, v...)
	atomic.AddInt64(&cnPanic, 1)
}
