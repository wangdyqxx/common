module github.com/wangdyqxx/common

replace go.uber.org/zap => github.com/uber-go/zap v1.13.0

require (
	github.com/BurntSushi/toml v0.3.1
	github.com/bradfitz/gomemcache v0.0.0-20190913173617-a41fca850d0b // indirect
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/gin-gonic/gin v1.6.1
	github.com/go-redis/redis v6.15.9+incompatible
	github.com/gomodule/redigo v2.0.0+incompatible // indirect
	github.com/jinzhu/gorm v1.9.16
	github.com/json-iterator/go v1.1.9
	github.com/kataras/iris/v12 v12.1.8
	github.com/natefinch/lumberjack v2.0.0+incompatible
	github.com/opentracing-contrib/go-stdlib v0.0.0-20190519235532-cf7a6c988dc9
	github.com/opentracing/opentracing-go v1.1.0
	github.com/pkg/errors v0.8.1
	github.com/silenceper/wechat v2.0.1+incompatible
	github.com/uber/jaeger-client-go v2.22.1+incompatible
	github.com/uber/jaeger-lib v2.2.0+incompatible // indirect
	go.uber.org/zap v1.10.0
)

go 1.13
