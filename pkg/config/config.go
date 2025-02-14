package config

import (
    "fmt"
    "gopkg.in/yaml.v3"
    "io"
    "os"
)

type Config struct {
	AmqpConnectUrl string `yaml:"amqp_connect_url"`
	StartUrl string `yaml:"start_url"`
	Deep int `yaml:"deep"`
	LogsDir string `yaml:"logs_dir"`
}

func Load(path string) (Config, error) {
	defconf := Config{
		AmqpConnectUrl: "amqp://localhost:5672",
		Deep: 3,
		LogsDir: "logs",
	}

	file, err := os.Open(path)
	if err != nil{
		return defconf, fmt.Errorf("cant open config file: %w", err)
	}

	body, err := io.ReadAll(file)
	if err != nil{
		return defconf, fmt.Errorf("can not to read config: %w", err)
	}

	var cfg Config
	
	if err = yaml.Unmarshal(body, &cfg);err != nil{
		return defconf, fmt.Errorf("can not unmarshal bytes: %w", err)
	}

	return cfg, nil
}