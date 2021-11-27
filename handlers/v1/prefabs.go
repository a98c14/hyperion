package v1

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/a98c14/hyperion/utils"
)

type PrefabsHandler struct{}

func (h *PrefabsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var head string
	head, r.URL.Path = utils.ShiftPath(r.URL.Path)

	if head == "" {
		// If no id is specified
		switch r.Method {
		case utils.GET:
			handleGetPrefabs(w)
		case utils.POST:
			handleCreatePrefab(w, r)
		}
	} else {
		// If there is an id available
		id, err := strconv.Atoi(head)
		if err != nil {
			http.Error(w, fmt.Sprintf("Invalid prefab id %q", head), http.StatusBadRequest)
			return
		}

		switch r.Method {
		case utils.GET:
			handleGetPrefab(id, w)
		case utils.PUT:
			handleUpdatePrefab(id, w, r)
		}
	}
}

func handleGetPrefab(id int, w http.ResponseWriter) {

}

func handleGetPrefabs(w http.ResponseWriter) {

}

func handleCreatePrefab(w http.ResponseWriter, r *http.Request) {

}

func handleUpdatePrefab(id int, w http.ResponseWriter, r *http.Request) {

}
