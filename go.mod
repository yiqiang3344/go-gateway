module github.com/yiqiang3344/go-gateway

go 1.14

// This can be removed once etcd becomes go gettable, version 3.4 and 3.5 is not,
// see https://github.com/etcd-io/etcd/issues/11154 and https://github.com/etcd-io/etcd/issues/11931.
replace google.golang.org/grpc => google.golang.org/grpc v1.26.0

//本地调试时可通过此方法引入本地代码
//replace github.com/yiqiang3344/go-lib => /Users/xinfei/docker/code/go/src/github.com/yiqiang3344/go-lib

require (
	github.com/gin-gonic/gin v1.7.1
	github.com/micro/go-micro/v2 v2.9.1
	github.com/micro/go-plugins/registry/etcdv3/v2 v2.9.1
	github.com/micro/go-plugins/wrapper/trace/opentracing/v2 v2.9.1
	github.com/opentracing/opentracing-go v1.2.0
	github.com/prometheus/client_golang v1.10.0
	github.com/yiqiang3344/go-lib v0.0.2-0.20210430083149-a3453ef04e62
)
