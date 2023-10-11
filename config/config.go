package config

import "github.com/joeshaw/envdecode"

type Config struct {
	Addr         string `env:"ADDR,default=:8000"`
	OtelEndpoint string `env:"OTEL_ENDPOINT,default=0.0.0.0:4317"`
	ServiceName  string `env:"SERVICE_NAME,default=ExampleService"`
	Environment  string `env:"ENVIRONMENT,default=staging"`
}

func Load() (Config, error) {
	var cfg Config
	if err := envdecode.Decode(&cfg); err != nil {
		return Config{}, err
	}

	return cfg, nil
}
