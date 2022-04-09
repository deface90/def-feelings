package storage

type Config struct {
	Port     string `env:"BACKEND_PORT" json:"port" default:"8088"`
	Timezone string `env:"TIMEZONE" json:"timezone" default:"Asia/Yekaterinburg"`
	BaseURL  string `env:"BASE_URL" json:"base_url" default:"http://127.0.0.1:3000/"`

	TelegramBotToken string `env:"TELEGRAM_BOT_TOKEN" json:"telegram_bot_token"`

	Storage struct {
		Type     string `env:"STORAGE_TYPE" json:"type" default:"postgres"`
		Postgres struct {
			DSN string `env:"POSTGRES_DSN" json:"dsn"`
		}
	}
}
