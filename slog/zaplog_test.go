package slog_test

import (
	"context"
	"github.com/opentracing/opentracing-go"
	"testing"
	"util/scontext"
	"util/slog"
	"util/strace"
)

var ctx = context.Background()

func init() {
	//先初始化trace，后初始化log
	_ = strace.InitDefaultTracer("slog test")
	tracer := opentracing.GlobalTracer()
	span := tracer.StartSpan("testlog")
	ctx = context.Background()
	ctx = opentracing.ContextWithSpan(ctx, span)
	ctx = context.WithValue(ctx, scontext.ContextKeyHead, &testHead{
		uid:     1234,
		source:  5678,
		ip:      "192.168.0.1",
		region:  "asia",
		dt:      1560499340,
		unionid: "7494ab07987ba112bd5c4f9857ccfb3f",
	})
	//dir := fileDir
	dir := ""
	slog.Init(dir, "INFO")
}

const (
	fileDir = "/Users/wangdy/go/src/util"
)

type testHead struct {
	uid     int64
	source  int32
	ip      string
	region  string
	dt      int32
	unionid string
}

func (th *testHead) ToKV() map[string]interface{} {
	return map[string]interface{}{
		"uid":     th.uid,
		"source":  th.source,
		"ip":      th.ip,
		"region":  th.region,
		"dt":      th.dt,
		"unionid": th.unionid,
	}
}

func Test_Infof(t *testing.T) {
	defer slog.Sync()
	slog.Infof(ctx, "%s", "a test log")
	slog.Errorf(ctx, "%s", "a test log")
}
