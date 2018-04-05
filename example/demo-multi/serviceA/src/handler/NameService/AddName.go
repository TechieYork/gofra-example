
/**********************************
 * Author : foo
 * Time : 2018-04-05 23:01:38
 **********************************/

package NameService 

import (
	"context"

	//Log package
	//log "github.com/cihub/seelog"

	//Monitor package
	//monitor "github.com/DarkMetrix/gofra/common/monitor/statsd"

	//Tracing package
	//tracing "github.com/DarkMetrix/gofra/common/tracing/zipkin"

	pb "github.com/TechieYork/gofra-example/example/demo-multi/serviceA/src/proto/name"
)

func (this NameServiceImpl) AddName (ctx context.Context, req *pb.AddNameRequest) (*pb.AddNameResponse, error) {
	//Log Example:traceid must be logged
	//log.Infof("AddName begin, traceid=%v, req=%v", tracing.GetTracingId(ctx), req)

	resp := new(pb.AddNameResponse)

	return resp, nil
}
