package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/a98c14/hyperion/db"
	"github.com/a98c14/hyperion/db/query"
	"github.com/a98c14/hyperion/model/render"
	"github.com/jackc/pgx/v4"
)

func CreateSprites(w http.ResponseWriter, r *http.Request) {
	type createSpritesRequest struct {
		TextureId int
		Sprites   []render.Sprite
	}
	// Parse request body
	var req createSpritesRequest
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

	var exists bool
	err = conn.QueryRow(ctx, "select exists(select 1 from texture where id=$1)", req.TextureId).Scan(&exists)
	if err != nil || !exists {
		http.Error(w, "Error:"+err.Error(), http.StatusBadRequest)
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
	batch := &pgx.Batch{}
	for _, sprite := range req.Sprites {
		query.InsertSpriteIfNotExists(ctx, batch, req.TextureId, &sprite)
	}
	br := conn.SendBatch(ctx, batch)
	ct, err := br.Exec()
	if err != nil {
		http.Error(w, "Error during batch insert!", http.StatusInternalServerError)
		return
	}
	fmt.Printf("Inserted rows: %d", ct.RowsAffected())

	err = tx.Commit(ctx)
	if err != nil {
		http.Error(w, "Error during commit, Message:"+err.Error(), http.StatusInternalServerError)
		return
	}
	// for _, value := range req.Modules {
	// 	query.InsertPrefab()
	// }

}
