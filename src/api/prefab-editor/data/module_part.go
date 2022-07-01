package data

import (
	"database/sql"
	"encoding/json"
	"errors"

	"github.com/a98c14/hyperion/common"
	e "github.com/a98c14/hyperion/common/errors"
	"github.com/jackc/pgx/v4"
)

type ModulePartNode struct {
	Name      string
	ValueType int
	Tooltip   string
	IsArray   bool
	ParentId  sql.NullInt32
	Value     json.RawMessage
}

type ModulePart struct {
	Id         int
	Name       string
	ValueType  int
	Tooltip    string
	IsArray    bool
	ParentId   int
	ParentName string
}

type ModulePartDB struct {
	Id        int
	Name      string
	ValueType int
	Tooltip   string
	IsArray   bool
	ParentId  sql.NullInt32
}

type RootModule struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type ModulePartInfo struct {
	ValueType int
	Tooltip   string
	IsArray   bool
	Children  json.RawMessage
}

func DoesModulePartExist(state common.State, moduleId int) (bool, error) {
	var exists bool
	err := state.Conn.QueryRow(state.Context, "select exists(select 1 from module_part where id=$1 and deleted_date is null)", moduleId).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func DoesModulePartWithNameExist(state common.State, moduleName string) (bool, error) {
	var exists bool
	err := state.Conn.QueryRow(state.Context, "select exists(select 1 from module_part where name=$1 and deleted_date is null)", moduleName).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

// Returns the Id of the module for given `name` and `parent_id`. `parent_id` is used because
// names are only unique between siblings
func GetModulePartIdWithName(state common.State, moduleName string, parentId sql.NullInt32) (sql.NullInt32, error) {
	var id sql.NullInt32
	var err error

	// TODO(selim): Is there a better way to do this?
	if parentId.Valid {
		err = state.Conn.QueryRow(state.Context,
			`select id from module_part 
			where name=$1 and parent_id=$2 and deleted_date is null`, moduleName, parentId.Int32).Scan(&id)
	} else {
		err = state.Conn.QueryRow(state.Context,
			`select id from module_part 
			where name=$1 and parent_id is null and deleted_date is null`, moduleName).Scan(&id)
	}

	if err != nil && err == pgx.ErrNoRows {
		return id, nil
	} else if err != nil {
		return id, err
	}

	return id, nil
}

// Creates module hashmap for all module parts for given root module id.
// `Module Parent Name+Module Name` is used as key since module names are only unique
// between siblings.
func GetModulePartMap(state common.State, moduleName string) (map[string]*ModulePart, error) {
	rows, err := state.Conn.Query(state.Context, `
		with recursive module_part_recursive as (
			select 
				id, 
				name, 
				value_type, 
				parent_id, 
				case when parent_id is not null then name else null end as parent_name,
				is_array,
				tooltip
			from module_part
			where name=$1 and parent_id is null and deleted_date is null
			union select 
				c.id, 
				c.name, 
				c.value_type, 
				c.parent_id, 
				(select name from module_part where id=cp.id),
				c.is_array,
				c.tooltip
			from module_part c inner join module_part_recursive cp on cp.id=c.parent_id
			where deleted_date is null)
		select * from module_part_recursive;`, moduleName)

	if err != nil {
		return nil, e.Wrap("GetModulePartMap", err)
	}
	defer rows.Close()

	modulePartMap := make(map[string]*ModulePart)
	var parentId sql.NullInt32
	var parentName, tooltip sql.NullString
	for rows.Next() {
		module := &ModulePart{}
		err = rows.Scan(&module.Id, &module.Name, &module.ValueType, &parentId, &parentName, &module.IsArray, &tooltip)
		if err != nil {
			return nil, e.Wrap("GetModulePartMap", err)
		}

		if tooltip.Valid {
			module.Tooltip = tooltip.String
		}

		if parentId.Valid && parentName.Valid {
			module.ParentId = int(parentId.Int32)
			module.ParentName = parentName.String
			modulePartMap[GetModulePartKey(module.ParentName, module.Name)] = module
		} else {
			modulePartMap[module.Name] = module
		}
	}

	moduleCount := rows.CommandTag().RowsAffected()
	if moduleCount == 0 {
		return nil, errors.New("given module has no part attached")
	}
	return modulePartMap, nil
}

func GetModuleParts(state common.State, moduleId int) ([]*ModulePart, error) {
	rows, err := state.Conn.Query(state.Context, `with recursive module_part_recursive as (
		select id, name, value_type, parent_id, is_array from module_part
		where id=$1 and parent_id is null and deleted_date is null
		union select c.id, c.name, c.value_type, c.parent_id, c.is_array from module_part c
		inner join module_part_recursive cp on cp.id=c.parent_id 
		where deleted_date is null 
	) select * from module_part_recursive;`, moduleId)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	moduleParts := make([]*ModulePart, 0, 10)
	var id int
	var name string
	var valueType int
	var parentId sql.NullInt32
	var isArray bool
	for rows.Next() {
		err = rows.Scan(&id, &name, &valueType, &parentId, &isArray)
		if err != nil {
			return nil, err
		}
		pid := 0
		if parentId.Valid {
			pid = int(parentId.Int32)
		}
		moduleParts = append(moduleParts, &ModulePart{Id: id, Name: name, ValueType: valueType, ParentId: pid, IsArray: isArray})
	}
	moduleCount := rows.CommandTag().RowsAffected()
	if moduleCount == 0 {
		return nil, errors.New("given module has no part attached")
	}
	return moduleParts, nil
}

func GetRootModuleParts(state common.State) ([]RootModule, error) {
	rows, err := state.Conn.Query(state.Context, `select id, name from "module_part" 
		where parent_id is null and deleted_date is null 
		order by name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	rootModules := make([]RootModule, 0, 100)

	var id int
	var name string
	for rows.Next() {
		err = rows.Scan(&id, &name)
		if err != nil {
			return nil, err
		}
		rootModules = append(rootModules, RootModule{Id: id, Name: name})
	}
	return rootModules, nil
}

// Delets a module part node and all of its children from database
func DeleteModulePartTree(state common.State, id int) error {
	_, err := state.Conn.Exec(state.Context, `
		with recursive module_part_recursive as (
			select 
				id, 
				parent_id
			from module_part
			where id=$1 and deleted_date is null
			union select 
				c.id, 
				c.parent_id
			from module_part c 
			inner join module_part_recursive cp on cp.id=c.parent_id
			where deleted_date is null
		) update module_part set deleted_date=now() where id in (select id from module_part_recursive);`, id)

	if err != nil {
		return e.Wrap("DeleteModulePartTree", err)
	}
	return nil
}

func UpdateModulePart(state common.State, id int, node *ModulePartNode) error {
	_, err := state.Conn.Exec(state.Context, `update module_part set is_array=$2, value_type=$3, tooltip=$4 where id=$1`, id, node.IsArray, node.ValueType, node.Tooltip)
	if err != nil {
		return e.Wrap("DeleteModulePartTree", err)
	}
	return nil
}

// Inserts a module part node with all of its children to database
func InsertModulePartTree(state common.State, node *ModulePartNode) error {
	// Start transaction. If all modules can not be added successfully, don't
	// insert anything
	tx, err := state.Conn.Begin(state.Context)
	if err != nil {
		return e.Wrap("InsertModulePartTree", err)
	}
	defer tx.Rollback(state.Context)
	c := make(chan *ModulePartNode, 500)
	c <- node

	// Process nodes in json object tree
	for n := range c {
		// Insert current node to database and store its id
		id, err := InsertModulePart(state, n)
		if err != nil {
			return e.Wrap("InsertModulePartTree", err)
		}

		// Check if current value is a json object
		m := make(map[string]json.RawMessage, 20)
		err = json.Unmarshal(n.Value, &m)

		if err != nil || m == nil {
			// If there is no more elements to process, close the channel
			if len(c) == 0 {
				close(c)
			}
			continue
		}

		// Add all values to process channel
		for k := range m {
			var pi ModulePartInfo
			err = json.Unmarshal(m[k], &pi)

			if err != nil {
				return e.Wrap("InsertModulePartTree", err)
			}
			c <- &ModulePartNode{
				ParentId:  id,
				Name:      k,
				ValueType: pi.ValueType,
				Value:     pi.Children,
				Tooltip:   pi.Tooltip,
				IsArray:   pi.IsArray,
			}
		}
	}

	// Commit transaction
	err = tx.Commit(state.Context)
	if err != nil {
		return e.Wrap("InsertModulePartTree", err)
	}

	return nil
}

// Inserts a single module part to database
func InsertModulePart(state common.State, node *ModulePartNode) (sql.NullInt32, error) {
	var id sql.NullInt32
	err := state.Conn.QueryRow(state.Context,
		`insert into "module_part" 
		 (name, value_type, parent_id) 
		 values($1, $2, $3) returning id`,
		node.Name, node.ValueType, node.ParentId).Scan(&id)
	return id, err
}

func GetModulePartKey(parent string, child string) string {
	return parent + "." + child
}
