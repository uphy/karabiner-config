package app

import (
	"errors"
	"fmt"
	"path/filepath"

	"github.com/urfave/cli"
)

var (
	version  = "unknown"
	revision = "unknown"
)

type CLI struct {
	app *cli.App
}

func New() *CLI {
	app := cli.NewApp()
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name: "w,watch",
		},
		cli.StringFlag{
			Name:  "p,profile",
			Value: "Default profile",
		},
		cli.BoolFlag{
			Name: "k,karabiner",
		},
	}
	app.Name = "karabiner-config"
	app.Usage = "Karabiner Config File Generator"
	app.Author = "uphy"
	app.Version = fmt.Sprintf("%s (%s)", version, revision)
	app.Action = func(ctx *cli.Context) error {
		watch := ctx.Bool("watch")
		profile := ctx.String("profile")

		if ctx.NArg() != 2 {
			return errors.New("specify input, output as arguments")
		}
		input := ctx.Args().Get(0)
		output := ctx.Args().Get(1)
		_, filename := filepath.Split(output)

		var karabinerMode bool
		if ctx.IsSet("karabiner") {
			karabinerMode = ctx.Bool("karabiner")
		} else {
			karabinerMode = filename == "karabiner.json"
		}

		var writer ConfigWriter
		if karabinerMode {
			writer = &KarabinierJSONWriter{output, profile}
		} else {
			writer = &OverwriteConfigWriter{output}
		}
		runner := &Runner{writer}
		if watch {
			return runner.Watch(input)
		}
		return runner.Run(input)
	}
	return &CLI{app}
}

func (c *CLI) Run(args []string) error {
	return c.app.Run(args)
}
