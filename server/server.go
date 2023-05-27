package server

import (
	"DSS-uploader/config"
	"DSS-uploader/upload"
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/yakiroren/dss-common/db"

	log "github.com/sirupsen/logrus"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Server struct {
	conn      *amqp.Connection
	channel   *amqp.Channel
	queue     amqp.Queue
	upload    upload.Client
	dataStore db.DataStore
}

func NewServer(conf config.RabbitConfig, client upload.Client, dataStore db.DataStore) (*Server, error) {
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
		conn:      conn,
		channel:   channel,
		queue:     queue,
		upload:    client,
		dataStore: dataStore,
	}, nil
}

func (server *Server) Close() {
	server.channel.Close()
	server.conn.Close()
}

func (server *Server) Consume() {
	msgs, err := server.channel.Consume(
		server.queue.Name, // queue
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
		go func(msg amqp.Delivery) {
			fragmentNumber := msg.Headers["fragment_number"].(string)
			id := msg.Headers["id"].(string)

			log.Infof("Got fragment number %s for object %s", fragmentNumber, id)

			err := server.upload.Upload(context.Background(), id, msg.Body, fragmentNumber)
			if err != nil {
				log.Error(err)
				msg.Nack(false, true)
				return
			}

			log.Infof("fragment %s object %s uploaded successfully", fragmentNumber, id)
			msg.Ack(false)

			hex, err := primitive.ObjectIDFromHex(id)
			if err != nil {
				log.Error(err)
				return
			}

			metadata, found := server.dataStore.GetMetadataByID(context.Background(), hex)
			if !found {
				log.Error("Couldn't find object with id ", id)
				return
			}

			if metadata.TotalFragments == len(metadata.Fragments) {
				log.Info("setting ishidden to false")
				err := server.dataStore.UpdateField(context.Background(), id, "ishidden", false)
				if err != nil {
					log.Error(err)
				}
			}
		}(msg)
	}

	<-forever
}
