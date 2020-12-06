package internal

import (
	"fmt"

	"github.com/urfave/cli"
)

func NewApp(version string) *cli.App {
	app := cli.NewApp()
	app.Name = "ppap"
	app.Version = version

	app.Usage = "ppap protocol"

	app.EnableBashCompletion = true
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "mode",
			Value: "client",
			Usage: "mode selected plz input client or server",
		},
		cli.StringFlag{
			Name:  "bindaddr",
			Value: "auto",
			Usage: "bind address(use srcaddr)",
		},
		cli.StringFlag{
			Name:  "gateway",
			Value: "172.27.1.2",
			Usage: "gateway address",
		},
		cli.StringFlag{
			Name:  "have1",
			Value: "Pen",
			Usage: "keys",
		},
		cli.StringFlag{
			Name:  "have2",
			Value: "have",
			Usage: "gateway address",
		},
	}
	app.Action = run
	return app
}

func run(ctx *cli.Context) error {
	gateway := ctx.String("gateway")
	mode := ctx.String("mode")
	bindaddr := ctx.String("bindaddr")

	h1 := ctx.String("have1")
	h2 := ctx.String("have2")

	ctl, err := NewController(gateway,bindaddr, h1, h2)

	if err != nil {
		return err
	}
	switch mode {
	case "client":
		return Client(ctl)
	case "server":
		return Server(ctl)
	default:
		return fmt.Errorf(fmt.Sprintf("%v is inv mode. plz check mode", mode))
	}

	return nil
}
