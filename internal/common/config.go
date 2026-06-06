package common

import "github.com/caarlos0/env/v11"

type Config struct {
	Database struct {
		Path string `env:"DATABASE_PATH" envDefault:"data/app.db"`
	}
}

func NewConfig() (*Config, error) {
	conf, err := env.ParseAs[Config]()
	if err != nil {
		return nil, err
	}

	return &conf, nil
}
