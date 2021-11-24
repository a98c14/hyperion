package v1

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/a98c14/hyperion/utils"
)

type ComponentsHandler struct{}

func (h *ComponentsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var head string
	head, r.URL.Path = utils.ShiftPath(r.URL.Path)

	if head == "" {
		switch r.Method {
		case utils.GET:
			handleGetComponents(w)
		case utils.POST:
			handleCreateComponent(w, r)
		}
		return
	}

	// Get by Id
	id, err := strconv.Atoi(head)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid component id %q", head), http.StatusBadRequest)
		return
	}

	switch r.Method {
	case utils.GET:
		handleGetComponent(id, w)
	}

}

func handleGetComponent(id int, w http.ResponseWriter) {

}

func handleGetComponents(w http.ResponseWriter) {

}

func handleCreateComponent(w http.ResponseWriter, r *http.Request) {

}
