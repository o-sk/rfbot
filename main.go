package main

import (
	"log"
	"os"

	"github.com/o-sk/slackbot-sample/config"
	"github.com/urfave/cli"
)

func main() {
	app := &cli.App{
		Name: "slackbot-sample",
		Action: func(c *cli.Context) error {
			_ = config.Load("config.yml")

			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
