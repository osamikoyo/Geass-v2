package config

type Config struct {
	AmqpConnectUrl string `yaml:"amqp_connect_url"`
	Deep uint8 `yaml:"deep"`
	LogsDir string `yaml:"logs_dir"`
}

func Load(path string) (Config, error) {
	
}