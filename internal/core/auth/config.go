package core_auth

import (
	"fmt"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Secret  string        `envconfig:"SECRET" required:"true"`
	Expires time.Duration `envconfig:"EXPIRES" default:"24h"`
}

func NewConfig() (Config, error) {
	var config Config
	if err := envconfig.Process("JWT", &config); err != nil {
		return Config{}, fmt.Errorf("process envconfig: %w", err)
	}
	return config, nil
}

func NewConfigMust() Config {
	config, err := NewConfig()
	if err != nil {
		panic(fmt.Errorf("get JWT config: %w", err))
	}
	return config
}
