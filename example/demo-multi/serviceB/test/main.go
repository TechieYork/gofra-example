
package main

import (
	"fmt"
	"time"
	"context"

    "google.golang.org/grpc"
	"google.golang.org/grpc/status"

    "github.com/grpc-ecosystem/go-grpc-middleware"

	health_check "github.com/TechieYork/gofra-example/example/demo-multi/serviceB/src/proto/health_check"

	logInterceptor "github.com/DarkMetrix/gofra/grpc-utils/interceptor/seelog_interceptor"
	monitorInterceptor "github.com/DarkMetrix/gofra/grpc-utils/interceptor/statsd_interceptor"
	tracingInterceptor "github.com/DarkMetrix/gofra/grpc-utils/interceptor/zipkin_interceptor"

    logger "github.com/DarkMetrix/gofra/common/logger/seelog"
	monitor "github.com/DarkMetrix/gofra/common/monitor/statsd"
    tracing "github.com/DarkMetrix/gofra/common/tracing/zipkin"

	pool "github.com/DarkMetrix/gofra/grpc-utils/pool"
	commonUtils "github.com/DarkMetrix/gofra/common/utils"
)

func main() {
	fmt.Println("====== Test [serviceB] begin ======")
	defer fmt.Println("====== Test [serviceB] end ======")

    // init log
    logger.Init("../conf/log.config", "serviceB_test")

	// init monitor
	monitor.Init("127.0.0.1:8125", "serviceB")

    // init tracing
    tracing.Init("http://127.0.0.1:9411/api/v1/spans", "false", ":60002", "serviceB")

	// dial remote server
	clientOpts := make([]grpc.DialOption, 0)

	clientOpts = append(clientOpts, grpc.WithUnaryInterceptor(
		grpc_middleware.ChainUnaryClient(
			tracingInterceptor.GetClientInterceptor(),
			logInterceptor.GetClientInterceptor(),
			monitorInterceptor.GetClientInterceptor())))

	pool.GetConnectionPool().Init(clientOpts, 5, 10, time.Second * 10)

	addr := commonUtils.GetRealAddrByNetwork(":60002")

	// begin test
	testHealthCheck(addr)

	time.Sleep(time.Second * 1)
}

func testHealthCheck(addr string) {
	// rpc call
	req := new(health_check.HealthCheckRequest)
	req.Message = "ping"

	for index := 0; index < 1; index++ {
		// get remote server connection
		conn, err := pool.GetConnectionPool().GetConnection(context.Background(), addr)
		defer conn.Recycle()

		// new client
		c := health_check.NewHealthCheckServiceClient(conn.Get())

		if err != nil {
			fmt.Printf("HealthCheck get connection failed! error%v", err.Error())
			continue
		}

		_, err = c.HealthCheck(context.Background(), req)

		if err != nil {
			stat, ok := status.FromError(err)

			if ok {
				fmt.Printf("HealthCheck request failed! code:%d, message:%v\r\n",
					stat.Code(), stat.Message())
			} else {
				fmt.Printf("HealthCheck request failed! err:%v\r\n", err.Error())
			}

			conn.Unhealhty()

			return
		}
	}
}
