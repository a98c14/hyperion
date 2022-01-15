package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/a98c14/hyperion/api/prefab-editor/data"
	"github.com/a98c14/hyperion/common"
	e "github.com/a98c14/hyperion/common/errors"
	"github.com/a98c14/hyperion/common/response"
	"github.com/go-chi/chi/v5"
)

/* Gets the base module_part and all of its children by id.
Can only query base components. Nodes that have null as children
represent leaf nodes */
func GetModuleById(w http.ResponseWriter, r *http.Request) {
	state, err := common.InitState(r)
	if err != nil {
		response.ErrorWhileInitializing(w, err)
		return
	}

	// Load components that have given name from database
	moduleIdString := chi.URLParam(r, "moduleId")
	moduleId, err := strconv.Atoi(moduleIdString)
	if err != nil {
		http.Error(w, "Could not parse id value!", http.StatusBadRequest)
		return
	}

	moduleParts, err := data.GetModuleParts(state, moduleId)
	if err != nil {
		response.InternalError(w, err)
		return
	}

	type modulePartResp struct {
		Id        int               `json:"id"`
		Name      string            `json:"name"`
		ValueType int               `json:"valueType"`
		Children  []*modulePartResp `json:"children"`
	}

	/* Create object tree from module part list. Components are always ordered from root to child*/
	root := modulePartResp{
		Id:        moduleParts[0].Id,
		Name:      moduleParts[0].Name,
		ValueType: moduleParts[0].ValueType,
		Children:  make([]*modulePartResp, 0),
	}
	nodeMap := make(map[int]*modulePartResp)
	nodeMap[root.Id] = &root
	for _, c := range moduleParts[1:] {
		if val, ok := nodeMap[c.ParentId]; ok {
			cr := modulePartResp{
				Id:        c.Id,
				Name:      c.Name,
				ValueType: c.ValueType,
				Children:  nil,
			}
			val.Children = append(val.Children, &cr)
			nodeMap[c.Id] = &cr
		} else {
			http.Error(w, "Could not create object tree from components! Error: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}

	response.Json(w, &root)
}

func ListComponents(w http.ResponseWriter, r *http.Request) {

}

// Returns all the root components.
// Root components are module_part that have no parent. Used to group
// other components and editor can only filter using root components
func GetRootModules(w http.ResponseWriter, r *http.Request) {
	state, err := common.InitState(r)
	if err != nil {
		response.ErrorWhileInitializing(w, err)
		return
	}

	components, err := data.GetRootModuleParts(state)
	if err != nil {
		response.InternalError(w, err)
		return
	}

	response.Json(w, &components)
}

func DeleteModule(w http.ResponseWriter, r *http.Request) {

}

// Add module if it doesn't exist, update module if it does
// if there is no difference do nothing
func SyncModule(w http.ResponseWriter, r *http.Request) {
	state, err := common.InitState(r)
	if err != nil {
		response.ErrorWhileInitializing(w, err)
	}

	// Parse request
	// TODO(selim): Add former name field to request.
	type request struct {
		Name      string
		Structure json.RawMessage
	}
	var req request
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		response.BadRequest(w, e.WrapMsg("SyncModule", response.ParseError, err))
		return
	}

	// Check if module id or not
	id, err := data.GetModulePartIdWithName(state, req.Name, sql.NullInt32{})
	if err != nil {

		response.InternalError(w, e.Wrap("SyncModule", err))
		return
	}

	// If module doesn't exist in database; create module and exit
	if !id.Valid {
		initialNode := data.ModulePartNode{
			ValueType: 0,
			ParentId:  sql.NullInt32{},
			Name:      req.Name,
			Value:     req.Structure,
		}

		data.InsertModulePartTree(state, &initialNode)
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Successfully created module part!")
		return
	}

	// Parts that currently exist in database for the module
	modulePartMapDb, err := data.GetModulePartMap(state, req.Name)
	if err != nil {
		response.InternalError(w, e.Wrap("SyncModule", err))
		return
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
		if err != nil {
			// If there is no more elements to process, close the channel
			if len(c) == 0 {
				close(c)
			}

			// If there is an error during json parsing for root node, exit early with error
			if processedObjects == 1 {
				response.InternalError(w, e.WrapMsg("SyncModule", "Root object doesn't have a valid json", err))
				return
			}
			continue
		}

		id, _ = data.GetModulePartIdWithName(state, n.Name, n.ParentId)

		// Add all values to process channel
		for child := range children {
			var pi data.ModulePartInfo
			err = json.Unmarshal(children[child], &pi)
			if err != nil {
				http.Error(w, "Invalid module part structure!", http.StatusBadRequest)
				return
			}
			childNode := data.ModulePartNode{
				Name:      child,
				ValueType: pi.ValueType,
				Value:     pi.Children,
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
				response.InternalError(w, e.WrapMsg("UpdateModule", "during module part tree deletion", err))
			}
		}
	}

	for k, unityModule := range modulePartMapUnity {
		if _, ok := modulePartMapDb[k]; !ok {
			err = data.InsertModulePartTree(state, unityModule)
			if err != nil {
				response.InternalError(w, e.WrapMsg("UpdateModule", "during module part tree insertion", err))
			}
		}
	}
}
