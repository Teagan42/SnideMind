package models

type ServerConfig struct {
	Port int    `json:"port"`
	Bind string `json:"bind"`
}

type Config struct {
	Server ServerConfig `json:"server"`
}
