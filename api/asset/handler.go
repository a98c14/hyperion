package asset

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/a98c14/hyperion/common"
	"github.com/a98c14/hyperion/common/errors"
	"github.com/a98c14/hyperion/common/response"
	"github.com/go-chi/chi/v5"
)

func GetAssets(w http.ResponseWriter, r *http.Request) error {
	state, err := common.InitState(r)
	if err != nil {
		return err
	}

	assetTypeString := chi.URLParam(r, "assetType")
	assetType, err := strconv.Atoi(assetTypeString)
	if err != nil {
		return errors.ErrNotFound
	}

	assets, err := DbGetAssets(state, AssetType(assetType))
	if err != nil {
		return err
	}

	response.Json(w, &assets)
	return nil
}

func SyncAssets(w http.ResponseWriter, r *http.Request) error {
	state, err := common.InitState(r)
	if err != nil {
		return err
	}

	// Parse request
	type reqModel struct {
		Type   AssetType
		Assets []AssetUnity
	}

	var req reqModel
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return err
	}

	resp := make([]AssetDb, len(req.Assets))
	for _, asset := range req.Assets {
		// TODO(selim): Get asset id if it exists
		var id int32
		id, err := DbSyncAsset(state, req.Type, &asset)
		if err != nil {
			return err
		}
		resp = append(resp, AssetDb{
			Id:        id,
			UnityGuid: "t",
			Name:      asset.Name,
			Type:      req.Type,
		})
	}
	if err != nil {
		return err
	}

	response.Json(w, resp)
	return nil
}
