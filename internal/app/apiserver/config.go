package apiserver

// Config ...
type Config struct {
	BindAddr    string      `toml:"bind_addr"`
	LogLevel    string      `toml:"log_level"`
	Database    database    `toml:"database"`
	TelegramBot telegramBot `toml:"telegram_bot"`
}

type telegramBot struct {
	Token   string `toml:"token"`
	Verbose bool   `toml:"verbose"`
}

type database struct {
	URL        string `toml:"url"`
	Migrations string `toml:"migrations"`
}

// NewConfig ...
func NewConfig() *Config {
	return &Config{
		BindAddr: ":3000",
		LogLevel: "info",
	}
}
