package config

import (
	"encoding/json"
	"io"
	"os"
)

type Config struct {
	Env             string          `json:"env"`
	HttpOrderServer HttpOrderServer `json:"http_order_server"`
	GrpcOrderServer GrpcOrderServer `json:"grpc_order"`
	Mongo           MongoDB         `json:"mongo"`
}

type HttpOrderServer struct {
	Port int `json:"port"`
}

type MongoDB struct {
	Uri string `json:"uri"`
}

type GrpcOrderServer struct {
	Port int `json:"port"`
}

func MustRead(configPath string) *Config {
	file, err := os.Open(configPath)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err = file.Close(); err != nil {
			panic(err)
		}
	}()

	b, err := io.ReadAll(file)
	if err != nil {
		panic(err)
	}

	var config Config
	if err = json.Unmarshal(b, &config); err != nil {
		panic(err)
	}

	return &config
}
