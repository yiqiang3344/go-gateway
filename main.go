package main

import (
	"github.com/gin-gonic/gin"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-micro/v2/web"
	"github.com/micro/go-plugins/registry/etcdv3/v2"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	"github.com/yiqiang3344/go-gateway/route"
	"github.com/yiqiang3344/go-lib/utils/config"
	cLog "github.com/yiqiang3344/go-lib/utils/log"
	"github.com/yiqiang3344/go-lib/utils/monitor"
	"github.com/yiqiang3344/go-lib/utils/trace"
	"log"

	wrapperTrace "github.com/micro/go-plugins/wrapper/trace/opentracing/v2"
)

func init() {
	config.InitCfg()
	cLog.InitLogger(config.GetCfgString("project"), config.GetCfgBool("showLogToConsole", false))
	route.HttpReqsHistory = monitor.InitPrometheus()
}

func main() {
	project := config.GetCfgString("project")

	//配置网关链路追踪
	_, closer, err := trace.InitJaegerTracer(
		"go.micro.api."+project,
		config.GetCfgString("jaeger.address"),
		jaeger.SamplerTypeConst,
		1,
	)
	if err != nil {
		log.Fatal(err)
	}
	defer closer.Close()

	//初始化micro客户端
	client := micro.NewService(
		micro.WrapClient(wrapperTrace.NewClientWrapper(opentracing.GlobalTracer())),             //配置微服务客户端链路追踪
		micro.Registry(etcdv3.NewRegistry(registry.Addrs(config.GetCfgString("etcd.address")))), // 服务发现
	)
	client.Init()
	gin.SetMode(config.GetCfgString("ginMode"))
	g := gin.Default()

	//初始化路由
	initRoute(g, client)

	//初始化micro web服务，绑定端口
	service := web.NewService(
		web.Name("go.micro.api."+project),
		web.Address(":8080"),
		web.Handler(g),
	)
	_ = service.Init()
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}

func initRoute(g *gin.Engine, client micro.Service) {
	//机器人消息服务
	new(route.RobotSrvRoute).InitRoute("/robot", g, client)

	//...
}
