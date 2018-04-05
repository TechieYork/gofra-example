
/**********************************
 * Author : foo
 * Time : 2018-04-05 23:01:38
 **********************************/

package application

import (
	"net"
	"time"

	"google.golang.org/grpc"
	"github.com/grpc-ecosystem/go-grpc-middleware"

	pool "github.com/DarkMetrix/gofra/grpc-utils/pool"
	logger "github.com/DarkMetrix/gofra/common/logger/seelog"
	monitor "github.com/DarkMetrix/gofra/common/monitor/statsd"
	tracing "github.com/DarkMetrix/gofra/common/tracing/zipkin"

	recoverInterceptor "github.com/DarkMetrix/gofra/grpc-utils/interceptor/recover_interceptor"
	logInterceptor "github.com/DarkMetrix/gofra/grpc-utils/interceptor/seelog_interceptor"
	monitorInterceptor "github.com/DarkMetrix/gofra/grpc-utils/interceptor/statsd_interceptor"
	tracingInterceptor "github.com/DarkMetrix/gofra/grpc-utils/interceptor/zipkin_interceptor"

	commonUtils "github.com/DarkMetrix/gofra/common/utils"

	"github.com/TechieYork/gofra-example/example/demo-multi/serviceC/src/common"
	"github.com/TechieYork/gofra-example/example/demo-multi/serviceC/src/config"

	//!!!DO NOT EDIT!!!
	health_check "github.com/TechieYork/gofra-example/example/demo-multi/serviceC/src/proto/health_check"
	email "github.com/TechieYork/gofra-example/example/demo-multi/serviceC/src/proto/email"
	/*@PROTO_STUB*/
	HealthCheckServiceHandler "github.com/TechieYork/gofra-example/example/demo-multi/serviceC/src/handler/HealthCheckService"
	EmailServiceHandler "github.com/TechieYork/gofra-example/example/demo-multi/serviceC/src/handler/EmailService"
	/*@HANDLER_STUB*/
)

type Application struct {
	ServerOpts []grpc.ServerOption
	ClientOpts []grpc.DialOption
}

//Init application
func (app *Application) Init(conf *config.Config) error {
	// process conf.Server.Addr
	conf.Server.Addr = commonUtils.GetRealAddrByNetwork(conf.Server.Addr)

	// init log
	logger.Init("../conf/log.config", common.ProjectName)

	// init monitor
	monitor.Init("127.0.0.1:8125", "serviceC", common.ProjectName)

	// init tracing
	tracing.Init("http://127.0.0.1:9411/api/v1/spans", "false", ":60003", "serviceC")

	// set server interceptor
	app.ServerOpts = append(app.ServerOpts, grpc.UnaryInterceptor(
		grpc_middleware.ChainUnaryServer(
			recoverInterceptor.GetServerInterceptor(),
			tracingInterceptor.GetServerInterceptor(),
			logInterceptor.GetServerInterceptor(),
			monitorInterceptor.GetServerInterceptor())))

	// set client interceptor
	app.ClientOpts = append(app.ClientOpts, grpc.WithUnaryInterceptor(
		grpc_middleware.ChainUnaryClient(
			tracingInterceptor.GetClientInterceptor(),
			logInterceptor.GetClientInterceptor(),
			monitorInterceptor.GetClientInterceptor())))

	err := pool.GetConnectionPool().Init(app.ClientOpts,
		conf.Client.Pool.InitConns,
		conf.Client.Pool.MaxConns,
		time.Second * time.Duration(conf.Client.Pool.IdleTime))

	if err != nil {
		return err
	}

	return nil
}

//Run application
func (app *Application) Run(address string) error {
	listen, err := net.Listen("tcp", address)

	if err != nil {
		return err
	}

	// init grpc server
	s := grpc.NewServer(app.ServerOpts ...)

	// register services
	//!!!DO NOT EDIT!!!
	health_check.RegisterHealthCheckServiceServer(s, HealthCheckServiceHandler.HealthCheckServiceImpl{})
	email.RegisterEmailServiceServer(s, EmailServiceHandler.EmailServiceImpl{})
	/*@REGISTER_STUB*/

	// run to serve
	err = s.Serve(listen)

	if err != nil {
		return err
	}

	return nil
}
