package router

import (
	"fmt"
	"net/http"
	"time"

	"github.com/a98c14/hyperion/api/asset"
	prefab "github.com/a98c14/hyperion/api/prefab-editor/handler"
	render "github.com/a98c14/hyperion/api/render/handler"
	"github.com/a98c14/hyperion/api/versioning/handler"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
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
	// Setup middlewares
	r.Use(HandleCors)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	// r.Use(LogRequest)

	// Versions
	r.Route("/versions", func(r chi.Router) {
		r.Get("/", handler.ListVersions)
		r.Get("/{versionId}", handler.GetVersionById)
	})

	// Components
	r.Route("/modules", func(r chi.Router) {
		r.Method("POST", "/", Handler(prefab.SyncModule))
		r.Delete("/", prefab.DeleteModule)
		r.Method("GET", "/{moduleId}", Handler(prefab.GetModuleById))
		r.Method("GET", "/", Handler(prefab.GetRootModules))
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
		r.Method("POST", "/", Handler(render.CreateTexture))
		r.Method("GET", "/", Handler(render.GetTextures))
		r.Method("GET", "/{textureId}", Handler(render.GetTextureFile))
	})

	// Sprites
	r.Route("/sprites", func(r chi.Router) {
		r.Method("POST", "/", Handler(render.CreateSprites))
		r.Method("GET", "/", Handler(render.GetSprites))
	})

	// Animations
	r.Route("/animations", func(r chi.Router) {
		r.Get("/", render.GetAnimations)
		r.Method("POST", "/generate", Handler(render.GenerateAnimationsFromSprites))
	})

	r.Route("/assets", func(r chi.Router) {
		r.Method("GET", "/", Handler(asset.GetAssets))
		r.Method("POST", "/", Handler(asset.SyncAssets))
	})

	r.Route("/health", func(r chi.Router) {
		r.Get("/status", func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "It's ok 😊")
			w.WriteHeader(http.StatusOK)
		})
	})

	return r
}
