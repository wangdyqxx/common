package cmd

import (
	"context"
	"github.com/wangdyqxx/util/slog"
	"testing"
)

func TestExecCommand(t *testing.T) {
	slog.Info(ExecCommand(context.TODO(), "pwd"))
}
