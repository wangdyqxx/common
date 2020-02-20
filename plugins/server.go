package plugins

import (
	"context"
	"github.com/wangdyqxx/util/sconfig"
	"github.com/wangdyqxx/util/slog"
)

var (
	PluginMap map[string]PluginServer
)

func init() {
	PluginMap = make(map[string]PluginServer)
	slog.Infoln(context.TODO(), "init RunPluginMap")
}

//添加运行服务需要实现的接口
type PluginServer interface {
	InitConfig(context.Context, *sconfig.StaticFileConfig) (PluginServer, error)
	Close(context.Context)
}

