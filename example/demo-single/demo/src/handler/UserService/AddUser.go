
/**********************************
 * Author : techieliu
 * Time : 2018-04-05 21:42:16
 **********************************/

package UserService 

import (
	"context"

	//Log package
	//log "github.com/cihub/seelog"

	//Monitor package
	//monitor "github.com/DarkMetrix/gofra/common/monitor/statsd"

	//Tracing package
	//tracing "github.com/DarkMetrix/gofra/common/tracing/zipkin"

	pb "github.com/DarkMetrix/gofra/tmp/demo/src/proto/user"
)

func (this UserServiceImpl) AddUser (ctx context.Context, req *pb.AddUserRequest) (*pb.AddUserResponse, error) {
	//Log Example:traceid must be logged
	//log.Infof("AddUser begin, traceid=%v, req=%v", tracing.GetTracingId(ctx), req)

	resp := new(pb.AddUserResponse)

	return resp, nil
}
