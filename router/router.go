package router

import (
	"net/http"

	"github.com/a98c14/hyperion/handler"
	"github.com/go-chi/chi/v5"
)

func HandleCors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		next.ServeHTTP(w, r)
	})
}

func New() *chi.Mux {
	r := chi.NewRouter()
	r.Use(HandleCors)
	// Versions
	r.Route("/versions", func(r chi.Router) {
		r.Get("/", handler.ListVersions)
		r.Get("/{versionId}", handler.GetVersionById)
	})

	// Components
	r.Route("/modules", func(r chi.Router) {
		r.Put("/", handler.UpdateComponent)
		r.Post("/", handler.CreateComponent)
		r.Delete("/", handler.DeleteComponent)
		r.Get("/{componentId}", handler.GetModuleById)
		r.Get("/", handler.GetRootComponents)
	})

	// Prefabs

	// Textures
	r.Route("/textures", func(r chi.Router) {
		r.Post("/", handler.CreateTexture)
		r.Get("/", handler.GetTextures)
		r.Get("/{textureId}", handler.GetTextureFile)
	})

	// Sprites
	r.Route("/sprites", func(r chi.Router) {
		r.Post("/", handler.CreateSprites)
	})

	// Animations
	r.Route("/animations", func(r chi.Router) {
		r.Post("/generate", handler.GenerateAnimationsFromSprites)
	})

	return r
}
