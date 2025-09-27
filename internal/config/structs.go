package config

type Config struct {
	Database DatabaseConfig `koanf:"db"`
	Broker   BrokerConfig

	Telegram TelegramConfig
	Discord  DiscordConfig
}

type DatabaseConfig struct {
	Driver string
	Url    string
}

type BrokerConfig struct {
	Addr string
}

type TelegramConfig struct {
	Token string `koanf:"token"`
}

type DiscordConfig struct {
	Token string `koanf:"token"`
}
