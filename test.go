package main

import (
	"context"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-plugins/registry/etcdv3/v2"
	"github.com/opentracing/opentracing-go"
	"github.com/yiqiang3344/go-gateway/route"
	"github.com/yiqiang3344/go-lib/helper"
	"log"
	"strconv"
	"sync"

	wrapperTrace "github.com/micro/go-plugins/wrapper/trace/opentracing/v2"
	xyfRobotSrvProto "github.com/yiqiang3344/go-lib/proto/xyf-robot-srv"
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

	//初始化路由
	for {
		initRoute1(client)
	}
}

func initRoute1(client micro.Service) {
	wg := sync.WaitGroup{}
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(n int) {
			xyfRobotSrv := xyfRobotSrvProto.NewXyfRobotSrvService("go.micro.service.xyfRobotSrv", client.Client())
			req := new(xyfRobotSrvProto.TestRequest)
			req.Test = "test" + strconv.Itoa(n)
			resp, err := xyfRobotSrv.Test(context.Background(), req)
			if err != nil {
				log.Printf("progress%d failed:%v", n, err)
			} else {
				log.Printf("progress%d success:%v", n, resp)
			}
			wg.Done()
		}(i)
	}
	wg.Wait()
}
