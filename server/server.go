package server

import (
	"DSS-uploader/config"
	"fmt"
	log "github.com/sirupsen/logrus"

	"github.com/dustin/go-humanize"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Server struct {
	Conn    *amqp.Connection
	Channel *amqp.Channel
	Queue   amqp.Queue
}

func NewServer(conf *config.Config) (*Server, error) {
	conn, err := amqp.Dial(conf.RabbitUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	channel, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to open a channel: %w", err)
	}

	queue, err := channel.QueueDeclare(
		conf.QueueName, // name
		true,           // durable
		false,          // delete when unused
		false,          // exclusive
		false,          // no-wait
		nil,            // arguments
	)
	if err != nil {
		return nil, fmt.Errorf("Failed to declare a queue: %w", err)
	}

	return &Server{
		Conn:    conn,
		Channel: channel,
		Queue:   queue,
	}, nil
}

func (server *Server) Close() {
	server.Channel.Close()
	server.Conn.Close()
}

func (server *Server) Consume() {
	msgs, err := server.Channel.Consume(
		server.Queue.Name, // queue
		"",                // consumer
		false,             // auto-ack
		false,             // exclusive
		false,             // no-local
		false,             // no-wait
		nil,               // args
	)

	if err != nil {
		log.Fatalf("could not start consumeing: %w", err)
	}

	var forever chan struct{}

	for msg := range msgs {
		log.Info(msg.Headers)
		log.Info(humanize.IBytes(uint64(len(msg.Body))))

		msg.Ack(false)
	}

	<-forever
}
