package logging

import (
	"fmt"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog/log"
)

type ZerologAdapter struct{}

func (z *ZerologAdapter) Print(v ...interface{}) {
	log.Info().Msg(fmt.Sprint(v...))
}

func init() {
	middleware.DefaultLogger = middleware.RequestLogger(&middleware.DefaultLogFormatter{
		Logger:  &ZerologAdapter{},
		NoColor: true,
	})
}
