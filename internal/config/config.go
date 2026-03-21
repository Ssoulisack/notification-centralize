package config

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	API      APIConfig      `mapstructure:"api"`
	Database DatabaseConfig `mapstructure:"database"`
	Redis    RedisConfig    `mapstructure:"redis"`
	Queue    QueueConfig    `mapstructure:"queue"`
	Worker   WorkerConfig   `mapstructure:"worker"`
	Provider ProviderConfig `mapstructure:"providers"`
	Logging  LoggingConfig  `mapstructure:"logging"`
}

type ServerConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
	Mode string `mapstructure:"mode"`
}

type APIConfig struct {
	Key       string `mapstructure:"key"`
	RateLimit int    `mapstructure:"rate_limit"`
}

type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Name     string `mapstructure:"name"`
	MaxConns int    `mapstructure:"max_conns"`
	MinConns int    `mapstructure:"min_conns"`
}

type RedisConfig struct {
	Addr     string `mapstructure:"addr"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

type QueueConfig struct {
	Engine     string         `mapstructure:"engine"`
	NATS       NATSConfig     `mapstructure:"nats"`
	RabbitMQ   RabbitMQConfig `mapstructure:"rabbitmq"`
	BufferSize int            `mapstructure:"buffer_size"`
}

type NATSConfig struct {
	URL     string `mapstructure:"url"`
	Subject string `mapstructure:"subject"`
}

type RabbitMQConfig struct {
	URL      string `mapstructure:"url"`
	Queue    string `mapstructure:"queue"`
	Exchange string `mapstructure:"exchange"`
}

type WorkerConfig struct {
	Concurrency int           `mapstructure:"concurrency"`
	MaxRetry    int           `mapstructure:"max_retry"`
	RetryDelay  time.Duration `mapstructure:"retry_delay"`
}

type ProviderConfig struct {
	Email    EmailProviderConfig    `mapstructure:"email"`
	SMS      SMSProviderConfig      `mapstructure:"sms"`
	Push     PushProviderConfig     `mapstructure:"push"`
	Slack    SlackProviderConfig    `mapstructure:"slack"`
	Telegram TelegramProviderConfig `mapstructure:"telegram"`
	Line     LineProviderConfig     `mapstructure:"line"`
}

type EmailProviderConfig struct {
	Enabled  bool   `mapstructure:"enabled"`
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	From     string `mapstructure:"from"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
}

type SMSProviderConfig struct {
	Enabled    bool   `mapstructure:"enabled"`
	AccountSID string `mapstructure:"account_sid"`
	AuthToken  string `mapstructure:"auth_token"`
	FromNumber string `mapstructure:"from_number"`
}

type PushProviderConfig struct {
	Enabled bool           `mapstructure:"enabled"`
	FCM     FCMConfig      `mapstructure:"fcm"`
	APNs    APNsConfig     `mapstructure:"apns"`
}

type FCMConfig struct {
	CredentialsFile string `mapstructure:"credentials_file"`
}

type APNsConfig struct {
	CertFile   string `mapstructure:"cert_file"`
	KeyFile    string `mapstructure:"key_file"`
	Production bool   `mapstructure:"production"`
}

type SlackProviderConfig struct {
	Enabled bool   `mapstructure:"enabled"`
	Token   string `mapstructure:"token"`
}

type TelegramProviderConfig struct {
	Enabled  bool   `mapstructure:"enabled"`
	BotToken string `mapstructure:"bot_token"`
}

type LineProviderConfig struct {
	Enabled       bool   `mapstructure:"enabled"`
	ChannelSecret string `mapstructure:"channel_secret"`
	ChannelToken  string `mapstructure:"channel_token"`
}

type LoggingConfig struct {
	Level  string `mapstructure:"level"`
	Format string `mapstructure:"format"`
}

// Load reads config from file and environment variables.
func Load(path string) (*Config, error) {
	viper.SetConfigFile(path)
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
