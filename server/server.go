package server

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/opentracing/opentracing-go"
	"github.com/wangdyqxx/common/backDoor"
	"github.com/wangdyqxx/common/config"
	"github.com/wangdyqxx/common/log"
	"github.com/wangdyqxx/common/processor"
)

type cmdArgs struct {
	confFile string
}

func parseFlag() (*cmdArgs, error) {
	var confFile string
	flag.StringVar(&confFile, "conf", "./main.toml", "start config")
	flag.Parse()
	return &cmdArgs{
		confFile: confFile,
	}, nil

}

var _ = processor.Processor(&Server{})

type Server struct {
	Name            string
	Conf            *config.ServerConfig
	MasterServer    processor.Processor
	SentinelServers map[string]processor.Processor
	StatLog         string
	RunStatus       bool
}

func (m *Server) Init(ctx context.Context, conf *config.ServerConfig) (interface{}, error) {
	fun := "Server.Init ->"
	fmt.Println(ctx, fun, " conf:", m.Conf)
	if m.Conf == nil || m.MasterServer == nil {
		log.Errorf(ctx, "%s check m:%+v", fun, m)
		return nil, fmt.Errorf("notFound")
	}
	args, err := parseFlag()
	if err != nil {
		fmt.Println(fun, " parse arg err:", err)
		panic(err)
	}
	m.Conf = new(config.ServerConfig)
	config.ConfigureToml(m.Conf, args.confFile, true)
	if m.Conf.BaseConfig == nil {
		return nil, errors.New("conf nil")
	}
	fmt.Println("conf3:", fmt.Sprintf("%v", &m.Conf))
	// 初始化日志
	initLog(ctx, m.Conf.LogConfig.LogDir, m.Conf.LogConfig.LogLevel)
	//初始化teace
	m.Conf.Tracer = opentracing.GlobalTracer()
	// NOTE: processor 在初始化 trace middleware 前需要保证 xtrace.GlobalTracer() 初始化完毕
	initTracer(m.Conf.BaseConfig.ServerGroup + m.Conf.BaseConfig.ServerName)
	// 初始化服务进程打点
	//stat.Init(sb.servGroup, sb.servName, "")
	//m.initMetric(sb)
	backDoorServer, err := backDoor.InitBackDoor(ctx, m.Conf, ":13141")
	if err != nil {
		log.Errorf(ctx, "%s parse arg err:%v", fun, err)
		return nil, err
	}

	_, err = m.MasterServer.Init(ctx, m.Conf)
	if err != nil {
		log.Errorf(ctx, "MasterServer start failed, server: [%s]", "master")
		return nil, err
	}
	masterConf, e := m.MasterServer.GetConfig(ctx)
	if e != nil {
		log.Errorf(ctx, "SentinelDrivers start failed, err: [%s]", e)
		return nil, err
	}
	log.Infof(ctx, "MasterServer Init success, server: [%s]", masterConf.BaseConfig.String())

	m.SentinelServers = make(map[string]processor.Processor)
	m.SentinelServers["back_door"] = backDoorServer
	for k, server := range m.SentinelServers {
		_, err = server.Init(ctx, m.Conf)
		if err != nil {
			log.Errorf(ctx, "SentinelDrivers start failed, server: [%s]", k)
			return nil, err
		}
		sentinelConf, e := server.GetConfig(ctx)
		if e != nil {
			log.Errorf(ctx, "SentinelDrivers start failed, err: [%s]", e)
			return nil, err
		}
		log.Infof(ctx, "SentinelDrivers Init success, server: [%s]", sentinelConf.BaseConfig.String())
	}
	return nil, nil
}

func (m *Server) GetConfig(ctx context.Context) (*config.ServerConfig, error) {
	fun := "Server.GetConfig ->"
	log.Infof(ctx, "%s name:%+v", fun, m.Name)
	return nil, nil
}

func (m *Server) GetDriver(ctx context.Context) interface{} {
	fun := "Server.GetDriver ->"
	log.Infof(ctx, "%s name:%+v", fun, m.Name)
	return nil
}

