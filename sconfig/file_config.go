package sconfig


type StaticFileConfig interface {
	GetServer() *ServerConfig
	GetCache() *RedisConfig
}

type YamlConfig struct {
	Server *ServerConfig `yaml:"server"`
	Redis  *RedisConfig  `yaml:"redis"`
}

func (m *YamlConfig) GetServer() *ServerConfig {
	return m.Server
}

func (m *YamlConfig) GetCache() *RedisConfig {
	return m.Redis
}

type ServerConfig struct {
	Addr     string `yaml:"addr"`
	FilePath string `yaml:"path"`
}

type RedisConfig struct {
	RedisAddr string `yaml:"addr"`
	RedisPwd  string `yaml:"pwd"`
}
