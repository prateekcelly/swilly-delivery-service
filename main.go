package main

import (
	"log"
	"os"
	"swilly-delivery-service/internal/app/server"
	"swilly-delivery-service/internal/app/worker"

	"github.com/urfave/cli/v2"
)

// main godoc
//
//	@title API Documentation for swilly-delivery-service
//	@version 1.0.0
//	@description Responsible for invoking message delivery for users
//	@contact.name Prateek Celly
//	@contact.email prateekcelly@gmail.com
//	@BasePath /
//	@query.collection.format multi
func main() {
	app := cli.NewApp()
	app.Name = "swilly-delivery-service"
	app.Version = "1.0.0"
	app.Commands = []*cli.Command{
		{
			Name:  "worker",
			Usage: "Worker mode for Swilly Delivery Service",
			Action: func(context *cli.Context) error {
				log.Println("Starting worker mode")
				return worker.StartWorker(context.Context)
			},
		},
		{
			Name:        "server",
			Description: "Server mode for Swilly Delivery Service",
			Action: func(context *cli.Context) error {
				log.Println("Starting server mode")
				server.StartServer()
				return nil
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
