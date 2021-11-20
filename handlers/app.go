package handlers

import (
	"net/http"

	v1 "github.com/a98c14/hyperion/handlers/v1"
)

type App struct {
	UserHandler *v1.UserHandler
}

func (h *App) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var head string
	head, req.URL.Path = ShiftPath(req.URL.Path)
	switch head {
	case "user":
		h.UserHandler.ServeHTTP(w, req)
	case "ws":
		// conn, err := upgrader.Upgrade(w, req, nil)
	}

	http.Error(w, "Not Found", http.StatusNotFound)
}

func New() *App {
	app := &App{
		UserHandler: &v1.UserHandler{},
	}

	return app
}
