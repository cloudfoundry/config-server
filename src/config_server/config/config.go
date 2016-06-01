package config

import (
	"io/ioutil"
	"yaml"
)

type ServerConfig struct {
	Port int
	Database DBConfig
}

type DBConfig struct {
	Adapter string
	User string
	Password string
	Host string
	Port int
	Name string `yaml:"db_name"`
	ConnectionOptions struct {
		MaxConnections int `yaml:"max_connections"`
		PoolTimeout int	`yaml:"pool_timeout"`
	} `yaml:"connection_options"`
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

 	return config, nil;
}
