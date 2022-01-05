package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/a98c14/hyperion/api/render/data"
	"github.com/a98c14/hyperion/common"
	"github.com/a98c14/hyperion/common/response"
	"github.com/jackc/pgx/v4"
)

func GetSprites(w http.ResponseWriter, r *http.Request) {
	state, err := common.InitState(r)
	if err != nil {
		response.ErrorWhileInitializing(w, err)
		return
	}
	rows, err := state.Conn.Query(state.Context,
		`select id, texture_id, unity_name, unity_pivot, 
				unity_rect, unity_border, unity_alignment
		 from sprite`)
	if err != nil {
		response.InternalError(w, err)
		return
	}
	defer rows.Close()
	type spriteResponse struct {
		Id        int       `json:"id"`
		TextureId int       `json:"textureId"`
		Name      string    `json:"name"`
		Pivot     data.Vec2 `json:"pivot"`
		Rect      data.Rect `json:"rect"`
		Border    data.Vec4 `json:"border"`
		Alignment int       `json:"alignment"`
	}

	sprites := make([]*spriteResponse, 0, 3000)
	var id int
	var textureId int
	var name string
	var pivot string
	var rect string
	var border string
	var alignment int
	for rows.Next() {
		err = rows.Scan(&id, &textureId, &name, &pivot, &rect, &border, &alignment)
		if err != nil {
			response.InternalError(w, err)
			return
		}
		sprite := &spriteResponse{
			Id:        id,
			TextureId: textureId,
			Name:      name,
			Alignment: alignment,
		}
		json.Unmarshal([]byte(pivot), &sprite.Pivot)
		json.Unmarshal([]byte(rect), &sprite.Rect)
		json.Unmarshal([]byte(border), &sprite.Border)
		sprites = append(sprites, sprite)
	}

	response.Json(w, sprites)
}

func CreateSprites(w http.ResponseWriter, r *http.Request) {
	state, err := common.InitState(r)
	if err != nil {
		response.ErrorWhileInitializing(w, err)
		return
	}
	ctx := state.Context
	conn := state.Conn

	type model struct {
		TextureId int
		Sprites   []data.Sprite
	}
	// Parse request body
	var req model
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Could not parse request body!", http.StatusBadRequest)
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
		data.InsertSpriteIfNotExists(ctx, batch, req.TextureId, &sprite)
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
}
