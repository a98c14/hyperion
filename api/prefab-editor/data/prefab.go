package data

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/a98c14/hyperion/api/asset"
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

type ColliderType int32

const (
	UndefinedCollider  ColliderType = 0
	RectColliderType   ColliderType = 1
	CircleColliderType ColliderType = 2
)

type Vec2 struct {
	X float32
	Y float32
}
type Vec3 struct {
	X float32
	Y float32
}

// TODO(selim): Add missing fields
type Renderer struct {
	IsVisible bool
}

type Transform struct {
	Position Vec3
	Scale    Vec2
	Rotation float32
}

// TODO(selim): Add tags
type CircleCollider struct {
	Center Vec2
	Radius float32
}

// TODO(selim): Add tags
type RectCollider struct {
	BL Vec2
	TR Vec2
}

// Data gets parsed accordin to the type
type Collider struct {
	Type ColliderType
	Data json.RawMessage
}

type ByIdPValue []PrefabModuleValueDB

func (b ByIdPValue) Id(i int) int { return b[i].ModulePartId }
func (b ByIdPValue) Len() int     { return len(b) }

type Prefab struct {
	Id        int                 `json:"id"`
	Name      string              `json:"name"`
	ParentId  int                 `json:"parentId"`
	Transform json.RawMessage     `json:"transform"`
	Renderer  json.RawMessage     `json:"renderer"`
	Colliders json.RawMessage     `json:"colliders"`
	Modules   []*PrefabModulePart `json:"modules"`
	Children  []*Prefab           `json:"children"`
}

type PrefabDB struct {
	Id       int
	Name     string
	ParentId sql.NullInt32

	// Stored as json values in database
	Transform json.RawMessage
	Renderer  json.RawMessage
	Colliders json.RawMessage
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
	err := conn.QueryRow(ctx, "select exists(select 1 from prefab p inner join asset a on p.asset_id=a.id where a.name=$1)", name).Scan(&exists)
	if err != nil {
		return false, xerrors.Wrap("DoesNameExist", err)
	}
	return exists, nil
}

func DoesIdExist(ctx context.Context, conn *pgxpool.Pool, id int) (bool, error) {
	var exists bool
	err := conn.QueryRow(ctx, "select exists(select 1 from prefab where id=$1)", id).Scan(&exists)
	if err != nil {
		return false, xerrors.Wrap("DoesNameExist", err)
	}
	return exists, nil
}

func UpdatePrefab(state common.State, prefabId int, name string, parentId sql.NullInt32, transform json.RawMessage, renderer json.RawMessage, colliders json.RawMessage) error {
	_, err := state.Conn.Exec(state.Context,
		`
		update asset
		set name=$2
		where id=(select asset_id from prefab where id=$1);
		`, prefabId, name)
	if err != nil {
		return xerrors.Wrap("InsertPrefab", err)
	}

	if parentId.Valid {
		_, err = state.Conn.Exec(state.Context,
			`
				update prefab
				set parent_id=$2,
					transform=$3,
					renderer=$4,
					colliders=$5
				where id=$1;
				`, prefabId, parentId, transform, renderer, colliders)
		if err != nil {
			return xerrors.Wrap("InsertPrefab", err)
		}
	} else {
		fmt.Println(colliders)
		_, err = state.Conn.Exec(state.Context,
			`
				update prefab
				set transform=$2,
					renderer=$3,
					colliders=$4
				where id=$1 and parent_id is null;
				`, prefabId, string(transform), string(renderer), string(colliders))
		if err != nil {
			return xerrors.Wrap("InsertPrefab", err)
		}
	}
	return nil
}

func InsertPrefab(ctx context.Context, conn *pgxpool.Pool, name string, parentId sql.NullInt32, transform json.RawMessage, renderer json.RawMessage, colliders json.RawMessage) (sql.NullInt32, error) {
	var id sql.NullInt32

	err := conn.QueryRow(ctx,
		`
		with ins as (
			insert into asset (name, unity_guid, type)
			select $1, '-', $3
			returning id
		)
		insert into "prefab" 
		(asset_id, parent_id, transform, renderer, colliders) 
		values((select id from ins), $2, $4, $5, $6) returning id`,
		name, parentId, asset.Prefab, transform, renderer, colliders).Scan(&id)
	if err != nil {
		return sql.NullInt32{}, xerrors.Wrap("InsertPrefab", err)
	}
	return id, nil
}

