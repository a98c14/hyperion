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

func GetAssets(state common.State, w http.ResponseWriter, r *http.Request) error {
	assetTypeString := r.URL.Query().Get("assetType")
	if assetTypeString == "" {
		return errors.ErrBadRequest
	}

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

func GetAssetName(state common.State, w http.ResponseWriter, r *http.Request) error {
	assetIdString := chi.URLParam(r, "assetId")
	assetId, err := strconv.Atoi(assetIdString)
	if err != nil {
		return errors.ErrNotFound
	}

	resp := struct {
		Name string `json:"name"`
	}{}
	resp.Name, err = DbGetAssetName(state, int32(assetId))
	if err != nil {
		return err
	}

	response.Json(w, &resp)
	return nil
}

func SyncAssets(state common.State, w http.ResponseWriter, r *http.Request) error {
	req := struct {
		Type   AssetType
		Assets []AssetUnity
	}{}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return err
	}

	resp := make([]AssetDb, len(req.Assets))
	for _, asset := range req.Assets {
		id, err := DbSyncAsset(state, req.Type, &asset)
		if err != nil {
			return err
		}
		resp = append(resp, AssetDb{
			Id:        id,
			UnityGuid: asset.Guid,
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
