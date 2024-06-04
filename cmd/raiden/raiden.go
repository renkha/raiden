package main

import (
	"fmt"
	"raiden/internal/bootstrap"

	"github.com/sev-2/raiden"
)

func main() {
	// load configuration
	config, err := raiden.LoadConfig(nil)
	if err != nil {
		raiden.Panic(fmt.Sprintf("%v", err))
	}

	// Setup server
	server := raiden.NewServer(config)

	// register route
	bootstrap.RegisterRoute(server)

	// run server
	server.Run()
}
