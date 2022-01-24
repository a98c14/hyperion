package data

import (
	"context"
	"encoding/json"

	"github.com/a98c14/hyperion/api/asset"
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

func InsertSpriteIfNotExists(ctx context.Context, batch *pgx.Batch, textureId int, sprite *Sprite, asset *asset.AssetDb) {
	batch.Queue(
		`
		with ins as (
			insert into asset (name, unity_guid, type)
			select $9, $10, $11
			where not exists (select name from asset where name=$9)
			returning id
		)
		insert into "sprite"
		(texture_id, unity_sprite_id, unity_internal_id, unity_name, unity_pivot, unity_rect, unity_border, unity_alignment, asset_id)
		 values($1, $2, $3, $4, $5, $6, $7, $8, )
		 on conflict do nothing`,
		textureId,
		sprite.SpriteId,
		sprite.InternalId,
		sprite.Name,
		sprite.Pivot,
		sprite.Rect,
		sprite.Border,
		sprite.Alignment,
		asset.Name,
		asset.UnityGuid,
		asset.Type)
}
