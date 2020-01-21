package slog_test

import (
	"context"
	"testing"
	"util/slog"
)

const (
	fileDir = "/Users/wangdy/go/src/util"
)

func Test_Infof(t *testing.T) {
	//dir := fileDir
	dir := ""
	slog.Init(dir, "INFO")
	defer slog.Sync()
	slog.Infof(context.TODO(),"%s", "a test log")
	slog.Errorf(context.TODO(),"%s", "a test log")
}
