package data

import (
	"context"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type PrefabModulePartValue struct {
	ArrayIndex       int
	ValueType        int
	ModulePartId     int
	BalanceVersionId int
	Value            string
}

func InsertPrefabModulePartValues(ctx context.Context, conn *pgxpool.Pool, prefabId int, values []PrefabModulePartValue) error {
	batch := &pgx.Batch{}
	for _, v := range values {
		batch.Queue(`insert into "prefab_module_part"
				(array_index, value_type, value, prefab_id, module_part_id, balance_version_id)
				values($1, $2, $3, $4, $5, $6) 
				on conflict do nothing`,
			v.ArrayIndex, v.ValueType, v.Value, prefabId, v.ModulePartId, v.BalanceVersionId)
	}
	br := conn.SendBatch(ctx, batch)
	_, err := br.Exec()
	if err != nil {
		return err
	}

	br.Close()
	return nil
}
