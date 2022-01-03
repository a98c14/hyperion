package query

import (
	"context"
	"database/sql"

	"github.com/a98c14/hyperion/model/prefab"
	"github.com/jackc/pgx/v4/pgxpool"
)

func DoesPrefabWithNameExist(ctx context.Context, conn *pgxpool.Pool, name string) (bool, error) {
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

func InsertPrefabModulePartValue(ctx context.Context, conn *pgxpool.Pool, prefabId int, balanceVersionId int, value *prefab.ModulePartValue) (sql.NullInt32, error) {
	var id sql.NullInt32
	// err := conn.QueryRow(ctx,
	// 	`insert into "prefab_module_part"
	// 	(array_index, value_type, numeric_value, string_value, prefab_id, module_part_id, balance_version_id)
	// 	values($1, $2) returning id`,
	// 	value.ArrayIndex, value.ValueType, parentId).Scan(&id)
	return id, nil
}
