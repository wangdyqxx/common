package processor

import (
	"context"
	"github.com/wangdyqxx/common/config"
)

//程序接口
type Processor interface {
	Init(ctx context.Context, conf *config.ServerConfig) (interface{}, error)
	GetConfig(ctx context.Context) (*config.ServerConfig, error)
	GetDriver(ctx context.Context) interface{}
	Start(ctx context.Context, conf *config.ServerConfig) error
	Close(ctx context.Context, conf *config.ServerConfig) error
}
