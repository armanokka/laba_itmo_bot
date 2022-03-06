// t.me/translobot source code
package main

import (
	"context"
	"encoding/json"
	"git.mills.io/prologic/bitcask"
	boto "github.com/armanokka/translobot/cmd/bot"
	"github.com/armanokka/translobot/cmd/botdb"
	"github.com/armanokka/translobot/cmd/server"
	"github.com/armanokka/translobot/internal/config"
	"github.com/k0kubun/pp"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"golang.org/x/text/language"
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

	ctx, cancel := context.WithCancel(context.Background())
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, syscall.SIGQUIT, syscall.SIGHUP)
	defer signal.Stop(stop)
	go func() {
		<-stop
		cancel()
	}()

	db := botdb.New(config.DB())
	botAPI := config.BotAPI()
	analytics := config.Analytics()
	bc, err := bitcask.Open("bitcask_db")
	if err != nil {
		panic(err)
	}

	bundle := i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)
	bundle.MustLoadMessageFile("resources/ru.json")
	bundle.MustLoadMessageFile("resources/en.json")
	bundle.MustLoadMessageFile("resources/de.json")
	bundle.MustLoadMessageFile("resources/es.json")
	bundle.MustLoadMessageFile("resources/uk.json")
	bundle.MustLoadMessageFile("resources/uz.json")
	bundle.MustLoadMessageFile("resources/id.json")
	bundle.MustLoadMessageFile("resources/it.json")
	bundle.MustLoadMessageFile("resources/pt.json")
	//bundle.MustLoadMessageFile("resources/ar.json")

	g, ctx := errgroup.WithContext(ctx)
	g.Go(func() error {
		return boto.New(botAPI, db, analytics, log, bundle, bc).Run(ctx)
	})

	g.Go(func() error {
		return server.Run()
	})

	if err := g.Wait(); err != nil {
		panic(err)
	}
}


