package configuration

// Config object represents a config
type Config struct {
	// LogLevel returns a string
	//
	//  NOTE: can be populated from toml (log_level), yaml (log_level) or env variable (LOG_LEVEL)
	LogLevel string `mapstructure:"log_level"`

	// DebugMode returns a bool
	//
	//  NOTE: can be populated from toml (debug_mode), yaml (debug_mode) or env variable (DEBUG_MODE)
	DebugMode bool `mapstructure:"debug_mode"`

	// Apiserver object represents an api server config
	Apiserver apiserver `mapstructure:"apiserver"`

	// AWSLambda object represents an AWS lambda config
	AWSLambda awslambda `mapstructure:"awslambda"`

	// Database object represents a database config
	Database database `mapstructure:"database"`

	// TelegramBot object represents a telegram bot config
	TelegramBot telegramBot `mapstructure:"telegram_bot"`

	// BuildDate returns a build date
	//
	// NOTE: populated only on build stage
	BuildDate string

	// Version returns a version
	//
	// NOTE: populted only on build stage
	Version string
}

type apiserver struct {
	// BindAddr returns an address string
	//
	// NOTE: can be populated from toml (apiserver.bind_addr), yaml (apiserver.bind_addr) or env variable (APISERVER_BIND_ADDR)
	BindAddr string `mapstructure:"bind_addr"`

	// Cron object represents a cron config
	Cron cron `mapstructure:"cron"`
}

type awslambda struct {
	// Enabled returns a bool
	//
	// NOTE: can be populated from toml (awslambda.enabled), yaml (awslambda.enabled) or env variable (AWSLAMBDA_ENABLED)
	Enabled bool `mapstructure:"enabled"`
}

type database struct {
	// URL returns a fully qualified url string
	//
	// NOTE: can be populated from toml (database.url), yaml (database.url) or env variable (DATABASE_URL)
	URL string `mapstructure:"url"`

	// Migrations returns a path to migrations directory
	//
	// NOTE: can be populated from toml (database.migrations), yaml (database.migrations) or env variable (DATABASE_MIGRATIONS)
	Migrations string `mapstructure:"migrations"`
}

type telegramBot struct {
	// Token returns a token for telegram bot
	//
	// NOTE: can be populated from toml (telegram_bot.token), yaml (telegram_bot.token) or env variable (TELEGRAM_BOT_TOKEN)
	Token string `mapstructure:"token"`

	// Verbos returns a bool
	//
	// NOTE: can be populated from toml (telegram_bot.verbose), yaml (telegram_bot.verbose) or env variable (TELEGRAM_BOT_VERBOSE)
	Verbose bool `mapstructure:"verbose"`
}

type cron struct {
	// Enabled returns a bool
	//
	// NOTE: can be populated from toml (apiserver.cron.enabled), yaml (apiserver.cron.enabled) or env variable (APISERVER_CRON_ENABLED)
	Enabled bool `mapstructure:"enabled"`

	// FullReport returns a bool
	//
	// NOTE: can be populated from toml (apiserver.cron.fullreport), yaml (apiserver.cron.fullreport) or env variable (APISERVER_CRON_FULLREPORT)
	Fullreport bool `mapstructure:"fullreport"`

	// Schedule returns a fully qualified string for cron
	//
	// NOTE: can be populated from toml (apiserver.cron.schedule), yaml (apiserver.cron.schedule) or env variable (APISERVER_CRON_SCHEDULE)
	Schedule string `mapstructure:"schedule"`
}

// NewConfig returns a pointer to Config object
func NewConfig() *Config {
	return &Config{}
}
