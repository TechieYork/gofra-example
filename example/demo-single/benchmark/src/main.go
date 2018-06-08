
package main

import (
	"flag"
	"time"
	"context"

	log "github.com/cihub/seelog"

    "google.golang.org/grpc"
	"google.golang.org/grpc/status"

    logger "github.com/DarkMetrix/gofra/common/logger/seelog"
	monitor "github.com/DarkMetrix/gofra/common/monitor/statsd"
    tracing "github.com/DarkMetrix/gofra/common/tracing/zipkin"

    logInterceptor "github.com/DarkMetrix/gofra/grpc-utils/interceptor/seelog_interceptor"
	monitorInterceptor "github.com/DarkMetrix/gofra/grpc-utils/interceptor/statsd_interceptor"
	tracingInterceptor "github.com/DarkMetrix/gofra/grpc-utils/interceptor/zipkin_interceptor"

	health_check "github.com/TechieYork/gofra-example/example/demo-single/demo/src/proto/health_check"
	"sync"
	"github.com/grpc-ecosystem/go-grpc-middleware"
)

var (
	threads = flag.Int("threads", 1, "Go routine number to send request")
	requests = flag.Int("requests", 10000, "Request number each go routine to send")
	with_interceptor = flag.Bool("with_interceptor", false, "If gRPC DailOption is with interceptors")
)

func main() {
	flag.Parse()

	defer log.Flush()

    // init log
    err := logger.Init("../conf/log.config", "demo_test")

	if err != nil {
		log.Warnf("Init logger failed! error:%v", err.Error())
	}

	log.Info("====== Test [demo-benchmark] begin ======")
	defer log.Info("====== Test [demo-benchmark] end ======")

	// init monitor
	err = monitor.Init("127.0.0.1:8125", "demo")

	if err != nil {
		log.Warnf("Init monitor failed! error:%v", err.Error())
	}

    // init tracing
    err = tracing.Init("http://127.0.0.1:9411/api/v1/spans", "false", "localhost:58888", "demo")

	if err != nil {
		log.Warnf("Init tracing failed! error:%v", err.Error())
	}

	// dial remote server
	clientOpts := make([]grpc.DialOption, 0)

	if *with_interceptor {
		clientOpts = append(clientOpts, grpc.WithUnaryInterceptor(
		grpc_middleware.ChainUnaryClient(
			tracingInterceptor.GetClientInterceptor(),
			logInterceptor.GetClientInterceptor(),
			monitorInterceptor.GetClientInterceptor())), grpc.WithInsecure())
	} else {
		clientOpts = append(clientOpts, grpc.WithInsecure())
	}

	//test benchmark
	testBenchmark(clientOpts)
}

func testBenchmark(clientOpts []grpc.DialOption) {
	wg := &sync.WaitGroup{}
	wg.Add(*threads)

	for index := 0; index < *threads; index++ {
		// init conn
		conn, err := grpc.Dial("localhost:58888", clientOpts...)

		if err != nil {
			log.Warnf("grpc.Dial failed! error:%v", err.Error())
			return
		}

		go testHealthCheck(conn, wg)
	}

	wg.Wait()
}

func testHealthCheck(conn *grpc.ClientConn, wg *sync.WaitGroup) {
	defer wg.Done()

	//calc thread time cost
	startTime := time.Now().UnixNano()
	defer func(start int64) {
		end := time.Now().UnixNano()

		log.Infof("thread cost:%vms", (end - start) / 1000000)
	}(startTime)

	//calc per req cost average
	timeTotalNano := int64(0)

	// rpc call
	req := new(health_check.HealthCheckRequest)
	req.Message = "ping"

	for index := 0; index < *requests; index++ {
		reqStart := time.Now().UnixNano()

		c := health_check.NewHealthCheckServiceClient(conn)

		_, err := c.HealthCheck(context.Background(), req)

		if err != nil {
			stat, ok := status.FromError(err)

			if ok {
				log.Warnf("HealthCheck request failed! code:%d, message:%v",
					stat.Code(), stat.Message())
			} else {
				log.Warnf("HealthCheck request failed! err:%v", err.Error())
			}
		}

		reqEnd := time.Now().UnixNano()

		timeTotalNano += reqEnd - reqStart
	}

	log.Infof("per request cost avg:%vms", timeTotalNano / int64(*requests) / 1000000)
}
