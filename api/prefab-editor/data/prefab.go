package data

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/a98c14/hyperion/common"
	xerrors "github.com/a98c14/hyperion/common/errors"
	"github.com/a98c14/hyperion/common/querystr"
	"github.com/jackc/pgx/v4/pgxpool"
)

type ValuePart struct {
	Id    int
	Name  string
	Value string
}

type PrefabModuleValueDB struct {
	Id           int
	Name         string
	ModulePartId int
	PrefabId     int
	ValueType    int
	ArrayIndex   int
	Value        json.RawMessage
}

type ByIdPValue []PrefabModuleValueDB

func (b ByIdPValue) Id(i int) int { return b[i].Id }
func (b ByIdPValue) Len() int     { return len(b) }

type Prefab struct {
	Id       int                 `json:"id"`
	Name     string              `json:"name"`
	ParentId int                 `json:"parentId"`
	Modules  []*PrefabModulePart `json:"modules"`
}

type PrefabDB struct {
	Id       int
	Name     string
	ParentId sql.NullInt32
}

type RootPrefab struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type PrefabModulePart struct {
	Id        int                 `json:"id"`
	Name      string              `json:"name"`
	ParentId  int                 `json:"parentId"`
	ValueType int                 `json:"valueType"`
	Value     json.RawMessage     `json:"value"`
	Children  []*PrefabModulePart `json:"children"`
}

type ById []*Prefab

func (b ById) Id(i int) int { return b[i].Id }
func (b ById) Len() int     { return len(b) }

func DoesNameExist(ctx context.Context, conn *pgxpool.Pool, name string) (bool, error) {
	var exists bool
	err := conn.QueryRow(ctx, "select exists(select 1 from prefab where name=$1)", name).Scan(&exists)
	if err != nil {
		return false, xerrors.Wrap("DoesNameExist", err)
	}
	return exists, nil
}

func InsertPrefab(ctx context.Context, conn *pgxpool.Pool, name string, parentId sql.NullInt32) (sql.NullInt32, error) {
	var id sql.NullInt32
	err := conn.QueryRow(ctx,
		`insert into "prefab" 
		 (name, parent_id) 
		 values($1, $2) returning id`,
		name, parentId).Scan(&id)
	if err != nil {
		return sql.NullInt32{}, xerrors.Wrap("InsertPrefab", err)
	}
	return id, nil
}

// Returns prefabs that have no parent
func GetRootPrefabs(state common.State) ([]RootPrefab, error) {
	rows, err := state.Conn.Query(state.Context, `select id, name from prefab where parent_id is null`)
	if err != nil {
		return nil, xerrors.WrapMsg("GetRootPrefabs", "Query", err)
	}
	prefabs := make([]RootPrefab, 0)
	for rows.Next() {
		prefab := RootPrefab{}
		err = rows.Scan(&prefab.Id, &prefab.Name)
		if err != nil {
			return nil, xerrors.WrapMsg("GetRootPrefabs", "RowScan", err)
		}
		prefabs = append(prefabs, prefab)
	}
	return prefabs, nil
}

func GetPrefabById(state common.State, prefabId int, balanceVersionId int) (*Prefab, error) {
	// Fetch all prefabs that have given prefabId as root
	rows, err := state.Conn.Query(state.Context, `with recursive prefab_recursive as (
		select id, name, parent_id from prefab
		where id=$1 and parent_id is null
		union select c.id, c.name, c.parent_id from prefab c inner join prefab_recursive cp on cp.id=c.parent_id
	) select id, name, parent_id from prefab_recursive;`, prefabId)
	if err != nil {
		return nil, xerrors.Wrap("GetPrefabById", err)
	}

	dbPrefabs := make([]*PrefabDB, 0)
	for rows.Next() {
		prefab := PrefabDB{}
		err = rows.Scan(&prefab.Id, &prefab.Name, &prefab.ParentId)
		if err != nil {
			return nil, xerrors.Wrap("GetPrefabById", err)
		}

		dbPrefabs = append(dbPrefabs, &prefab)
	}
	defer rows.Close()

	// For each prefab, get the module tree and values
	prefabs := make([]*Prefab, 0, len(dbPrefabs))
	for _, dbPrefab := range dbPrefabs {
		prefab, err := getPrefabDetails(state, balanceVersionId, dbPrefab)
		if err != nil {
			return nil, xerrors.Wrap("GetPrefabById", err)
		}
		prefabs = append(prefabs, prefab)
	}

	// TODO(selim): Child prefabs should be inside the original prefab as children
	return prefabs[0], nil
}

