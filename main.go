// t.me/translobot source code
package main

import (
	"context"
	boto "github.com/armanokka/translobot/cmd/bot"
	"github.com/armanokka/translobot/cmd/server"
	"github.com/armanokka/translobot/internal/config"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	if err := config.Load(); err != nil {
		panic(err)
	}

	log, _ := zap.NewProduction()
	defer log.Sync()
	sugar := log.Sugar()


	ctx, cancel := context.WithCancel(context.Background())
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, syscall.SIGQUIT, syscall.SIGHUP)
	defer signal.Stop(stop)
	go func() {
		<-stop
		cancel()
	}()

	db := config.DB()
	botAPI := config.BotAPI()
	analytics := config.Analytics()

	g, ctx := errgroup.WithContext(ctx)
	g.Go(func() error {
		return boto.Run(ctx, botAPI, db, analytics, sugar)
	})

	g.Go(func() error {
		return server.Run()
	})

	if err := g.Wait(); err != nil {
		panic(err)
	}
}


