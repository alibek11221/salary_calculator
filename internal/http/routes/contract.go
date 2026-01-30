package routes

import (
	"github.com/go-chi/chi/v5"
)

type RouterRegistrar interface {
	Register(router chi.Router)
	Name() string
}
