package handler

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/a98c14/hyperion/api/prefab-editor/data"
	"github.com/a98c14/hyperion/common"
	xerrors "github.com/a98c14/hyperion/common/errors"
	xjson "github.com/a98c14/hyperion/common/json"
	"github.com/a98c14/hyperion/common/response"
	"github.com/go-chi/chi/v5"
)

func CreatePrefabPreset(state common.State, w http.ResponseWriter, r *http.Request) error {
	// TODO(selim):
	// 1. Create the base prefab structure
	// 2. Add transform module to id
	// 3. Insert this prefab to db and get its id
	// 4. Return this prefab as json
	return nil
}

func GetPrefabById(w http.ResponseWriter, r *http.Request) {
	state, err := common.InitState(r)
	if err != nil {
		response.ErrorWhileInitializing(w, err)
		return
	}

	prefabIdStr := chi.URLParam(r, "prefabId")
	prefabId, err := strconv.Atoi(prefabIdStr)
	if err != nil {
		response.BadRequest(w, err)
		return
	}

	// Get version id. Use `BaseVersion` id as default
	var versionId int
	versionIdStr := chi.URLParam(r, "versionId")
	if versionIdStr == "" {
		versionId = 1 // TODO(selim): This should be the base version id.
	} else {
		versionId, err = strconv.Atoi(versionIdStr)
		if err != nil {
			response.BadRequest(w, err)
			return
		}
	}

	prefabs, err := data.GetPrefabById(state, prefabId, versionId)
	if err != nil {
		response.InternalError(w, err)
		return
	}

	response.Json(w, prefabs)
}

func ListPrefabs(state common.State, w http.ResponseWriter, r *http.Request) error {
	prefabs, err := data.GetRootPrefabs(state)
	if err != nil {
		return err
	}

	response.Json(w, prefabs)
	return nil
}

func CreatePrefab(state common.State, w http.ResponseWriter, r *http.Request) error {
	type prefabCreateRequest struct {
		ParentId  sql.NullInt32                `json:"parentId"`
		Name      string                       `json:"name"`
		Transform json.RawMessage              `json:"transform"`
		Renderer  json.RawMessage              `json:"renderer"`
		Colliders json.RawMessage              `json:"colliders"`
		Modules   []data.PrefabModulePartValue `json:"modules"`
		Children  json.RawMessage              `json:"children"`
	}

	var req prefabCreateRequest
	err := xjson.Decode(r, &req)
	if err != nil {
		return err
	}

	ctx := state.Context
	conn := state.Conn
	exists, err := data.DoesNameExist(ctx, conn, req.Name)
	if err != nil {
		return err
	} else if exists {
		return errors.New("given name already exists")
	}

	prefabChannel := make(chan prefabCreateRequest, 500)
	prefabChannel <- req

	// Start transcation. If all components can not be added successfully, don't
	// insert anything
	tx, err := conn.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	// Process nodes in json object tree
	for prefab := range prefabChannel {
		if prefab.Transform == nil || prefab.Colliders == nil || prefab.Renderer == nil {
			return xerrors.ErrBadRequest
		}

		prefabId, err := data.InsertPrefab(ctx, conn, req.Name, prefab.ParentId, prefab.Transform, prefab.Renderer, prefab.Colliders)
		if err != nil {
			return err
		}
		err = data.InsertPrefabModulePartValues(state, int(prefabId.Int32), 1, prefab.Modules)
		if err != nil {
			return err
		}

		// Check if current value is a json object
		children := make([]prefabCreateRequest, 0)
		err = json.Unmarshal(prefab.Children, &children)
		if err != nil || children == nil || len(children) == 0 {
			// If there is no more elements to process, close the channel
			if len(prefabChannel) == 0 {
				close(prefabChannel)
			}
			continue
		}

		// Add all values to process channel
		for _, child := range children {
			child.ParentId = prefabId
			prefabChannel <- child
		}
	}
	err = tx.Commit(state.Context)
	if err != nil {
		return err
	}

	response.Success(w, "Successfully created prefab!")
	return nil
}

func DeletePrefab(state common.State, w http.ResponseWriter, r *http.Request) error {
	prefabIdStr := chi.URLParam(r, "prefabId")
	prefabId, err := strconv.Atoi(prefabIdStr)
	if err != nil {
		return err
	}

	err = data.DeletePrefab(state, prefabId)
	if err != nil {
		return err
	}

	return nil
}

func UpdatePrefab(state common.State, w http.ResponseWriter, r *http.Request) error {
	type prefabUpdateRequest struct {
		Id        int                          `json:"id"`
		ParentId  sql.NullInt32                `json:"parentId"`
		Name      string                       `json:"name"`
		Transform json.RawMessage              `json:"transform"`
		Renderer  json.RawMessage              `json:"renderer"`
		Colliders json.RawMessage              `json:"colliders"`
		Modules   []data.PrefabModulePartValue `json:"modules"`
		Children  json.RawMessage              `json:"children"`
	}

	var req prefabUpdateRequest
	err := xjson.Decode(r, &req)
	if err != nil {
		return err
	}

	prefabChannel := make(chan prefabUpdateRequest, 500)
	prefabChannel <- req

	// Start transcation. If all components can not be added successfully, don't
	// insert anything
	tx, err := state.Conn.Begin(state.Context)
	if err != nil {
		return err
	}
	defer tx.Rollback(state.Context)
	// Process nodes in json object tree
	for prefab := range prefabChannel {
		if prefab.Transform == nil || prefab.Colliders == nil || prefab.Renderer == nil {
			return xerrors.ErrBadRequest
		}

		if prefab.Id == 0 {
			return errors.New("invalid prefab id")
		}

		err := data.UpdatePrefab(state, prefab.Id, prefab.Name, prefab.ParentId, prefab.Transform, prefab.Renderer, prefab.Colliders)
		if err != nil {
			return err
		}
		err = data.UpdatePrefabModulePartValues(state, prefab.Id, 1, prefab.Modules)
		if err != nil {
			return err
		}

		// Check if current value is a json object
		children := make([]prefabUpdateRequest, 0)
		err = json.Unmarshal(prefab.Children, &children)
		if err != nil || children == nil || len(children) == 0 {
			// If there is no more elements to process, close the channel
			if len(prefabChannel) == 0 {
				close(prefabChannel)
			}
			continue
		}

		// Add all values to process channel
		for _, child := range children {
			child.ParentId = prefab.ParentId
			prefabChannel <- child
		}
	}
	err = tx.Commit(state.Context)
	if err != nil {
		return err
	}

	response.Success(w, "Successfully updated prefab!")
	return nil
}