func getPrefabDetails(state common.State, balanceVersionId int, prefab *PrefabDB) (*Prefab, error) {
	parentId := 0
	if prefab.ParentId.Valid {
		parentId = int(prefab.ParentId.Int32)
	}
	result := Prefab{
		Id:       prefab.Id,
		ParentId: parentId,
		Name:     prefab.Name,
		Modules:  make([]*PrefabModulePart, 0),
	}

	// Load prefab module part values. These are the actual values entered from
	// editor app.
	rows, err := state.Conn.Query(state.Context,
		`select module_part.id, array_index, module_part_id, module_part.
		value_type, value, prefab_id
		from prefab_module_part 
		inner join module_part on module_part.id = prefab_module_part.module_part_id
		where balance_version_id=$1 and prefab_id=$2`, balanceVersionId, prefab.Id)
	if err != nil {
		return nil, xerrors.Wrap("getPrefabDetails", err)
	}
	defer rows.Close()

	prefabModuleValues := make([]PrefabModuleValueDB, 0)
	modulePartValueMap := make(map[int]*PrefabModuleValueDB)
	for rows.Next() {
		pvalue := PrefabModuleValueDB{}
		err = rows.Scan(&pvalue.Id, &pvalue.ArrayIndex, &pvalue.ModulePartId, &pvalue.ValueType, &pvalue.Value, &pvalue.PrefabId)
		if err != nil {
			return nil, xerrors.Wrap("getPrefabDetails", err)
		}

		prefabModuleValues = append(prefabModuleValues, pvalue)
		modulePartValueMap[pvalue.ModulePartId] = &pvalue

	}

	// Load prefab module trees.
	instr, params := querystr.GenerateInStringIdentifiable(ByIdPValue(prefabModuleValues))
	rows, err = state.Conn.Query(state.Context, `with recursive module_part_recursive as (
			select id, name, value_type, parent_id from module_part 
			where id in (`+instr+`)
			union select mp.id, mp.name, mp.value_type, mp.parent_id from module_part mp 
			inner join module_part_recursive mpr on mp.id=mpr.parent_id
		) select id, name, value_type, parent_id from module_part_recursive;`, params...)
	if err != nil {
		return nil, xerrors.Wrap("getPrefabDetails", err)
	}
	defer rows.Close()

	// Stores every children that has `key` as parent
	childMap := make(map[int][]*PrefabModulePart)
	processQueue := make(chan *PrefabModulePart, 500)
	var parentIdSql sql.NullInt32
	for rows.Next() {
		modulePart := PrefabModulePart{}
		rows.Scan(&modulePart.Id, &modulePart.Name, &modulePart.ValueType, &parentIdSql)

		// Store the parentId in map
		if parentIdSql.Valid {
			bucket := childMap[int(parentIdSql.Int32)]
			bucket = append(bucket, &modulePart)
			childMap[int(parentIdSql.Int32)] = bucket
			modulePart.ParentId = int(parentIdSql.Int32)
		} else {
			// Add every root module to process queue for later processing
			processQueue <- &modulePart

			// If a module has no parent id, it means it is a root node and should be
			// added to the result prefab module list. All other parts will be attached
			// to these root modules.
			result.Modules = append(result.Modules, &modulePart)
		}
	}

	for m := range processQueue {
		if children, ok := childMap[m.Id]; ok {
			for _, child := range children {
				processQueue <- child
				m.Children = append(m.Children, child)
			}
		} else {
			// If current module has no child it means it is a leaf/value node
			// and value is set.
			v, ok := modulePartValueMap[m.Id]
			if ok {
				m.Value = json.RawMessage(v.Value)
			} else {
				m.Value = nil
			}
		}

		if len(processQueue) == 0 {
			close(processQueue)
		}
	}

	return &result, nil
}
