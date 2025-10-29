package config

type RabbitMQClient struct {
	CancelOrder RabbitMQClientConfig `mapstructure:"cancel_order"`
}

type RabbitMQClientConfig struct {
	Queue    string `mapstructure:"queue"`
	Key      string `mapstructure:"key"`
	Exchange string `mapstructure:"exchange"`
}

type RabbitMQConfig struct {
	Host     string `mapstructure:"host"`
	VHost    string `mapstructure:"vhost"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"pass"`
}
