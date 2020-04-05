package main

import (
	"log"
	"os"

	"github.com/uphy/karabiner-config/app"
)

func main() {
	app := app.New()
	if err := app.Run(os.Args); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
