package data

import (
	"context"
	"database/sql"

	"github.com/jackc/pgx/v4/pgxpool"
)

func InsertTexture(ctx context.Context, conn *pgxpool.Pool, path string, guid string, name string) (sql.NullInt32, error) {
	var id sql.NullInt32
	err := conn.QueryRow(ctx,
		`insert into "texture" 
		 (image_path, unity_guid, unity_name) 
		 values($1, $2, $3) returning id`,
		path, guid, name).Scan(&id)
	return id, err
}
