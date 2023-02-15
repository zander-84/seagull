package project

var ConfTpl = `package config

import (
	"errors"
	"github.com/zander-84/seagull/drive/etcd"
	"github.com/zander-84/seagull/drive/goredis"
	"github.com/zander-84/seagull/drive/gormmysql"
	"github.com/zander-84/seagull/drive/mongo"
	"github.com/zander-84/seagull/think"
	"github.com/zander-84/seagull/tool"
)

var Data = new(Config)

type Config struct {
	Node      string     // 启动时候外部参数带入,并注册到注册中心，且唯一
	Mode      think.Mode // prod dev local
	Debug     bool
	WebServer struct {
		LocalIp      string ` + "`" + `default:"0.0.0.0"` + "`" + ` // 本地运行
		RemoteIp     string // 用于注册中心
		Port         int    ` + "`" + `default:"${httpPort}"` + "`" + ` // 端口
		ReadTimeout  int    ` + "`" + `default:"60"` + "`" + `
		WriteTimeout int    ` + "`" + `default:"60"` + "`" + `
	}
	GrpcServer struct {
		LocalIp      string ` + "`" + `default:"0.0.0.0"` + "`" + `// 本地运行
		RemoteIp     string // 用于注册中心
		Port         int    ` + "`" + `default:"${grpcPort}"` + "`" + ` // 端口
		ReadTimeout  int    ` + "`" + `default:"60"` + "`" + `
		WriteTimeout int    ` + "`" + `default:"60"` + "`" + `
	}

	Etcd  etcd.Conf
	Redis goredis.Conf
	Mysql gormmysql.Conf
	Mongo mongo.Conf
}

func (c *Config) Validate() error {
	if c.Node == "" {
		return errors.New("请输入node")
	}
	return nil
}

func (c *Config) GinMode() string {
	if c.Debug {
		return "debug"
	}
	return "release"
}

func (c *Config) Println() {
	if !c.Mode.IsProd() {
		tool.PrettyPrint(c)
	}
}
`
