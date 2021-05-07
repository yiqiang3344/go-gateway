package route

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/yiqiang3344/go-lib/utils/trace"
	"strconv"
	"time"
)

type Func func(g *gin.Engine)

type Route struct {
	InitRoute Func
}

type handlerFunc func(c *gin.Context, ctx context.Context) int

var HttpReqsHistory *prometheus.HistogramVec

func RunFunc(c *gin.Context, handlerFunc handlerFunc) {
	//链路追踪
	ctx, sp := trace.NewGinContextAndSpan(c)
	defer sp.Finish()
	start := time.Now()

	code := handlerFunc(c, ctx)

	//监控
	costTime := time.Since(start)
	HttpReqsHistory.WithLabelValues(c.Request.RequestURI, strconv.Itoa(code)).Observe(costTime.Seconds())
}
