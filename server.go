package main

import (
	"DSS-uploader/config"
	"DSS-uploader/server"
	discord "DSS-uploader/upload/discord/webhooks"
	"fmt"

	"github.com/yakiroren/dss-common/db"

	"github.com/caarlos0/env/v7"
	"github.com/gofiber/adaptor/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/hellofresh/health-go/v5"
	log "github.com/sirupsen/logrus"
)

func main() {
	conf := &config.Config{}
	opts := env.Options{UseFieldNameByDefault: true}

	if err := env.Parse(conf, opts); err != nil {
		log.Fatal(err)
	}

	log.SetLevel(conf.LogLevel)

	store, err := db.NewMongoDataStore(&conf.Mongo)
	if err != nil {
		log.Fatal(err)
	}

	client, err := discord.New(store, conf.Discord)
	if err != nil {
		log.Fatal(err)
	}

	srv, err := server.NewServer(conf.Rabbit, client, store)
	if err != nil {
		log.Fatal(err)
	}

	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
	})

	app.Use(recover.New())
	app.Use(logger.New())

	go srv.Consume()

	h, _ := health.New(health.WithSystemInfo(), health.WithComponent(health.Component{
		Name:    "Uploader",
		Version: "v1",
	}))

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello from Uploader!")
	})

	app.Get("/status", adaptor.HTTPHandler(h.Handler()))

	if err = app.Listen(fmt.Sprintf(":%s", conf.Port)); err != nil {
		log.Fatal(err)
	}
}
