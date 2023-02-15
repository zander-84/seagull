package project

var MainTpl = `package main

import (
	"github.com/gin-gonic/gin"
	"github.com/zander-84/seagull/app"
	"github.com/zander-84/seagull/contrib/endpoint/grpc_router"
	"github.com/zander-84/seagull/contrib/endpoint/http_gin"
	"github.com/zander-84/seagull/endpoint"
	"github.com/zander-84/seagull/transport/grpc"
	"github.com/zander-84/seagull/transport/http"
	"${project}/apps/${server}/internal/config"
	"${project}/apps/${server}/internal/transport"
	"log"
	"time"
)

func main() {
	bs := boot()
	bs.Init()

	rmc := transport.NewRmc(config.Data.Mode)
	httpProxy := http_gin.NewRouter(gin.New())
	rmc.Proxy(httpProxy.Endpoint, endpoint.Http)

	grpcProxy := grpc_router.NewRouter(rmc)
	rmc.Proxy(grpcProxy.Endpoint, endpoint.Grpc)

	hs := http.NewServer("http-order", config.Data.WebServer.LocalIp, config.Data.WebServer.RemoteIp, config.Data.WebServer.Port, http.ServerHandler(httpProxy))
	gs := grpc.NewServer("grpc-order", config.Data.GrpcServer.LocalIp, config.Data.GrpcServer.RemoteIp, config.Data.GrpcServer.Port, grpc.ServerHandler(grpcProxy.ServerHandler()))

	apps := app.New(
		app.Name("${server}-api"),
		app.Version("v1.0.0"),
		app.Server(hs, gs),
		app.Bs(bs),
		app.RegistrarTimeout(time.Second*10),
	//	app.Registrar(&registry.Registry{Engine: etcd.New(pkg.Etcd.Engine())}),
	)

	if err := apps.Run(); err != nil {
		log.Fatal(err)
	}
}

`

var BsTpl = `package main

import (
	"github.com/jinzhu/configor"
	"github.com/zander-84/seagull/app"
	"github.com/zander-84/seagull/think"
	"${project}/apps/${server}/internal/config"
	"${project}/apps/${server}/internal/pkg"
	"${project}/apps/${server}/internal/usecase"
	"log"
)

func boot() *app.Bootstrap {
	bs := app.NewBootstrap()
	bs.RegisterInitEvents(0,
		app.NewEvent("load", func() error {
			think.Bootstrap()
			if err := configor.Load(config.Data, "../config/config.yml"); err != nil {
				log.Fatal("配置文件错误" + err.Error())
			}
			pkg.InitLib()
			usecase.InitUseCase()
			return nil
		}),
	)
	bs.RegisterAfterStopEvents(0, app.NewEvent("ProcessorCancel", func() error {
		return pkg.ProcessorCancel()
	}))
	return bs
}

`
