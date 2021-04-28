module xyf-gateway

go 1.13

replace xyf-lib => ../xyf-lib

require (
	github.com/garyburd/redigo v1.6.2
	github.com/gin-gonic/gin v1.7.1
	github.com/jmoiron/sqlx v1.3.3
	github.com/micro/go-micro/v2 v2.9.1
	github.com/micro/go-plugins/registry/etcdv3/v2 v2.9.1
	github.com/micro/go-plugins/wrapper/trace/opentracing/v2 v2.9.1
	github.com/opentracing/opentracing-go v1.2.0
	github.com/prometheus/client_golang v1.1.0
	xyf-lib v0.0.0-00010101000000-000000000000
)
