package main

import (
	"github.com/armanokka/laba_itmo_bot/config"
	"github.com/armanokka/laba_itmo_bot/internal/app"
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
