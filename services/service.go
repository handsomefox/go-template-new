package services

import "github.com/go-chi/chi/v5"

type Interface interface {
	Bind() func(r *chi.Mux)
}
