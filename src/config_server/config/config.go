package config

import (
	"io/ioutil"
	"yaml"
	"strings"
)

type ServerConfig struct {
	Port int
	Database DBConfig
}

type DBConnectionConfig struct {
	MaxOpenConnections int `yaml:"max_open_connections"`
	MaxIdleConnections int `yaml:"max_idle_connections"`
}

type DBConfig struct {
	Adapter           string
	User              string
	Password          string
	Host              string
	Port              int
	Name              string `yaml:"db_name"`
	ConnectionOptions DBConnectionConfig `yaml:"connection_options"`
}

func ParseConfig(filename string) (ServerConfig, error) {

	config := ServerConfig {}

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return config, err
	}

	err = yaml.Unmarshal([]byte(data), &config)
	if err != nil {
		return config, err
	}

	if (&config.Database != nil) && (&config.Database.Adapter != nil) {
		config.Database.Adapter = strings.ToLower(config.Database.Adapter)
	}

 	return config, nil;
}
