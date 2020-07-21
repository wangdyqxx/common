package backDoor

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/opentracing-contrib/go-stdlib/nethttp"
	"github.com/wangdyqxx/common/config"
	"github.com/wangdyqxx/common/log"
	"github.com/wangdyqxx/common/processor"
	"github.com/wangdyqxx/common/util"
	"net/http"
	"os"
	"time"
)

var DefBackDoor *BackDoor

func InitBackDoor(ctx context.Context, conf *config.ServerConfig, addr string) (*BackDoor, error) {
	DefBackDoor = new(BackDoor)
	DefBackDoor.conf = conf
	DefBackDoor.addr = addr
	_, err := DefBackDoor.Init(ctx, conf)
	return DefBackDoor, err
}

var _ = processor.Processor(&BackDoor{})

type BackDoor struct {
	addr string
	conf *config.ServerConfig
	driver *http.Server
	serviceMD5  string
	startUpTime string
}

func (m *BackDoor) Init(ctx context.Context, conf *config.ServerConfig) (interface{}, error) {
	m.conf = conf
	filePath, err := os.Executable()
	if err == nil {
		md5, err := util.MD5Sum(filePath)
		if err == nil {
			m.serviceMD5 = fmt.Sprintf("%x", md5)
		}
	}
	m.startUpTime = time.Now().Format("2006-01-02 15:04:05")
	router := m.initRouter()
	// tracing
	mw := nethttp.Middleware(
		m.conf.GetTracer(),
		httpTrafficLogMiddleware(router),
		nethttp.OperationNameFunc(func(r *http.Request) string {
			return "HTTP " + r.Method + ": " + r.URL.Path
		}),
		//nethttp.MWSpanFilter(tracer.UrlSpanFilter),
	)
	m.driver = &http.Server{Handler: mw}
	return m.driver, nil
}

func (m *BackDoor) GetConfig(ctx context.Context) (conf *config.ServerConfig, err error) {
	conf = &(*m.conf)
	return conf, nil
}

func (m *BackDoor) GetDriver(ctx context.Context) interface{} {
	return m.driver
}

func (m *BackDoor) initRouter() (router *gin.Engine) {
	router = gin.New()
	router.Use(gin.Recovery())
	router.POST("/", index)
	// admin
	admin := router.Group("/admin")
	admin.POST("/ping", ping)
	return
}

func (m *BackDoor) Start(ctx context.Context, conf *config.ServerConfig) error {
	//fun := "BackDoor.Start->"
	//paddr, err := util.GetListenAddr(m.addr)
	//if err != nil {
	//	return err
	//}
	//tcpAddr, err := net.ResolveTCPAddr("tcp", paddr)
	//if err != nil {
	//	return err
	//}
	//netListen, err := net.Listen(tcpAddr.Network(), tcpAddr.String())
	//if err != nil {
	//	return err
	//}
	//go func() {
	//	err := m.server.Serve(netListen)
	//	if err != nil {
	//		log.Panicf(ctx, "%s paddr[%s]", fun, paddr)
	//	}
	//}()
	return nil
}

func (m *BackDoor) Close(ctx context.Context, conf *config.ServerConfig) error {
	fun := "BackDoor.Close->"
	log.Infof(ctx, "%s ", fun)
	return m.driver.Shutdown(ctx)
}


func httpTrafficLogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
	})
}

func index(c *gin.Context) {
	log.Info(c.Request.URL.String())
	c.JSON(http.StatusOK, c.Request.URL.String())
}

func ping(c *gin.Context) {
	log.Info(c.Request.ContentLength, " ping: ", c.Request.URL.String())
	c.JSON(http.StatusOK, c.Request.URL.String())
}
