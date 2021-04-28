package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-micro/v2/web"
	"github.com/micro/go-plugins/registry/etcdv3/v2"
	"github.com/opentracing/opentracing-go"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
	"strconv"
	"time"
	"xyf-lib/helper"
	pb "xyf-lib/proto/xyf-robot-srv"

	wrapperTrace "github.com/micro/go-plugins/wrapper/trace/opentracing/v2"
	"github.com/prometheus/client_golang/prometheus"
)

type handlerFunc func(c *gin.Context, ctx context.Context) int

var httpReqsHistory = prometheus.NewHistogramVec(prometheus.HistogramOpts{
	Namespace:   "micro",
	Subsystem:   "",
	Name:        "srv_gateway_req_history",
	Help:        "Histogram of response latency (seconds) of http handlers.",
	ConstLabels: nil,
	Buckets:     nil,
}, []string{"method", "code"})

func runFunc(c *gin.Context, handlerFunc handlerFunc) {
	//链路追踪
	ctx, sp := helper.NewGinContextAndSpan(c)
	defer sp.Finish()
	start := time.Now()

	code := handlerFunc(c, ctx)

	//监控
	costTime := time.Since(start)
	httpReqsHistory.WithLabelValues(c.Request.RequestURI, strconv.Itoa(code)).Observe(costTime.Seconds())
}

func init() {
	//配置网关请求监控
	http.Handle("/metrics", promhttp.Handler())
	prometheus.MustRegister(
		httpReqsHistory,
	)
	go func() {
		err := http.ListenAndServe(":1010", nil)
		if err != nil {
			helper.FatalLog(err.Error(), "")
		}
	}()
}

func main() {
	//配置网关链路追踪
	_, closer, err := helper.NewJaegerTracer("go.micro.service.gateway", ":6831")
	if err != nil {
		log.Fatal(err)
	}
	defer closer.Close()

	g := gin.Default()
	service := web.NewService(
		web.Name("go.micro.api.xyfRobotSrv"),
		web.Address(":8888"),
		web.Handler(g),
	)

	//服务发现
	reg := etcdv3.NewRegistry(
		registry.Addrs("127.0.0.1:2379"),
	)

	//配置微服务客户端链路追踪
	client := micro.NewService(
		micro.WrapClient(wrapperTrace.NewClientWrapper(opentracing.GlobalTracer())),
		micro.Registry(reg),
	)
	client.Init()
	xyfRobotSrv := pb.NewXyfRobotSrvService("go.micro.service.xyfRobotSrv", client.Client())

	v1 := g.Group("/robot")
	{
		v1.POST("/send-msg", func(c *gin.Context) {
			runFunc(c, func(c *gin.Context, ctx context.Context) int {
				req := new(pb.Request)
				code := 200
				if err := c.ShouldBind(req); err != nil {
					code = 500
					c.JSON(code, gin.H{
						"status":  "500",
						"message": "参数异常:" + err.Error(),
					})
				} else if resp, err := xyfRobotSrv.SendMsg(ctx, req); err != nil {
					code = 500
					c.JSON(code, gin.H{
						"status":  "-1",
						"message": err.Error(),
					})
				} else {
					c.JSON(code, gin.H{
						"status":  resp.Status,
						"message": resp.Msg,
					})
				}
				return code
			})
		})
		v1.POST("/test", func(c *gin.Context) {
			runFunc(c, func(c *gin.Context, ctx context.Context) int {
				req := new(pb.TestRequest)
				code := 200
				if err := c.ShouldBind(req); err != nil {
					code = 500
					c.JSON(code, gin.H{
						"status":  "500",
						"message": "参数异常:" + err.Error(),
					})
				} else if resp, err := xyfRobotSrv.Test(ctx, req); err != nil {
					code = 500
					c.JSON(code, gin.H{
						"status":  "-1",
						"message": err.Error(),
					})
				} else {
					c.JSON(code, gin.H{
						"status":  resp.Status,
						"message": resp.Msg,
					})
				}
				return code
			})
		})
		v1.POST("/test1", func(c *gin.Context) {
			runFunc(c, func(c *gin.Context, ctx context.Context) int {
				code := 200
				c.JSON(code, gin.H{
					"status":  "1",
					"message": "success",
				})
				return code
			})
		})
	}
	_ = service.Init()
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
