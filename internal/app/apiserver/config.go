package apiserver

type Config struct {
	BindAddr    string      `toml:"bind_addr"`
	LogLevel    string      `toml:"log_level"`
	DatabaseURL string      `toml:"database_url"`
	TelegramBot telegramBot `toml:"telegram_bot"`
}

type telegramBot struct {
	Token   string `toml:"token"`
	Verbose bool   `toml:"verbose"`
}

func NewConfig() *Config {
	return &Config{
		BindAddr: ":3000",
		LogLevel: "info",
	}
}
