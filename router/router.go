package router

import (
	"github.com/a98c14/hyperion/handler"
	"github.com/go-chi/chi/v5"
)

func New() *chi.Mux {
	r := chi.NewRouter()

	// Versions
	r.Route("/versions", func(r chi.Router) {
		r.Get("/", handler.ListVersions)
		r.Get("/{versionId}", handler.GetVersionById)
	})

	// Components
	r.Route("/components", func(r chi.Router) {
		r.Get("/", handler.ListComponents)
		r.Put("/", handler.UpdateComponent)
		r.Post("/", handler.CreateComponent)
		r.Delete("/", handler.DeleteComponent)

		r.Get("/{componentId}", handler.GetComponentById)
		r.Get("/{componentName}", handler.GetComponentByName)

		r.Get("/roots", handler.GetRootComponents)
	})

	// Prefabs

	return r
}
