package data

import (
	"context"
	"encoding/json"

	xerrors "github.com/a98c14/hyperion/common/errors"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type PrefabModulePartValue struct {
	ArrayIndex   int
	ValueType    int
	ModulePartId int
	Value        json.RawMessage
}

func InsertPrefabModulePartValues(ctx context.Context, conn *pgxpool.Pool, prefabId int, versionId int, values []PrefabModulePartValue) error {
	batch := &pgx.Batch{}
	for _, v := range values {
		batch.Queue(`insert into "prefab_module_part"
				(array_index, value_type, value, prefab_id, module_part_id, balance_version_id)
				values($1, $2, $3, $4, $5, $6) 
				on conflict do nothing`,
			v.ArrayIndex, v.ValueType, v.Value, prefabId, v.ModulePartId, versionId)
	}
	br := conn.SendBatch(ctx, batch)
	_, err := br.Exec()
	if err != nil {
		return xerrors.Wrap("InsertPrefabModulePartValues", err)
	}

	br.Close()
	return nil
}
