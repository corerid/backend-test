package main

import (
	"github.com/corerid/backend-test/config"
	"github.com/corerid/backend-test/db"
	"github.com/corerid/backend-test/router"
)

func main() {

	config.Init()

	dbCon := db.New(config.V())

	handler := injection(dbCon, config.V())

	r := router.AddRoute(handler)

	go func() {
		err := r.Run()
		if err != nil {
			panic(err)
		}
	}()

	go handler.MonitorBlockEthereumHandler(config.V().MonitoredAddress)

	select {}
}
