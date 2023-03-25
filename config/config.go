package config

import log "github.com/sirupsen/logrus"

type Config struct {
	Port                   string    `env:"PORT,required"`
	RabbitUrl              string    `env:"RABBIT_URL,required"`
	QueueName              string    `env:"QUEUE_NAME,required"`
	LogLevel               log.Level `env:"LOG_LEVEL,required"`
	DiscordToken           string    `env:"DISCORD_TOKEN,required"`
	DiscordStorageChannels []string  `env:"STORAGE_CHANNELS,required"`
	DBName                 string    `env:"DB_NAME,required"`
	MongoURL               string    `env:"MONGODB_URL,required"`
	MongoDbUsername        string    `env:"ME_CONFIG_MONGODB_ADMINUSERNAME,required"`
	MongoDbPassword        string    `env:"ME_CONFIG_MONGODB_ADMINPASSWORD,required"`
	FileCollection         string    `env:"FILE_COLLECTION,required"`
}
