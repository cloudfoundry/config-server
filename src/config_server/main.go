package main

import (
	"config_server/config"
	"config_server/server"
	"fmt"
	"os"
)

func main() {

	if len(os.Args) != 2 {
		fmt.Printf("Usage: %s <config-file>\n", os.Args[0])
		os.Exit(1)
	}

	config, err := config.ParseConfig(os.Args[1])
	if err != nil {
		panic("Unable to parse configuration file\n" + err.Error())
	}

	server := server.NewConfigServer(config)
	err = server.Start()
	if err != nil {
		panic("Unable to start server\n" + err.Error())
	}
}
