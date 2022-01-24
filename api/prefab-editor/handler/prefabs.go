package handler

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	"github.com/a98c14/hyperion/api/prefab-editor/data"
	"github.com/a98c14/hyperion/common"
	xerrors "github.com/a98c14/hyperion/common/errors"
	"github.com/a98c14/hyperion/common/json"
	"github.com/a98c14/hyperion/common/response"
	"github.com/go-chi/chi/v5"
)

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

func ListPrefabs(w http.ResponseWriter, r *http.Request) {
	state, err := common.InitState(r)
	if err != nil {
		response.ErrorWhileInitializing(w, err)
		return
	}

	prefabs, err := data.GetRootPrefabs(state)
	if err != nil {
		response.InternalError(w, err)
		return
	}

	response.Json(w, prefabs)
}

func CreatePrefab(w http.ResponseWriter, r *http.Request) {
	state, err := common.InitState(r)
	if err != nil {
		response.ErrorWhileInitializing(w, err)
		return
	}

	// TODO(selim): Should insert list of prefabs
	type model struct {
		Name     string
		ParentId sql.NullInt32
		Modules  []data.PrefabModulePartValue
	}

	var req model
	err = json.Decode(r, &req)
	if err != nil {
		response.BadRequest(w, fmt.Errorf("CreatePrefab:85 | Decode Request -> %w", err))
		return
	}

	ctx := state.Context
	conn := state.Conn

	exists, err := data.DoesNameExist(ctx, conn, req.Name)
	if err != nil {
		response.InternalError(w, err)
		return
	}

	if exists {
		response.BadRequest(w, xerrors.ErrExists)
		return
	}

	// Start transcation. If all components can not be added successfully, don't
	// insert anything
	tx, err := conn.Begin(ctx)
	if err != nil {
		http.Error(w, "Could not start transaction!", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback(ctx)
	prefabId, err := data.InsertPrefab(ctx, conn, req.Name, req.ParentId)
	if err != nil {
		response.InternalError(w, err)
		return
	}

	// TODO(selim): Given version id should be for the `BaseVersion` id
	err = data.InsertPrefabModulePartValues(ctx, conn, int(prefabId.Int32), 1, req.Modules)
	if err != nil {
		response.InternalError(w, err)
		return
	}

}

func UpdatePrefab(id int, w http.ResponseWriter, r *http.Request) {

}
