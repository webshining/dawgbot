package config

type Config struct {
	Telegram TelegramConfig
	Discord  DiscordConfig
	Database DatabaseConfig `koanf:"db"`
}

type TelegramConfig struct {
	Token string
}

type DiscordConfig struct {
	Token string
}

type DatabaseConfig struct {
	Driver string
	Url    string
}
