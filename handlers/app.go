package handlers

import (
	"net/http"

	v1 "github.com/a98c14/hyperion/handlers/v1"
	"github.com/a98c14/hyperion/utils"
)

type App struct {
	UserHandler    *v1.UserHandler
	VersionHandler *v1.VersionsHandler
}

func (h *App) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var head string
	head, req.URL.Path = utils.ShiftPath(req.URL.Path)
	switch head {
	case "user":
		h.UserHandler.ServeHTTP(w, req)
		return
	case "versions":
		h.VersionHandler.ServeHTTP(w, req)
		return
	case "ws":
		// conn, err := upgrader.Upgrade(w, req, nil)
	}

	http.Error(w, "Not Found", http.StatusNotFound)
}

func New() *App {
	app := &App{
		UserHandler:    &v1.UserHandler{},
		VersionHandler: &v1.VersionsHandler{},
	}

	return app
}
