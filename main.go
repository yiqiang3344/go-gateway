package main

import (
	"github.com/gin-gonic/gin"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-micro/v2/web"
	"github.com/micro/go-plugins/registry/etcdv3/v2"
	"github.com/opentracing/opentracing-go"
	"github.com/yiqiang3344/go-gateway/route"
	"github.com/yiqiang3344/go-lib/helper"
	"log"

	wrapperTrace "github.com/micro/go-plugins/wrapper/trace/opentracing/v2"
)

func init() {
	helper.InitCfg()
	helper.InitLogger()
	route.HttpReqsHistory = helper.InitPrometheus()
}

func main() {
	project := helper.GetCfgString("project")

	//配置网关链路追踪
	_, closer, err := helper.InitJaegerTracer("go.micro.api."+project, helper.GetCfgString("jaeger.address"))
	if err != nil {
		log.Fatal(err)
	}
	defer closer.Close()

	//初始化micro客户端
	client := micro.NewService(
		micro.WrapClient(wrapperTrace.NewClientWrapper(opentracing.GlobalTracer())),             //配置微服务客户端链路追踪
		micro.Registry(etcdv3.NewRegistry(registry.Addrs(helper.GetCfgString("etcd.address")))), // 服务发现
	)
	client.Init()
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
	new(route.XyfRobotSrvRoute).InitRoute("/robot", g, client)

	//...
}
