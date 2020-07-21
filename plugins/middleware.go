package plugins

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/wangdyqxx/common/log"
	"github.com/wangdyqxx/common/tracer"
	"net/http"
	"time"
)

//todo 全局中间件
func GinPublicMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		t1 := time.Now()
		defer func() {
			log.Infof(ctx, "request method:%v, uri:%v, sub time:%v", c.Request.Method, c.Request.RequestURI, time.Now().Sub(t1))
		}()
		ctx = extContextFromRequestHead(ctx, c.Request)
		c.Request = c.Request.WithContext(ctx)
		c.Request.ParseMultipartForm(32 << 20)
		c.Next()
	}
}

func GinNotFound(c *gin.Context) {
	ctx := c.Request.Context()
	log.Warnf(ctx, "load 404 page failed:%v",)
	c.Redirect(http.StatusSeeOther, "http://www.google.com/")
}

func extContextFromRequestHead(ctx context.Context, req *http.Request) context.Context {
	fun := "extContextFromRequestHead ->"
	head := &tracer.Head{}
	val := ctx.Value(tracer.ContextKeyHead)
	if val != nil {
		if vh, ok := val.(*tracer.Head); ok {
			head = vh
		}
	} else {
		tokenStr := req.Header.Get("token")
		if tokenStr != "" {
			head.Token = tokenStr
			head.Ip = req.RemoteAddr
		}
	}
	ctx = context.WithValue(ctx, tracer.ContextKeyHead, head)
	log.Infof(ctx, "%s head:%+v", fun, head)
	return ctx
}
