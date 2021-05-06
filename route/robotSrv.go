package route

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/micro/go-micro/v2"
	xyfRobotSrvProto "github.com/yiqiang3344/go-lib/proto/xyf-robot-srv"
	"time"
)

type RobotSrvRoute struct {
	*Route
}

func (e *RobotSrvRoute) InitRoute(rootRoute string, g *gin.Engine, client micro.Service) {
	xyfRobotSrv := xyfRobotSrvProto.NewXyfRobotSrvService("go.micro.service.xyfRobotSrv", client.Client())
	v1 := g.Group(rootRoute)
	{
		v1.POST("/send-msg", func(c *gin.Context) {
			RunFunc(c, func(c *gin.Context, ctx context.Context) int {
				req := new(xyfRobotSrvProto.Request)
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
			RunFunc(c, func(c *gin.Context, ctx context.Context) int {
				req := new(xyfRobotSrvProto.TestRequest)
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
			RunFunc(c, func(c *gin.Context, ctx context.Context) int {
				time.Sleep(5 * time.Second)

				code := 200
				c.JSON(code, gin.H{
					"status":  "1",
					"message": "success1",
				})
				return code
			})
		})
	}
}
