package route

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"strconv"
	"time"
	"xyf-lib/helper"
)

type Func func(g *gin.Engine)

type Route struct {
	InitRoute Func
}

type handlerFunc func(c *gin.Context, ctx context.Context) int

var HttpReqsHistory *prometheus.HistogramVec

func RunFunc(c *gin.Context, handlerFunc handlerFunc) {
	//链路追踪
	ctx, sp := helper.NewGinContextAndSpan(c)
	defer sp.Finish()
	start := time.Now()

	code := handlerFunc(c, ctx)

	//监控
	costTime := time.Since(start)
	HttpReqsHistory.WithLabelValues(c.Request.RequestURI, strconv.Itoa(code)).Observe(costTime.Seconds())
}
