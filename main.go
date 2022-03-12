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
	"golang.org/x/sync/errgroup"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

func extract(key, html string) string {
	q := key + ": '"
	i1 := strings.Index(html, q) +  len(q)
	i2 := strings.Index(html[i1:], "'")
	return html[i1:i1+i2]
}

func reverseByDots(s string) string {
	out := ""
	for i, part := range strings.Split(s, ".") {
		if i > 0 {
			out += "."
		}
		for i := len(part)-1; i > -1; i-- {
			out += string(part[i])
		}
	}
	return out
}

func main() {
	pp.WithLineInfo = true
	if err := config.Load(); err != nil {
		panic(err)
	}

	log, _ := zap.NewProduction()
	defer log.Sync()

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, os.Interrupt, syscall.SIGQUIT, syscall.SIGHUP)
	defer cancel()

	db := config.DB()
	botAPI := config.BotAPI()
	analytics := config.Analytics()
	bc, err := bitcask.Open("bitcask_db")
	if err != nil {
		panic(err)
	}

	g, ctx := errgroup.WithContext(ctx)
	g.Go(func() error {
		return boto.New(botAPI, botdb.New(db), analytics, log, bc).Run(ctx)
	})

	g.Go(func() error {
		return server.Run()
	})

	if err := g.Wait(); err != nil {
		panic(err)
	}
}


