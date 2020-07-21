package config

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/opentracing/opentracing-go"
)

type ServerConfig struct {
	BaseConfig *BaseConfig `json:"base_config" toml:"base_config"`
	LogConfig  *LogConfig  `json:"log_config" toml:"log_config"`
	//PluginConfig *PluginConfig `json:"plugin_config"`
	//DbConfig     *DbConfig     `json:"db_config"`
	//CacheConfig  *CacheConfig  `json:"cache_config"`
	Tracer opentracing.Tracer
}

func (m *ServerConfig) GetTracer() opentracing.Tracer {
	return m.Tracer
}

type BaseConfig struct {
	ServerId    string                 `json:"server_id" toml:"server_id"`
	ServerGroup string                 `json:"server_group" toml:"server_group"`
	ServerName  string                 `json:"server_name" toml:"server_name"`
	ServerType  string                 `json:"server_type" toml:"server_type"`
	//MetaConfig  map[string]interface{} `json:"meta_config" toml:"meta_config"`
}

func (m *BaseConfig) String() string {
	return fmt.Sprintf("type:%v, group:%v, name:%v, id:%v",
		m.ServerType, m.ServerGroup, m.ServerName, m.ServerId)
}

type LogConfig struct {
	LogDir     string `json:"log_dir" toml:"log_dir"`
	LogLevel   string `json:"log_level" toml:"log_level"`
	StatLogDir string `json:"stat_log_dir" toml:"stat_log_dir"`
}

type PluginConfig struct {
	PluginName string
	MetaConfig map[string]interface{}
}
type DbConfig struct {
	DbName     string
	MetaConfig map[string]interface{}
}
type CacheConfig struct {
	CacheName  string
	MetaConfig map[string]interface{}
}

type MiddlewareConfig struct {
	MiddlewareName string
	MetaConfig     map[string]interface{}
}

// Configure .
func ConfigureToml(obj interface{}, file string, must bool) {
	if obj == nil {
		fmt.Println("obj:", obj)
		return
	}
	if file == "" {
		fmt.Println("file:", file)
		return
	}
	_, err := toml.DecodeFile(file, obj)
	if err != nil && !must {
		panic(err)
	}
	fmt.Println("obj:", fmt.Sprintf("%+v",obj))
}
