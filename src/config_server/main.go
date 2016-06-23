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
		fmt.Printf("Usage: %s <config-file>\n", os.Args[0])
		os.Exit(1)
	}

	config, err := config.ParseConfig(os.Args[1])
	if err != nil {
		panic("Unable to parse configuration file\n" + err.Error())
	}

	store := store.CreateStore(config)
	server := server.NewServer(store)

	err = server.Start(config.Port, config.CertificateFilePath, config.PrivateKeyFilePath)
	if err != nil {
		panic("Unable to start server\n" + err.Error())
	}
}
