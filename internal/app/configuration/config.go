package configuration

// Config ...
type Config struct {
	LogLevel    string      `mapstructure:"log_level"`
	Apiserver   apiserver   `mapstructure:"apiserver"`
	AWSLambda   awslambda   `mapstructure:"awslambda"`
	Database    database    `mapstructure:"database"`
	TelegramBot telegramBot `mapstructure:"telegram_bot"`
	BuildDate   string
	Version     string
}

type apiserver struct {
	BindAddr string `mapstructure:"bind_addr"`
	Cron     cron   `mapstructure:"cron"`
}

type awslambda struct {
	Enabled bool `mapstructure:"enabled"`
}

type database struct {
	URL        string `mapstructure:"url"`
	Migrations string `mapstructure:"migrations"`
}

type telegramBot struct {
	Token   string `mapstructure:"token"`
	Verbose bool   `mapstructure:"verbose"`
}

type cron struct {
	Enabled    bool   `mapstructure:"enabled"`
	Fullreport bool   `mapstructure:"fullreport"`
	Schedule   string `mapstructure:"schedule"`
}

// NewConfig ...
func NewConfig() *Config {
	return &Config{}
}
