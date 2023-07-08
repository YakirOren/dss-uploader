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
		go server.consumeMessage(msg)
	}

	<-forever
}

func (server *Server) consumeMessage(msg amqp.Delivery) {
	fragmentNumber := msg.Headers["fragment_number"].(string)
	id := msg.Headers["id"].(string)

	cl := log.WithFields(log.Fields{
		"fragment_number": fragmentNumber,
		"object_id":       id,
	})

	cl.Infof("Got fragment")
	hex, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		cl.Error("invalid object id ", err)
		DiscardMsg(msg, cl)

		return
	}

	metadata, found := server.dataStore.GetMetadataByID(context.Background(), hex)
	if !found {
		cl.Error("Couldn't find object with id ", id)
		DiscardMsg(msg, cl)
		return
	}

	err = server.upload.Upload(context.Background(), id, msg.Body, fragmentNumber)
	if err != nil {
		cl.Error(err)
		msg.Nack(false, true)
		return
	}

	cl.Info("fragment uploaded successfully")
	msg.Ack(false)

	if metadata.TotalFragments == len(metadata.Fragments) {
		cl.Info("setting ishidden to false")
		err := server.dataStore.UpdateField(context.Background(), id, "ishidden", false)
		if err != nil {
			cl.Error(err)
		}
	}
}

func DiscardMsg(msg amqp.Delivery, logger *log.Entry) {
	msg.Nack(false, false)
	logger.Info("discarding message")
}
