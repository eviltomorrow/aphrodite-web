package config

import (
	"fmt"

	"github.com/BurntSushi/toml"
	jsoniter "github.com/json-iterator/go"
)

// Config config
type Config struct {
	Redis  Redis  `json:"redis" toml:"redis"`
	MySQL  MySQL  `json:"mysql" toml:"mysql"`
	Log    Log    `json:"log" toml:"log"`
	System System `json:"system" toml:"system"`
}

// Redis redis
type Redis struct {
	DSN string `json:"dsn" toml:"dsn"`
}

// MySQL mysql
type MySQL struct {
	DSN     string `json:"dsn" toml:"dsn"`
	MinOpen int    `json:"min-open" toml:"min-open"`
	MaxOpen int    `json:"max-open" toml:"max-open"`
}

// System system
type System struct {
	PProfListenPort int    `json:"pprof-listen-port" toml:"pprof-listen-port"`
	HTTPServerPort  int    `json:"http-server-port" toml:"http-server-port"`
	PathHTML        string `json:"path-html" toml:"path-html"`
}

// Log 日志配置项
type Log struct {
	DisableTimestamp bool   `json:"disable-timestamp" toml:"disable-timestamp"`
	Level            string `json:"level" toml:"level"`
	Format           string `json:"format" toml:"format"`
	FileName         string `json:"filename" toml:"filename"`
	MaxSize          int    `json:"maxsize" toml:"maxsize"`
}

// Load 加载配置文件
func (cg *Config) Load(path string, f func(*Config)) error {
	_, err := toml.DecodeFile(path, cg)
	if err != nil {
		return err
	}
	f(cg)
	return nil
}

func (cg *Config) String() string {
	buf, err := jsoniter.ConfigCompatibleWithStandardLibrary.Marshal(cg)
	if err != nil {
		return fmt.Sprintf("Marshal config to json failure, nest error: %v", err)
	}
	return string(buf)
}

// DefaultGlobalConfig default global config
var DefaultGlobalConfig = &Config{
	Redis: Redis{
		DSN: "localhost:6379",
	},
	MySQL: MySQL{
		DSN:     "root:root@tcp(localhost:3306)/aphrodite?charset=utf8mb4&parseTime=true&loc=Local",
		MinOpen: 5,
		MaxOpen: 10,
	},
	Log: Log{
		DisableTimestamp: false,
		Level:            "info",
		Format:           "text",
		FileName:         "/tmp/aphrodite-web/data.log",
		MaxSize:          20,
	},
	System: System{
		PProfListenPort: 6070,
		PathHTML:        "/usr/share/aphrodite/html",
	},
}
