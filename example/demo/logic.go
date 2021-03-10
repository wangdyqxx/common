package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/wangdyqxx/common/config"
	"github.com/wangdyqxx/common/log"
	"github.com/wangdyqxx/common/plugins"
	"github.com/wangdyqxx/common/processor"
	"net/http"
	"time"
)

func initLogic(ctx context.Context, conf *config.ServerConfig) error {
	fun := "initLogic ->"
	log.Infof(ctx, "%s Conf:%+v", fun, conf)
	return nil
}

var _ = processor.Processor(&Demo{})

type Demo struct {
	Name    string
	Conf    *config.ServerConfig
	driver  *http.Server
	StatLog string
}

func (m *Demo) Init(ctx context.Context, conf *config.ServerConfig) (interface{}, error) {
	fun := "Demo.Init ->"
	m.Conf = conf
	log.Infof(ctx, "%s Conf:%+v", fun, m.Conf)
	if m.Conf == nil {
		log.Errorf(ctx, "%s Conf nil:%+v", fun, m.Conf)
		return nil, fmt.Errorf(m.Name,"notFound")
	}
	engine := gin.New()
	//m.StatLog = m.Conf.LogConfig.StatLogDir+"/stat.log"
	//ginLog, err := os.Create(m.StatLog)
	//if err != nil {
	//	panic(err)
	//}
	//gin.DefaultWriter = io.MultiWriter(ginLog)
	loadEngine(engine)
	m.driver = &http.Server{
		Addr:           ":13140",
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		IdleTimeout:    10 * time.Second,
		MaxHeaderBytes: 1024 * 10,
		Handler:        engine,
	}
	return m.driver, nil
}

func (m *Demo) GetConfig(ctx context.Context) (conf *config.ServerConfig, err error) {
	fun := "Demo.GetConfig ->"
	log.Infof(ctx, "%s name:%+v", fun, m.Name)
	if m.Conf == nil {
		log.Errorf(ctx, "%s Conf nil:%+v", fun, m)
		return nil, fmt.Errorf("notFound")
	}
	conf = &(*m.Conf)
	return conf, nil
}

func (m *Demo) GetDriver(ctx context.Context) interface{} {
	return m.driver
}

func (m *Demo) Start(ctx context.Context, conf *config.ServerConfig) error {
	//fun := "Demo.Start ->"
	//log.Infof(ctx, "%s Conf:%+v", fun, conf)
	//if m.driver == nil {
	//	log.Errorf(ctx, "%s server nil Conf:%+v", fun, conf)
	//	panic("server is nil")
	//}
	//go func() {
	//	err := m.driver.ListenAndServe()
	//	if err != nil && err != http.ErrServerClosed {
	//		log.Errorf(ctx, "%s err:%v", fun, err)
	//		return
	//	}
	//}()
	//log.Infof(ctx, "%s init server success: http://localhost%s/tool", fun, m.driver.Addr)
	return nil
}

func (m *Demo) Close(ctx context.Context, conf *config.ServerConfig) error {
	fun := "Demo.Close ->"
	log.Infof(ctx, "%s Conf:%+v", fun, conf)
	return m.driver.Shutdown(ctx)
}

func loadEngine(engine *gin.Engine) {
	engine.Use(
		gin.Logger(),
		plugins.GinPublicMiddleware(),
		gin.Recovery(),
	)
	engine.NoRoute(func(c *gin.Context) {
		plugins.GinNotFound(c)
	})
	engine.GET("/", Index)
}

func Index(c *gin.Context) {
	_, _ = c.Writer.Write([]byte(c.Request.RequestURI))
}