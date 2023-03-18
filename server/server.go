package server

import (
	"DSS-uploader/config"
	"context"
	kafka "github.com/segmentio/kafka-go"
	log "github.com/sirupsen/logrus"
)

type Server struct {
	shouldConsume bool
	consumer      *kafka.Reader
}

func NewServer(conf *config.Config) *Server {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{"kafka:9092"},
		GroupID:  "mygroup",
		Topic:    "mytopic",
		MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB
	})

	return &Server{
		shouldConsume: true,
		consumer:      reader,
	}
}

func (server *Server) Run() {
	defer server.consumer.Close()

	for server.shouldConsume {
		msg, err := server.consumer.ReadMessage(context.Background())
		if err != nil {
			log.Error(err)
		} else {
			log.Infof("message at topic:%v partition:%v offset:%v	%s = %s\n", msg.Topic, msg.Partition, msg.Offset, string(msg.Key), string(msg.Value))
		}
	}
}
