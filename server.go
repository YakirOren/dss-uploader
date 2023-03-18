package main

import (
	"DSS-uploader/config"
	"DSS-uploader/server"
	"fmt"
	"github.com/caarlos0/env"
	"github.com/gofiber/fiber/v2"
	log "github.com/sirupsen/logrus"
)

func main() {
	conf := new(config.Config)
	if err := env.Parse(conf); err != nil {
		log.Fatal(err)
	}

	log.SetLevel(conf.LogLevel)

	srv := server.NewServer(conf)

	app := fiber.New()

	go srv.Run()

	setupRoute(app, srv)

	if err := app.Listen(fmt.Sprintf(":%s", conf.Port)); err != nil {
		log.Fatal(err)
	}
}

func setupRoute(app *fiber.App, srv *server.Server) {
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})
}
