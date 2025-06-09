package models

import (
	"github.com/teagan42/snidemind/internal/types"
)

type ServerConfig struct {
	Port types.Port `json:"port" yaml:"port" validate:"required"`
	Bind types.Host `json:"bind" yaml:"bind" validate:"required"`
}