func DeletePrefab(state common.State, prefabId int) error {
	_, err := state.Conn.Exec(state.Context, `
		with recursive prefab_recursive as (
			select p.id, a.name, p.parent_id from prefab p
				inner join asset a on a.id=p.asset_id
				where p.id=$1 and parent_id is null
			union select c.id, ca.name, c.parent_id from prefab c 
				inner join asset ca on ca.id=c.asset_id 
				inner join prefab_recursive cp on cp.id=c.parent_id
		) 
		delete from prefab_module_part where prefab_id in (select id from prefab_recursive);`, prefabId)
	if err != nil {
		return err
	}
	_, err = state.Conn.Exec(state.Context, `
		with recursive prefab_recursive as (
			select p.id, a.name, p.parent_id from prefab p
				inner join asset a on a.id=p.asset_id
				where p.id=$1 and parent_id is null
			union select c.id, ca.name, c.parent_id from prefab c 
				inner join asset ca on ca.id=c.asset_id 
				inner join prefab_recursive cp on cp.id=c.parent_id
		) 
		delete from prefab where id in (select id from prefab_recursive);`, prefabId)
	if err != nil {
		return err
	}
	return nil
}

// Returns prefabs that have no parent
func GetRootPrefabs(state common.State) ([]RootPrefab, error) {
	rows, err := state.Conn.Query(state.Context, `select p.id, a.name from prefab p inner join asset a on a.id=p.asset_id where parent_id is null`)
	if err != nil {
		return nil, xerrors.WrapMsg("GetRootPrefabs", "Query", err)
	}
	defer rows.Close()
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
		select p.id, a.name, p.parent_id from prefab p
			inner join asset a on a.id=p.asset_id
			where p.id=$1 and parent_id is null
		union select c.id, ca.name, c.parent_id from prefab c 
			inner join asset ca on ca.id=c.asset_id 
			inner join prefab_recursive cp on cp.id=c.parent_id
	) select pr.id, pr.name, pr.parent_id, p.transform, p.renderer, p.colliders from prefab_recursive pr 
	  inner join prefab p on pr.id=p.id;`, prefabId)
	if err != nil {
		return nil, xerrors.Wrap("GetPrefabById", err)
	}
	defer rows.Close()

	dbPrefabs := make([]*PrefabDB, 0)
	for rows.Next() {
		prefab := PrefabDB{}
		err = rows.Scan(&prefab.Id, &prefab.Name, &prefab.ParentId, &prefab.Transform, &prefab.Renderer, &prefab.Colliders)
		if err != nil {
			return nil, xerrors.Wrap("GetPrefabById", err)
		}

		dbPrefabs = append(dbPrefabs, &prefab)
	}
	defer rows.Close()

	// For each prefab, get the module tree and values
	prefabMap := make(map[int]*Prefab, 0)
	var rootId int
	for _, dbPrefab := range dbPrefabs {
		prefab, err := getPrefabDetails(state, balanceVersionId, dbPrefab)
		if err != nil {
			return nil, xerrors.Wrap("GetPrefabById", err)
		}

		prefabMap[prefab.Id] = prefab
		if prefab.ParentId != 0 {
			parent := prefabMap[prefab.ParentId]
			parent.Children = append(parent.Children, prefab)
		} else {
			rootId = prefab.Id
		}
	}

	if rootId > 0 {
		return prefabMap[rootId], nil
	} else {
		return nil, xerrors.ErrNotFound
	}
}

func getPrefabDetails(state common.State, balanceVersionId int, prefab *PrefabDB) (*Prefab, error) {
	parentId := 0
	if prefab.ParentId.Valid {
		parentId = int(prefab.ParentId.Int32)
	}
	result := Prefab{
		Id:        prefab.Id,
		ParentId:  parentId,
		Name:      prefab.Name,
		Transform: prefab.Transform,
		Renderer:  prefab.Renderer,
		Colliders: prefab.Colliders,
		Modules:   make([]*PrefabModulePart, 0),
		Children:  make([]*Prefab, 0),
	}

	// Load prefab module part values. These are the actual values entered from
	// editor app.
	rows, err := state.Conn.Query(state.Context,
		`select pmp.id, pmp.array_index, pmp.module_part_id, mp.value_type, pmp.value, pmp.prefab_id
		from prefab_module_part pmp
		inner join module_part mp on mp.id = pmp.module_part_id
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
	instr, params := querystr.GenerateInStringIdentifiable(ByIdPValue(prefabModuleValues), 0)
	rows, err = state.Conn.Query(state.Context, `with recursive module_part_recursive as (
			select module_part.id, name, value_type, parent_id from module_part 
			where module_part.id in (`+instr+`)
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
				fmt.Println(v.Value)
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
