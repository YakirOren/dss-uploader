package config

import log "github.com/sirupsen/logrus"

type Config struct {
	Port     string    `env:"PORT"`
	LogLevel log.Level `env:"LOG_LEVEL"`
}
