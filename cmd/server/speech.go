package server

import (
	"bytes"
	"github.com/armanokka/translobot/pkg/translate"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func speechHandler(log *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		log := log.With()
		defer func() {
			if err := recover(); err != nil {
				log.Error("panic", zap.Any("error", err))
			}
		}()
		warn := func(code int, err error) {
			c.AbortWithStatusJSON(code, err)
			log.Error("", zap.Error(err))
		}
		lang := c.Query("lang")
		if lang == "" {
			warn(400, ErrEmpty("lang"))
			return
		}
		text := c.Query("text")
		if text == "" {
			warn(400, ErrEmpty("text"))
			return
		}
		data, err := translate.TTS(lang, text)
		if err != nil {
			warn(500, ErrInternal)
			log.Error("", zap.Error(err))
			return
		}

		c.DataFromReader(200, int64(len(data)), "audio/mp3", bytes.NewBuffer(data), nil)
	}
}
