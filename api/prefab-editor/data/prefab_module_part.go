package data

import (
	"encoding/json"
	"fmt"

	"github.com/a98c14/hyperion/common"
	xerrors "github.com/a98c14/hyperion/common/errors"
	"github.com/jackc/pgx/v4"
)

type PrefabModulePartValue struct {
	ArrayIndex   int             `json:"arrayIndex"`
	ModulePartId int             `json:"modulePartId"`
	Value        json.RawMessage `json:"value"`
}

func InsertPrefabModulePartValues(state common.State, prefabId int, versionId int, values []PrefabModulePartValue) error {
	batch := &pgx.Batch{}
	for _, v := range values {
		fmt.Println(v)
		batch.Queue(`insert into "prefab_module_part"
				(array_index, value, prefab_id, module_part_id, balance_version_id)
				values($1, $2, $3, $4, $5)
				on conflict do nothing`,
			v.ArrayIndex, v.Value, prefabId, v.ModulePartId, versionId)
	}
	br := state.Conn.SendBatch(state.Context, batch)
	_, err := br.Exec()
	if err != nil {
		return xerrors.Wrap("InsertPrefabModulePartValues", err)
	}

	br.Close()
	return nil
}
