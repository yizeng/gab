package config

import (
	"context"

	"github.com/sethvargo/go-envconfig"
)

func Load() (*AppConfig, error) {
	var conf AppConfig
	if err := envconfig.Process(context.TODO(), &conf); err != nil {
		return nil, err
	}

	return &conf, nil
}
