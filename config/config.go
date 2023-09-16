package config

import (
	discord "DSS-uploader/upload/discord/webhooks"

	log "github.com/sirupsen/logrus"
	"github.com/yakiroren/dss-common/db"
)

type RabbitConfig struct {
	RabbitURL string `env:"RABBIT_URL,required,notEmpty"`
	QueueName string `env:"QUEUE_NAME,required,notEmpty"`
}

type Config struct {
	Port     string    `env:"PORT,required,notEmpty"`
	LogLevel log.Level `env:"LOG_LEVEL,required,notEmpty"`
	Rabbit   RabbitConfig
	Mongo    db.MongoConfig
	Discord  discord.DiscordWebhookConfig
}
