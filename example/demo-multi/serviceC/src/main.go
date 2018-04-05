
/**********************************
 * Author : foo
 * Time : 2018-04-05 23:01:38
 **********************************/

package main

import (
	"fmt"
	"os"

	config "github.com/TechieYork/gofra-example/example/demo-multi/serviceC/src/config"
	application "github.com/TechieYork/gofra-example/example/demo-multi/serviceC/src/application"
)

func main() {
	// start
	fmt.Println("====== Server [serviceC] Start ======")

	// init config
	conf := config.GetConfig()

	err := conf.Init("../conf/config.json")

	if err != nil {
		fmt.Printf("Init config failed! error:%v\r\n", err.Error())
		os.Exit(-1)
	}

	// init application
	var application application.Application

	err = application.Init(conf)

	if err != nil {
		fmt.Printf("Application init failed! error:%v\r\n", err.Error())
		os.Exit(-2)
	}

	fmt.Printf("Listen on port [%v]\r\n", conf.Server.Addr)

	err = application.Run(conf.Server.Addr)

	if err != nil {
		fmt.Printf("Application run failed! error:%v\r\n", err.Error())
	}
}
