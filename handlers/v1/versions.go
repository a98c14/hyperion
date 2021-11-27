package v1

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/a98c14/hyperion/db"
	"github.com/a98c14/hyperion/utils"
	"github.com/jackc/pgx/v4"
)

type VersionsHandler struct{}

type VersionResponse struct {
	Id      int
	Name    string
	Content string
}

type VersionsRequest struct {
	Id int
}

func (h *VersionsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var head string
	head, r.URL.Path = utils.ShiftPath(r.URL.Path)
	// Get All
	if head == "" {
		switch r.Method {
		case utils.GET:
			handleGetAll(w)
		case utils.POST:
			handlePost(w, r)
		default:
			http.Error(w, "Only GET and POST are allowed", http.StatusMethodNotAllowed)
		}
		return
	}

	// Get by Id
	id, err := strconv.Atoi(head)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid user id %q", head), http.StatusBadRequest)
		return
	}

	switch r.Method {
	case utils.GET:
		handleGet(id, w)
	case utils.POST:
		handlePost(w, r)
	default:
		http.Error(w, "Only GET and POST are allowed", http.StatusMethodNotAllowed)
		return
	}
}

func handleGetAll(w http.ResponseWriter) {
	conn, err := pgx.Connect(context.Background(), db.ConnectionString)
	if err != nil {
		http.Error(w, "Could not connect to database", http.StatusInternalServerError)
		return
	}
	defer conn.Close(context.Background())
	
	var v []VersionResponse
	err = conn.QueryRow(context.Background(), "select (id, name, content) from versions").Scan(&v)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(v)
}

func handleGet(id int, w http.ResponseWriter) {
}

func handlePost(w http.ResponseWriter, r *http.Request) {

}
