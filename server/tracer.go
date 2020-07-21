package server

import (
	"context"
	"github.com/wangdyqxx/common/log"
	"github.com/wangdyqxx/common/tracer"
)

func initTracer(name string) error {
	fun := "Server.initTracer -->"
	ctx := context.Background()
	err := tracer.InitDefaultTracer(name)
	if err != nil {
		log.Errorf(ctx, "%s init tracer err: %v", fun, err)
	}
	return err
}