package handlers

import "net/http"

type App struct {
	UserHandler *UserHandler
}

func (h *App) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var head string
	head, req.URL.Path = shiftPath(req.URL.Path)
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
		UserHandler: &UserHandler{},
	}

	return app
}
