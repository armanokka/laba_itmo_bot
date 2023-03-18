package main

import (
	"github.com/armanokka/translobot/config"
	"github.com/armanokka/translobot/internal/app"
)

func main() {
	// Loading config
	cfg, err := config.NewConfig()
	if err != nil {
		panic(err)
	}

	err = app.Run(cfg)
	if err != nil {
		panic(err)
	}
}
