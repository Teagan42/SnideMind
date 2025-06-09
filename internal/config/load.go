package config

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
	"github.com/teagan42/snidemind/internal/models"
	"github.com/teagan42/snidemind/internal/types"
)

func LoadConfig(path string, bindAddress *types.Host, port *types.Port) (*models.Config, error) {
	v := viper.New()

	// Environment variable override
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	v.SetConfigFile(path)
	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	var cfgRaw map[string]interface{}
	if err := v.Unmarshal(&cfgRaw); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	rawBytes, err := json.Marshal(cfgRaw)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal config to JSON: %w", err)
	}
	var cfg models.Config
	if err := json.Unmarshal(rawBytes, &cfg); err != nil {
		return nil, err
	}

	validate := validator.New()
	if err := validate.Struct(cfg); err != nil {
		return nil, fmt.Errorf("validation error: %w", err)
	}

	if bindAddress != nil {
		cfg.Server.Bind = *bindAddress
	}
	if port != nil {
		cfg.Server.Port = *port
	}

	return &cfg, nil
}
