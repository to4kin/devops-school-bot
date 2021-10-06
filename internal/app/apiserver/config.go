package apiserver

// Config ...
type Config struct {
	BindAddr    string      `mapstructure:"bind_addr"`
	LogLevel    string      `mapstructure:"log_level"`
	Database    database    `mapstructure:"database"`
	TelegramBot telegramBot `mapstructure:"telegram_bot"`
	Cron        cron        `mapstructure:"cron"`
}

type telegramBot struct {
	Token   string `mapstructure:"token"`
	Verbose bool   `mapstructure:"verbose"`
}

type database struct {
	URL        string `mapstructure:"url"`
	Migrations string `mapstructure:"migrations"`
}

type cron struct {
	Enable     bool   `mapstructure:"enable"`
	Fullreport bool   `mapstructure:"fullreport"`
	Schedule   string `mapstructure:"schedule"`
}

// NewConfig ...
func NewConfig() *Config {
	return &Config{
		BindAddr: ":3000",
		LogLevel: "info",
	}
}
