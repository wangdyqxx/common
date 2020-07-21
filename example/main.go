package main

import (
	"context"
	"fmt"
	"github.com/wangdyqxx/common/config"
	"github.com/wangdyqxx/common/server"
)

var defCtx = context.TODO()

func init() {}

func main() {
	serverLogic(defCtx)
}

func serverLogic(ctx context.Context) {
	pro := &Demo{}
	ser := new(server.Server)
	ser.Conf = &config.ServerConfig{}
	fmt.Println("conf1:", fmt.Sprintf("%v", &ser.Conf))
	pro.Conf = *(&ser.Conf)
	fmt.Println("conf1:", fmt.Sprintf("%v", &pro.Conf))
	ser.MasterServer = pro
	ser.Init(ctx, nil)
	ser.Start(ctx, nil)
	ser.AwaitSignal(ctx)
}
