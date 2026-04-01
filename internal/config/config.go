package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"
)

// Config holds all application configuration.
type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	RabbitMQ RabbitMQConfig `mapstructure:"rabbitmq"`
	Keycloak KeycloakConfig `mapstructure:"keycloak"`
	Worker   WorkerConfig   `mapstructure:"worker"`
	Logging  LoggingConfig  `mapstructure:"logging"`
}

// ServerConfig holds HTTP server settings.
type ServerConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
	Mode string `mapstructure:"mode"`
}

// DatabaseConfig holds PostgreSQL settings.
type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Name     string `mapstructure:"name"`
	SSLMode  string `mapstructure:"ssl_mode"`
	MaxConns int32  `mapstructure:"max_conns"`
	MinConns int32  `mapstructure:"min_conns"`
}

// RabbitMQConfig holds RabbitMQ settings.
type RabbitMQConfig struct {
	URL           string `mapstructure:"url"`
	PrefetchCount int    `mapstructure:"prefetch_count"`
}

// KeycloakConfig holds Keycloak authentication settings.
type KeycloakConfig struct {
	BaseURL      string `mapstructure:"base_url"`
	Realm        string `mapstructure:"realm"`
	ClientID     string `mapstructure:"client_id"`
	ClientSecret string `mapstructure:"client_secret"`
}

// WorkerConfig holds background worker settings.
type WorkerConfig struct {
	Concurrency int           `mapstructure:"concurrency"`
	MaxRetry    int           `mapstructure:"max_retry"`
	RetryDelay  time.Duration `mapstructure:"retry_delay"`
}

// LoggingConfig holds logging settings.
type LoggingConfig struct {
	Level  string `mapstructure:"level"`
	Format string `mapstructure:"format"`
}

// Load reads configuration from file and environment.
func Load(configPath string) (*Config, error) {
	v := viper.New()

	// Set defaults
	v.SetDefault("server.host", "0.0.0.0")
	v.SetDefault("server.port", 8080)
	v.SetDefault("server.mode", "debug")

	v.SetDefault("database.host", "localhost")
	v.SetDefault("database.port", 5432)
	v.SetDefault("database.user", "postgres")
	v.SetDefault("database.password", "postgres")
	v.SetDefault("database.name", "notification_center")
	v.SetDefault("database.ssl_mode", "disable")
	v.SetDefault("database.max_conns", 25)
	v.SetDefault("database.min_conns", 5)

	v.SetDefault("rabbitmq.url", "amqp://guest:guest@localhost:5672/")
	v.SetDefault("rabbitmq.prefetch_count", 10)

	v.SetDefault("keycloak.base_url", "http://localhost:8180")
	v.SetDefault("keycloak.realm", "notification-center")
	v.SetDefault("keycloak.client_id", "notification-api")

	v.SetDefault("worker.concurrency", 10)
	v.SetDefault("worker.max_retry", 3)
	v.SetDefault("worker.retry_delay", "5s")

	v.SetDefault("logging.level", "info")
	v.SetDefault("logging.format", "json")

	// Read config file
	if configPath != "" {
		v.SetConfigFile(configPath)
	} else {
		v.SetConfigName("config")
		v.SetConfigType("yaml")
		v.AddConfigPath(".")
		v.AddConfigPath("./config")
	}

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
	}

	// Environment variables override
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	// Explicit env bindings
	envBindings := map[string]string{
		"SERVER_HOST":        "server.host",
		"SERVER_PORT":        "server.port",
		"DB_HOST":            "database.host",
		"DB_PORT":            "database.port",
		"DB_USER":            "database.user",
		"DB_PASSWORD":        "database.password",
		"DB_NAME":            "database.name",
		"DB_SSL_MODE":        "database.ssl_mode",
		"RABBITMQ_URL":       "rabbitmq.url",
		"KEYCLOAK_BASE_URL":  "keycloak.base_url",
		"KEYCLOAK_REALM":     "keycloak.realm",
		"KEYCLOAK_CLIENT_ID": "keycloak.client_id",
		"KEYCLOAK_SECRET":    "keycloak.client_secret",
	}

	for env, key := range envBindings {
		if err := v.BindEnv(key, env); err != nil {
			return nil, fmt.Errorf("failed to bind env %s: %w", env, err)
		}
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &cfg, nil
}

// Address returns the server address.
func (c *ServerConfig) Address() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

// DSN returns the PostgreSQL connection string.
func (c *DatabaseConfig) DSN() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		c.User, c.Password, c.Host, c.Port, c.Name, c.SSLMode,
	)
}

// JWKSURL returns the Keycloak JWKS endpoint.
func (c *KeycloakConfig) JWKSURL() string {
	return fmt.Sprintf("%s/realms/%s/protocol/openid-connect/certs", c.BaseURL, c.Realm)
}

// TokenURL returns the Keycloak token endpoint.
func (c *KeycloakConfig) TokenURL() string {
	return fmt.Sprintf("%s/realms/%s/protocol/openid-connect/token", c.BaseURL, c.Realm)
}

// UserInfoURL returns the Keycloak userinfo endpoint.
func (c *KeycloakConfig) UserInfoURL() string {
	return fmt.Sprintf("%s/realms/%s/protocol/openid-connect/userinfo", c.BaseURL, c.Realm)
}
