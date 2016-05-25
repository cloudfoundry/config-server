package main

import (
	"config_server/server"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Printf("Usage: %s <config-file>", os.Args[0])
		os.Exit(1)
	}

	serverConfig, err := server.ParseConfig(os.Args[1])
	if err != nil {
		panic("Unable to parse configuration file")
	}

	server.StartServer(serverConfig.Port)
}
