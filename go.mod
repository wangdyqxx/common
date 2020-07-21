module github.com/wangdyqxx/common

replace go.uber.org/zap => github.com/uber-go/zap v1.13.0

require (
	github.com/BurntSushi/toml v0.3.1
	github.com/gin-gonic/gin v1.6.1
	github.com/natefinch/lumberjack v2.0.0+incompatible
	github.com/opentracing-contrib/go-stdlib v0.0.0-20190519235532-cf7a6c988dc9
	github.com/opentracing/opentracing-go v1.1.0
	github.com/pkg/errors v0.8.1
	github.com/shawnfeng/sutil v1.3.36
	github.com/uber/jaeger-client-go v2.22.1+incompatible
	go.uber.org/zap v1.10.0
)

go 1.13
