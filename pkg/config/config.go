package config

type NotificationsConfig struct {
	Service      `mapstructure:"service"`
	Webhook      `mapstructure:"webhook"`
	IFTTTWebhook `mapstructure:"ifttt_webhook"`
	SlackWebhook `mapstructure:"slack_webhook"`
}

type Service struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

type Webhook struct {
	URL string `mapstructure:"url"`
}

type IFTTTWebhook struct {
	URL string `mapstructure:"url"`
}

type SlackWebhook struct {
	URL string `mapstructure:"url"`
}
