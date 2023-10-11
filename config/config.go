package config

import (
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Environment string `yaml:"BOT_TOKEN" env:"BOT_TOKEN" env-default:"debug"`

	BotToken string `yaml:"BOT_TOKEN" env:"BOT_TOKEN" env-required:"true"`
	AdminID  int64  `yaml:"ADMIN_ID" env:"ADMIN_ID" env-required:"true"`

	PostgresHost     string `yaml:"POSTGRES_HOST" env:"POSTGRES_HOST" env-required:"true"`
	PostgresPort     int    `yaml:"POSTGRES_PORT" env:"POSTGRES_PORT" env-required:"true"`
	PostgresUser     string `yaml:"POSTGRES_USER" env:"POSTGRES_USER" env-required:"true"`
	PostgresDB       string `yaml:"POSTGRES_DB" env:"POSTGRES_DB" env-required:"true"`
	PostgresPassword string `yaml:"POSTGRES_PASSWORD" env:"POSTGRES_PASSWORD" env-required:"true"`
}

func NewConfig() (*Config, error) {
	cfg := &Config{}

	if err := cleanenv.ReadEnv(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
