package config

type NotificationsConfig struct {
	Service `mapstructure:"service"`
	Webhook `mapstructure:"webhook"`
}

type Service struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

type Webhook struct {
	URL string `mapstructure:"url"`
}
