
/**********************************
 * Author : foo
 * Time : 2018-04-05 21:57:29
 **********************************/

package AgeService 

import (
	"context"
	"fmt"

	//Log package
	//log "github.com/cihub/seelog"

	//Monitor package
	//monitor "github.com/DarkMetrix/gofra/common/monitor/statsd"

	//Tracing package
	//tracing "github.com/DarkMetrix/gofra/common/tracing/zipkin"

	pb "github.com/TechieYork/gofra-example/example/demo-multi/serviceB/src/proto/age"

	pool "github.com/DarkMetrix/gofra/grpc-utils/pool"

	email "github.com/TechieYork/gofra-example/example/demo-multi/serviceC/src/proto/email"
	addr "github.com/TechieYork/gofra-example/example/demo-multi/serviceD/src/proto/addr"
)

func AddEmail(ctx context.Context) error {
	// get remote server connection
	conn, err := pool.GetConnectionPool().GetConnection(context.Background(),":60003")
	defer conn.Close()

	// new client
	c := email.NewEmailServiceClient(conn.Get())

	if err != nil {
		fmt.Printf("AddEmail get connection failed! error%v", err.Error())
		return err
	}

	req := new(email.AddEmailRequest)

	_, err = c.AddEmail(ctx, req)

	if err != nil {
		fmt.Printf("AddEmail request failed! error%v", err.Error())
		conn.Unhealthy()
		return err
	}

	return nil
}

func AddAddr(ctx context.Context) error {
	// get remote server connection
	conn, err := pool.GetConnectionPool().GetConnection(context.Background(),":60004")
	defer conn.Close()

	// new client
	c := addr.NewAddrServiceClient(conn.Get())

	if err != nil {
		fmt.Printf("AddAddr get connection failed! error%v", err.Error())
		return err
	}

	req := new(addr.AddAddrRequest)

	_, err = c.AddAddr(ctx, req)

	if err != nil {
		fmt.Printf("AddAddr request failed! error%v", err.Error())
		conn.Unhealthy()
		return err
	}

	return nil
}

func (this AgeServiceImpl) AddAge (ctx context.Context, req *pb.AddAgeRequest) (*pb.AddAgeResponse, error) {
	//Log Example:traceid must be logged
	//log.Infof("AddAge begin, traceid=%v, req=%v", tracing.GetTracingId(ctx), req)

	resp := new(pb.AddAgeResponse)

	err := AddEmail(ctx)

	if err != nil {
		return nil, err
	}

	err = AddAddr(ctx)

	if err != nil {
		return nil, err
	}

	return resp, nil
}