func (m *Server) Start(ctx context.Context, conf1 *config.ServerConfig) error {
	fun := "Server.Start ->"
	log.Infof(ctx, "%s conf:%+v", fun, m.Conf)
	err := runDriver(ctx, m.MasterServer)
	if err == nil {
		log.Infof(ctx, "MasterServer start success, server: [%s]", m.Conf.BaseConfig.String())
	} else {
		log.Errorf(ctx, "MasterServer start failed, server: [%s]", m.Conf.BaseConfig.String())
		panic(err)
	}
	for _, server := range m.SentinelServers {
		err := runDriver(ctx, server)
		if err == nil {
			log.Infof(ctx, "SentinelDrivers start success, server: [%s]", m.Conf.BaseConfig.String())
		} else {
			log.Errorf(ctx, "SentinelDrivers start failed, server: [%s]", m.Conf.BaseConfig.String())
		}
	}
	m.RunStatus = true
	return nil
}

func (m *Server) Close(ctx context.Context, conf *config.ServerConfig) error {
	fun := "Server.Close ->"
	log.Infof(ctx, "%s conf:%+v", fun, m.Conf)
	err := m.MasterServer.Close(ctx, m.Conf)
	if err != nil {
		log.Errorf(ctx, "MasterServer start failed, server: [%s]", m.Conf.BaseConfig.String())
	}
	for _, server := range m.SentinelServers {
		err := server.Close(ctx, m.Conf)
		if err != nil {
			log.Errorf(ctx, "SentinelDrivers start failed, server: [%s]", m.Conf.BaseConfig.String())
		}
	}
	CloseServer(ctx)
	m.RunStatus = false
	return nil
}

func runDriver(ctx context.Context, server processor.Processor) (err error) {
	driver := server.GetDriver(ctx)
	switch d := driver.(type) {
	case *http.Server:
		err = startHttpDriver(ctx, server, d)
	default:
		return fmt.Errorf("processor:%+v driver:%+v not recognition", server, driver)
	}
	return err
}

func (m *Server) AwaitSignal(ctx context.Context) {
	c := make(chan os.Signal, 1)
	signals := []os.Signal{
		syscall.SIGHUP,
		syscall.SIGUSR1,
		syscall.SIGUSR2,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
		syscall.SIGPIPE}
	//signal.Reset(signals...)
	signal.Notify(c, signals...)
	for s := range c {
		log.Infof(ctx, "receive a signal:%s", s.String())
		switch s {
		case syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
			fmt.Println("退出", s)
			m.Close(ctx, m.Conf)
		case syscall.SIGUSR1:
			fmt.Println("usr1", s)
		case syscall.SIGUSR2:
			fmt.Println("usr2", s)
		default:
			fmt.Println("other", s)
		}
		if !m.RunStatus {
			break
		}
	}
	os.Exit(0)
}

func CloseServer(ctx context.Context) error {
	defer log.Sync()
	return nil
}

//func StartServer(ctx context.Context, server processor.Processor, initFunc func(ctx context.Context, conf *config.ServerConfig) error) error {
//	fun := "StartServer->"
//	args, err := parseFlag()
//	if err != nil {
//		fmt.Println(fun, " parse arg err:", err)
//		panic(err)
//	}
//	conf := new(config.ServerConfig)
//	config.ConfigureToml(conf, args.confFile, true)
//	if conf.BaseConfig == nil {
//		return errors.New("conf nil")
//	}
//	// 初始化日志
//	initLog(ctx, conf.LogConfig.LogDir, conf.LogConfig.LogLevel)
//	//初始化teace
//	conf.Tracer = opentracing.GlobalTracer()
//	// NOTE: processor 在初始化 trace middleware 前需要保证 xtrace.GlobalTracer() 初始化完毕
//	initTracer(conf.BaseConfig.ServerGroup + conf.BaseConfig.ServerName)
//
//	//应用初始化
//	err = initFunc(ctx, conf)
//	if err != nil {
//		log.Panicf(ctx, "%s callInitFunc err: %v", fun, err)
//		return err
//	}
//
//	driver, err := server.Init(ctx, conf)
//	if err != nil {
//		log.Panicf(ctx, "%s initProcessor err: %v", fun, err)
//		return err
//	}
//	err = runDriver(ctx, server, driver)
//	//m.initMetric(sb)
//	if err == nil {
//		log.Infof(ctx, "server start success, server: [%s]", conf.BaseConfig.String())
//	} else {
//		log.Infof(ctx, "server start failed, server: [%s]", conf.BaseConfig.String())
//	}
//	return err
//}
