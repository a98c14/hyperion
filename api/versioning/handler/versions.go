package handler

import (
	"context"
	"net/http"

	"github.com/a98c14/hyperion/common/response"
	"github.com/a98c14/hyperion/db"
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

func ListVersions(w http.ResponseWriter, r *http.Request) {
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

	response.Json(w, v)
}

func GetVersionById(w http.ResponseWriter, r *http.Request) {
}

func CreateVersion(w http.ResponseWriter, r *http.Request) {

}
