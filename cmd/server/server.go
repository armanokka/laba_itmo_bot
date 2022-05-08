package server

import (
	"context"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/k0kubun/pp"
	"go.uber.org/zap"
	"net/http"
	"os"
)

func Run(ctx context.Context, log *zap.Logger) error {
	r := gin.Default()
	pprof.Register(r)
	r.GET("/", func(c *gin.Context) {
		c.String(200, "ok")
	})
	r.GET("/speech", speechHandler(log))
	// Ports for Heroku
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	server := &http.Server{Addr: ":" + port, Handler: r}
	go func() {
		<-ctx.Done()
		err := server.Shutdown(ctx)
		if err != nil && err != http.ErrServerClosed && err != context.Canceled {
			pp.Println("server.Run: Error:", err)
		}
	}()
	err := server.ListenAndServe()
	if err != http.ErrServerClosed {
		return err
	}
	return nil
}
