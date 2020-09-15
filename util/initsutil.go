package util

import (
	"context"
	"github.com/opentracing/opentracing-go"
	"github.com/wangdyqxx/common/log"
	"github.com/wangdyqxx/common/tracer"
	"io"
	"os"
)

var GinDefaultLog io.Writer

func init() {
	//先初始化trace，后初始化log
	_ = tracer.InitDefaultTracer("log test")
	tracer := opentracing.GlobalTracer()
	span := tracer.StartSpan("testlog")
	ctx := context.Background()
	ctx = opentracing.ContextWithSpan(ctx, span)
	//dir := fileDir
	dir := ""
	log.Init(dir, "INFO")
	ginLog, err := os.Create("stat.log")
	if err != nil {
		panic(err)
	}
	GinDefaultLog = io.MultiWriter(ginLog)
	//GinDefaultLog = io.MultiWriter(ginLog, os.Stdout, os.Stderr)
	log.Info("init util success")
}
