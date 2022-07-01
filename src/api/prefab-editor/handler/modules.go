package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/a98c14/hyperion/api/prefab-editor/data"
	"github.com/a98c14/hyperion/common"
	"github.com/a98c14/hyperion/common/response"
	"github.com/go-chi/chi/v5"
)

/* Gets the base module_part and all of its children by id.
Can only query base components. Nodes that have null as children
represent leaf nodes */
func GetModuleById(state common.State, w http.ResponseWriter, r *http.Request) error {
	// Load components that have given name from database
	moduleIdString := chi.URLParam(r, "moduleId")
	moduleId, err := strconv.Atoi(moduleIdString)
	if err != nil {
		return err
	}

	moduleParts, err := data.GetModuleParts(state, moduleId)
	if err != nil {
		return err
	}

	type modulePart struct {
		Id        int           `json:"id"`
		Name      string        `json:"name"`
		ValueType int           `json:"valueType"`
		Children  []*modulePart `json:"children"`
		IsArray   bool          `json:"isArray"`
	}

	/* Create object tree from module part list. Components are always ordered from root to child*/
	root := modulePart{
		Id:        moduleParts[0].Id,
		Name:      moduleParts[0].Name,
		ValueType: moduleParts[0].ValueType,
		IsArray:   moduleParts[0].IsArray,
		Children:  make([]*modulePart, 0),
	}
	nodeMap := make(map[int]*modulePart)
	nodeMap[root.Id] = &root
	for _, c := range moduleParts[1:] {
		if val, ok := nodeMap[c.ParentId]; ok {
			cr := modulePart{
				Id:        c.Id,
				Name:      c.Name,
				ValueType: c.ValueType,
				IsArray:   c.IsArray,
				Children:  nil,
			}
			val.Children = append(val.Children, &cr)
			nodeMap[c.Id] = &cr
		} else {
			return err
		}
	}

	response.Json(w, &root)
	return nil
}

func ListComponents(w http.ResponseWriter, r *http.Request) {

}

// Returns all the root components.
// Root components are module_part that have no parent. Used to group
// other components and editor can only filter using root components
func GetRootModules(state common.State, w http.ResponseWriter, r *http.Request) error {
	components, err := data.GetRootModuleParts(state)
	if err != nil {
		return err
	}

	response.Json(w, &components)
	return nil
}

func DeleteModule(w http.ResponseWriter, r *http.Request) {

}

// Add module if it doesn't exist, update module if it does
// if there is no difference do nothing
func SyncModule(state common.State, w http.ResponseWriter, r *http.Request) error {
	// Parse request
	// TODO(selim): Add former name field to request.
	req := struct {
		Name      string
		Structure json.RawMessage
	}{}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return err
	}

	// Check if module id or not
	id, err := data.GetModulePartIdWithName(state, req.Name, sql.NullInt32{})
	if err != nil {
		return err
	}

	// If module doesn't exist in database; create module and exit
	if !id.Valid {
		initialNode := data.ModulePartNode{
			ValueType: 0,
			ParentId:  sql.NullInt32{},
			Name:      req.Name,
			Value:     req.Structure,
		}

		err = data.InsertModulePartTree(state, &initialNode)
		if err != nil {
			return err
		}
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Successfully created module part!")
		return nil
	}

	// Parts that currently exist in database for the module
	modulePartMapDb, err := data.GetModulePartMap(state, req.Name)
	if err != nil {
		return err
	}

	// Parts that exists in incoming request for the module
	modulePartMapUnity := make(map[string]*data.ModulePartNode)

	// Create processing channel
	// TODO(selim): Is 500 correct for channel size?
	rootNode := data.ModulePartNode{
		ValueType: 0,
		ParentId:  sql.NullInt32{},
		Name:      req.Name,
		Value:     req.Structure,
	}

	c := make(chan data.ModulePartNode, 500)
	c <- rootNode

	// Process nodes in json object tree
	processedObjects := 0
	for n := range c {
		processedObjects++

		// Check if current value is a json object
		children := make(map[string]json.RawMessage)
		err = json.Unmarshal(n.Value, &children)
		if err != nil || children == nil {
			// If there is no more elements to process, close the channel
			if len(c) == 0 {
				close(c)
			}

			// If there is an error during json parsing for root node, exit early with error
			if processedObjects == 1 {
				return err
			}
			continue
		}

		id, _ = data.GetModulePartIdWithName(state, n.Name, n.ParentId)

		// Add all values to process channel
		for child := range children {
			var pi data.ModulePartInfo
			err = json.Unmarshal(children[child], &pi)
			if err != nil {
				return err
			}
			childNode := data.ModulePartNode{
				Name:      child,
				ValueType: pi.ValueType,
				Value:     pi.Children,
				IsArray:   pi.IsArray,
				Tooltip:   pi.Tooltip,
				ParentId:  id,
			}
			c <- childNode
			modulePartMapUnity[data.GetModulePartKey(n.Name, child)] = &childNode
		}
	}

	for k, dbModule := range modulePartMapDb {
		// Skip the root module
		if dbModule.ParentId == 0 {
			continue
		}

		// Part exists in database but not Unity, delete part
		if _, ok := modulePartMapUnity[k]; !ok {
			err = data.DeleteModulePartTree(state, dbModule.Id)
			if err != nil {
				return err
			}
		}
	}

	for k, unityModule := range modulePartMapUnity {
		if dbModule, ok := modulePartMapDb[k]; ok {
			// Part exists in unity and database, update fields
			data.UpdateModulePart(state, dbModule.Id, unityModule)
		} else {
			// Part exists in unity but not database, create part
			err = data.InsertModulePartTree(state, unityModule)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
