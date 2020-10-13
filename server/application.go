package server

import (
	stdCtx "context"
	"github.com/go-redis/redis"
	"github.com/kataras/golog"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"sync"
	"time"
)

var (
	globalAppOnce sync.Once
	globalApp     *Application
)

// NewApp .
func NewApp() *Application {
	globalAppOnce.Do(func() {
		globalApp = new(Application)
		globalApp.WebApp = iris.New()
		globalApp.WebApp.Logger().SetTimeFormat("2006-01-02 15:04:05.000")
	})
	return globalApp
}

type Application struct {
	WebApp     *iris.Application
	Middleware []context.Handler

	Database struct {
		Debug   bool
		client  interface{}
		Install func() (db interface{})
	}

	Cache struct {
		Debug   bool
		client  redis.Cmdable
		Install func() (client redis.Cmdable)
	}
	//msgsBus     *EventBus
	//other         *other
	//Prometheus    *Prometheus
	//ControllerDep []interface{}
	//eventInfra    DomainEventInfra
	//unmarshal     func(data []byte, v interface{}) error
	//marshal       func(v interface{}) ([]byte, error)
}

func (m *Application) InstallDB(debug bool, f func() interface{}) {
	m.Database.Install = f
	m.Database.Debug = debug
}

func (m *Application) InstallRedis(debug bool, f func() (client redis.Cmdable)) {
	m.Cache.Install = f
	m.Cache.Debug = debug
}

func (m *Application) RunDb() {
	if m.Cache.Install != nil {
		m.Cache.client = m.Cache.Install()
	}
	if m.Database.Install != nil {
		m.Database.client = m.Database.Install()
	}
}

// Logger .
func (m *Application) Logger() *golog.Logger {
	return m.WebApp.Logger()
}

func (m *Application) InstallMiddleware(handler iris.Handler) {
	m.Middleware = append(m.Middleware, handler)
}

func (m *Application) RunMiddleware() {
	m.WebApp.Use(m.Middleware...)
}

func (m *Application) close(timeout int64) {
	iris.RegisterOnInterrupt(func() {
		//读取配置的关闭最长时间
		ctx, cancel := stdCtx.WithTimeout(stdCtx.Background(), time.Duration(timeout)*time.Second)
		defer cancel()
		defer func() {
			if err := recover(); err != nil {
				m.WebApp.Logger().Error(err)
			}
		}()
		//通知组件服务即将关闭
		m.WebApp.Shutdown(ctx)
	})
}
