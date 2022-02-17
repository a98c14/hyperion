package data

import (
	"encoding/json"
	"fmt"

	"github.com/a98c14/hyperion/common"
	xerrors "github.com/a98c14/hyperion/common/errors"
	"github.com/a98c14/hyperion/common/querystr"
	"github.com/jackc/pgx/v4"
)

type PrefabModulePartValue struct {
	ArrayIndex   int             `json:"arrayIndex"`
	ModulePartId int             `json:"modulePartId"`
	Value        json.RawMessage `json:"value"`
}

type ByIdPMPValue []PrefabModulePartValue

func (b ByIdPMPValue) Id(i int) int { return b[i].ModulePartId }
func (b ByIdPMPValue) Len() int     { return len(b) }

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

func UpdatePrefabModulePartValues(state common.State, prefabId int, versionId int, values []PrefabModulePartValue) error {
	batch := &pgx.Batch{}
	for _, v := range values {
		batch.Queue(`
				update prefab_module_part
				set value=$2::json
				where prefab_id=$3 and module_part_id=$4::int and balance_version_id=$5 and array_index=$1;
				`,
			v.ArrayIndex, v.Value, prefabId, v.ModulePartId, versionId)

		batch.Queue(`
			insert into prefab_module_part
			(array_index, value, prefab_id, module_part_id, balance_version_id)
			select $1, $2::json, $3, $4::int, $5
			where not exists (
				select module_part_id from prefab_module_part 
				where prefab_id=$3 and module_part_id=$4::int and balance_version_id=$5 and array_index=$1);
			`,
			v.ArrayIndex, v.Value, prefabId, v.ModulePartId, versionId)
	}

	instr, params := querystr.GenerateInStringIdentifiable(ByIdPMPValue(values), 2)
	params = append([]interface{}{prefabId, versionId}, params...)
	batch.Queue(`delete from prefab_module_part where prefab_id=$1 and balance_version_id=$2 and module_part_id not in (`+instr+`)`, params...)
	br := state.Conn.SendBatch(state.Context, batch)
	_, err := br.Exec()
	if err != nil {
		fmt.Println("ERrror")
		return err
	}
	br.Close()
	return nil
}
