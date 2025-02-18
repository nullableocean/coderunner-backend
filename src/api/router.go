package api

import (
	"github.com/go-chi/chi/v5"
)

func NewChiRouter(h *CompileHandler) *chi.Mux {
	r := chi.NewRouter()

	registerCompilerRoutes(r, h)

	return r
}

func registerCompilerRoutes(router chi.Router, h *CompileHandler) {
	router.Post("/task", h.Create)
	router.Get("/status/{task_id}", h.Status)
	router.Get("/result/{task_id}", h.Result)
}
