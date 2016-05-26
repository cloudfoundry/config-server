package server

import (
	"io/ioutil"
	"yaml"
)

type ServerConfig struct {
	Port     int
	Database struct {
		Adapter           string
		User              string
		Password          string
		Host              string
		Port              int
		database          string
		ConnectionOptions struct {
			MaxConnections int `yaml:"max_connections"`
			PoolTimeout    int `yaml:"pool_timeout"`
		} `yaml:"connection_options"`
	}
}

func ParseConfig(filepath string) (ServerConfig, error) {
	config := ServerConfig{}

	data, err := ioutil.ReadFile(filepath)
	if err != nil {
		return config, err
	}

	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return config, err
	}

	return config, nil
}
