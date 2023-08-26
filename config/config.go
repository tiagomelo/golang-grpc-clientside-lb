package config

import (
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
)

// Config holds all configuration needed by this app.
type Config struct {
	ServiceTargetAddress        string `envconfig:"SERVICE_TARGET_ADDRESS" required:"true"`
	PromTemplateFile            string `envconfig:"PROM_TEMPLATE_FILE" required:"true"`
	PromOutputFile              string `envconfig:"PROM_OUTPUT_FILE" required:"true"`
	PromTargetGrpcServerOnePort int    `envconfig:"PROM_TARGET_GRPC_SERVER_ONE_PORT" required:"true"`
	PromTargetGrpcServerTwoPort int    `envconfig:"PROM_TARGET_GRPC_SERVER_TWO_PORT" required:"true"`
	DsTemplateFile              string `envconfig:"DS_TEMPLATE_FILE" required:"true"`
	DsOutputFile                string `envconfig:"DS_OUTPUT_FILE" required:"true"`
	DsServerPort                int    `envconfig:"DS_SERVER_PORT" required:"true"`
}

// Read reads the environment variables from the given file and returns a Config.
func Read() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, errors.Wrap(err, "loading env vars")
	}
	config := new(Config)
	if err := envconfig.Process("", config); err != nil {
		return nil, errors.Wrap(err, "processing env vars")
	}
	return config, nil
}
