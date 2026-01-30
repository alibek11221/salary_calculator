package logging

import (
	"os"

	"github.com/rs/zerolog"
)

type Logger interface {
	Info() *zerolog.Event
	Debug() *zerolog.Event
	Warn() *zerolog.Event
	Error() *zerolog.Event
	Fatal() *zerolog.Event
	With() zerolog.Context
}

type logger struct {
	zl zerolog.Logger
}

func New(isProduction bool) Logger {
	var zl zerolog.Logger
	if isProduction {
		zl = zerolog.New(os.Stdout).With().Timestamp().Logger()
	} else {
		zl = zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: "15:04:05"}).
			With().
			Timestamp().
			Logger()
	}
	return &logger{zl: zl}
}

func (l *logger) Info() *zerolog.Event {
	return l.zl.Info()
}

func (l *logger) Debug() *zerolog.Event {
	return l.zl.Debug()
}

func (l *logger) Warn() *zerolog.Event {
	return l.zl.Warn()
}

func (l *logger) Error() *zerolog.Event {
	return l.zl.Error()
}

func (l *logger) Fatal() *zerolog.Event {
	return l.zl.Fatal()
}

func (l *logger) With() zerolog.Context {
	return l.zl.With()
}
