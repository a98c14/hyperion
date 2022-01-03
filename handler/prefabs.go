package handler

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/a98c14/hyperion/db"
	"github.com/a98c14/hyperion/db/query"
	"github.com/a98c14/hyperion/model/prefab"
)

func GetPrefabById(id int, w http.ResponseWriter) {

}

func ListPrefabs(w http.ResponseWriter) {

}

func CreatePrefab(w http.ResponseWriter, r *http.Request) {
	type requestModel struct {
		Name     string                   `json:name`
		ParentId sql.NullInt32            `json:parentId`
		Modules  []prefab.ModulePartValue `json:modules`
	}

	// Parse request body
	var req requestModel
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Could not parse request body!", http.StatusBadRequest)
		return
	}
	ctx := r.Context()
	conn, err := db.GetConnectionPool(ctx)
	if err != nil {
		http.Error(w, "Could not connect to database!", http.StatusInternalServerError)
		return
	}

	// Check if prefab with given name already exists
	var exists bool
	err = conn.QueryRow(ctx, "select exists(select 1 from prefab where name=$1)", req.Name).Scan(&exists)
	if exists || err != nil {
		http.Error(w, "Prefab with given name already exists!", http.StatusBadRequest)
		return
	}

	// Start transcation. If all components can not be added successfully, don't
	// insert anything
	tx, err := conn.Begin(ctx)
	if err != nil {
		http.Error(w, "Could not start transaction!", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback(ctx)

	_, err = query.InsertPrefab(ctx, conn, req.Name, req.ParentId)
	if err != nil {
		http.Error(w, "Error during prefab insert, Error: "+err.Error(), http.StatusInternalServerError)
	}

}

func UpdatePrefab(id int, w http.ResponseWriter, r *http.Request) {

}
