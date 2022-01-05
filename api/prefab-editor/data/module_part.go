package data

import (
	"database/sql"
	"errors"

	"github.com/a98c14/hyperion/common"
)

type ModulePart struct {
	Id        int
	Name      string
	ValueType int
	ParentId  int
}

type RootModule struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

func GetModuleParts(state common.State, moduleId int) ([]*ModulePart, error) {
	rows, err := state.Conn.Query(state.Context, `with recursive module_part_recursive as (
		select id, name, value_type, parent_id from module_part
		where id=$1 and parent_id is null
		union select c.id, c.name, c.value_type, c.parent_id from module_part c inner join module_part_recursive cp on cp.id=c.parent_id 
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
	for rows.Next() {
		err = rows.Scan(&moduleId, &name, &valueType, &parentId)
		if err != nil {
			return nil, err
		}
		pid := 0
		if parentId.Valid {
			pid = int(parentId.Int32)
		}
		moduleParts = append(moduleParts, &ModulePart{Id: id, Name: name, ValueType: valueType, ParentId: pid})
	}
	componentCount := rows.CommandTag().RowsAffected()
	if componentCount == 0 {
		return nil, errors.New("given module has no part attached")
	}
	return moduleParts, nil
}

func GetRootModuleParts(state common.State) ([]RootModule, error) {
	rows, err := state.Conn.Query(state.Context, `select id, name from "module_part" where parent_id is null`)
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
