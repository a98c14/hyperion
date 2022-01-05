package data

import (
	"context"
	"database/sql"

	"github.com/jackc/pgx/v4/pgxpool"
)

func DoesNameExist(ctx context.Context, conn *pgxpool.Pool, name string) (bool, error) {
	var exists bool
	err := conn.QueryRow(ctx, "select exists(select 1 from prefab where name=$1)", name).Scan(&exists)
	return exists, err
}

func InsertPrefab(ctx context.Context, conn *pgxpool.Pool, name string, parentId sql.NullInt32) (sql.NullInt32, error) {
	var id sql.NullInt32
	err := conn.QueryRow(ctx,
		`insert into "prefab_part" 
		 (name, parent_id) 
		 values($1, $2) returning id`,
		name, parentId).Scan(&id)
	return id, err
}
