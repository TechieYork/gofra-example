
/**********************************
 * Author : foo
 * Time : 2018-04-05 23:01:38
 **********************************/

package HealthCheckService 

import (
	"context"

	//Log package
	//log "github.com/cihub/seelog"

	//Monitor package
	//monitor "github.com/DarkMetrix/gofra/common/monitor/statsd"

	//Tracing package
	//tracing "github.com/DarkMetrix/gofra/common/tracing/zipkin"

	pb "github.com/TechieYork/gofra-example/example/demo-multi/serviceB/src/proto/health_check"
)

func (this HealthCheckServiceImpl) HealthCheck (ctx context.Context, req *pb.HealthCheckRequest) (*pb.HealthCheckResponse, error) {
	//Log Example:traceid must be logged
	//log.Infof("HealthCheck begin, traceid=%v, req=%v", tracing.GetTracingId(ctx), req)

	resp := new(pb.HealthCheckResponse)

	return resp, nil
}
