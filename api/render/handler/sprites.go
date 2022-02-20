package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/a98c14/hyperion/api/asset"
	"github.com/a98c14/hyperion/api/render/data"
	"github.com/a98c14/hyperion/common"
	"github.com/a98c14/hyperion/common/response"
	"github.com/jackc/pgx/v4"
)

func GetSprites(state common.State, w http.ResponseWriter, r *http.Request) error {
	rows, err := state.Conn.Query(state.Context,
		`select asset.id, 
				texture_id, 
				asset.name, 
				unity_pivot, 
				unity_rect, 
				unity_border, 
				unity_alignment, 
				unity_sprite_id,
				asset.id,
				asset.unity_internal_id,
				asset.unity_guid
		 from sprite
		 inner join asset on sprite.asset_id=asset.id `)
	if err != nil {
		return err
	}
	defer rows.Close()
	type resp struct {
		Id              int       `json:"id"`
		TextureId       int       `json:"textureId"`
		Name            string    `json:"name"`
		Pivot           data.Vec2 `json:"pivot"`
		Rect            data.Rect `json:"rect"`
		Border          data.Vec4 `json:"border"`
		Alignment       int       `json:"alignment"`
		SpriteId        string    `json:"spriteId"`
		AssetId         int       `json:"assetId"`
		UnityInternalId int64     `json:"unityInternalId"`
		UnityGuid       string    `json:"unityGuid"`
	}

	sprites := make([]*resp, 0, 3000)
	var pivot string
	var rect string
	var border string
	for rows.Next() {
		sprite := resp{}
		err = rows.Scan(&sprite.Id, &sprite.TextureId, &sprite.Name, &pivot, &rect, &border, &sprite.Alignment, &sprite.SpriteId, &sprite.AssetId, &sprite.UnityInternalId, &sprite.UnityGuid)
		if err != nil {
			return err
		}
		json.Unmarshal([]byte(pivot), &sprite.Pivot)
		json.Unmarshal([]byte(rect), &sprite.Rect)
		json.Unmarshal([]byte(border), &sprite.Border)
		sprites = append(sprites, &sprite)
	}

	response.Json(w, sprites)
	return nil
}

func CreateSprites(state common.State, w http.ResponseWriter, r *http.Request) error {
	ctx := state.Context
	conn := state.Conn

	// Parse request body

	req := struct {
		TextureId int
		Sprites   []data.Sprite
	}{}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return err
	}

	var exists bool
	err = conn.QueryRow(ctx, "select exists(select 1 from texture where id=$1)", req.TextureId).Scan(&exists)
	if err != nil || !exists {
		return err
	}

	// Start transaction. If all components can not be added successfully, don't
	// insert anything
	tx, err := conn.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	batch := &pgx.Batch{}
	for _, sprite := range req.Sprites {
		asset := asset.AssetDb{
			UnityGuid:       "0",
			UnityInternalId: sprite.InternalId,
			Name:            sprite.Name,
			Type:            asset.Sprite,
		}
		data.InsertSpriteIfNotExists(ctx, batch, req.TextureId, &sprite, &asset)
	}
	br := conn.SendBatch(ctx, batch)
	ct, err := br.Exec()
	if err != nil {
		return err
	}
	fmt.Printf("Inserted rows: %d", ct.RowsAffected())
	err = tx.Commit(ctx)
	if err != nil {
		return err
	}

	return nil
}
