package cache

import (
	"context"
	"encoding/json"
	"github.com/go-redis/redis"
	"github.com/wangdyqxx/common/config"
	"github.com/wangdyqxx/common/log"
	"github.com/wangdyqxx/common/util"
	"time"
)

type RedisOptions struct {
	Addr         string        `json:"addr"`
	Network      string        `json:"network"`
	Password     string        `json:"password"`
	DialTimeout  time.Duration `json:"dial_timeout"`
	ReadTimeout  time.Duration `json:"read_timeout"`
	WriteTimeout time.Duration `json:"write_timeout"`
	PoolSize     int           `json:"pool_size"`
	PoolTimeout  time.Duration `json:"pool_timeout"`

	DB                 int           `json:"db"`
	MaxRetries         int           `json:"max_retries"`
	MinRetryBackoff    time.Duration `json:"min_retry_backoff"`
	MaxRetryBackoff    time.Duration `json:"max_retry_backoff"`
	MinIdleConns       int           `json:"min_idle_conns"`
	MaxConnAge         time.Duration `json:"max_conn_age"`
	IdleTimeout        time.Duration `json:"idle_timeout"`
	IdleCheckFrequency time.Duration `json:"idle_check_frequency"`
}

func initRedis(ctx context.Context, config *config.CacheConfig) error {
	fun := "initRedis->"
	if config == nil {
		return util.ErrorNil
	}
	redisOptions := new(RedisOptions)
	bs, _ := json.Marshal(config.MetaConfig)
	json.Unmarshal(bs, redisOptions)
	client := redis.NewClient(&redis.Options{
		Addr:         redisOptions.Addr,
		DialTimeout:  redisOptions.DialTimeout,
		ReadTimeout:  redisOptions.ReadTimeout,
		WriteTimeout: redisOptions.WriteTimeout,
		PoolSize:     redisOptions.PoolSize,
		PoolTimeout:  redisOptions.PoolTimeout,
	})
	pong, err := client.Ping().Result()
	if err != nil {
		log.Errorf(ctx, "%s ping:%s err:%s", fun, pong, err)
		return err
	}
	return nil
}
