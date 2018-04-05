
/**********************************
 * Author : foo
 * Time : 2018-04-05 23:01:38
 **********************************/

package AddrService 

import (
	"context"

	//Log package
	//log "github.com/cihub/seelog"

	//Monitor package
	//monitor "github.com/DarkMetrix/gofra/common/monitor/statsd"

	//Tracing package
	//tracing "github.com/DarkMetrix/gofra/common/tracing/zipkin"

	pb "github.com/TechieYork/gofra-example/example/demo-multi/serviceD/src/proto/addr"
)

func (this AddrServiceImpl) AddAddr (ctx context.Context, req *pb.AddAddrRequest) (*pb.AddAddrResponse, error) {
	//Log Example:traceid must be logged
	//log.Infof("AddAddr begin, traceid=%v, req=%v", tracing.GetTracingId(ctx), req)

	resp := new(pb.AddAddrResponse)

	return resp, nil
}
