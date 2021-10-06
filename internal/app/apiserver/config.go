package apiserver

// Config ...
type Config struct {
	BindAddr    string      `mapstructure:"bind_addr"`
	LogLevel    string      `mapstructure:"log_level"`
	Database    database    `mapstructure:"database"`
	TelegramBot telegramBot `mapstructure:"telegram_bot"`
}

type telegramBot struct {
	Token   string `mapstructure:"token"`
	Verbose bool   `mapstructure:"verbose"`
}

type database struct {
	URL        string `mapstructure:"url"`
	Migrations string `mapstructure:"migrations"`
}

// NewConfig ...
func NewConfig() *Config {
	return &Config{
		BindAddr: ":3000",
		LogLevel: "info",
	}
}
