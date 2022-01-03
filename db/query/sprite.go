package query

import (
	"context"

	"github.com/a98c14/hyperion/model/render"
	"github.com/jackc/pgx/v4"
)

func InsertSpriteIfNotExists(ctx context.Context, batch *pgx.Batch, textureId int, sprite *render.Sprite) {
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
