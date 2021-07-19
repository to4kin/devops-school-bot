package apiserver

import "gitlab.devops.telekom.de/anton.bastin/devops-school-bot/internal/app/store"

// Config ...
type Config struct {
	BindAddr string `toml:"bind_addr"`
	LogLevel string `toml:"log_level"`
	Store    *store.Config
}

// NewConfig ...
func NewConfig() *Config {
	return &Config{
		BindAddr: ":3000",
		LogLevel: "info",
		Store:    store.NewConfig(),
	}
}
