// t.me/translobot source code
package main

import (
	"context"
	"git.mills.io/prologic/bitcask"
	boto "github.com/armanokka/translobot/cmd/bot"
	"github.com/armanokka/translobot/cmd/botdb"
	"github.com/armanokka/translobot/cmd/server"
	"github.com/armanokka/translobot/internal/config"
	"github.com/k0kubun/pp"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"golang.org/x/sync/errgroup"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	conf := zap.NewDevelopmentConfig()
	conf.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	conf.DisableStacktrace = true
	log, _ := conf.Build()
	defer log.Sync()

	pp.WithLineInfo = true
	if err := config.Load(); err != nil {
		panic(err)
	}

	signalChanel := make(chan os.Signal, 1)
	signal.Notify(signalChanel,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)
	defer signal.Stop(signalChanel)

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		<-signalChanel
		pp.Println("exiting...")
		cancel()
	}()

	db := config.DB()
	arangodb := config.ArangoDB()
	botAPI := config.BotAPI()

	//botAPI.SetAPIEndpoint("http://127.0.0.1:8081/bot%s/%s")
	analytics := config.Analytics()
	bc, err := bitcask.Open("bitcask_db")
	if err != nil {
		panic(err)
	}

	g, ctx := errgroup.WithContext(ctx)
	g.Go(func() error {
		Bot, err := boto.New(botAPI, botdb.New(db), analytics, log, bc /*, d*/, arangodb)
		if err != nil {
			return err
		}
		return Bot.Run(ctx)
	})

	g.Go(func() error {
		return server.Run(ctx, log)
	})

	if err := g.Wait(); err != nil && err != context.Canceled {
		panic(err)
	}
}
