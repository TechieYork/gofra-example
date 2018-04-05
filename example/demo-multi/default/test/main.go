
package main

import (
	"fmt"
	"time"
	"context"

    "google.golang.org/grpc"
	"google.golang.org/grpc/status"

	"github.com/grpc-ecosystem/go-grpc-middleware"

	user "github.com/TechieYork/gofra-example/example/demo-multi/default/src/proto/user"

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
	fmt.Println("====== Test [default] begin ======")
	defer fmt.Println("====== Test [default] end ======")

    // init log
    logger.Init("../conf/log.config", "default_test")

	// init monitor
	monitor.Init("127.0.0.1:8125", "default")

    // init tracing
    tracing.Init("http://127.0.0.1:9411/api/v1/spans", "false", "localhost:58888", "default")

	// dial remote server
	clientOpts := make([]grpc.DialOption, 0)

	clientOpts = append(clientOpts, grpc.WithUnaryInterceptor(
		grpc_middleware.ChainUnaryClient(
			tracingInterceptor.GetClientInterceptor(),
			logInterceptor.GetClientInterceptor(),
			monitorInterceptor.GetClientInterceptor())))

	pool.GetConnectionPool().Init(clientOpts, 5, 10, time.Second * 10)

	addr := commonUtils.GetRealAddrByNetwork("localhost:58888")

	// begin test
	testAddUser(addr)

	time.Sleep(time.Second * 1)
}

func testAddUser(addr string) {
	// rpc call
	req := new(user.AddUserRequest)

	     for index := 0; index < 1; index++ {
		     // get remote server connection
                conn, err := pool.GetConnectionPool().GetConnection(context.Background(),":58888")
                defer conn.Close()

                // new client
                c := user.NewUserServiceClient(conn.Get())

                if err != nil {
                        fmt.Printf("AddUser get connection failed! error%v", err.Error())
                        continue
                }

                _, err = c.AddUser(context.Background(), req)

                if err != nil {
                        stat, ok := status.FromError(err)

                        if ok {
                                fmt.Printf("AddUser request failed! code:%d, message:%v\r\n",
                                        stat.Code(), stat.Message())
                        } else {
                                fmt.Printf("AddUser request failed! err:%v\r\n", err.Error())
                        }

                        conn.Unhealhty()

                        return
                }
        }
}
