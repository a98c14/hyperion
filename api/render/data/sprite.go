package data

import (
	"context"
	"encoding/json"

	"github.com/jackc/pgx/v4"
)

type Sprite struct {
	Name       string          `json:"name"`
	SpriteId   string          `json:"spriteId"`
	InternalId string          `json:"internalId"`
	Pivot      json.RawMessage `json:"pivot"`
	Border     json.RawMessage `json:"border"`
	Rect       json.RawMessage `json:"rect"`
	Alignment  int             `json:"alignment"`
}

func InsertSpriteIfNotExists(ctx context.Context, batch *pgx.Batch, textureId int, sprite *Sprite) {
	batch.Queue(
		`insert into "sprite"
		(texture_id, 
		 unity_sprite_id, 
		 unity_internal_id,
		 unity_name, 
		 unity_pivot, 
		 unity_rect,
		 unity_border,
		 unity_alignment)
		values($1, $2, $3, $4, $5, $6, $7, $8)
		on conflict do nothing`,
		textureId, sprite.SpriteId, sprite.InternalId, sprite.Name, sprite.Pivot, sprite.Rect, sprite.Border, sprite.Alignment)
}
