package main

import (
	"log"
	"os"

	"github.com/takehaya/PPAP_Protocol/internal"
	"github.com/takehaya/PPAP_Protocol/pkg/version"
)

func main() {
	app := internal.NewApp(version.Version)
	if err := app.Run(os.Args); err != nil {
		log.Fatalf("%+v", err)
	}
}

