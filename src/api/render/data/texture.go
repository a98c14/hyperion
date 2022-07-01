package data

import (
	"database/sql"

	"github.com/a98c14/hyperion/api/asset"
	"github.com/a98c14/hyperion/common"
)

func InsertTexture(state common.State, path string, guid string, internalId int64, name string) (sql.NullInt32, error) {
	var id sql.NullInt32
	err := state.Conn.QueryRow(state.Context,
		`
		with ins as (
			insert into asset (name, unity_guid, unity_internal_id, type)
			values ($3, $2, $4, $5)
			returning id
		)
		insert into "texture" (image_path, asset_id) 
		values($1, (select id from ins)) returning id`,
		path, guid, name, internalId, asset.Texture).Scan(&id)
	return id, err
}
