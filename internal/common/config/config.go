package config

import "github.com/caarlos0/env/v11"

type Config struct {
	Database struct {
		Path string `env:"DATABASE_PATH" envDefault:"data/app.db"`
	}
	Entrypoints struct {
		AdminServer struct {
			Host string `env:"ENTRYPOINTS_ADMIN_SERVER_HOST" envDefault:"127.0.0.1"`
			Port int    `env:"ENTRYPOINTS_ADMIN_SERVER_PORT" envDefault:"4321"`
		}
		ReleasesFeedServer struct {
			Host string `env:"ENTRYPOINTS_RELEASES_FEED_SERVER_HOST" envDefault:"127.0.0.1"`
			Port int    `env:"ENTRYPOINTS_RELEASES_FEED_SERVER_PORT" envDefault:"4322"`
		}
	}
}

func New() (*Config, error) {
	conf, err := env.ParseAs[Config]()
	if err != nil {
		return nil, err
	}

	return &conf, nil
}
