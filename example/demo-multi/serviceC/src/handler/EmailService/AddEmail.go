
/**********************************
 * Author : foo
 * Time : 2018-04-05 23:01:38
 **********************************/

package EmailService 

import (
	"context"

	//Log package
	//log "github.com/cihub/seelog"

	//Monitor package
	//monitor "github.com/DarkMetrix/gofra/common/monitor/statsd"

	//Tracing package
	//tracing "github.com/DarkMetrix/gofra/common/tracing/zipkin"

	pb "github.com/TechieYork/gofra-example/example/demo-multi/serviceC/src/proto/email"
)

func (this EmailServiceImpl) AddEmail (ctx context.Context, req *pb.AddEmailRequest) (*pb.AddEmailResponse, error) {
	//Log Example:traceid must be logged
	//log.Infof("AddEmail begin, traceid=%v, req=%v", tracing.GetTracingId(ctx), req)

	resp := new(pb.AddEmailResponse)

	return resp, nil
}
