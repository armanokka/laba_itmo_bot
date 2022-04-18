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
	"strings"
	"syscall"
)

func extract(key, html string) string {
	q := key + ": '"
	i1 := strings.Index(html, q) + len(q)
	i2 := strings.Index(html[i1:], "'")
	return html[i1 : i1+i2]
}

func reverseByDots(s string) string {
	out := ""
	for i, part := range strings.Split(s, ".") {
		if i > 0 {
			out += "."
		}
		for i := len(part) - 1; i > -1; i-- {
			out += string(part[i])
		}
	}
	return out
}

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
		signal := <-signalChanel
		//switch t := signal.(type) {
		//case syscall.SIGHUP:
		//	//reread config...
		//}
		pp.Println("Received", signal.String())
		cancel()
	}()

	db := config.DB()
	botAPI := config.BotAPI()
	botAPI.SetAPIEndpoint("https://api.telegram.org/bot%s/%s")
	analytics := config.Analytics()
	bc, err := bitcask.Open("bitcask_db")
	if err != nil {
		panic(err)
	}
	//d, err := translate.NewDeepl()
	//if err != nil {
	//	panic(err)
	//}

	g, ctx := errgroup.WithContext(ctx)
	g.Go(func() error {
		return boto.New(botAPI, botdb.New(db), analytics, log, bc /*, d*/).Run(ctx)
	})

	g.Go(func() error {
		return server.Run(ctx)
	})

	if err := g.Wait(); err != nil && err != context.Canceled {
		panic(err)
	}
	pp.Println("Program stopped gracefully.")
}
