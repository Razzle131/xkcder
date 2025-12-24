package config

import (
	"log"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	LogLevel    string `yaml:"log_level" env:"LOG_LEVEL" env-default:"DEBUG"`
	ApiAddress  string `yaml:"api_address" env:"API_ADDRESS" env-default:"api:80"`
	TokenCookie string `yaml:"token_cookie" env:"TOKEN_COOKIE" env-default:"XKCDtoken"`
}

func MustLoad(configPath string) Config {
	var cfg Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("cannot read config %q: %s", configPath, err)
	}
	return cfg
}
