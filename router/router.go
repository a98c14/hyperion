package router

import (
	"fmt"
	"net/http"
	"time"

	prefab "github.com/a98c14/hyperion/api/prefab-editor/handler"
	render "github.com/a98c14/hyperion/api/render/handler"
	"github.com/a98c14/hyperion/api/versioning/handler"

	"github.com/go-chi/chi/v5"
)

func HandleCors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		next.ServeHTTP(w, r)
	})
}

func LogRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(r.Method, r.URL, r.RemoteAddr, time.Now().Format("01-02-2006 15:04:05"))
		next.ServeHTTP(w, r)
	})
}

func New() *chi.Mux {
	r := chi.NewRouter()
	r.Use(HandleCors)
	r.Use(LogRequest)
	// Versions
	r.Route("/versions", func(r chi.Router) {
		r.Get("/", handler.ListVersions)
		r.Get("/{versionId}", handler.GetVersionById)
	})

	// Components
	r.Route("/modules", func(r chi.Router) {
		r.Put("/", prefab.UpdateModule)
		r.Post("/", prefab.CreateModule)
		r.Delete("/", prefab.DeleteModule)
		r.Get("/{moduleId}", prefab.GetModuleById)
		r.Get("/", prefab.GetRootModules)
	})

	// Prefabs
	r.Route("/prefabs", func(r chi.Router) {
		r.Post("/", prefab.CreatePrefab)
		r.Get("/", prefab.ListPrefabs)
		r.Get("/{prefabId}", prefab.GetPrefabById)
		r.Get("/{prefabId}/versions/{versionId}", prefab.GetPrefabById)
	})

	// Textures
	r.Route("/textures", func(r chi.Router) {
		r.Post("/", render.CreateTexture)
		r.Get("/", render.GetTextures)
		r.Get("/{textureId}", render.GetTextureFile)
	})

	// Sprites
	r.Route("/sprites", func(r chi.Router) {
		r.Post("/", render.CreateSprites)
		r.Get("/", render.GetSprites)
	})

	// Animations
	r.Route("/animations", func(r chi.Router) {
		r.Get("/", render.GetAnimations)
		r.Post("/generate", render.GenerateAnimationsFromSprites)
	})

	return r
}
