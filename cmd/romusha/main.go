package main

import (
	"os"

	"github.com/urfave/cli/v2"

	"github.com/parinpan/romusha/internal/app/master"
	"github.com/parinpan/romusha/internal/app/worker"
)

func main() {
	app := &cli.App{
		Name:  "romusha",
		Usage: "MapReduce opinionated implementation by @parinpan",
		Commands: []*cli.Command{
			{
				Name:  "master",
				Usage: "start the master app",
				Action: func(c *cli.Context) error {
					return master.Start(c.Context)
				},
			},
			{
				Name:  "worker",
				Usage: "start the worker app",
				Action: func(c *cli.Context) error {
					return worker.Start(c.Context)
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		panic(err)
	}
}
