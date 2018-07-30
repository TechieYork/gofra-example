
package main

import (
	"flag"
	"time"
	"context"
	"sync"

	log "github.com/cihub/seelog"

    "google.golang.org/grpc"
	"google.golang.org/grpc/status"

	"github.com/grpc-ecosystem/go-grpc-middleware"

    logger "github.com/DarkMetrix/gofra/common/logger/seelog"
	monitor "github.com/DarkMetrix/gofra/common/monitor/statsd"
    tracing "github.com/DarkMetrix/gofra/common/tracing/jaeger"
    pool "github.com/DarkMetrix/gofra/grpc-utils/pool"

    logInterceptor "github.com/DarkMetrix/gofra/grpc-utils/interceptor/seelog_interceptor"
	monitorInterceptor "github.com/DarkMetrix/gofra/grpc-utils/interceptor/statsd_interceptor"
	tracingInterceptor "github.com/DarkMetrix/gofra/grpc-utils/interceptor/opentracing_interceptor"

	config "github.com/TechieYork/gofra-example/example/demo-single/benchmark/src/config"
	health_check "github.com/TechieYork/gofra-example/example/demo-single/demo/src/proto/health_check"
)

var (
	addr = flag.String("addr", "localhost:58888", "Service's address, default is localhost:58888")
	threads = flag.Int("threads", 1, "Go routine number to send request, default is 1")
	requests = flag.Int("requests", 10000, "Request number each go routine to send, default is 10000")
	with_interceptor = flag.Bool("with_interceptor", false, "If gRPC DailOption is with interceptors, default is false")
)

func main() {
	flag.Parse()

	defer log.Flush()

    // init log
    err := logger.Init("../conf/log.config", "demo_benchmark")

	if err != nil {
		log.Warnf("Init logger failed! error:%v", err.Error())
		return
	}

	log.Info("====== Test [demo-benchmark] begin ======")
	defer log.Info("====== Test [demo-benchmark] end ======")

	// init config
	conf := config.GetConfig()

	err = conf.Init("../conf/config.toml")

	if err != nil {
		panic("Init config failed! error:%v",)
		log.Warnf("Init config failed! error:%v", err.Error())
		return
	}

	// init monitor
	err = monitor.Init(conf.Monitor.Params...)

	if err != nil {
		log.Warnf("Init monitor failed! error:%v", err.Error())
	}

    // init tracing
    err = tracing.Init(conf.Tracing.Params...)

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

	err = pool.GetConnectionPool().Init(clientOpts)

	if err != nil {
		log.Warnf("Init pool failed! error:%v", err.Error())
		return
	}

	//test benchmark
	testBenchmark()
}

func testBenchmark() {
	wg := &sync.WaitGroup{}
	wg.Add(*threads)

	for index := 0; index < *threads; index++ {
		go testHealthCheck(wg)
	}

	wg.Wait()
}

func testHealthCheck(wg *sync.WaitGroup) {
	defer wg.Done()

	//calc thread time cost
	startTime := time.Now().UnixNano()
	defer func(start int64) {
		end := time.Now().UnixNano()

		log.Infof("thread cost:%vms", (end - start) / 1000000)
	}(startTime)

	//calc per req cost average
	timeTotalNano := int64(0)
	timeMaxNano := int64(0)
	timeMinNano := int64(0xFFFFFFFF)

	// rpc call
	req := new(health_check.HealthCheckRequest)
	req.Message = "ping"

	for index := 0; index < *requests; index++ {
		reqStart := time.Now().UnixNano()

		conn, err := pool.GetConnectionPool().GetConnection(context.Background(), *addr)

		if err != nil {
			log.Warnf("pool.GetConnection failed! error:%v", err.Error())
			continue
		}

		c := health_check.NewHealthCheckServiceClient(conn)

		_, err = c.HealthCheck(context.Background(), req)

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

		timeCost := reqEnd - reqStart

		timeTotalNano += timeCost

		if timeCost > timeMaxNano {
			timeMaxNano = timeCost
		}

		if timeCost < timeMinNano {
			timeMinNano = timeCost
		}
	}

	log.Infof("per request cost avg:%vms, min:%vms, max:%vms",
		timeTotalNano / int64(*requests) / 1000000, timeMinNano / 1000000, timeMaxNano / 1000000)
}
