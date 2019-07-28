package config

type NotificationsConfig struct {
	Service      `mapstructure:"service"`
	Webhook      `mapstructure:"webhook"`
	IFTTTWebhook `mapstructure:"ifttt_webhook"`
	SlackWebhook `mapstructure:"slack_webhook"`
	Inspec       `mapstructure:"inspec"`
	Automate     `mapstructure:"automate"`
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

type Inspec struct {
	MinImpact float32 `mapstructure:"min_impact_to_notify"`
}

type Automate struct {
	EnableForwarding bool   `mapstructure:"enable_forwarding"`
	URL              string `mapstructure:"url"`
	Token            string `mapstructure:"token"`
}
