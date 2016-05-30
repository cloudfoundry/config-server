package main

import (
	"config_server/config"
	"config_server/server"
	"config_server/store"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Printf("Usage: %s <config-file>", os.Args[0])
		os.Exit(1)
	}

	serverConfig, err := config.ParseConfig(os.Args[1])
	if err != nil {
		panic("Unable to parse configuration file")
	}

	server := server.NewServer(store.NewMemoryStore())
	server.Start(serverConfig.Port)
}
