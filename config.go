package wonderwall

import "io"

type ServerConfig struct {
	Addr string
}

type LogConfig struct {
}

type Config struct {
	Debug  bool
	Store  StoreConfig
	Server ServerConfig
	Email  EmailConfig
}

func ReadConfig(r io.Reader) (*Config, error) {
	return nil, nil
}
