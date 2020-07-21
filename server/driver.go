package server

import (
	"context"
	"github.com/opentracing-contrib/go-stdlib/nethttp"
	"github.com/wangdyqxx/common/log"
	"github.com/wangdyqxx/common/processor"
	"github.com/wangdyqxx/common/util"
	"net"
	"net/http"
)

func startHttpDriver(ctx context.Context, pro processor.Processor, driver *http.Server) error {
	fun := "startHttpDriver ->"
	conf, err := pro.GetConfig(ctx)
	if err != nil {
		return err
	}
	paddr, err := util.GetListenAddr(driver.Addr)
	if err != nil {
		return err
	}
	tcpAddr, err := net.ResolveTCPAddr("tcp", paddr)
	if err != nil {
		return err
	}
	netListen, err := net.Listen(tcpAddr.Network(), tcpAddr.String())
	if err != nil {
		return err
	}

	// tracing
	mw := nethttp.Middleware(
		conf.GetTracer(),
		httpTrafficLogMiddleware(driver.Handler),
		nethttp.OperationNameFunc(func(r *http.Request) string {
			return "HTTP " + r.Method + ": " + r.URL.Path
		}),
		//nethttp.MWSpanFilter(tracer.UrlSpanFilter),
	)
	server := &http.Server{Handler: mw}
	go func() {
		err := server.Serve(netListen)
		if err != nil {
			log.Panicf(ctx, "%s paddr[%s]", fun, paddr)
		}
	}()
	return nil
}

func httpTrafficLogMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handler.ServeHTTP(w, r)
	})
}
