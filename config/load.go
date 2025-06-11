package config

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
	"go.uber.org/fx"
)

type Params struct {
	fx.In
	ConfigPath  string `name:"configPath"`
	BindAddress *Host  `name:"bindAddress"`
	Port        *Port  `name:"port"`
}

type Result struct {
	fx.Out
	Config *Config
}

func LoadConfig(p Params) (Result, error) {
	result := Result{}
	v := viper.New()

	// Environment variable override
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	v.SetConfigFile(p.ConfigPath)
	if err := v.ReadInConfig(); err != nil {
		return result, fmt.Errorf("failed to read config: %w", err)
	}

	var cfgRaw map[string]interface{}
	if err := v.Unmarshal(&cfgRaw); err != nil {
		return result, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	rawBytes, err := json.Marshal(cfgRaw)
	if err != nil {
		return result, fmt.Errorf("failed to marshal config to JSON: %w", err)
	}
	var cfg Config
	if err := json.Unmarshal(rawBytes, &cfg); err != nil {
		return result, err
	}

	validate := validator.New()
	if err := validate.Struct(cfg); err != nil {
		return result, fmt.Errorf("validation error: %w", err)
	}

	if p.BindAddress != nil {
		cfg.Server.Bind = *p.BindAddress
	}
	if p.Port != nil {
		cfg.Server.Port = *p.Port
	}

	return Result{
		Config: &cfg,
	}, nil
}
