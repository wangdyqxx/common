package server

import (
	"context"
	"github.com/wangdyqxx/common/log"
)

const (
	level0 = "TRACE"
	level1 = "DEBUG"
	level2 = "INFO"
	level3 = "WARN"
	level4 = "ERROR"
	level5 = "FATAL"
	level6 = "PANIC"
)

func initLog(ctx context.Context, dir, level string) {
	log.Init(dir, level)
}
