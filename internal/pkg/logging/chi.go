package logging

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
)

type ZerologAdapter struct {
	logger Logger
}

func (z *ZerologAdapter) Print(v ...interface{}) {
	z.logger.Info().Msg(fmt.Sprint(v...))
}

func GetChiMiddleware(logger Logger) func(http.Handler) http.Handler {
	return middleware.RequestLogger(&middleware.DefaultLogFormatter{
		Logger:  &ZerologAdapter{logger: logger},
		NoColor: true,
	})
}
