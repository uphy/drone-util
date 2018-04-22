package main

import (
	"fmt"
	"os"

	"github.com/uphy/drone-util/manager"
	"github.com/uphy/drone-util/model"
	"github.com/urfave/cli"
)

var Version = "0.0.1"

func main() {
	app := cli.NewApp()
	app.Name = "drone-util"
	app.Version = Version
	app.Commands = []cli.Command{
		{
			Name:  "export",
			Usage: "Export repository settings",
			Action: func(c *cli.Context) error {
				m, err := manager.NewFromEnv()
				if err != nil {
					return err
				}
				repos, err := m.Export()
				if err != nil {
					return err
				}
				conf := &model.Config{
					Repos: repos,
				}
				conf.Write(os.Stdout)
				return nil
			},
		},
		{
			Name:      "import",
			Usage:     "Import repository settings",
			ArgsUsage: "<file>",
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "dry-run,d",
					Usage: "if true, command will not update repository",
				},
			},
			Action: func(c *cli.Context) error {
				dryRun := c.Bool("dry-run")
				conf, err := model.ParseFile(c.Args().First())
				if err != nil {
					return err
				}
				if dryRun {
					// dump
					conf.Resolve().Write(os.Stdout)
					return nil
				}
				// update repository
				m, err := manager.NewFromEnv()
				if err != nil {
					return err
				}
				return m.Import(conf.Resolve())
			},
		},
	}
	if err := app.Run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to execute command: %v\n", err)
		os.Exit(1)
	}
}
