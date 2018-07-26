
/**********************************
 * Author : foo
 * Time : 2018-04-05 21:48:52
 **********************************/

package UserService 

import (
	"context"
	"fmt"

	//Log package
	//log "github.com/cihub/seelog"

	//Monitor package
	//monitor "github.com/DarkMetrix/gofra/common/monitor/statsd"

	//Tracing package
	//tracing "github.com/DarkMetrix/gofra/common/tracing/zipkin"

	pb "github.com/TechieYork/gofra-example/example/demo-multi/default/src/proto/user"

	pool "github.com/DarkMetrix/gofra/grpc-utils/pool"

	name "github.com/TechieYork/gofra-example/example/demo-multi/serviceA/src/proto/name"
	age "github.com/TechieYork/gofra-example/example/demo-multi/serviceB/src/proto/age"
)

func AddName(ctx context.Context) error {
	// get remote server connection
	conn, err := pool.GetConnectionPool().GetConnection(context.Background(),":60001")

	// new client
	c := name.NewNameServiceClient(conn)

	if err != nil {
		fmt.Printf("AddName get connection failed! error%v", err.Error())
		return err
	}

	req := new(name.AddNameRequest)

	_, err = c.AddName(ctx, req)

	if err != nil {
		fmt.Printf("AddName request failed! error%v", err.Error())
		return err
	}

	return nil
}

func AddAge(ctx context.Context) error {
	// get remote server connection
	conn, err := pool.GetConnectionPool().GetConnection(context.Background(),":60002")

	// new client
	c := age.NewAgeServiceClient(conn)

	if err != nil {
		fmt.Printf("AddAge get connection failed! error%v", err.Error())
		return err
	}

	req := new(age.AddAgeRequest)

	_, err = c.AddAge(ctx, req)

	if err != nil {
		fmt.Printf("AddAge request failed! error%v", err.Error())
		return err
	}

	return nil
}

func (this UserServiceImpl) AddUser (ctx context.Context, req *pb.AddUserRequest) (*pb.AddUserResponse, error) {
	//Log Example:traceid must be logged
	//log.Infof("AddUser begin, traceid=%v, req=%v", tracing.GetTracingId(ctx), req)

	resp := new(pb.AddUserResponse)

	err := AddName(ctx)

	if err != nil {
		return nil, err
	}

	err = AddAge(ctx)

	if err != nil {
		return nil, err
	}

	return resp, nil
}
